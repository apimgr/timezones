.PHONY: build test docker docker-dev clean release help

# Project configuration
PROJECTNAME := timezones
PROJECTORG := casjay
VERSION := $(shell cat release.txt 2>/dev/null || echo "0.0.1")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -w -s

# Platforms to build for
PLATFORMS := \
	linux-amd64 \
	linux-arm64 \
	windows-amd64 \
	windows-arm64 \
	darwin-amd64 \
	darwin-arm64 \
	freebsd-amd64 \
	freebsd-arm64

# Default target
all: build

# Help target
help:
	@echo "Timezones API - Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build binaries for all platforms"
	@echo "  test         - Run tests"
	@echo "  docker       - Build and push multi-platform Docker images"
	@echo "  docker-dev   - Build development Docker image (local only)"
	@echo "  clean        - Remove build artifacts"
	@echo "  release      - Create GitHub release with binaries"
	@echo ""
	@echo "Current version: $(VERSION)"

# Build binaries for all platforms
build:
	@echo "Building $(PROJECTNAME) v$(VERSION)..."
	@mkdir -p binaries releases
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%-*} GOARCH=$${platform#*-} ; \
		EXT="" ; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi ; \
		echo "  Building $$GOOS/$$GOARCH..." ; \
		CGO_ENABLED=1 GOOS=$$GOOS GOARCH=$$GOARCH go build \
			-ldflags "$(LDFLAGS)" \
			-o binaries/$(PROJECTNAME)-$$GOOS-$$GOARCH$$EXT \
			./src ; \
		cp binaries/$(PROJECTNAME)-$$GOOS-$$GOARCH$$EXT releases/ ; \
	done
	@echo "✓ Build complete: binaries/"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -timeout 5m ./...
	@echo "✓ Tests passed"

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

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf binaries/ releases/
	@echo "✓ Clean complete"

# Create GitHub release
release: build
	@echo "Creating GitHub release v$(VERSION)..."
	@if command -v gh >/dev/null 2>&1; then \
		gh release delete $(VERSION) -y 2>/dev/null || true ; \
		gh release create $(VERSION) \
			--title "$(VERSION)" \
			--notes "Release $(VERSION)" \
			./releases/* ; \
		echo "✓ GitHub release created: $(VERSION)" ; \
		NEW_VERSION=$$(echo $(VERSION) | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g') ; \
		echo $$NEW_VERSION > release.txt ; \
		echo "✓ Version incremented to $$NEW_VERSION" ; \
	else \
		echo "❌ gh CLI not found. Install from https://cli.github.com" ; \
		exit 1 ; \
	fi

# Build for host platform only (quick build for testing)
quick:
	@echo "Building for host platform..."
	@mkdir -p binaries
	@CGO_ENABLED=1 go build \
		-ldflags "$(LDFLAGS)" \
		-o binaries/$(PROJECTNAME) \
		./src
	@echo "✓ Binary: binaries/$(PROJECTNAME)"
