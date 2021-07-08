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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangePeerStatusRequest(t *testing.T) {
	mockID := "testid"
	mockGateway := true
	mockBlock := true
	mockUnblock := false

	data, err := EncodeChangePeerStatusRequest(mockID, mockGateway, mockBlock, mockUnblock)
	assert.Empty(t, err)
	assert.Equal(t, "7b226e6f64655f6964223a22746573746964222c2267617465776179223a747275652c22626c6f636b223a747275652c22756e626c6f636b223a66616c73657d", hex.EncodeToString(data))

	resID, resGateway, resBlock, resUnblock, err := DecodeChangePeerStatusRequest(data)
	assert.Empty(t, err)
	assert.Equal(t, mockID, resID)
	assert.Equal(t, mockGateway, resGateway)
	assert.Equal(t, mockBlock, resBlock)
	assert.Equal(t, mockUnblock, resUnblock)
}
