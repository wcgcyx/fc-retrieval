package p2papi

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/wcgcyx/fc-retrieval/client/pkg/core"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
)

// OfferQueryRequester sends an offer query request.
func OfferQueryRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
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
	pieceCID, ok := args[1].(*cid.ContentID)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a piece CID in string")
		logging.Error(err.Error())
		return nil, err
	}
	maxOfferRequested, ok := args[2].(uint32)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a max offer requested in uint32")
		logging.Error(err.Error())
		return nil, err
	}

	// Get core structure
	c := core.GetSingleInstance()

	// Generate random nonce
	nonce := uint64(rand.Int63())

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
		err := fmt.Errorf("Gateway %v is not active", targetID)
		logging.Error(err.Error())
		return nil, err
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

	expected := big.NewInt(0).Add(c.SearchPrice, big.NewInt(0).Mul(c.OfferPrice, big.NewInt(int64(maxOfferRequested))))
	voucher, create, topup, err := c.PaymentMgr.Pay(recipientAddr, 0, expected)
	if err != nil {
		err = fmt.Errorf("Error in paying gateway %v with expected amount of %v: %v", targetID, expected.String(), err.Error())
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
	request, err := fcrmessages.EncodeStandardOfferDiscoveryRequest(nonce, c.NodeID, pieceCID, maxOfferRequested, c.WalletAddr, voucher)
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
	response, err := reader.Read(c.TCPInactivityTimeout)
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

	// Decode response
	nonceRecv, offers, refundVoucher, err := fcrmessages.DecodeStandardOfferDiscoveryResponse(response)
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
	remain := int64(maxOfferRequested)
	duplicateCheck := make(map[string]bool)
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
		// Check duplicates
		_, ok := duplicateCheck[offer.GetMessageDigest()]
		if ok {
			err = fmt.Errorf("Received duplicated offers")
			logging.Error(err.Error())
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
			return nil, err
		}
		// Check offer expiry, reject if less than 1 hour
		if offer.GetExpiry()-time.Now().Unix() < 3600 {
			// Offer is soon to expire
			err = fmt.Errorf("Received soon to expire offer")
			logging.Error(err.Error())
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
			return nil, err
		}
		// Offer verified
		remain--
		c.OfferMgr.AddSubOffer(&offer)
		c.ReputationMgr.UpdateGWRecord(targetID, reputation.StandardOfferRetrieved.Copy(), 0)
	}

	if remain > 0 {
		// Check refund
		refunded, err := c.PaymentMgr.ReceiveRefund(recipientAddr, refundVoucher)
		if err != nil {
			// Refund is wrong, but we can still respond to client, no need to return error
			c.ReputationMgr.UpdateGWRecord(targetID, reputation.InvalidRefund.Copy(), 0)
			c.ReputationMgr.PendGW(targetID)
		} else {
			expectedRefund := big.NewInt(0).Mul(c.OfferPrice, big.NewInt(remain))
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
