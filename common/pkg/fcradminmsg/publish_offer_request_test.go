/*
Package fcradminmsg - stores all the admin messages.
*/
package fcradminmsg

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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishOfferRequest(t *testing.T) {
	mockFiles := []string{"test1", "test2"}
	mockPrice := big.NewInt(100)
	mockExpiry := int64(1000)
	mockQoS := uint64(44)

	data, err := EncodePublishOfferRequest(mockFiles, mockPrice, mockExpiry, mockQoS)
	assert.Empty(t, err)
	assert.Equal(t, "7b2266696c6573223a5b227465737431222c227465737432225d2c227072696365223a22313030222c22657870697279223a313030302c22716f73223a34347d", hex.EncodeToString(data))

	resFiles, resPrice, resExpiry, resQoS, err := DecodePublishOfferRequest(data)
	assert.Empty(t, err)
	assert.Equal(t, mockFiles, resFiles)
	assert.Equal(t, mockPrice.String(), resPrice.String())
	assert.Equal(t, mockExpiry, resExpiry)
	assert.Equal(t, mockQoS, resQoS)

	_, _, err = DecodeACK([]byte{100, 100, 100})
	assert.NotEmpty(t, err)
}
