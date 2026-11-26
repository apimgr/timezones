# Timezones API

Timezones API is a public reference service for bundled world timezone records. It provides public, read-only versioned REST and GraphQL endpoints, a server-rendered web interface, owner and linked bearer-token management, an authenticated MCP endpoint, and the `timezones-cli` companion client. Record offsets and daylight-saving flags describe the static dataset; the current-time endpoint uses IANA timezone rules at request time.

The embedded translation catalogs currently support English, Spanish, French, Arabic (right-to-left), German, and Italian. Browser requests can select a locale with `?lang=es`; otherwise the server checks a `lang` cookie and the request's `Accept-Language` header before falling back to English.

[![GitHub Actions](https://github.com/apimgr/timezones/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/timezones/actions/workflows/ci.yml) [![License](https://img.shields.io/github/license/apimgr/timezones)](LICENSE.md) [![Documentation](https://readthedocs.org/projects/apimgr-timezones/badge/?version=latest)](https://apimgr-timezones.readthedocs.io/en/latest/)

## Install

Download the latest server and CLI binaries from [GitHub Releases](https://github.com/apimgr/timezones/releases/latest). Release binaries are built for these platforms:

| Platform | Architecture | Server binary | Client binary |
|---|---|---|---|
| Linux | amd64 | `timezones-linux-amd64` | `timezones-cli-linux-amd64` |
| Linux | arm64 | `timezones-linux-arm64` | `timezones-cli-linux-arm64` |
| macOS | Intel (amd64) | `timezones-darwin-amd64` | `timezones-cli-darwin-amd64` |
| macOS | Apple Silicon (arm64) | `timezones-darwin-arm64` | `timezones-cli-darwin-arm64` |
| Windows | amd64 | `timezones-windows-amd64.exe` | `timezones-cli-windows-amd64.exe` |
| Windows | arm64 | `timezones-windows-arm64.exe` | `timezones-cli-windows-arm64.exe` |
| FreeBSD | amd64 | `timezones-freebsd-amd64` | `timezones-cli-freebsd-amd64` |
| FreeBSD | arm64 | `timezones-freebsd-arm64` | `timezones-cli-freebsd-arm64` |

For Linux amd64:

```sh
curl -q -LSsf https://github.com/apimgr/timezones/releases/latest/download/timezones-linux-amd64 -o /usr/local/bin/timezones && chmod +x /usr/local/bin/timezones
curl -q -LSsf https://github.com/apimgr/timezones/releases/latest/download/timezones-cli-linux-amd64 -o /usr/local/bin/timezones-cli && chmod +x /usr/local/bin/timezones-cli
```

On macOS, choose the `darwin-amd64` or `darwin-arm64` assets, install them in a directory on `PATH`, and remove quarantine if required by Gatekeeper:

```sh
xattr -d com.apple.quarantine /usr/local/bin/timezones /usr/local/bin/timezones-cli 2>/dev/null || true
```

On Windows, download the matching `.exe` files and add their directory to `%PATH%`.

## Docker

The production image listens on port 80 and includes Tor and I2P packages. When `tor` is available, the server starts an isolated Tor hidden service. I2P is enabled by default in application containers and can be enabled or disabled with `I2P_ENABLED=true|false` or `--i2p[=true|false]`. Both services use separate loopback listeners; their addresses are shown by health and discovery only when configured for exposure.

```sh
docker run --rm -p 8080:80 ghcr.io/apimgr/timezones:latest
```

## Client CLI

The CLI queries a configured Timezones API server by display name, abbreviation, or IANA identifier. It also provides list/search/current-time commands, interactive terminal and desktop interfaces, and `--output` formats for machine-readable output.

```sh
timezones-cli America/New_York
timezones-cli --output json EST
timezones-cli --server http://127.0.0.1:8080 list --page 1
timezones-cli --help
timezones-cli --version
```

## Server

Run the server in the foreground:

```sh
timezones --address 127.0.0.1 --port 8080
```

The server creates a default `server.yml` in its configuration directory. `CONFIG_DIR`, `PORT`, and `ADDRESS` override the corresponding configuration or listener setting. The browser includes timezone list, search, random, offset, abbreviation, identifier and name lookups at `/timezones`, current-time views at `/timezones/current/{iana-id}`, dataset summaries at `/stats` and `/count`, plus server information, preferences, security, contact, and documentation pages.

## Configuration

The server loads `server.yml` and its configured YAML drop-ins, creating the default configuration on first start. `CONFIG_DIR`, `DATA_DIR`, `PORT`, and `ADDRESS` select the configuration/data paths and listener settings. Owner credentials are managed through the token preferences page or API; the CLI accepts a token from its config, environment, `--token`, or `--token-file`. The theme preference and request-protection state use first-party cookies. Requests are size-limited and rate-limited by client address.

## API

Public timezone reads do not require a token. Owner-token and linked-token operations use the API routes under `/api/v1/server/preferences/tokens`; MCP requires a valid Bearer token. The versioned REST API is rooted at `/api/v1`.

| Method and path | Purpose |
|---|---|
| `GET /api/v1/timezones` | List records with `page`/`limit`; optional `offset` and `isdst` filters |
| `GET /api/v1/timezones.json` | Return the embedded JSON dataset |
| `GET /api/v1/timezones/search?q=...` | Search by name, description, abbreviation, or IANA identifier |
| `GET /api/v1/timezones/offset/{offset}` | Filter by UTC offset in hours |
| `GET /api/v1/timezones/abbr/{abbr}` | Look up an abbreviation |
| `GET /api/v1/timezones/utc/{iana-id}` | Look up an IANA identifier |
| `GET /api/v1/timezones/value/{name}` | Look up a display name |
| `GET /api/v1/timezones/random` | Return one random record |
| `GET /api/v1/timezones/current/{iana-id}` | Return the current local time and UTC offset |
| `GET /api/v1/stats` and `/api/v1/count` | Return dataset statistics or record count |
| `GET /api/v1/server/healthz` | Check service health |
| `GET /api/v1/server/swagger` | Retrieve the OpenAPI JSON description |
| `GET` and `POST /api/v1/server/graphql` | Query timezone records, current time, or dataset statistics |
| `GET /api/v1/server/preferences/tokens` | List the caller's token metadata |
| `POST /api/v1/server/preferences/tokens` | Create a linked token |
| `POST /api/v1/server/preferences/tokens/primary` | Create the caller's primary owner token |
| `GET /api/autodiscover` | Discover server API and published CLI release metadata |
| `GET /cli/binaries/{binary}` | Download a verified CLI release artifact when present |

Add `.txt` to the list, search, random, statistics, and count endpoints for plain-text responses. Search queries are limited to 128 bytes. `/api/swagger` aliases the current OpenAPI document, and `/api/graphql` aliases the versioned GraphQL endpoint. GraphQL supports POST JSON requests and GET requests with `query`, optional JSON-object `variables`, and optional `operationName` query parameters; a GET without `query` returns usage text.

JSON responses use an `ok` and `data` envelope. List and search responses include page, limit, total, and page-count metadata. The default page size is 250.

```sh
curl -q -LSsf 'http://127.0.0.1:8080/api/v1/timezones/search?q=America%2FNew_York'
curl -q -LSsf 'http://127.0.0.1:8080/api/v1/timezones/current/America/New_York'
curl -q -LSsf -H 'Content-Type: application/json' -d '{"query":"{ timezone(identifier: \"America/New_York\") { value abbr utc } }"}' http://127.0.0.1:8080/api/v1/server/graphql
curl -q -G -LSsf --data-urlencode 'query={ timezone(identifier: "America/New_York") { value abbr utc } }' http://127.0.0.1:8080/api/v1/server/graphql
```

## Other

- Full [ReadTheDocs documentation](https://apimgr-timezones.readthedocs.io/en/latest/)
- [Source repository](https://github.com/apimgr/timezones)
- `/server` redirects to the about page; privacy, terms, security, security policy, acknowledgments, and contact pages are served under `/server`.
- MCP is available at `/mcp`, `/api/mcp`, and `/api/v1/server/mcp`; all aliases require `Authorization: Bearer …`.
- Security reporting starts at `/server/security`; see the [security policy](https://apimgr-timezones.readthedocs.io/en/latest/security/) and [deployment security.txt](https://apimgr-timezones.readthedocs.io/en/latest/integrations/#public-endpoints).

## Development

The timezone dataset and web assets are embedded in the Go binaries. Go commands run in the Docker toolchain configured by the Makefile. `make dev` builds the current platform, `make build` builds the release platforms, and `make test` runs Go vet, vulnerability checks, tests, and coverage enforcement. Run `./tests/run_tests.sh` for Docker-based binary and HTTP integration checks, and `./tests/e2e.sh` for browser checks using isolated Docker containers, Chromium, and the Go toolchain image. Both require Docker and image/network access. Run `./tests/incus.sh [systemd|openrc|all]` for full service lifecycle checks in disposable Debian/systemd and Alpine/OpenRC Incus containers; it defaults to both modes and requires an initialized Incus host and Docker for `make build`.

GitHub Actions build server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64. Docker images target amd64 and arm64.

## Disclaimer

This software is provided “as is” without warranty of any kind. Timezone data may not reflect current civil-time rules; verify time-critical decisions against authoritative sources. No software is guaranteed to be free of vulnerabilities. Evaluate it before production use.

## License

Project code is MIT licensed. The embedded timezone dataset is copied from [dmfilipenko/timezones.json](https://github.com/dmfilipenko/timezones.json) at commit `658abf632bae9c0313360702e3018024eb88a874` and is separately attributed under its MIT license in [LICENSE.md](LICENSE.md).
