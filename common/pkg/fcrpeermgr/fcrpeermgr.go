/*
Package fcrpeermgr - peer manager manages all retrieval peers.
*/
package fcrpeermgr

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

// FCRPeerMgr represents the manager that manages all peers.
type FCRPeerMgr interface {
	// Start starts the manager's routine.
	Start() error

	// Shutdown ends the manager's routine safely.
	Shutdown()

	// Sync forces the manager to do a sync to the register.
	Sync()

	// SyncGW forces the manager to do a quick sync to the register for a specific gateway.
	SyncGW(gwID string)

	// SyncPVD forces the manager to do a quick sync to the register for a specific provider.
	SyncPVD(pvdID string)

	// GetGWInfo gets the data of a gateway, it queries the local storage, rather than the remote register.
	GetGWInfo(gwID string) (*Peer, error)

	// GetPVDInfo gets the data of a provider, it queries the local storage rather than the remote register.
	GetPVDInfo(pvdID string) (*Peer, error)

	// ListGWS lists all the gateways
	ListGWS() ([]Peer, error)

	// GetGWSNearCID gets 16 gateways that are near given CID. Called only by gateways.
	GetGWSNearCIDHash(hash string, except string) ([]Peer, error)

	// GetCurrentCIDHashRange gets the cid min hash and cid max hash that a gateway should store based on current network. Called only by gateways.
	GetCurrentCIDHashRange() (string, string, error)
}

// Peer represents a peer in the system.
type Peer struct {
	// RootKey is the peer's public key,
	// and can be used to derive the filecoin wallet address for payment.
	// It is a 65 bytes hex string.
	RootKey string

	// NodeID is derived from the root key, it is set by the smart contract.
	// It is a 32 bytes hex string.
	NodeID string

	// MsgSigningKey is the message signing public key.
	MsgSigningKey string

	// MsgSigningKeyVer is the message signing public key version.
	MsgSigningKeyVer byte

	// OfferSigningKey is the offer signing public key.
	// It is a 32 bytes hex string. Empty for gateway peer.
	OfferSigningKey string

	// RegionCode is the region code of this peer.
	RegionCode string

	// NetworkAddr is the network address of this peer.
	NetworkAddr string

	// Deregistering indicates whether or not this peer is in the middle of deregistering itself.
	Deregistering bool

	// DeregisteringHeight is the height of the block which contains the deregistering transaction.
	DeregisteringHeight uint64
}
