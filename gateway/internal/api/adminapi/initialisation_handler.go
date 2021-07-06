/*
Package adminapi contains the API code for the admin client - gateway communication.
*/
package adminapi

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
	"fmt"
	"net/http"
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrlotusmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcroffermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpeermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrreputationmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// InitialisationHandler handles initialisation.
func InitialisationHandler(data []byte) (byte, []byte, error) {
	// Get core
	c := core.GetSingleInstance()
	if c.Initialised {
		// Already initialised.
		err := errors.New("Already initialised")
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Decoding payload
	p2pPrvKey, p2pPort, networkAddr, rootPrvKey, lotusAPIAddr, lotusAuthToken, _, registerAPIAddr, _, regionCode, err := fcradminmsg.DecodeInitialisationRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Obtaining the root key
	rootKey, nodeID, err := fcrcrypto.GetPublicKey(rootPrvKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining the public key: %v", err.Error())
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	c.NodeID = nodeID
	c.WalletAddr, err = fcrcrypto.GetWalletAddress(rootKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining the wallet address: %v", err.Error())
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Generating msg signing key
	msgKey, msgSigningKey, _, err := fcrcrypto.GenerateRetrievalKeyPair()
	if err != nil {
		err = fmt.Errorf("Error in generating message signing key: %v", err.Error())
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	c.MsgSigningKey = msgKey
	c.MsgSigningKeyVer = 0

	// Initialise P2P Server
	c.P2PServer = fcrserver.NewFCRServerImplV1(p2pPrvKey, uint(p2pPort), c.Settings.TCPInactivityTimeout)

	// Initialise peer manager
	registerMgr := fcrregistermgr.NewFCRRegisterMgrImplV1(registerAPIAddr, &http.Client{Timeout: 180 * time.Second})
	c.StoreFullOffer = c.Settings.StoreFullOffer
	c.PeerMgr = fcrpeermgr.NewFCRPeerMgrImplV1(registerMgr, true, true, !c.StoreFullOffer, nodeID, c.Settings.SyncDuration)

	// Initialise payment manager
	lotusMgr := fcrlotusmgr.NewFCRLotusMgrImplV1(lotusAPIAddr, lotusAuthToken, nil)
	c.PaymentMgr = fcrpaymentmgr.NewFCRPaymentMgrImplV1(rootPrvKey, lotusMgr)

	// Initialise offer manager
	c.OfferMgr = fcroffermgr.NewFCROfferMgrImplV1(true)

	// Initialise reputation manager
	c.ReputationMgr = fcrreputationmgr.NewFCRReputationMgrImpV1()

	// Ask the server to start
	c.Ready <- true
	if !<-c.Ready {
		// Initialisation failed.
		err = errors.New("Initialisation failed")
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Initialisation succeed. Start register this gateway.
	err = registerMgr.RegisterGateway(nodeID, &register.GatewayRegisteredInfo{
		RootKey:             rootKey,
		NodeID:              nodeID,
		MsgSigningKey:       msgSigningKey,
		MsgSigningKeyVer:    0,
		RegionCode:          regionCode,
		NetworkAddr:         networkAddr,
		Deregistering:       false,
		DeregisteringHeight: 0,
	})
	if err != nil {
		c.Ready <- false
		err = fmt.Errorf("Error in registering the gateway: %v", err.Error())
		ack, _ := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Succeed.
	c.Ready <- true
	c.Initialised = true
	ack, _ := fcradminmsg.EncodeACK(true, "Succeed.")
	return fcradminmsg.ACKType, ack, nil
}