package p2papi

import (
	"fmt"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

func EstablishmentHandler(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	// Response
	var response *fcrmessages.FCRMessage

	// Message decoding
	challenge, err := fcrmessages.DecodeEstablishmentRequest(request)
	if err != nil {
		response, _ = fcrmessages.EncodeACK(false, 0, fmt.Sprintf("Error in decoding payload: %v", err.Error()))
		response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}

	// Respond
	response, err = fcrmessages.EncodeACK(false, 0, challenge)
	if err != nil {
		response, _ = fcrmessages.EncodeACK(false, 0, fmt.Sprintf("Internal error in encoding response: %v", err.Error()))
		return writer.Write(response, c.Settings.TCPInactivityTimeout)
	}
	err = response.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		logging.Error("Error in signing response: %v", err.Error())
	}

	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
