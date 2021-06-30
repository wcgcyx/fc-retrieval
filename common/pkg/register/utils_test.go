/*
Package register - location for smart contract registration structs.
*/
package register

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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidGW(t *testing.T) {
	gwInfo := &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res := ValidateGatewayInfo(gwInfo)
	assert.True(t, res)
}

func TestInValidGW(t *testing.T) {
	gwInfo := &GatewayRegisteredInfo{
		RootKey:             "p04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res := ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "p59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "0159e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "69e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "p04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)

	gwInfo = &GatewayRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		RegionCode:          "",
		NetworkAddr:         "",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateGatewayInfo(gwInfo)
	assert.False(t, res)
}

func TestValidPVD(t *testing.T) {
	pvdInfo := &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res := ValidateProviderInfo(pvdInfo)
	assert.True(t, res)
}

func TestInValidPVD(t *testing.T) {
	pvdInfo := &ProviderRegisteredInfo{
		RootKey:             "p04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res := ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "p59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "2059e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "69e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "p04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "6c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "p04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "6c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "us",
		NetworkAddr:         "addr0",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)

	pvdInfo = &ProviderRegisteredInfo{
		RootKey:             "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		NodeID:              "59e548312e1cc4eeb25dc145ea458996441ad2898b5bf42487174456b80415fe",
		MsgSigningKey:       "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		MsgSigningKeyVer:    0,
		OfferSigningKey:     "04a66c41de8ad19f109fc4fc504d21ac376ddb32b8f3fcf60354a7a29e97bcb3d96146f992a60e53a511ec44a3bbbf719d524d863233452a7e9238efb271efe62d",
		RegionCode:          "",
		NetworkAddr:         "",
		Deregistering:       false,
		DeregisteringHeight: 0,
	}
	res = ValidateProviderInfo(pvdInfo)
	assert.False(t, res)
}
