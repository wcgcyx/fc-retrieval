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

import (
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
)

// FCRServerReader is a reader for reading message.
type FCRServerReader interface {
	// Read reads a message for a given timeout.
	// It returns the message, a boolean indicates whether a timeout occurs and error.
	Read(timeout time.Duration) (*fcrmessages.FCRMessage, bool, error)
}
