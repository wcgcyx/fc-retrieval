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

func TestACK(t *testing.T) {
	mockACK := true
	mockData := "testdata"

	data, err := EncodeACK(mockACK, mockData)
	assert.Empty(t, err)
	assert.Equal(t, "7b2261636b223a747275652c2264617461223a227465737464617461227d", hex.EncodeToString(data))

	resACK, resData, err := DecodeACK(data)
	assert.Empty(t, err)
	assert.Equal(t, mockACK, resACK)
	assert.Equal(t, mockData, resData)

	_, _, err = DecodeACK([]byte{100, 100, 100})
	assert.NotEmpty(t, err)
}
