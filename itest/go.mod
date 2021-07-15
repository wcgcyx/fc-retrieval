module github.com/wcgcyx/fc-retrieval/itest

go 1.16

replace github.com/wcgcyx/fc-retrieval/common => ../common

replace github.com/wcgcyx/fc-retrieval/client => ../client

replace github.com/wcgcyx/fc-retrieval/gateway-admin => ../gateway-admin

replace github.com/wcgcyx/fc-retrieval/provider-admin => ../provider-admin

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-crypto v0.0.0-20191218222705-effae4ea9f03
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/lotus v1.10.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/ipfs/go-cid v0.0.7
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/wcgcyx/fc-retrieval/client v0.0.0-00010101000000-000000000000
	github.com/wcgcyx/fc-retrieval/common v0.0.0-20210713022117-13ad25a4fd06 // indirect
	github.com/wcgcyx/fc-retrieval/gateway-admin v0.0.0-00010101000000-000000000000
	github.com/wcgcyx/fc-retrieval/provider-admin v0.0.0-00010101000000-000000000000
)
