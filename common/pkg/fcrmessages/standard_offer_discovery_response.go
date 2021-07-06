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
	Offers        []string `json:"offers"`
	RefundVoucher string   `json:"refund_voucher"`
}

// EncodeStandardOfferDiscoveryResponse is used to get the FCRMessage of standardOfferDiscoveryResponseJson.
func EncodeStandardOfferDiscoveryResponse(
	nonce uint64,
	offers []cidoffer.SubCIDOffer,
	refundVoucher string,
) (*FCRACKMsg, error) {
	offersStr := make([]string, 0)
	for _, offer := range offers {
		data, err := offer.ToBytes()
		if err != nil {
			return nil, err
		}
		offersStr = append(offersStr, hex.EncodeToString(data))
	}
	body, err := json.Marshal(standardOfferDiscoveryResponseJson{
		Offers:        offersStr,
		RefundVoucher: refundVoucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRACKMsg(nonce, body), nil
}

// DecodeStandardOfferDiscoveryResponse is used to get the fields from FCRMessage of standardOfferDiscoveryResponseJson.
// It returns nonce, a list of offers, refund account address and voucher.
func DecodeStandardOfferDiscoveryResponse(fcrMsg *FCRACKMsg) (
	uint64,
	[]cidoffer.SubCIDOffer,
	string,
	error,
) {
	if !fcrMsg.ACK() {
		return 0, nil, "", fmt.Errorf("ACK is false")
	}
	msg := standardOfferDiscoveryResponseJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, nil, "", err
	}
	offers := make([]cidoffer.SubCIDOffer, 0)
	offersStr := msg.Offers
	for _, offerStr := range offersStr {
		data, err := hex.DecodeString(offerStr)
		if err != nil {
			return 0, nil, "", err
		}
		offer := cidoffer.SubCIDOffer{}
		err = offer.FromBytes(data)
		if err != nil {
			return 0, nil, "", err
		}
		offers = append(offers, offer)
	}
	return fcrMsg.Nonce(), offers, msg.RefundVoucher, nil
}
