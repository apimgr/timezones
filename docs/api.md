# API reference

The canonical REST API is versioned under `/api/v1`. Public timezone and server
information reads do not need a token. Token-management routes use the owner
token after the initial primary-token claim. MCP and metrics use their
configured bearer credentials. All protected native API credentials are sent
in `Authorization: Bearer <token>`; query-string tokens are not credentials.

## Responses

Most successful REST responses use `{ "ok": true, "data": ... }`. List
responses also include a `pagination` object with `page`, `limit`, `total`, and
`pages`. REST errors use `{ "ok": false, "error": "CODE", "message": "..." }`
and may include `details`.

Health routes return the health document directly. GraphQL uses the standard
GraphQL response shape (`data` and optional `errors`), and MCP uses JSON-RPC.
Binary CLI downloads return the verified artifact bytes rather than a JSON
envelope.

## Timezone data

| Method | Path | Purpose |
|---|---|---|
| GET | `/api/v1/` | API name, version, and endpoint entry points |
| GET | `/api/v1/timezones` | List records; supports `page`, `limit`, `offset`, and `isdst` |
| GET | `/api/v1/timezones.json` | Return the bundled JSON dataset |
| GET | `/api/v1/timezones/search?q={query}` | Search names, descriptions, abbreviations, and IANA identifiers |
| GET | `/api/v1/timezones/offset/{hours}` | Find records with a UTC offset |
| GET | `/api/v1/timezones/abbr/{abbr}` | Find records by abbreviation |
| GET | `/api/v1/timezones/utc/{iana-id}` | Find a record by IANA identifier; slashes are accepted in the path |
| GET | `/api/v1/timezones/value/{name}` | Find a record by display name |
| GET | `/api/v1/timezones/random` | Return a random record |
| GET | `/api/v1/timezones/current/{iana-id}` | Return current time and offset for an IANA identifier |
| GET | `/api/v1/stats` | Dataset statistics |
| GET | `/api/v1/count` | Dataset record count |

API requests default to JSON. Append `.txt` to any registered API route or send
`Accept: text/plain` to request its text representation; for example,
`/api/v1/timezones.txt`, `/api/v1/server/about.txt`, and
`/api/v1/server/healthz.txt`. Binary download routes keep their binary response
and do not have a text form. HTTP tools receive text by default; text browsers
receive no-JavaScript HTML on browser routes and text on API routes. Frontend
`Accept: application/json` requests use the paired API response when that route
is registered.

## Server API

| Method | Path | Purpose |
|---|---|---|
| GET | `/api/v1/server/about` | About information |
| GET | `/api/v1/server/auth/check` | Validate the supplied Bearer token for client setup; requires a valid credential and returns its kind |
| GET | `/api/v1/server/help` | Help content and registered help pages |
| GET | `/api/v1/server/help/mcp` | MCP integration help |
| GET | `/api/v1/server/privacy` | Privacy policy facts, cookie categories, tracking, retention, sharing, third parties, CCPA status, and localized content |
| GET | `/api/v1/server/terms` | Terms information |
| GET | `/api/v1/server/dpo` | Data protection contact information |
| GET | `/api/v1/server/security` | Security contacts and reporting information |
| GET | `/api/v1/server/security/policy` | Coordinated disclosure policy in Markdown |
| GET | `/api/v1/server/security/thanks` | Public acknowledgments with credit, disclosure date, and severity |
| GET, POST | `/api/v1/server/contact` | Read contact-form fields, then submit a contact or security report |
| GET, POST | `/api/v1/server/preferences` | Read and save guest preferences |
| GET | `/api/v1/server/preferences/sync` | Create/export a preference-sync code |
| GET | `/api/v1/server/preferences/sync/import` | Import preferences by code or URL parameters |
| GET | `/api/v1/server/healthz` | Health status |
| GET | `/api/v1/server/metrics` | Prometheus metrics |
| GET | `/api/v1/server/metrics/{service}` | Prometheus, Grafana, or Loki integration output |
| GET | `/api/v1/server/swagger` | OpenAPI JSON document |
| GET, POST | `/api/v1/server/graphql` | Read-only GraphQL queries |
| POST | `/api/v1/server/mcp` | MCP Streamable HTTP endpoint; bearer token required |

About and Terms responses contain `title`, raw `content`, and `format`. The
content is the request-language source after configured placeholders have been
replaced; HTML is returned as source and is never returned as rendered page
markup. Operator content comes from `server.branding.content.about` and
`server.branding.content.terms` in `yaml.d/branding.yml`. Formats are
`markdown`, `html`, or `text`. Without a usable override, About uses the
project description and Terms uses the built-in terms text; both default to
`markdown`.

