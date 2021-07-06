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

// ackJson represents the a ack to message.
type ackJson struct {
	ACK  bool   `json:"ack"`
	Data string `json:"data"`
}

// EncodeACK is used to get the byte array of ackJson
func EncodeACK(
	ack bool,
	data string,
) []byte {
	res, _ := json.Marshal(&ackJson{
		ACK:  ack,
		Data: data,
	})
	return res
}

// DecodeACK is used to get the fields from byte array of ackJson
// It returns the ack, data and error
func DecodeACK(data []byte) (
	bool,
	string,
	error,
) {
	msg := ackJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return false, "", err
	}
	return msg.ACK, msg.Data, nil
}
