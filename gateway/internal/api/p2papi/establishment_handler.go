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

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// EstablishmentHandler handles dht offer establishment.
func EstablishmentHandler(reader fcrserver.FCRServerRequestReader, writer fcrserver.FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error {
	logging.Debug("Handle establishment")
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Message decoding
	nonce, senderID, challenge, err := fcrmessages.DecodeEstablishmentRequest(request)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Verify signature
	if request.VerifyByID(senderID) != nil {
		// Verify by signing key
		gwInfo := c.PeerMgr.GetGWInfo(senderID)
		if gwInfo == nil {
			// Not found, try sync once
			gwInfo = c.PeerMgr.SyncGW(senderID)
			if gwInfo == nil {
				err = fmt.Errorf("Error in obtaining information for gateway %v", senderID)
				logging.Error(err.Error())
				return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
			}
		}
		if request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
			// Try update
			gwInfo = c.PeerMgr.SyncGW(senderID)
			if gwInfo == nil || request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
				err = fmt.Errorf("Error in verifying request from gateway %v: %v", senderID, err.Error())
				logging.Error(err.Error())
				return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
			}
		}
	}

	// Respond
	response, err := fcrmessages.EncodeEstablishmentResponse(nonce, challenge)
	if err != nil {
		err = fmt.Errorf("Internal error in encoding response: %v", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	return writer.Write(response, c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
}
