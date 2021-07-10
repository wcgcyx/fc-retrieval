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
	privKeyStr string
	port       uint
	start      bool
	timeout    time.Duration

	handlers   map[byte]func(reader FCRServerRequestReader, writer FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error
	requesters map[byte]func(reader FCRServerResponseReader, writer FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error)

	shutdown chan bool
	host     host.Host
	cancel   context.CancelFunc
}

func NewFCRServerImplV1(privKeyStr string, port uint, timeout time.Duration) FCRServer {
	return &FCRServerImplV1{
		privKeyStr: privKeyStr,
		port:       port,
		start:      false,
		timeout:    timeout,
		handlers:   make(map[byte]func(reader FCRServerRequestReader, writer FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error),
		requesters: make(map[byte]func(reader FCRServerResponseReader, writer FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error)),
		shutdown:   make(chan bool),
	}
}

func (s *FCRServerImplV1) Start() error {
	// Start server
	if s.start {
		return errors.New("Server already started")
	}
	privKeyBytes, err := hex.DecodeString(s.privKeyStr)
	if err != nil {
		return err
	}
	privKey, err := crypto.UnmarshalRsaPrivateKey(privKeyBytes)
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
		libp2p.Identity(privKey),
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

func (s *FCRServerImplV1) AddHandler(msgType byte, handler func(reader FCRServerRequestReader, writer FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error) FCRServer {
	if s.start {
		return s
	}
	s.handlers[msgType] = handler
	return s
}

func (s *FCRServerImplV1) AddRequester(msgType byte, requester func(reader FCRServerResponseReader, writer FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error)) FCRServer {
	if s.start {
		return s
	}
	s.requesters[msgType] = requester
	return s
}

func (s *FCRServerImplV1) Request(multiaddrStr string, msgType byte, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
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

	writer := &FCRServerRequestWriterImplV1{conn: conn}
	reader := &FCRServerResponseReaderImplV1{conn: conn}
	return requester(reader, writer, args...)
}

func (s *FCRServerImplV1) handleIncomingConnection(conn network.Stream) {
	// New incoming connection
	logging.Info("P2P server has incoming connection from :%s", conn.ID())

	// Close connection on exit.
	defer conn.Close()

	reader := &FCRServerRequestReaderImplV1{conn}
	request, err := reader.Read(s.timeout)
	if err != nil {
		// Error in tcp communication, drop the connection.
		logging.Error("P2P Server has error reading message from %s: %s - Connection dropped", conn.ID(), err.Error())
		return
	}
	handler := s.handlers[request.Type()]
	if handler != nil {
		// Call handler to handle the request
		writer := &FCRServerResponseWriterImplV1{conn: conn}
		err = handler(reader, writer, request)
		if err != nil {
			// Error that couldn't ignore, drop the connection.
			logging.Error("P2P Server has error handling message from %s: %s - Connection dropped", conn.ID(), err.Error())
			return
		}
	} else {
		// Message is invalid, drop the connection.
		logging.Error("P2P Server received unsupported message type %v from %s - Connection dropped", request.Type(), conn.ID())
		return
	}
}

// read read a message bytes from a given connection.
func read(conn network.Stream, timeout time.Duration) ([]byte, error) {
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
	return data, err
}

// write writes a message bytes array to a given connection.
func write(conn network.Stream, data []byte, timeout time.Duration) error {
	// Get data
	// Initialise a writer
	writer := bufio.NewWriter(conn)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	_, err := writer.Write(append(length, data...))
	if err != nil {
		return err
	}
	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		panic(err)
	}
	return writer.Flush()
}

// FCRServerRequestReaderImplV1 implements FCRServerRequestReader.
type FCRServerRequestReaderImplV1 struct {
	conn network.Stream
}

func (r *FCRServerRequestReaderImplV1) Read(timeout time.Duration) (*fcrmessages.FCRReqMsg, error) {
	res := &fcrmessages.FCRReqMsg{}
	data, err := read(r.conn, timeout)
	if err == nil {
		err = res.FromBytes(data)
	}
	return res, err
}

// FCRServerResponseReaderImplV1 implements FCRServerResponseReader.
type FCRServerResponseReaderImplV1 struct {
	conn network.Stream
}

func (r *FCRServerResponseReaderImplV1) Read(timeout time.Duration) (*fcrmessages.FCRACKMsg, error) {
	res := &fcrmessages.FCRACKMsg{}
	data, err := read(r.conn, timeout)
	if err == nil {
		err = res.FromBytes(data)
	}
	return res, err
}

// FCRServerRequestWriterImplV1 implements FCRServerRequestWriter.
type FCRServerRequestWriterImplV1 struct {
	conn network.Stream
}

func (w *FCRServerRequestWriterImplV1) Write(msg *fcrmessages.FCRReqMsg, privKey string, keyVer byte, timeout time.Duration) error {
	err := msg.Sign(privKey, keyVer)
	if err == nil {
		data, err := msg.ToBytes()
		if err == nil {
			err = write(w.conn, data, timeout)
		}
	}
	return err
}

// FCRServerResponseWriterImplV1 implements FCRServerResponseWriter.
type FCRServerResponseWriterImplV1 struct {
	conn network.Stream
}

func (w *FCRServerResponseWriterImplV1) Write(msg *fcrmessages.FCRACKMsg, privKey string, keyVer byte, timeout time.Duration) error {
	err := msg.Sign(privKey, keyVer)
	if err == nil {
		data, err := msg.ToBytes()
		if err == nil {
			err = write(w.conn, data, timeout)
		}
	}
	return err
}

// IsTimeoutError checks if the given error is a timeout error
func IsTimeoutError(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}

// GetMultiAddr returns the supposed multiaddr string from given private key, ip address and port.
func GetMultiAddr(privKeyStr string, ip string, port uint) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyStr)
	if err != nil {
		return "", err
	}
	privKey, err := crypto.UnmarshalRsaPrivateKey(privKeyBytes)
	if err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(privKey.GetPublic())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/ip4/%v/tcp/%v/p2p/%v", ip, port, pid), nil
}
