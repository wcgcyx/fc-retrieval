/*
Package fcrpeermgr - peer manager manages all retrieval peers.
*/
package fcrpeermgr

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
)

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
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000000",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000010",
			MsgSigningKeyVer:    2,
			RegionCode:          "au",
			NetworkAddr:         "addr0",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000001",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000011",
			MsgSigningKeyVer:    3,
			RegionCode:          "au",
			NetworkAddr:         "addr1",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000002",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000012",
			MsgSigningKeyVer:    4,
			RegionCode:          "us",
			NetworkAddr:         "addr2",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000003",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000013",
			MsgSigningKeyVer:    5,
			RegionCode:          "ca",
			NetworkAddr:         "addr3",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[1] = []register.GatewayRegisteredInfo{
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000004",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000014",
			MsgSigningKeyVer:    6,
			RegionCode:          "ca",
			NetworkAddr:         "addr4",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000005",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000015",
			MsgSigningKeyVer:    7,
			RegionCode:          "us",
			NetworkAddr:         "addr5",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000006",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000016",
			MsgSigningKeyVer:    8,
			RegionCode:          "cn",
			NetworkAddr:         "addr6",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000007",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000017",
			MsgSigningKeyVer:    9,
			RegionCode:          "cn",
			NetworkAddr:         "addr7",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[2] = []register.GatewayRegisteredInfo{
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000008",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000018",
			MsgSigningKeyVer:    10,
			RegionCode:          "ca",
			NetworkAddr:         "addr8",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000009",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000019",
			MsgSigningKeyVer:    11,
			RegionCode:          "us",
			NetworkAddr:         "addr9",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000A",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001A",
			MsgSigningKeyVer:    12,
			RegionCode:          "cn",
			NetworkAddr:         "addr10",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000B",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001B",
			MsgSigningKeyVer:    13,
			RegionCode:          "cn",
			NetworkAddr:         "addr11",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[3] = []register.GatewayRegisteredInfo{
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000C",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001C",
			MsgSigningKeyVer:    14,
			RegionCode:          "ca",
			NetworkAddr:         "addr12",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000D",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001D",
			MsgSigningKeyVer:    15,
			RegionCode:          "us",
			NetworkAddr:         "addr13",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000E",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001E",
			MsgSigningKeyVer:    16,
			RegionCode:          "cn",
			NetworkAddr:         "addr14",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "000000000000000000000000000000000000000000000000000000000000000F",
			MsgSigningKey:       "000000000000000000000000000000000000000000000000000000000000001F",
			MsgSigningKeyVer:    17,
			RegionCode:          "cn",
			NetworkAddr:         "addr15",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.gws[4] = []register.GatewayRegisteredInfo{
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000010",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000020",
			MsgSigningKeyVer:    18,
			RegionCode:          "ca",
			NetworkAddr:         "addr16",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000011",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000021",
			MsgSigningKeyVer:    19,
			RegionCode:          "us",
			NetworkAddr:         "addr17",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000012",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000022",
			MsgSigningKeyVer:    20,
			RegionCode:          "cn",
			NetworkAddr:         "addr18",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000013",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000023",
			MsgSigningKeyVer:    21,
			RegionCode:          "cn",
			NetworkAddr:         "addr19",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
	}
	res.pvds[0] = []register.ProviderRegisteredInfo{
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000014",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000024",
			MsgSigningKeyVer:    22,
			OfferSigningKey:     "0000000000000000000000000000000000000000000000000000000000000034",
			RegionCode:          "ru",
			NetworkAddr:         "addr20",
			Deregistering:       false,
			DeregisteringHeight: 0,
		},
		{
			NodeID:              "0000000000000000000000000000000000000000000000000000000000000015",
			MsgSigningKey:       "0000000000000000000000000000000000000000000000000000000000000025",
			MsgSigningKeyVer:    23,
			OfferSigningKey:     "0000000000000000000000000000000000000000000000000000000000000035",
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

func (m *mockRegisterMgr) GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error) {
	return nil, nil
}

func (m *mockRegisterMgr) GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error) {
	return nil, nil
}

func TestAutoSync(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	// No effect before starting manager routine
	peerMgr.Sync()
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000000")
	peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	peerMgr.Shutdown()
	_, err := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NotEmpty(t, err)
	_, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.NotEmpty(t, err)
	err = peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	err = peerMgr.Start()
	assert.NotEmpty(t, err)
	time.Sleep(2 * time.Second)
	min, max, err := peerMgr.GetCurrentCIDHashRange()
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)
	gws, err := peerMgr.GetGWSNearCIDHash("0000000000000000000000000000000000000000000000000000000000000008", "0000000000000000000000000000000000000000000000000000000000000009")
	assert.Empty(t, err)
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
	assert.Equal(t, gws[9].NodeID, "000000000000000000000000000000000000000000000000000000000000000A")
	assert.Equal(t, gws[10].NodeID, "000000000000000000000000000000000000000000000000000000000000000B")
	assert.Equal(t, gws[11].NodeID, "000000000000000000000000000000000000000000000000000000000000000C")
	assert.Equal(t, gws[12].NodeID, "000000000000000000000000000000000000000000000000000000000000000D")
	assert.Equal(t, gws[13].NodeID, "000000000000000000000000000000000000000000000000000000000000000E")
	assert.Equal(t, gws[14].NodeID, "000000000000000000000000000000000000000000000000000000000000000F")
	assert.Equal(t, gws[15].NodeID, "0000000000000000000000000000000000000000000000000000000000000010")
}

func TestSyncUpgrade(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	err := peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	peerMgr.Sync()
	// Test update gw entry
	min, max, err := peerMgr.GetCurrentCIDHashRange()
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)
	gws, err := peerMgr.GetGWSNearCIDHash("0000000000000000000000000000000000000000000000000000000000000008", "0000000000000000000000000000000000000000000000000000000000000009")
	assert.Empty(t, err)
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
	assert.Equal(t, gws[9].NodeID, "000000000000000000000000000000000000000000000000000000000000000A")
	assert.Equal(t, gws[10].NodeID, "000000000000000000000000000000000000000000000000000000000000000B")
	assert.Equal(t, gws[11].NodeID, "000000000000000000000000000000000000000000000000000000000000000C")
	assert.Equal(t, gws[12].NodeID, "000000000000000000000000000000000000000000000000000000000000000D")
	assert.Equal(t, gws[13].NodeID, "000000000000000000000000000000000000000000000000000000000000000E")
	assert.Equal(t, gws[14].NodeID, "000000000000000000000000000000000000000000000000000000000000000F")
	assert.Equal(t, gws[15].NodeID, "0000000000000000000000000000000000000000000000000000000000000010")
	peer, err := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, err)
	assert.Equal(t, byte(4), peer.MsgSigningKeyVer)
	mockRegisterMgr.gws[0][2].MsgSigningKeyVer = 1
	peerMgr.Sync()
	peer, err = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, err)
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove gw entry
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0][:2], mockRegisterMgr.gws[0][3:]...)
	peerMgr.Sync()
	_, err = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.NotEmpty(t, err)
	min, max, err = peerMgr.GetCurrentCIDHashRange()
	assert.Empty(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", min)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000011", max)

	// Test update pvd entry
	peer, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, err)
	assert.Equal(t, byte(22), peer.MsgSigningKeyVer)
	mockRegisterMgr.pvds[0][0].MsgSigningKeyVer = 1
	peerMgr.Sync()
	peer, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, err)
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove pvd entry
	mockRegisterMgr.pvds[0] = mockRegisterMgr.pvds[0][1:]
	peerMgr.Sync()
	_, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.NotEmpty(t, err)
}

