# IPMI API

A lightweight REST API wrapper for the `go-ipmi` library that exposes all IPMI functionality over HTTP with minimal dependencies and maximum simplicity.

## Features

- Exposes IPMI functionality through a RESTful API
- Connection pooling for efficient IPMI communication
- Header-based authentication for flexible credential management
- Standardized JSON responses
- Health check endpoints
- Docker support for easy deployment
- Interactive API documentation with Scalar

## API Endpoints

- `GET /openapi.json` - OpenAPI specification for the service
- `GET /docs` - Scalar API reference, based on the `openapi.json` specification

### Health Checks
- `GET /health` - Basic health check
- `GET /health/ready` - Readiness probe
- `GET /health/live` - Liveness probe

### System Information
- `GET /api/v1/ipmi/system/info` - Get device ID and system info
- `GET /api/v1/ipmi/system/guid` - Get system GUID
- `GET /api/v1/ipmi/system/boot-options` - Get boot device options *(Not implemented)*
- `POST /api/v1/ipmi/system/boot-device` - Set boot device *(Partially implemented)*

### Power Management
- `GET /api/v1/ipmi/power/status` - Get power status
- `POST /api/v1/ipmi/power/on` - Power on
- `POST /api/v1/ipmi/power/off` - Power off
- `POST /api/v1/ipmi/power/cycle` - Power cycle
- `POST /api/v1/ipmi/power/reset` - Hard reset

### Sensors & Monitoring
- `GET /api/v1/ipmi/sensors` - Get all sensors
- `GET /api/v1/ipmi/sdr` - Get sensor data records *(Not implemented)*
- `GET /api/v1/ipmi/sel` - Get system event log
- `POST /api/v1/ipmi/sel/clear` - Clear event log

### Chassis Management
- `GET /api/v1/ipmi/chassis/status` - Get chassis status
- `POST /api/v1/ipmi/chassis/identify` - Chassis identify LED *(Not implemented)*
- `GET /api/v1/ipmi/chassis/bootdev` - Get boot device
- `POST /api/v1/ipmi/chassis/bootdev/set` - Set boot device *(Partially implemented)*

### FRU (Field Replaceable Unit)
- `GET /api/v1/ipmi/fru` - Get all FRU data
- `GET /api/v1/ipmi/fru/{fru-id}` - Get specific FRU *(Not implemented)*

### User Management
- `GET /api/v1/ipmi/users` - List users
- `POST /api/v1/ipmi/users/create` - Create user
- `GET /api/v1/ipmi/users/{user-id}` - Get user
- `PUT /api/v1/ipmi/users/{user-id}` - Update user
- `DELETE /api/v1/ipmi/users/{user-id}` - Delete user

### LAN Configuration
- `GET /api/v1/ipmi/lan` - Get LAN configuration
- `POST /api/v1/ipmi/lan` - Set LAN configuration *(Partially implemented)*

## Authentication

The API supports two methods for providing IPMI credentials:

### Environment Variables (Default)
Set these environment variables to provide default credentials for all requests:
```bash
DEFAULT_IPMI_HOST=<default-host>
DEFAULT_IPMI_PORT=623            # Optional, defaults to 623
DEFAULT_IPMI_USERNAME=<default-username>
DEFAULT_IPMI_PASSWORD=<default-password>
```

### Request Headers (Override)
Provide these headers in your requests to override default credentials:
```
X-IPMI-Host: <target-host>
X-IPMI-Port: <target-port>       # Optional, defaults to 623
X-IPMI-Username: <username>
X-IPMI-Password: <password>
```

## Configuration

The following environment variables can be used to configure the API:

```bash
# API Server Configuration
IPMI_API_PORT=8080              # API server port
IPMI_API_HOST=0.0.0.0           # API server bind address
IPMI_TIMEOUT=30s                # IPMI command timeout
IPMI_POOL_SIZE=100              # Max concurrent connections
LOG_LEVEL=info                  # Logging level

# Default IPMI Authentication (Optional)
DEFAULT_IPMI_HOST=              # Default IPMI host for all requests
DEFAULT_IPMI_PORT=623           # Default IPMI port
DEFAULT_IPMI_USERNAME=          # Default IPMI username
DEFAULT_IPMI_PASSWORD=          # Default IPMI password
```

## Installation

### Using Go

```bash
# Clone the repository
git clone https://github.com/yourusername/ipmi-api.git
cd ipmi-api

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

### Using Docker

```bash
# Build the Docker image
docker build -t ipmi-api .

# Run the container
docker run -d \
  --name ipmi-api \
  -p 8080:8080 \
  -e DEFAULT_IPMI_HOST=your.ipmi.host \
  -e DEFAULT_IPMI_USERNAME=admin \
  -e DEFAULT_IPMI_PASSWORD=secret \
  ipmi-api
```

## API Documentation

Interactive API documentation is available at `/docs` endpoint when the server is running. This documentation is powered by [Scalar](https://scalar.com/), a modern and lightweight alternative to Swagger UI.

## Usage Examples

### Get Power Status
```bash
curl -H "X-IPMI-Host: 192.168.1.100" \
     -H "X-IPMI-Username: admin" \
     -H "X-IPMI-Password: password" \
     http://localhost:8080/api/v1/ipmi/power/status
```

### Power On System
```bash
curl -X POST \
     -H "X-IPMI-Host: 192.168.1.100" \
     -H "X-IPMI-Username: admin" \
     -H "X-IPMI-Password: password" \
     http://localhost:8080/api/v1/ipmi/power/on
```

## Development

### Prerequisites
- Go 1.21 or later
- Docker (for containerization)

### Building
```bash
# Build the binary
go build -o ipmi-api main.go

# Run tests
go test ./...
```

## License

This project is licensed under the Apache-2.0 License - see the LICENSE file for details.