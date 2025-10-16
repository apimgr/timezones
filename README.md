# Timezones API

Comprehensive timezone data API with UTC offsets, DST information, and regional identifiers. Single static binary with embedded JSON data.

## About

**Timezones API** provides a fast, lightweight HTTP API for querying timezone information. All timezone data is embedded in the binary for zero-dependency deployment.

### Features

- **Comprehensive Data**: Complete timezone information including UTC offsets, DST status, abbreviations, and IANA identifiers
- **Fast Search**: Search by timezone name, abbreviation, UTC identifier, or offset
- **Single Binary**: All assets (JSON data, HTML templates, CSS, JS) embedded via go:embed
- **REST API**: Clean JSON API with standard response format
- **Web UI**: Modern dark-themed web interface for browsing API documentation
- **Admin Panel**: Protected admin endpoints for settings management
- **Docker Ready**: Multi-arch Docker images (amd64, arm64) with Alpine runtime
- **Production Grade**: Health checks, logging, CORS support, authentication

## Production Installation

### Download Binary

```bash
# Linux AMD64
curl -LO https://github.com/casjay/timezones/releases/latest/download/timezones-linux-amd64
chmod +x timezones-linux-amd64
sudo mv timezones-linux-amd64 /usr/local/bin/timezones

# Linux ARM64
curl -LO https://github.com/casjay/timezones/releases/latest/download/timezones-linux-arm64
chmod +x timezones-linux-arm64
sudo mv timezones-linux-arm64 /usr/local/bin/timezones
```

### Run Server

```bash
# Start with default settings (port 8080)
timezones

# Specify port and admin credentials
timezones --port 8080 --admin-user admin --admin-password secretpass

# With custom directories
timezones \
  --port 8080 \
  --config /etc/timezones \
  --data /var/lib/timezones \
  --logs /var/log/timezones
```

### Environment Variables

```bash
export PORT=8080
export ADDRESS=0.0.0.0
export ADMIN_USER=administrator
export ADMIN_PASSWORD=changeme
export CONFIG_DIR=/etc/timezones
export DATA_DIR=/var/lib/timezones
export LOGS_DIR=/var/log/timezones
```

## Docker Deployment

### Docker Compose (Production)

```yaml
# docker-compose.yml
services:
  timezones:
    image: ghcr.io/casjay/timezones:latest
    container_name: timezones
    restart: unless-stopped
    environment:
      - PORT=80
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=changeme
    volumes:
      - ./rootfs/config/timezones:/config
      - ./rootfs/data/timezones:/data
      - ./rootfs/logs/timezones:/logs
    ports:
      - "172.17.0.1:64180:80"
```

```bash
# Start service
docker-compose up -d

# View logs
docker-compose logs -f

# Access API
curl http://172.17.0.1:64180/api/v1/timezones
```

### Docker Run

```bash
docker run -d \
  --name timezones \
  --restart unless-stopped \
  -e PORT=80 \
  -e ADMIN_USER=administrator \
  -e ADMIN_PASSWORD=changeme \
  -v $(pwd)/config:/config \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/logs \
  -p 8080:80 \
  ghcr.io/casjay/timezones:latest
```

## API Usage

### Endpoints

#### Get Raw JSON Data
```bash
curl http://localhost:8080/api/v1/timezones.json
```

#### Get All Timezones (API Response)
```bash
curl http://localhost:8080/api/v1/timezones
```

#### Search Timezones
```bash
# Search for "pacific"
curl http://localhost:8080/api/v1/timezones/search?q=pacific

# Search for "new york"
curl http://localhost:8080/api/v1/timezones/search?q=new+york
```

#### Get by UTC Offset
```bash
# Get timezones with offset -5 (EST)
curl http://localhost:8080/api/v1/timezones/offset/-5

# Get timezones with offset +9 (JST)
curl http://localhost:8080/api/v1/timezones/offset/9
```

#### Get by Abbreviation
```bash
curl http://localhost:8080/api/v1/timezones/abbr/EST
curl http://localhost:8080/api/v1/timezones/abbr/PST
```

#### Get by UTC Identifier
```bash
curl http://localhost:8080/api/v1/timezones/utc/America/New_York
curl http://localhost:8080/api/v1/timezones/utc/Europe/London
```

#### Get by Value
```bash
curl http://localhost:8080/api/v1/timezones/value/Eastern%20Standard%20Time
```

#### Get Statistics
```bash
curl http://localhost:8080/api/v1/stats
```

### Response Format

```json
{
  "success": true,
  "data": [
    {
      "value": "Eastern Standard Time",
      "abbr": "EST",
      "offset": -5,
      "isdst": false,
      "text": "(UTC-05:00) Eastern Time (US & Canada)",
      "utc": [
        "America/Detroit",
        "America/New_York",
        "America/Toronto"
      ]
    }
  ]
}
```

### Admin Endpoints (Protected)

Admin endpoints require Bearer token authentication.

```bash
# Get all settings
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/admin/settings

# Set a setting
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"key": "maintenance_mode", "value": "false"}' \
  http://localhost:8080/api/v1/admin/settings

# Delete a setting
curl -X DELETE \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/admin/settings/maintenance_mode
```

**Finding Your Token**: After first run, check `{CONFIG_DIR}/credentials.txt` for your admin token.

## Development

### Requirements

- Go 1.23+
- Make
- Docker (optional, for containerized builds)

### Build from Source

```bash
# Clone repository
git clone https://github.com/casjay/timezones.git
cd timezones

# Build for all platforms
make build

# Build for host platform only (quick)
make quick

# Run tests
make test
```

### Build System & Testing

**Makefile Targets**:
- `make build` - Build binaries for all 8 platforms (Linux, Windows, macOS, FreeBSD × amd64, arm64)
- `make test` - Run tests
- `make docker` - Build and push multi-arch Docker images
- `make docker-dev` - Build development Docker image (local only)
- `make clean` - Remove build artifacts
- `make release` - Create GitHub release with binaries
- `make quick` - Build for host platform only (fast iteration)

**Platforms**: linux-amd64, linux-arm64, windows-amd64, windows-arm64, darwin-amd64, darwin-arm64, freebsd-amd64, freebsd-arm64

**Versioning**: Version stored in `release.txt`. Auto-increments after successful `make release`.

### Development Mode

```bash
# Build dev Docker image
make docker-dev

# Run with docker-compose (uses /tmp for ephemeral storage)
docker-compose -f docker-compose.test.yml up -d

# Test API
curl http://localhost:64181/api/v1/timezones

# Cleanup
docker-compose -f docker-compose.test.yml down
sudo rm -rf /tmp/timezones/rootfs
```

### CI/CD

**GitHub Actions** (`.github/workflows/`):
- `release.yml` - Build binaries and create GitHub releases
- `docker.yml` - Build and push multi-arch Docker images

**Jenkins** (`Jenkinsfile`):
- Multi-architecture pipeline (amd64, arm64)
- Parallel builds and tests
- Docker image builds and registry push
- GitHub release creation

**Triggers**: Push to main/master, monthly schedule (1st of month), manual dispatch

## License & Credits

**License**: MIT

**Author**: casjay

**Repository**: https://github.com/casjay/timezones

---

Built with Go • Single Static Binary • Production Ready
