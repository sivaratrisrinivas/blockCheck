# BlockCheck

A lightweight Golang RESTful API service for validating blockchain wallet addresses, with a focus on Ethereum.

## Features

✅ Ethereum address format validation
✅ EIP-55 checksum validation and conversion
✅ ENS resolution support
✅ In-memory caching for improved performance
❌ Contract detection (coming soon)
❌ Rate limiting and API key support
❌ Dark mode UI (coming soon)

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `go mod download` to install dependencies
4. Start the server with `go run cmd/server/main.go`

## API Endpoints

### Validate Ethereum Address
```
GET /v1/validate/{address}
```

Validates an Ethereum address format and checksum.

**Response:**
```json
{
    "address": "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
    "isValid": true,
    "hasValidChecksum": true,
    "checksumAddress": "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
}
```

### Resolve ENS Name
```
GET /v1/resolveEns/{name}
```

Resolves an ENS name to an Ethereum address.

**Response:**
```json
{
    "name": "vitalik.eth",
    "address": "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
}
```

### Health Check
```
GET /health
```
Returns "OK" if the service is running.

### Coming Soon
- `GET /v1/addressType/{address}` - Check if an address is a contract or EOA

## Configuration

See `.env.example` for available configuration options:
- Server settings (host, port)
- ENS provider configuration
- Rate limiting options
- Cache settings

## Development

### Dependencies
- Go 1.22+
- Required packages:
  - `github.com/go-chi/chi/v5`: HTTP routing
  - `github.com/joho/godotenv`: Environment configuration
  - `github.com/sirupsen/logrus`: Structured logging
  - `github.com/ethereum/go-ethereum`: Ethereum client and utilities
  - `golang.org/x/crypto`: Cryptographic functions

### Testing
Coming soon

## License

MIT