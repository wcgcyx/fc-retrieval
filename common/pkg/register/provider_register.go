package register

// ProviderRegisteredInfo stores information of a registered provider
type ProviderRegisteredInfo struct {
	NodeID string

	MsgSigningKey string

	MsgSigningKeyVer byte

	OfferSigningKey string

	RegionCode string

	NetworkAddr string

	Deregistering bool
}
