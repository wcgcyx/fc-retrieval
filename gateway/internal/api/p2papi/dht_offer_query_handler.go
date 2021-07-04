/*
Package p2papi contains the API code for the p2p communication.
*/
package p2papi

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

	"math/big"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// DHTOfferQueryHandler handles dht offer query.
func DHTOfferQueryHandler(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Response
	var response *fcrmessages.FCRMessage

	// Message decoding
	nodeID, pieceCID, nonce, numDHT, maxOfferRequestedPerDHT, accountAddr, voucher, err := fcrmessages.DecodeDHTOfferDiscoveryRequest(request)
	if err != nil {
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in decoding payload: %v", err.Error()))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	if request.VerifyByID(nodeID) != nil {
		// Message fails to verify
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in verifying msg"))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// Check payment
	received, lane, err := c.PaymentMgr.Receive(accountAddr, voucher)
	if lane != 0 {
		logging.Warn("Payment not in correct lane, should be 0 got %v", lane)
	}
	// expected is 1 * search price + numDHT * (search price + max offer per DHT * offer price)
	expected := big.NewInt(0).Add(c.Settings.SearchPrice, big.NewInt(0).Mul(big.NewInt(0).Add(c.Settings.SearchPrice, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(maxOfferRequestedPerDHT))), big.NewInt(numDHT)))
	if expected.Cmp(received) < 0 {
		// Short payment
		voucher, err := c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Short payment received: %v, expected: %v, refund voucher: %v", received.String(), expected.String(), voucher))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	cidHash, err := pieceCID.CalculateHash()
	if err != nil {
		refundVoucher, err := c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in calculating cid hash: %v, refund voucher: %v", err.Error(), refundVoucher))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// Payment is fine, search.
	refundVoucher := ""

	gws, err := c.PeerMgr.GetGWSNearCIDHash(hex.EncodeToString(cidHash), c.NodeID)
	if err != nil {
		// Internal error in generating sub offers
		refundVoucher, err := c.PaymentMgr.Refund(accountAddr, lane, received)
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Internal error, refund voucher: %v", refundVoucher))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// TODO: Cuncurrently
	supposed := big.NewInt(0).Set(c.Settings.SearchPrice)
	contacted := make(map[string]*fcrmessages.FCRMessage)
	for _, gw := range gws {
		resp, err := c.P2PServer.Request(gw.NetworkAddr, fcrmessages.StandardOfferDiscoveryRequestType, gw.NodeID, pieceCID, maxOfferRequestedPerDHT)
		if err != nil {
			continue
		}
		found := maxOfferRequestedPerDHT
		offers, _, _, _ := fcrmessages.DecodeStandardOfferDiscoveryResponse(resp)
		if len(offers) < int(maxOfferRequestedPerDHT) {
			found = int64(len(offers))
		}
		supposed.Add(c.Settings.SearchPrice, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(found)))
		contacted[gw.NodeID] = resp
	}
	if supposed.Cmp(expected) < 0 {
		refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(expected, supposed))
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
	}

	response, err = fcrmessages.EncodeDHTOfferDiscoveryResponse(contacted, nonce, refundVoucher)
	if err != nil {
		// Internal error in encoding
		refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, received)
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Internal error, refund voucher: %v", refundVoucher))
	}
	err = response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		logging.Error("Error in signing response: %v", err.Error())
	}

	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
