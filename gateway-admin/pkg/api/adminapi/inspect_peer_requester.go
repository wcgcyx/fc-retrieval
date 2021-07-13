/*
Package adminapi - contains the the adminapi code.
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
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// RequestInspectPeer inspects a given peer from a managed gateway
func RequestInspectPeer(adminURL string, adminKey string, targetID string) (
	int64, // score
	bool, // pending
	bool, // blocked
	[]string, // history
	error, // error
) {
	request, err := fcradminmsg.EncodeInspectPeerRequest(targetID)
	if err != nil {
		err = fmt.Errorf("Error in encoding request: %v", request)
		logging.Error(err.Error())
		return 0, false, false, nil, err
	}

	respType, respData, err := fcradminserver.Request(adminURL, adminKey, fcradminmsg.InspectPeerRequestType, request)
	if err != nil {
		err = fmt.Errorf("Error in sending request: %v", err.Error())
		logging.Error(err.Error())
		return 0, false, false, nil, err
	}

	if respType != fcradminmsg.InspectPeerResponseType {
		err = fmt.Errorf("Getting response of wrong type expect %v, got %v", fcradminmsg.InspectPeerResponseType, respType)
		logging.Error(err.Error())
		return 0, false, false, nil, err
	}

	return fcradminmsg.DecodeInspectPeerResponse(respData)
}
