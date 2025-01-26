# Ethereum Address Tools

A simple web app that helps you work with Ethereum addresses - validate them, look up ENS names, and check if an address is a smart contract.

## Why?

Working with Ethereum addresses can be tricky:
- It's hard to tell if an address is valid or has typos
- ENS names like 'vitalik.eth' need to be converted to actual addresses
- You might need to know if you're dealing with a smart contract or a regular address

This tool makes these common tasks easy and quick, with a simple web interface.

## Quick Start

1. Make sure you have:
   - Go 1.21+
   - An Ethereum node access (like Infura)

2. Create a `config.yaml`:
```yaml
server:
  host: localhost
  port: 8080

ethereum:
  provider_url: your-ethereum-node-url
  cache_duration: 3600  # 1 hour

jwt:
  secret: your-secret-key
  expiry: 3600  # 1 hour
```

3. Run it:
```bash
go run cmd/server/main.go
```

4. Open `http://localhost:8080` in your browser

## Usage

### 1. Get an API Token
First, click "Generate Token" in the UI or use:
```bash
curl -X POST http://localhost:8080/v1/token
```

### 2. Check an Ethereum Address
```bash
# Example of a valid address
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/v1/validate/0x742d35Cc6634C0532925a3b844Bc454e4438f44e
```

### 3. Look up an ENS Name
```bash
# Convert vitalik.eth to its address
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/v1/resolveEns/vitalik.eth
```

### 4. Check if Address is a Contract
```bash
# Check if an address is a smart contract
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/v1/isContract/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
```

### Features
- ✅ Validates addresses using the EIP-55 standard
- 🔍 Converts ENS names to addresses
- 🤖 Detects smart contracts
- 🔒 Secure API with tokens
- 🌓 Dark/Light mode
- 📱 Works on mobile

## Contributing

1. Fork the repo
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

Found a bug or have a suggestion? Open an issue!