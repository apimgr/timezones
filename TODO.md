# Timezones API - SPEC.md Implementation Task List

**Last Updated**: 2025-10-16 08:15:00

## 🔄 In Progress (ONLY ONE TASK AT A TIME)
- [ ] Migrate from Gorilla Mux to Chi router
  - Update server package to use chi/v5
  - Apply security middleware (rate limiting, headers)
  - Maintain all existing endpoints

## 📋 Pending

### Core Infrastructure (High Priority)
- [ ] Add scheduler package (src/scheduler/)
  - Cron-like task scheduling
  - Background jobs support
  - No external cron dependency

- [ ] Add enhanced logging system
  - access.log (HTTP requests)
  - error.log (errors/warnings)
  - audit.log (admin actions)
  - Log rotation support
  - Configurable formats (Apache Combined, JSON)

- [ ] Add debug mode support
  - DEBUG environment variable
  - Debug endpoints (/debug/routes, /debug/config, /debug/db, /debug/memory)
  - Template hot reload
  - Verbose logging

### Admin & Configuration (High Priority)
- [ ] Add admin WebUI (/admin page)
  - Settings management interface
  - Live configuration reload
  - Server statistics dashboard
  - Database status view

- [ ] Expand admin API endpoints
  - GET /api/v1/admin/settings (list all)
  - POST /api/v1/admin/settings (update)
  - DELETE /api/v1/admin/settings/{key}
  - GET /api/v1/admin/stats
  - GET /api/v1/admin/health

### Service Management (High Priority)
- [ ] Add service commands
  - service start/stop/restart/reload/status
  - Signal handling (SIGTERM, SIGHUP, SIGUSR1, SIGUSR2)
  - Graceful shutdown
  - PID file management

- [ ] Add maintenance mode
  - Read-only mode on database failure
  - Auto-recovery attempts
  - Admin UI diagnostics
  - Self-healing logic

### Installation & Deployment (Medium Priority)
- [ ] Create installation scripts
  - scripts/install-linux.sh (systemd/OpenRC/SysVinit/runit auto-detect)
  - scripts/install-macos.sh (launchd)
  - scripts/install-bsd.sh (rc.d)
  - scripts/install-windows.ps1 (NSSM)
  - scripts/uninstall.sh (universal)

### Documentation (Medium Priority)
- [ ] Create docs/ directory with ReadTheDocs
  - docs/index.md (documentation home)
  - docs/API.md (complete API reference)
  - docs/SERVER.md (server administration)
  - docs/mkdocs.yml (Dracula theme)
  - docs/requirements.txt
  - docs/stylesheets/dracula.css
  - .readthedocs.yml

### API Documentation UI (Medium Priority)
- [ ] Add Swagger/OpenAPI UI
  - /openapi route (Swagger UI)
  - /api/v1/openapi route
  - /api/v1/openapi.json (OpenAPI 3.0 spec)

- [ ] Add GraphQL Playground
  - /graphql route (playground UI)
  - /api/v1/graphql endpoint

### Optional Features
- [ ] HTTPS/TLS & Let's Encrypt support
  - Automatic certificate obtainment
  - DNS-01, TLS-ALPN-01, HTTP-01 challenges
  - Auto-renewal
  - Admin UI for SSL configuration

- [ ] GeoIP Integration - **SKIP**
  - Timezones API has no location-based features
  - No "find nearby" or geographic queries
  - Per SPEC Exception 2, GeoIP not required

## ✅ Completed
- [x] Copy full SPEC.md to CLAUDE.md
  - Files: CLAUDE.md

- [x] Fix Dockerfile (CGO_ENABLED=0)
  - Changed from CGO_ENABLED=1 to CGO_ENABLED=0
  - Switched from github.com/mattn/go-sqlite3 to modernc.org/sqlite (pure Go)
  - Removed sqlite-libs from runtime
  - Files: Dockerfile, go.mod, src/database/database.go

- [x] Update Makefile to use Docker builds
  - All builds now use Docker (golang:alpine)
  - CGO_ENABLED=0 in all build commands
  - Added test-docker target
  - Files: Makefile

- [x] Add security package foundations
  - Created src/security/security.go
  - Rate limiting middleware
  - Security headers middleware
  - Brute force protection
  - Files: src/security/security.go, go.mod

---

**Progress**: 4/25 tasks completed (16%)
