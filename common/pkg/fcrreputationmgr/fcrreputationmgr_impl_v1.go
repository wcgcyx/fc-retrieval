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
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

// FCRReputationMgrImplV1 implements FCRReputationMgr, it is an in-memory version.
type FCRReputationMgrImplV1 struct {
	gwsLock  sync.RWMutex
	pvdsLock sync.RWMutex

	gws  map[string]*Reputation
	pvds map[string]*Reputation

	gwHistory     map[string]([]reputation.Record)
	pvdHistory    map[string]([]reputation.Record)
	gwViolations  map[string]([]reputation.Record)
	pvdViolations map[string]([]reputation.Record)

	pendingGWS  map[string]bool
	pendingPVDS map[string]bool
	blockedGWS  map[string]bool
	blockedPVDS map[string]bool
}

func NewFCRReputationMgrImpV1() FCRReputationMgr {
	return &FCRReputationMgrImplV1{
		gwsLock:       sync.RWMutex{},
		pvdsLock:      sync.RWMutex{},
		gws:           make(map[string]*Reputation),
		pvds:          make(map[string]*Reputation),
		gwHistory:     make(map[string]([]reputation.Record)),
		pvdHistory:    make(map[string][]reputation.Record),
		gwViolations:  make(map[string][]reputation.Record),
		pvdViolations: make(map[string][]reputation.Record),
		pendingGWS:    make(map[string]bool),
		pendingPVDS:   make(map[string]bool),
		blockedGWS:    make(map[string]bool),
		blockedPVDS:   make(map[string]bool),
	}
}

func (mgr *FCRReputationMgrImplV1) Start() error {
	return nil
}

func (mgr *FCRReputationMgrImplV1) Shutdown() {
}

func (mgr *FCRReputationMgrImplV1) AddGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	_, ok := mgr.gws[gwID]
	if ok {
		return
	}
	mgr.gws[gwID] = &Reputation{
		NodeID:  gwID,
		Score:   0,
		Pending: false,
		Blocked: false,
	}
	mgr.gwHistory[gwID] = make([]reputation.Record, 0)
	mgr.gwViolations[gwID] = make([]reputation.Record, 0)
}

func (mgr *FCRReputationMgrImplV1) ListGWS() []string {
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	res := make([]string, 0)
	for id := range mgr.gws {
		res = append(res, id)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) RemoveGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	delete(mgr.gws, gwID)
	delete(mgr.gwHistory, gwID)
	delete(mgr.gwViolations, gwID)
	delete(mgr.pendingGWS, gwID)
	delete(mgr.blockedGWS, gwID)
}

func (mgr *FCRReputationMgrImplV1) AddPVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	_, ok := mgr.pvds[pvdID]
	if ok {
		return
	}
	mgr.pvds[pvdID] = &Reputation{
		NodeID:  pvdID,
		Score:   0,
		Pending: false,
		Blocked: false,
	}
	mgr.pvdHistory[pvdID] = make([]reputation.Record, 0)
	mgr.pvdViolations[pvdID] = make([]reputation.Record, 0)
}

func (mgr *FCRReputationMgrImplV1) ListPVDS() []string {
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	res := make([]string, 0)
	for id := range mgr.pvds {
		res = append(res, id)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) RemovePVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	delete(mgr.pvds, pvdID)
	delete(mgr.pvdHistory, pvdID)
	delete(mgr.pvdViolations, pvdID)
	delete(mgr.pendingPVDS, pvdID)
	delete(mgr.blockedPVDS, pvdID)
}

func (mgr *FCRReputationMgrImplV1) GetGWReputation(gwID string) *Reputation {
	// Return a copy
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return nil
	}
	return &Reputation{
		NodeID:  rep.NodeID,
		Score:   rep.Score,
		Pending: rep.Pending,
		Blocked: rep.Blocked,
	}
}

func (mgr *FCRReputationMgrImplV1) GetPVDReputation(pvdID string) *Reputation {
	// Return a copy
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return nil
	}
	return &Reputation{
		NodeID:  rep.NodeID,
		Score:   rep.Score,
		Pending: rep.Pending,
		Blocked: rep.Blocked,
	}
}

