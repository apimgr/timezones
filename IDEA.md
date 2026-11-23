## Project description

Timezones is a full-stack Go web application providing a comprehensive reference dataset of world timezones through a versioned REST API, GraphQL endpoint, and a server-side rendered web UI. Each timezone entry includes its display name, abbreviation, UTC offset, DST status, human-readable label, and the full list of IANA timezone identifiers it covers. All data is embedded in the binary at build time. A companion CLI tool enables timezone lookup from the terminal. Deployed as a single self-contained static binary.

## Project variables

project_name: timezones
project_org: apimgr
internal_name: timezones
internal_org: apimgr
app_name: Timezones API
repo: https://github.com/apimgr/timezones
license: MIT
binary: timezones
client_binary: timezones-cli

## Business logic

### Product scope & non-goals

**In scope:**
- Full world timezone dataset with IANA identifiers, UTC offsets, abbreviations, DST flags, and human-readable labels
- Single timezone lookup by display name, abbreviation, or IANA identifier (case-insensitive)
- List all timezones (with optional UTC offset or DST filter)
- Current time in any timezone
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Server pages: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`
- CLI client (`timezones-cli`) for terminal lookup: `timezones-cli "America/New_York"`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/api/{api_version}/server/graphql` with `/api/graphql` as the latest-version alias

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No future daylight-saving transition schedule or forecast
- No calendar or scheduling features
- No paid tiers, no API keys, no rate-limited access tiers

### Roles & permissions

There are no user roles. All endpoints are public and require no authentication.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to all web pages and API endpoints |
| **Anonymous API client (curl/CLI)** | Full read access to all API endpoints |
| **Server operator** | Configures server via `server.yml` only; no web management interface |

### Data model & sensitivity

**Timezone record** (embedded at build time, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `value` | string — timezone display name (e.g., `Eastern Standard Time`) | Public |
| `abbr` | string — abbreviation (e.g., `EST`) | Public |
| `offset` | integer — UTC offset in hours | Public |
| `isdst` | boolean — whether DST is currently observed | Public |
| `text` | string — human-readable label (e.g., `(UTC-05:00) Eastern Time (US & Canada)`) | Public |
| `utc` | string[] — IANA timezone identifiers covered (e.g., `["America/New_York"]`) | Public |

No PII stored or served.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Timezone dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | All query parameters validated and size-capped |

No external services called at runtime.

### Threat model & abuse cases

**Primary assets:** service availability.

**Attacker/abuser goals:**
- DoS via high-rate requests

**Defenses:**
- Rate limiting on all endpoints
- All query parameters validated and length-capped
- No user accounts eliminates credential stuffing and privilege escalation entirely

### Security decisions & exceptions

- **No authentication on any endpoint**: intentional. Public read-only reference API.
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public data API designed for cross-origin browser use.

### API routes

All routes are public and require no authentication.

| Method | Path | Purpose and request/response summary |
|--------|------|--------------------------------------|
| GET | `/api/v1/timezones` | List timezone records as JSON; accepts optional `offset` and `isdst` filters. |
| GET | `/api/v1/timezones.json` | Return the embedded JSON dataset. |
| GET | `/api/v1/timezones/search?q={query}` | Search by display name, abbreviation, or IANA identifier; returns matching records. |
| GET | `/api/v1/timezones/offset/{offset}` | Return records matching a UTC offset in hours. |
| GET | `/api/v1/timezones/abbr/{abbr}` | Return records matching an abbreviation. |
| GET | `/api/v1/timezones/utc/{iana-id}` | Return a record matching an IANA identifier. |
| GET | `/api/v1/timezones/value/{name}` | Return a record matching a display name. |
| GET | `/api/v1/timezones/random` | Return one random timezone record. |
| GET | `/api/v1/stats` | Return aggregate dataset counts. |
| GET | `/api/v1/timezones/current/{iana-id}` | Return the current local time for an IANA timezone identifier. |
| GET | `/api/v1/timezones.txt` | Return the full list as plain text. |
| GET | `/api/v1/timezones/search.txt?q={query}` | Return matching records as plain text. |
| GET | `/api/v1/timezones/random.txt` | Return one random record as plain text. |
| GET | `/api/v1/stats.txt` | Return dataset statistics as plain text. |
| GET | `/api/v1/count` | Return the record count as JSON. |
| GET | `/api/v1/count.txt` | Return the record count as plain text. |
| GET | `/api/v1/` | Return API name, version, and endpoint entry points. |
| GET | `/api/v1/server/healthz` | Return public service health and version information. |
| GET | `/api/v1/server/swagger` | Return the OpenAPI description. |
| GET | `/api/swagger` | Return the current OpenAPI description for API tooling. |
| POST | `/api/v1/server/graphql` | Execute a public read-only GraphQL query over timezone data. |
| POST | `/api/graphql` | Execute the same GraphQL handler for the latest API version. |
| GET | `/server/docs/graphql` | Open the GraphiQL interface for the GraphQL API. |
| GET | `/server/about` | Render project description and features. |
| GET | `/server/help` | Render real API usage examples and help information. |
| GET | `/server/healthz` | Render public service health. |
| GET | `/timezones/current/{iana-id}` | Render the current local time for an IANA timezone identifier. |
| GET | `/server/privacy` | Render the deployment privacy notice. |
| GET | `/server/terms` | Render terms for the public reference service. |
