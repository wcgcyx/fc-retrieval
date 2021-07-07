package p2papi

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/wcgcyx/fc-retrieval/client/pkg/core"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

// DHTOfferQueryRequester sends an offer query request.
func DHTOfferQueryRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
	// Get parameters
	if len(args) != 4 {
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
	pieceCID, ok := args[1].(*cid.ContentID)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a piece CID in string")
		logging.Error(err.Error())
		return nil, err
	}
	numDHT, ok := args[2].(uint32)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a numDHT in uint32")
		logging.Error(err.Error())
		return nil, err
	}
	maxOfferRequestedPerDHT, ok := args[2].(uint32)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a max offer requested per DHT in uint32")
		logging.Error(err.Error())
		return nil, err
	}

	// Get core structure
	c := core.GetSingleInstance()

	// Generate random nonce
	nonce := uint64(rand.Int63())

	// Maximum numDHT is 16.
	if numDHT > 16 {
		numDHT = 16
	}

	// Get gateway information
	gwInfo := c.PeerMgr.GetGWInfo(targetID)
	if gwInfo == nil {
		// Not found, try sync once
		gwInfo = c.PeerMgr.SyncGW(targetID)
		if gwInfo == nil {
			err := fmt.Errorf("Error in obtaining information for gateway %v", targetID)
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Check if the gateway is blocked/pending
	rep := c.ReputationMgr.GetGWReputation(targetID)
	if rep == nil {
		c.ReputationMgr.AddGW(targetID)
		rep = c.ReputationMgr.GetGWReputation(targetID)
	}
	if rep.Pending || rep.Blocked {
		err := fmt.Errorf("Gateway %v is in pending %v, blocked %v", targetID, rep.Pending, rep.Blocked)
		logging.Error(err.Error())
		return nil, err
	}

	// Pay the recipient
	recipientAddr, err := fcrcrypto.GetWalletAddress(gwInfo.RootKey)
	if err != nil {
		err = fmt.Errorf("Error in obtaining wallet addreess for gateway %v with root key %v: %v", targetID, gwInfo.RootKey, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// expected is 1 * search price + numDHT * (search price + max offer per DHT * offer price)
	expected := big.NewInt(0).Add(c.SearchPrice, big.NewInt(0).Mul(big.NewInt(0).Add(c.SearchPrice, big.NewInt(0).Mul(c.OfferPrice, big.NewInt(int64(maxOfferRequestedPerDHT)))), big.NewInt(int64(numDHT))))
	voucher, _, topup, err := c.PaymentMgr.Pay(recipientAddr, 0, expected)
	if err != nil {
		err = fmt.Errorf("Error in paying gateway %v with expected amount of %v: %v", targetID, expected.String(), err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	if topup {
		// Need to topup
		err = c.PaymentMgr.Topup(recipientAddr, c.TopupAmount)
		if err != nil {
			err = fmt.Errorf("Error in topup a payment channel to %v with wallet address %v with topup amount of %v: %v", targetID, recipientAddr, c.TopupAmount.String(), err.Error())
			logging.Error(err.Error())
			return nil, err
		}
		voucher, _, topup, err = c.PaymentMgr.Pay(recipientAddr, 0, expected)
		if topup {
			// This should never happen
			err = fmt.Errorf("Error in paying gateway %v, needs to create/topup after just topup", targetID)
			logging.Error(err.Error())
			return nil, err
		}
		if err != nil {
			err = fmt.Errorf("Error in paying gateway %v with expected amount of %v: %v after just topup", targetID, expected.String(), err.Error())
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Now we have got a voucher
	// Encode request
	request, err := fcrmessages.EncodeDHTOfferDiscoveryRequest(nonce, c.NodeID, pieceCID, numDHT, maxOfferRequestedPerDHT, c.WalletAddr, voucher)
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
		// Pend GW
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(targetID)
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.LongTCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in receiving response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		// Pend GW
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(targetID)
		return nil, err
	}

	// Verify the response
	if response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
		// Try update
		gwInfo = c.PeerMgr.SyncGW(targetID)
		if gwInfo == nil || response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
			err = fmt.Errorf("Error in verifying response from %v: %v", targetID, err.Error())
			logging.Error(err.Error())
			// Pend GW
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
			return nil, err
		}
	}

	// Check response
	if !response.ACK() {
		err = fmt.Errorf("Reponse contains an error: %v", response.Error())
		logging.Error(err.Error())
		// Pend GW
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(targetID)
		return nil, err
	}

	nonceRecv, contacted, refundVoucher, err := fcrmessages.DecodeDHTOfferDiscoveryResponse(response)
	if err != nil {
		err = fmt.Errorf("Error in decoding response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		// Pend GW
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(targetID)
		return nil, err
	}

	if nonceRecv != nonce {
		err = fmt.Errorf("Nonce mismatch: expected %v got %v", nonce, nonceRecv)
		logging.Error(err.Error())
		// Pend GW
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(targetID)
		return nil, err
	}

	// Check payment and offer
	remainTotal := int64(numDHT * maxOfferRequestedPerDHT)
	for subID, resp := range contacted {
		// Verify the response one by one
		subGWInfo := c.PeerMgr.GetGWInfo(subID)
		if subGWInfo == nil {
			// Not found, try sync once
			subGWInfo = c.PeerMgr.SyncGW(subID)
			if subGWInfo == nil {
				err := fmt.Errorf("Error in obtaining information for sub gateway %v", subID)
				logging.Error(err.Error())
				// Pend GW
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
				return nil, err
			}
		}
		if resp.Verify(subGWInfo.MsgSigningKey, subGWInfo.MsgSigningKeyVer) != nil {
			// Try update
			subGWInfo = c.PeerMgr.SyncGW(subID)
			if subGWInfo == nil || resp.Verify(subGWInfo.MsgSigningKey, subGWInfo.MsgSigningKeyVer) != nil {
				err = fmt.Errorf("Error in verifying sub response from %v: %v", subID, err.Error())
				logging.Error(err.Error())
				// Pend GW
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
				return nil, err
			}
		}
		// Check response
		if !resp.ACK() {
			err = fmt.Errorf("Reponse contains an error: %v", resp.Error())
			logging.Error(err.Error())
			// Pend GW
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
			return nil, err
		}
		// Decode response
		_, offers, _, err := fcrmessages.DecodeStandardOfferDiscoveryResponse(&resp)
		if err != nil {
			err = fmt.Errorf("Error in decoding sub response from %v: %v", subID, err.Error())
			logging.Error(err.Error())
			// Pend GW
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
			return nil, err
		}

		// Check offer
		remainSub := int(maxOfferRequestedPerDHT)
		for _, offer := range offers {
			// Verify offer one by one
			// Get offer signing key
			pvdID := offer.GetProviderID()
			pvdInfo := c.PeerMgr.GetPVDInfo(pvdID)
			if err == nil {
				// Not found, try again
				pvdInfo = c.PeerMgr.SyncPVD(pvdID)
				if pvdInfo == nil {
					// Not found, return error
					err = fmt.Errorf("Error in obtaining information for provider %v", pvdID)
					logging.Error(err.Error())
					c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
					c.ReputationMgr.PendGW(targetID)
					return nil, err
				}
			}
			// Verify sub cid.
			if offer.GetSubCID().ToString() != pieceCID.ToString() {
				err = fmt.Errorf("Received offer that doesn't contain requested cid, expect: %v, got: %v", pieceCID.ToString(), offer.GetSubCID().ToString())
				logging.Error(err.Error())
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
				return nil, err
			}
			// Verify offer signature
			if offer.Verify(pvdInfo.OfferSigningKey) != nil {
				err = fmt.Errorf("Received offer fails to verify against signature of provider %v", pvdID)
				logging.Error(err.Error())
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
				return nil, err
			}
			// Verify offer merkle proof
			if offer.VerifyMerkleProof() != nil {
				err = fmt.Errorf("Received offer fails to verify merkle proof")
				logging.Error(err.Error())
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
				return nil, err
			}

			// Offer verified
			remainSub--
			c.OfferMgr.AddSubOffer(&offer)
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.DHTOfferRetrieved.Copy(), 0)
		}
		if remainSub < 0 {
			remainSub = 0
		}
		remainTotal -= (int64(maxOfferRequestedPerDHT) - int64(remainSub))
	}

	// Check remain total
	if remainTotal > 0 {
		// Check refund
		refunded, err := c.PaymentMgr.ReceiveRefund(recipientAddr, refundVoucher)
		if err != nil {
			// Refund is wrong, but we can still respond to client, no need to return error
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidRefund.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
		} else {
			expectedRefund := big.NewInt(0).Mul(c.OfferPrice, big.NewInt(remainTotal))
			if refunded.Cmp(expectedRefund) < 0 {
				// Refund is wrong, but we can still respond to client, no need to return error
				c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidRefund.Copy(), 0)
				c.ReputationMgr.PendGW(targetID)
			}
		}
	}

	// Return response
	return response, nil
}
