module github.com/wcgcyx/fc-retrieval/itest

go 1.16

replace github.com/ConsenSys/fc-retrieval/common => ../common

replace github.com/ConsenSys/fc-retrieval/client => ../client

replace github.com/ConsenSys/fc-retrieval/gateway-admin => ../gateway-admin

replace github.com/ConsenSys/fc-retrieval/provider-admin => ../provider-admin

require (
	github.com/docker/docker v20.10.7+incompatible
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-crypto v0.0.0-20191218222705-effae4ea9f03
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/lotus v1.10.1
	github.com/google/uuid v1.2.0
	github.com/ipfs/go-cid v0.0.7
	github.com/testcontainers/testcontainers-go v0.11.1
	github.com/wcgcyx/fc-retrieval/common v0.0.0-20210713022117-13ad25a4fd06
	github.com/wcgcyx/testcontainers-go v0.10.1-0.20210511154849-504eecefabe0
)
