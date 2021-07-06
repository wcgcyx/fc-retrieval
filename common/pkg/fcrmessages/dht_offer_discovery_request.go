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

// dhtOfferDiscoveryRequestJson represents the request to ask for offers in DHT.
type dhtOfferDiscoveryRequestJson struct {
	NodeID                  string `json:"node_id"`
	PieceCID                string `json:"piece_cid"`
	NumDHT                  int64  `json:"num_dht"`
	MaxOfferRequestedPerDHT int64  `json:"max_offer_requested_per_dht"`
	AccountAddr             string `json:"account_addr"`
	Voucher                 string `json:"voucher"`
}

// EncodeDHTOfferDiscoveryRequest is used to get the FCRMessage of dhtOfferDiscoveryRequestJson
func EncodeDHTOfferDiscoveryRequest(
	nonce uint64,
	NodeID string,
	pieceCID *cid.ContentID,
	numDHT int64,
	maxOfferRequestedPerDHT int64,
	accountAddr string,
	voucher string,
) (*FCRReqMsg, error) {
	body, err := json.Marshal(dhtOfferDiscoveryRequestJson{
		NodeID:                  NodeID,
		PieceCID:                pieceCID.ToString(),
		NumDHT:                  numDHT,
		MaxOfferRequestedPerDHT: maxOfferRequestedPerDHT,
		AccountAddr:             accountAddr,
		Voucher:                 voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRReqMsg(DHTOfferDiscoveryRequestType, nonce, body), nil
}

// DecodeDHTOfferDiscoveryRequest is used to get the fields from FCRMessage of dhtOfferDiscoveryRequestJson
// It returns the nonce, nodeID, pieceCID, numDHT, maxOfferRequestedPerDHT, account address and voucher.
func DecodeDHTOfferDiscoveryRequest(fcrMsg *FCRReqMsg) (
	uint64,
	string,
	*cid.ContentID,
	int64,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.Type() != DHTOfferDiscoveryRequestType {
		return 0, "", nil, 0, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", DHTOfferDiscoveryRequestType, fcrMsg.Type())
	}
	msg := dhtOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", nil, 0, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return 0, "", nil, 0, 0, "", "", err
	}
	return fcrMsg.Nonce(), msg.NodeID, pieceCID, msg.NumDHT, msg.MaxOfferRequestedPerDHT, msg.AccountAddr, msg.Voucher, nil
}
