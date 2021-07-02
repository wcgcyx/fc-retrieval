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
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
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

	s.host = h
	s.start = true
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

func (s *FCRServerImplV1) Request(multiaddrStr string, msgType byte, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	if !s.start {
		return nil, errors.New("Server not started")
	}
	requester := s.requesters[msgType]
	if requester == nil {
		return nil, errors.New("No available requester found for given type")
	}
	// Get multiaddr
	maddr, err := multiaddr.NewMultiaddr(multiaddrStr)
	if err != nil {
		return nil, err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}
	// Store peer
	s.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// Start new stream
	conn, err := s.host.NewStream(context.Background(), info.ID, "/fc-retrieval/0.0.1")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	logging.Info("Established connection to %v", info.ID)

	writer := &FCRServerWriterImplV1{conn: conn}
	reader := &FCRServerReaderImplV1{conn: conn}
	response, err := requester(reader, writer, args...)
	return response, nil
}

func (s *FCRServerImplV1) handleIncomingConnection(conn network.Stream) {
	// New incoming connection
	logging.Info("P2P server has incoming connection from :%s", conn.ID())

	// Close connection on exit.
	defer conn.Close()

	message, err := readFCRMessage(conn, time.Second)
	if err != nil {
		// Error in tcp communication, drop the connection.
		logging.Error("P2P Server has error reading message from %s: %s - Connection dropped", conn.ID(), err.Error())
		return
	}
	handler := s.handlers[message.GetMessageType()]
	if handler != nil {
		// Call handler to handle the request
		writer := &FCRServerWriterImplV1{conn: conn}
		reader := &FCRServerReaderImplV1{conn: conn}
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

// readFCRMessage read a fcr message from a given connection.
func readFCRMessage(conn network.Stream, timeout time.Duration) (*fcrmessages.FCRMessage, error) {
	// Initialise a reader
	reader := bufio.NewReader(conn)
	// Read the length
	length := make([]byte, 4)
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	_, err := io.ReadFull(reader, length)
	if err != nil {
		return nil, err
	}
	// Read the data
	data := make([]byte, int(binary.BigEndian.Uint32(length)))
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return fcrmessages.FromBytes(data)
}

// sendFCRMessage sends a fcr message to a given connection.
func sendFCRMessage(conn network.Stream, fcrMsg *fcrmessages.FCRMessage, timeout time.Duration) error {
	// Get data
	data, err := fcrMsg.ToBytes()
	if err != nil {
		return err
	}
	// Initialise a writer
	writer := bufio.NewWriter(conn)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	_, err = writer.Write(append(length, data...))
	if err != nil {
		return err
	}
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	return writer.Flush()
}

// FCRServerReaderImplV1 implements FCRServerReader.
type FCRServerReaderImplV1 struct {
	conn network.Stream
}

func (r *FCRServerReaderImplV1) Read(timeout time.Duration) (*fcrmessages.FCRMessage, error) {
	return readFCRMessage(r.conn, timeout)
}

// FCRServerWriterImplV1 implements FCRServerWriter.
type FCRServerWriterImplV1 struct {
	conn network.Stream
}

func (w *FCRServerWriterImplV1) Write(msg *fcrmessages.FCRMessage, timeout time.Duration) error {
	return sendFCRMessage(w.conn, msg, timeout)
}

// IsTimeoutError checks if the given error is a timeout error
func IsTimeoutError(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}

// GetMultiAddr returns the supposed multiaddr string from given private key, ip address and port.
func GetMultiAddr(prvKeyStr string, ip string, port uint) (string, error) {
	prvKeyBytes, err := hex.DecodeString(prvKeyStr)
	if err != nil {
		return "", err
	}
	prvKey, err := crypto.UnmarshalRsaPrivateKey(prvKeyBytes)
	if err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(prvKey.GetPublic())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/ip4/%v/tcp/%v/p2p/%v", ip, port, pid), nil
}
