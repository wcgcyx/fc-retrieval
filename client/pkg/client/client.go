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
	"math/big"
	"net/http"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/wcgcyx/fc-retrieval/client/pkg/api/p2papi"
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
)

// FilecoinRetrievalClient is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Client with Filecoin Retrieval Gateways/Providers.
type FilecoinRetrievalClient struct {
	// Msg Key
	MsgKey string
	// Node ID, calculated from msg key
	NodeID string

	// The P2P Server
	P2PServer fcrserver.FCRServer

	// The Register Manager
	RegisterMgr fcrregistermgr.FCRRegisterMgr

	// The Peer Manager
	PeerMgr fcrpeermgr.FCRPeerMgr

	// The Payment Manager
	PaymentMgr fcrpaymentmgr.FCRPaymentMgr

	// The Offer Manager
	OfferMgr fcroffermgr.FCROfferMgr

	// The Reputation Manager
	ReputationMgr fcrreputationmgr.FCRReputationMgr

	// Payment related
	SearchPrice *big.Int
	OfferPrice  *big.Int
	TopupAmount *big.Int
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
	// Initialise client
	c := &FilecoinRetrievalClient{
		SearchPrice: big.NewInt(1_000_000_000_000_000),
		OfferPrice:  big.NewInt(1_000_000_000_000_000),
		TopupAmount: big.NewInt(100_000_000_000_000_000),
	}

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
		AddRequester(fcrmessages.EstablishmentType, p2papi.EstablishmentRequester).
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

	return c, nil
}

// Search searches gateway that is in given location.
func (c *FilecoinRetrievalClient) Search(location string) ([]string, error) {
	// TODO, Search by location
	res := make([]string, 0)
	infos, err := c.RegisterMgr.GetAllRegisteredGateway(0, 0)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		res = append(res, info.NodeID)
	}
	return res, nil
}

// AddActive adds an active gateway ID
func (c *FilecoinRetrievalClient) AddActive(nodeID string) error {
	// Get gw info
	gw, err := c.PeerMgr.GetGWInfo(nodeID)
	if err != nil {
		// Try again
		c.PeerMgr.SyncGW(nodeID)
		gw, err = c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			return err
		}
	}
	_, err = c.P2PServer.Request(gw.NetworkAddr, fcrmessages.EstablishmentType, gw.NodeID)
	if err != nil {
		return err
	}
	// Create payment channel
	recipientAddr, err := fcrcrypto.GetWalletAddress(gw.RootKey)
	if err != nil {
		return err
	}

	err = c.PaymentMgr.Create(recipientAddr, c.TopupAmount)
	if err != nil {
		return err
	}

	// Add gateway entry to reputation
	c.ReputationMgr.AddGW(gw.NodeID)

	return nil
}

// ListActive lists all active gateways
func (c *FilecoinRetrievalClient) ListActive() ([]string, error) {
	return c.ReputationMgr.ListGWS(), nil
}

// StandardDiscovery performs a standard discovery.
func (c *FilecoinRetrievalClient) StandardDiscovery(cidStr string) ([]cidoffer.SubCIDOffer, error) {
	res := make([]cidoffer.SubCIDOffer, 0)
	return nil, nil
}

// DHTDiscovery performs a DHT discovery.
func (c *FilecoinRetrievalClient) DHTDiscovery(cidStr string, gwID string) ([]cidoffer.SubCIDOffer, error) {
	return nil, nil
}
