#!/bin/sh

var=$(go run main.go client)
docker run -it --network=shared --env DEVINIT=$var wcgcyx/fc-retrieval/client ./main