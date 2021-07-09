# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.16-alpine as builder

RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/wcgcyx/fc-retrieval
COPY common common
COPY gateway-admin gateway-admin
WORKDIR /go/src/github.com/wcgcyx/fc-retrieval/gateway-admin

# Get all dependancies, but don't install.
RUN go get -d -v github.com/wcgcyx/fc-retrieval/gateway-admin/cmd/gatewayadmin-cli
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/gatewayadmin ./cmd/gatewayadmin-cli

# Pull build gateway into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/gatewayadmin /main
WORKDIR /