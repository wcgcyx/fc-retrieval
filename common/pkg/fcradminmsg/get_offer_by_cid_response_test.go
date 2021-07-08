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

func TestGetOfferByCIDResponse(t *testing.T) {
	mockDigest := []string{"testdigest0", "testdigest1"}
	mockProviders := []string{"provider0", "provider1"}
	mockPrices := []string{"1000", "2000"}
	mockExpiry := []int64{123, 234}
	mockQoS := []uint64{345, 456}

	data, err := EncodeGetOfferByCIDResponse(mockDigest, mockProviders, mockPrices, mockExpiry, mockQoS)
	assert.Empty(t, err)
	assert.Equal(t, "7b2244696765737473223a5b227465737464696765737430222c227465737464696765737431225d2c2250726f766964657273223a5b2270726f766964657230222c2270726f766964657231225d2c22507269636573223a5b2231303030222c2232303030225d2c22457870697279223a5b3132332c3233345d2c22516f53223a5b3334352c3435365d7d", hex.EncodeToString(data))

	resDigest, resProviders, resPrices, resExpiry, resQoS, err := DecodeGetOfferByCIDResponse(data)
	assert.Empty(t, err)
	assert.Equal(t, mockDigest, resDigest)
	assert.Equal(t, mockProviders, resProviders)
	assert.Equal(t, mockPrices, resPrices)
	assert.Equal(t, mockExpiry, resExpiry)
	assert.Equal(t, mockQoS, resQoS)
}
