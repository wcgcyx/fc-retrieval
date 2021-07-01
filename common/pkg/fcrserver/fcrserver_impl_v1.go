/*
Package fcrserver - provides an interface to do networking.
*/
package fcrserver

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

const (
	Lowwater    = 100
	Highwater   = 400
	GracePeriod = time.Hour
)

// FCRServerImplV1 implements FCRServer, it is built on top of libp2p.
type FCRServerImplV1 struct {
	prvKeyStr string
	port      uint
	start     bool
	timeout   time.Duration

	handlers   map[byte]func(reader FCRServerReader, writer FCRServerWriter, request *fcrmessages.FCRMessage) error
	requesters map[byte]func(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)

	shutdown chan bool
	host     host.Host
	cancel   context.CancelFunc
}

func NewFCRServerImplV1(prvKeyStr string, port uint, timeout time.Duration) FCRServer {
	return &FCRServerImplV1{
		prvKeyStr:  prvKeyStr,
		port:       port,
		start:      false,
		timeout:    timeout,
		handlers:   make(map[byte]func(reader FCRServerReader, writer FCRServerWriter, request *fcrmessages.FCRMessage) error),
		requesters: make(map[byte]func(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)),
		shutdown:   make(chan bool),
	}
}

func (s *FCRServerImplV1) Start() error {
	// Start server
	if s.start {
		return errors.New("Server already started")
	}
	prvKeyBytes, err := hex.DecodeString(s.prvKeyStr)
	if err != nil {
		return err
	}
	prvKey, err := crypto.UnmarshalRsaPrivateKey(prvKeyBytes)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", s.port))
	if err != nil {
		return err
	}

	h, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
		libp2p.ConnectionManager(connmgr.NewConnManager(
			Lowwater,
			Highwater,
			GracePeriod,
		)),
	)
	if err != nil {
		return err
	}
	h.SetStreamHandler("/fc-retrieval/0.0.1", s.handleIncomingConnection)
	logging.Info("P2P Server starts listening on /ip4/127.0.0.1/tcp/%v/p2p/%s", s.port, h.ID())

	go func() {
		defer s.cancel()
		select {
		case <-s.shutdown:
			logging.Info("P2P Server shutdown.")
			s.shutdown <- true
		}
	}()

	return nil
}

func (s *FCRServerImplV1) Shutdown() {
	if !s.start {
		return
	}
	s.start = false
	s.shutdown <- true
	<-s.shutdown
}

func (s *FCRServerImplV1) AddHandler(msgType byte, handler func(reader FCRServerReader, writer FCRServerWriter, request *fcrmessages.FCRMessage) error) FCRServer {
	if s.start {
		return s
	}
	s.handlers[msgType] = handler
	return s
}

func (s *FCRServerImplV1) AddRequester(msgType byte, requester func(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)) FCRServer {
	if s.start {
		return s
	}
	s.requesters[msgType] = requester
	return s
}

func (s *FCRServerImplV1) Request(multiaddrStr string, msgType byte, args ...interface{}) (*fcrmessages.FCRMessage, bool, error) {
	if !s.start {
		return nil, false, errors.New("Server not started")
	}
	requester := s.requesters[msgType]
	if requester == nil {
		return nil, false, errors.New("No available requester found for given type")
	}
	// Get multiaddr
	maddr, err := multiaddr.NewMultiaddr(multiaddrStr)
	if err != nil {
		return nil, false, err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, false, err
	}
	// Store peer
	s.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// Start new stream
	conn, err := s.host.NewStream(context.Background(), info.ID, "/fc-retrieval/0.0.1")
	if err != nil {
		return nil, false, err
	}
	logging.Info("Established connection to %v", info.ID)

	writer := &FCRServerWriterImplV1{conn2: conn}
	reader := &FCRServerReaderImplV1{conn2: conn}
	response, err := requester(reader, writer, args...)
	return response, isTimeoutError(err), nil
}

