/*
Package client - contains the client code.
*/
package client

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
	"net/http"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/wcgcyx/fc-retrieval/client/pkg/api/p2papi"
	"github.com/wcgcyx/fc-retrieval/client/pkg/core"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrlotusmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcroffermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpeermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrreputationmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// FilecoinRetrievalClient is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Client with Filecoin Retrieval Gateways/Providers.
type FilecoinRetrievalClient struct {
	core *core.Core
}

// NewFilecoinRetrievalClient initialise the Filecoin Retrieval Client
func NewFilecoinRetrievalClient(
	walletPrvKey string,
	lotusAPIAddr string,
	lotusAuthToken string,
	registerPrvKey string,
	registerAPIAddr string,
	registerAuthToken string,
) (*FilecoinRetrievalClient, error) {
	// Logging init
	logging.InitWithoutConfig("debug", "STDOUT", "client", "RFC3339")

	// Initialise client
	c := core.GetSingleInstance()

	// Generating msg signing key
	msgKey, _, _, err := fcrcrypto.GenerateRetrievalKeyPair()
	if err != nil {
		err = fmt.Errorf("Error in generating msg signing key: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	c.MsgKey = msgKey
	_, nodeID, err := fcrcrypto.GetPublicKey(msgKey)
	if err != nil {
		err = fmt.Errorf("Error in generating nodeID from msg signing key %v: %v", msgKey, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	c.NodeID = nodeID

	// Initialise P2P Server
	// Generate Keypair
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		err = fmt.Errorf("Error in generating P2P key: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	prvKeyBytes, err := prvKey.Raw()
	if err != nil {
		err = fmt.Errorf("Error in getting P2P key bytes: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	// Initialise components
	c.P2PServer = fcrserver.NewFCRServerImplV1(hex.EncodeToString(prvKeyBytes), 0, time.Second*60)
	c.P2PServer.
		AddRequester(fcrmessages.EstablishmentRequestType, p2papi.EstablishmentRequester).
		AddRequester(fcrmessages.StandardOfferDiscoveryRequestType, p2papi.OfferQueryRequester).
		AddRequester(fcrmessages.DHTOfferDiscoveryRequestType, p2papi.DHTOfferQueryRequester)
	err = c.P2PServer.Start()
	if err != nil {
		err = fmt.Errorf("Error in starting P2P server: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	c.ReputationMgr = fcrreputationmgr.NewFCRReputationMgrImpV1()
	err = c.ReputationMgr.Start()
	if err != nil {
		err = fmt.Errorf("Error in starting reputation manager: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	c.RegisterMgr = fcrregistermgr.NewFCRRegisterMgrImplV1(registerAPIAddr, &http.Client{Timeout: 180 * time.Second})
	c.PeerMgr = fcrpeermgr.NewFCRPeerMgrImplV1(c.RegisterMgr, c.ReputationMgr, false, false, false, nodeID, time.Hour)
	err = c.PeerMgr.Start()
	if err != nil {
		err = fmt.Errorf("Error in starting peer manager: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	lotusMgr := fcrlotusmgr.NewFCRLotusMgrImplV1(lotusAPIAddr, lotusAuthToken, nil)
	c.PaymentMgr = fcrpaymentmgr.NewFCRPaymentMgrImplV1(walletPrvKey, lotusMgr)
	err = c.PaymentMgr.Start()
	if err != nil {
		err = fmt.Errorf("Error in starting payment manager: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	c.OfferMgr = fcroffermgr.NewFCROfferMgrImplV1(true)
	err = c.OfferMgr.Start()
	if err != nil {
		err = fmt.Errorf("Error in starting offer manager: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// At start-up, updating all active gateways and providers
	for _, gwID := range c.ReputationMgr.ListGWS() {
		c.PeerMgr.SyncGW(gwID)
	}
	for _, pvdID := range c.ReputationMgr.ListPVDS() {
		c.PeerMgr.SyncPVD(pvdID)
	}

	return &FilecoinRetrievalClient{
		core: c,
	}, nil
}

// Search searches gateways that are in given location.
func (c *FilecoinRetrievalClient) Search(location string) ([]string, error) {
	// TODO, Search by location
	res := make([]string, 0)
	infos, err := c.core.RegisterMgr.GetAllRegisteredGateway(0, 0)
	if err != nil {
		err = fmt.Errorf("Error in getting all registered gateways: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	for _, info := range infos {
		res = append(res, info.NodeID)
	}
	return res, nil
}

// AddActiveGW adds an active gateway ID
func (c *FilecoinRetrievalClient) AddActiveGW(targetID string) error {
	if c.core.ReputationMgr.GetGWReputation(targetID) != nil {
		err := fmt.Errorf("Gateway %v is already active", targetID)
		logging.Error(err.Error())
		return err
	}

	// Get gw info
	gwInfo := c.core.PeerMgr.GetGWInfo(targetID)
	if gwInfo == nil {
		// Not found, try sync once
		gwInfo = c.core.PeerMgr.SyncGW(targetID)
		if gwInfo == nil {
			err := fmt.Errorf("Error in obtaining information for gateway %v", targetID)
			logging.Error(err.Error())
			return err
		}
	}
	_, err := c.core.P2PServer.Request(gwInfo.NetworkAddr, fcrmessages.EstablishmentRequestType, targetID)
	if err != nil {
		err = fmt.Errorf("Error in sending establishment request to %v with addr %v: %v", targetID, gwInfo.NetworkAddr, err.Error())
		logging.Error(err.Error())
		return err
	}
	// Create payment channel
	recipientAddr, err := fcrcrypto.GetWalletAddress(gwInfo.RootKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining wallet addreess for gateway %v with root key %v: %v", targetID, gwInfo.RootKey, err.Error())
		logging.Error(err.Error())
		return err
	}

	err = c.core.PaymentMgr.Create(recipientAddr, c.core.TopupAmount)
	if err != nil {
		err = fmt.Errorf("Error in creating a payment channel to %v with wallet address %v with topup amount of %v: %v", targetID, recipientAddr, c.core.TopupAmount.String(), err.Error())
		logging.Error(err.Error())
		return err
	}
	// Add gateway entry to reputation
	c.core.ReputationMgr.AddGW(gwInfo.NodeID)
	return nil
}

// ListActiveGWS lists all active gateways
func (c *FilecoinRetrievalClient) ListActiveGWS() ([]string, error) {
	return c.core.ReputationMgr.ListGWS(), nil
}

// AddActivePVD adds an active provider ID
func (c *FilecoinRetrievalClient) AddActivePVD(targetID string) error {
	if c.core.ReputationMgr.GetPVDReputation(targetID) != nil {
		err := fmt.Errorf("Gateway %v is already active", targetID)
		logging.Error(err.Error())
		return err
	}

	// Get pvd info
	pvdInfo := c.core.PeerMgr.GetPVDInfo(targetID)
	if pvdInfo == nil {
		// Not found, try sync once
		pvdInfo = c.core.PeerMgr.SyncPVD(targetID)
		if pvdInfo == nil {
			err := fmt.Errorf("Error in obtaining information for gateway %v", targetID)
			logging.Error(err.Error())
			return err
		}
	}
	_, err := c.core.P2PServer.Request(pvdInfo.NetworkAddr, fcrmessages.EstablishmentRequestType, targetID)
	if err != nil {
		err = fmt.Errorf("Error in sending establishment request to %v with addr %v: %v", targetID, pvdInfo.NetworkAddr, err.Error())
		logging.Error(err.Error())
		return err
	}
	// Create payment channel
	recipientAddr, err := fcrcrypto.GetWalletAddress(pvdInfo.RootKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining wallet addreess for gateway %v with root key %v: %v", targetID, pvdInfo.RootKey, err.Error())
		logging.Error(err.Error())
		return err
	}

	err = c.core.PaymentMgr.Create(recipientAddr, c.core.TopupAmount)
	if err != nil {
		err = fmt.Errorf("Error in creating a payment channel to %v with wallet address %v with topup amount of %v: %v", targetID, recipientAddr, c.core.TopupAmount.String(), err.Error())
		logging.Error(err.Error())
		return err
	}
	// Add provider entry to reputation
	c.core.ReputationMgr.AddPVD(pvdInfo.NodeID)
	return nil
}

// ListActivePVDS lists all active providers
func (c *FilecoinRetrievalClient) ListActivePVDS() ([]string, error) {
	return c.core.ReputationMgr.ListPVDS(), nil
}

// StandardDiscovery performs a standard discovery.
func (c *FilecoinRetrievalClient) StandardDiscovery(cidStr string, toContact map[string]uint32) ([]cidoffer.SubCIDOffer, error) {
	pieceCID, err := cid.NewContentID(cidStr)
	if err != nil {
		err = fmt.Errorf("Error in decoding cid: %v: %v", cidStr, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	temp := make(map[string]*cidoffer.SubCIDOffer, 0)
	// TODO, Concurrency
	for targetID, maxOfferRequested := range toContact {
		// Get gw info
		gwInfo := c.core.PeerMgr.GetGWInfo(targetID)
		if gwInfo == nil {
			// Not found, try sync once
			gwInfo = c.core.PeerMgr.SyncGW(targetID)
			if gwInfo == nil {
				logging.Error("Error in obtaining information for gateway %v", targetID)
				continue
			}
		}
		response, err := c.core.P2PServer.Request(gwInfo.NetworkAddr, fcrmessages.StandardOfferDiscoveryRequestType, targetID, pieceCID, maxOfferRequested)
		if err != nil {
			logging.Error("Error in requesting gateway %v for offers: %v", targetID, err.Error())
			continue
		}
		_, offers, _, _ := fcrmessages.DecodeStandardOfferDiscoveryResponse(response)
		for _, offer := range offers {
			temp[offer.GetMessageDigest()] = &offer
		}
	}
	res := make([]cidoffer.SubCIDOffer, 0)
	for _, offer := range temp {
		res = append(res, *offer)
	}
	return res, nil
}

// DHTDiscovery performs a DHT discovery.
func (c *FilecoinRetrievalClient) DHTDiscovery(cidStr string, targetID string, numDHT uint32, maxOfferRequestedPerDHT uint32) ([]cidoffer.SubCIDOffer, error) {
	pieceCID, err := cid.NewContentID(cidStr)
	if err != nil {
		err = fmt.Errorf("Error in decoding cid: %v: %v", cidStr, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	// Get gw info
	gwInfo := c.core.PeerMgr.GetGWInfo(targetID)
	if gwInfo == nil {
		// Not found, try sync once
		gwInfo = c.core.PeerMgr.SyncGW(targetID)
		if gwInfo == nil {
			err = fmt.Errorf("Error in obtaining information for gateway %v", targetID)
			logging.Error(err.Error())
			return nil, err
		}
	}
	temp := make(map[string]*cidoffer.SubCIDOffer, 0)
	response, err := c.core.P2PServer.Request(gwInfo.NetworkAddr, fcrmessages.DHTOfferDiscoveryRequestType, targetID, pieceCID, numDHT, maxOfferRequestedPerDHT)
	if err != nil {
		err = fmt.Errorf("Error in requesting gateway %v for offers in DHT: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	_, contacted, _, _ := fcrmessages.DecodeDHTOfferDiscoveryResponse(response)
	for _, resp := range contacted {
		_, offers, _, _ := fcrmessages.DecodeStandardOfferDiscoveryResponse(&resp)
		for _, offer := range offers {
			temp[offer.GetMessageDigest()] = &offer
		}
	}
	res := make([]cidoffer.SubCIDOffer, 0)
	for _, offer := range temp {
		res = append(res, *offer)
	}
	return res, nil
}
