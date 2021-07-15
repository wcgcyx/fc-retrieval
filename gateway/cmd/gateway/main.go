/*
Package main - program entry point for a Retrieval Gateway node.

Retrieval Gateway is a type of nodes in FileCoin blockchain network, which serves purpose of being first point of contact
for a client, who is trying to find and retrieve their files.
Retrieval Gateway is responsible for providing the best way for the client to get their files back from the network.
*/
package main

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
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminserver"
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
	"github.com/wcgcyx/fc-retrieval/gateway/internal/api/adminapi"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/api/p2papi"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/config"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// Start Gateway service
func main() {
	// Configure what should be called if Control-C is hit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		gracefulExit()
		os.Exit(0)
	}()

	// Load config
	conf := config.NewConfig()
	appSettings := config.Map(conf)

	// Initialise logging
	logging.Init(conf)

	// Initialise gateway core instance
	c := core.GetSingleInstance(&appSettings)

	// Attempt to load token
	var token [32]byte
	os.MkdirAll(c.Settings.SystemDir, os.ModePerm)
	os.MkdirAll(c.Settings.RetrievalDir, os.ModePerm)
	f, err := os.Open(c.Settings.AdminKeyFile)
	if err == nil {
		n, err := f.Read(token[:])
		if n != 32 || err != nil {
			rand.Read(token[:])
		}
		f.Close()
	} else {
		rand.Read(token[:])
	}
	f, err = os.Create(c.Settings.AdminKeyFile)
	if err != nil {
		logging.Error("Error in creating token file: %v", err.Error())
		return
	}
	f.Write(token[:])
	f.Close()

	adminKey := hex.EncodeToString(token[:])
	logging.Info("Admin access token is %v and it has been saved to %v", adminKey, c.Settings.AdminKeyFile)

	// Start the Admin API, waiting for initialisation.
	c.AdminServer = fcradminserver.NewFCRAdminServerImplV1(fmt.Sprintf(":%v", c.Settings.BindAdminAPI), adminKey)
	c.AdminServer.
		// Handlers
		AddHandler(fcradminmsg.InitialisationRequestType, adminapi.InitialisationHandler).
		AddHandler(fcradminmsg.CacheOfferByDigestRequestType, adminapi.CacheOfferByDigestHandler).
		AddHandler(fcradminmsg.ChangePeerStatusRequestType, adminapi.ChangePeerStatusHandler).
		AddHandler(fcradminmsg.GetOfferByCIDRequestType, adminapi.GetOfferByCIDHandler).
		AddHandler(fcradminmsg.InspectPeerRequestType, adminapi.InspectPeerHandler).
		AddHandler(fcradminmsg.ListCIDFrequencyRequestType, adminapi.ListCIDFrequencyHandler).
		AddHandler(fcradminmsg.ListPeersRequestType, adminapi.ListPeersHandler).
		AddHandler(fcradminmsg.ForceSyncRequestType, adminapi.ForceSyncHandler)

	err = c.AdminServer.Start()
	if err != nil {
		logging.Error("Error in starting admin server: %v", err)
		return
	}
	logging.Info("Admin server starts listening on [::]:%v", c.Settings.BindAdminAPI)

	// Attempt to load config file and initialise this gateway
	go func() {
		f, err = os.Open(c.Settings.ConfigFile)
		if err != nil {
			return
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return
		}
		config := strings.Split(string(data), ";")
		if len(config) != 10 {
			return
		}
		p2pPrivKey := config[0]
		p2pPort, err := strconv.ParseInt(config[1], 10, 32)
		if err != nil {
			return
		}
		rootPrivKey := config[2]
		lotusAPIAddr := config[3]
		lotusAuthToken := config[4]
		// registerPrivKey := config[5]
		registerAPIAddr := config[6]
		// registerAuthToken := config[7]
		msgSigningKey := config[8]
		msgSigningKeyVer, err := strconv.ParseInt(config[9], 10, 32)
		if err != nil {
			return
		}
		rootKey, nodeID, err := fcrcrypto.GetPublicKey(rootPrivKey)
		if err != nil {
			return
		}
		c.NodeID = nodeID
		c.WalletAddr, err = fcrcrypto.GetWalletAddress(rootKey)
		if err != nil {
			return
		}
		temp, err := hex.DecodeString(msgSigningKey)
		if err != nil {
			return
		}
		if len(temp) != 32 {
			return
		}
		c.MsgSigningKey = msgSigningKey
		c.MsgSigningKeyVer = byte(msgSigningKeyVer)
		c.P2PServer = fcrserver.NewFCRServerImplV1(p2pPrivKey, uint(p2pPort), c.Settings.TCPInactivityTimeout)
		c.ReputationMgr = fcrreputationmgr.NewFCRReputationMgrImpV1()
		registerMgr := fcrregistermgr.NewFCRRegisterMgrImplV1(registerAPIAddr, &http.Client{Timeout: 180 * time.Second})
		c.StoreFullOffer = c.Settings.StoreFullOffer
		c.PeerMgr = fcrpeermgr.NewFCRPeerMgrImplV1(registerMgr, c.ReputationMgr, true, true, !c.StoreFullOffer, nodeID, c.Settings.SyncDuration)
		lotusMgr := fcrlotusmgr.NewFCRLotusMgrImplV1(lotusAPIAddr, lotusAuthToken, nil)
		c.PaymentMgr = fcrpaymentmgr.NewFCRPaymentMgrImplV1(rootPrivKey, lotusMgr)
		c.OfferMgr = fcroffermgr.NewFCROfferMgrImplV1(true)
		c.Ready <- true
		if !<-c.Ready {
			return
		}
		c.Initialised = true
		c.Ready <- true
	}()

	// Wait for admin to initialise this gateway
	for !<-c.Ready {
	}

	// Gateway has been initialised.
	c.P2PServer.
		// Handlers
		AddHandler(fcrmessages.EstablishmentRequestType, p2papi.EstablishmentHandler).
		AddHandler(fcrmessages.StandardOfferDiscoveryRequestType, p2papi.OfferQueryHandler).
		AddHandler(fcrmessages.DHTOfferDiscoveryRequestType, p2papi.DHTOfferQueryHandler).
		AddHandler(fcrmessages.OfferPublishRequestType, p2papi.OfferPublishHandler).
		// Requesters
		AddRequester(fcrmessages.StandardOfferDiscoveryRequestType, p2papi.OfferQueryRequester).
		AddRequester(fcrmessages.EstablishmentRequestType, p2papi.EstablishmentRequester).
		AddRequester(fcrmessages.DataRetrievalRequestType, p2papi.DataRetrievalRequester)

	err = c.P2PServer.Start()
	if err != nil {
		logging.Error("Error in starting P2P Server: %v", err)
		c.Ready <- false
		gracefulExit()
		return
	}

	err = c.PeerMgr.Start()
	if err != nil {
		logging.Error("Error in starting Peer Manager: %v", err)
		c.Ready <- false
		gracefulExit()
		return
	}

	err = c.PaymentMgr.Start()
	if err != nil {
		logging.Error("Error in starting Payment Manager: %v", err)
		c.Ready <- false
		gracefulExit()
		return
	}

	err = c.OfferMgr.Start()
	if err != nil {
		logging.Error("Error in starting Offer Manager: %v", err)
		c.Ready <- false
		gracefulExit()
		return
	}

	err = c.ReputationMgr.Start()
	if err != nil {
		logging.Error("Error in starting Reputation Manager: %v", err)
		c.Ready <- false
		gracefulExit()
		return
	}

	// Everything has been started.
	c.Ready <- true
	// Wait for this gateway to be registered.
	if !<-c.Ready {
		// Register failed.
		logging.Error("Error in registering this gateway.")
		gracefulExit()
		return
	}
	// Register succeed. Run gateway
	logging.Info("Filecoin Gateway Start-up Complete")
	c.PeerMgr.Sync()

	// Wait forever
	// TODO: Start message signing key update routine.
	select {}
}

// gracefulExit handles exit
func gracefulExit() {
	logging.Info("Filecoin Gateway Shutdown: Start")
	// Delay 3 seconds to let admin knows any error.
	time.Sleep(3 * time.Second)

	c := core.GetSingleInstance()
	if c.AdminServer != nil {
		c.AdminServer.Shutdown()
	}
	if c.P2PServer != nil {
		c.P2PServer.Shutdown()
	}
	if c.PeerMgr != nil {
		c.PeerMgr.Shutdown()
	}
	if c.PaymentMgr != nil {
		c.PaymentMgr.Shutdown()
	}
	if c.OfferMgr != nil {
		c.OfferMgr.Shutdown()
	}
	if c.ReputationMgr != nil {
		c.ReputationMgr.Shutdown()
	}

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
