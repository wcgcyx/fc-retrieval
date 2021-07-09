/*
Package core - structure representing a Provider's current state, including setting, configuration, references to
all running Provider APIs of this instance.
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
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcroffermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpaymentmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpeermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/settings"
)

type Core struct {
	// Settings
	Settings *settings.AppSettings

	// Boolean indicating whether or not this provider has been initialised
	Initialised bool
	Ready       chan bool

	// Node ID
	NodeID string
	// Wallet Address
	WalletAddr string

	// Message signing key and a lock protecting the access
	MsgSigningKey     string
	MsgSigningKeyVer  byte
	MsgSigningKeyLock sync.RWMutex

	OfferSigningKey    string
	OfferSigningPubKey string

	// The Admin Server
	AdminServer fcradminserver.FCRAdminServer

	// The P2P Server
	P2PServer fcrserver.FCRServer

	// The Peer Manager
	PeerMgr fcrpeermgr.FCRPeerMgr

	// The Payment Manager
	PaymentMgr fcrpaymentmgr.FCRPaymentMgr

	// The Offer Manager
	OfferMgr fcroffermgr.FCROfferMgr
}

// Single instance of the provider
var instance *Core
var doOnce sync.Once

// GetSingleInstance returns the single instance of the provider
func GetSingleInstance(confs ...*settings.AppSettings) *Core {
	doOnce.Do(func() {
		if len(confs) == 0 {
			logging.Panic("No settings supplied to Provider start-up")
		}
		if len(confs) != 1 {
			logging.Panic("More than one sets of settings supplied to Provider start-up")
		}
		instance = &Core{
			Settings:          confs[0],
			Initialised:       false,
			Ready:             make(chan bool),
			NodeID:            "",
			WalletAddr:        "",
			MsgSigningKey:     "",
			MsgSigningKeyVer:  0,
			MsgSigningKeyLock: sync.RWMutex{},
			OfferSigningKey:   "",
			AdminServer:       nil,
			P2PServer:         nil,
			OfferMgr:          nil,
			PeerMgr:           nil,
			PaymentMgr:        nil,
		}
	})
	return instance
}
