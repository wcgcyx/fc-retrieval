package p2papi

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/wcgcyx/fc-retrieval/common/pkg/cidoffer"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

func OfferQueryHandler(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Response
	var response *fcrmessages.FCRMessage

	// Message decoding
	client, nodeID, pieceCID, nonce, maxOfferRequested, accountAddr, voucher, err := fcrmessages.DecodeStandardOfferDiscoveryRequest(request)
	if err != nil {
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in decoding payload: %v", err.Error()))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	var verify bool
	if client {
		verify = request.VerifyByID(nodeID) == nil
	} else {
		// Get GW Info
		gwInfo, err := c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			// Not found, try again
			c.PeerMgr.SyncGW(nodeID)
			gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
			if err != nil {
				// Not found, return error
				response, _ = fcrmessages.EncodeACK(false, nonce, "Error in finding gateway infomation")
				response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
				return writer.Write(response, c.Settings.TCPInactivityTimeout)
			}
		}
		verify = request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) == nil
		if !verify {
			// Try again
			c.PeerMgr.SyncGW(nodeID)
			gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
			if err != nil {
				// Not found, return error
				response, _ = fcrmessages.EncodeACK(false, nonce, "Error in finding gateway infomation")
				response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
				return writer.Write(response, c.Settings.TCPInactivityTimeout)
			}
			verify = request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) == nil
		}
	}
	if !verify {
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
	expected := big.Zero().Add(c.Settings.SearchPrice, big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(maxOfferRequested).Int))
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

	// Payment is fine, search.
	refundVoucher := ""
	offers := c.OfferMgr.GetOffers(pieceCID)

	// Generating sub CID offers
	res := make([]cidoffer.SubCIDOffer, 0)
	toRefund := maxOfferRequested
	for _, offer := range offers {
		if toRefund == 0 {
			break
		}
		subOffer, err := offer.GenerateSubCIDOffer(pieceCID)
		if err != nil {
			// Internal error in generating sub offers
			refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, received)
			if err != nil {
				// This should never happen
				logging.Error("Error in refunding %v", err.Error())
			}
			response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Internal error, refund voucher: %v", refundVoucher))
			response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
			return writer.Write(response, c.Settings.TCPInactivityTimeout)
		}
		res = append(res, *subOffer)
		toRefund--
	}
	if toRefund > 0 {
		refundVoucher, err = c.PaymentMgr.Refund(accountAddr, lane, big.Zero().Sub(received, big.Zero().Mul(c.Settings.OfferPrice, big.NewInt(toRefund).Int)))
		if err != nil {
			// This should never happen
			logging.Error("Error in refunding %v", err.Error())
		}
	}

	response, err = fcrmessages.EncodeStandardOfferDiscoveryResponse(res, nonce, refundVoucher)
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
