package fcrregistermgr

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
		assert.Equal(t, "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", target.NodeID)
		assert.Equal(t, "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1", target.MsgSigningKey)
		assert.Equal(t, byte(1), target.MsgSigningKeyVer)
		assert.Equal(t, "au", target.RegionCode)
		assert.Equal(t, "testaddr", target.NetworkAddr)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	err := mgr.RegisterGateway("00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", &register.GatewayRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
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
		assert.Equal(t, "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", target.NodeID)
		assert.Equal(t, "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1", target.MsgSigningKey)
		assert.Equal(t, byte(1), target.MsgSigningKeyVer)
		assert.Equal(t, "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef2", target.OfferSigningKey)
		assert.Equal(t, "au", target.RegionCode)
		assert.Equal(t, "testaddr", target.NetworkAddr)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	err := mgr.RegisterProvider("00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", &register.ProviderRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef2",
		RegionCode:       "au",
		NetworkAddr:      "testaddr",
	})
	assert.Empty(t, err)
}

func TestGetAllRegisteredGateway(t *testing.T) {
	gwInfo0 := &register.GatewayRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
		MsgSigningKeyVer: 1,
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}
	gwInfo1 := &register.GatewayRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef2",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef3",
		MsgSigningKeyVer: 2,
		RegionCode:       "us",
		NetworkAddr:      "testaddr1",
	}
	gwInfo2 := &register.GatewayRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef4",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef5",
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
		assert.Equal(t, 817, len)
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
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeefa",
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}
	pvdInfo1 := &register.ProviderRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef2",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef3",
		MsgSigningKeyVer: 2,
		OfferSigningKey:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeefb",
		RegionCode:       "us",
		NetworkAddr:      "testaddr1",
	}
	pvdInfo2 := &register.ProviderRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef4",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef5",
		MsgSigningKeyVer: 3,
		OfferSigningKey:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeefc",
		RegionCode:       "ca",
		NetworkAddr:      "testaddr2",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/provider", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal([]register.ProviderRegisteredInfo{*pvdInfo0, *pvdInfo1, *pvdInfo2})
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 1072, len)
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
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
		MsgSigningKeyVer: 1,
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/gateway/00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal(gwInfo)
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 271, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	gw, err := mgr.GetRegisteredGatewayByID("00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0")
	assert.Empty(t, err)
	assert.Equal(t, gwInfo, gw)
}

func TestGetAllRegisteredProviderByID(t *testing.T) {
	pvdInfo := &register.ProviderRegisteredInfo{
		NodeID:           "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0",
		MsgSigningKey:    "00112233445566778899aabbccddeeff00112233445566778899aabbccddeef1",
		MsgSigningKeyVer: 1,
		OfferSigningKey:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeefa",
		RegionCode:       "au",
		NetworkAddr:      "testaddr0",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registers/provider/00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		data, err := json.Marshal(pvdInfo)
		assert.Empty(t, err)
		len, err := w.Write(data)
		assert.Equal(t, 356, len)
		assert.Empty(t, err)
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})

	pvd, err := mgr.GetRegisteredProviderByID("00112233445566778899aabbccddeeff00112233445566778899aabbccddeef0")
	assert.Empty(t, err)
	assert.Equal(t, pvdInfo, pvd)
}

func TestUnimplemented(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()
	// Initialise a manager
	mgr := NewFCRRegisterMgrImplV1(ts.URL, &http.Client{Timeout: 180 * time.Second})
	_, err := mgr.GetHeight()
	assert.NotEmpty(t, err)
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
