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

// standardOfferDiscoveryResponseJson represents the response to a request of asking for offers.
type standardOfferDiscoveryResponseJson struct {
	Offers            []string `json:"offers"`
	Nonce             int64    `json:"nonce"`
	RefundAccountAddr string   `json:"refund_account_addr"`
	RefundVoucher     string   `json:"refund_voucher"`
}

// EncodeStandardOfferDiscoveryResponse is used to get the FCRMessage of standardOfferDiscoveryResponseJson.
func EncodeStandardOfferDiscoveryResponse(
	offers []cidoffer.SubCIDOffer,
	nonce int64,
	refundAccountAddr string,
	refundVoucher string,
) (*FCRMessage, error) {
	offersStr := make([]string, 0)
	for _, offer := range offers {
		data, err := offer.ToBytes()
		if err != nil {
			return nil, err
		}
		offersStr = append(offersStr, hex.EncodeToString(data))
	}
	body, err := json.Marshal(standardOfferDiscoveryResponseJson{
		Offers:            offersStr,
		Nonce:             nonce,
		RefundAccountAddr: refundAccountAddr,
		RefundVoucher:     refundVoucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(StandardOfferDiscoveryResponseType, body), nil
}

// DecodeStandardOfferDiscoveryResponse is used to get the fields from FCRMessage of standardOfferDiscoveryResponseJson.
// It returns a list of offers, nonce, refund account address and voucher.
func DecodeStandardOfferDiscoveryResponse(fcrMsg *FCRMessage) (
	[]cidoffer.SubCIDOffer,
	int64,
	string,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != StandardOfferDiscoveryResponseType {
		return nil, 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", StandardOfferDiscoveryResponseType, fcrMsg.GetMessageType())
	}
	msg := standardOfferDiscoveryResponseJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, "", "", err
	}
	offers := make([]cidoffer.SubCIDOffer, 0)
	offersStr := msg.Offers
	for _, offerStr := range offersStr {
		data, err := hex.DecodeString(offerStr)
		if err != nil {
			return nil, 0, "", "", err
		}
		offer := cidoffer.SubCIDOffer{}
		err = offer.FromBytes(data)
		if err != nil {
			return nil, 0, "", "", err
		}
		offers = append(offers, offer)
	}
	return offers, msg.Nonce, msg.RefundAccountAddr, msg.RefundVoucher, nil
}
