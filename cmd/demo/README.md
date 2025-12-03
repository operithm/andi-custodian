# andi-custodian

A Multi-chain, MPC-simulated custody service built in Go — designed to reflect the engineering principles of institutional-grade digital asset custody at firms like **Anchorage Digital**.

> “Ship code that will impact the global economy.”  
> — Anchorage Digital, Member of Technical Staff Role

This project demonstrates:
- ✅ BIP-39 mnemonic recovery
- ✅ HD wallet derivation (Bitcoin Testnet + Ethereum Sepolia)
- ✅ UTXO selection simulation
- ✅ Live Ethereum nonce fetching
- ✅ Multi-chain abstraction from a single seed

> ⚠️ **For educational use only** — uses testnet keys. **Not for production or mainnet use.**

## 🔧 Setup

1. **Get a Sepolia RPC URL** from [Alchemy](https://www.alchemy.com/) or [Infura](https://infura.io/)
2. **Set environment variable**:
   ```bash
   export SEPOLIA_RPC_URL="https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"


## Run log:

=== andi-custodian: Multi-Chain Custody Simulation ===

1. Generating BIP-39 mnemonic (12 words)...
   Mnemonic: inquiry ready express sudden always hammer brave acquire leaf neglect never cash

2. Deriving Bitcoin Testnet (Bech32) address...
   Bitcoin Address (Testnet): tb1qlgkq0jz2l37g8alaxw8923d9wfgv0tfc04ux5p

3. Deriving Ethereum (Sepolia) address...
   Ethereum Address: 0x567D41E38336F3681B28a6C5f1De8Bc4762Ef776

4. Simulating UTXO selection (target: 1.15 BTC)...
   → tx_a (6.000000 BTC)
   Total selected: 6.000000 BTC | Change: 4.850000 BTC

5. Fetching Ethereum nonce on Sepolia...
   Next nonce for 0x567D41E38336F3681B28a6C5f1De8Bc4762Ef776: 0
