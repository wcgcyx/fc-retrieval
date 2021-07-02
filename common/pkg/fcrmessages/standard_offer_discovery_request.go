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
	NodeID            string `json:"node_id"`
	PieceCID          string `json:"piece_cid"`
	Nonce             int64  `json:"nonce"`
	MaxOfferRequested int64  `json:"max_offer_requested"`
	AccountAddr       string `json:"account_addr"`
	Voucher           string `json:"voucher"`
}

// EncodeStandardOfferDiscoveryRequest is used to get the FCRMessage of standardOfferDiscoveryRequestJson.
func EncodeStandardOfferDiscoveryRequest(
	NodeID string,
	pieceCID *cid.ContentID,
	nonce int64,
	maxOfferRequested int64,
	accountAddr string,
	voucher string,
) (*FCRMessage, error) {
	body, err := json.Marshal(standardOfferDiscoveryRequestJson{
		NodeID:            NodeID,
		PieceCID:          pieceCID.ToString(),
		Nonce:             nonce,
		MaxOfferRequested: maxOfferRequested,
		AccountAddr:       accountAddr,
		Voucher:           voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(StandardOfferDiscoveryRequestType, body), nil
}

// DecodeStandardOfferDiscoveryRequest is used to get the fields from FCRMessage of standardOfferDiscoveryRequestJson.
// It returns the nodeID, pieceCID, nonce, maxOfferRequested, account address and voucher.
func DecodeStandardOfferDiscoveryRequest(fcrMsg *FCRMessage) (
	string,
	*cid.ContentID,
	int64,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != StandardOfferDiscoveryRequestType {
		return "", nil, 0, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", StandardOfferDiscoveryRequestType, fcrMsg.GetMessageType())
	}
	msg := standardOfferDiscoveryRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", nil, 0, 0, "", "", err
	}
	pieceCID, err := cid.NewContentID(msg.PieceCID)
	if err != nil {
		return "", nil, 0, 0, "", "", err
	}
	return msg.NodeID, pieceCID, msg.Nonce, msg.MaxOfferRequested, msg.AccountAddr, msg.Voucher, nil
}
