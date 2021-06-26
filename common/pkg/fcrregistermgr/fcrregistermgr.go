package fcrregistermgr

type FCRRegisterMgr interface {
	RegisterGateway(gwInfo GatewayRegisteredInfo) error

	RegisterProvider(pvdInfo ProviderRegisteredInfo) error

	// TODO: Paging
	GetAllRegisteredGateway() ([]GatewayRegisteredInfo, error)

	// TODO: Paging
	GetAllRegisteredProvider() ([]ProviderRegisteredInfo, error)

	GetRegisteredGatewayByID(id string) (GatewayRegisteredInfo, error)

	GetRegisteredProviderByID(id string) (ProviderRegisteredInfo, error)

	// TODO: Paging
	GetRegisteredGatewaysByRegion(region string) ([]GatewayRegisteredInfo, error)

	// TODO: Paging
	GetRegisteredProvidersByRegion(region string) ([]ProviderRegisteredInfo, error)
}
