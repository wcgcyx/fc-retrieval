/*
Package fcradminmsg - stores all the admin messages.
*/
package fcradminmsg

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
	"errors"
	"math/big"
)

// publishOfferRequestJson represents the request to publish offer.
type publishOfferRequestJson struct {
	Files  []string `json:"files"`
	Price  string   `json:"price"`
	Expiry int64    `json:"expiry"`
	QoS    uint64   `json:"qos"`
}

// EncodePublishOfferRequest is used to get the byte array of publishOfferRequestJson
func EncodePublishOfferRequest(
	files []string,
	price *big.Int,
	expiry int64,
	qos uint64,
) ([]byte, error) {
	return json.Marshal(&publishOfferRequestJson{
		Files:  files,
		Price:  price.String(),
		Expiry: expiry,
		QoS:    qos,
	})
}

// DecodePublishOfferRequest is used to get the fields from byte array of publishOfferRequestJson
func DecodePublishOfferRequest(data []byte) (
	[]string, // files
	*big.Int, // price
	int64, // expiry
	uint64, // qos
	error, // error
) {
	msg := publishOfferRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, 0, 0, err
	}
	price, ok := big.NewInt(0).SetString(msg.Price, 10)
	if !ok {
		return nil, nil, 0, 0, errors.New("Error in decoding price")
	}
	return msg.Files, price, msg.Expiry, msg.QoS, nil
}
