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
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
)

// FCRServer represents a server handling p2p connection.
type FCRServer interface {
	// Start starts the server.
	Start() error

	// Shutdown stops the server.
	Shutdown()

	// AddHandler adds a handler to the server, which handles a given message type.
	AddHandler(msgType byte, handler func(reader FCRServerRequestReader, writer FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error) FCRServer

	// AddRequester adds a requester to the server, which is used to send a request for a given message type.
	AddRequester(msgType byte, requester func(reader FCRServerResponseReader, writer FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error)) FCRServer

	// Request uses a requester corresponding to the given message type to send a request to given multiaddr.
	Request(multiaddrStr string, msgType byte, args ...interface{}) (*fcrmessages.FCRACKMsg, error)
}

// FCRServerRequestReader is a reader for reading message.
type FCRServerRequestReader interface {
	// Read reads a message for a given timeout.
	// It returns the message, and error.
	Read(timeout time.Duration) (*fcrmessages.FCRReqMsg, error)
}

// FCRServerResponseReader is a reader for reading message.
type FCRServerResponseReader interface {
	// Read reads a message for a given timeout.
	// It returns the message, and error.
	Read(timeout time.Duration) (*fcrmessages.FCRACKMsg, error)
}

// FCRServerRequesterWriter is a writer for writer request.
type FCRServerRequestWriter interface {
	// Write writes a message for a given timeout.
	// It returns error.
	Write(msg *fcrmessages.FCRReqMsg, privKey string, keyVer byte, timeout time.Duration) error
}

// FCRServerResponseWriter is a writer for writer response.
type FCRServerResponseWriter interface {
	// Write writes a message for a given timeout.
	// It returns error.
	Write(msg *fcrmessages.FCRACKMsg, privKey string, keyVer byte, timeout time.Duration) error
}
