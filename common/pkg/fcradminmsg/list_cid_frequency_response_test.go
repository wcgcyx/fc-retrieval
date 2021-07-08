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

func TestListCIDFrequencyResponse(t *testing.T) {
	mockCIDs := []string{"cid0", "cid1"}
	mockCount := []int{200, 100}

	data, err := EncodeListCIDFrequencyResponse(mockCIDs, mockCount)
	assert.Empty(t, err)
	assert.Equal(t, "7b2263696473223a5b2263696430222c2263696431225d2c22636f756e74223a5b3230302c3130305d7d", hex.EncodeToString(data))

	resCIDs, resCount, err := DecodeListCIDFrequencyResponse(data)
	assert.Empty(t, err)
	assert.Equal(t, mockCIDs, resCIDs)
	assert.Equal(t, mockCount, resCount)
}
