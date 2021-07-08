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

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// ListPeersHandler handles get offer by cid request
func ListPeersHandler(data []byte) (byte, []byte, error) {
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	gatewayIDs := make([]string, 0)
	gatewayScore := make([]int64, 0)
	gatewayPending := make([]bool, 0)
	gatewayBlocked := make([]bool, 0)
	gatewayRecent := make([]string, 0)
	providerIDs := make([]string, 0)
	providerScore := make([]int64, 0)
	providerPending := make([]bool, 0)
	providerBlocked := make([]bool, 0)
	providerRecent := make([]string, 0)

	for _, gw := range c.ReputationMgr.ListGWS() {
		rep := c.ReputationMgr.GetGWReputation(gw)
		gatewayIDs = append(gatewayIDs, gw)
		gatewayScore = append(gatewayScore, rep.Score)
		gatewayPending = append(gatewayPending, rep.Pending)
		gatewayBlocked = append(gatewayBlocked, rep.Blocked)
		recent := c.ReputationMgr.GetGWHistory(gw, 0, 1)
		if len(recent) == 0 {
			gatewayRecent = append(gatewayRecent, "No record found")
		} else {
			gatewayRecent = append(gatewayRecent, recent[0].Reason())
		}
	}
	for _, pvd := range c.ReputationMgr.ListPVDS() {
		rep := c.ReputationMgr.GetGWReputation(pvd)
		providerIDs = append(providerIDs, pvd)
		providerScore = append(providerScore, rep.Score)
		providerPending = append(providerPending, rep.Pending)
		providerBlocked = append(providerBlocked, rep.Blocked)
		recent := c.ReputationMgr.GetPVDHistory(pvd, 0, 1)
		if len(recent) == 0 {
			providerRecent = append(providerRecent, "No record found")
		} else {
			providerRecent = append(providerRecent, recent[0].Reason())
		}
	}

	// Succeed
	response, err := fcradminmsg.EncodeListPeersResponse(gatewayIDs, gatewayScore, gatewayPending, gatewayBlocked, gatewayRecent, providerIDs, providerScore, providerPending, providerBlocked, providerRecent)
	if err != nil {
		err = fmt.Errorf("Error in encoding response: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	return fcradminmsg.ListPeersResponseType, response, nil
}
