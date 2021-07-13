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
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// ListPeersHandler handles list peers request
func ListPeersHandler(data []byte) (byte, []byte, error) {
	logging.Debug("Handle list peers from admin")
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	peerIDs := make([]string, 0)
	peerScore := make([]int64, 0)
	peerPending := make([]bool, 0)
	peerBlocked := make([]bool, 0)
	peerRecent := make([]string, 0)

	for _, peerID := range c.ReputationMgr.ListPeers() {
		rep := c.ReputationMgr.GetPeerReputation(peerID)
		peerIDs = append(peerIDs, peerID)
		peerScore = append(peerScore, rep.Score)
		peerPending = append(peerPending, rep.Pending)
		peerBlocked = append(peerBlocked, rep.Blocked)
		recent := c.ReputationMgr.GetPeerHistory(peerID, 0, 1)
		if len(recent) == 0 {
			peerRecent = append(peerRecent, "No record found")
		} else {
			peerRecent = append(peerRecent, recent[0].Reason())
		}
	}

	// Succeed
	response, err := fcradminmsg.EncodeListPeersResponse(peerIDs, peerScore, peerPending, peerBlocked, peerRecent)
	if err != nil {
		err = fmt.Errorf("Error in encoding response: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	return fcradminmsg.ListPeersResponseType, response, nil
}
