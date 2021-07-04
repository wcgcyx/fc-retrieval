package p2papi

import (
	"encoding/hex"
	"errors"
	"math/rand"

	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
	"github.com/wcgcyx/fc-retrieval/gateway/internal/core"
)

func EstablishmentRequester(reader fcrserver.FCRServerReader, writer fcrserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 1 {
		return nil, errors.New("wrong arguments")
	}
	nodeID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("wrong arguments")
	}

	// Get core structure
	c := core.GetSingleInstance()
	c.MsgSigningKeyLock.RLock()
	defer c.MsgSigningKeyLock.RUnlock()

	challengeBytes := make([]byte, 32)
	rand.Read(challengeBytes)
	challenge := hex.EncodeToString(challengeBytes)
	request, err := fcrmessages.EncodeEstablishmentRequest(challenge)
	if err != nil {
		return nil, err
	}

	// Sign request
	err = request.Sign(c.MsgSigningKey, c.MsgSigningKeyVer)
	if err != nil {
		// This should never happen.
		logging.Error("Error in signing request.")
		return nil, err
	}

	// Write request
	err = writer.Write(request, c.Settings.TCPInactivityTimeout)
	if err != nil {
		logging.Error("Error in sending request.")
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		logging.Error("Error in reading response.")
		return nil, err
	}

	// Verify the response
	gwInfo, err := c.PeerMgr.GetGWInfo(nodeID)
	if err != nil {
		// Not found, try sync once
		c.PeerMgr.SyncGW(nodeID)
		gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			return nil, err
		}
	}
	verify := request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) == nil
	if !verify {
		// Sync the pvd once
		c.PeerMgr.SyncGW(nodeID)
		gwInfo, err = c.PeerMgr.GetGWInfo(nodeID)
		if err != nil {
			return nil, err
		}
		verify = request.Verify(gwInfo.MsgSigningKey, gwInfo.MsgSigningKeyVer) == nil
	}
	if !verify {
		return nil, errors.New("Fail to verify")
	}

	ack, _, challengeRecv, err := fcrmessages.DecodeACK(response)
	if err != nil {
		return nil, err
	}
	if !ack || challenge != challengeRecv {
		return nil, errors.New("Fail to verify")
	}

	return nil, nil
}
