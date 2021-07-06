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

type establishmentResponseJson struct {
	Challenge string `json:"challenge"`
}

// EncodeEstablishmentResponse is used to get the FCRMessage of establishmentResponseJson
func EncodeEstablishmentResponse(
	nonce uint64,
	challenge string,
) (*FCRACKMsg, error) {
	body, err := json.Marshal(establishmentRequestJson{
		Challenge: challenge,
	})
	if err != nil {
		return nil, err
	}
	return CreateFCRACKMsg(nonce, body), nil
}

// DecodeEstablishmentResponse is used to get the fields from FCRMessage of establishmentResponseJson
// It returns the nonce and challenge string in this establishment request.
func DecodeEstablishmentResponse(fcrMsg *FCRACKMsg) (
	uint64,
	string,
	error,
) {
	if !fcrMsg.ACK() {
		return 0, "", fmt.Errorf("ACK is false")
	}
	msg := establishmentResponseJson{}
	err := json.Unmarshal(fcrMsg.Body(), &msg)
	if err != nil {
		return 0, "", err
	}
	return fcrMsg.Nonce(), msg.Challenge, nil
}
