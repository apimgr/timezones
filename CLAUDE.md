# 🚀 SPEC.md Updates - Production Specifications

**Last Updated**: 2025-01-14
**Version**: 2.0
**Changes**: URL display, Dockerfile runtime, docker-compose standards, Makefile improvements

---

## 📋 Critical Updates to Apply

When using the generic SPEC.md template, apply these production-grade updates:

### 1. **URL Display Standards** ⭐ NEW

**Rule**: Never show `localhost`, `127.0.0.1`, or `0.0.0.0` to users. Always display the most relevant accessible URL.

**Priority**:
1. **FQDN** (if hostname resolves)
2. **Public IP** (outbound IP)
3. **Hostname** (if available)
4. **Fallback** (`<your-host>`)

**Implementation** - Add to `src/database/credentials.go`:

```go
// getAccessibleURL returns the most relevant URL for accessing the server
// Priority: FQDN > hostname > public IP > fallback
// NEVER shows localhost, 127.0.0.1, or 0.0.0.0
func getAccessibleURL(port string) string {
	// Try to get hostname
	hostname, err := os.Hostname()
	if err == nil && hostname != "" && hostname != "localhost" {
		// Try to resolve hostname to see if it's a valid FQDN
		if addrs, err := net.LookupHost(hostname); err == nil && len(addrs) > 0 {
			return fmt.Sprintf("http://%s:%s", hostname, port)
		}
	}

	// Try to get outbound IP (most likely accessible IP)
	if ip := getOutboundIP(); ip != "" {
		return fmt.Sprintf("http://%s:%s", ip, port)
	}

	// Fallback to hostname if we have one
	if hostname != "" && hostname != "localhost" {
		return fmt.Sprintf("http://%s:%s", hostname, port)
	}

	// Last resort: use a generic message
	return fmt.Sprintf("http://<your-host>:%s", port)
}

// getOutboundIP gets the preferred outbound IP of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
```

**Update** `SaveCredentialsToFile` signature to accept port:
```go
func SaveCredentialsToFile(creds *AdminCredentials, configDir, port string) error {
	serverURL := getAccessibleURL(port)

	content := fmt.Sprintf(`{PROJECT_NAME} - ADMIN CREDENTIALS
========================================
WEB UI LOGIN:
  URL:      %s/admin
  Username: %s

API ACCESS:
  URL:      %s/api/v1/admin
  Header:   Authorization: Bearer %s

CREDENTIALS:
  Username: %s
  Password: %s
  Token:    %s

Created: %s
========================================`, serverURL, creds.Username, serverURL, creds.Token,
    creds.Username, creds.Password, creds.Token, time.Now().Format("2006-01-02 15:04:05"))
	//...
}
```

**In main.go** - Save credentials AFTER port is determined:
```go
// Port resolution first
if port == "" {
	// ... port determination logic
}

// THEN save credentials with actual port
if creds.Token != "" {
	if err := database.SaveCredentialsToFile(creds, configDir, port); err != nil {
		log.Printf("Warning: Failed to save credentials file: %v", err)
	} else {
		log.Printf("⚠️  Access URL: %s", getAccessibleURL(port))
	}
}
```

---

### 2. **Dockerfile - Alpine Runtime** ⭐ UPDATED

**Replace**: `FROM scratch` runtime
**With**: `FROM alpine:latest` runtime

**Why**: Need curl and bash for health checks and debugging

```dockerfile
# ============================================
# Build stage
# ============================================
FROM golang:alpine AS builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/

# Build static binary with all assets embedded
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE} -w -s" \
    -a -installsuffix cgo \
    -o {projectname} \
    ./src

# ============================================
# Runtime stage - Alpine with minimal tools
# ============================================
FROM alpine:latest

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Install runtime dependencies (curl, bash)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Copy binary to /usr/local/bin
COPY --from=builder /build/{projectname} /usr/local/bin/{projectname}
RUN chmod +x /usr/local/bin/{projectname}

# Environment variables
ENV PORT=80 \
    CONFIG_DIR=/config \
    DATA_DIR=/data \
    LOGS_DIR=/logs \
    ADDRESS=0.0.0.0 \
    DB_PATH=/data/db/{projectname}.db

# Create directories
RUN mkdir -p /config /data /data/db /logs && \
    chown -R 65534:65534 /config /data /logs

# Metadata labels (OCI standard)
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.authors="{organization}" \
      org.opencontainers.image.url="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.source="https://github.com/{organization}/{projectname}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.vendor="{organization}" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="{projectname}" \
      org.opencontainers.image.description="{Project description} - Single static binary" \
      org.opencontainers.image.documentation="https://github.com/{organization}/{projectname}/blob/main/docs/README.md" \
      org.opencontainers.image.base.name="alpine:latest"

# Expose default port
EXPOSE 80

# Create mount points for volumes
VOLUME ["/config", "/data", "/logs"]

# Run as non-root user (nobody)
USER 65534:65534

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/{projectname}", "--status"]

# Run
ENTRYPOINT ["/usr/local/bin/{projectname}"]
CMD ["--port", "80"]
```

**Key Changes**:
- Runtime base: `alpine:latest` (not scratch)
- Includes: curl, bash, ca-certificates, tzdata
- Binary location: `/usr/local/bin/{projectname}` (not root)
- SQLite DB location: `/data/db/{projectname}.db`
- Internal port: 80
- All OCI metadata labels

---

### 3. **docker-compose.yml - Production Standards** ⭐ UPDATED

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition
- ✅ Use pre-built images from registry
- ✅ Custom network: `{projectname}` (external: false)
- ✅ Volume structure: `./rootfs/{type}/{servicename}`
- ✅ Production port: `172.17.0.1:{randomport}:80`

```yaml
# Production Docker Compose Configuration
# Uses ./rootfs for persistent storage
# External port: 172.17.0.1:64180:80 (Docker bridge network)

services:
  {projectname}:
    image: ghcr.io/{organization}/{projectname}:latest
    container_name: {projectname}
    restart: unless-stopped

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      # Uncomment and set for first deployment
      #- ADMIN_USER=administrator
      #- ADMIN_PASSWORD=changeme
      #- ADMIN_TOKEN=your-token-here

    volumes:
      - ./rootfs/config/{projectname}:/config
      - ./rootfs/data/{projectname}:/data
      - ./rootfs/logs/{projectname}:/logs

    ports:
      - "172.17.0.1:64180:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL database
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres
  #   restart: unless-stopped
  #
  #   environment:
  #     POSTGRES_DB: {projectname}
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: changeme
  #
  #   volumes:
  #     - ./rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD-SHELL", "pg_isready -U {projectname}"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: MariaDB database
  # mariadb:
  #   image: mariadb:11-alpine
  #   container_name: {projectname}-mariadb
  #   restart: unless-stopped
  #
  #   environment:
  #     MYSQL_DATABASE: {projectname}
  #     MYSQL_USER: {projectname}
  #     MYSQL_PASSWORD: changeme
  #     MYSQL_RANDOM_ROOT_PASSWORD: "yes"
  #
  #   volumes:
  #     - ./rootfs/db/mariadb:/var/lib/mysql
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "healthcheck.sh", "--connect", "--innodb_initialized"]
  #     interval: 10s
  #     timeout: 5s
  #     retries: 5

  # Optional: Redis cache
  # redis:
  #   image: redis:7-alpine
  #   container_name: {projectname}-redis
  #   restart: unless-stopped
  #
  #   volumes:
  #     - ./rootfs/db/redis:/data
  #
  #   networks:
  #     - {projectname}
  #
  #   healthcheck:
  #     test: ["CMD", "redis-cli", "ping"]
  #     interval: 10s
  #     timeout: 3s
  #     retries: 3

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

### 4. **docker-compose.test.yml - Development Standards** ⭐ UPDATED

**Standards**:
- ❌ NO `version:` field
- ❌ NO `build:` definition (use `{projectname}:dev` image)
- ✅ Ephemeral storage: `/tmp/{projectname}/rootfs`
- ✅ Development port: `{randomport}:80` (e.g., 64181:80)
- ✅ Same network name: `{projectname}`

```yaml
# Development/Testing Docker Compose Configuration
# Uses /tmp for ephemeral storage
# External port: 64181:80 (random unused port in 64xxx range)

