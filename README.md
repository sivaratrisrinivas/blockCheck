# Ethereum Address Validator & ENS Resolver

A high-performance API service for validating Ethereum addresses, resolving ENS names, and detecting smart contracts.

## Features

✅ Ethereum Address Format Validation
- EIP-55 checksum validation
- Response time: ~1.89ms

✅ ENS Name Resolution
- Resolves ENS names to addresses
- Caches results for improved performance
- Response time: ~127ms (first request), ~80ms (cached)

✅ Contract Detection
- Detects if an address is a contract or EOA
- Response time: ~62ms

✅ Security & Authentication
- JWT-based authentication
- API key generation
- Protected endpoints
- Response time: ~759μs

✅ Caching System
- In-memory caching
- Redis support (optional)
- Configurable TTL

❌ Rate Limiting (Coming Soon)
❌ Prometheus Metrics (Coming Soon)
❌ Front-end Interface (Coming Soon)

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```
Response: `200 OK` if service is healthy
Response time: ~29.8μs

#### Generate API Token
```
POST /v1/token
```
Response:
```json
{
  "api_key": "520db2e2-141c-463e-bc7e-0246f5b09ecd",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Endpoints
All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

#### Validate Address
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

#### Resolve ENS Name
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

#### Check Contract Status
```
GET /v1/isContract/{address}
```
Example Response:
```json
{
  "address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
  "isContract": true
}
```

## Setup

### Requirements
- Go 1.22+
- Redis (optional, for distributed caching)

### Environment Variables
```
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# ENS Configuration
ENS_PROVIDER_URL=https://mainnet.infura.io/v3/your-project-id
ENS_TIMEOUT_SECONDS=10
ENS_RETRY_ATTEMPTS=3

# Cache Configuration
CACHE_TYPE=memory  # or "redis"
CACHE_TTL_MINUTES=60

# Redis Configuration (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET_KEY=your-256-bit-secret
JWT_DURATION_MINUTES=60

# API Configuration
ENABLE_RATE_LIMIT=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION_SECONDS=60
```

### Installation
```bash
git clone https://github.com/sivaratrisrinivas/web3/blockCheck
cd blockCheck
go mod download
go run cmd/server/main.go
```

## Performance
Latest benchmarks from testing:
- Health Check: ~29.8μs
- Address Validation: ~1.89ms
- ENS Resolution: ~127ms (first request), ~80ms (cached)
- Contract Detection: ~62ms
- Token Generation: ~759μs

## Development
```bash
# Run tests
go test ./...

# Run with hot reload
go install github.com/cosmtrek/air@latest
air
```

## Contributing
Pull requests welcome! Please read CONTRIBUTING.md first.

## License
MIT