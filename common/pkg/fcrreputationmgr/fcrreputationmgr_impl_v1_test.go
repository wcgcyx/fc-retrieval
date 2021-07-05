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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

func TestAddRemoveGW(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(0), rep.Score)
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000000")
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000002")
	res := mgr.ListGWS()
	assert.Equal(t, 3, len(res))
	_, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.Empty(t, err)
	_, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, err)
	_, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, err)
	mgr.RemoveGW("0000000000000000000000000000000000000000000000000000000000000001")
	_, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	_, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.NotEmpty(t, err)
	_, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.NotEmpty(t, err)
}

func TestAddRemovePVD(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(0), rep.Score)
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000000")
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000001")
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000002")
	res := mgr.ListPVDS()
	assert.Equal(t, 3, len(res))
	_, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.Empty(t, err)
	_, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, err)
	_, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.Empty(t, err)
	mgr.RemovePVD("0000000000000000000000000000000000000000000000000000000000000001")
	_, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	_, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.NotEmpty(t, err)
	_, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000001", 0, 10)
	assert.NotEmpty(t, err)
}

func TestPendingBlockGW(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Pending)
	pending, err := mgr.GetPendingGWS()
	assert.Empty(t, err)
	assert.Empty(t, pending)
	err = mgr.PendGW("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.PendGW("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, true, rep.Pending)
	pending, err = mgr.GetPendingGWS()
	assert.Empty(t, err)
	assert.Equal(t, 1, len(pending))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", pending[0])
	err = mgr.ResumeGW("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.ResumeGW("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Pending)
	pending, err = mgr.GetPendingGWS()
	assert.Empty(t, err)
	assert.Empty(t, pending)

	blocked, err := mgr.GetBlockedGWS()
	assert.Empty(t, err)
	assert.Empty(t, blocked)
	err = mgr.BlockGW("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.BlockGW("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, true, rep.Blocked)
	blocked, err = mgr.GetBlockedGWS()
	assert.Empty(t, err)
	assert.Equal(t, 1, len(blocked))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", blocked[0])
	err = mgr.UnBlockGW("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.UnBlockGW("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Blocked)
	blocked, err = mgr.GetBlockedGWS()
	assert.Empty(t, err)
	assert.Empty(t, blocked)
}

func TestPendingBlockPVD(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Pending)
	pending, err := mgr.GetPendingPVDS()
	assert.Empty(t, err)
	assert.Empty(t, pending)
	err = mgr.PendPVD("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.PendPVD("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, true, rep.Pending)
	pending, err = mgr.GetPendingPVDS()
	assert.Empty(t, err)
	assert.Equal(t, 1, len(pending))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", pending[0])
	err = mgr.ResumePVD("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.ResumePVD("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Pending)
	pending, err = mgr.GetPendingPVDS()
	assert.Empty(t, err)
	assert.Empty(t, pending)

	blocked, err := mgr.GetBlockedPVDS()
	assert.Empty(t, err)
	assert.Empty(t, blocked)
	err = mgr.BlockPVD("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.BlockPVD("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, true, rep.Blocked)
	blocked, err = mgr.GetBlockedPVDS()
	assert.Empty(t, err)
	assert.Equal(t, 1, len(blocked))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", blocked[0])
	err = mgr.UnBlockPVD("0000000000000000000000000000000000000000000000000000000000000001")
	assert.NotEmpty(t, err)
	err = mgr.UnBlockPVD("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, false, rep.Blocked)
	blocked, err = mgr.GetBlockedPVDS()
	assert.Empty(t, err)
	assert.Empty(t, blocked)
}

func TestAddRecordGW(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddGW("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(0), rep.Score)

	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000001", &reputation.MockGoodRecord, 0)
	assert.NotEmpty(t, err)
	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockGoodRecord, 0)
	assert.Empty(t, err)

	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(1000), rep.Score)

	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)

	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(999), rep.Score)

	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdateGWRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	rep, err = mgr.GetGWReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(995), rep.Score)

	// Now we have 6 entires of history, 5 entries of violations
	_, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.NotEmpty(t, err)

	violations, err := mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Empty(t, err)
	assert.Equal(t, 2, len(violations))

	violations, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Empty(t, err)
	assert.Equal(t, 5, len(violations))

	violations, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(violations))

	violations, err = mgr.GetGWViolations("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, err)
	assert.Equal(t, 0, len(violations))

	_, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.NotEmpty(t, err)

	history, err := mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Empty(t, err)
	assert.Equal(t, 2, len(history))

	history, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Empty(t, err)
	assert.Equal(t, 6, len(history))
	assert.Equal(t, false, history[5].Violation())

	history, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Empty(t, err)
	assert.Equal(t, 4, len(history))
	assert.Equal(t, false, history[3].Violation())

	history, err = mgr.GetGWHistory("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, err)
	assert.Equal(t, 0, len(history))
}

func TestAddRecordPVD(t *testing.T) {
	mgr := NewFCRReputationMgrImpV1()
	err := mgr.Start()
	defer mgr.Shutdown()
	assert.Empty(t, err)
	mgr.AddPVD("0000000000000000000000000000000000000000000000000000000000000000")
	rep, err := mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(0), rep.Score)

	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000001", &reputation.MockGoodRecord, 0)
	assert.NotEmpty(t, err)
	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockGoodRecord, 0)
	assert.Empty(t, err)

	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(1000), rep.Score)

	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)

	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(999), rep.Score)

	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	err = mgr.UpdatePVDRecord("0000000000000000000000000000000000000000000000000000000000000000", &reputation.MockBadRecord, 0)
	assert.Empty(t, err)
	rep, err = mgr.GetPVDReputation("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, err)
	assert.Equal(t, int64(995), rep.Score)

	// Now we have 6 entires of history, 5 entries of violations
	_, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.NotEmpty(t, err)

	violations, err := mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Empty(t, err)
	assert.Equal(t, 2, len(violations))

	violations, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Empty(t, err)
	assert.Equal(t, 5, len(violations))

	violations, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(violations))

	violations, err = mgr.GetPVDViolations("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, err)
	assert.Equal(t, 0, len(violations))

	_, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 1)
	assert.NotEmpty(t, err)

	history, err := mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 2)
	assert.Empty(t, err)
	assert.Equal(t, 2, len(history))

	history, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000000", 0, 10)
	assert.Empty(t, err)
	assert.Equal(t, 6, len(history))
	assert.Equal(t, false, history[5].Violation())

	history, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000000", 2, 10)
	assert.Empty(t, err)
	assert.Equal(t, 4, len(history))
	assert.Equal(t, false, history[3].Violation())

	history, err = mgr.GetPVDHistory("0000000000000000000000000000000000000000000000000000000000000000", 9, 10)
	assert.Empty(t, err)
	assert.Equal(t, 0, len(history))
}
