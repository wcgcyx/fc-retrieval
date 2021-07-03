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
	Challenge string `json:"challenge"`
}

// EncodeEstablishmentRequest is used to get the FCRMessage of establishmentRequestJson
func EncodeEstablishmentRequest(
	challenge string,
) (*FCRMessage, error) {
	body, err := json.Marshal(establishmentRequestJson{
		Challenge: challenge,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRMessage(EstablishmentType, body), nil
}

// DecodeEstablishmentRequest is used to get the fields from FCRMessage of establishmentRequestJson
// It returns the challenge string in this establishment request.
func DecodeEstablishmentRequest(fcrMsg *FCRMessage) (
	string,
	error,
) {
	if fcrMsg.GetMessageType() != EstablishmentType {
		return "", fmt.Errorf("Message type mismatch, expect %v, got %v", EstablishmentType, fcrMsg.GetMessageType())
	}
	msg := establishmentRequestJson{}
	err := json.Unmarshal(fcrMsg.GetMessageBody(), &msg)
	if err != nil {
		return "", err
	}
	return msg.Challenge, nil
}
