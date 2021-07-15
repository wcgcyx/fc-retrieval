# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.16-alpine as builder

RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/wcgcyx/fc-retrieval
COPY common common
COPY client client
COPY gateway-admin gateway-admin
COPY provider-admin provider-admin
COPY itest itest
WORKDIR /go/src/github.com/wcgcyx/fc-retrieval/itest

# Get all dependancies, but don't install.
RUN go get -d -v github.com/wcgcyx/fc-retrieval/itest/pkg/e2e
RUN go get -d -v github.com/wcgcyx/fc-retrieval/itest/pkg/util

# Run test
CMD ["go", "test", "-v", "-p=1", "--count=1", "./..."]