Metrics are protected by their service-specific bearer token unless the
operator enables unauthenticated metrics. The `service` value is `prometheus`,
`grafana`, or `loki`.

### Health document

`GET /api/v1/server/healthz` returns a bare JSON health document rather than
the standard REST envelope. It includes project branding, the overall status,
version and build information, Go version, mode, uptime and UTC timestamp,
enabled features, component and dataset checks, and public-safe aggregate
request statistics. The server-rendered `/server/healthz` page presents the
same health fields. Frontend text requests to `/server/healthz`,
`/server/healthz.txt`, or `/healthz` use human-readable labels, durations, and
timestamps. API text requests to `/api/v1/server/healthz.txt` or with
`Accept: text/plain` use dotted machine-readable keys and compact values.

Health status is `healthy`, `degraded`, `restart_required`, `unhealthy`,
`maintenance`, or `shutting_down`. Unhealthy, maintenance, and shutdown
responses return HTTP 503; healthy, degraded, and restart-required responses
return HTTP 200. Overlay features report lifecycle status separately, and
disabled Tor or I2P services do not add component checks. Public health output
contains no credentials or filesystem paths.

The privacy response is generated from the active privacy, tracking, retention,
third-party, DPO, and branding configuration. Its `summary`, `cookies`,
`data`, `tracking`, `retention`, `third_party`, `ccpa`, and `content` members
are always present. The CCPA opt-out link and state are included even when the
opt-out UI is not applicable; `ccpa.applicable` follows `server.privacy.data.sold`.
The consent and data-usage copy is localized and reflects the sold-data setting.

### Owner-token routes

| Method | Path | Purpose |
|---|---|---|
| GET | `/api/v1/server/preferences/tokens` | List token metadata |
| POST | `/api/v1/server/preferences/tokens/primary` | Claim the primary owner token when no valid owner credential is presented |
| POST | `/api/v1/server/preferences/tokens/primary/reveal` | Reveal the primary token |
| POST | `/api/v1/server/preferences/tokens/primary/rotate` | Rotate the primary token |
| DELETE | `/api/v1/server/preferences/tokens/primary` | Revoke the primary token |
| POST | `/api/v1/server/preferences/tokens` | Create a linked token |
| POST | `/api/v1/server/preferences/tokens/{id}/reveal` | Reveal a linked token |
| POST | `/api/v1/server/preferences/tokens/{id}/rotate` | Rotate a linked token |
| PATCH | `/api/v1/server/preferences/tokens/{id}` | Update linked-token metadata |
| DELETE | `/api/v1/server/preferences/tokens/{id}` | Revoke a linked token |

These routes require the owner token after bootstrap. The raw token is returned
only during issuance or an explicit reveal operation; list responses contain
metadata.

## GraphQL

The schema is read-only and supports timezone lookup, current-time lookup, and
dataset statistics. Use GET for queries that fit in the URL, or POST a JSON
object containing `query` and optional `variables` and `operationName`.

```sh
curl -q -G -LSsf --data-urlencode \
  'query={ timezone(identifier: "America/New_York") { value abbr utc } }' \
  https://YOUR_SERVER/api/v1/server/graphql
```

For local development, replace the host with the address and port where the
server is listening.

GET accepts `query` plus optional JSON-object `variables` and `operationName`
parameters. A GET without `query` returns a short usage response for the
no-JavaScript explorer link. The request limit is 8 KiB and query depth is
limited to four fields. The explorer is at `/server/docs/graphql`.

## Aliases and discovery

These unversioned paths directly mount the same handlers as their canonical
versioned routes; they do not redirect:

| Alias | Canonical route |
|---|---|
| `/api/swagger` | `/api/v1/server/swagger` |
| `/api/graphql` | `/api/v1/server/graphql` |
| `/api/healthz` | `/api/v1/server/healthz` |
| `/api/healthz.txt` | `/api/v1/server/healthz.txt` |
| `/api/metrics[/{service}]` | `/api/v1/server/metrics[/{service}]` |
| `/metrics[/{service}]` | `/api/v1/server/metrics[/{service}]` |
| `/api/mcp` | `/api/v1/server/mcp` |
| `/mcp` | `/api/v1/server/mcp` |

`GET /api/autodiscover` is unversioned by design. It returns server identity,
the current API version, public origins, and CLI release metadata. Verified
binary downloads use `GET /cli/binaries/{binary}`; only recognized artifacts
with a valid release manifest and SHA-256 checksum are served.

The interactive OpenAPI page is `/server/docs/swagger`; its JSON document is
served at `/api/v1/server/swagger`. The API route inventory in
[`IDEA.md`](https://github.com/apimgr/timezones/blob/main/IDEA.md) defines the project-specific contract.
