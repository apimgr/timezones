PROJECT_NAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)(\.git)?$$|\1|' || basename "$(CURDIR)")
PROJECT_ORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+(\.git)?$$|\1|' || basename "$(dir $(CURDIR))")
INTERNAL_NAME := $(PROJECT_NAME)
VERSION := $(shell if [ -f release.txt ]; then cat release.txt; else printf '%s' "$${VERSION:-devel}"; fi)
COMMIT_ID := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo N/A)
BUILD_EPOCH := $(shell date -u +%s)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
OFFICIAL_SITE := $(shell if [ -f site.txt ]; then cat site.txt; else printf '%s' "$${OFFICIAL_SITE:-}"; fi)
BINDIR ?= binaries
RELDIR ?= releases
PLATFORMS ?= linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64,windows/arm64,freebsd/amd64,freebsd/arm64
REGISTRY ?= ghcr.io/$(PROJECT_ORG)/$(PROJECT_NAME)
GO_CACHE ?= $(HOME)/go/pkg/mod
GO_BUILD ?= $(HOME)/.cache/go-build/$(PROJECT_NAME)
DOCKER_MEM ?= 4g
DOCKER_CPUS ?= 2
HOST_GOOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOST_GOARCH := $(shell uname -m | sed -e "s/^x86_64$$/amd64/" -e "s/^aarch64$$/arm64/")

LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT_ID)' -X 'main.BuildDate=$(BUILD_DATE)'
GO_DOCKER_RUN := docker run --rm \
	--name $(PROJECT_NAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) \
	--memory=$(DOCKER_MEM) --cpus=$(DOCKER_CPUS) \
	-v "$(CURDIR):/app" \
	-v "$(GO_CACHE):/usr/local/share/go/pkg/mod" \
	-v "$(GO_BUILD):/usr/local/share/go/cache" \
	-w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false
GO_DOCKER := $(GO_DOCKER_RUN) casjaysdev/go:latest

.PHONY: dev local build test release docker

dev: ## Build the current platform into an isolated temporary directory
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)" "$${TMPDIR:-/tmp}/$(PROJECT_ORG)"
	@$(GO_DOCKER) go mod download
	@BUILD_DIR=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECT_ORG)/$(INTERNAL_NAME)-XXXXXX") && \
		$(GO_DOCKER_RUN) -v "$$BUILD_DIR:/build" casjaysdev/go:latest \
		sh -c 'GOOS="$(HOST_GOOS)" GOARCH="$(HOST_GOARCH)" go build -trimpath -o /build/$(PROJECT_NAME) ./src' && \
		echo "Built $$BUILD_DIR/$(PROJECT_NAME)"

local: ## Build a production binary for the host platform
	@rm -rf "$(BINDIR)" "$(RELDIR)"
	@mkdir -p "$(BINDIR)" "$(GO_CACHE)" "$(GO_BUILD)"
	@$(GO_DOCKER) go mod download
	@OUTPUT="$(BINDIR)/$(PROJECT_NAME)"; \
		[ "$(HOST_GOOS)" != windows ] || OUTPUT="$$OUTPUT.exe"; \
		$(GO_DOCKER) sh -c 'GOOS="$(HOST_GOOS)" GOARCH="$(HOST_GOARCH)" go build -trimpath -ldflags "$(LDFLAGS)" -o "'"$$OUTPUT"'" ./src'

build: ## Build the host binary and all eight release platforms
	@rm -rf "$(BINDIR)" "$(RELDIR)"
	@mkdir -p "$(BINDIR)" "$(GO_CACHE)" "$(GO_BUILD)"
	@$(GO_DOCKER) go mod download
	@$(GO_DOCKER) sh -c 'GOOS="$(HOST_GOOS)" GOARCH="$(HOST_GOARCH)" go build -trimpath -ldflags "$(LDFLAGS)" -o "$(BINDIR)/$(PROJECT_NAME)" ./src'
	@set -eu; for platform in $$(printf '%s' "$(PLATFORMS)" | tr ',' ' '); do \
		GOOS=$${platform%/*}; GOARCH=$${platform#*/}; OUTPUT="$(BINDIR)/$(PROJECT_NAME)-$$GOOS-$$GOARCH"; \
		[ "$$GOOS" != windows ] || OUTPUT="$$OUTPUT.exe"; \
		echo "Building $$GOOS/$$GOARCH"; \
		$(GO_DOCKER) sh -c "GOOS=$$GOOS GOARCH=$$GOARCH go build -trimpath -ldflags \"$(LDFLAGS)\" -o $$OUTPUT ./src"; \
	done

test: ## Run vet, vulnerability checks, and tests with coverage
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)"
	@$(GO_DOCKER) go mod download
	@$(GO_DOCKER) go vet ./...
	@$(GO_DOCKER) govulncheck ./...
	@$(GO_DOCKER) go test -v -cover ./...

release: build ## Prepare release binaries and source archive
	@mkdir -p "$(RELDIR)"
	@printf '%s\n' "$(VERSION)" > "$(RELDIR)/version.txt"
	@tar --exclude=.git --exclude=.github --exclude=.gitea --exclude=$(BINDIR) --exclude=$(RELDIR) -czf "$(RELDIR)/$(PROJECT_NAME)-$(VERSION)-source.tar.gz" .
	@sha256sum "$(BINDIR)"/$(PROJECT_NAME)-* "$(RELDIR)/$(PROJECT_NAME)-$(VERSION)-source.tar.gz" > "$(RELDIR)/sha256.txt"
	@sha512sum "$(BINDIR)"/$(PROJECT_NAME)-* "$(RELDIR)/$(PROJECT_NAME)-$(VERSION)-source.tar.gz" > "$(RELDIR)/sha512.txt"
	@RELEASE_TAG="$(VERSION)"; case "$(VERSION)" in [0-9]*.[0-9]*.[0-9]*) RELEASE_TAG="v$(VERSION)" ;; esac; \
		gh release create "$$RELEASE_TAG" "$(BINDIR)"/$(PROJECT_NAME)-* "$(RELDIR)"/* \
		--title "$(PROJECT_NAME) $$RELEASE_TAG" --notes "Release $(VERSION)" --latest

docker: ## Build and publish the multi-architecture container image
	@docker buildx build -f docker/Dockerfile --platform linux/amd64,linux/arm64 \
		--build-arg VERSION="$(VERSION)" --build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg BUILD_EPOCH="$(BUILD_EPOCH)" --build-arg COMMIT_ID="$(COMMIT_ID)" \
		-t "$(REGISTRY):$(VERSION)" -t "$(REGISTRY):latest" --push .
