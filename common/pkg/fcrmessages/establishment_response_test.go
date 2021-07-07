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

func TestEstablishmentResponse(t *testing.T) {
	mockNonce := uint64(100)
	mockChallenge := "test challenge"

	msg, err := EncodeEstablishmentResponse(mockNonce, mockChallenge)
	assert.Empty(t, err)
	assert.Equal(t, true, msg.ack)
	assert.Equal(t, uint64(100), msg.nonce)
	assert.Equal(t, "7b226e6f64655f6964223a22222c226368616c6c656e6765223a2274657374206368616c6c656e6765227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resNonce, resChallenge, err := DecodeEstablishmentResponse(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockChallenge, resChallenge)

	msg.ack = false
	_, _, err = DecodeEstablishmentResponse(msg)
	assert.NotEmpty(t, err)
	msg.ack = true

	msg.messageBody = []byte{100, 100, 100}
	_, _, err = DecodeEstablishmentResponse(msg)
	assert.NotEmpty(t, err)
}
