# Timezones API

[![CI](https://github.com/apimgr/timezones/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/timezones/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/release/apimgr/timezones)](https://github.com/apimgr/timezones/releases) [![License](https://img.shields.io/github/license/apimgr/timezones)](LICENSE.md) [![Documentation](https://readthedocs.org/projects/apimgr-timezones/badge/?version=latest)](https://apimgr-timezones.readthedocs.io/)

## About

Timezones API serves a bundled reference dataset of timezone display names,
abbreviations, UTC offsets, daylight-saving flags, labels, and IANA identifiers.
It provides versioned REST and read-only GraphQL APIs, a server-rendered web
interface, and the `timezones-cli` companion client. Dataset offsets and flags
are static; current local time is calculated from IANA timezone rules.

## Features

- Public timezone search, lookup, filtering, random selection, and dataset statistics
- Current-time lookup by IANA timezone identifier
- Server-rendered, mobile-friendly pages and an installable PWA shell
- Owner and linked bearer-token management without user accounts
- Read-only GraphQL API and query explorer
- Authenticated Streamable HTTP MCP endpoint
- Tor hidden service when a Tor binary is available; optional I2P eepsite
- Localized, sanitized operator branding for the shared header, About, Privacy, Terms, and footer
- Server and CLI interfaces in English, Spanish, French, German, Italian, and Arabic

## Production

### Docker Compose

The provided Compose setup runs the API and Valkey cache. It stores configuration
and data in `docker/volumes/` and binds the published port to Docker's bridge
gateway as configured in the compose file.

```sh
docker compose -f docker/docker-compose.yml up -d
```

The production image listens on port 80 and includes `tor` and `i2pd`. Tor is
managed automatically when its binary is available. I2P is enabled by default
in an application container; set `I2P_ENABLED=true|false` or use `--i2p` to
override the setting. Overlay addresses appear in discovery only when the
respective `server.tor.expose` or `server.i2p.expose` setting is enabled; the
same exposure flags control the health response, Help page, and shared footer.

### Binary

Download release assets from [GitHub Releases](https://github.com/apimgr/timezones/releases).
The server and CLI are available for Linux, macOS, Windows, and FreeBSD on
amd64 and arm64. Optional desktop CLI builds are also published for Linux,
macOS, and Windows.

For a local Linux amd64 installation:

```sh
curl -q -LSsf -o timezones https://github.com/apimgr/timezones/releases/latest/download/timezones-linux-amd64
chmod +x timezones
sudo install -m 0755 timezones /usr/local/bin/timezones
```

### systemd

After installing `timezones` on a Linux host managed by systemd, install,
enable, and start the service:

```sh
sudo timezones --service --install
```

## Client

The companion executable is `timezones-cli`. Install the matching release
asset, then set the server URL with `--server` or in its `cli.yml` profile.
This project does not define an official server URL; use your deployment's URL.

```sh
timezones-cli --server https://YOUR_SERVER search America/New_York
timezones-cli --server https://YOUR_SERVER list --page 1
timezones-cli --help
```

The client also supports `current`, `utc`, `offset`, `abbr`, `value`, `random`,
`stats`, `count`, token management, table/JSON/plain output, and an interactive
terminal interface. Optional `timezones-cli-gui-{os}-{arch}` builds provide a
desktop interface. The GUI toolkit does not yet provide platform accessibility
adapters for screen readers; use the TUI or CLI for accessible operation.

## Configuration

The server generates `server.yml` on first start and reads area settings from
`yaml.d/` and credentials from `secrets.yml`. It watches those files for
changes. Use `CONFIG_DIR` or `--config` to choose the configuration directory,
`DATA_DIR` or `--data` for application data, `PORT` or `--port` for the HTTP
port, and `LISTEN` or `--address` for the listener address. With no port set,
the server chooses an unused port on a host and uses port 80 in an application
container.

Operator branding content is stored separately in `yaml.d/branding.yml`. It
supports localized About and Privacy extensions, Terms overrides, a shared
header fragment, and a sanitized footer branding block. See the
[configuration reference](docs/configuration.md) for formats and placement.

The CLI profile is `cli.yml` in the user configuration directory. Set the
server with `--server` or `TIMEZONES_SERVER_PRIMARY`. For available settings,
environment inputs, and flags, see the [configuration reference](docs/configuration.md).
`server.mode` defaults to `production`; an empty `server.port` is selected at
startup (an unused host port, or port 80 in an application container).

## API

The versioned REST API is rooted at `/api/v1`. Public timezone reads do not
require a token. Token management requires the owner token after primary-token
bootstrap. MCP and protected metrics use bearer credentials. When configuring
your client, replace `https://YOUR_SERVER` with your deployment URL.

| Method and path | Purpose |
|---|---|
| `GET /api/v1/timezones` | List records with pagination and optional offset/DST filters |
| `GET /api/v1/timezones/search?q=...` | Search timezone names and identifiers |
| `GET /api/v1/timezones/utc/{iana-id}` | Look up an IANA identifier |
| `GET /api/v1/timezones/current/{iana-id}` | Return current local time and offset |
| `GET /api/v1/timezones/random` | Return a random record |
| `GET /api/v1/stats` and `/api/v1/count` | Dataset statistics and count |
| `GET /api/v1/server/healthz` | Bare JSON service health document (project/build/runtime, features, component checks, and aggregate stats) |
| `GET` or `POST /api/v1/server/graphql` | Read-only timezone and statistics queries |
| `POST /api/v1/server/preferences/tokens` | Create a linked token (owner token required) |
| `GET /api/autodiscover` | Discover API version and CLI release metadata |

OpenAPI JSON is available at `/api/v1/server/swagger`; interactive docs are at
`/server/docs/swagger`. GraphQL has a query explorer at `/server/docs/graphql`.
API routes default to JSON and support plain text with `.txt` or
`Accept: text/plain`; HTTP tools receive text by default. Browser routes return
server-rendered HTML, while text browsers receive HTML without JavaScript.
See the [API reference](docs/api.md) for routes, aliases, formats, and examples.

## Other

- [Full documentation](https://apimgr-timezones.readthedocs.io/)
- [Source repository](https://github.com/apimgr/timezones)
- `/mcp`, `/api/mcp`, and `/api/v1/server/mcp` expose the authenticated MCP endpoint.
- `/.well-known/security.txt`, `/server/security`, and `/server/contact` describe and accept security reports.
- The bundled dataset is static; this service does not publish future DST-transition schedules.

## Development

Development and release builds use the Docker Go toolchain configured by the
Makefile. `make build` creates debug binaries for the local platform and all
eight distribution targets; `make release` creates optimized production
artifacts, a source archive, and checksums without publishing. `make test` runs
the lint, vulnerability, unit-test, coverage, and basic binary checks.
`make docker` builds the host-platform image locally, applies the required OCI
metadata, and checks the runtime commands and Tor health. The check allows
offline runners to skip overlay bootstrap unless `REQUIRE_OVERLAY=1` is set.
`./tests/run_tests.sh` checks binaries and HTTP behavior in Docker;
`./tests/e2e.sh` runs browser
checks. `./tests/docker.sh` reuses executable server and CLI binaries supplied
via `BINDIR`, or builds them with the Docker toolchain otherwise; it checks
operator and owner/linked-token authorization, authenticated CLI reads and
token create/list/rename/revoke operations, private secret storage, debug-route
authorization, and Tor/I2P health and bootstrap state. Set `TOR_WAIT`
and `I2P_WAIT` to change overlay startup ceilings, or `REQUIRE_OVERLAY=1` to
make unavailable overlays fail the run. `./tests/incus.sh [systemd|openrc|all]`
exercises service lifecycle in disposable Incus containers. See
[Development](docs/development.md) for prerequisites and resource requirements.

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it at your
own risk.

- **No warranty:** The authors are not responsible for damages, data loss, or other issues arising from use of this software.
- **Data accuracy:** Timezone data may not reflect current civil-time rules. Verify time-critical decisions against authoritative sources.
- **Third-party services:** External service terms apply when a deployment configures those services.
- **Security:** No software is guaranteed to be free of vulnerabilities.
- **Production use:** Evaluate the software and deployment before production use.

## License

Project code is MIT licensed. The embedded timezone dataset is copied from
[`dmfilipenko/timezones.json`](https://github.com/dmfilipenko/timezones.json)
at commit `658abf632bae9c0313360702e3018024eb88a874`; its separate MIT license
and attribution are included in [LICENSE.md](LICENSE.md).
