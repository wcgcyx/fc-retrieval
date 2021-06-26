/*
Package register - location for smart contract registration structs.
*/
package register

/*
 * Copyright 2021 ConsenSys Software Inc.
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

// ValidateGatewayInfo check if a given gateway info is valid.
// It is used when before registering and updating.
func ValidateGatewayInfo(gwInfo *GatewayRegisteredInfo) bool {
	return true
}

// ValidateGatewayInfo check if a given provider info is valid.
// It is used when before registering and updating.
func ValidateProviderInfo(pvdInfo *ProviderRegisteredInfo) bool {
	return true
}
