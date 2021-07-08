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

func TestListPeersResponse(t *testing.T) {
	mockGWIDs := []string{"id0", "id1"}
	mockGWScore := []int64{100, 200}
	mockGWPending := []bool{true, false}
	mockGWBlocked := []bool{false, true}
	mockGWRecent := []string{"recent0", "recent1"}
	mockPVDIDs := []string{"id2", "id3"}
	mockPVDScore := []int64{101, 201}
	mockPVDPending := []bool{false, true}
	mockPVDBlocked := []bool{true, false}
	mockPVDRecent := []string{"recent2", "recent3"}

	data, err := EncodeListPeersResponse(mockGWIDs, mockGWScore, mockGWPending, mockGWBlocked, mockGWRecent, mockPVDIDs, mockPVDScore, mockPVDPending, mockPVDBlocked, mockPVDRecent)
	assert.Empty(t, err)
	assert.Equal(t, "7b22676174657761795f696473223a5b22696430222c22696431225d2c22676174657761795f73636f7265223a5b3130302c3230305d2c22676174657761795f70656e64696e67223a5b747275652c66616c73655d2c22676174657761795f626c6f636b6564223a5b66616c73652c747275655d2c22676174657761795f726563656e74223a5b22726563656e7430222c22726563656e7431225d2c2270726f76696465725f696473223a5b22696432222c22696433225d2c2270726f76696465725f73636f7265223a5b3130312c3230315d2c2270726f76696465725f70656e64696e67223a5b66616c73652c747275655d2c2270726f76696465725f626c6f636b6564223a5b747275652c66616c73655d2c2270726f76696465725f726563656e74223a5b22726563656e7432222c22726563656e7433225d7d", hex.EncodeToString(data))

	resGWIDs, resGWScore, resGWPending, resGWBlocked, resGWRecent, resPVDIDs, resPVDScore, resPVDPending, resPVDBlocked, resPVDRecent, err := DecodeListPeersResponse(data)
	assert.Empty(t, err)
	assert.Equal(t, mockGWIDs, resGWIDs)
	assert.Equal(t, mockGWScore, resGWScore)
	assert.Equal(t, mockGWPending, resGWPending)
	assert.Equal(t, mockGWBlocked, resGWBlocked)
	assert.Equal(t, mockGWRecent, resGWRecent)
	assert.Equal(t, mockPVDIDs, resPVDIDs)
	assert.Equal(t, mockPVDScore, resPVDScore)
	assert.Equal(t, mockPVDPending, resPVDPending)
	assert.Equal(t, mockPVDBlocked, resPVDBlocked)
	assert.Equal(t, mockPVDRecent, resPVDRecent)
}