func TestSyncSingle(t *testing.T) {
	mockRegisterMgr := newMockRegister()
	peerMgr := NewFCRPeerMgrImplV1(mockRegisterMgr, true, true, true, "0000000000000000000000000000000000000000000000000000000000000009", time.Second)
	err := peerMgr.Start()
	assert.Empty(t, err)
	defer peerMgr.Shutdown()
	peerMgr.Sync()
	peer, err := peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, err)
	assert.Equal(t, byte(4), peer.MsgSigningKeyVer)
	mockRegisterMgr.gws[0][2].MsgSigningKeyVer = 1
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	peer, err = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, err)
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove gw entry
	tempgw := mockRegisterMgr.gws[0][2]
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0][:2], mockRegisterMgr.gws[0][3:]...)
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	_, err = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.NotEmpty(t, err)
	// Test New gw entry
	mockRegisterMgr.gws[0] = append(mockRegisterMgr.gws[0], tempgw)
	peerMgr.SyncGW("0000000000000000000000000000000000000000000000000000000000000002")
	_, err = peerMgr.GetGWInfo("0000000000000000000000000000000000000000000000000000000000000002")
	assert.Empty(t, err)

	// Test pvd
	peer, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, err)
	assert.Equal(t, byte(22), peer.MsgSigningKeyVer)
	mockRegisterMgr.pvds[0][0].MsgSigningKeyVer = 1
	peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	peer, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, err)
	assert.Equal(t, byte(1), peer.MsgSigningKeyVer)
	// Test remove pvd entry
	temppvd := mockRegisterMgr.pvds[0][0]
	mockRegisterMgr.pvds[0] = mockRegisterMgr.pvds[0][1:]
	peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	_, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.NotEmpty(t, err)
	// Test new pvd entry
	mockRegisterMgr.pvds[0] = append(mockRegisterMgr.pvds[0], temppvd)
	peerMgr.SyncPVD("0000000000000000000000000000000000000000000000000000000000000014")
	_, err = peerMgr.GetPVDInfo("0000000000000000000000000000000000000000000000000000000000000014")
	assert.Empty(t, err)
}
