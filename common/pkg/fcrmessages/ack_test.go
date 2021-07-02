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

func TestACK(t *testing.T) {
	mockACK := true
	mockNonce := int64(42)
	mockData := "testdata"

	msg, err := EncodeACK(mockACK, mockNonce, mockData)
	assert.Empty(t, err)
	assert.Equal(t, byte(10), msg.messageType)
	assert.Equal(t, "7b2261636b223a747275652c226e6f6e6365223a34322c2264617461223a227465737464617461227d", hex.EncodeToString(msg.messageBody))
	assert.Equal(t, "", msg.signature)

	resACK, resNonce, resData, err := DecodeACK(msg)
	assert.Empty(t, err)
	assert.Equal(t, mockACK, resACK)
	assert.Equal(t, mockNonce, resNonce)
	assert.Equal(t, mockData, resData)

	msg.messageType = 100
	_, _, _, err = DecodeACK(msg)
	assert.NotEmpty(t, err)
	msg.messageType = 10

	msg.messageBody = []byte{100, 100, 100}
	_, _, _, err = DecodeACK(msg)
	assert.NotEmpty(t, err)
}
