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

import "encoding/json"

// getOfferByCIDResponseJson represents the response of getting offer.
type getOfferByCIDResponseJson struct {
	Digests   []string
	Providers []string
	Prices    []string
	Expiry    []int64
	QoS       []uint64
}

// EncodeGetOfferByCIDResponse is used to get the byte array of getOfferByCIDResponseJson
func EncodeGetOfferByCIDResponse(
	digests []string,
	providers []string,
	prices []string,
	expiry []int64,
	qos []uint64,
) ([]byte, error) {
	return json.Marshal(&getOfferByCIDResponseJson{
		Digests:   digests,
		Providers: providers,
		Prices:    prices,
		Expiry:    expiry,
		QoS:       qos,
	})
}

// DecodeGetOfferByCIDResponse is used to get the fields from byte array of getOfferByCIDResponseJson
func DecodeGetOfferByCIDResponse(data []byte) (
	[]string, // digests
	[]string, // providers
	[]string, // prices
	[]int64, // expiry
	[]uint64, // qos
	error, // error
) {
	msg := getOfferByCIDResponseJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return msg.Digests, msg.Providers, msg.Prices, msg.Expiry, msg.QoS, nil
}
