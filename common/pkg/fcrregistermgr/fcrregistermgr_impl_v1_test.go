/*
Package fcrregistermgr - register manager handles the interaction with the register.
*/
package fcrregistermgr

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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
)

func TestRegisterGateway(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/gateway", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		target := register.GatewayRegisteredInfo{}
		err := json.NewDecoder(r.Body).Decode(&target)
		assert.Empty(t, err)
		assert.Equal(t, "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16", target.RootKey)
		assert.Equal(t, "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", target.NodeID)
		assert.Equal(t, "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16", target.MsgSigningKey)
		assert.Equal(t, byte(1), target.MsgSigningKeyVer)
		assert.Equal(t, "au", target.RegionCode)
		assert.Equal(t, "testaddr", target.NetworkAddr)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	err := mgr.RegisterGateway("256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", &register.GatewayRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		RegionCode:       "au",
		NetworkAddr:      "testaddr",
	})
	assert.Empty(t, err)
}

func TestRegisterProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/provider", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		target := register.ProviderRegisteredInfo{}
		err := json.NewDecoder(r.Body).Decode(&target)
		assert.Empty(t, err)
		assert.Equal(t, "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16", target.RootKey)
		assert.Equal(t, "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", target.NodeID)
		assert.Equal(t, "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16", target.MsgSigningKey)
		assert.Equal(t, byte(1), target.MsgSigningKeyVer)
		assert.Equal(t, "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16", target.OfferSigningKey)
		assert.Equal(t, "au", target.RegionCode)
		assert.Equal(t, "testaddr", target.NetworkAddr)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	err := mgr.RegisterProvider("256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", &register.ProviderRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		RegionCode:       "au",
		NetworkAddr:      "testaddr",
	})
	assert.Empty(t, err)
}

func TestGetAllRegisteredGateway(t *testing.T) {
	gwInfo0 := &register.GatewayRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}
	gwInfo1 := &register.GatewayRegisteredInfo{
		RootKey:          "04dacfa291dfcf4b04c0936a5b2ec4e253af7032be38366de3d3ac5406345a1817999316c4a8cd7c3a8496e54886e7f3a47da0a962adebe569a42360d193316082",
		NodeID:           "ac1f490e923852ffc1a99d11b60b5ef378ff16f3cc71a1ce1f6983f064696ac5",
		MsgSigningKey:    "04dacfa291dfcf4b04c0936a5b2ec4e253af7032be38366de3d3ac5406345a1817999316c4a8cd7c3a8496e54886e7f3a47da0a962adebe569a42360d193316082",
		MsgSigningKeyVer: 2,
		RegionCode:       "us",
		NetworkAddr:      "testaddr1",
	}
	gwInfo2 := &register.GatewayRegisteredInfo{
		RootKey:          "04934a397c4c61f6c86c9630f25b668f0081798e8cd399442f44029f92fef11c345132d461b71db7e2abe6234edd51e8036f7530bc3f5abb9c8f539cbb23387b97",
		NodeID:           "844488bb57b7f6ae52f1797b8fea7663d3eb27efae2bea87875a44c37dbe983e",
		MsgSigningKey:    "04934a397c4c61f6c86c9630f25b668f0081798e8cd399442f44029f92fef11c345132d461b71db7e2abe6234edd51e8036f7530bc3f5abb9c8f539cbb23387b97",
		MsgSigningKeyVer: 3,
		RegionCode:       "ca",
		NetworkAddr:      "testaddr2",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/gateway", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal([]register.GatewayRegisteredInfo{*gwInfo0, *gwInfo1, *gwInfo2})
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 1444, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	gws, err := mgr.GetAllRegisteredGateway(0, 0)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(gws))
	assert.Equal(t, *gwInfo0, gws[0])
	assert.Equal(t, *gwInfo1, gws[1])
	assert.Equal(t, *gwInfo2, gws[2])
}

