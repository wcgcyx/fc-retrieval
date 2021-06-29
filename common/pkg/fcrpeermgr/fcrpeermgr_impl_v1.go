/*
Package fcrpeermgr - peer manager manages all retrieval peers.
*/
package fcrpeermgr

import (
	"sync"
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/dhtring"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrregistermgr"
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

// FCRPeerMgrImplV1 implements FCRPeerMgr, it is an in-memory version.
type FCRPeerMgrImplV1 struct {
	registerMgr *fcrregistermgr.FCRRegisterMgr

	// Duration to wait between two updates
	refreshDuration time.Duration

	// Boolean indicates if to discover gateway/provider
	gatewayDiscv  bool
	providerDiscv bool

	// Channels to control the threads
	gatewayShutdownCh  chan bool
	providerShutdownCh chan bool
	gatewayRefreshCh   chan bool
	providerRefreshCh  chan bool

	discoveredGWS     map[string]Peer
	discoveredGWSLock sync.RWMutex

	// closestGateways stores the mapping from gateway closest for DHT network sorted clockwise
	closestGatewaysIDs     *dhtring.Ring
	closestGatewaysIDsLock sync.RWMutex

	discoveredPVDS     map[string]Peer
	discoveredPVDSLock sync.RWMutex
}

func NewFCRPeerMgrImplV1(registerMgr *fcrregistermgr.FCRRegisterMgr, gatewayDiscv bool, providerDiscv bool, refreshDuration time.Duration) FCRPeerMgr {
	return nil
}
