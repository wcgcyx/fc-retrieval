/*
Package p2papi contains the API code for the p2p communication.
*/
package p2papi

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
	"errors"
	"math/rand"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// OfferPublishRequester sends an offer publish request.
func OfferPublishRequester(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 1 {
		return nil, errors.New("wrong arguments")
	}
	offer, ok := args[0].(*cidoffer.CIDOffer)
	if !ok {
		return nil, errors.New("wrong arguments")
	}

	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Generate random nonce
	nonce := rand.Int63()

	request, err := fcrmessages.EncodeOfferPublishRequest(c.NodeID, nonce, offer)
	if err != nil {
		// Error in encoding
		return nil, err
	}
	err = request.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		// Error in signing
		return nil, err
	}

	err = writer.Write(request, c.Settings.TCPInactivityTimeout)
	if err != nil {
		// Error in writing
		return nil, err
	}

	return reader.Read(c.Settings.TCPInactivityTimeout)
}
