package register

// GatewayRegisteredInfo stores information of a registered gateway
type GatewayRegisteredInfo struct {
	NodeID string

	MsgSigningKey string

	MsgSigningKeyVer byte

	RegionCode string

	NetworkAddr string

	Deregistering bool
}
