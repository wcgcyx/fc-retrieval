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

// ackJson represents the a ack to message.
type ackJson struct {
	ACK   bool   `json:"ack"`
	Nonce int64  `json:"nonce"`
	Data  string `json:"data"`
}

// EncodeACK is used to get the FCRMessage of ackJson
func EncodeACK(
	ack bool,
	nonce int64,
	data string,
) (*FCRMessage, error) {
	body, err := json.Marshal(ackJson{
		ACK:   ack,
		Nonce: nonce,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(ACKType, body), nil
}

// DecodeACK is used to get the fields from FCRMessage of ackJson
// It returns the ack, nonce and the data related to this ack.
func DecodeACK(fcrMsg *FCRMessage) (
	bool,
	int64,
	string,
	error,
) {
	if fcrMsg.GetMessageType() != ACKType {
		return false, 0, "", fmt.Errorf("Message type mismatch, expect %v, got %v", ACKType, fcrMsg.GetMessageType())
	}
	msg := ackJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return false, 0, "", err
	}
	return msg.ACK, msg.Nonce, msg.Data, nil
}
