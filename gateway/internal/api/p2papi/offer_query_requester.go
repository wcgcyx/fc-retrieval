package p2papi

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cid"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrcrypto"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/common/pkg/reputation"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

func OfferQueryRequester(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 3 {
		return nil, errors.New("wrong arguments")
	}
	nodeID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	pieceCID, ok := args[1].(*cid.ContentID)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	maxOfferRequested, ok := args[2].(int64)
	if !ok {
		return nil, errors.New("wrong arguments")
	}

	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Generate random nonce
	nonce := rand.Int63()

	// Get GW Info
	gwInfo, err := c.PeerMgr.GetGWInfo(nodeID)
	if err != nil {
		// Not found, try again
		c.PeerMgr.SyncGW(nodeID)
		gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			// Not found, return error
			return nil, err
		}
	}

	// Check if this gw is blocked/pending
	rep, err := c.ReputationMgr.GetGWReputation(nodeID)
	if err != nil {
		c.ReputationMgr.AddGW(nodeID)
		rep, err = c.ReputationMgr.GetGWReputation(nodeID)
		if err != nil {
			return nil, err
		}
	}

	if rep.Pending || rep.Blocked {
		return nil, errors.New("This gateway is in pending/blocked")
	}

	// Pay the recipient
	recipientAddr, err := fcrcrypto.GetWalletAddress(gwInfo.RootKey)
	if err != nil {
		return nil, err
	}
	expected := big.Zero().Add(c.Settings.SearchPrice, big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(maxOfferRequested).Int))
	voucher, create, topup, err := c.PaymentMgr.Pay(recipientAddr, 0, expected)
	if err != nil {
		return nil, err
	}
	if create {
		// Need to create
		// First do an establishment to see if the target is alive.
		_, err := c.P2PServer.Request(gwInfo.NetworkAddr, fcrmessages.EstablishmentType, nodeID)
		if err != nil {
			return nil, err
		}
		err = c.PaymentMgr.Create(recipientAddr, c.Settings.TopupAmount)
		if err != nil {
			return nil, err
		}
		voucher, create, topup, err = c.PaymentMgr.Pay(recipientAddr, 0, expected)
		if create || topup {
			// This should never happen
			return nil, errors.New("Error, needs to create/topup channel after just creation")
		}
		if err != nil {
			return nil, err
		}
	} else if topup {
		// Need to topup
		err = c.PaymentMgr.Topup(recipientAddr, c.Settings.TopupAmount)
		if err != nil {
			return nil, err
		}
		voucher, create, topup, err = c.PaymentMgr.Pay(recipientAddr, 0, expected)
		if create || topup {
			// This should never happen
			return nil, errors.New("Error, needs to create/topup channel after just topup")
		}
		if err != nil {
			return nil, err
		}
	}
	// Now we have got a voucher
	// Encode request
	request, err := fcrmessages.EncodeStandardOfferDiscoveryRequest(false, c.NodeID, pieceCID, nonce, maxOfferRequested, c.WalletAddr, voucher)
	if err != nil {
		// This should never happen.
		logging.Error("Error in encoding request with voucher created.")
		return nil, err
	}

	// Sign request
	err = request.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		// This should never happen.
		logging.Error("Error in signing request with voucher created.")
		return nil, err
	}

	// Write request
	err = writer.Write(request, c.Settings.TCPInactivityTimeout)
	if err != nil {
		// Has error.
		c.ReputationMgr.UpdateGWRecord(nodeID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(nodeID)
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		// Has error.
		c.ReputationMgr.UpdateGWRecord(nodeID, reputation.NetworkErrorAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(nodeID)
		return nil, err
	}

	// Verify the response
	if response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
		c.PeerMgr.SyncGW(nodeID)
		gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			logging.Error("Error in getting gateway information")
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.NetworkErrorAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, err
		}
		if response.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) != nil {
			// Violation
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, errors.New("Message fail to verify")
		}
	}

	// Decode the response
	offers, nonceRecv, refundVoucher, err := fcrmessages.DecodeStandardOfferDiscoveryResponse(response)
	if err != nil {
		// Violation
		c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
		c.ReputationMgr.PendGW(nodeID)
		return nil, err
	}
	if nonceRecv != nonce {
		// Nonce mismatch
		// TODO: Warning?
		c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponse.Copy(), 0)
	}
	needToRefund := maxOfferRequested
	for _, offer := range offers {
		// Verify offer one by one
		// Get offer signing key
		pvdID := offer.GetProviderID()
		pvdInfo, err := c.PeerMgr.GetPVDInfo(pvdID)
		if err != nil {
			// Not found, try again
			c.PeerMgr.SyncPVD(nodeID)
			pvdInfo, err = c.PeerMgr.GetPVDInfo(nodeID)
			if err != nil {
				// Not found, return error
				c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
				c.ReputationMgr.PendGW(nodeID)
				return nil, err
			}
		}
		// Verify sub cid.
		if offer.GetSubCID().ToString() != pieceCID.ToString() {
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, fmt.Errorf("Offer does not contain correct cid, expect: %v, got: %v", pieceCID.ToString(), offer.GetSubCID().ToString())
		}
		// Verify offer signature
		err = offer.Verify(pvdInfo.OfferSigningKey)
		if err != nil {
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, err
		}
		// Verify offer merkle proof
		err = offer.VerifyMerkleProof()
		if err != nil {
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, err
		}
		// Offer verified
		needToRefund--
	}

	if needToRefund > 0 {
		// Check refund
		refunded, err := c.PaymentMgr.ReceiveRefund(recipientAddr, refundVoucher)
		if err != nil {
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
			return nil, err
		}
		expectedRefund := big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(needToRefund).Int)
		if refunded.Cmp(expectedRefund) < 0 {
			// Refund is wrong, but we can still respond to client, no need to return error
			c.ReputationMgr.UpdateGWRecord(nodeID, reputation.InvalidResponseAfterPayment.Copy(), 0)
			c.ReputationMgr.PendGW(nodeID)
		}
	}

	// Return response
	return response, nil
}
