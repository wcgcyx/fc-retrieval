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
	Nonce                   int64  `json:"nonce"`
	NumDHT                  int64  `json:"num_dht"`
	MaxOfferRequestedPerDHT int64  `json:"max_offer_requested_per_dht"`
	AccountAddr             string `json:"account_addr"`
	Voucher                 string `json:"voucher"`
}

// EncodeDHTOfferDiscoveryRequest is used to get the FCRMessage of dhtOfferDiscoveryRequestJson
func EncodeDHTOfferDiscoveryRequest(
	NodeID string,
	pieceCID *cid.ContentID,
	nonce int64,
	numDHT int64,
	maxOfferRequestedPerDHT int64,
	accountAddr string,
	voucher string,
) (*FCRMessage, error) {
	body, err := json.Marshal(dhtOfferDiscoveryRequestJson{
		NodeID:                  NodeID,
		PieceCID:                pieceCID.ToString(),
		Nonce:                   nonce,
		NumDHT:                  numDHT,
		MaxOfferRequestedPerDHT: maxOfferRequestedPerDHT,
		AccountAddr:             accountAddr,
		Voucher:                 voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(DHTOfferDiscoveryRequestType, body), nil
}

// DecodeDHTOfferDiscoveryRequest is used to get the fields from FCRMessage of dhtOfferDiscoveryRequestJson
// It returns the nodeID, pieceCID, nonce, numDHT, maxOfferRequestedPerDHT, account address and voucher.
func DecodeDHTOfferDiscoveryRequest(fcrMsg *FCRMessage) (
	string,
	*cid.ContentID,
	int64,
	int64,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != DHTOfferDiscoveryRequestType {
		return "", nil, 0, 0, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", DHTOfferDiscoveryRequestType, fcrMsg.GetMessageType())
	}
	msg := dhtOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", nil, 0, 0, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return "", nil, 0, 0, 0, "", "", err
	}
	return msg.NodeID, pieceCID, msg.Nonce, msg.NumDHT, msg.MaxOfferRequestedPerDHT, msg.AccountAddr, msg.Voucher, nil
}
