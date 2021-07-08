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

// RequestGetOfferByCID gets offers containing given cid from a managed gateway
func RequestGetOfferByCID(adminURL string, adminKey string, cid string) (
	[]string, // digests
	[]string, // providers
	[]string, // prices
	[]int64, // expiry
	[]uint64, // qos
	error, // error
) {
	request, err := fcradminmsg.EncodeGetOfferByCIDRequest(cid)
	if err != nil {
		err = fmt.Errorf("Error in encoding request: %v", request)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}

	respType, respData, err := fcradminserver.Request(adminURL, adminKey, fcradminmsg.GetOfferByCIDRequestType, request)
	if err != nil {
		err = fmt.Errorf("Error in sending request: %v", err.Error())
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}

	if respType != fcradminmsg.GetOfferByCIDResponseType {
		err = fmt.Errorf("Getting response of wrong type expect %v, got %v", fcradminmsg.GetOfferByCIDResponseType, respType)
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, err
	}

	return fcradminmsg.DecodeGetOfferByCIDResponse(respData)
}
