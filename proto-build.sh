#!/bin/bash

export PATH="$PATH:$(go env GOPATH)/bin"

protoc \
  --go_out=. \
  --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  api/custody/v1/custody.proto

go build -o custody-server ./cmd/server