func TestGetAllRegisteredProvider(t *testing.T) {
	pvdInfo0 := &register.ProviderRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}
	pvdInfo1 := &register.ProviderRegisteredInfo{
		RootKey:          "04dacfa291dfcf4b04c0936a5b2ec4e253af7032be38366de3d3ac5406345a1817999316c4a8cd7c3a8496e54886e7f3a47da0a962adebe569a42360d193316082",
		NodeID:           "ac1f490e923852ffc1a99d11b60b5ef378ff16f3cc71a1ce1f6983f064696ac5",
		MsgSigningKey:    "04dacfa291dfcf4b04c0936a5b2ec4e253af7032be38366de3d3ac5406345a1817999316c4a8cd7c3a8496e54886e7f3a47da0a962adebe569a42360d193316082",
		MsgSigningKeyVer: 2,
		OfferSigningKey:  "04dacfa291dfcf4b04c0936a5b2ec4e253af7032be38366de3d3ac5406345a1817999316c4a8cd7c3a8496e54886e7f3a47da0a962adebe569a42360d193316082",
		RegionCode:       "us",
		NetworkAddr:      "testaddr1",
	}
	pvdInfo2 := &register.ProviderRegisteredInfo{
		RootKey:          "04934a397c4c61f6c86c9630f25b668f0081798e8cd399442f44029f92fef11c345132d461b71db7e2abe6234edd51e8036f7530bc3f5abb9c8f539cbb23387b97",
		NodeID:           "844488bb57b7f6ae52f1797b8fea7663d3eb27efae2bea87875a44c37dbe983e",
		MsgSigningKey:    "04934a397c4c61f6c86c9630f25b668f0081798e8cd399442f44029f92fef11c345132d461b71db7e2abe6234edd51e8036f7530bc3f5abb9c8f539cbb23387b97",
		MsgSigningKeyVer: 3,
		OfferSigningKey:  "04934a397c4c61f6c86c9630f25b668f0081798e8cd399442f44029f92fef11c345132d461b71db7e2abe6234edd51e8036f7530bc3f5abb9c8f539cbb23387b97",
		RegionCode:       "ca",
		NetworkAddr:      "testaddr2",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/provider", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal([]register.ProviderRegisteredInfo{*pvdInfo0, *pvdInfo1, *pvdInfo2})
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 1897, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	pvds, err := mgr.GetAllRegisteredProvider(0, 0)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(pvds))
	assert.Equal(t, *pvdInfo0, pvds[0])
	assert.Equal(t, *pvdInfo1, pvds[1])
	assert.Equal(t, *pvdInfo2, pvds[2])
}

func TestGetAllRegisteredGatewayByID(t *testing.T) {
	gwInfo := &register.GatewayRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/gateway/256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal(gwInfo)
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 480, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	gw, err := mgr.GetRegisteredGatewayByID("256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11")
	assert.Empty(t, err)
	assert.Equal(t, gwInfo, gw)
}

func TestGetAllRegisteredProviderByID(t *testing.T) {
	pvdInfo := &register.ProviderRegisteredInfo{
		RootKey:          "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		NodeID:           "256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11",
		MsgSigningKey:    "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "0496a1a3c388b63a577d7c6661cf615ede5c1d7c5545dbd9f3745ef81e8dbfbdb63139b78ecdc2d44982782f00bdaa1f77463a052debe93c86947f81d59bff2d16",
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/provider/256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal(pvdInfo)
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 631, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	pvd, err := mgr.GetRegisteredProviderByID("256a237ce1f8abac72728ac8f2edbe4a436ff1f898cd2e8ff869899e9bd92d11")
	assert.Empty(t, err)
	assert.Equal(t, pvdInfo, pvd)
}

func TestUnimplemented(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})
	height, err := mgr.GetHeight()
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), height)
	maxPage, err := mgr.GetGWMaxPage(height)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), maxPage)
	maxPage, err = mgr.GetPVDMaxPage(height)
	assert.Empty(t, err)
	assert.Equal(t, uint64(0), maxPage)
	assert.NotEmpty(t, mgr.UpdateGateway("test", nil))
	assert.NotEmpty(t, mgr.RequestDeregisterGateway("test"))
	assert.NotEmpty(t, mgr.DeregisterGateway("test"))
	assert.NotEmpty(t, mgr.UpdateProvider("test", nil))
	assert.NotEmpty(t, mgr.RequestDeregisterProvider("test"))
	assert.NotEmpty(t, mgr.DeregisterProvider("test"))
	_, err = mgr.GetRegisteredGatewaysByRegion(0, "test", 0)
	assert.NotEmpty(t, err)
	_, err = mgr.GetRegisteredProvidersByRegion(0, "test", 0)
	assert.NotEmpty(t, err)
}
