/*
Package fcrp2pserver - provides an interface to interact with libp2p.

FCRP2PServer is a wrapper over libp2p.
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

import "github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"

// FCRP2PServer represents a server handling p2p connection.
type FCRP2PServer interface {
	// Start starts the server.
	Start()

	// Shutdown stops the server.
	Shutdown()

	// AddHandler adds a handler to the server, which handles a given message type.
	AddHandler(msgType int32, handler func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error) *FCRP2PServer

	// AddRequester adds a requester to the server, which is used to send a request for a given message type.
	AddRequester(msgType int32, requester func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)) *FCRP2PServer
}
