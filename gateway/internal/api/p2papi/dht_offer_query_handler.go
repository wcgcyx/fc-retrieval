package p2papi

import (
	"encoding/hex"
	"fmt"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

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
	expected := big.Zero().Add(c.Settings.SearchPrice, big.Zero().Mul(big.Zero().Add(c.Settings.SearchPrice, big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(maxOfferRequestedPerDHT).Int)), big.NewInt(numDHT).Int))
	if expected.Cmp(received) < 0 {
		// Short payment
		voucher, err := c.PaymentMgr.Refund(accountAddr, lane, big.Zero().Sub(received, c.Settings.SearchPrice))
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
		refundVoucher, err := c.PaymentMgr.Refund(accountAddr, lane, big.Zero().Sub(received, c.Settings.SearchPrice))
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
	supposed := big.Zero().Set(c.Settings.SearchPrice)
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
		supposed.Add(c.Settings.SearchPrice, big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(found).Int))
		contacted[gw.NodeID] = resp
	}
	if supposed.Cmp(expected) < 0 {
		refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, big.Zero().Sub(expected, supposed))
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
