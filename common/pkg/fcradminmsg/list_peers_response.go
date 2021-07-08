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
	GatewayIDs      []string `json:"gateway_ids"`
	GatewayScore    []int64  `json:"gateway_score"`
	GatewayPending  []bool   `json:"gateway_pending"`
	GatewayBlocked  []bool   `json:"gateway_blocked"`
	GatewayRecent   []string `json:"gateway_recent"`
	ProviderIDs     []string `json:"provider_ids"`
	ProviderScore   []int64  `json:"provider_score"`
	ProviderPending []bool   `json:"provider_pending"`
	ProviderBlocked []bool   `json:"provider_blocked"`
	ProviderRecent  []string `json:"provider_recent"`
}

// EncodeListPeersResponse is used to get the byte array of listPeersResponseJson
func EncodeListPeersResponse(
	gatewayIDs []string,
	gatewayScore []int64,
	gatewayPending []bool,
	gatewayBlocked []bool,
	gatewayRecent []string,
	providerIDs []string,
	providerScore []int64,
	providerPending []bool,
	providerBlocked []bool,
	providerRecent []string,
) ([]byte, error) {
	return json.Marshal(&listPeersResponseJson{
		GatewayIDs:      gatewayIDs,
		GatewayScore:    gatewayScore,
		GatewayPending:  gatewayPending,
		GatewayBlocked:  gatewayBlocked,
		GatewayRecent:   gatewayRecent,
		ProviderIDs:     providerIDs,
		ProviderScore:   providerScore,
		ProviderPending: providerPending,
		ProviderBlocked: providerBlocked,
		ProviderRecent:  providerRecent,
	})
}

// DecodeListPeersResponse is used to get the fields from byte array of listPeersResponseJson
func DecodeListPeersResponse(data []byte) (
	[]string, // gateway IDs
	[]int64, // score
	[]bool, // pending
	[]bool, // blocked
	[]string, // most recent activity
	[]string, // provider IDs
	[]int64, // score
	[]bool, // pending
	[]bool, // blocked
	[]string, // most recent activity
	error, // error
) {
	msg := listPeersResponseJson{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	return msg.GatewayIDs, msg.GatewayScore, msg.GatewayPending, msg.GatewayBlocked, msg.GatewayRecent, msg.ProviderIDs, msg.ProviderScore, msg.ProviderPending, msg.ProviderBlocked, msg.ProviderRecent, nil
}
