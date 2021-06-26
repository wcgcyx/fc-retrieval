package fcrregistermgr

import "github.com/wcgcyx/fc-retrieval/common/pkg/register"

type FCRRegisterMgr interface {
	GetHeight() (uint64, error)

	RegisterGateway(id string, gwInfo *register.GatewayRegisteredInfo) error

	UpdateGateway(id string, gwInfo *register.GatewayRegisteredInfo) error

	RequestDeregisterGateway(id string) error

	DeregisterGateway(id string) error

	RegisterProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error

	UpdateProvider(id string, pvdInfo *register.ProviderRegisteredInfo) error

	RequestDeregisterProvider(id string) error

	DeregisterProvider(id string) error

	GetAllRegisteredGateway(height uint64, page uint64) ([]register.GatewayRegisteredInfo, error)

	GetAllRegisteredProvider(height uint64, page uint64) ([]register.ProviderRegisteredInfo, error)

	GetRegisteredGatewayByID(id string) (*register.GatewayRegisteredInfo, error)

	GetRegisteredProviderByID(id string) (*register.ProviderRegisteredInfo, error)

	GetRegisteredGatewaysByRegion(height uint64, region string, page uint64) ([]register.GatewayRegisteredInfo, error)

	GetRegisteredProvidersByRegion(height uint64, region string, page uint64) ([]register.ProviderRegisteredInfo, error)
}