func (mgr *FCRReputationMgrImplV1) UpdateGWRecord(gwID string, record *reputation.Record, replica uint) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return
	}
	for i := uint(0); i < replica+1; i++ {
		rep.Score += record.Point()
		mgr.gwHistory[gwID] = append([]reputation.Record{*record.Copy()}, mgr.gwHistory[gwID]...)
		if record.Violation() {
			mgr.gwViolations[gwID] = append([]reputation.Record{*record.Copy()}, mgr.gwViolations[gwID]...)
		}
	}
}

func (mgr *FCRReputationMgrImplV1) UpdatePVDRecord(pvdID string, record *reputation.Record, replica uint) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return
	}
	for i := uint(0); i < replica+1; i++ {
		rep.Score += record.Point()
		mgr.pvdHistory[pvdID] = append([]reputation.Record{*record.Copy()}, mgr.pvdHistory[pvdID]...)
		if record.Violation() {
			mgr.pvdViolations[pvdID] = append([]reputation.Record{*record.Copy()}, mgr.pvdViolations[pvdID]...)
		}
	}
}

func (mgr *FCRReputationMgrImplV1) PendGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return
	}
	rep.Pending = true
	mgr.pendingGWS[gwID] = true
}

func (mgr *FCRReputationMgrImplV1) PendPVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return
	}
	rep.Pending = true
	mgr.pendingPVDS[pvdID] = true
}

func (mgr *FCRReputationMgrImplV1) ResumeGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return
	}
	rep.Pending = false
	delete(mgr.pendingGWS, gwID)
}

func (mgr *FCRReputationMgrImplV1) ResumePVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return
	}
	rep.Pending = false
	delete(mgr.pendingPVDS, pvdID)
}

func (mgr *FCRReputationMgrImplV1) GetPendingGWS() []string {
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.pendingGWS {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetPendingPVDS() []string {
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.pendingPVDS {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) BlockGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return
	}
	rep.Blocked = true
	mgr.blockedGWS[gwID] = true
}

func (mgr *FCRReputationMgrImplV1) BlockPVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return
	}
	rep.Blocked = true
	mgr.blockedPVDS[pvdID] = true
}

func (mgr *FCRReputationMgrImplV1) UnBlockGW(gwID string) {
	mgr.gwsLock.Lock()
	defer mgr.gwsLock.Unlock()
	rep, ok := mgr.gws[gwID]
	if !ok {
		return
	}
	rep.Blocked = false
	delete(mgr.blockedGWS, gwID)
}

func (mgr *FCRReputationMgrImplV1) UnBlockPVD(pvdID string) {
	mgr.pvdsLock.Lock()
	defer mgr.pvdsLock.Unlock()
	rep, ok := mgr.pvds[pvdID]
	if !ok {
		return
	}
	rep.Blocked = false
	delete(mgr.blockedPVDS, pvdID)
}

func (mgr *FCRReputationMgrImplV1) GetBlockedGWS() []string {
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.blockedGWS {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetBlockedPVDS() []string {
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.blockedPVDS {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetGWViolations(gwID string, from uint, to uint) []reputation.Record {
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	res := make([]reputation.Record, 0)
	violations, ok := mgr.gwViolations[gwID]
	if !ok || from > to || from > uint(len(violations)) {
		return res
	}
	for i := from; i < to && i < uint(len(violations)); i++ {
		res = append(res, *violations[i].Copy())
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetPVDViolations(pvdID string, from uint, to uint) []reputation.Record {
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	res := make([]reputation.Record, 0)
	violations, ok := mgr.pvdViolations[pvdID]
	if !ok || from > to || from > uint(len(violations)) {
		return res
	}
	for i := from; i < to && i < uint(len(violations)); i++ {
		res = append(res, *violations[i].Copy())
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetGWHistory(gwID string, from uint, to uint) []reputation.Record {
	mgr.gwsLock.RLock()
	defer mgr.gwsLock.RUnlock()
	res := make([]reputation.Record, 0)
	history, ok := mgr.gwHistory[gwID]
	if !ok || from > to || from > uint(len(history)) {
		return res
	}
	for i := from; i < to && i < uint(len(history)); i++ {
		res = append(res, *history[i].Copy())
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetPVDHistory(pvdID string, from uint, to uint) []reputation.Record {
	mgr.pvdsLock.RLock()
	defer mgr.pvdsLock.RUnlock()
	res := make([]reputation.Record, 0)
	history, ok := mgr.pvdHistory[pvdID]
	if !ok || from > to || from > uint(len(history)) {
		return res
	}
	for i := from; i < to && i < uint(len(history)); i++ {
		res = append(res, *history[i].Copy())
	}
	return res
}
