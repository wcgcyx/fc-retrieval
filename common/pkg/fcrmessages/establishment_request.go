/*
Package fcrmessages - stores all the p2p messages.
*/
package fcrmessages

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
	"encoding/json"
	"fmt"
)

// establishmentRequestJson represents an establishment.
type establishmentRequestJson struct {
	NodeID    string `json:"node_id"`
	Challenge string `json:"challenge"`
}

// EncodeEstablishmentRequest is used to get the FCRMessage of establishmentRequestJson
func EncodeEstablishmentRequest(
	nonce uint64,
	nodeID string,
	challenge string,
) (*FCRReqMsg, error) {
	body, err := json.Marshal(establishmentRequestJson{
		NodeID:    nodeID,
		Challenge: challenge,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRReqMsg(EstablishmentRequestType, nonce, body), nil
}

// DecodeEstablishmentRequest is used to get the fields from FCRMessage of establishmentRequestJson
// It returns the nonce and challenge string in this establishment request.
func DecodeEstablishmentRequest(fcrMsg *FCRReqMsg) (
	uint64,
	string,
	string,
	error,
) {
	if fcrMsg.Type() != EstablishmentRequestType {
		return 0, "", "", fmt.Errorf("Message type mismatch, expect %v, got %v", EstablishmentRequestType, fcrMsg.Type())
	}
	msg := establishmentRequestJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", "", err
	}
	return fcrMsg.Nonce(), msg.NodeID, msg.Challenge, nil
}
