# Andi-Custodian: Multi-Chain MPC-Simulated Custody Service

A Multi-chain, MPC-simulated digital asset custodian built in Go—designed to reflect the engineering principles of institutional crypto custody platforms like **Anchorage Digital**.

Inspired by institutional custody challenges at firms like Anchorage Digital—where security, determinism, and multi-chain support are non-negotiable.

Built to demonstrate deep understanding of:
- Multi-chain transaction lifecycle (Bitcoin UTXO + Ethereum account model)
- Secure key handling & BIP-39 recovery
- Idempotent, nonce-aware, replay-safe transfers
- Extensible `Signer` interface (MPC-ready)
- Deterministic testing & auditability

> ⚠️ **For educational use only** — uses testnet keys. **Not for production or mainnet use.**

## 🔧 Features

- ✅ Generate BIP-39 mnemonic & HD wallet
- ✅ Derive Bitcoin (Testnet) & Ethereum (Sepolia) addresses
- ✅ Simulate UTXO selection (greedy algorithm)
- ✅ Fetch/assign Ethereum nonce safely
- ✅ Abstract signing via `Signer` interface (MPC-pluggable)
- ✅ Idempotency key support (via in-memory store)

## 🚀 Quick Start

1. Get a **Sepolia RPC URL** from [Alchemy](https://www.alchemy.com/) or [Infura](https)
2. Set environment variable:
   ```bash
   export SEPOLIA_RPC_URL="https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY"
3. Run demo: go run cmd/demo.main.go
4. Run Docker: 
   docker run --rm \
   -e SEPOLIA_RPC_URL="https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY" \
   andi-custodian