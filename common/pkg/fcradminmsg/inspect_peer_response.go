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

// inspectPeerResponseJson represents the response of inspecting a peer.
type inspectPeerResponseJson struct {
	Score   int64    `json:"score"`
	Pending bool     `json:"pending"`
	Blocked bool     `json:"blocked"`
	History []string `json:"history"`
}

// EncodeInspectPeerResponse is used to get the byte array of inspectPeerResponseJson
func EncodeInspectPeerResponse(
	score int64,
	pending bool,
	blocked bool,
	history []string,
) ([]byte, error) {
	return json.Marshal(&inspectPeerResponseJson{
		Score:   score,
		Pending: pending,
		Blocked: blocked,
		History: history,
	})
}

// DecodeInspectPeerResponse is used to get the fields from byte array of inspectPeerResponseJson
func DecodeInspectPeerResponse(data []byte) (
	int64, // score
	bool, // pending
	bool, // blocked
	[]string, // history
	error, // error
) {
	msg := inspectPeerResponseJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return 0, false, false, nil, err
	}
	return msg.Score, msg.Pending, msg.Blocked, msg.History, nil
}
