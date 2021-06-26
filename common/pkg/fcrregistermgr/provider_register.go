package fcrregistermgr

// ProviderRegisteredInfo stores information of a registered provider
type ProviderRegisteredInfo interface {
	GetNodeID() string

	GetMsgSigningKey() string

	GetMsgSigningKeyVer() byte

	GetOfferSigningKey() string

	GetRegionCode() string

	GetNetworkAddr() string
}
