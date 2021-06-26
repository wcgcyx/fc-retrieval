package fcrregistermgr

import (
	"net/http"
	"sync"

	"github.com/wcgcyx/fc-retrieval/common/pkg/register"
)

type FCRRegisterMgrImplV1 struct {
	client *http.Client
	lock   sync.RWMutex
}

func NewFCRRegisterMgrVImplV1(client *http.Client) FCRRegisterMgr {
	return &FCRRegisterMgrImplV1{client: client, lock: sync.RWMutex{}}
}

func (mgr *FCRRegisterMgrImplV1) GetHeight() (uint64, error) {
	return 0, nil
}

func (mgr *FCRRegisterMgrImplV1) RegisterGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) UpdateGateway(id string, gwInfo *register.GatewayRegisteredInfo) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) DeregisterGateway(id string) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) RegisterProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) UpdateProvider(id string, gwInfo *register.ProviderRegisteredInfo) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) DeregisterProvider(id string) error {
	return nil
}

func (mgr *FCRRegisterMgrImplV1) GetAllRegisteredGateway(height uint64, page uint64) ([]register.GatewayRegisteredInfo, error) {
	return nil, nil
}

func (mgr *FCRRegisterMgrImplV1) GetAllRegisteredProvider(height uint64, page uint64) ([]register.ProviderRegisteredInfo, error) {
	return nil, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredGatewayByID(id string) (*register.GatewayRegisteredInfo, error) {
	return nil, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredProviderByID(id string) (*register.ProviderRegisteredInfo, error) {
	return nil, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error) {
	return nil, nil
}

func (mgr *FCRRegisterMgrImplV1) GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error) {
	return nil, nil
}
