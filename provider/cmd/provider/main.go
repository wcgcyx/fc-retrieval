/*
Package main - program entry point for a Retrieval Provider node.

Retrieval Provider is a type of nodes in FileCoin blockchain network, which serves purpose of being a way to
communicate with a Storage Miner.

Retrieval Provider is used by Retrieval Gateways in order to get their files back from the particular Storage Miner
in the network.
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
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/api/adminapi"
	"github.com/wcgcyx/fc-retrieval/provider/internal/api/p2papi"
	"github.com/wcgcyx/fc-retrieval/provider/internal/config"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// Start Provider service
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

	// Initialise provider core instance
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
		AddHandler(fcradminmsg.GetOfferByCIDRequestType, adminapi.GetOfferByCIDHandler).
		AddHandler(fcradminmsg.ListFilesRequestType, adminapi.ListFilesHandler).
		AddHandler(fcradminmsg.PublishOfferRequestType, adminapi.OfferPublishHandler).
		AddHandler(fcradminmsg.UploadFileRequestType, adminapi.UploadFileHandler)

	err = c.AdminServer.Start()
	if err != nil {
		logging.Error("Error in starting admin server: %v", err)
		return
	}
	logging.Info("Admin server starts listening on [::]:%v", c.Settings.BindAdminAPI)

	// Wait for admin to initialise this provider
	for !<-c.Ready {
	}

	// Provider has been initialised.
	c.P2PServer.
		// Handlers
		AddHandler(fcrmessages.EstablishmentRequestType, p2papi.EstablishmentHandler).
		AddHandler(fcrmessages.DataRetrievalRequestType, p2papi.DataRetrievalHandler).
		// Requesters
		AddRequester(fcrmessages.OfferPublishRequestType, p2papi.OfferPublishRequester)

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

	// Everything has been started.
	c.Ready <- true
	// Wait for this provider to be registered.
	if !<-c.Ready {
		// Register failed.
		logging.Error("Error in registering this provider.")
		gracefulExit()
		return
	}
	logging.Info("Filecoin Provider Start-up Complete")

	// Wait forever
	// TODO: Start message signing key update routine.
	select {}
}

// gracefulExit handles exit
func gracefulExit() {
	logging.Info("Filecoin Provider Shutdown: Start")
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

	logging.Info("Filecoin Provider Shutdown: Completed")
}
