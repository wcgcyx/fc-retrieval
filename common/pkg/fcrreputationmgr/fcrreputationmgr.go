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

	// AddPeer starts tracking given peer's reputation.
	AddPeer(peerID string)

	// ListPeers lists all peers' ids in tracking.
	ListPeers() []string

	// RemvoePeer stops tracking given peer's reputation.
	RemovePeer(peerID string)

	// GetPeerReputation gets the reputation of a given peer.
	GetPeerReputation(peerID string) *Reputation

	// UpdatePeerRecord updates the peer's reputation with a given record and a given replica (replica = 1 means one additional application of the same record).
	UpdatePeerRecord(peerID string, record *reputation.Record, replica uint)

	// PendPeer puts a given peer into pending.
	PendPeer(peerID string)

	// ResumePeer puts a given peer out of pending.
	ResumePeer(peerID string)

	// GetPendingPeers gets a list of peers currently in pending.
	GetPendingPeers() []string

	// BlockPeer blocks a peer.
	BlockPeer(peerID string)

	// UnBlockPeer unblocks a peer.
	UnBlockPeer(peerID string)

	// GetBlockedPeers gets a list of blocked peers.
	GetBlockedPeers() []string

	// GetPeerViolations gets a list of violations from given index to given index for a given peer.
	GetPeerViolations(peerID string, from uint, to uint) []reputation.Record

	// GetPeerHistory gets a list of history from given index to given index for a given peer.
	GetPeerHistory(peerID string, from uint, to uint) []reputation.Record
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
