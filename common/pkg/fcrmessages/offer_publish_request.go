/*
Package fcrmessages - stores all the p2p messages.
*/
package fcrmessages

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
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
)

// offerPublishRequestJson represents the request to publish an offer.
type offerPublishRequestJson struct {
	NodeID string `json:"node_id"`
	Offer  string `json:"offer"`
}

// EncodeOfferPublishRequest is used to get the FCRMessage of offerPublishRequestJson.
func EncodeOfferPublishRequest(
	nonce uint64,
	nodeID string,
	offer *cidoffer.CIDOffer,
) (*FCRReqMsg, error) {
	data, err := offer.ToBytes()
	if err != nil {
		return nil, err
	}
	offerStr := hex.EncodeToString(data)
	body, err := json.Marshal(offerPublishRequestJson{
		NodeID: nodeID,
		Offer:  offerStr,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRReqMsg(OfferPublishRequestType, nonce, body), nil
}

// DecodeOfferPublishRequest is used to get the fields from FCRMessage of offerPublishRequestJson.
// It returns the nonce, nodeID and the offer.
func DecodeOfferPublishRequest(fcrMsg *FCRReqMsg) (
	uint64,
	string,
	*cidoffer.CIDOffer,
	error,
) {
	if fcrMsg.Type() != OfferPublishRequestType {
		return 0, "", nil, fmt.Errorf("Message type mismatch, expect %v, got %v", OfferPublishRequestType, fcrMsg.Type())
	}
	msg := offerPublishRequestJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", nil, err
	}
	data, err := hex.DecodeString(msg.Offer)
	if err != nil {
		return 0, "", nil, err
	}
	offer := cidoffer.CIDOffer{}
	err = offer.FromBytes(data)
	if err != nil {
		return 0, "", nil, err
	}
	return fcrMsg.Nonce(), msg.NodeID, &offer, nil
}
