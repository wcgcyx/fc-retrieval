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

import "github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"

// FCRServerImplV1 implements FCRServer, it is built on top of libp2p
type FCRServerImplV1 struct {
}

func NewFCRServerImplV1() FCRServer {
	return nil
}

func (s *FCRServerImplV1) Start() error {
	return nil
}

func (s *FCRServerImplV1) Shutdown() {
}

func (s *FCRServerImplV1) AddHandler(msgType int32, handler func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error) FCRServer {
	return s
}

func (s *FCRServerImplV1) AddRequester(msgType int32, requester func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)) FCRServer {
	return s
}
