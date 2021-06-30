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

import (
	"errors"
	"sync"
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/dhtring"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// FCRPeerMgrImplV1 implements FCRPeerMgr, it is an in-memory version.
type FCRPeerMgrImplV1 struct {
	// Boolean indicates if the manager has started
	start bool

	// Register manager
	registerMgr fcrregistermgr.FCRRegisterMgr

	// Duration to wait between two updates
	refreshDuration time.Duration

	// Boolean indicates if to discover gateway/provider
	gatewayDiscv  bool
	providerDiscv bool

	// trackCIDRange indicates if track current cid range
	trackCIDRange bool

	// Channels to control the threads
	gatewayShutdownCh  chan bool
	providerShutdownCh chan bool
	gatewayRefreshCh   chan bool
	providerRefreshCh  chan bool

	discoveredGWS     map[string]*Peer
	discoveredGWSLock sync.RWMutex

	discoveredPVDS     map[string]*Peer
	discoveredPVDSLock sync.RWMutex

	// closestGateways stores the mapping from gateway closest for DHT network sorted clockwise
	closestGatewaysIDs     *dhtring.Ring
	closestGatewaysIDsLock sync.RWMutex

	// The following fields apply only when tracking cid hash range.
	anchor    string
	hashMin   string
	hashMax   string
	rangeLock sync.RWMutex
}

func NewFCRPeerMgrImplV1(registerMgr fcrregistermgr.FCRRegisterMgr, gatewayDiscv bool, providerDiscv bool, trackCIDRange bool, trackAnchor string, refreshDuration time.Duration) FCRPeerMgr {
	return &FCRPeerMgrImplV1{
		start:                  false,
		registerMgr:            registerMgr,
		refreshDuration:        refreshDuration,
		gatewayDiscv:           gatewayDiscv,
		providerDiscv:          providerDiscv,
		trackCIDRange:          trackCIDRange,
		gatewayShutdownCh:      make(chan bool),
		providerShutdownCh:     make(chan bool),
		gatewayRefreshCh:       make(chan bool),
		providerRefreshCh:      make(chan bool),
		discoveredGWS:          make(map[string]*Peer),
		discoveredGWSLock:      sync.RWMutex{},
		discoveredPVDS:         make(map[string]*Peer),
		discoveredPVDSLock:     sync.RWMutex{},
		closestGatewaysIDs:     dhtring.CreateRing(),
		closestGatewaysIDsLock: sync.RWMutex{},
		anchor:                 trackAnchor,
		hashMin:                "0000000000000000000000000000000000000000000000000000000000000000",
		hashMax:                "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		rangeLock:              sync.RWMutex{},
	}
}

func (mgr *FCRPeerMgrImplV1) Start() error {
	if mgr.start {
		return errors.New("FCRPeerManager has already started")
	}
	mgr.start = true
	if mgr.gatewayDiscv {
		go mgr.gwSyncRoutine()
	}
	if mgr.providerDiscv {
		go mgr.pvdSyncRoutine()
	}
	return nil
}

func (mgr *FCRPeerMgrImplV1) Shutdown() {
	if !mgr.start {
		return
	}
	if mgr.gatewayDiscv {
		mgr.gatewayShutdownCh <- true
		<-mgr.gatewayShutdownCh
	}
	if mgr.providerDiscv {
		mgr.providerShutdownCh <- true
		<-mgr.providerShutdownCh
	}
	mgr.start = false
}

func (mgr *FCRPeerMgrImplV1) Sync() {
	if !mgr.start {
		return
	}
	if mgr.gatewayDiscv {
		mgr.gatewayRefreshCh <- true
		<-mgr.gatewayRefreshCh
	}
	if mgr.providerDiscv {
		mgr.providerRefreshCh <- true
		<-mgr.providerRefreshCh
	}
}

func (mgr *FCRPeerMgrImplV1) SyncGW(gwID string) {
	if !mgr.start {
		return
	}
	gwReg, err := mgr.registerMgr.GetRegisteredGatewayByID(gwID)
	mgr.discoveredGWSLock.Lock()
	defer mgr.discoveredGWSLock.Unlock()
	mgr.closestGatewaysIDsLock.Lock()
	defer mgr.closestGatewaysIDsLock.Unlock()
	if err != nil {
		delete(mgr.discoveredGWS, gwID)
		if mgr.gatewayDiscv {
			mgr.closestGatewaysIDs.Remove(gwID)
		}
		return
	}
	// Check if there is an existing entry
	gwPeer, ok := mgr.discoveredGWS[gwID]
	if !ok {
		gwPeer = &Peer{
			NodeID: gwID,
		}
		mgr.discoveredGWS[gwID] = gwPeer
		mgr.closestGatewaysIDs.Insert(gwID)
	}
	// Mostly used to updating msg key
	gwPeer.RootKey = gwReg.RootKey
	gwPeer.MsgSigningKey = gwReg.MsgSigningKey
	gwPeer.MsgSigningKeyVer = gwReg.MsgSigningKeyVer
	gwPeer.RegionCode = gwReg.RegionCode
	gwPeer.NetworkAddr = gwReg.NetworkAddr
	gwPeer.Deregistering = gwReg.Deregistering
	gwPeer.DeregisteringHeight = gwReg.DeregisteringHeight
}

func (mgr *FCRPeerMgrImplV1) SyncPVD(pvdID string) {
	if !mgr.start {
		return
	}
	pvdReg, err := mgr.registerMgr.GetRegisteredProviderByID(pvdID)
	mgr.discoveredPVDSLock.Lock()
	defer mgr.discoveredPVDSLock.Unlock()
	if err != nil {
		delete(mgr.discoveredPVDS, pvdID)
		return
	}
	// Check if there is an existing entry
	pvdPeer, ok := mgr.discoveredPVDS[pvdID]
	if !ok {
		pvdPeer = &Peer{
			NodeID: pvdID,
		}
		mgr.discoveredPVDS[pvdID] = pvdPeer
	}
	// Mostly used to updating msg key
	pvdPeer.RootKey = pvdReg.RootKey
	pvdPeer.MsgSigningKey = pvdReg.MsgSigningKey
	pvdPeer.MsgSigningKeyVer = pvdReg.MsgSigningKeyVer
	pvdPeer.OfferSigningKey = pvdReg.OfferSigningKey
	pvdPeer.RegionCode = pvdReg.RegionCode
	pvdPeer.NetworkAddr = pvdReg.NetworkAddr
	pvdPeer.Deregistering = pvdReg.Deregistering
	pvdPeer.DeregisteringHeight = pvdReg.DeregisteringHeight
}

func (mgr *FCRPeerMgrImplV1) GetGWInfo(gwID string) (*Peer, error) {
	if !mgr.start {
		return nil, errors.New("FCRPeerManager has not been started")
	}
	// Return a copy
	mgr.discoveredGWSLock.RLock()
	defer mgr.discoveredGWSLock.RUnlock()
	peer, ok := mgr.discoveredGWS[gwID]
	if !ok {
		return nil, errors.New("Gateway not found locally")
	}
	return &Peer{
		RootKey:             peer.RootKey,
		NodeID:              peer.NodeID,
		MsgSigningKey:       peer.MsgSigningKey,
		MsgSigningKeyVer:    peer.MsgSigningKeyVer,
		RegionCode:          peer.RegionCode,
		NetworkAddr:         peer.NetworkAddr,
		Deregistering:       peer.Deregistering,
		DeregisteringHeight: peer.DeregisteringHeight,
	}, nil
}

func (mgr *FCRPeerMgrImplV1) GetPVDInfo(pvdID string) (*Peer, error) {
	if !mgr.start {
		return nil, errors.New("FCRPeerManager has not been started")
	}
	// Return a copy
	mgr.discoveredPVDSLock.RLock()
	defer mgr.discoveredPVDSLock.RUnlock()
	peer, ok := mgr.discoveredPVDS[pvdID]
	if !ok {
		return nil, errors.New("Provider not found locally")
	}
	return &Peer{
		RootKey:             peer.RootKey,
		NodeID:              peer.NodeID,
		MsgSigningKey:       peer.MsgSigningKey,
		MsgSigningKeyVer:    peer.MsgSigningKeyVer,
		OfferSigningKey:     peer.OfferSigningKey,
		RegionCode:          peer.RegionCode,
		NetworkAddr:         peer.NetworkAddr,
		Deregistering:       peer.Deregistering,
		DeregisteringHeight: peer.DeregisteringHeight,
	}, nil
}

func (mgr *FCRPeerMgrImplV1) GetGWSNearCIDHash(hash string, except string) ([]Peer, error) {
	if !mgr.start {
		return nil, errors.New("FCRPeerManager has not been started")
	}
	if !mgr.gatewayDiscv {
		return nil, errors.New("FCRPeerManager does not sync gateways")
	}
	mgr.discoveredGWSLock.RLock()
	defer mgr.discoveredGWSLock.RUnlock()
	mgr.closestGatewaysIDsLock.RLock()
	defer mgr.closestGatewaysIDsLock.RUnlock()
	ids := mgr.closestGatewaysIDs.GetClosest(hash, 16, except)
	// return copies
	res := make([]Peer, 0)
	for _, id := range ids {
		peer := mgr.discoveredGWS[id]
		res = append(res, Peer{
			RootKey:             peer.RootKey,
			NodeID:              peer.NodeID,
			MsgSigningKey:       peer.MsgSigningKey,
			MsgSigningKeyVer:    peer.MsgSigningKeyVer,
			RegionCode:          peer.RegionCode,
			NetworkAddr:         peer.NetworkAddr,
			Deregistering:       peer.Deregistering,
			DeregisteringHeight: peer.DeregisteringHeight,
		})
	}
	return res, nil
}

func (mgr *FCRPeerMgrImplV1) GetCurrentCIDHashRange() (string, string, error) {
	if !mgr.start {
		return "", "", errors.New("FCRPeerManager has not been started")
	}
	if !mgr.gatewayDiscv || !mgr.trackCIDRange {
		return "", "", errors.New("FCRPeerManager does not sync gateways/does not track range")
	}
	mgr.rangeLock.RLock()
	defer mgr.rangeLock.RUnlock()
	return mgr.hashMin, mgr.hashMax, nil
}

func (mgr *FCRPeerMgrImplV1) gwSyncRoutine() {
	refreshForce := false
	for {
		afterChan := time.After(mgr.refreshDuration)
		select {
		case <-mgr.gatewayRefreshCh:
			// Need to refresh
			logging.Info("FCRPeerManager force sync gateways.")
			refreshForce = true
		case <-afterChan:
			// Need to refresh
		case <-mgr.gatewayShutdownCh:
			// Need to shutdown
			logging.Info("FCRPeerManager shutdown gateway syncing routine.")
			mgr.gatewayShutdownCh <- true
			return
		}

		// Get current height
		height, err := mgr.registerMgr.GetHeight()
		if err != nil {
			logging.Warn("FCRPeerManager gateway sync fail to get current height: %v", err.Error())
			continue
		}
		toRemove := make(map[string]bool)
		mgr.discoveredGWSLock.RLock()
		for key := range mgr.discoveredGWS {
			toRemove[key] = true
		}
		mgr.discoveredGWSLock.RUnlock()
		maxPage, err := mgr.registerMgr.GetGWMaxPage(height)
		if err != nil {
			logging.Warn("FCRPeerManager gateway sync fail to get max page at height %v: %v", height, err.Error())
			continue
		}
		refreshRange := false
		for page := uint64(0); page <= maxPage; page++ {
			gwInfos, err := mgr.registerMgr.GetAllRegisteredGateway(height, page)
			if err != nil {
				logging.Warn("FCRPeerManager gateway sync fail to get registered gateways at page %v at height %v: %v. Try again", page, height, err.Error())
				page--
				continue
			}
			for _, gwInfo := range gwInfos {
				delete(toRemove, gwInfo.NodeID)
				update := false
				mgr.discoveredGWSLock.RLock()
				storedInfo, ok := mgr.discoveredGWS[gwInfo.NodeID]
				if !ok {
					// Not exist, we need to add a new entry
					mgr.closestGatewaysIDsLock.Lock()
					mgr.closestGatewaysIDs.Insert(gwInfo.NodeID)
					mgr.closestGatewaysIDsLock.Unlock()
					refreshRange = true
					update = true
				} else {
					if storedInfo.MsgSigningKey != gwInfo.MsgSigningKey ||
						storedInfo.MsgSigningKeyVer != gwInfo.MsgSigningKeyVer ||
						storedInfo.Deregistering != gwInfo.Deregistering ||
						storedInfo.DeregisteringHeight != gwInfo.DeregisteringHeight {
						update = true
					}
				}
				mgr.discoveredGWSLock.RUnlock()
				if update {
					mgr.discoveredGWSLock.Lock()
					mgr.discoveredGWS[gwInfo.NodeID] = &Peer{
						RootKey:             gwInfo.RootKey,
						NodeID:              gwInfo.NodeID,
						MsgSigningKey:       gwInfo.MsgSigningKey,
						MsgSigningKeyVer:    gwInfo.MsgSigningKeyVer,
						RegionCode:          gwInfo.RegionCode,
						NetworkAddr:         gwInfo.NetworkAddr,
						Deregistering:       gwInfo.Deregistering,
						DeregisteringHeight: gwInfo.DeregisteringHeight,
					}
					mgr.discoveredGWSLock.Unlock()
				}
			}
		}
		for key := range toRemove {
			refreshRange = true
			mgr.discoveredGWSLock.Lock()
			mgr.closestGatewaysIDsLock.Lock()
			delete(mgr.discoveredGWS, key)
			mgr.closestGatewaysIDs.Remove(key)
			mgr.closestGatewaysIDsLock.Unlock()
			mgr.discoveredGWSLock.Unlock()
		}
		if refreshRange {
			mgr.updateCIDHashRange()
		}
		if refreshForce {
			mgr.gatewayRefreshCh <- true
			refreshForce = false
		}
	}
}

func (mgr *FCRPeerMgrImplV1) pvdSyncRoutine() {
	refreshForce := false
	for {
		afterChan := time.After(mgr.refreshDuration)
		select {
		case <-mgr.providerRefreshCh:
			// Need to refresh
			logging.Info("FCRPeerManager force sync providers.")
			refreshForce = true
		case <-afterChan:
			// Need to refresh
		case <-mgr.providerShutdownCh:
			// Need to shutdown
			logging.Info("FCRPeerManager shutdown provider syncing routine.")
			mgr.providerShutdownCh <- true
			return
		}

		// Get current height
		height, err := mgr.registerMgr.GetHeight()
		if err != nil {
			logging.Warn("FCRPeerManager provider sync fail to get current height: %v", err.Error())
			continue
		}
		toRemove := make(map[string]bool)
		mgr.discoveredPVDSLock.RLock()
		for key := range mgr.discoveredPVDS {
			toRemove[key] = true
		}
		mgr.discoveredPVDSLock.RUnlock()
		maxPage, err := mgr.registerMgr.GetPVDMaxPage(height)
		if err != nil {
			logging.Warn("FCRPeerManager provider sync fail to get max page at height %v: %v", height, err.Error())
			continue
		}
		for page := uint64(0); page <= maxPage; page++ {
			pvdInfos, err := mgr.registerMgr.GetAllRegisteredProvider(height, page)
			if err != nil {
				logging.Warn("FCRPeerManager provider sync fail to get registered providers at page %v at height %v: %v. Try again", page, height, err.Error())
				page--
				continue
			}
			for _, pvdInfo := range pvdInfos {
				delete(toRemove, pvdInfo.NodeID)
				update := false
				mgr.discoveredPVDSLock.RLock()
				storedInfo, ok := mgr.discoveredPVDS[pvdInfo.NodeID]
				if !ok {
					// Not exist, we need to add a new entry
					update = true
				} else {
					if storedInfo.MsgSigningKey != pvdInfo.MsgSigningKey ||
						storedInfo.MsgSigningKeyVer != pvdInfo.MsgSigningKeyVer ||
						storedInfo.OfferSigningKey != pvdInfo.OfferSigningKey ||
						storedInfo.Deregistering != pvdInfo.Deregistering ||
						storedInfo.DeregisteringHeight != pvdInfo.DeregisteringHeight {
						update = true
					}
				}
				mgr.discoveredPVDSLock.RUnlock()
				if update {
					mgr.discoveredPVDSLock.Lock()
					mgr.discoveredPVDS[pvdInfo.NodeID] = &Peer{
						RootKey:             pvdInfo.RootKey,
						NodeID:              pvdInfo.NodeID,
						MsgSigningKey:       pvdInfo.MsgSigningKey,
						MsgSigningKeyVer:    pvdInfo.MsgSigningKeyVer,
						OfferSigningKey:     pvdInfo.OfferSigningKey,
						RegionCode:          pvdInfo.RegionCode,
						NetworkAddr:         pvdInfo.NetworkAddr,
						Deregistering:       pvdInfo.Deregistering,
						DeregisteringHeight: pvdInfo.DeregisteringHeight,
					}
					mgr.discoveredPVDSLock.Unlock()
				}
			}
		}
		for key := range toRemove {
			mgr.discoveredPVDSLock.Lock()
			delete(mgr.discoveredPVDS, key)
			mgr.discoveredPVDSLock.Unlock()
		}
		if refreshForce {
			mgr.providerRefreshCh <- true
			refreshForce = false
		}
	}
}

func (mgr *FCRPeerMgrImplV1) updateCIDHashRange() {
	mgr.closestGatewaysIDsLock.RLock()
	defer mgr.closestGatewaysIDsLock.RUnlock()
	mgr.rangeLock.Lock()
	defer mgr.rangeLock.Unlock()
	res := mgr.closestGatewaysIDs.GetClosest(mgr.anchor, 16, mgr.anchor)
	if len(res) < 16 {
		return
	}
	mgr.hashMin = res[0]
	mgr.hashMax = res[15]
}
