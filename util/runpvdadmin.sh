#!/bin/sh

var=$(go run main.go pvd)
docker run -it --network=shared --env DEVINIT=$var wcgcyx/fc-retrieval/provider-admin ./main