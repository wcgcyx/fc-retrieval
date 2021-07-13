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

// listPeersResponseJson represents the response of listing peers.
type listPeersResponseJson struct {
	PeerIDs     []string `json:"peer_ids"`
	PeerScore   []int64  `json:"peer_score"`
	PeerPending []bool   `json:"peer_pending"`
	PeerBlocked []bool   `json:"peer_blocked"`
	PeerRecent  []string `json:"peer_recent"`
}

// EncodeListPeersResponse is used to get the byte array of listPeersResponseJson
func EncodeListPeersResponse(
	peerIDs []string,
	peerScore []int64,
	peerPending []bool,
	peerBlocked []bool,
	peerRecent []string,
) ([]byte, error) {
	return json.Marshal(&listPeersResponseJson{
		PeerIDs:     peerIDs,
		PeerScore:   peerScore,
		PeerPending: peerPending,
		PeerBlocked: peerBlocked,
		PeerRecent:  peerRecent,
	})
}

// DecodeListPeersResponse is used to get the fields from byte array of listPeersResponseJson
func DecodeListPeersResponse(data []byte) (
	[]string, // peer IDs
	[]int64, // score
	[]bool, // pending
	[]bool, // blocked
	[]string, // most recent activity
	error, // error
) {
	msg := listPeersResponseJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return msg.PeerIDs, msg.PeerScore, msg.PeerPending, msg.PeerBlocked, msg.PeerRecent, nil
}
