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
	AddHandler(msgType byte, handler func(reader FCRServerReader, writer FCRServerWriter, request *fcrmessages.FCRMessage) error) FCRServer

	// AddRequester adds a requester to the server, which is used to send a request for a given message type.
	AddRequester(msgType byte, requester func(reader FCRServerReader, writer FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)) FCRServer

	// Request uses a requester corresponding to the given message type to send a request to given multiaddr.
	Request(multiaddrStr string, msgType byte, args ...interface{}) (*fcrmessages.FCRMessage, error)
}

// FCRServerReader is a reader for reading message.
type FCRServerReader interface {
	// Read reads a message for a given timeout.
	// It returns the message, and error.
	Read(timeout time.Duration) (*fcrmessages.FCRMessage, error)
}

// FCRServerWriter is a reader for writer message.
type FCRServerWriter interface {
	// Write writes a message for a given timeout.
	// It returns error.
	Write(msg *fcrmessages.FCRMessage, timeout time.Duration) error
}
