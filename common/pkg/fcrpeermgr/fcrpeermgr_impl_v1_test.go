/*
Package fcrpeermgr - peer manager manages all retrieval peers.
*/
package fcrpeermgr

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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/fcrreputationmgr"
	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
)

type mockRegisterMgr struct {
	gws  map[uint64]([]register.GatewayRegisteredInfo)
	pvds map[uint64]([]register.ProviderRegisteredInfo)
}

func newMockRegister() *mockRegisterMgr {
	res := &mockRegisterMgr{
		gws:  make(map[uint64][]register.GatewayRegisteredInfo),
		pvds: make(map[uint64][]register.ProviderRegisteredInfo),
	}
	res.gws[0] = []register.GatewayRegisteredInfo{
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000000",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010",
			MsgSigningKeyVer:    2,
			RegionCode:          "au",
			NetworkAddr:         "addr0",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000001",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011",
			MsgSigningKeyVer:    3,
			RegionCode:          "au",
			NetworkAddr:         "addr1",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000002",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012",
			MsgSigningKeyVer:    4,
			RegionCode:          "us",
			NetworkAddr:         "addr2",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000003",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013",
			MsgSigningKeyVer:    5,
			RegionCode:          "ca",
			NetworkAddr:         "addr3",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[1] = []register.GatewayRegisteredInfo{
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000004",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014",
			MsgSigningKeyVer:    6,
			RegionCode:          "ca",
			NetworkAddr:         "addr4",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000005",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000015",
			MsgSigningKeyVer:    7,
			RegionCode:          "us",
			NetworkAddr:         "addr5",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000006",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016",
			MsgSigningKeyVer:    8,
			RegionCode:          "cn",
			NetworkAddr:         "addr6",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000007",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000017",
			MsgSigningKeyVer:    9,
			RegionCode:          "cn",
			NetworkAddr:         "addr7",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[2] = []register.GatewayRegisteredInfo{
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000008",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018",
			MsgSigningKeyVer:    10,
			RegionCode:          "ca",
			NetworkAddr:         "addr8",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000009",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000019",
			MsgSigningKeyVer:    11,
			RegionCode:          "us",
			NetworkAddr:         "addr9",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000a",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a",
			MsgSigningKeyVer:    12,
			RegionCode:          "cn",
			NetworkAddr:         "addr10",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000b",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001b",
			MsgSigningKeyVer:    13,
			RegionCode:          "cn",
			NetworkAddr:         "addr11",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[3] = []register.GatewayRegisteredInfo{
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000c",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001c",
			MsgSigningKeyVer:    14,
			RegionCode:          "ca",
			NetworkAddr:         "addr12",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000d",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001d",
			MsgSigningKeyVer:    15,
			RegionCode:          "us",
			NetworkAddr:         "addr13",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000e",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e",
			MsgSigningKeyVer:    16,
			RegionCode:          "cn",
			NetworkAddr:         "addr14",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f",
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000f",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001f",
			MsgSigningKeyVer:    17,
			RegionCode:          "cn",
			NetworkAddr:         "addr15",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[4] = []register.GatewayRegisteredInfo{
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000010",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020",
			MsgSigningKeyVer:    18,
			RegionCode:          "ca",
			NetworkAddr:         "addr16",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000011",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021",
			MsgSigningKeyVer:    19,
			RegionCode:          "us",
			NetworkAddr:         "addr17",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000012",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000022",
			MsgSigningKeyVer:    20,
			RegionCode:          "cn",
			NetworkAddr:         "addr18",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000013",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000023",
			MsgSigningKeyVer:    21,
			RegionCode:          "cn",
			NetworkAddr:         "addr19",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.pvds[0] = []register.ProviderRegisteredInfo{
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000014",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024",
			MsgSigningKeyVer:    22,
			OfferSigningKey:     "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034",
			RegionCode:          "ru",
			NetworkAddr:         "addr20",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			RootKey:             "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000015",
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000015",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000025",
			MsgSigningKeyVer:    23,
			OfferSigningKey:     "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000035",
			RegionCode:          "fr",
			NetworkAddr:         "addr21",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	return res
}

func (m *mockRegisterMgr) GetHeight() (uint64, error) {
	return 0, nil
}

func (m *mockRegisterMgr) GetGWMaxPage(height uint64) (uint64, error) {
	return 4, nil
}

func (m *mockRegisterMgr) GetPVDMaxPage(height uint64) (uint64, error) {
	return 0, nil
}

func (m *mockRegisterMgr) RegisterGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	return nil
}

func (m *mockRegisterMgr) UpdateGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	return nil
}

func (m *mockRegisterMgr) RequestDeregisterGateway(id string) error {
	return nil
}

func (m *mockRegisterMgr) DeregisterGateway(id string) error {
	return nil
}

func (m *mockRegisterMgr) RegisterProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error {
	return nil
}

func (m *mockRegisterMgr) UpdateProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error {
	return nil
}

func (m *mockRegisterMgr) RequestDeregisterProvider(id string) error {
	return nil
}

func (m *mockRegisterMgr) DeregisterProvider(id string) error {
	return nil
}

func (m *mockRegisterMgr) GetAllRegisteredGateway(height uint64, page uint64) ([]register.GatewayRegisteredInfo, error) {
	res, ok := m.gws[page]
	if !ok {
		return nil, errors.New("Not found")
	}
	return res, nil
}

func (m *mockRegisterMgr) GetAllRegisteredProvider(height uint64, page uint64) ([]register.ProviderRegisteredInfo, error) {
	res, ok := m.pvds[page]
	if !ok {
		return nil, errors.New("Not found")
	}
	return res, nil
}

func (m *mockRegisterMgr) GetRegisteredGatewayByID(id string) (*register.GatewayRegisteredInfo, error) {
	for _, val := range m.gws {
		for _, info := range val {
			if info.NodeID == id {
				return &info, nil
			}
		}
	}
	return nil, errors.New("Not found")
}

func (m *mockRegisterMgr) GetRegisteredProviderByID(id string) (*register.ProviderRegisteredInfo, error) {
	for _, val := range m.pvds {
		for _, info := range val {
			if info.NodeID == id {
				return &info, nil
			}
		}
	}
	return nil, errors.New("Not found")
}

func (m *mockRegisterMgr) GetGWMaxPageByRegion(height uint64, region string) (uint64, error) {
	return 0, nil
}

func (m *mockRegisterMgr) GetPVDMaxPageByRegion(height uint64, region string) (uint64, error) {
	return 0, nil
}

func (m *mockRegisterMgr) GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error) {
	return nil, nil
}

func (m *mockRegisterMgr) GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error) {
	return nil, nil
}

func TestAutoSync(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	mockReputationMgr := fcrreputationmgr.NewFCRReputationMgrImpV1()
	err := mockReputationMgr.Start()
	assert.Empty(t, err)
	defer mockReputationMgr.Shutdown()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, mockReputationMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	// No effect before starting manager routine
	peerMgr.Sync()
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000000")
	peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	peerMgr.Shutdown()
	peer := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Empty(t, peer)
	peer = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, peer)
	err = peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	err = peerMgr.Start()
	assert.NotEmpty(t, err)
	time.Sleep(2 * time.Second)
	min, max := peerMgr.GetCurrentCIDHashRange()
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)
	gws := peerMgr.GetGWSNearCIDHash("0000000000000000000000000000000000000000000000000000000000000008", 16, "0000000000000000000000000000000000000000000000000000000000000009")
	assert.Equal(t, 16, len(gws))
	assert.Equal(t, gws[0].NodeID, "0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, gws[1].NodeID, "0000000000000000000000000000000000000000000000000000000000000001")
	assert.Equal(t, gws[2].NodeID, "0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, gws[3].NodeID, "0000000000000000000000000000000000000000000000000000000000000003")
	assert.Equal(t, gws[4].NodeID, "0000000000000000000000000000000000000000000000000000000000000004")
	assert.Equal(t, gws[5].NodeID, "0000000000000000000000000000000000000000000000000000000000000005")
	assert.Equal(t, gws[6].NodeID, "0000000000000000000000000000000000000000000000000000000000000006")
	assert.Equal(t, gws[7].NodeID, "0000000000000000000000000000000000000000000000000000000000000007")
	assert.Equal(t, gws[8].NodeID, "0000000000000000000000000000000000000000000000000000000000000008")
	assert.Equal(t, gws[9].NodeID, "000000000000000000000000000000000000000000000000000000000000000a")
	assert.Equal(t, gws[10].NodeID, "000000000000000000000000000000000000000000000000000000000000000b")
	assert.Equal(t, gws[11].NodeID, "000000000000000000000000000000000000000000000000000000000000000c")
	assert.Equal(t, gws[12].NodeID, "000000000000000000000000000000000000000000000000000000000000000d")
	assert.Equal(t, gws[13].NodeID, "000000000000000000000000000000000000000000000000000000000000000e")
	assert.Equal(t, gws[14].NodeID, "000000000000000000000000000000000000000000000000000000000000000f")
	assert.Equal(t, gws[15].NodeID, "0000000000000000000000000000000000000000000000000000000000000010")
}

func TestSyncUpgrade(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	mockReputationMgr := fcrreputationmgr.NewFCRReputationMgrImpV1()
	err := mockReputationMgr.Start()
	assert.Empty(t, err)
	defer mockReputationMgr.Shutdown()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, mockReputationMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	err = peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	peerMgr.Sync()
	// Test update gw entry
	min, max := peerMgr.GetCurrentCIDHashRange()
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)
	gws := peerMgr.GetGWSNearCIDHash("0000000000000000000000000000000000000000000000000000000000000008", 16, "0000000000000000000000000000000000000000000000000000000000000009")
	assert.Equal(t, 16, len(gws))
	assert.Equal(t, gws[0].NodeID, "0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, gws[1].NodeID, "0000000000000000000000000000000000000000000000000000000000000001")
	assert.Equal(t, gws[2].NodeID, "0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, gws[3].NodeID, "0000000000000000000000000000000000000000000000000000000000000003")
	assert.Equal(t, gws[4].NodeID, "0000000000000000000000000000000000000000000000000000000000000004")
	assert.Equal(t, gws[5].NodeID, "0000000000000000000000000000000000000000000000000000000000000005")
	assert.Equal(t, gws[6].NodeID, "0000000000000000000000000000000000000000000000000000000000000006")
	assert.Equal(t, gws[7].NodeID, "0000000000000000000000000000000000000000000000000000000000000007")
	assert.Equal(t, gws[8].NodeID, "0000000000000000000000000000000000000000000000000000000000000008")
	assert.Equal(t, gws[9].NodeID, "000000000000000000000000000000000000000000000000000000000000000a")
	assert.Equal(t, gws[10].NodeID, "000000000000000000000000000000000000000000000000000000000000000b")
	assert.Equal(t, gws[11].NodeID, "000000000000000000000000000000000000000000000000000000000000000c")
	assert.Equal(t, gws[12].NodeID, "000000000000000000000000000000000000000000000000000000000000000d")
	assert.Equal(t, gws[13].NodeID, "000000000000000000000000000000000000000000000000000000000000000e")
	assert.Equal(t, gws[14].NodeID, "000000000000000000000000000000000000000000000000000000000000000f")
	assert.Equal(t, gws[15].NodeID, "0000000000000000000000000000000000000000000000000000000000000010")
	peer := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, byte(4), peer.MsgSigningKeyVer)
	mockRegisterMgr.gws[0][2].MsgSigningKeyVer = 1
	peerMgr.Sync()
	peer = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove gw entry
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0][:2], mockRegisterMgr.gws[0][3:]...)
	peerMgr.Sync()
	peer = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, peer)
	min, max = peerMgr.GetCurrentCIDHashRange()
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)

	// Test update pvd entry
	peer = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Equal(t, byte(22), peer.MsgSigningKeyVer)
	mockRegisterMgr.pvds[0][0].MsgSigningKeyVer = 1
	peerMgr.Sync()
	peer = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove pvd entry
	mockRegisterMgr.pvds[0] = mockRegisterMgr.pvds[0][1:]
	peerMgr.Sync()
	peer = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, peer)
}

