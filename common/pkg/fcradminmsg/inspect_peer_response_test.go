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

func TestInspectPeerResponse(t *testing.T) {
	mockScore := int64(156)
	mockPending := true
	mockBlocked := true
	mockHistory := []string{"good", "good2", "good3"}

	data, err := EncodeInspectPeerResponse(mockScore, mockPending, mockBlocked, mockHistory)
	assert.Empty(t, err)
	assert.Equal(t, "7b2273636f7265223a3135362c2270656e64696e67223a747275652c22626c6f636b6564223a747275652c22686973746f7279223a5b22676f6f64222c22676f6f6432222c22676f6f6433225d7d", hex.EncodeToString(data))

	resScore, resPending, resBlocked, resHistory, err := DecodeInspectPeerResponse(data)
	assert.Empty(t, err)
	assert.Equal(t, mockScore, resScore)
	assert.Equal(t, mockBlocked, resBlocked)
	assert.Equal(t, mockPending, resPending)
	assert.Equal(t, mockHistory, resHistory)
}
