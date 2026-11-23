# Timezones API

Timezones API is a public, read-only reference service for bundled world timezone records. It provides versioned REST endpoints, a server-rendered web interface, and a companion terminal client. Record offsets and daylight-saving flags describe the static dataset; the current-time endpoint uses IANA timezone rules at request time.

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

The production image listens on port 80. Its entrypoint enables Tor when available; set `ENABLE_TOR=false` to disable it.

```sh
docker run --rm -p 8080:80 -e ENABLE_TOR=false ghcr.io/apimgr/timezones:latest
```

## Client CLI

The CLI looks up a timezone by display name, abbreviation, or IANA identifier. Use `--json` for machine-readable output.

```sh
timezones-cli America/New_York
timezones-cli --json EST
timezones-cli --help
timezones-cli --version
```

## Server

Run the server in the foreground:

```sh
timezones --address 127.0.0.1 --port 8080
```

The server creates a default `server.yml` in its configuration directory. `CONFIG_DIR`, `PORT`, and `ADDRESS` override the corresponding configuration or listener setting. The browser includes timezone list, search, random, offset, abbreviation, identifier and name lookups at `/timezones`, current-time views at `/timezones/current/{iana-id}`, dataset summaries at `/stats` and `/count`, plus server information and documentation pages.

## Configuration

The server loads `server.yml` and creates a default file on first start. Current overrides are `CONFIG_DIR`, `PORT`, and `ADDRESS`. The theme preference and request-protection state use first-party cookies. Requests are size-limited and rate-limited by client address.

## API

All endpoints are public and read-only. The versioned REST API is rooted at `/api/v1`.

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
| `POST /api/v1/server/graphql` | Query timezone records, current time, or dataset statistics |

Add `.txt` to the list, search, random, statistics, and count endpoints for plain-text responses. Search queries are limited to 128 bytes. `/api/swagger` aliases the current OpenAPI document, and `/api/graphql` aliases the versioned GraphQL endpoint.

JSON responses use an `ok` and `data` envelope. List and search responses include page, limit, total, and page-count metadata. The default page size is 250.

```sh
curl -q -LSsf 'http://127.0.0.1:8080/api/v1/timezones/search?q=America%2FNew_York'
curl -q -LSsf 'http://127.0.0.1:8080/api/v1/timezones/current/America/New_York'
curl -q -LSsf -H 'Content-Type: application/json' -d '{"query":"{ timezone(identifier: \"America/New_York\") { value abbr utc } }"}' http://127.0.0.1:8080/api/v1/server/graphql
```

## Other

- Full [ReadTheDocs documentation](https://apimgr-timezones.readthedocs.io/en/latest/)
- [Source repository](https://github.com/apimgr/timezones)
- Privacy and terms pages are served by each deployment at `/server/privacy` and `/server/terms`.

## Development

The timezone dataset and web assets are embedded in the Go binaries. Go commands run in the Docker toolchain configured by the Makefile. `make dev` builds the current platform, `make build` builds the release platforms, and `make test` runs Go vet, vulnerability checks, and tests with coverage.

GitHub Actions build server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64. Docker images target amd64 and arm64.

## Disclaimer

This software is provided “as is” without warranty of any kind. Timezone data may not reflect current civil-time rules; verify time-critical decisions against authoritative sources. No software is guaranteed to be free of vulnerabilities. Evaluate it before production use.

## License

Project code is MIT licensed. The embedded timezone dataset has separate source and redistribution terms. Review [LICENSE.md](LICENSE.md) before redistribution.