func (s *FCRServerImplV1) handleIncomingConnection(conn network.Stream) {
	// New incoming connection
	logging.Info("P2P server has incoming connection from :%s", conn.ID())

	// Close connection on exit.
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	// Loop until error occurs and connection is dropped.
	for {
		message, err := readFCRMessage(conn, nil, s.timeout)
		if err != nil && !isTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error("P2P Server has error reading message from %s: %s - Connection dropped", conn.ID(), err.Error())
			return
		}
		// TODO: discard a connection if it doesn’t give a valid response for a really long time
		if err != nil && isTimeoutError(err) {
			continue
		}
		handler := s.handlers[message.GetMessageType()]
		if handler != nil {
			// Call handler to handle the request
			writer := &FCRServerWriterImplV1{conn1: conn}
			reader := &FCRServerReaderImplV1{conn1: conn}
			err = handler(reader, writer, message)
			if err != nil {
				// Error that couldn't ignore, drop the connection.
				logging.Error("P2P Server has error handling message from %s: %s - Connection dropped", conn.ID(), err.Error())
				return
			}
		} else {
			// Message is invalid, drop the connection.
			logging.Error("P2P Server received unsupported message type %v from %s - Connection dropped", message.GetMessageType(), conn.ID())
			return
		}
	}
}

// isTimeoutError checks if the given error is a timeout error
func isTimeoutError(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}

// readFCRMessage read a fcr message from a given connection.
func readFCRMessage(conn1 network.Stream, conn2 inet.Stream, timeout time.Duration) (*fcrmessages.FCRMessage, error) {
	// Initialise a reader
	var reader *bufio.Reader
	if conn1 != nil {
		reader = bufio.NewReader(conn1)
	} else {
		reader = bufio.NewReader(conn2)
	}
	// Read the length
	length := make([]byte, 4)
	// Set timeout
	if conn1 != nil {
		if err := conn1.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	} else {
		if err := conn2.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	}
	_, err := io.ReadFull(reader, length)
	if err != nil {
		return nil, err
	}
	// Read the data
	data := make([]byte, int(binary.BigEndian.Uint32(length)))
	// Set timeout
	if conn1 != nil {
		if err := conn1.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	} else {
		if err := conn2.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	}
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return fcrmessages.FromBytes(data)
}

// sendFCRMessage sends a fcr message to a given connection.
func sendFCRMessage(conn1 network.Stream, conn2 inet.Stream, fcrMsg *fcrmessages.FCRMessage, timeout time.Duration) error {
	// Get data
	data, err := fcrMsg.ToBytes()
	if err != nil {
		return err
	}
	// Initialise a writer
	var writer *bufio.Writer
	if conn1 != nil {
		writer = bufio.NewWriter(conn1)
	} else {
		writer = bufio.NewWriter(conn2)
	}
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))
	// Set timeout
	if conn1 != nil {
		if err := conn1.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	} else {
		if err := conn2.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	}
	_, err = writer.Write(append(length, data...))
	if err != nil {
		return err
	}
	// Set timeout
	if conn1 != nil {
		if err := conn1.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	} else {
		if err := conn2.SetDeadline(time.Now().Add(timeout)); err != nil {
			panic(err)
		}
	}
	return writer.Flush()
}

// FCRServerReaderImplV1 implements FCRServerReader.
type FCRServerReaderImplV1 struct {
	conn1 network.Stream
	conn2 inet.Stream
}

func (r *FCRServerReaderImplV1) Read(timeout time.Duration) (*fcrmessages.FCRMessage, bool, error) {
	msg, err := readFCRMessage(r.conn1, r.conn2, timeout)
	return msg, isTimeoutError(err), err
}

// FCRServerWriterImplV1 implements FCRServerWriter.
type FCRServerWriterImplV1 struct {
	conn1 network.Stream
	conn2 inet.Stream
}

func (w *FCRServerWriterImplV1) Write(msg *fcrmessages.FCRMessage, timeout time.Duration) (bool, error) {
	err := sendFCRMessage(w.conn1, w.conn2, msg, timeout)
	return isTimeoutError(err), err
}
