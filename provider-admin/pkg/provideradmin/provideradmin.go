/*
Package provideradmin - contains the provideradmin code.
*/
package provideradmin

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
	"math/big"
	"sync"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider-admin/pkg/api/adminapi"
)

type FilecoinRetrievalProviderAdmin struct {
	// admin lock
	lock sync.RWMutex
	// activeProviders is a map from nodeID to provider's info
	activeProviders map[string]*provider
}

// provider represents a provider the admin manages
type provider struct {
	adminURL   string
	adminKey   string
	regionCode string
	alias      string
}

// NewFilecoinRetrievalProviderAdmin initialises the Filecoin Retrieval Provider Admin.
func NewFilecoinRetrievalProviderAdmin() *FilecoinRetrievalProviderAdmin {
	// Logging init
	logging.InitWithoutConfig("debug", "STDOUT", "provideradmin", "RFC3339")

	return &FilecoinRetrievalProviderAdmin{
		lock:            sync.RWMutex{},
		activeProviders: make(map[string]*provider),
	}
}

// InitialiseProvider initialises given provider.
func (a *FilecoinRetrievalProviderAdmin) InitialiseProvider(
	adminURL string,
	adminKey string,
	p2pPort int,
	providerIP string,
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
	networkAddr, err := fcrserver.GetMultiAddr(p2pPrivKey, providerIP, uint(p2pPort))
	if err != nil {
		err = fmt.Errorf("Error in getting libp2p multiaddr: %v", err.Error())
		logging.Error(err.Error())
		return err
	}
	_, providerID, err := fcrcrypto.GetPublicKey(rootPrivKey)
	if err != nil {
		err = fmt.Errorf("Error in generating provider ID: %v", err.Error())
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
	// Record this provider
	a.activeProviders[providerID] = &provider{
		adminURL:   adminURL,
		adminKey:   adminKey,
		regionCode: regionCode,
		alias:      alias,
	}
	return nil
}

// ListProviders lists the list of active providers.
// It returns a slice of provider ids and a corresponding slice of region code and alias.
func (a *FilecoinRetrievalProviderAdmin) ListProviders() (
	[]string, // id
	[]string, // region code
	[]string, // alias
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	ids := make([]string, 0)
	region := make([]string, 0)
	alias := make([]string, 0)
	for k, v := range a.activeProviders {
		ids = append(ids, k)
		region = append(region, v.regionCode)
		alias = append(alias, v.alias)
	}
	return ids, region, alias
}

// ListFiles lists all the files a given provider is monitoring.
func (a *FilecoinRetrievalProviderAdmin) ListFiles(targetID string) (
	[]string, // files
	[]string, // cids
	[]int64, // sizes
	[]bool, // published
	[]int, // frequency
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}
	return adminapi.RequestListFiles(p.adminURL, p.adminKey)
}

// GetOfferByCID gets offers containing given cid from a managed provider
func (a *FilecoinRetrievalProviderAdmin) GetOfferByCID(targetID string, cid string) (
	[]string, // digests
	[]string, // providers
	[]string, // prices
	[]int64, // expiry
	[]uint64, // qos
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}
	return adminapi.RequestGetOfferByCID(p.adminURL, p.adminKey, cid)
}

// UploadFile uploads a file to given managed provider
func (a *FilecoinRetrievalProviderAdmin) UploadFile(targetID string, filename string, tag string) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestUploadFile(p.adminURL, p.adminKey, filename, tag)
}

// PublishOffer asks given managed provider to publish an offer
func (a *FilecoinRetrievalProviderAdmin) PublishOffer(targetID string, files []string, price *big.Int, expiry int64, qos uint64) (
	bool, // Success
	string, // Information
	error, // error
) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return false, "", err
	}
	return adminapi.RequestPublishOffer(p.adminURL, p.adminKey, files, price, expiry, qos)
}

// FastPublishOffer uploads a file to given provider then asks it to publish the offer
func (a *FilecoinRetrievalProviderAdmin) FastPublishOffer(targetID string, filename string, tag string, price *big.Int, expiry int64, qos uint64) error {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return err
	}
	ok, msg, err := adminapi.RequestUploadFile(p.adminURL, p.adminKey, filename, tag)
	if err != nil {
		err = fmt.Errorf("Error in uploading file %v to %v: %v", filename, targetID, err.Error())
		logging.Error(err.Error())
		return err
	}
	if !ok {
		err = fmt.Errorf("Fail to upload file %v to %v: %v", filename, targetID, msg)
		logging.Error(err.Error())
		return err
	}
	// Then ask it to publish the offer
	ok, msg, err = adminapi.RequestPublishOffer(p.adminURL, p.adminKey, []string{tag}, price, expiry, qos)
	if err != nil {
		err = fmt.Errorf("Error in publishing offer for %v: %v", tag, err.Error())
		logging.Error(err.Error())
		return err
	}
	if !ok {
		err = fmt.Errorf("Fail to publish offer for %v: %v", tag, msg)
		logging.Error(err.Error())
		return err
	}
	return nil
}

// ForceSync forces a given managed provider to sync
func (a *FilecoinRetrievalProviderAdmin) ForceSync(targetID string) error {
	a.lock.RLock()
	defer a.lock.RUnlock()
	p, ok := a.activeProviders[targetID]
	if !ok {
		err := fmt.Errorf("Provider %v is not in active providers", targetID)
		logging.Error(err.Error())
		return err
	}

	// Decode request
	ok, msg, err := adminapi.RequestForceSync(p.adminURL, p.adminKey)
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
