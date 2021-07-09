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
func DHTOfferQueryHandler(reader fcrserver.FCRServerRequestReader, writer fcrserver.FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error {
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Message decoding
	nonce, senderID, pieceCID, numDHT, maxOfferRequestedPerDHT, accountAddr, voucher, err := fcrmessages.DecodeDHTOfferDiscoveryRequest(request)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Verify signature
	if request.VerifyByID(senderID) != nil {
		err = fmt.Errorf("Error in verifying request from %v: %v", senderID, err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Check numDHT
	if numDHT > 16 {
		err = fmt.Errorf("Error exceeding maximum numDHT 16 from %v, got %v", senderID, numDHT)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Check payment
	refundVoucher := ""
	received, lane, err := c.PaymentMgr.Receive(accountAddr, voucher)
	if err != nil {
		err = fmt.Errorf("Error in receiving voucher %v:", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	if lane != 0 {
		err = fmt.Errorf("Not correct lane received expect 0 got %v:", lane)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	// expected is 1 * search price + numDHT * (search price + max offer per DHT * offer price)
	expected := big.NewInt(0).Add(c.Settings.SearchPrice, big.NewInt(0).Mul(big.NewInt(0).Add(c.Settings.SearchPrice, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(int64(maxOfferRequestedPerDHT)))), big.NewInt(int64(numDHT))))
	if received.Cmp(expected) < 0 {
		// Short payment
		// Refund money
		if received.Cmp(c.Settings.SearchPrice) <= 0 {
			// No refund
		} else {
			refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
			if err != nil {
				// This should never happen
				logging.Error("Error in refunding: %v", err.Error())
			}
		}
		err = fmt.Errorf("Short payment received, expect %v got %v, refund voucher %v", expected.String(), received.String(), refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Payment is fine, search.
	c.OfferMgr.IncrementCIDAccessCount(pieceCID)
	cidHash, err := pieceCID.CalculateHash()
	if err != nil {
		// Internal error in calculating cid hash
		var ierr error
		refundVoucher, ierr := c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", ierr.Error())
		}
		err = fmt.Errorf("Error in calculating cid hash: %v, refund voucher: %v", err.Error(), refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Get gateways
	gws := c.PeerMgr.GetGWSNearCIDHash(hex.EncodeToString(cidHash), int(numDHT), c.NodeID)
	if err != nil {
		// Internal error in getting near gateways
		var ierr error
		refundVoucher, ierr := c.PaymentMgr.Refund(accountAddr, lane, received)
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding %v", ierr.Error())
		}
		err = fmt.Errorf("Internal error in getting near gateways to requested cid: %v with hash: %v, refund voucher: %v", pieceCID.ToString(), hex.EncodeToString(cidHash), refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// TODO: Concurrency
	supposed := big.NewInt(0).Set(c.Settings.SearchPrice)
	contacted := make(map[string]*fcrmessages.FCRACKMsg)
	for _, gw := range gws {
		resp, err := c.P2PServer.Request(gw.NetworkAddr, fcrmessages.StandardOfferDiscoveryRequestType, gw.NodeID, pieceCID, maxOfferRequestedPerDHT)
		if err != nil {
			continue
		}
		found := int64(maxOfferRequestedPerDHT)
		_, offers, _, _ := fcrmessages.DecodeStandardOfferDiscoveryResponse(resp)
		if len(offers) < int(maxOfferRequestedPerDHT) {
			found = int64(len(offers))
		}
		supposed.Add(c.Settings.SearchPrice, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(found)))
		contacted[gw.NodeID] = resp
	}
	if supposed.Cmp(expected) < 0 {
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(expected, supposed))
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding %v", ierr.Error())
		}
	}

	// Respond
	response, err := fcrmessages.EncodeDHTOfferDiscoveryResponse(nonce, contacted, refundVoucher)
	if err != nil {
		// Internal error in encoding
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, received)
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding %v", ierr.Error())
		}
		err = fmt.Errorf("Internal error in encoding response: %v, refund voucher %v", err.Error(), refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	return writer.Write(response, c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
}
