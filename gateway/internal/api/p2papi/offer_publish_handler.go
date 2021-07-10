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

// OfferPublishHandler handles offer publication.
func OfferPublishHandler(reader fcrserver.FCRServerRequestReader, writer fcrserver.FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error {
	logging.Debug("Handle offer publish")
	// Get core response
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Message decoding
	nonce, senderID, offer, err := fcrmessages.DecodeOfferPublishRequest(request)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Verify the signature
	pvdInfo := c.PeerMgr.GetPVDInfo(senderID)
	if pvdInfo == nil {
		// Not found, try sync once
		pvdInfo = c.PeerMgr.SyncPVD(senderID)
		if pvdInfo == nil {
			err = fmt.Errorf("Error in obtaining information for provider %v", senderID)
			logging.Error(err.Error())
			return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
		}
	}
	if request.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) != nil {
		// Try update
		pvdInfo = c.PeerMgr.SyncGW(senderID)
		if pvdInfo == nil || request.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) != nil {
			err = fmt.Errorf("Error in verifying request from provider %v: %v", senderID, err.Error())
			logging.Error(err.Error())
			return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
		}
	}

	// Check offer signature
	if offer.Verify(pvdInfo.OfferSigningKey) != nil {
		err = fmt.Errorf("Received offer fails to verify against signature of provider %v", senderID)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Offer verified, add to storage
	// TODO: if c.StoreFullOffer is set to false, only store offer if it contains an cid that is within the range
	c.OfferMgr.AddOffer(offer)
	return writer.Write(fcrmessages.CreateFCRACKMsg(nonce, []byte{0}), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
}
