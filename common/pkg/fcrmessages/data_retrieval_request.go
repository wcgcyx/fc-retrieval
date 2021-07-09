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

// dataRetrievalRequestJson represents the request to retrieve a piece of data.
type dataRetrievalRequestJson struct {
	SenderID    string `json:"sender_id"`
	Offer       string `json:"offer"`
	AccountAddr string `json:"account_addr"`
	Voucher     string `json:"voucher"`
}

// EncodeDataRetrievalRequest is used to get the FCRMessage of dataRetrievalRequest.
func EncodeDataRetrievalRequest(
	nonce uint64,
	senderID string,
	offer *cidoffer.SubCIDOffer,
	accountAddr string,
	voucher string,
) (*FCRReqMsg, error) {
	data, err := offer.ToBytes()
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(dataRetrievalRequestJson{
		SenderID:    senderID,
		Offer:       hex.EncodeToString(data),
		AccountAddr: accountAddr,
		Voucher:     voucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRReqMsg(DataRetrievalRequestType, nonce, body), nil
}

// DecodeDataRetrievalRequest is used to get the fields from FCRMessage of dataRetrievalRequest.
// It returns the nonce, offer, account address and voucher.
func DecodeDataRetrievalRequest(fcrMsg *FCRReqMsg) (
	uint64,
	string,
	*cidoffer.SubCIDOffer,
	string,
	string,
	error,
) {
	if fcrMsg.Type() != DataRetrievalRequestType {
		return 0, "", nil, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", DataRetrievalRequestType, fcrMsg.Type())
	}
	msg := dataRetrievalRequestJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", nil, "", "", err
	}
	data, err := hex.DecodeString(msg.Offer)
	if err != nil {
		return 0, "", nil, "", "", err
	}
	offer := cidoffer.SubCIDOffer{}
	err = offer.FromBytes(data)
	if err != nil {
		return 0, "", nil, "", "", err
	}
	return fcrMsg.Nonce(), msg.SenderID, &offer, msg.AccountAddr, msg.Voucher, nil
}
