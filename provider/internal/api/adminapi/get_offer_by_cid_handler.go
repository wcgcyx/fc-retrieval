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

	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcradminmsg"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// GetOfferByCIDHandler handles get offer by cid request
func GetOfferByCIDHandler(data []byte) (byte, []byte, error) {
	// Get core
	c := core.GetSingleInstance()
	if !c.Initialised {
		// Not initialised.
		err := errors.New("Not initialised")
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	// Decode payload
	cidStr, err := fcradminmsg.DecodeGetOfferByCIDRequest(data)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	cid, err := cid.NewContentID(cidStr)
	if err != nil {
		err = fmt.Errorf("Error in decoding cid string: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}

	digests := make([]string, 0)
	providers := make([]string, 0)
	prices := make([]string, 0)
	expiry := make([]int64, 0)
	qos := make([]uint64, 0)

	offers := c.OfferMgr.GetOffers(cid)
	for _, offer := range offers {
		digests = append(digests, offer.GetMessageDigest())
		providers = append(providers, offer.GetProviderID())
		prices = append(prices, offer.GetPrice().String())
		expiry = append(expiry, offer.GetExpiry())
		qos = append(qos, offer.GetQoS())
	}

	// Succeed
	response, err := fcradminmsg.EncodeGetOfferByCIDResponse(digests, providers, prices, expiry, qos)
	if err != nil {
		err = fmt.Errorf("Error in encoding response: %v", err.Error())
		ack := fcradminmsg.EncodeACK(false, err.Error())
		return fcradminmsg.ACKType, ack, err
	}
	return fcradminmsg.GetOfferByCIDResponseType, response, nil
}
