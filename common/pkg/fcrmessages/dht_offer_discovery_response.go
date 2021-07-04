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
	"sort"
)

// dhtOfferDiscoveryResponseJson represents the response to a request of asking for offers in DHT.
type dhtOfferDiscoveryResponseJson struct {
	Contacted     []string `json:"contacted"`
	Responses     []string `json:"responses"`
	Nonce         int64    `json:"nonce"`
	RefundVoucher string   `json:"refund_voucher"`
}

// EncodeDHTOfferDiscoveryResponse is used to get the FCRMessage of dhtOfferDiscoveryResponseJson.
func EncodeDHTOfferDiscoveryResponse(
	contacted map[string]*FCRMessage,
	nonce int64,
	refundVoucher string,
) (*FCRMessage, error) {
	keys := make([]string, len(contacted))
	i := 0
	for k := range contacted {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	contactedStr := make([]string, 0)
	responsesStr := make([]string, 0)
	for _, key := range keys {
		contactedStr = append(contactedStr, key)
		data, err := contacted[key].ToBytes()
		if err != nil {
			return nil, err
		}
		responsesStr = append(responsesStr, hex.EncodeToString(data))
	}
	body, err := json.Marshal(dhtOfferDiscoveryResponseJson{
		Contacted:     contactedStr,
		Responses:     responsesStr,
		Nonce:         nonce,
		RefundVoucher: refundVoucher,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(DHTOfferDiscoveryResponseType, body), nil
}

// DecodeDHTOfferDiscoveryResponse is used to get the fields from FCRMessage of dhtOfferDiscoveryResponseJson.
// It returns a map of contacted nodes -> contacted messages, nonce, refund account address and voucher.
func DecodeDHTOfferDiscoveryResponse(fcrMsg *FCRMessage) (
	map[string]FCRMessage,
	int64,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != DHTOfferDiscoveryResponseType {
		return nil, 0, "", fmt.Errorf("Message type mismatch, expect %v, got %v", DHTOfferDiscoveryResponseType, fcrMsg.GetMessageType())
	}
	msg := dhtOfferDiscoveryResponseJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return nil, 0, "", err
	}
	if len(msg.Contacted) != len(msg.Responses) {
		return nil, 0, "", fmt.Errorf("Contacted length %v mismatches response length %v", len(msg.Contacted), len(msg.Responses))
	}
	contacted := make(map[string]FCRMessage)
	for i := 0; i < len(msg.Contacted); i++ {
		data, err := hex.DecodeString(msg.Responses[i])
		if err != nil {
			return nil, 0, "", err
		}
		resp, err := FromBytes(data)
		if err != nil {
			return nil, 0, "", err
		}
		_, ok := contacted[msg.Contacted[i]]
		if ok {
			return nil, 0, "", fmt.Errorf("Node %v appears at least twice in the response", msg.Contacted[i])
		}
		contacted[msg.Contacted[i]] = *resp
	}
	return contacted, msg.Nonce, msg.RefundVoucher, nil
}
