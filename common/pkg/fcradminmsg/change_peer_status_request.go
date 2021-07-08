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

// changePeerStatusRequestJson represents the request to change peer status.
type changePeerStatusRequestJson struct {
	NodeID  string `json:"node_id"`
	Gateway bool   `json:"gateway"`
	Block   bool   `json:"block"`
	Unblock bool   `json:"unblock"`
}

// EncodeChangePeerStatusRequest is used to get the byte array of changePeerStatusRequestJson
func EncodeChangePeerStatusRequest(
	nodeID string,
	gateway bool,
	block bool,
	unblock bool,
) ([]byte, error) {
	return json.Marshal(&changePeerStatusRequestJson{
		NodeID:  nodeID,
		Gateway: gateway,
		Block:   block,
		Unblock: unblock,
	})
}

// DecodeChangePeerStatusRequest is used to get the fields from changePeerStatusRequestJson
func DecodeChangePeerStatusRequest(data []byte) (
	string, // node id
	bool, // gateway
	bool, // block
	bool, // unblock
	error, // error
) {
	msg := changePeerStatusRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", false, false, false, err
	}
	return msg.NodeID, msg.Gateway, msg.Block, msg.Unblock, nil
}
