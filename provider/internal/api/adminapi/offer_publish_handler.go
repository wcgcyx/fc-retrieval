/*
Package adminapi contains the API code for the admin client - provider communication.
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
	"os"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// OfferPublishHandler handles offer publish.
func OfferPublishHandler(data []byte) (byte, []byte, error) {
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	files, price, expiry, qos, err := fcradminmsg.DecodePublishOfferRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// TODO: Replace files to be the real cid.
	cids := make([]cid.ContentID, 0)
	for _, file := range files {
		reader, err := os.Open(filepath.Join(c.Settings.RetrievalDir, file))
		if err != nil {
			err = fmt.Errorf("Fail to open file %v: %v", file, err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		cid, err := cid.NewContentIDFromFile(reader)
		if err != nil {
			err = fmt.Errorf("Invalid CID: %v", err.Error())
			ack := fcradminmsg.EncodeACK(false, err.Error())
			return fcradminmsg.ACKType, ack, err
		}
		c.OfferMgr.AddCIDTag(cid, file)
		cids = append(cids, *cid)
	}

	// Create offer
	offer, err := cidoffer.NewCIDOffer(c.NodeID, cids, price, expiry, qos)
	if err != nil {
		err = fmt.Errorf("Error creating offer: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Sign offer
	err = offer.Sign(c.OfferSigningKey)
	if err != nil {
		err = fmt.Errorf("Error signing offer: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Send offer
	// TODO, concurrency and memory
	gws := c.PeerMgr.ListGWS()
	for _, gw := range gws {
		c.P2PServer.Request(gw.NetworkAddr, fcrmessages.OfferPublishRequestType, offer)
	}

	// Add offer
	c.OfferMgr.AddOffer(offer)

	// Succeed
	ack := fcradminmsg.EncodeACK(true, "Succeed")
	return fcradminmsg.ACKType, ack, nil
}
