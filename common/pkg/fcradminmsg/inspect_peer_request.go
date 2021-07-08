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

// inspectPeerRequestJson represents the request to inspect a peer.
type inspectPeerRequestJson struct {
	NodeID  string `json:"node_id"`
	Gateway bool   `json:"gateway"`
}

// EncodeInspectPeerRequest is used to get the byte array of inspectPeerRequestJson
func EncodeInspectPeerRequest(
	nodeID string,
	gateway bool,
) ([]byte, error) {
	return json.Marshal(&inspectPeerRequestJson{
		NodeID:  nodeID,
		Gateway: gateway,
	})
}

// DecodeInspectPeerRequest is used to get the fields from byte array of inspectPeerRequestJson
func DecodeInspectPeerRequest(data []byte) (
	string, // node id
	bool, // gateway
	error, // error
) {
	msg := inspectPeerRequestJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return "", false, err
	}
	return msg.NodeID, msg.Gateway, nil
}
