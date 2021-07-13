/*
Package fcrreputationmgr - reputation manager manages the reputation of all retrieval Peers.
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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

func TestAddRemovePeer(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep := mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, int64(0), rep.Score)
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000000")
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000002")
	res := mgr.ListPeers()
	assert.Equal(t, 3, len(res))
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, rep)
	reps := mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, reps)
	reps = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, reps)
	mgr.RemovePeer("0000000000000000000000000000000000000000000000000000000000000001")
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.Empty(t, rep)
	reps = mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, reps)
	reps = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, reps)
}

func TestPendingBlockPeer(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep := mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, false, rep.Pending)
	pending := mgr.GetPendingPeers()
	assert.Empty(t, pending)
	mgr.PendPeer("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.PendPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, true, rep.Pending)
	pending = mgr.GetPendingPeers()
	assert.Equal(t, 1, len(pending))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", pending[0])
	mgr.ResumePeer("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.ResumePeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, false, rep.Pending)
	pending = mgr.GetPendingPeers()
	assert.Empty(t, pending)

	blocked := mgr.GetBlockedPeers()
	assert.Empty(t, blocked)
	mgr.BlockPeer("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.BlockPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, true, rep.Blocked)
	blocked = mgr.GetBlockedPeers()
	assert.Equal(t, 1, len(blocked))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", blocked[0])
	mgr.UnBlockPeer("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.UnBlockPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, false, rep.Blocked)
	blocked = mgr.GetBlockedPeers()
	assert.Empty(t, blocked)
}

func TestAddRecordPeer(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPeer("0000000000000000000000000000000000000000000000000000000000000000")
	rep := mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, int64(0), rep.Score)

	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000001", &reputation.MockGoodRecord, 0)
	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockGoodRecord, 0)

	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, int64(1000), rep.Score)

	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)

	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, int64(999), rep.Score)

	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	mgr.UpdatePeerRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	rep = mgr.GetPeerReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, int64(995), rep.Score)

	// Now we have 6 entires of history, 5 entries of violations
	violations := mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.Empty(t, violations)

	violations = mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Equal(t, 2, len(violations))

	violations = mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Equal(t, 5, len(violations))

	violations = mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Equal(t, 3, len(violations))

	violations = mgr.GetPeerViolations("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, violations)

	history := mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.Empty(t, history)

	history = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Equal(t, 2, len(history))

	history = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Equal(t, 6, len(history))
	assert.Equal(t, false, history[5].Violation())

	history = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Equal(t, 4, len(history))
	assert.Equal(t, false, history[3].Violation())

	history = mgr.GetPeerHistory("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, history)
}
