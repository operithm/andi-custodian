#!/bin/bash

# 1. Clean out broken module cache
go clean -modcache

# 2. Initialize clean module
rm go.mod go.sum
go mod init andi-custodian

# 3. Add correct dependencies
go get github.com/btcsuite/btcd@v0.23.3
go get github.com/ethereum/go-ethereum@v1.13.15
go get github.com/tyler-smith/go-bip39@v1.0.0

# 4. Verify
go mod tidy