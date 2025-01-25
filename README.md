# BlockCheck - Ethereum Address & ENS Validation Service

A high-performance Go service for validating Ethereum addresses and resolving ENS names.

## Features

✅ Ethereum Address Format Validation
- Validates address format and checksum
- EIP-55 compliant validation
- Fast regex-based validation

✅ ENS Resolution Support
- Resolves ENS names to Ethereum addresses
- Caches resolutions for improved performance
- Handles non-existent names gracefully

✅ Plugin Architecture
- Extensible validator system
- Support for multiple chains (Ethereum implemented)
- Easy to add new validators

✅ Caching System
- Redis support for distributed caching
- In-memory fallback cache
- Configurable TTL and cache strategies

✅ Contract Detection
- Check if an address is a contract or EOA
- Fast response times (~64ms for contracts)
- Proper error handling for invalid addresses

✅ Security & Authentication
- JWT-based authentication
- API key generation
- Protected endpoints
- Configurable token expiration

❌ Rate Limiting & Security
- Request rate limiting
- IP-based throttling
- API key authentication

## API Endpoints

### Health Check (Public)
```
GET /health
```
Returns `200 OK` if service is healthy

### Generate API Token (Public)
```
POST /v1/token
```
Generates a new API key and JWT token.

Example Response:
```json
{
  "api_key": "550e8400-e29b-41d4-a716-446655440000",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Protected Endpoints
All other endpoints require authentication using the JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

Example:
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." http://localhost:8080/v1/validate/0x...
```

### Validate Ethereum Address
```
GET /v1/validate/{address}
```
Example Response:
```json
{
  "address": "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
  "isValid": true
}
```

### Resolve ENS Name
```
GET /v1/resolveEns/{name}
```
Example Response:
```json
{
  "name": "vitalik.eth",
  "address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
}
```

### Check Contract Status
```
GET /v1/isContract/{address}
```
Checks if an Ethereum address is a contract (smart contract) or an EOA (Externally Owned Account).

Example Response (Contract):
```json
{
  "address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
  "isContract": true
}
```

Example Response (EOA):
```json
{
  "address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
  "isContract": false
}
```

Error Response:
```json
{
  "address": "0xinvalid",
  "error": "invalid address format"
}
```

Response Codes:
- `200 OK`: Successfully checked contract status
- `400 Bad Request`: Invalid address format
- `500 Internal Server Error`: RPC or network error

## Setup

### Requirements
- Go 1.22+
- Redis (optional, for distributed caching)
- Infura API Key (for ENS resolution)

### Environment Variables
Copy `.env.example` to `.env` and configure:
```env
SERVER_PORT=8080
SERVER_HOST=localhost
ENS_PROVIDER_URL=https://mainnet.infura.io/v3/your-project-id
CACHE_TYPE=redis  # or "memory" for in-memory cache

# JWT Configuration
JWT_SECRET_KEY=your-256-bit-secret
JWT_DURATION_MINUTES=60
```

### Development
1. Install dependencies:
```bash
go mod download
```

2. Run tests:
```bash
go test ./...
```

3. Start server:
```bash
go run cmd/server/main.go
```

## Performance
- Sub-millisecond response times for address validation
- ~100ms average for ENS resolution (with caching)
- Scales horizontally with Redis caching

## Contributing
Pull requests welcome! Please read CONTRIBUTING.md first.

## License
MIT