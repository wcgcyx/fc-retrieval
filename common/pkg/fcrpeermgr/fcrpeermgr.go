package fcrpeermgr

type FCRPeerMgr interface {
	Start()

	Shutdown()

	Refresh()

	RefreshGW(gwID string)

	RefreshPVD(pvdID string)

	GetGWNetworkInfo(gwID string) (string, error)

	GetGWMessageSigningKey(gwID string) (string, byte, error)

	GetPVDNetworkInfo(pvdID string) (string, error)

	GetPVDMessageSigningKey(pvdID string) (string, byte, error)

	GetPVDOfferSigningKey(pvdID string) (string, error)
}
