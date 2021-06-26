package fcrmessages

type FCRMessage interface {
	GetMessageType() int32
	GetProtocolVersion() int32
	GetProtocolSupported() int32
	GetMessageBody() []byte
	GetSignature() string
	Sign(privKey string, keyVer byte) error
	Verify(pubKey string, keyVer byte) error

	ToBytes() ([]byte, error)
	FromBytes(data []byte) error
}
