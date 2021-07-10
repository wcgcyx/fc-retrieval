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
	"io/ioutil"
	"math/big"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/provider/internal/core"
)

// DataRetrievalHandler handles data retrieval request.
func DataRetrievalHandler(reader fcrserver.FCRServerRequestReader, writer fcrserver.FCRServerResponseWriter, request *fcrmessages.FCRReqMsg) error {
	logging.Debug("Handle data retrieval")
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Message decoding
	nonce, senderID, offer, accountAddr, voucher, err := fcrmessages.DecodeDataRetrievalRequest(request)
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
	if lane != 1 {
		err = fmt.Errorf("Not correct lane received expect 1 got %v:", lane)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	expected := big.NewInt(0).Add(c.Settings.SearchPrice, offer.GetPrice())
	if received.Cmp(expected) < 0 {
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

	// Payment is fine, verify offer
	if offer.Verify(c.OfferSigningPubKey) != nil {
		// Refund money
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding: %v", ierr.Error())
		}
		err = fmt.Errorf("Fail to verify the offer signature, refund voucher %v", refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	// Verify offer merkle proof
	if offer.VerifyMerkleProof() != nil {
		// Refund money
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding: %v", ierr.Error())
		}
		err = fmt.Errorf("Fail to verify the offer merkle proof, refund voucher %v", refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	// Verify offer expiry
	if offer.HasExpired() {
		// Refund money
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, big.NewInt(0).Sub(received, c.Settings.SearchPrice))
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding: %v", ierr.Error())
		}
		err = fmt.Errorf("Offer has expired, refund voucher %v", refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	// Offer is verified. Respond
	// First get the tag
	tag := c.OfferMgr.GetTagByCID(offer.GetSubCID())
	// Second read the data
	data, err := ioutil.ReadFile(filepath.Join(c.Settings.RetrievalDir, tag))
	if err != nil {
		// Refund money, internal error, refund all
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, received)
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding: %v", ierr.Error())
		}
		err = fmt.Errorf("Internal error in finding the content, refund voucher %v", refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	// Third encoding response
	response, err := fcrmessages.EncodeDataRetrievalResponse(nonce, tag, data)
	if err != nil {
		// Refund money, internal error, refund all
		var ierr error
		refundVoucher, ierr = c.PaymentMgr.Refund(accountAddr, lane, received)
		if ierr != nil {
			// This should never happen
			logging.Error("Error in refunding: %v", ierr.Error())
		}
		err = fmt.Errorf("Internal error in encoding the response, refund voucher %v", refundVoucher)
		logging.Error(err.Error())
		return writer.Write(fcrmessages.CreateFCRACKErrorMsg(nonce, err), c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
	}
	c.OfferMgr.IncrementCIDAccessCount(offer.GetSubCID())

	return writer.Write(response, c.MsgSigningKey, c.MsgSigningKeyVer, c.Settings.TCPInactivityTimeout)
}