func TestSyncSingle(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	mockReputationMgr := fcrreputationmgr.NewFCRReputationMgrImpV1()
	err := mockReputationMgr.Start()
	assert.Empty(t, err)
	defer mockReputationMgr.Shutdown()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, mockReputationMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	err = peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	peerMgr.Sync()
	peer := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, byte(4), peer.MsgSigningKeyVer)
	mockRegisterMgr.gws[0][2].MsgSigningKeyVer = 1
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	peer = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove gw entry
	tempgw := mockRegisterMgr.gws[0][2]
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0][:2], mockRegisterMgr.gws[0][3:]...)
	peer = peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, peer)
	// Test New gw entry
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0], tempgw)
	peer = peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	assert.NotEmpty(t, peer)

	// Test pvd
	peer = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Equal(t, byte(22), peer.MsgSigningKeyVer)
	mockRegisterMgr.pvds[0][0].MsgSigningKeyVer = 1
	peer = peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove pvd entry
	temppvd := mockRegisterMgr.pvds[0][0]
	mockRegisterMgr.pvds[0] = mockRegisterMgr.pvds[0][1:]
	peer = peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, peer)
	// Test new pvd entry
	mockRegisterMgr.pvds[0] = append(mockRegisterMgr.pvds[0], temppvd)
	peer = peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	assert.NotEmpty(t, peer)

	// Test List GWS
	peers := peerMgr.ListGWS()
	assert.Equal(t, 20, len(peers))
}
