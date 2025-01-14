# AegisGate

AegisGate is a powerful and flexible API Gateway written in Go, designed to manage and secure your microservices architecture. It provides robust routing, service discovery, and request handling capabilities with a simple YAML configuration.

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/yourusername/aegisgate)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🚀 **High Performance**: Built with Go for optimal performance and low resource usage
- 🛠 **Simple Configuration**: Easy-to-understand YAML configuration
- 🔄 **Dynamic Routing**: Flexible path-based routing with method filtering
- ⚡ **Hot Reload**: Configuration changes without restart
- 🛡 **Path Stripping**: Optional path stripping for clean forwarding
- 🔒 **Safe Shutdown**: Graceful shutdown with request draining and connection handling

## Quick Start

### Using Docker

1. Create your configuration file (e.g., `config.yaml`):

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  debug: false

services:
  - name: "example-api"
    base_path: "/api"
    target_url: "https://api.example.com"
    routes:
      - path: "/*"
        methods: ["FULL"]
        strip_path: true
```

2. Run using Docker Compose:

```bash
docker-compose up -d
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/aegisgate.git
cd aegisgate
```

2. Build the binary:
```bash
go build -o aegisgate cmd/aegisgate/main.go
```

3. Run the gateway:
```bash
CONFIG_PATH=./config.yaml ./aegisgate
```

## Configuration

### Server Configuration

```yaml
server:
  host: "0.0.0.0"  # Bind address
  port: 8080       # Listen port
  debug: false     # Enable debug mode
```

### Service Configuration

```yaml
services:
  - name: "service-name"           # Service identifier
    base_path: "/path"            # Base path for routing
    target_url: "http://target"   # Target service URL
    routes:
      - path: "/*"                # Route path pattern
        methods: ["GET", "POST"]   # Allowed methods
        strip_path: true          # Strip base path
```

Available method configurations:
- `FULL`: All HTTP methods
- `CRUD`: GET, POST, PUT, PATCH, DELETE
- `RO`: GET, HEAD
- `RW`: GET, POST, PUT, PATCH
- Individual methods: `["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "TRACE", "CONNECT"]`

## Docker Support

The project includes Docker support out of the box:

```yaml
version: '3.8'

services:
  aegisgate:
    build:
      context: .
      dockerfile: build/Docker/Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
    environment:
      - CONFIG_PATH=/app/config.yaml
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    restart: unless-stopped
```

## Health Check

AegisGate provides a health check endpoint at `/health` that returns HTTP 200 when the service is healthy.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 