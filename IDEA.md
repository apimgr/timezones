## Project description

Timezones provides a comprehensive reference dataset of world timezones through a versioned REST API, GraphQL endpoint, and web interface. Each timezone entry includes its display name, abbreviation, UTC offset, daylight-saving flag, human-readable label, and associated IANA identifiers. A companion terminal client enables timezone lookup and access to the service's public API.

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
owner_token: timezones_owner_token_wxaarpwb

## Business logic

### Product scope & non-goals

**In scope:**
- Full world timezone dataset with IANA identifiers, UTC offsets, abbreviations, DST flags, and human-readable labels
- Single timezone lookup by display name, abbreviation, or IANA identifier (case-insensitive)
- List all timezones (with optional UTC offset or DST filter)
- Current time in any timezone
- Full web frontend with dark/light/automatic themes, installable offline support, and mobile-friendly layouts
- Server information, help, health, privacy, terms, preferences, and API documentation pages required by AI.md
- Anonymous owner and linked-token management without user accounts
- CLI client (`timezones-cli`) for terminal lookup: `timezones-cli "America/New_York"`
- Interactive OpenAPI documentation
- A read-only GraphQL API and query explorer for timezone lookups
- An authenticated Streamable HTTP MCP endpoint for timezone tools

**Non-goals:**
- No user accounts, registration, or password-based login; bearer tokens do not represent user accounts
- No browser-based server administration
- No future daylight-saving transition schedule or forecast
- No calendar or scheduling features
- No paid tiers or access tiers; bearer tokens provide ownership and privileged access, not a subscription level

**Data sources:** The bundled timezone record dataset is sourced from
[`dmfilipenko/timezones.json`](https://github.com/dmfilipenko/timezones.json).
Its source and license are recorded in `LICENSE.md`.

### Roles & permissions

There are no user accounts or password-based roles. Public REST reads and public pages require no token. The MCP endpoint, token-management actions, metrics, and operator-only operations require the Bearer credentials and permissions defined by AI.md.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Public web pages and public read-only API access |
| **Anonymous API client (curl/CLI)** | Public read-only timezone REST access |
| **Owner token holder** | Uses the owner's resources and manages linked tokens where supported |
| **MCP client** | Uses an owner or linked Bearer token; receives only the operations allowed to that token |
| **Server operator** | Configures the server through `server.yml` and uses the operator token for protected operations |

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

Timezone records contain no PII. Owner and token metadata is stored as described by AI.md; raw tokens are not returned by token-list endpoints.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Bundled timezone dataset | Trusted reference input | Source and redistribution terms are documented in `LICENSE.md` |
| Incoming HTTP requests | **Untrusted** | Requests may contain malformed or oversized values |
| Configured outbound integrations | Trusted only for the configured purpose | Optional mail and keyserver integrations follow AI.md failure-handling rules; public timezone reads do not depend on them |

### Threat model & abuse cases

**Primary assets:** service availability, operator credentials, and owner-token integrity.

**Attacker/abuser goals:**
- DoS via high-rate requests

**Defenses:**
- Rate limiting on all endpoints
- All query parameters validated and length-capped
- No passwords or user accounts eliminates password-based credential stuffing; Bearer-token handling follows AI.md

### Security decisions & exceptions

- **Public timezone API**: reads require no token. MCP and privileged token-management/operational routes require the Bearer tokens specified in AI.md.
- Public timezone data is available to cross-origin browser clients under the CORS policy required by AI.md.

### API routes

These are the timezone-domain and GraphQL routes defined by this project. Public data reads need no credential.

| Method | Path | Auth | Purpose and request/response summary |
|--------|------|------|--------------------------------------|
| GET | `/api/v1/timezones` | Public | List records; accepts pagination and optional `offset` and `isdst` filters; returns records and pagination metadata. |
| GET | `/api/v1/timezones.json` | Public | Return the embedded dataset in a JSON response envelope. |
| GET | `/api/v1/timezones/search?q={query}` | Public | Search display name, abbreviation and IANA identifiers; returns matches and pagination metadata. |
| GET | `/api/v1/timezones/offset/{offset}` | Public | Return paginated records matching a UTC offset. |
| GET | `/api/v1/timezones/abbr/{abbr}` | Public | Return paginated records matching an abbreviation. |
| GET | `/api/v1/timezones/utc/{iana-id}` | Public | Return the record matching an IANA identifier. |
| GET | `/api/v1/timezones/value/{name}` | Public | Return the record matching a display name. |
| GET | `/api/v1/timezones/random` | Public | Return one randomly selected timezone record. |
| GET | `/api/v1/timezones/current/{iana-id}` | Public | Return current local time for an IANA identifier. |
| GET | `/api/v1/timezones.txt` | Public | Return the full timezone list as plain text. |
| GET | `/api/v1/timezones/search.txt?q={query}` | Public | Return matching records as plain text. |
| GET | `/api/v1/timezones/random.txt` | Public | Return one random record as plain text. |
| GET | `/api/v1/stats` | Public | Return aggregate dataset statistics. |
| GET | `/api/v1/stats.txt` | Public | Return aggregate dataset statistics as plain text. |
| GET | `/api/v1/count` | Public | Return the record count in the standard JSON envelope. |
| GET | `/api/v1/count.txt` | Public | Return the record count as plain text. |
