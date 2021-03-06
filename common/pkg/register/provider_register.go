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

// ProviderRegisteredInfo represents the state or the stored information of a registered provider.
type ProviderRegisteredInfo struct {
	// RootKey is the provider's public key,
	// and can be used to derive the filecoin wallet address for payment.
	// It is a 65 bytes hex string.
	RootKey string `json:"root_key"`

	// NodeIDD is derived from the root key, it is set by the smart contract.
	// It is a 32 bytes hex string.
	NodeID string `json:"node_id"`

	// MsgSigningKey is the message signing public key.
	// It is a 65 bytes hex string.
	// This can be updated via smart contract.
	MsgSigningKey string `json:"msg_signing_key"`

	// MsgSigningKeyVer is the message signing public key version.
	// This can be updated via smart contract.
	MsgSigningKeyVer byte `json:"msg_signing_key_ver"`

	// OfferSigningKey is the offer signing public key.
	// It is a 65 bytes hex string.
	OfferSigningKey string `json:"offer_signing_key"`

	// RegionCode is the region code of this provider.
	// It is a ISO 3166-1 alpha-2 string.
	RegionCode string `json:"region_code"`

	// NetworkAddr is the network address of this provider.
	// It should be a valid libp2p address.
	NetworkAddr string `json:"network_addr"`

	// Deregistering indicates whether or not this provider is in the middle of deregistering itself.
	// It is set by the smart contract.
	Deregistering bool `json:"deregistering"`

	// DeregisteringHeight is the height of the block which contains the deregistering transaction.
	// It is set by the smart contract.
	DeregisteringHeight uint64 `json:"deregistering_height"`
}
