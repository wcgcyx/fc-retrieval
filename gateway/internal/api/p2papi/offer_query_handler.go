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
	"fmt"
	"time"

	"math/big"

	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

// OfferQueryHandler handles standard offer query.
func OfferQueryHandler(reader fcrserver.FCRServerRequestReader, writer fcrserver.FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error {
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Message decoding
	nonce, senderID, pieceCID, maxOfferRequested, accountAddr, voucher, err := fcrmessages.DecodeStandardOfferDiscoveryRequest(request)
	if err != nil {
		err = fmt.Errorf("Error in decoding payload: %v", err.Error())
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Verify signature
	if request.VerifyByID(senderID) != nil {
		// Verify by signing key
		gwInfo := c.PeerMgr.GetGWInfo(senderID)
		if gwInfo == nil {
			// Not found, try sync once
			gwInfo = c.PeerMgr.SyncGW(senderID)
			if gwInfo == nil {
				err = fmt.Errorf("Error in obtaining information for gateway %v", senderID)
				logging.Error(err.Error())
				return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
			}
		}
		if request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
			// Try update
			gwInfo = c.PeerMgr.SyncGW(senderID)
			if gwInfo == nil || request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
				err = fmt.Errorf("Error in verifying request from gateway %v: %v", senderID, err.Error())
				logging.Error(err.Error())
				return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
			}
		}
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
	expected := big.NewInt(0).Add(c.Settings.SearchPrice, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(int64(maxOfferRequested))))
	if expected.Cmp(received) < 0 {
		// Short payment
		// Refund money
		if received.Cmp(c.Settings.SearchPrice) <= 0 {
			// No refund
		} else {
			var ierr error
			refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
			if ierr != nil {
				// This should never happen
				logging.Error("Error in refunding: %v", ierr.Error())
			}
		}
		err = fmt.Errorf("Short payment received, expect %v got %v, refund voucher %v", expected.String(), received.String(), refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}

	// Payment is fine, search.
	c.OfferMgr.IncrementCIDAccessCount(pieceCID)
	offers := c.OfferMgr.GetOffers(pieceCID)

	// Generating sub CID offers
	res := make([]cidoffer.SubCIDOffer, 0)
	remain := int64(maxOfferRequested)
	for _, offer := range offers {
		if remain == 0 {
			break
		}
		// Check offer expiry, remove if less than 1 hour + 1 hour room
		if offer.GetExpiry()-time.Now().Unix() < 7200 {
			// Offer is soon to expire
			c.OfferMgr.RemoveOffer(offer.GetMessageDigest())
			continue
		}

		subOffer, err := offer.GenerateSubCIDOffer(pieceCID)
		if err != nil {
			// Internal error in generating sub offers
			var ierr error
			refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, received)
			if ierr != nil {
				// This should never happen
				logging.Error("Error in refunding: %v", ierr.Error())
			}
			err = fmt.Errorf("Internal error in generating sub cid offer: %v, refund voucher %v", err.Error(), refundVoucher)
			logging.Error(err.Error())
			return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
		}
		res = append(res, *subOffer)
		remain--
	}
	if remain > 0 {
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, big.NewInt(0).Mul(c.Settings.OfferPrice, big.NewInt(remain))))
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding %v", ierr.Error())
		}
	}

	// Respond
	response, err := fcrmessages.EncodeStandardOfferDiscoveryResponse(nonce, res, refundVoucher)
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
