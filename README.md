# Ethereum Address Tools

A modern web application for Ethereum address validation, ENS resolution, and contract detection.

## Features

- **Address Validation**: Validate Ethereum addresses with EIP-55 checksum verification
- **ENS Resolution**: Resolve Ethereum Name Service (ENS) domains to addresses
- **Contract Detection**: Check if an address is a smart contract
- **JWT Authentication**: Secure API endpoints with JWT token-based authentication
- **Dark/Light Mode**: Automatic theme detection with manual toggle
- **Modern UI**: Clean, responsive interface with real-time feedback

## API Endpoints

- `POST /v1/token`: Generate a new API token
- `GET /v1/validate/{address}`: Validate an Ethereum address
- `GET /v1/resolveEns/{name}`: Resolve an ENS name to an address
- `GET /v1/isContract/{address}`: Check if an address is a contract

## Example Usage

### Address Validation
```bash
# First, generate a token
curl -X POST http://localhost:8080/v1/token

# Then use the token to validate an address
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/v1/validate/0x742d35Cc6634C0532925a3b844Bc454e4438f44e
```

### ENS Resolution
```bash
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/v1/resolveEns/vitalik.eth
```

### Contract Detection
```bash
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/v1/isContract/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
```

## Development

### Prerequisites
- Go 1.21 or higher
- Access to an Ethereum node (e.g., Infura)

### Configuration
Create a `config.yaml` file:
```yaml
server:
  host: localhost
  port: 8080
  env: development

ethereum:
  provider_url: https://mainnet.infura.io/v3/your-project-id
  cache_duration: 3600  # in seconds

jwt:
  secret: your-secret-key
  expiry: 3600  # in seconds
```

### Running the Server
```bash
go run cmd/server/main.go
```

### Testing
```bash
go test ./...
```

## Architecture

- **Plugin Architecture**: Modular design for easy extension
- **Caching**: Built-in caching for ENS and contract detection
- **Middleware**: JWT authentication and request timeout handling
- **Structured Logging**: Detailed logging with different levels
- **Error Handling**: Comprehensive error handling and user feedback

## Security Features

- JWT-based authentication
- Request timeouts
- Input validation
- Secure response headers
- Rate limiting (configurable)

## Performance

- Response caching
- Connection pooling
- Efficient error handling
- Optimized API response formats