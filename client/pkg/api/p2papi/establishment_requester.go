package p2papi

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"github.com/wcgcyx/fc-retrieval/client/pkg/core"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrmessages"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrpeermgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrserver"
	"github.com/wcgcyx/fc-retrieval/common/pkg/logging"
)

// EstablishmentRequester sends an establishment request.
func EstablishmentRequester(reader fcrserver.FCRServerResponseReader, writer fcrserver.FCRServerRequestWriter, args ...interface{}) (*fcrmessages.FCRACKMsg, error) {
	// Get parameters
	if len(args) != 2 {
		err := fmt.Errorf("Wrong arguments, expect length 2, got length %v", len(args))
		logging.Error(err.Error())
		return nil, err
	}
	targetID, ok := args[0].(string)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a target ID in string")
		logging.Error(err.Error())
		return nil, err
	}
	gateway, ok := args[1].(bool)
	if !ok {
		err := fmt.Errorf("Wrong arguments, expect a gateway in bool")
		logging.Error(err.Error())
		return nil, err
	}

	// Get core structure
	c := core.GetSingleInstance()

	// Generate random nonce
	nonce := uint64(rand.Int63())

	challengeBytes := make([]byte, 32)
	rand.Read(challengeBytes)
	challenge := hex.EncodeToString(challengeBytes)
	request, err := fcrmessages.EncodeEstablishmentRequest(nonce, c.NodeID, challenge)
	if err != nil {
		err = fmt.Errorf("Internal error in encoding response: %v", err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Write request
	err = writer.Write(request, c.MsgKey, 0, c.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in sending request to %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.TCPInactivityTimeout)
	if err != nil {
		err = fmt.Errorf("Error in receiving response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Verify the response
	var nodeInfo *fcrpeermgr.Peer
	if gateway {
		nodeInfo = c.PeerMgr.GetGWInfo(targetID)
		if nodeInfo == nil {
			// Not found, try sync once
			nodeInfo = c.PeerMgr.SyncGW(targetID)
			if nodeInfo == nil {
				err = fmt.Errorf("Error in obtaining information for gateway %v", targetID)
				logging.Error(err.Error())
				return nil, err
			}
		}
	} else {
		nodeInfo = c.PeerMgr.GetPVDInfo(targetID)
		if nodeInfo == nil {
			// Not found, try sync once
			nodeInfo = c.PeerMgr.SyncPVD(targetID)
			if nodeInfo == nil {
				err = fmt.Errorf("Error in obtaining information for provider %v", targetID)
				logging.Error(err.Error())
				return nil, err
			}
		}
	}
	if response.Verify(nodeInfo.MsgSigningKey, nodeInfo.MsgSigningKeyVer) != nil {
		// Try update
		if gateway {
			nodeInfo = c.PeerMgr.SyncGW(targetID)
		} else {
			nodeInfo = c.PeerMgr.SyncPVD(targetID)
		}
		if nodeInfo == nil || response.Verify(nodeInfo.MsgSigningKey, nodeInfo.MsgSigningKeyVer) != nil {
			err = fmt.Errorf("Error in verifying response from %v: %v", targetID, err.Error())
			logging.Error(err.Error())
			return nil, err
		}
	}

	// Check response
	if !response.ACK() {
		err = fmt.Errorf("Reponse contains an error: %v", response.Error())
		logging.Error(err.Error())
		return nil, err
	}

	// Decode response
	nonceRecv, challengeRecv, err := fcrmessages.DecodeEstablishmentResponse(response)
	if err != nil {
		err = fmt.Errorf("Error in decoding response from %v: %v", targetID, err.Error())
		logging.Error(err.Error())
		return nil, err
	}

	if nonceRecv != nonce {
		err = fmt.Errorf("Nonce mismatch: expected %v got %v", nonce, nonceRecv)
		logging.Error(err.Error())
		return nil, err
	}

	if challengeRecv != challenge {
		err = fmt.Errorf("Challenge mismatch: expected %v got %v", challenge, challengeRecv)
		logging.Error(err.Error())
		return nil, err
	}

	return response, nil
}
