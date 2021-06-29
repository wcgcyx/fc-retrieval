/*
Package fcrpeermgr - peer manager manages all retrieval peers.
*/
package fcrpeermgr

import (
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

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

type FCRPeerMgr interface {
	// Start starts the manager's routine.
	Start()

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

	// GetGWSNearCID gets 16 gateways that are near given CID. Called only by gateways.
	GetGWSNearCID(id *cid.ContentID) ([]Peer, error)

	// GetCIDHashRange gets the cid min hash and cid max hash that a gateway should store based on current network. Called only by gateways.
	GetCIDHashRange() (string, string, error)

	/* Reputation related functions */

	// UpdateGWRecord updates the gateway's reputation with a given record and a given replica (replica = 1 means one additional application of the same record).
	UpdateGWRecord(gwID string, record *reputation.Record, replica int) error

	// UpdatePVDRecord updates the provider's reputation with a given record and a given replica (replica = 1 means one additional application of the same record).
	UpdatePVDRecord(gwID string, record *reputation.Record, replica int) error

	// PendGW puts a given gateway into pending.
	PendGW(gwID string) error

	// PendPVD puts a given provider into pending.
	PendPVD(pvdID string) error

	// ResumeGW puts a given gateway out of pending.
	ResumeGW(gwID string) error

	// ResumePVD puts a provider out of pending.
	ResumePVD(pvdID string) error

	// GetPendingGWS gets a list of gateways currently in pending.
	GetPendingGWS() ([]Peer, error)

	// GetPendingPVDS gets a list of providers currently in pending.
	GetPendingPVDS() ([]Peer, error)

	// BlockGW blocks a gateway.
	BlockGW(gwID string) error

	// BlockPVD blocks a provider.
	BlockPVD(pvdID string) error

	// UnBlockGW unblocks a gateway.
	UnBlockGW(gwID string) error

	// UnBlockPVD unblocks a provider.
	UnBlockPVD(pvdID string) error

	// GetBlockedGWS gets a list of blocked gateways.
	GetBlockedGWS() ([]Peer, error)

	// GetBlockedPVDS gets a list of blocked providers.
	GetBlockedPVDS() ([]Peer, error)
}

// Peer represents a peer in the system.
type Peer struct {
	NodeID           string
	MsgSigningKey    string
	MsgSigningKeyVER string
	OfferSigningKey  string
	RegionCode       string
	NetworkAddr      string

	/* Reputation related fields */

	// Reputation is the overall reputation score of this peer
	Reputation int64
	// Pending indicates whether this peer is in pending
	Pending bool
	// Blocked indicates whether this peer is blocked
	Blocked bool
	// Violations store the list of recent violations (recent 50 entries)
	Violations []reputation.Record
	// History stores the list of recent record updating activies (recent 500 entries)
	History []reputation.Record
}
