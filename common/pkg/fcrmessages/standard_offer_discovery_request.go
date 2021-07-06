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
	"encoding/json"
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
)

// standardOfferDiscoveryRequestJson represents the request to ask for offers.
type standardOfferDiscoveryRequestJson struct {
	Client            bool   `json:"client"`
	NodeID            string `json:"node_id"`
	PieceCID          string `json:"piece_cid"`
	MaxOfferRequested int64  `json:"max_offer_requested"`
	AccountAddr       string `json:"account_addr"`
	Voucher           string `json:"voucher"`
}

// EncodeStandardOfferDiscoveryRequest is used to get the FCRMessage of standardOfferDiscoveryRequestJson.
func EncodeStandardOfferDiscoveryRequest(
	nonce uint64,
	client bool,
	NodeID string,
	pieceCID *cid.ContentID,
	maxOfferRequested int64,
	accountAddr string,
	voucher string,
) (*FCRReqMsg, error) {
	body, err := json.Marshal(standardOfferDiscoveryRequestJson{
		Client:            client,
		NodeID:            NodeID,
		PieceCID:          pieceCID.ToString(),
		MaxOfferRequested: maxOfferRequested,
		AccountAddr:       accountAddr,
		Voucher:           voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRReqMsg(StandardOfferDiscoveryRequestType, nonce, body), nil
}

// DecodeStandardOfferDiscoveryRequest is used to get the fields from FCRMessage of standardOfferDiscoveryRequestJson.
// It returns the nonce, nodeID, pieceCID, maxOfferRequested, account address and voucher.
func DecodeStandardOfferDiscoveryRequest(fcrMsg *FCRReqMsg) (
	uint64,
	bool,
	string,
	*cid.ContentID,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.Type() != StandardOfferDiscoveryRequestType {
		return 0, false, "", nil, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", StandardOfferDiscoveryRequestType, fcrMsg.Type())
	}
	msg := standardOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, false, "", nil, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return 0, false, "", nil, 0, "", "", err
	}
	return fcrMsg.Nonce(), msg.Client, msg.NodeID, pieceCID, msg.MaxOfferRequested, msg.AccountAddr, msg.Voucher, nil
}
