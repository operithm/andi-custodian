ando-custodian/
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docs/
│   └── architecture.md
├── cmd/
│   ├── andi-custodian/
│   │   └── main.go
│   └── demo/
│       └── main.go                 # Your utility + expanded entrypoint
├── internal/
│   ├── wallet/
│   │   ├── keygen.go               # Key/Mnemonic generation
│   │   ├── address.go              # BTC/ETH address derivation
│   │   ├── threshold.go            # Threshold key and policy
│   │   ├── verifier.go             # BTC/ETH Signature verification
│   │   ├── signer.go               # Signer interface + simulated MPC
│   │   └── wallet_demo.go          # Simple demo 
│   ├── chain/
│   │   ├── bitcoin.go              # UTXO selection, tx building
│   │   ├── ethereum.go             # Nonce, tx building
│   │   └── types.go                # Chain enum, tx types
│   ├── custody/
│   │   ├── service.go              # Core custody logic (transfer, monitor)
│   │   ├── nonce_manager.go        # Safe nonce assignment
│   │   └── utxo_selector.go        # Coin selection algorithms
│   └── store/
│       ├── store.go                # Interface (PostgreSQL impl)
│       └── inmemory.go             # In-memory store for demo
├── pkg/
│   └── types/                      # Public types (if you expose SDK)
├── test/
│   └── integration/                # Sepolia + Bitcoin testnet tests
├── deployments/                    # Kubernetes/Fly.io config
└── LICENSE