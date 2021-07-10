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
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrreputationmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// InspectPeerHandler handles inspect peer request
func InspectPeerHandler(data []byte) (byte, []byte, error) {
	logging.Debug("Handle inspect peer from admin")
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Decode payload
	nodeID, gateway, err := fcradminmsg.DecodeInspectPeerRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	var rep *fcrreputationmgr.Reputation
	var history []reputation.Record
	if gateway {
		rep = c.ReputationMgr.GetGWReputation(nodeID)
		history = c.ReputationMgr.GetGWHistory(nodeID, 0, 10)
	} else {
		rep = c.ReputationMgr.GetPVDReputation(nodeID)
		history = c.ReputationMgr.GetPVDHistory(nodeID, 0, 10)
	}
	if rep == nil {
		err = fmt.Errorf("Cannot find reputation for: %v", nodeID)
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	reasons := make([]string, 0)
	for _, entry := range history {
		reasons = append(reasons, entry.Reason())
	}

	// Succeed
	response, err := fcradminmsg.EncodeInspectPeerResponse(rep.Score, rep.Pending, rep.Blocked, reasons)
	if err != nil {
		err = fmt.Errorf("Error in encoding response: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	return fcradminmsg.InspectPeerResponseType, response, nil
}
