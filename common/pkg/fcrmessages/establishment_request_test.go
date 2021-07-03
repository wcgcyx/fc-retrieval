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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstablishment(t *testing.T) {
	mockChallenge := "test challenge"

	msg, err := EncodeEstablishmentRequest(mockChallenge)
	assert.Empty(t, err)
	assert.Equal(t, byte(11), msg.messageType)
	assert.Equal(t, "7b226368616c6c656e6765223a2274657374206368616c6c656e6765227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resChallenge, err := DecodeEstablishmentRequest(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockChallenge, resChallenge)

	msg.messageType = 100
	_, err = DecodeEstablishmentRequest(msg)
	assert.NotEmpty(t, err)
	msg.messageType = 11

	msg.messageBody = []byte{100, 100, 100}
	_, err = DecodeEstablishmentRequest(msg)
	assert.NotEmpty(t, err)
}
