#!/bin/bash

go mod init andi-custodian
go get github.com/btcsuite/btcd@v0.23.3
go get github.com/ethereum/go-ethereum@v1.13.15
go get github.com/tyler-smith/go-bip39@v1.0.0