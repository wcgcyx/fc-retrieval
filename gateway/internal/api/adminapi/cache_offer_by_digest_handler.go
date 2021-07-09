/*
Package adminapi contains the API code for the admin client - gateway communication.
*/
package adminapi

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
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// CacheOfferByDigestHandler handles cache offer request
func CacheOfferByDigestHandler(data []byte) (byte, []byte, error) {
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Decode payload
	digest, cidStr, err := fcradminmsg.DecodeCacheOfferByDigestRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	cid, err := cid.NewContentID(cidStr)
	if err != nil {
		err = fmt.Errorf("Error in decoding cid: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	offer := c.OfferMgr.GetOfferByDigest(digest)
	if offer == nil {
		err = fmt.Errorf("Cannot find offer with digest: %v", digest)
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	suboffer, err := offer.GenerateSubCIDOffer(cid)
	if err != nil {
		err = fmt.Errorf("Error in generating sub cid offer: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Get provider information
	pvdInfo := c.PeerMgr.GetPVDInfo(suboffer.GetProviderID())
	if pvdInfo == nil {
		// Not found, try sync once
		pvdInfo = c.PeerMgr.SyncPVD(suboffer.GetProviderID())
		if pvdInfo == nil {
			err = fmt.Errorf("Cannot find provider %v that supplied the offer", suboffer.GetProviderID())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
	}
	// Do caching
	_, err = c.P2PServer.Request(pvdInfo.NetworkAddr, fcrmessages.DataRetrievalRequestType, pvdInfo.NodeID, suboffer)
	if err != nil {
		err = fmt.Errorf("Error in data retrieval: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Succeed
	ack := fcradminmsg.EncodeACK(true, "Succeed.")
	return fcradminmsg.ACKType, ack, nil
}
