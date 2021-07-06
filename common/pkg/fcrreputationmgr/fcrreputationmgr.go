/*
Package fcrreputationmgr - reputation manager manages the reputation of all retrieval peers.
*/
package fcrreputationmgr

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

import "github.com/wcgcyx/fc-retrieval/common/pkg/reputation"

// FCRReputationMgr represents the manager that manages all reputation.
type FCRReputationMgr interface {
	// Start starts the manager's routine.
	Start() error

	// Shutdown ends the manager's routine safely.
	Shutdown()

	// AddGW starts tracking given gateway's reputation.
	AddGW(gwID string)

	// ListGWS lists all gateways in tracking.
	ListGWS() []string

	// RemvoeGW stops tracking given gateway's reputation.
	RemoveGW(gwID string)

	// AddPVD starts tracking given provider's reputation.
	AddPVD(pvdID string)

	// ListPVDS lists all providers in tracking.
	ListPVDS() []string

	// RemovePVD stops tracking given provider's reputation.
	RemovePVD(pvdID string)

	// GetGWReputation gets the reputation of a given gateway.
	GetGWReputation(gwID string) *Reputation

	// GetPVDReputation gets the reputation of a given provider.
	GetPVDReputation(pvdID string) *Reputation

	// UpdateGWRecord updates the gateway's reputation with a given record and a given replica (replica = 1 means one additional application of the same record).
	UpdateGWRecord(gwID string, record *reputation.Record, replica uint)

	// UpdatePVDRecord updates the provider's reputation with a given record and a given replica (replica = 1 means one additional application of the same record).
	UpdatePVDRecord(pvdID string, record *reputation.Record, replica uint)

	// PendGW puts a given gateway into pending.
	PendGW(gwID string)

	// PendPVD puts a given provider into pending.
	PendPVD(pvdID string)

	// ResumeGW puts a given gateway out of pending.
	ResumeGW(gwID string)

	// ResumePVD puts a provider out of pending.
	ResumePVD(pvdID string)

	// GetPendingGWS gets a list of gateways currently in pending.
	GetPendingGWS() []string

	// GetPendingPVDS gets a list of providers currently in pending.
	GetPendingPVDS() []string

	// BlockGW blocks a gateway.
	BlockGW(gwID string)

	// BlockPVD blocks a provider.
	BlockPVD(pvdID string)

	// UnBlockGW unblocks a gateway.
	UnBlockGW(gwID string)

	// UnBlockPVD unblocks a provider.
	UnBlockPVD(pvdID string)

	// GetBlockedGWS gets a list of blocked gateways.
	GetBlockedGWS() []string

	// GetBlockedPVDS gets a list of blocked providers.
	GetBlockedPVDS() []string

	// GetGWViolations gets a list of violations from given index to given index for a given gateway.
	GetGWViolations(gwID string, from uint, to uint) []reputation.Record

	// GetPVDViolations gets a list of violations from given index to given index for a given provider.
	GetPVDViolations(pvdID string, from uint, to uint) []reputation.Record

	// GetGWHistory gets a list of history from given index to given index for a given gateway.
	GetGWHistory(gwID string, from uint, to uint) []reputation.Record

	// GetPVDHistory gets a list of history from given index to given index for a given provider.
	GetPVDHistory(pvdID string, from uint, to uint) []reputation.Record
}

// Reputation represents the reputation of a peer in the system.
type Reputation struct {
	// NodeID is the peer's ID
	NodeID string

	// Score is the overall reputation score of this peer
	Score int64

	// Pending indicates whether this peer is in pending
	Pending bool

	// Blocked indicates whether this peer is blocked
	Blocked bool
}
