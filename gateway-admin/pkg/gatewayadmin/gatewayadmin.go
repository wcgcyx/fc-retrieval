/*
Package gatewayadmin - contains the gatewayadmin code.
*/
package gatewayadmin

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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway-admin/pkg/api/adminapi"
)

// FilecoinRetrievalGatewayAdmin is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Gateway Admin with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdmin struct {
	// admin lock
	lock sync.RWMutex
	// activeGateways is a map from nodeID to gateway's info
	activeGateways map[string]*gateway
}

// gateway represents a gateway the admin manages
type gateway struct {
	adminURL   string
	adminKey   string
	regionCode string
	alias      string
}

// NewFilecoinRetrievalGatewayAdmin initialises the Filecoin Retrieval Gateway Admin.
func NewFilecoinRetrievalGatewayAdmin() *FilecoinRetrievalGatewayAdmin {
	// Logging init
	logging.InitWithoutConfig("debug", "STDOUT", "gatewayadmin", "RFC3339")

	return &FilecoinRetrievalGatewayAdmin{
		lock:           sync.RWMutex{},
		activeGateways: make(map[string]*gateway),
	}
}

// InitialiseGateway initialises given gateway.
func (a *FilecoinRetrievalGatewayAdmin) InitialiseGateway(
	adminURL string,
	adminKey string,
	p2pPort int,
	gatewayIP string,
	rootPrivKey string,
	lotusAPIAddr string,
	lotusAuthToken string,
	registerPrivKey string,
	registerAPIAddr string,
	registerAuthToken string,
	regionCode string,
	alias string,
) error {
	// Generate Keypair
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		err = fmt.Errorf("Error in generating P2P key: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	privKeyBytes, err := privKey.Raw()
	if err != nil {
		err = fmt.Errorf("Error in getting P2P key bytes: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	p2pPrivKey := hex.EncodeToString(privKeyBytes)
	networkAddr, err := fcrserver.GetMultiAddr(p2pPrivKey, gatewayIP, uint(p2pPort))
	if err != nil {
		err = fmt.Errorf("Error in getting libp2p multiaddr: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	_, gatewayID, err := fcrcrypto.GetPublicKey(rootPrivKey)
	if err != nil {
		err = fmt.Errorf("Error in generating gateway ID: %v", err.Error())
		logging.Error(err.Error())
		return err
	}

	// Decode request
	ok, msg, err := adminapi.RequestInitiasation(adminURL, adminKey, p2pPrivKey, p2pPort, networkAddr, rootPrivKey, lotusAPIAddr, lotusAuthToken, registerPrivKey, registerAPIAddr, registerAuthToken, regionCode)
	if err != nil {
		err = fmt.Errorf("Error in decoding response: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	if !ok {
		err = fmt.Errorf("Initialisation failed with message: %v", msg)
		logging.Error(err.Error())
		return err
	}
	// Succeed
	a.lock.Lock()
	defer a.lock.Unlock()
	// Record this gateway
	a.activeGateways[gatewayID] = &gateway{
		adminURL:   adminURL,
		adminKey:   adminKey,
		regionCode: regionCode,
		alias:      alias,
	}
	return nil
}

// ListGateways lists the list of active gateways.
// It returns a slice of gateway ids and a corresponding slice of region code and alias.
func (a *FilecoinRetrievalGatewayAdmin) ListGateways() (
	[]string, // id
	[]string, // region code
	[]string, // alias
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	ids := make([]string, 0)
	region := make([]string, 0)
	alias := make([]string, 0)
	for k, v := range a.activeGateways {
		ids = append(ids, k)
		region = append(region, v.regionCode)
		alias = append(alias, v.alias)
	}
	return ids, region, alias
}

// ListPeers lists all the peers a given gateway is having a business relationship with.
func (a *FilecoinRetrievalGatewayAdmin) ListPeers(targetID string) (
	[]string, // gateway IDs
	[]int64, // score
	[]bool, // pending
	[]bool, // blocked
	[]string, // most recent activity
	[]string, // provider IDs
	[]int64, // score
	[]bool, // pending
	[]bool, // blocked
	[]string, // most recent activity
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	return adminapi.RequestListPeers(g.adminURL, g.adminKey)
}

// InspectGateway inspects a given peer gateway from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) InspectGateway(targetID string, peerID string) (
	int64, // score
	bool, // pending
	bool, // blocked
	[]string, // history
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return 0, false, false, nil, err
	}
	return adminapi.RequestInspectPeer(g.adminURL, g.adminKey, peerID, true)
}

// InspectProvider inspects a given peer provider from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) InspectProvider(targetID string, peerID string) (
	int64, // score
	bool, // pending
	bool, // blocked
	[]string, // history
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return 0, false, false, nil, err
	}
	return adminapi.RequestInspectPeer(g.adminURL, g.adminKey, peerID, false)
}

// BlockGateway blocks a given peer gateway from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) BlockGateway(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, true, true, false)
}

// UnblockGateway unblocks a given peer gateway from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) UnblockGateway(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, true, false, true)
}

// ResumeGateway resumes a given peer gateway from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) ResumeGateway(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, true, false, false)
}

// BlockProvider blocks a given peer provider from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) BlockProvider(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, false, true, false)
}

// UnblockProvider unblocks a given peer provider from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) UnblockProvider(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, false, false, true)
}

// ResumeProvider resumes a given peer provider from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) ResumeProvider(targetID string, peerID string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestChangePeerStatus(g.adminURL, g.adminKey, peerID, false, false, false)
}

// ListCIDFrequency lists the cid frequency from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) ListCIDFrequency(targetID string, page uint) (
	[]string, // cids
	[]int, // count
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return nil, nil, err
	}
	// page 0 is from 0 to 10
	// page 1 is from 10 to 20...
	return adminapi.RequestListCIDFrequency(g.adminURL, g.adminKey, 10*page, 10*(page+1))
}

// GetOfferByCID gets offers containing given cid from a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) GetOfferByCID(targetID string, cid string) (
	[]string, // digests
	[]string, // providers
	[]string, // prices
	[]int64, // expiry
	[]uint64, // qos
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}
	return adminapi.RequestGetOfferByCID(g.adminURL, g.adminKey, cid)
}

// CacheOfferByDigest caches an offer by given digest, given cid by a managed gateway
func (a *FilecoinRetrievalGatewayAdmin) CacheOfferByDigest(targetID string, digest string, cid string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Gateway %v is not in active gateways", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestCacheOfferByDigest(g.adminURL, g.adminKey, digest, cid)
}

// ForceSync forces a given managed provider to sync
func (a *FilecoinRetrievalGatewayAdmin) ForceSync(targetID string) error {
	a.lock.RLock()
	defer a.lock.RUnlock()
	g, ok := a.activeGateways[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return err
	}

	// Decode request
	ok, msg, err := adminapi.RequestForceSync(g.adminURL, g.adminKey)
	if err != nil {
		err = fmt.Errorf("Error in decoding response: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	if !ok {
		err = fmt.Errorf("Initialisation failed with message: %v", msg)
		logging.Error(err.Error())
		return err
	}
	return nil
}
