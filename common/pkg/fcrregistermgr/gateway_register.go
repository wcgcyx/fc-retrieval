package fcrregistermgr

// GatewayRegisteredInfo stores information of a registered gateway
type GatewayRegisteredInfo interface {
	GetNodeID() string

	GetMsgSigningKey() string

	GetMsgSigningKeyVer() byte

	GetRegionCode() string

	GetNetworkAddr() string
}
