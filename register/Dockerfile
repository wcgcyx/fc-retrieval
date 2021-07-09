# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.16-alpine as builder
RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/wcgcyx/fc-retrieval
COPY common common
COPY register register
WORKDIR /go/src/github.com/wcgcyx/fc-retrieval/register

# Get all dependancies, but don't install.
RUN go get -d -v github.com/wcgcyx/fc-retrieval/register/cmd/register-server
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/register-server ./cmd/register-server

# Pull build register into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/register-server /register-server
COPY register/.env.example /.env

# Run the binary when the container starts.
WORKDIR /
CMD ["/register-server", "--host", "0.0.0.0", "--port", "9020"]
EXPOSE 9020
