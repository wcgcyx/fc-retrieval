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

func TestDataRetrievalResponse(t *testing.T) {
	mockNonce := uint64(100)
	mockTag := "mocktag"
	mockData := []byte{1, 2, 3}

	msg, err := EncodeDataRetrievalResponse(mockNonce, mockTag, mockData)
	assert.Empty(t, err)
	assert.Equal(t, true, msg.ack)
	assert.Equal(t, uint64(100), msg.nonce)
	assert.Equal(t, "7b22746167223a226d6f636b746167222c2264617461223a2241514944227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resNonce, resTag, resData, err := DecodeDataRetrievalResponse(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockTag, resTag)
	assert.Equal(t, mockData, resData)

	msg.ack = false
	_, _, _, err = DecodeDataRetrievalResponse(msg)
	assert.NotEmpty(t, err)
	msg.ack = true

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, err = DecodeDataRetrievalResponse(msg)
	assert.NotEmpty(t, err)
}