services:
  {projectname}:
    image: {projectname}:dev
    container_name: {projectname}-test
    restart: "no"

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/{projectname}.db
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=testpass123
      - DEV=true

    volumes:
      - /tmp/{projectname}/rootfs/config/{projectname}:/config
      - /tmp/{projectname}/rootfs/data/{projectname}:/data
      - /tmp/{projectname}/rootfs/logs/{projectname}:/logs

    ports:
      - "64181:80"

    networks:
      - {projectname}

    healthcheck:
      test: ["CMD", "/usr/local/bin/{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

  # Optional: PostgreSQL for testing
  # postgres:
  #   image: postgres:16-alpine
  #   container_name: {projectname}-postgres-test
  #   restart: "no"
  #
  #   environment:
  #     POSTGRES_DB: {projectname}_test
  #     POSTGRES_USER: {projectname}
  #     POSTGRES_PASSWORD: testpass
  #
  #   volumes:
  #     - /tmp/{projectname}/rootfs/db/postgres:/var/lib/postgresql/data
  #
  #   networks:
  #     - {projectname}

networks:
  {projectname}:
    name: {projectname}
    external: false
    driver: bridge
```

---

### 5. **Makefile - Docker Improvements** ⭐ UPDATED

Add `docker-dev` target for local development:

```makefile
# Build and push multi-platform Docker images (release)
docker:
	@echo "Building multi-platform Docker images..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) \
		--push \
		.
	@echo "✓ Docker images pushed to ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION)"

# Build Docker image for development (local only, not pushed)
docker-dev:
	@echo "Building development Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(PROJECTNAME):dev \
		.
	@echo "✓ Docker development image built: $(PROJECTNAME):dev"
```

---

### 6. **Jenkinsfile** ⭐ NEW

Multi-architecture CI/CD pipeline:

```groovy
// Jenkinsfile for {projectname}
// Multi-architecture build pipeline (amd64, arm64)
// Server: jenkins.casjay.cc

pipeline {
    agent none

    environment {
        PROJECTNAME = '{projectname}'
        PROJECTORG = '{organization}'
        REGISTRY = 'ghcr.io'
        IMAGE_NAME = "${REGISTRY}/${PROJECTORG}/${PROJECTNAME}"

        VERSION = sh(script: 'cat release.txt 2>/dev/null || echo "0.0.1"', returnStdout: true).trim()
        COMMIT = sh(script: 'git rev-parse --short HEAD 2>/dev/null || echo "unknown"', returnStdout: true).trim()
        BUILD_DATE = sh(script: 'date -u +%Y-%m-%dT%H:%M:%SZ', returnStdout: true).trim()
    }

    stages {
        stage('Checkout') {
            agent { label 'amd64' }
            steps {
                checkout scm
                echo "Building ${PROJECTNAME} ${VERSION} (${COMMIT})"
            }
        }

        stage('Test') {
            parallel {
                stage('Test on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Running tests on AMD64...'
                        sh 'make test'
                    }
                }
                stage('Test on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Running tests on ARM64...'
                        sh 'make test'
                    }
                }
            }
        }

        stage('Build Binaries') {
            parallel {
                stage('Build on AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building binaries on AMD64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-amd64'
                        stash includes: 'releases/**/*', name: 'release-amd64'
                    }
                }
                stage('Build on ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building binaries on ARM64...'
                        sh 'make build'
                        stash includes: 'binaries/**/*', name: 'binaries-arm64'
                        stash includes: 'releases/**/*', name: 'release-arm64'
                    }
                }
            }
        }

        stage('Build Docker Images') {
            parallel {
                stage('Build Docker AMD64') {
                    agent { label 'amd64' }
                    steps {
                        echo 'Building Docker image for AMD64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/amd64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-amd64 \
                                    -t ${IMAGE_NAME}:latest-amd64 \
                                    .
                            """
                        }
                    }
                }
                stage('Build Docker ARM64') {
                    agent { label 'arm64' }
                    steps {
                        echo 'Building Docker image for ARM64...'
                        script {
                            sh """
                                docker build \
                                    --platform linux/arm64 \
                                    --build-arg VERSION=${VERSION} \
                                    --build-arg COMMIT=${COMMIT} \
                                    --build-arg BUILD_DATE=${BUILD_DATE} \
                                    -t ${IMAGE_NAME}:${VERSION}-arm64 \
                                    -t ${IMAGE_NAME}:latest-arm64 \
                                    .
                            """
                        }
                    }
                }
            }
        }

        stage('Push Docker Images') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Pushing Docker images to registry...'
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'github-registry',
                        usernameVariable: 'REGISTRY_USER',
                        passwordVariable: 'REGISTRY_TOKEN'
                    )]) {
                        sh """
                            echo \$REGISTRY_TOKEN | docker login ${REGISTRY} -u \$REGISTRY_USER --password-stdin

                            docker push ${IMAGE_NAME}:${VERSION}-amd64
                            docker push ${IMAGE_NAME}:${VERSION}-arm64
                            docker push ${IMAGE_NAME}:latest-amd64
                            docker push ${IMAGE_NAME}:latest-arm64

                            docker manifest create ${IMAGE_NAME}:${VERSION} \
                                ${IMAGE_NAME}:${VERSION}-amd64 \
                                ${IMAGE_NAME}:${VERSION}-arm64

                            docker manifest create ${IMAGE_NAME}:latest \
                                ${IMAGE_NAME}:latest-amd64 \
                                ${IMAGE_NAME}:latest-arm64

                            docker manifest push ${IMAGE_NAME}:${VERSION}
                            docker manifest push ${IMAGE_NAME}:latest

                            docker logout ${REGISTRY}
                        """
                    }
                }
            }
        }

        stage('GitHub Release') {
            agent { label 'amd64' }
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                echo 'Creating GitHub release...'
                script {
                    unstash 'release-amd64'
                    withCredentials([string(
                        credentialsId: 'github-token',
                        variable: 'GITHUB_TOKEN'
                    )]) {
                        sh 'make release'
                    }
                }
            }
        }
    }

    post {
        success {
            echo "✅ Build successful for ${PROJECTNAME} ${VERSION}"
        }
        failure {
            echo "❌ Build failed for ${PROJECTNAME} ${VERSION}"
        }
        always {
            cleanWs()
        }
    }
}
```

---

### 7. **src/data Directory - JSON Data Files** ⭐ NEW

**Rule**: `src/data/` directory contains ONLY JSON files. No Go code, no other file types.

**Structure**:
```
src/
├── main.go                      # Embeds: //go:embed data/*.json
├── data/
│   ├── {datafile}.json          # JSON data ONLY
│   └── {other}.json             # Additional JSON files (if needed)
└── {service}/
    └── service.go               # NewService(jsonData []byte)
```

**Embedding from main.go:**

```go
// src/main.go
package main

import (
    _ "embed"
    "github.com/{org}/{project}/src/{service}"
)

//go:embed data/{datafile}.json
var jsonData []byte

func main() {
    svc, err := service.NewService(jsonData)
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

**Service receives embedded data:**

```go
// src/{service}/service.go
package service

import "encoding/json"

type Service struct {
    data MyDataType
}

func NewService(jsonData []byte) (*Service, error) {
    var data MyDataType
    err := json.Unmarshal(jsonData, &data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    return &Service{data: data}, nil
}
```

**Why This Pattern:**
- ✅ `src/data/` contains ONLY JSON files (no .go code)
- ✅ JSON is embedded in single static binary
- ✅ No copies, no symlinks, no duplicates
- ✅ Embedding happens from `main.go` at `src/` level
- ✅ Services receive data as parameter (clean dependency injection)
- ✅ True single binary with all assets embedded

**Example** (airports project):
```
src/
├── main.go                  # //go:embed data/airports.json
├── data/
│   └── airports.json        # 8.7MB (ONLY JSON, no .go files)
└── airports/
    └── data.go              # NewService(jsonData []byte)
```

**Important**:
- The `src/data/` directory contains ONLY JSON files
- NO .go files in `src/data/`
- NO symlinks or copies
- Embedding MUST be done from `src/main.go`
- Services accept `[]byte` parameter with embedded data

---

### 8. **README.md Structure** ⭐ UPDATED

**Order**: About → Production → Docker → API Usage → Development

```markdown
# {PROJECT_NAME}

Brief description

## About
- Features list
- Key capabilities

## Production Installation
- Binary installation
- Systemd service
- Environment variables

## Docker Deployment
- Docker Compose (production)
- Docker Compose (development)
- Docker run examples

## API Usage
- Quick examples
- Admin panel
- Documentation links

## Development

- Requirements
- Build System & Testing (Makefile targets, platforms, versioning)
- Development Mode (dev flags, debug features)
- CI/CD (GitHub Actions, Jenkins, ReadTheDocs)

## License & Credits
```

---

### 9. **Complete Project Layout** ⭐ REFERENCE

**Standard Go API Server Project Structure** (matches citylist/airports):

```
{projectname}/
├── .claude/
│   └── settings.local.json      # Claude Code settings
├── .github/
│   └── workflows/
│       └── build.yml            # GitHub Actions (build on push & monthly)
├── .gitattributes               # Git attributes
├── .gitignore                   # Git ignore patterns
├── .readthedocs.yml             # ReadTheDocs configuration
├── CLAUDE.md                    # Project specification (this file)
├── Dockerfile                   # Alpine-based multi-stage build
├── docker-compose.yml           # Production compose
├── docker-compose.test.yml      # Development/testing compose
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Jenkinsfile                  # CI/CD pipeline (jenkins.casjay.cc)
├── LICENSE.md                   # License file
├── Makefile                     # Build system (4 targets)
├── README.md                    # User documentation
├── release.txt                  # Version tracking (auto-increment)
├── binaries/                    # Built binaries (gitignored)
│   ├── {projectname}-linux-amd64
│   ├── {projectname}-linux-arm64
│   ├── {projectname}-windows-amd64.exe
│   ├── {projectname}-windows-arm64.exe
│   ├── {projectname}-macos-amd64
│   ├── {projectname}-macos-arm64
│   ├── {projectname}-bsd-amd64
│   ├── {projectname}-bsd-arm64
│   └── {projectname}            # Host platform binary
├── releases/                    # Release artifacts (gitignored)
├── rootfs/                      # Docker volumes (gitignored)
│   ├── config/
│   │   └── {projectname}/       # Service config
│   ├── data/
│   │   └── {projectname}/       # Service data
│   ├── logs/
│   │   └── {projectname}/       # Service logs
│   └── db/                      # External databases
│       ├── postgres/
│       ├── mariadb/
│       └── redis/
├── docs/                        # API documentation (optional)
│   ├── index.md                 # Documentation home
│   ├── API.md                   # Complete API reference
│   ├── SERVER.md                # Server administration guide
│   ├── README.md                # Documentation index
│   ├── mkdocs.yml               # MkDocs configuration (Dracula theme)
│   ├── requirements.txt         # Python dependencies for RTD
│   ├── stylesheets/
│   │   └── dracula.css          # Dracula theme CSS
│   ├── javascripts/
│   │   └── extra.js             # Custom JavaScript
│   └── overrides/               # MkDocs theme overrides
├── scripts/                     # Production scripts (optional)
│   ├── install.sh               # Installation script
│   ├── backup.sh                # Backup script
│   └── uninstall.sh             # Uninstallation script
├── test/                        # Test files (optional)
│   ├── test-docker.sh           # Docker testing script
│   ├── unit/                    # Unit tests
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end tests
└── src/                         # Source code
    ├── data/
    │   └── {datafile}.json      # JSON data ONLY (no .go files)
    ├── {service}/
    │   └── service.go           # Service logic
    ├── auth/                    # Authentication (optional)
    │   └── auth.go
    ├── database/
    │   ├── database.go          # DB connection & schema
    │   ├── auth.go              # Admin authentication
    │   ├── credentials.go       # Credential management
    │   └── settings.go          # Settings CRUD
    ├── paths/
    │   └── paths.go             # OS-specific directory detection
    ├── scheduler/               # Task scheduler (optional)
    │   └── scheduler.go
    ├── server/
    │   ├── server.go            # Server setup & routing
    │   ├── handlers.go          # Public handlers
    │   ├── admin_handlers.go    # Admin handlers
    │   ├── auth_middleware.go   # Auth middleware
    │   ├── templates.go         # Template embedding
    │   ├── static/              # Static assets (embedded)
    │   │   ├── css/
    │   │   │   └── main.css
    │   │   ├── js/
    │   │   │   └── main.js
    │   │   ├── images/
    │   │   └── favicon.png
    │   └── templates/           # HTML templates (embedded)
    │       ├── base.html
    │       ├── home.html
    │       └── {other}.html
    ├── utils/                   # Utility functions (optional)
    │   └── network.go
    └── main.go                  # Entry point
```

**Key Directories:**

| Directory | Purpose | Contents |
|-----------|---------|----------|
| `src/data/` | JSON data files | ONLY `.json` files (no code) |
| `src/{service}/` | Service logic | Service implementation, loads from `src/data/` |
| `src/database/` | Database layer | DB connection, auth, settings |
| `src/paths/` | OS detection | Platform-specific directory paths |
| `src/server/` | HTTP server | Routing, handlers, templates, static files |
| `src/utils/` | Utilities | Helper functions (network, etc.) |
| `docs/` | Documentation | API.md, SERVER.md, README.md (optional) |
| `scripts/` | Production scripts | install.sh, backup.sh, uninstall.sh (optional) |
| `test/` | Test files | test-docker.sh, unit/, integration/, e2e/ (optional) |
| `binaries/` | Build output | Platform-specific binaries (gitignored) |
| `releases/` | Release artifacts | Files for GitHub releases (gitignored) |
| `rootfs/` | Docker volumes | Persistent storage for containers (gitignored) |

**Required Root Files:**
- ✅ `CLAUDE.md` - Project specification
- ✅ `README.md` - User documentation
- ✅ `LICENSE.md` - License file
- ✅ `Makefile` - Build system (4 targets)
- ✅ `Dockerfile` - Alpine-based container definition
- ✅ `docker-compose.yml` - Production compose
- ✅ `docker-compose.test.yml` - Development compose
- ✅ `Jenkinsfile` - Jenkins CI/CD pipeline (jenkins.casjay.cc)
- ✅ `.github/workflows/release.yml and docker.yml` - GitHub Actions (push & monthly)
- ✅ `.readthedocs.yml` - ReadTheDocs configuration
- ✅ `release.txt` - Version tracking (semantic versioning)
- ✅ `go.mod` / `go.sum` - Go modules
- ✅ `.gitignore` - Git ignore patterns
- ✅ `.gitattributes` - Git attributes (optional)

**Required Source Packages:**
- ✅ `src/main.go` - Application entry point
- ✅ `src/database/` - Database layer with auth
- ✅ `src/paths/` - OS-specific path detection
- ✅ `src/server/` - HTTP server with templates & static files
- ✅ `src/{service}/` - Primary service logic
- ✅ `src/data/` - JSON data files (if applicable)

**Optional Source Packages:**
- `src/auth/` - Additional authentication logic
- `src/utils/` - Utility functions (network helpers, etc.)
- `src/scheduler/` - Task scheduler for periodic jobs

**Optional Directories:**
- `docs/` - **All API type documentation** (API.md, SERVER.md, etc.)
  - ReadTheDocs compatible
  - MkDocs with Material theme (Dracula color scheme)
  - Includes: index.md, mkdocs.yml, requirements.txt
- `scripts/` - **All production and install scripts** (install.sh, backup.sh, etc.)
- `test/` - **All test files** (test-docker.sh, unit/, integration/, e2e/)

---

### 10. **ReadTheDocs Configuration** ⭐ NEW

**Purpose**: Automatic documentation hosting with Dracula theme

**Files Required:**

1. **`.readthedocs.yml`** (root):
```yaml
version: 2

mkdocs:
  configuration: docs/mkdocs.yml
  fail_on_warning: false

formats:
  - pdf
  - epub

python:
  version: "3.11"
  install:
    - requirements: docs/requirements.txt
```

2. **`docs/mkdocs.yml`**:
```yaml
site_name: {projectname} Documentation
site_description: {Project description}
site_author: {organization}
site_url: https://{projectname}.readthedocs.io

repo_name: {organization}/{projectname}
repo_url: https://github.com/{organization}/{projectname}
edit_uri: edit/main/docs/

copyright: Copyright &copy; 2024 {organization}

theme:
  name: material
  custom_dir: overrides
  palette:
    scheme: dracula
    primary: deep purple
    accent: pink
  font:
    text: Roboto
    code: Fira Code
  features:
    - navigation.instant
    - navigation.tracking
    - navigation.tabs
    - navigation.sections
    - navigation.expand
    - navigation.top
    - search.suggest
    - search.highlight
    - search.share
    - content.code.copy
    - content.code.annotate
  icon:
    repo: fontawesome/brands/github

nav:
  - Home: index.md
  - API Reference: API.md
  - Server Guide: SERVER.md
  - GitHub: https://github.com/{organization}/{projectname}

markdown_extensions:
  - admonition
  - pymdownx.details
  - pymdownx.superfences
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite
  - pymdownx.snippets
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.tasklist:
      custom_checkbox: true
  - pymdownx.emoji
  - toc:
      permalink: true
  - tables
  - attr_list

plugins:
  - search
  - minify:
      minify_html: true

extra_css:
  - stylesheets/dracula.css

extra_javascript:
  - javascripts/extra.js
```

3. **`docs/requirements.txt`**:
```txt
mkdocs>=1.5.0
mkdocs-material>=9.5.0
pymdown-extensions>=10.7
mkdocs-minify-plugin>=0.8.0
mike>=2.0.0
```

4. **`docs/stylesheets/dracula.css`** - Dracula color scheme overrides
5. **`docs/javascripts/extra.js`** - Optional custom JavaScript

**Documentation Structure:**
```
docs/
├── index.md              # Documentation home
├── API.md                # API reference
├── SERVER.md             # Server administration
├── mkdocs.yml            # MkDocs configuration
├── requirements.txt      # Python dependencies
├── stylesheets/
│   └── dracula.css       # Dracula theme CSS
├── javascripts/
│   └── extra.js          # Custom JavaScript
└── overrides/            # Theme overrides (optional)
```

**Dracula Color Palette:**
- Background: `#282a36`
- Current Line: `#44475a`
- Foreground: `#f8f8f2`
- Comment: `#6272a4`
- Cyan: `#8be9fd`
- Green: `#50fa7b`
- Orange: `#ffb86c`
- Pink: `#ff79c6`
- Purple: `#bd93f9`
- Red: `#ff5555`
- Yellow: `#f1fa8c`

**Local Preview:**
```bash
cd docs
pip install -r requirements.txt
mkdocs serve
# Open http://localhost:8000
```

---

### 11. **GitHub Actions Workflows** ⭐ NEW

**Purpose**: Automated builds on push and monthly schedule (split into separate workflows)

**Files**:
- `.github/workflows/release.yml` - Binary builds and GitHub releases
- `.github/workflows/docker.yml` - Docker image builds and registry push

**Triggers** (both workflows):
- Push to `main` branch
- Monthly schedule (1st of month, 3:00 AM UTC)
- Manual trigger (workflow_dispatch)

---

**Workflow 1: release.yml** - Binary Release

**Jobs:**

1. **test**
   - Runs `make test` with Go 1.23
   - Validates all tests pass before building

2. **build-and-release**
   - Reads version from `release.txt` (does NOT update it)
   - Runs `make build` for all 8 platforms
   - Deletes existing release if exists (always recreate)
   - Creates new GitHub release `{VERSION}`
   - Attaches all platform binaries
   - Uploads artifacts (90 day retention)

---

**Workflow 2: docker.yml** - Docker Release

**Jobs:**

1. **build-and-push**
   - Reads version from `release.txt`
   - Builds multi-arch Docker images (amd64, arm64)
   - Pushes to `ghcr.io/{organization}/{projectname}`
   - Tags: `latest`, `{VERSION}`, `{branch}-{sha}`
   - Uses GitHub cache for faster builds
   - Verifies images after push

**Configuration:**

```yaml
name: Build and Release

on:
  push:
    branches:
      - main
      - master
  schedule:
    - cron: '0 3 1 * *'  # Monthly on 1st at 3 AM UTC
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-binaries:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Read version
        run: |
          VERSION=$(cat release.txt)
          echo "VERSION=${VERSION}" >> $GITHUB_OUTPUT
      - name: Build
        run: VERSION=$VERSION make build
      - uses: actions/upload-artifact@v4
        with:
          name: binaries
          path: releases/*

  build-docker:
    runs-on: ubuntu-latest
    permissions:
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-qemu-action@v3
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:latest
            ghcr.io/${{ github.repository }}:${{ steps.version.outputs.VERSION }}

  create-release:
    needs: build-binaries
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          name: binaries
          path: releases/
      - name: Delete existing release
        run: gh release delete $VERSION -y || true
      - name: Create release
        run: gh release create $VERSION ./releases/* --title "$VERSION"
```

**Important:**
- Version is read from `release.txt`, NOT modified by Actions
- Release is always recreated (deleted first if exists)
- Builds all platforms (8 binaries)
- Multi-arch Docker images (amd64, arm64)
- Runs on every push to main and monthly

---

## 🎯 Quick Reference

### Container Tags
- **Production**: `ghcr.io/{organization}/{projectname}:latest`
- **Versioned**: `ghcr.io/{organization}/{projectname}:{VERSION}`
- **Development**: `{projectname}:dev`

### Port Mappings
- **Production**: `172.17.0.1:{randomport}:80`
- **Development**: `{randomport}:80`
- **Internal**: Always `80`

### Volume Structure
```
./rootfs/
├── config/{servicename}/
├── data/{servicename}/
├── logs/{servicename}/
└── db/{dbtype}/
```

### Releases Directory Structure



**./releases/** contains files for GitHub releases:



```

./releases/

├── {projectname}-linux-amd64

├── {projectname}-linux-arm64

├── {projectname}-macos-amd64

├── {projectname}-macos-arm64

├── {projectname}-windows-amd64.exe

├── {projectname}-windows-arm64.exe

├── {projectname}-bsd-amd64

├── {projectname}-bsd-arm64

├── {projectname}-{VERSION}-src.tar.gz    # Source archive (no VCS)

└── {projectname}-{VERSION}-src.zip       # Source archive for Windows

```



**Source Archives:**

- Created with `git archive` (excludes .git, .github, etc.)

- Both tar.gz (Linux/macOS) and zip (Windows) formats

- Includes all source code, docs, configs

- No binaries, no VCS files, no build artifacts




### Binary Requirements
- ✅ Static binary (CGO_ENABLED=0)
- ✅ All assets embedded:
  - HTML templates (go:embed in server package)
  - CSS, JS, images (go:embed in server package)
  - JSON data files (go:embed in main.go)
- ✅ True single binary - no external files needed
- ✅ Location: `/usr/local/bin/{projectname}`
- ✅ Output: `./binaries/` or `./releases/`

### URL Display
- ❌ Never: localhost, 127.0.0.1, 0.0.0.0
- ✅ Priority: FQDN → hostname → public IP → fallback

### Data Loading (JSON Files)
```
src/
├── main.go                # Embeds: //go:embed data/{file}.json
├── data/
│   └── {datafile}.json    # JSON files ONLY (no .go files)
└── {service}/
    └── service.go         # NewService(jsonData []byte)
```
- ✅ `src/data/` contains ONLY JSON files
- ✅ NO Go code in `src/data/`
- ✅ Embedded in main.go, passed to services as []byte
- ✅ True single binary - JSON embedded, not external files

---

**Apply these updates to SPEC.md for production-grade specifications** 🚀

---

### 12. **Web UI / Frontend Standards** ⭐ NEW

**Purpose**: Modern, responsive dark-themed web interface with Go html/template

**Technology Stack:**
- Go `html/template` for server-side rendering
- CSS3 with CSS custom properties (variables)
- Vanilla JavaScript (no frameworks)
- Embedded assets (go:embed static/* and templates/*)

**Note**: Static assets (CSS, JS, images), HTML templates, AND JSON data files are ALL embedded via `go:embed` in the single static binary. JSON files are embedded from `main.go` and passed to services.

---

**Structure:**
```
src/server/
├── static/
│   ├── css/
│   │   └── main.css        # Main stylesheet (~900 lines)
│   ├── js/
│   │   └── main.js         # JavaScript utilities (~130 lines)
│   ├── images/
│   │   └── favicon.png
│   └── manifest.json       # PWA manifest
├── templates/
│   ├── base.html           # Base template with header/footer
│   ├── home.html           # Homepage
│   ├── search.html         # Search page
│   ├── {feature}.html      # Feature-specific pages
│   └── config.html         # Admin configuration page
└── templates.go            # Template embedding logic
```

---

**CSS Framework** - `src/server/static/css/main.css`

**Design System:**

```css
/* CSS Variables (Dark Theme Default) */
:root {
  /* Colors */
  --bg-primary: #1a1a1a;        /* Main background */
  --bg-secondary: #2d2d2d;      /* Cards, sections */
  --bg-tertiary: #404040;       /* Hover states */
  --text-primary: #ffffff;      /* Main text */
  --text-secondary: #b0b0b0;    /* Secondary text */
  --text-tertiary: #808080;     /* Disabled text */
  
  /* Accent Colors */
  --accent-primary: #0066cc;    /* Primary actions */
  --accent-success: #00aa00;    /* Success states */
  --accent-warning: #ff9900;    /* Warnings */
  --accent-danger: #cc0000;     /* Errors */
  
  /* Layout */
  --border-color: #404040;
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  
  /* Typography */
  --font-primary: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  
  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 300ms ease;
  
  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.12);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.16);
  --shadow-lg: 0 10px 20px rgba(0,0,0,0.2);
}

/* Light Theme Override */
[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #e0e0e0;
  --text-primary: #1a1a1a;
  --text-secondary: #404040;
  --border-color: #e0e0e0;
}
```

---

**Components:**

**1. Buttons**
```css
.btn-primary {
  padding: var(--space-md) var(--space-xl);
  background: var(--accent-primary);
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-primary:hover {
  background: #0052a3;
  transform: translateY(-2px);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
```

**2. Cards**
```css
.card {
  background: var(--bg-secondary);
  padding: var(--space-xl);
  border-radius: 12px;
  transition: all var(--transition-normal);
  border: 2px solid transparent;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--accent-primary);
}
```

**3. Modals**
```css
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
  display: none;
}

.modal.active { display: block; }

.modal-backdrop {
  background: rgba(0,0,0,0.7);
}

.modal-content {
  background: var(--bg-secondary);
  max-width: 600px;
  margin: 100px auto;
  border-radius: 12px;
  padding: var(--space-xl);
  box-shadow: var(--shadow-lg);
}
```

**4. Toast Notifications**
```css
#toast-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 10000;
}

.toast {
  background: var(--bg-secondary);
  padding: var(--space-md) var(--space-lg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  border-left: 4px solid var(--accent-primary);
  min-width: 300px;
  animation: slideIn 0.3s ease;
}

.toast.success { border-left-color: var(--accent-success); }
.toast.error { border-left-color: var(--accent-danger); }
.toast.warning { border-left-color: var(--accent-warning); }
```

**5. Forms**
```css
.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.form-group input, .form-group select {
  padding: var(--space-md);
  background: var(--bg-tertiary);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
}
```

**6. Grid Layouts**
```css
/* Features/Stats Grid */
.features-grid, .stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-lg);
}

/* Results Grid */
.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-lg);
}
```

---

**JavaScript Utilities** - `src/server/static/js/main.js`

**Features:**
- Theme toggle (dark/light with localStorage persistence)
- Toast notifications
- Modal dialogs
- API helper functions
- Mobile menu toggle
- Keyboard shortcuts (Escape to close modal)
- Utility formatters (distance, coordinates, etc.)

**Example Functions:**
```javascript
// Toast Notifications
function showToast(message, type = 'info', duration = 3000) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), duration);
}

// Modal Dialogs
function showModal(title, content) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-backdrop" onclick="closeModal()"></div>
        <div class="modal-content">
            <div class="modal-header">
                <h2>${title}</h2>
                <button class="modal-close" onclick="closeModal()">×</button>
            </div>
            <div class="modal-body">${content}</div>
        </div>
    `;
    document.getElementById('modal-container').appendChild(modal);
}

// API Helpers
async function apiGet(endpoint) {
    const response = await fetch(endpoint);
    const data = await response.json();
    if (!data.success) {
        throw new Error(data.error?.message || 'API request failed');
    }
    return data.data;
}

// Theme Toggle
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}
```

---

**HTML Templates** - `src/server/templates/`

**Base Template** (`base.html`):
```html
<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.ProjectName}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="icon" href="/static/images/favicon.png">
</head>
<body>
    <header id="main-header">
        <div class="header-container">
            <div class="header-left">
                <button class="mobile-menu-toggle" onclick="toggleMobileMenu()">☰</button>
                <a class="logo" href="/">{{.ProjectName}}</a>
            </div>
            <nav id="main-nav" class="header-center">
                <a href="/">Home</a>
                <a href="/search">Search</a>
                <a href="/docs">API</a>
                <a href="/stats">Stats</a>
            </nav>
            <div class="header-right">
                <button class="theme-toggle" onclick="toggleTheme()">
                    <span class="theme-icon">🌙</span>
                </button>
            </div>
        </div>
    </header>
    
    <main id="main-content">
        {{template "content" .}}
    </main>
    
    <footer id="main-footer">
        <div class="footer-container">
            <p>&copy; 2024 {{.ProjectName}} - MIT License</p>
        </div>
    </footer>
    
    <div id="toast-container"></div>
    <div id="modal-container"></div>
    <script src="/static/js/main.js"></script>
</body>
</html>
```

**Template Embedding** (`templates.go`):
```go
package server

import (
    "embed"
    "html/template"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates *template.Template

func initTemplates() error {
    var err error
    templates, err = template.ParseFS(templateFS, "templates/*.html")
    return err
}
```

---

**Admin Interface Standards:**

**Admin Dashboard** (`/admin`):
- Clean, modern layout
- Server statistics cards
- Settings management
- Database status
- Log viewer
- System health monitoring

**Features:**
- ✅ Dark theme by default (light theme available)
- ✅ Responsive design (mobile-friendly)
- ✅ Toast notifications for feedback
- ✅ Modal dialogs for confirmations
- ✅ Real-time updates
- ✅ Form validation
- ✅ Keyboard shortcuts
- ✅ Accessible (ARIA labels)

**Color Scheme:**
- Dark Mode: `#1a1a1a` background, `#0066cc` primary
- Light Mode: `#ffffff` background, `#0066cc` primary
- Consistent with professional, modern UI/UX

---

**Responsive Design:**

```css
/* Mobile (< 720px) */
@media (max-width: 720px) {
  .mobile-menu-toggle { display: block; }
  .header-center {
    position: absolute;
    top: 70px;
    flex-direction: column;
    display: none;
  }
  .header-center.active { display: flex; }
  .hero h1 { font-size: 2rem; }
  .form-row { grid-template-columns: 1fr; }
}

/* Desktop (> 720px) */
@media (min-width: 720px) {
  #main-content { width: 90%; }
}
```

---

**Key Principles:**
- ✅ **No external dependencies** - No Bootstrap, Tailwind, React, etc.
- ✅ **Vanilla JS** - Pure JavaScript, no jQuery
- ✅ **Embedded assets** - All CSS/JS/images built into binary
- ✅ **Progressive enhancement** - Works without JavaScript
- ✅ **Dark theme first** - Default dark, supports light
- ✅ **Mobile responsive** - Works on all screen sizes
- ✅ **Fast loading** - Minimal CSS/JS, no external requests
- ✅ **Accessible** - Semantic HTML, ARIA labels

---

**File Sizes:**
- CSS: ~867 lines / ~25KB
- JavaScript: ~130 lines / ~4KB
- Total frontend: ~30KB (gzipped: ~10KB)

**Browser Support:**
- Modern browsers with CSS Grid support
- Chrome/Edge 88+
- Firefox 78+
- Safari 14+


---

### 13. **Testing & Debugging Environment Priority** ⭐ NEW

**Critical Rule**: ALWAYS use `/tmp/` for all temporary files, test data, and debugging. NEVER write test data to production directories.


**Rule**: Always use isolated environments for testing/debugging. Never test directly on the host OS.

**Priority Order:**
1. **Docker** (if available) - Preferred
2. **Incus** (if available) - Alternative containerization
3. **Host OS** - Last resort only

---

**Testing Workflow:**

**1. Docker (Preferred)**

```bash
# Build development image
make docker-dev

# Run with docker-compose (test configuration)
docker-compose -f docker-compose.test.yml up -d

# Access service
curl http://localhost:64181/healthz

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Cleanup
docker-compose -f docker-compose.test.yml down
sudo rm -rf /tmp/{projectname}/rootfs
```

**Benefits:**
- ✅ Isolated environment
- ✅ Consistent across all developers
- ✅ Easy cleanup (ephemeral /tmp storage)
- ✅ No pollution of host system
- ✅ Matches production environment

---

**2. Incus (Alternative)**

If Docker is not available, use Incus containers:

```bash
# Launch container
incus launch images:alpine/3.19 {projectname}-test

# Copy binary
incus file push ./binaries/{projectname} {projectname}-test/usr/local/bin/

# Execute
incus exec {projectname}-test -- /usr/local/bin/{projectname} --port 8080

# Shell access
incus shell {projectname}-test

# Cleanup
incus delete -f {projectname}-test
```

**Benefits:**
- ✅ System container (not Docker)
- ✅ LXD/LXC based
- ✅ Full OS environment
- ✅ Easy snapshot/restore

---

**3. Host OS (Last Resort)**

Only use when containers are unavailable:

```bash
# Build for host
make build

# Run directly
./binaries/{projectname} --dev --port 8080

# With custom directories (avoid polluting system)
./binaries/{projectname} \
  --config /tmp/{projectname}/config \
  --data /tmp/{projectname}/data \
  --logs /tmp/{projectname}/logs \
  --port 8080

# Cleanup
rm -rf /tmp/{projectname}
```

**⚠️ Warnings:**
- May pollute host system directories
- Different from production environment
- Harder to cleanup
- Not reproducible across systems

---

**Make Targets for Testing:**

```makefile
# Run tests in Docker (always)
test:
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:alpine \
		sh -c 'go test -v -race -timeout 5m ./...'

# Build development Docker image
docker-dev:
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		-t $(PROJECTNAME):dev \
		.

# Test with docker-compose
test-docker:
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for service..."
	@timeout 30 bash -c 'until curl -sf http://localhost:64181/healthz; do sleep 1; done'
	@echo "✓ Service is running"
	@docker-compose -f docker-compose.test.yml logs
	@docker-compose -f docker-compose.test.yml down
	@sudo rm -rf /tmp/{projectname}/rootfs

# Test with Incus (if Docker not available)
test-incus:
	@incus launch images:alpine/3.19 $(PROJECTNAME)-test
	@incus file push ./binaries/$(PROJECTNAME) $(PROJECTNAME)-test/usr/local/bin/
	@incus exec $(PROJECTNAME)-test -- /usr/local/bin/$(PROJECTNAME) --status
	@incus delete -f $(PROJECTNAME)-test
```

---

**CI/CD Testing:**

**GitHub Actions:**
- Always runs in Docker containers (ubuntu-latest)
- Uses `make test` which runs tests in Docker

**Jenkins:**
- Runs on dedicated agents (amd64, arm64)
- Each agent runs in isolated environment
- Uses Docker for builds

**Local Development:**
- Developers MUST use Docker or Incus
- Host OS testing only when troubleshooting container issues
- Document when host testing is used (unusual case)

---

**Best Practices:**

✅ **DO:**
- Always test in Docker first
- Use docker-compose.test.yml for integration tests
- Use ephemeral storage (/tmp) for test data
- Clean up after tests
- Document test environment used

❌ **DON'T:**
- Test directly on host OS without containers
- Leave test data in system directories
- Mix development and production configs
- Skip cleanup after tests
- Assume host environment matches production


### Testing Environment Priority
1. **Docker** (preferred) - `docker-compose -f docker-compose.test.yml up -d`
2. **Incus** (alternative) - `incus launch images:alpine/3.19 {project}-test`
3. **Host OS** (last resort) - `./binaries/{projectname} --dev`

Always use isolated environments. Never test directly on host OS unless containers unavailable.


---

### 14. **AI Assistant Guidelines** ⭐ NEW

**Critical Rules**:
- ✅ ALWAYS use `/tmp/` for all temporary files and test data
- ✅ NEVER write to production directories (/etc, /var/lib, /var/log)
- ✅ ALWAYS use random ports (64000-64999)
- ✅ NEVER use common ports (80, 443, 8080, 3000, 5000)
- ✅ ALWAYS use Docker or Incus for testing
- ✅ NEVER test on host OS unless explicitly requested


**When helping with development, testing, or debugging:**

**Environment Rules:**
1. **ALWAYS use Docker first** (if available)
2. **Use Incus second** (if Docker unavailable)
3. **NEVER test directly on host OS** unless explicitly requested

**Port Selection Rules:**
- ✅ **ALWAYS use random ports** for testing
- ✅ **Use port range**: 64000-64999 (random selection)
- ❌ **NEVER use**: 80, 443, 8080, 3000, 5000, or other common ports
- ✅ **Generate random**: `shuf -i 64000-64999 -n 1` or similar

**Testing Commands:**

```bash
# CORRECT - Docker with random port
TESTPORT=$(shuf -i 64000-64999 -n 1)
docker run -d --name test-container -p ${TESTPORT}:80 projectname:dev
curl http://localhost:${TESTPORT}/healthz

# CORRECT - Incus with random port
TESTPORT=$(shuf -i 64000-64999 -n 1)
incus launch images:alpine/3.19 test-container
incus exec test-container -- /usr/local/bin/projectname --port ${TESTPORT}

# WRONG - Direct host testing with common port
./binaries/projectname --port 8080  # ❌ DON'T DO THIS

# CORRECT - Host only if necessary, with random port
TESTPORT=$(shuf -i 64000-64999 -n 1)
./binaries/projectname --port ${TESTPORT}
```

**Example Testing Session:**

```bash
# 1. Build development image
make docker-dev

# 2. Generate random port
TESTPORT=$(shuf -i 64000-64999 -n 1)
echo "Using test port: ${TESTPORT}"

# 3. Run in Docker
docker run -d \
  --name projectname-test-${TESTPORT} \
  -p ${TESTPORT}:80 \
  -v /tmp/projectname-test:/data \
  -e ADMIN_PASSWORD=testpass123 \
  projectname:dev

# 4. Wait for startup
sleep 3

# 5. Test
curl http://localhost:${TESTPORT}/healthz
curl http://localhost:${TESTPORT}/api/v1

# 6. Cleanup
docker stop projectname-test-${TESTPORT}
docker rm projectname-test-${TESTPORT}
rm -rf /tmp/projectname-test
```

**Why These Rules:**

✅ **Docker/Incus first:**
- Isolated from host system
- Consistent environment
- Easy cleanup
- No system pollution
- Matches production

✅ **Random ports:**
- Avoids port conflicts
- No interference with running services
- Safe for multi-project development
- No accidental production port usage

❌ **Never use common ports when testing:**
- Port 80/443 - May conflict with web servers
- Port 8080 - Common development port (conflicts)
- Port 3000 - Node.js default (conflicts)
- Port 5000 - Flask/Python default (conflicts)

**AI Code Generation:**

When generating test scripts or examples, ALWAYS:
- Use Docker or docker-compose
- Generate random port: `$(shuf -i 64000-64999 -n 1)`
- Include cleanup steps
- Show ephemeral storage (/tmp)
- Never hardcode common ports

**Example AI-Generated Test Script:**

```bash
#!/bin/bash
# AI-generated test script following SPEC guidelines

set -e

PROJECTNAME="projectname"
TESTPORT=$(shuf -i 64000-64999 -n 1)

echo "🧪 Testing ${PROJECTNAME} using Docker"
echo "📡 Port: ${TESTPORT}"

# Build dev image
make docker-dev

# Run container
docker run -d \
  --name ${PROJECTNAME}-test-${TESTPORT} \
  -p ${TESTPORT}:80 \
  -v /tmp/${PROJECTNAME}-test:/data \
  ${PROJECTNAME}:dev

# Wait and test
sleep 3
curl http://localhost:${TESTPORT}/healthz || exit 1

# Cleanup
docker stop ${PROJECTNAME}-test-${TESTPORT}
docker rm ${PROJECTNAME}-test-${TESTPORT}
rm -rf /tmp/${PROJECTNAME}-test

echo "✅ Tests passed"
```


### Temporary Files & Testing
- ✅ ALWAYS use `/tmp/{projectname}/` for all test data
- ✅ NEVER use production directories (/etc, /var/lib, /var/log) for testing
- ✅ Cleanup after tests: `rm -rf /tmp/{projectname}`

### Port Selection for Testing
- ✅ ALWAYS random: `$(shuf -i 64000-64999 -n 1)`
- ❌ NEVER: 80, 443, 8080, 3000, 5000, or other common ports


### Version Format Standards
**release.txt**: `0.0.1` (no "v" prefix)
**Git tags**: `0.0.1` (no "v" prefix)
**GitHub releases**: `0.0.1` (no "v" prefix)
**Docker tags**: `ghcr.io/{org}/{project}:0.0.1` (no "v" prefix)
**Docker latest**: `ghcr.io/{org}/{project}:latest`
**CLI output** (`--version`): `0.0.1` (ONLY the version number)
**CLI output** (`--help`): Shows version without "v" prefix

**Rule**: "v" prefix is NOT used - all tags use plain version numbers. Everywhere else uses plain version number.


### Version Management Workflow

**release.txt** stores the current version for releases.

**Workflow:**
1. `make build` - Reads version from `release.txt`, does NOT modify it
2. Developer manually edits `release.txt` when ready for new version
3. `make release` - Creates GitHub release with current version
4. AFTER successful `gh release create`, auto-increments `release.txt`

**Example:**
```bash
# release.txt contains: 1.0.0
make build              # Builds version 1.0.0
make release            # Releases 1.0.0 to GitHub
                        # After success, release.txt → 1.0.1

make build              # Next build will be 1.0.1
```

**Version Override:**
```bash
VERSION=2.0.0 make build    # Build specific version
VERSION=2.0.0 make release  # Release specific version
```

**Why This Approach:**
- ✅ Version only increments on successful release
- ✅ Failed releases don't increment version
- ✅ Manual control over version numbers
- ✅ Clear workflow: build → test → release → auto-increment


---

**API Documentation Features** ⭐ REQUIRED

**Swagger/OpenAPI UI:**

Routes:
- `/openapi` - Swagger UI (web interface)
- `/api/v1/openapi` - Swagger UI (API version)
- `/api/v1/openapi.json` - OpenAPI 3.0 specification

Implementation uses Swagger UI CDN (unpkg.com):
```html
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
<script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
```

**GraphQL Playground:**

Routes:
- `/graphql` - GraphQL Playground (web interface)
- `/api/v1/graphql` - GraphQL endpoint (GET for playground, POST for queries)

Implementation uses GraphQL Playground CDN:
```html
<link rel="stylesheet" href="https://unpkg.com/graphql-playground-react@1.7.28/build/static/css/index.css">
<script src="https://unpkg.com/graphql-playground-react@1.7.28/build/GraphQLPlayground.js"></script>
```

**Route Matching Philosophy:**

Frontend routes MUST mirror API routes:
- `/openapi` ↔ `/api/v1/openapi`
- `/graphql` ↔ `/api/v1/graphql`
- `/search` ↔ `/api/v1/search`
- `/{resource}` ↔ `/api/v1/{resource}`

This makes the API predictable and consistent.

**Why Required:**
- All projects must include interactive API documentation
- Adds external CDN dependencies (or increase bundle size if self-hosted)
- GraphQL is essential for simple REST APIs

**When to Include:**
- ✅ Public APIs with many endpoints
- ✅ APIs with complex query parameters
- ✅ Projects with external integrations
- ✅ Developer-focused tools


---

### 15. **IPv6 Support** ⭐ REQUIRED

**Purpose**: Full dual-stack IPv4/IPv6 support for modern network environments

**Requirements:**
- ✅ Listen on both IPv4 and IPv6
- ✅ Accept connections from both protocols
- ✅ Handle IPv6 addresses in URLs
- ✅ GeoIP lookups for IPv6
- ✅ Proper IPv6 address formatting

---

#### **Listening Addresses**

**Dual-Stack (Recommended):**
```bash
--address ::              # Listen on all IPv6 (includes IPv4 via dual-stack)
--address 0.0.0.0         # Listen on all IPv4 only
```

**Specific Addresses:**
```bash
--address ::1             # IPv6 localhost
--address 127.0.0.1       # IPv4 localhost
--address 2001:db8::1     # Specific IPv6 address
--address 192.168.1.100   # Specific IPv4 address
```

**Environment Variables:**
```bash
ADDRESS=::                # Dual-stack (IPv4 + IPv6)
ADDRESS=0.0.0.0           # IPv4 only
ADDRESS=::1               # IPv6 localhost only
```

---

#### **URL Display with IPv6**

**IPv6 addresses must be enclosed in brackets:**

```go
func getAccessibleURL(port string) string {
    // ... detection logic ...
    
    // For IPv6, wrap in brackets
    if strings.Contains(ip, ":") {
        return fmt.Sprintf("http://[%s]:%s", ip, port)
    }
    return fmt.Sprintf("http://%s:%s", ip, port)
}
```

**Examples:**
- IPv4: `http://192.168.1.100:64555`
- IPv6: `http://[2001:db8::1]:64555`
- IPv6 localhost: `http://[::1]:64555`
- Hostname: `http://server.example.com:64555`

---

#### **Docker IPv6 Support**

**Enable IPv6 in Docker daemon** (`/etc/docker/daemon.json`):
```json
{
  "ipv6": true,
  "fixed-cidr-v6": "2001:db8:1::/64"
}
```

**docker-compose.yml:**
```yaml
services:
  {projectname}:
    image: ghcr.io/{organization}/{projectname}:latest
    environment:
      - ADDRESS=::              # Dual-stack
    networks:
      - {projectname}

networks:
  {projectname}:
    name: {projectname}
    enable_ipv6: true
    ipam:
      config:
        - subnet: 172.18.0.0/16
        - subnet: 2001:db8:1::/64
```

**Port Mappings:**
```yaml
# IPv4 only
ports:
  - "172.17.0.1:64180:80"

# Dual-stack (binds to both)
ports:
  - "64180:80"

# IPv6 explicit
ports:
  - "[::1]:64180:80"
```

---

#### **GeoIP with IPv6**

**MaxMind GeoLite2 databases support IPv6:**

```go
// Lookup IPv6 address
record, err := geoipReader.City(net.ParseIP("2001:4860:4860::8888"))

// Works with both IPv4 and IPv6
func lookupIP(ipStr string) (*geoip2.City, error) {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return nil, fmt.Errorf("invalid IP address")
    }
    
    // Automatically handles IPv4 or IPv6
    return geoipReader.City(ip)
}
```

**API Endpoints:**
```bash
# IPv4
curl http://localhost:64180/api/v1/geoip/8.8.8.8

# IPv6
curl http://localhost:64180/api/v1/geoip/2001:4860:4860::8888

# IPv6 (URL encoded)
curl "http://localhost:64180/api/v1/geoip/2001%3A4860%3A4860%3A%3A8888"
```

---

#### **Code Implementation**

**Listen on Dual-Stack:**

```go
func StartServer(address, port string) error {
    // Default to dual-stack if not specified
    if address == "" {
        address = "::"  // IPv6 dual-stack
    }
    
    addr := fmt.Sprintf("%s:%s", address, port)
    
    // Create listener
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("failed to listen on %s: %w", addr, err)
    }
    
    log.Printf("Server listening on %s", addr)
    return http.Serve(listener, handler)
}
```

**IPv6 Address Detection:**

```go
func getOutboundIP() string {
    // Try IPv4 first
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err == nil {
        defer conn.Close()
        localAddr := conn.LocalAddr().(*net.UDPAddr)
        return localAddr.IP.String()
    }
    
    // Try IPv6
    conn, err = net.Dial("udp", "[2001:4860:4860::8888]:80")
    if err == nil {
        defer conn.Close()
        localAddr := conn.LocalAddr().(*net.UDPAddr)
        return localAddr.IP.String()
    }
    
    return ""
}

func formatURLWithIP(ip, port string) string {
    // IPv6 addresses need brackets
    if strings.Contains(ip, ":") {
        return fmt.Sprintf("http://[%s]:%s", ip, port)
    }
    return fmt.Sprintf("http://%s:%s", ip, port)
}
```

**IP Parsing and Validation:**

```go
func parseIP(ipStr string) (net.IP, error) {
    // Remove brackets if present (from URL)
    ipStr = strings.Trim(ipStr, "[]")
    
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return nil, fmt.Errorf("invalid IP address: %s", ipStr)
    }
    
    return ip, nil
}

func isIPv6(ip net.IP) bool {
    return ip.To4() == nil
}
```

---

#### **Testing IPv6**

**Local Testing:**
```bash
# Start server on IPv6 localhost
./binaries/{projectname} --address ::1 --port 64555

# Test with curl
curl http://[::1]:64555/healthz
curl -g http://[::1]:64555/api/v1/airports/JFK

# Test with IPv6 GeoIP lookup
curl http://[::1]:64555/api/v1/geoip/2001:4860:4860::8888
```

**Docker Testing:**
```bash
# Enable IPv6 in test compose
docker-compose -f docker-compose.test.yml up -d

# Test both protocols
curl http://localhost:64181/healthz              # IPv4
curl http://[::1]:64181/healthz                  # IPv6
```

---

#### **Health Check with IPv6**

```go
func healthCheck(address, port string) error {
    // Detect if IPv6
    url := fmt.Sprintf("http://%s:%s/healthz", address, port)
    if strings.Contains(address, ":") && address != "::1" {
        url = fmt.Sprintf("http://[%s]:%s/healthz", address, port)
    }
    
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return fmt.Errorf("unhealthy: status %d", resp.StatusCode)
    }
    
    return nil
}
```

---

#### **Best Practices**

✅ **DO:**
- Default to `::` for dual-stack support
- Wrap IPv6 addresses in brackets for URLs
- Support both IPv4 and IPv6 in all endpoints
- Test with both protocols
- Use `net.ParseIP()` for validation
- Handle IPv6 in logs and error messages

✅ **Documentation:**
- Document IPv6 support in README
- Show IPv6 examples in API docs
- Include IPv6 in curl examples
- Document dual-stack configuration

❌ **DON'T:**
- Force IPv4-only mode by default
- Forget brackets around IPv6 in URLs
- Assume all IPs are IPv4
- Skip IPv6 testing

---

#### **Configuration Examples**

**systemd service:**
```ini
[Service]
Environment="ADDRESS=::"      # Dual-stack
```

**Docker Compose:**
```yaml
environment:
  - ADDRESS=::                # Dual-stack
  - PORT=80
```

**Manual:**
```bash
# IPv4 only
./binaries/{projectname} --address 0.0.0.0 --port 8080

# IPv6 only
./binaries/{projectname} --address :: --port 8080

# IPv6 localhost
./binaries/{projectname} --address ::1 --port 8080

# Dual-stack (recommended)
./binaries/{projectname} --address :: --port 8080
```

---

**IPv6 Support Summary:**
- ✅ Dual-stack by default (IPv4 + IPv6)
- ✅ Proper URL formatting with brackets
- ✅ GeoIP lookups for both protocols
- ✅ Docker networking with IPv6
- ✅ Testing on both protocols
- ✅ Documentation with examples


---

### 16. **GeoIP Databases (sapics/ip-location-db)** ⭐ REQUIRED

**Purpose**: IP geolocation using sapics/ip-location-db aggregated data sources

**Repository**: https://github.com/sapics/ip-location-db
**CDN**: jsdelivr (https://cdn.jsdelivr.net/npm/@ip-location-db/)
**Update Frequency**: Daily (automatically via CDN)

---

#### **Database Configuration**

**4 Required Databases:**

1. **geolite2-city-ipv4.mmdb** (~50MB)
   - City-level geolocation for IPv4 addresses
   - Coordinates, timezone, postal codes
   - MaxMind GeoLite2 data
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb`

2. **geolite2-city-ipv6.mmdb** (~40MB)
   - City-level geolocation for IPv6 addresses
   - Coordinates, timezone, postal codes
   - MaxMind GeoLite2 data
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb`

3. **geo-whois-asn-country.mmdb** (~8MB)
   - Country-level data (combined IPv4/IPv6)
   - Aggregated from WHOIS and ASN sources
   - **Public domain** (no attribution required)
   - Daily updates
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb`

4. **asn.mmdb** (~5MB)
   - ASN/ISP information (combined IPv4/IPv6)
   - Autonomous System Numbers
   - Daily updates
   - URL: `https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb`

**Total Size**: ~103MB (4 databases)

---

#### **Implementation**

**Service Structure:**

```go
package geoip

import (
    "net"
    "github.com/oschwald/geoip2-golang"
)

const (
    cityIPv4URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb"
    cityIPv6URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb"
    countryURL   = "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb"
    asnURL       = "https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb"
)

type Service struct {
    cityIPv4DB *geoip2.Reader  // City database for IPv4
    cityIPv6DB *geoip2.Reader  // City database for IPv6
    countryDB  *geoip2.Reader  // Country (combined IPv4/IPv6)
    asnDB      *geoip2.Reader  // ASN (combined IPv4/IPv6)
    dataDir    string
}
```

**IPv4/IPv6 Detection:**

```go
func (s *Service) Lookup(ip net.IP) (*GeoLocation, error) {
    // Determine which city database to use
    var cityDB *geoip2.Reader
    if ip.To4() != nil {
        // IPv4 address - use IPv4 city database
        cityDB = s.cityIPv4DB
    } else {
        // IPv6 address - use IPv6 city database
        cityDB = s.cityIPv6DB
    }
    
    // Try city lookup first
    city, err := cityDB.City(ip)
    if err == nil {
        return buildLocation(city, ip), nil
    }
    
    // Fallback to country database (works for both IPv4/IPv6)
    country, err := s.countryDB.Country(ip)
    if err != nil {
        return nil, err
    }
    
    return buildCountryLocation(country, ip), nil
}
```

**Download Function:**

```go
func (s *Service) DownloadDatabases() error {
    databases := map[string]string{
        "geolite2-city-ipv4.mmdb":    cityIPv4URL,
        "geolite2-city-ipv6.mmdb":    cityIPv6URL,
        "geo-whois-asn-country.mmdb": countryURL,
        "asn.mmdb":                   asnURL,
    }
    
    for filename, url := range databases {
        path := filepath.Join(s.dataDir, filename)
        if err := downloadFile(path, url); err != nil {
            return fmt.Errorf("failed to download %s: %w", filename, err)
        }
    }
    
    return s.LoadDatabases()
}
```

---

#### **Update Schedule**

**Automatic Updates (via Scheduler):**

```go
// In main.go
sched := scheduler.New()
sched.AddTask("geoip-update", "0 3 * * 0", func() error {
    return geoipSvc.UpdateDatabases()
})
sched.Start()
```

**Update Frequency:**
- **geolite2-city**: Twice weekly (MaxMind release schedule)
- **geo-whois-asn-country**: Daily (aggregated sources)
- **asn**: Daily

**Weekly scheduler** catches all updates (runs Sunday 3:00 AM).

---

#### **Benefits of sapics/ip-location-db**

✅ **Daily Updates**
- Country and ASN data updated daily
- City data updated twice weekly
- Vs weekly from P3TERX

✅ **Multiple Data Sources**
- Aggregates MaxMind GeoLite2, WHOIS, ASN, GeoFeed
- Higher accuracy through data fusion
- More comprehensive IP coverage

✅ **Public Domain Options**
- geo-whois-asn-country: Public domain (CC0/PDDL)
- No attribution required for country data
- Vs CC BY-SA 4.0 from MaxMind

✅ **CDN Delivery**
- jsdelivr CDN (fast, global distribution)
- Automatic caching and compression
- High availability (99.9% uptime)

✅ **IPv4/IPv6 Optimization**
- Separate city databases reduce memory usage
- Faster lookups (smaller files)
- Better performance for IPv6-heavy traffic

✅ **File Sizes**
- Separate databases are smaller individually
- Can download only what's needed
- Total: ~103MB for all 4

---

#### **Licensing**

**sapics/ip-location-db:**

| Database | License | Attribution Required |
|----------|---------|---------------------|
| geolite2-city-ipv4.mmdb | CC BY-SA 4.0 | ✅ MaxMind |
| geolite2-city-ipv6.mmdb | CC BY-SA 4.0 | ✅ MaxMind |
| geo-whois-asn-country.mmdb | CC0/PDDL | ❌ Public domain |
| asn.mmdb | Various | ⚠️ Check sources |

**Attribution:**

```
GeoIP data from:
- MaxMind GeoLite2 (https://www.maxmind.com/)
- Aggregated via sapics/ip-location-db (https://github.com/sapics/ip-location-db)
- Country data: Public domain (WHOIS aggregation)
```

---

#### **Storage Location**

**Databases stored in:**
- `{CONFIG_DIR}/geoip/geolite2-city-ipv4.mmdb`
- `{CONFIG_DIR}/geoip/geolite2-city-ipv6.mmdb`
- `{CONFIG_DIR}/geoip/geo-whois-asn-country.mmdb`
- `{CONFIG_DIR}/geoip/asn.mmdb`

**Example paths:**
- Linux (root): `/etc/{projectname}/geoip/*.mmdb`
- Linux (user): `~/.config/{projectname}/geoip/*.mmdb`
- macOS (user): `~/Library/Application Support/{ProjectName}/geoip/*.mmdb`
- Windows: `%APPDATA%\{ProjectName}\geoip\*.mmdb`
- Docker: `/config/geoip/*.mmdb`

---

#### **First Run Behavior**

1. Check if databases exist in `{CONFIG_DIR}/geoip/`
2. If not found, download all 4 databases from jsdelivr CDN
3. Display progress for each database
4. Load databases into memory
5. Ready to serve GeoIP lookups

**Console Output:**
```
Loading GeoIP databases...
GeoIP databases not found, downloading...
  Downloading geolite2-city-ipv4.mmdb...
  Downloading geolite2-city-ipv6.mmdb...
  Downloading geo-whois-asn-country.mmdb...
  Downloading asn.mmdb...
GeoIP databases loaded successfully
```

---

#### **API Endpoints**

**GeoIP Lookup Routes:**

```yaml
# Lookup request IP
GET /api/v1/geoip
Response: { "ip": "1.2.3.4", "country": "US", "city": "New York", ... }

# Lookup specific IPv4
GET /api/v1/geoip/8.8.8.8

# Lookup specific IPv6
GET /api/v1/geoip/2001:4860:4860::8888

# Find nearby resources (e.g., airports) near IP location
GET /api/v1/geoip/{resource}/nearby?radius=100&limit=10
```

**Response Format:**

```json
{
  "success": true,
  "data": {
    "ip": "8.8.8.8",
    "country": "US",
    "country_name": "United States",
    "region": "CA",
    "region_name": "California",
    "city": "Mountain View",
    "postal_code": "94035",
    "latitude": 37.386,
    "longitude": -122.0838,
    "timezone": "America/Los_Angeles",
    "asn": 15169,
    "organization": "Google LLC"
  }
}
```

---

#### **Testing GeoIP**

**IPv4 Testing:**
```bash
curl http://localhost:64555/api/v1/geoip/8.8.8.8
curl http://localhost:64555/api/v1/geoip/1.1.1.1
```

**IPv6 Testing:**
```bash
curl http://localhost:64555/api/v1/geoip/2001:4860:4860::8888  # Google IPv6
curl http://localhost:64555/api/v1/geoip/2606:4700:4700::1111  # Cloudflare IPv6
```

**Auto-detection (request IP):**
```bash
curl http://localhost:64555/api/v1/geoip
```

---

#### **Manual Database Updates**

**Delete and re-download:**
```bash
# Remove old databases
rm -rf {CONFIG_DIR}/geoip/*.mmdb

# Restart service (auto-downloads)
systemctl restart {projectname}
```

**Programmatic update:**
```bash
# Via API (admin auth required)
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:64555/api/v1/admin/geoip/update
```

---

**sapics/ip-location-db provides superior GeoIP data with daily updates, multiple sources, and public domain options!**

