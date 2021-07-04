package p2papi

import (
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

func OfferPublishHandler(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get core response
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Response
	var response *fcrmessages.FCRMessage

	nodeID, nonce, offer, err := fcrmessages.DecodeOfferPublishRequest(request)
	if err != nil {
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in decoding payload: %v", err.Error()))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// First verify the message
	pvdInfo, err := c.PeerMgr.GetPVDInfo(nodeID)
	if err != nil {
		// Not found, try sync once
		c.PeerMgr.SyncPVD(nodeID)
		pvdInfo, err = c.PeerMgr.GetPVDInfo(nodeID)
		if err != nil {
			response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in getting provider info: %v", err.Error()))
			return writer.Write(response, c.Settings.TCPInactivityTimeout)
		}
	}
	verify := request.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) == nil
	if !verify {
		// Sync the pvd once
		c.PeerMgr.SyncPVD(nodeID)
		pvdInfo, err = c.PeerMgr.GetPVDInfo(nodeID)
		if err != nil {
			response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in getting provider information: %v", err.Error()))
			return writer.Write(response, c.Settings.TCPInactivityTimeout)
		}
		verify = request.Verify(pvdInfo.MsgSigningKey, pvdInfo.MsgSigningKeyVer) == nil
	}
	if !verify {
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in verifying msg: %v", err.Error()))
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// Check offer signature
	verify = offer.Verify(pvdInfo.OfferSigningKey) == nil
	if !verify {
		response, _ = fcrmessages.EncodeACK(false, nonce, fmt.Sprintf("Error in verifying offer: %v", err.Error()))
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// Offer verified, add to storage
	// TODO: c.StoreFullOffer
	c.OfferMgr.AddOffer(offer)

	response, _ = fcrmessages.EncodeACK(true, nonce, "Offer added")
	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
