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

func TestInitialisationRequest(t *testing.T) {
	mockP2PPrvKey := "p2pprivatekey"
	mockP2PPort := 60
	mockNetworkAddr := "testaddr"
	mockRootPrvKey := "testkey"
	mockLotusAPIAddr := "lotusaddr"
	mockLotusAuthToken := "lotusauthtoken"
	mockRegisterPrvKey := "registerprvkey"
	mockRegisterAPIAddr := "testregisteraddr"
	mockRegisterAuthToken := "registerauthtoken"
	mockRegionCode := "testregioncode"

	data, err := EncodeInitialisationRequest(mockP2PPrvKey, mockP2PPort, mockNetworkAddr, mockRootPrvKey, mockLotusAPIAddr, mockLotusAuthToken, mockRegisterPrvKey, mockRegisterAPIAddr, mockRegisterAuthToken, mockRegionCode)
	assert.Empty(t, err)
	assert.Equal(t, "7b227032705f707269766174655f6b6579223a22703270707269766174656b6579222c227032705f706f7274223a36302c226e6574776f726b5f61646472223a227465737461646472222c22726f6f745f707269766174655f6b6579223a22746573746b6579222c226c6f7475735f6170695f61646472223a226c6f74757361646472222c226c6f7475735f617574685f746f6b656e223a226c6f74757361757468746f6b656e222c2272656769737465725f707269766174655f6b6579223a2272656769737465727072766b6579222c2272656769737465725f6170695f61646472223a2274657374726567697374657261646472222c2272656769737465725f617574685f746f6b656e223a22726567697374657261757468746f6b656e222c22726567696f6e5f636f6465223a2274657374726567696f6e636f6465227d", hex.EncodeToString(data))

	resP2PPrvKey, resP2PPort, resNetworkAddr, resRootPrvKey, resLotusAPIAddr, resLotusAuthToken, resRegisterPrvKey, resRegisterAPIAddr, resRegisterAuthToken, resRegionCode, err := DecodeInitialisationRequest(data)
	assert.Empty(t, err)
	assert.Equal(t, mockP2PPrvKey, resP2PPrvKey)
	assert.Equal(t, mockP2PPort, resP2PPort)
	assert.Equal(t, mockNetworkAddr, resNetworkAddr)
	assert.Equal(t, mockRootPrvKey, resRootPrvKey)
	assert.Equal(t, mockLotusAPIAddr, resLotusAPIAddr)
	assert.Equal(t, mockLotusAuthToken, resLotusAuthToken)
	assert.Equal(t, mockRegisterPrvKey, resRegisterPrvKey)
	assert.Equal(t, mockRegisterAPIAddr, resRegisterAPIAddr)
	assert.Equal(t, mockRegisterAuthToken, resRegisterAuthToken)
	assert.Equal(t, mockRegionCode, resRegionCode)

	_, _, _, _, _, _, _, _, _, _, err = DecodeInitialisationRequest([]byte{100, 100, 100})
	assert.NotEmpty(t, err)
}
