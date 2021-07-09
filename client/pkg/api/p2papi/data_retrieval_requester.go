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
	"math/big"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/wcgcyx/fc-retrieval/client/pkg/core"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

// DataRetrievalRequester requests a data retrieval
func DataRetrievalRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
	// Get parameters
	if len(args) != 3 {
		err := fmt.Errorf("Wrong arguments, expect length 3, got length %v", len(args))
		logging.Error(err.Error())
		return nil, err
	}
	targetID, ok := args[0].(string)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a target ID in string")
		logging.Error(err.Error())
		return nil, err
	}
	offer, ok := args[1].(*cidoffer.SubCIDOffer)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a Sub CID Offer in *cidoffer.SubCIDOffer")
		logging.Error(err.Error())
		return nil, err
	}
	retrievalPath, ok := args[2].(string)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a retrievalPath in string")
		logging.Error(err.Error())
		return nil, err
	}

	// Get core structure
	c := core.GetSingleInstance()

	// Generate random nonce
	nonce := uint64(rand.Int63())

	// Get provider information
	pvdInfo := c.PeerMgr.GetPVDInfo(targetID)
	if pvdInfo == nil {
		// Not found, try sync once
		pvdInfo = c.PeerMgr.SyncPVD(targetID)
		if pvdInfo == nil {
			err := fmt.Errorf("Error in obtaining information for provider %v", targetID)
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Check if the provider is blocked/pending
	rep := c.ReputationMgr.GetPVDReputation(targetID)
	if rep == nil {
		err := fmt.Errorf("Provider %v is not active", targetID)
		logging.Error(err.Error())
		return nil, err
	}
	if rep.Pending || rep.Blocked {
		err := fmt.Errorf("Provider %v is in pending %v, blocked %v", targetID, rep.Pending, rep.Blocked)
		logging.Error(err.Error())
		return nil, err
	}

	// Pay the recipient
	recipientAddr, err := fcrcrypto.GetWalletAddress(pvdInfo.RootKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining wallet addreess for provider %v with root key %v: %v", targetID, pvdInfo.RootKey, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	expected := big.NewInt(0).Add(c.SearchPrice, offer.GetPrice())
	voucher, create, topup, err := c.PaymentMgr.Pay(recipientAddr, 0, expected)
	if err != nil {
		err = fmt.Errorf("Error in paying provider %v with expected amount of %v: %v", targetID, expected.String(), err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	if create {
		err = fmt.Errorf("No payment channel to %v", targetID)
		logging.Error(err.Error())
		return nil, err
	} else if topup {
		// Need to topup
		err = c.PaymentMgr.Topup(recipientAddr, c.TopupAmount)
		if err != nil {
			err = fmt.Errorf("Error in topup a payment channel to %v with wallet address %v with topup amount of %v: %v", targetID, recipientAddr, c.TopupAmount.String(), err.Error())
			logging.Error(err.Error())
			return nil, err
		}
		voucher, create, topup, err = c.PaymentMgr.Pay(recipientAddr, 0, expected)
		if create || topup {
			// This should never happen
			err = fmt.Errorf("Error in paying provider %v, needs to create/topup after just topup", targetID)
			logging.Error(err.Error())
			return nil, err
		}
		if err != nil {
			err = fmt.Errorf("Error in paying provider %v with expected amount of %v: %v after just topup", targetID, expected.String(), err.Error())
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Now we have got a voucher
	// Encode request
	request, err := fcrmessages.EncodeDataRetrievalRequest(nonce, c.NodeID, offer, c.WalletAddr, voucher)
	if err != nil {
		c.PaymentMgr.RevertPay(recipientAddr, 0)
		err = fmt.Errorf("Internal error in encoding response: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Write request
	err = writer.Write(request, c.MsgKey, 0, c.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in sending request to %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in receiving response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	// Verify the response
	if response.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) != nil {
		// Try update
		pvdInfo = c.PeerMgr.SyncPVD(targetID)
		if pvdInfo == nil || response.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) != nil {
			err = fmt.Errorf("Error in verifying response from %v: %v", targetID, err.Error())
			logging.Error(err.Error())
			// Pend PVD
			c.ReputationMgr.UpdatePVDRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendPVD(targetID)
			return nil, err
		}
	}

	// Check response
	if !response.ACK() {
		err = fmt.Errorf("Reponse contains an error: %v", response.Error())
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	// Decode response
	nonceRecv, tag, data, err := fcrmessages.DecodeDataRetrievalResponse(response)
	if err != nil {
		err = fmt.Errorf("Error in decoding response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	if nonceRecv != nonce {
		err = fmt.Errorf("Nonce mismatch: expected %v got %v", nonce, nonceRecv)
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	// Save file
	if _, err := os.Stat(filepath.Join(retrievalPath, tag)); os.IsNotExist(err) {
		// Not exist, save
		f, err := os.Create(filepath.Join(retrievalPath, tag))
		if err == nil {
			_, err = f.Write(data)
			f.Close()
		}
		if err != nil {
			err = fmt.Errorf("Error saving file: %v", err.Error())
			logging.Error(err.Error())
			return nil, err
		}
	} else {
		// Exist
		err = fmt.Errorf("Filename already existed %v", tag)
		logging.Error(err.Error())
		return nil, err
	}

	// Read file
	fileReader, err := os.Open(filepath.Join(retrievalPath, tag))
	if err != nil {
		err = fmt.Errorf("Fail to open file for cid calculation %v: %v", tag, err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	cid, err := cid.NewContentIDFromFile(fileReader)
	if err != nil {
		err = fmt.Errorf("Invalid CID: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}
	// Check file cid
	if cid.ToString() != offer.GetSubCID().ToString() {
		err = fmt.Errorf("Received data with wrong cid expected: %v got: %v", offer.GetSubCID().ToString(), cid.ToString())
		logging.Error(err.Error())
		// Pend PVD
		c.ReputationMgr.UpdatePVDRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendPVD(targetID)
		return nil, err
	}

	// Succeed
	c.ReputationMgr.UpdatePVDRecord(targetID, reputation.ContentRetrieved.Copy(), 0)
	return response, nil
}
