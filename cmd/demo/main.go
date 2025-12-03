package main

import (
	"andi-custodian/internal/wallet"
	"fmt"
	"log"
	"os"
)

func main() {
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("SEPOLIA_RPC_URL environment variable is required (e.g., Alchemy or Infura endpoint)")
	}

	fmt.Println("=== andi-custodian: Multi-Chain Custody Simulation ===")

	if err := wallet.RunCustodyDemo(rpcURL); err != nil {
		log.Fatal(err)
	}
}
