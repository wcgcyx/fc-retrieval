#!/bin/sh

var=$(go run main.go gw)
docker run -it --network=shared --env DEVINIT=$var wcgcyx/fc-retrieval/gateway-admin ./main