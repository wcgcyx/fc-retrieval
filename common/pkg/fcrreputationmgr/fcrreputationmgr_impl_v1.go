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
	lock sync.RWMutex

	peers map[string]*Reputation

	peerHistory    map[string]([]reputation.Record)
	peerViolations map[string]([]reputation.Record)

	pendingPeers map[string]bool
	blockedPeers map[string]bool
}

func NewFCRReputationMgrImpV1() FCRReputationMgr {
	return &FCRReputationMgrImplV1{
		lock:           sync.RWMutex{},
		peers:          make(map[string]*Reputation),
		peerHistory:    make(map[string][]reputation.Record),
		peerViolations: make(map[string][]reputation.Record),
		pendingPeers:   make(map[string]bool),
		blockedPeers:   map[string]bool{},
	}
}

func (mgr *FCRReputationMgrImplV1) Start() error {
	return nil
}

func (mgr *FCRReputationMgrImplV1) Shutdown() {
}

func (mgr *FCRReputationMgrImplV1) AddPeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	_, ok := mgr.peers[peerID]
	if ok {
		return
	}
	mgr.peers[peerID] = &Reputation{
		NodeID:  peerID,
		Score:   0,
		Pending: false,
		Blocked: false,
	}
	mgr.peerHistory[peerID] = make([]reputation.Record, 0)
	mgr.peerViolations[peerID] = make([]reputation.Record, 0)
}

func (mgr *FCRReputationMgrImplV1) ListPeers() []string {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]string, 0)
	for id := range mgr.peers {
		res = append(res, id)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) RemovePeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	delete(mgr.peers, peerID)
	delete(mgr.peerHistory, peerID)
	delete(mgr.peerViolations, peerID)
}

func (mgr *FCRReputationMgrImplV1) GetPeerReputation(peerID string) *Reputation {
	// Return a copy
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	rep, ok := mgr.peers[peerID]
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
func (mgr *FCRReputationMgrImplV1) UpdatePeerRecord(peerID string, record *reputation.Record, replica uint) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	rep, ok := mgr.peers[peerID]
	if !ok {
		return
	}
	for i := uint(0); i < replica+1; i++ {
		rep.Score += record.Point()
		mgr.peerHistory[peerID] = append([]reputation.Record{*record.Copy()}, mgr.peerHistory[peerID]...)
		if record.Violation() {
			mgr.peerViolations[peerID] = append([]reputation.Record{*record.Copy()}, mgr.peerViolations[peerID]...)
		}
	}
}

func (mgr *FCRReputationMgrImplV1) PendPeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	rep, ok := mgr.peers[peerID]
	if !ok {
		return
	}
	rep.Pending = true
	mgr.pendingPeers[peerID] = true
}

func (mgr *FCRReputationMgrImplV1) ResumePeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	rep, ok := mgr.peers[peerID]
	if !ok {
		return
	}
	rep.Pending = false
	delete(mgr.pendingPeers, peerID)
}

func (mgr *FCRReputationMgrImplV1) GetPendingPeers() []string {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.pendingPeers {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) BlockPeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	rep, ok := mgr.peers[peerID]
	if !ok {
		return
	}
	rep.Blocked = true
	mgr.blockedPeers[peerID] = true
}

func (mgr *FCRReputationMgrImplV1) UnBlockPeer(peerID string) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	rep, ok := mgr.peers[peerID]
	if !ok {
		return
	}
	rep.Blocked = false
	delete(mgr.blockedPeers, peerID)
}

func (mgr *FCRReputationMgrImplV1) GetBlockedPeers() []string {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]string, 0)
	for key := range mgr.blockedPeers {
		res = append(res, key)
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetPeerViolations(gwID string, from uint, to uint) []reputation.Record {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]reputation.Record, 0)
	violations, ok := mgr.peerViolations[gwID]
	if !ok || from > to || from > uint(len(violations)) {
		return res
	}
	for i := from; i < to && i < uint(len(violations)); i++ {
		res = append(res, *violations[i].Copy())
	}
	return res
}

func (mgr *FCRReputationMgrImplV1) GetPeerHistory(gwID string, from uint, to uint) []reputation.Record {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()
	res := make([]reputation.Record, 0)
	history, ok := mgr.peerHistory[gwID]
	if !ok || from > to || from > uint(len(history)) {
		return res
	}
	for i := from; i < to && i < uint(len(history)); i++ {
		res = append(res, *history[i].Copy())
	}
	return res
}
