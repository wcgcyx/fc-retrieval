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
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// FilecoinRetrievalGatewayAdmin is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Gateway Admin with Filecoin Retrieval Gateways.
type FilecoinRetrievalGatewayAdmin struct {
	adminKey string

	activeGateways     map[string]string
	activeGatewaysLock sync.RWMutex
}

// NewFilecoinRetrievalGatewayAdmin initialises the Filecoin Retrieval Gateway Admin.
func NewFilecoinRetrievalGatewayAdmin(adminKey string) (*FilecoinRetrievalGatewayAdmin, error) {
	// Logging init
	logging.InitWithoutConfig("debug", "STDOUT", "gatewayadmin", "RFC3339")

	// Check admin key is a valid 32 bytes hex string
	token, err := hex.DecodeString(adminKey)
	if err != nil || len(token) != 32 {
		err = fmt.Errorf("Provided admin key is not 32 bytes hex string")
		logging.Error(err.Error())
		return nil, err
	}

	return &FilecoinRetrievalGatewayAdmin{
		adminKey:           adminKey,
		activeGateways:     make(map[string]string),
		activeGatewaysLock: sync.RWMutex{},
	}, nil
}

func (a *FilecoinRetrievalGatewayAdmin) InitialiseGateway(
	adminURL string,
	p2pPort int,
	gatewayIP string,
	rootPrivKey string,
	lotusAPIAddr string,
	lotusAuthToken string,
	registerPrivKey string,
	registerAPIAddr string,
	registerAuthToken string,
) error {
	return nil
}
