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
	"fmt"
	"math/rand"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// OfferPublishRequester sends an offer publish request.
func OfferPublishRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
	logging.Debug("Request offer publish")
	// Get parameters
	if len(args) != 2 {
		err := fmt.Errorf("Wrong arguments, expect length 2, got length %v", len(args))
		logging.Error(err.Error())
		return nil, err
	}
	targetID, ok := args[0].(string)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a target ID in string")
		logging.Error(err.Error())
		return nil, err
	}
	offer, ok := args[1].(*cidoffer.CIDOffer)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a offer in *cidoffer.CIDOffer")
		logging.Error(err.Error())
		return nil, err
	}

	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Generate random nonce
	nonce := uint64(rand.Int63())

	request, err := fcrmessages.EncodeOfferPublishRequest(nonce, c.NodeID, offer)
	if err != nil {
		err = fmt.Errorf("Internal error in encoding response: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	err = request.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		// Error in signing
		return nil, err
	}

	// Write request
	err = writer.Write(request, c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in sending request to %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in receiving response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Verify the response
	gwInfo := c.PeerMgr.GetGWInfo(targetID)
	if gwInfo == nil {
		// Not found, try sync once
		gwInfo = c.PeerMgr.SyncGW(targetID)
		if gwInfo == nil {
			err = fmt.Errorf("Error in obtaining information for gateway %v", targetID)
			logging.Error(err.Error())
			return nil, err
		}
	}
	if response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
		// Try update
		gwInfo = c.PeerMgr.SyncGW(targetID)
		if gwInfo == nil || response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
			err = fmt.Errorf("Error in verifying response from %v: %v", targetID, err.Error())
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Check response
	if !response.ACK() {
		err = fmt.Errorf("Reponse contains an error: %v", response.Error())
		logging.Error(err.Error())
		return nil, err
	} else {
		logging.Info("Successfully published offer to gateway %v", targetID)
	}

	return response, nil
}
