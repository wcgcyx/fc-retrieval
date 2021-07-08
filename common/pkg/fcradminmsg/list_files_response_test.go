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

func TestListFilesResponse(t *testing.T) {
	mockFiles := []string{"file0", "file1"}
	mockCIDs := []string{"cid0", "cid1"}
	mockSizes := []int64{1000, 2000}
	mockPublished := []bool{true, false}
	mockFrequency := []int{100, 200}

	data, err := EncodeListFilesResponse(mockFiles, mockCIDs, mockSizes, mockPublished, mockFrequency)
	assert.Empty(t, err)
	assert.Equal(t, "7b2266696c6573223a5b2266696c6530222c2266696c6531225d2c2263696473223a5b2263696430222c2263696431225d2c227075626c6973686564223a5b747275652c66616c73655d2c226672657175656e6379223a5b3130302c3230305d7d", hex.EncodeToString(data))

	resFiles, resCIDs, resSizes, resPublished, resFrequency, err := DecodeListFilesResponse(data)
	assert.Empty(t, err)
	assert.Equal(t, mockFiles, resFiles)
	assert.Equal(t, mockCIDs, resCIDs)
	assert.Equal(t, mockSizes, resSizes)
	assert.Equal(t, mockPublished, resPublished)
	assert.Equal(t, mockFrequency, resFrequency)
}
