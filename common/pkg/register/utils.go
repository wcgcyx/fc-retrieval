/*
Package register - location for smart contract registration structs.
*/
package register

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

import "encoding/hex"

const (
	validNodeIDLen = 32
	validKeyLen    = 32
)

// ValidateGatewayInfo check if a given gateway info is valid.
func ValidateGatewayInfo(gwInfo *GatewayRegisteredInfo) bool {
	nodeID, err := hex.DecodeString(gwInfo.NodeID)
	if err != nil {
		return false
	}
	if len(nodeID) != validNodeIDLen {
		return false
	}
	key, err := hex.DecodeString(gwInfo.MsgSigningKey)
	if err != nil {
		return false
	}
	if len(key) != validKeyLen {
		return false
	}
	if gwInfo.RegionCode == "" || gwInfo.NetworkAddr == "" {
		return false
	}
	// TODO, Need to check if the region code and the network addr is valid.
	return true
}

// ValidateGatewayInfo check if a given provider info is valid.
func ValidateProviderInfo(pvdInfo *ProviderRegisteredInfo) bool {
	nodeID, err := hex.DecodeString(pvdInfo.NodeID)
	if err != nil {
		return false
	}
	if len(nodeID) != validNodeIDLen {
		return false
	}
	key, err := hex.DecodeString(pvdInfo.MsgSigningKey)
	if err != nil {
		return false
	}
	if len(key) != validKeyLen {
		return false
	}
	if pvdInfo.RegionCode == "" || pvdInfo.NetworkAddr == "" {
		return false
	}
	// TODO, Need to check if the region code and the network addr is valid.
	return true
}
