/*
Package core - structure representing a Client's current state, including setting, configuration, references to
all running Client APIs of this instance.
*/
package core

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
	"math/big"
	"sync"
	"time"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcroffermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpeermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrregistermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrreputationmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
)

type Core struct {
	// Msg Key
	MsgKey string
	// Node ID, calculated from msg key
	NodeID string

	// Wallet address
	WalletAddr string

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

	// Timeout constants
	TCPInactivityTimeout     time.Duration
	LongTCPInactivityTimeout time.Duration

	// Payment related
	SearchPrice *big.Int
	OfferPrice  *big.Int
	TopupAmount *big.Int
}

// Single instance of the gateway
var instance *Core
var doOnce sync.Once

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance() *Core {
	doOnce.Do(func() {
		instance = &Core{
			MsgKey:                   "",
			NodeID:                   "",
			WalletAddr:               "",
			P2PServer:                nil,
			RegisterMgr:              nil,
			PeerMgr:                  nil,
			PaymentMgr:               nil,
			OfferMgr:                 nil,
			ReputationMgr:            nil,
			TCPInactivityTimeout:     5000 * time.Millisecond,
			LongTCPInactivityTimeout: 300000 * time.Millisecond,
			SearchPrice:              big.NewInt(1_000_000_000_000_000),
			OfferPrice:               big.NewInt(1_000_000_000_000_000),
			TopupAmount:              big.NewInt(100_000_000_000_000_000),
		}
	})
	return instance
}
