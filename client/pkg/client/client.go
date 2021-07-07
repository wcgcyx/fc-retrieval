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
		return nil, err
	}
	c.MsgKey = msgKey
	_, nodeID, err := fcrcrypto.GetPublicKey(msgKey)
	if err != nil {
		return nil, err
	}
	c.NodeID = nodeID

	// Initialise P2P Server
	// Generate Keypair
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}
	prvKeyBytes, err := prvKey.Raw()
	if err != nil {
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
		return nil, err
	}

	c.RegisterMgr = fcrregistermgr.NewFCRRegisterMgrImplV1(registerAPIAddr, &http.Client{Timeout: 180 * time.Second})
	c.PeerMgr = fcrpeermgr.NewFCRPeerMgrImplV1(c.RegisterMgr, false, false, false, nodeID, time.Hour)
	err = c.PeerMgr.Start()
	if err != nil {
		return nil, err
	}

	lotusMgr := fcrlotusmgr.NewFCRLotusMgrImplV1(lotusAPIAddr, lotusAuthToken, nil)
	c.PaymentMgr = fcrpaymentmgr.NewFCRPaymentMgrImplV1(walletPrvKey, lotusMgr)
	err = c.PaymentMgr.Start()
	if err != nil {
		return nil, err
	}

	c.OfferMgr = fcroffermgr.NewFCROfferMgrImplV1(true)
	err = c.OfferMgr.Start()
	if err != nil {
		return nil, err
	}

	c.ReputationMgr = fcrreputationmgr.NewFCRReputationMgrImpV1()
	err = c.ReputationMgr.Start()
	if err != nil {
		return nil, err
	}

	return &FilecoinRetrievalClient{
		core: c,
	}, nil
}

// Search searches gateway that is in given location.
func (c *FilecoinRetrievalClient) Search(location string) ([]string, error) {
	// TODO, Search by location
	res := make([]string, 0)
	infos, err := c.core.RegisterMgr.GetAllRegisteredGateway(0, 0)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		res = append(res, info.NodeID)
	}
	return res, nil
}

// AddActive adds an active gateway ID
func (c *FilecoinRetrievalClient) AddActive(targetID string) error {
	// Get gw info
	gwInfo := c.core.PeerMgr.GetGWInfo(targetID)
	if gwInfo == nil {
		// Not found, try sync once
		gwInfo = c.core.PeerMgr.SyncGW(targetID)
		if gwInfo == nil {
			return fmt.Errorf("Error in obtaining information for gateway %v", targetID)
		}
	}
	_, err := c.core.P2PServer.Request(gwInfo.NetworkAddr, fcrmessages.EstablishmentRequestType, targetID)
	if err != nil {
		return fmt.Errorf("Error in sending establishment request to %v with addr %v: %v", targetID, gwInfo.NetworkAddr, err.Error())
	}
	// Create payment channel
	recipientAddr, err := fcrcrypto.GetWalletAddress(gwInfo.RootKey)
	if err != nil {
		return fmt.Errorf("Error in obtaining wallet addreess for gateway %v with root key %v: %v", targetID, gwInfo.RootKey, err.Error())
	}

	err = c.core.PaymentMgr.Create(recipientAddr, c.core.TopupAmount)
	if err != nil {
		return fmt.Errorf("Error in creating a payment channel to %v with wallet address %v with topup amount of %v: %v", targetID, recipientAddr, c.core.TopupAmount.String(), err.Error())
	}
	// Add gateway entry to reputation
	c.core.ReputationMgr.AddGW(gwInfo.NodeID)
	return nil
}

// ListActive lists all active gateways
func (c *FilecoinRetrievalClient) ListActive() ([]string, error) {
	return c.core.ReputationMgr.ListGWS(), nil
}

// StandardDiscovery performs a standard discovery.
func (c *FilecoinRetrievalClient) StandardDiscovery(cidStr string, toContact map[string]uint32) ([]cidoffer.SubCIDOffer, error) {
	return nil, nil
}

// DHTDiscovery performs a DHT discovery.
func (c *FilecoinRetrievalClient) DHTDiscovery(cidStr string, targetID string, numDHT uint32, maxOfferRequestedPerDHT uint32) ([]cidoffer.SubCIDOffer, error) {
	return nil, nil
}
