# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.16-alpine as builder

RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/wcgcyx/fc-retrieval
COPY common common
COPY gateway gateway
WORKDIR /go/src/github.com/wcgcyx/fc-retrieval/gateway


# Get all dependancies, but don't install.
RUN go get -d -v github.com/wcgcyx/fc-retrieval/gateway/cmd/gateway
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/gateway ./cmd/gateway

# Pull build gateway into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/gateway /main
COPY gateway/docker-entrypoint.sh /docker-entrypoint.sh
COPY gateway/.env.example /.env

# Run the binary when the container starts.
WORKDIR /
CMD ["./docker-entrypoint.sh"]
EXPOSE 9010 9011
