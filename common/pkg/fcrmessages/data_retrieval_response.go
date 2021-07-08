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
)

// dataRetrievalResponseJson represents the response to a request of asking for offers.
type dataRetrievalResponseJson struct {
	Tag  string `json:"tag"`
	Data []byte `json:"data"`
}

// EncodeDataRetrievalResponse is used to get the FCRMessage of dataRetrievalResponseJson.
func EncodeDataRetrievalResponse(
	nonce uint64,
	tag string,
	data []byte,
) (*FCRACKMsg, error) {
	body, err := json.Marshal(dataRetrievalResponseJson{
		Tag:  tag,
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRACKMsg(nonce, body), nil
}

// DecodeDataRetrievalResponse is used to get the fields from FCRMessage of dataRetrievalResponseJson.
// It returns the nonce, tag, and file data.s
func DecodeDataRetrievalResponse(fcrMsg *FCRACKMsg) (
	uint64,
	string,
	[]byte,
	error,
) {
	if !fcrMsg.ACK() {
		return 0, "", nil, fmt.Errorf("ACK is false")
	}
	msg := dataRetrievalResponseJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", nil, err
	}
	return fcrMsg.Nonce(), msg.Tag, msg.Data, nil
}
