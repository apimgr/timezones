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
- Server pages: `/server/about`, `/server/help`, `/server/help/mcp`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/preferences`, `/server/preferences/tokens`, `/server/docs/swagger`, and `/server/docs/graphql`
- Owner and linked-token management at `/server/preferences/tokens` and its API mirror
- CLI client (`timezones-cli`) for terminal lookup: `timezones-cli "America/New_York"`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/api/{api_version}/server/graphql` with `/api/graphql` as the latest-version alias
- Authenticated Streamable HTTP MCP at `/api/{api_version}/server/mcp` with `/api/mcp` and `/mcp` aliases
- Bearer-token management at `/server/preferences/tokens` and its API mirror, as required by AI.md

**Non-goals:**
- No user accounts, registration, or password-based login; bearer tokens do not represent user accounts
- No admin web panel (server configured via `server.yml` only)
- No future daylight-saving transition schedule or forecast
- No calendar or scheduling features
- No paid tiers or access tiers; bearer tokens provide ownership and privileged access, not a subscription level

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
| Timezone dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | All query parameters validated and size-capped |

No external services called at runtime.

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
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public data API designed for cross-origin browser use.

### API routes

The routes below are project-defined. Public timezone reads need no credential. Protected routes require the credential and role shown in the table; cookie fallback applies only to the specified token-management web forms.

| Method | Path | Auth | Purpose and request/response summary |
|--------|------|------|--------------------------------------|
| GET | `/api/v1/` | Public | Return API name, version and endpoint links in the standard success envelope. |
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
| GET | `/api/v1/server/healthz` | Public | Return public service health and version information. |
| GET | `/api/v1/server/healthz.txt` | Public | Return public service health as plain text. |
| GET | `/api/healthz` | Public | Current-version alias for the public health response. |
| GET | `/api/healthz.txt` | Public | Current-version plain-text health alias. |
| GET | `/api/autodiscover` | Public | Return supported server settings, config schema and CLI options for clients. |
| GET | `/api/v1/server/about` | Public | Return about-page content in the standard JSON envelope. |
| GET | `/api/v1/server/help` | Public | Return API usage and page links in the standard JSON envelope. |
| GET | `/api/v1/server/help/mcp` | Public | Return MCP connection guidance, endpoint URLs, authentication and tool metadata as JSON. |
| GET | `/api/v1/server/preferences` | Public | Return current guest preferences and the synchronization URL/code. |
| GET | `/api/v1/server/preferences/sync` | Public | Return preferences plus a shareable synchronization URL and short code. |
| GET | `/api/v1/server/preferences/sync/import` | Public | Import valid preference values from query parameters and return applied/ignored fields. |
| GET | `/api/v1/server/privacy` | Public | Return the privacy notice in the standard JSON envelope. |
| GET | `/api/v1/server/terms` | Public | Return terms in the standard JSON envelope. |
| GET | `/api/v1/server/dpo` | Public | Return the configured data protection officer name, email and privacy URL. |
| GET | `/api/v1/server/contact` | Public | Describe the report form fields, issue a fresh captcha and form token, and return the security URL and abuse contact without exposing recipients. |
| POST | `/api/v1/server/contact` | Public; form token, captcha and rate limits | Validate a report in a review/confirmation flow and deliver it only after confirmation. |
| GET | `/api/v1/server/security` | Public | Return configured disclosure contacts in RFC 9116 order, expiration, PGP availability, policy, thanks and report URLs. |
| GET | `/api/v1/server/security/policy` | Public | Return the disclosure policy Markdown. |
| GET | `/api/v1/server/security/thanks` | Public | Return disclosed researcher acknowledgments without contact information. |
| GET | `/api/v1/server/swagger` | Public | Return the OpenAPI description as JSON. |
| GET | `/api/swagger` | Public | Current-version alias for the OpenAPI description. |
| POST | `/api/v1/server/graphql` | Public | Accept a GraphQL JSON request and return read-only timezone query results. |
| POST | `/api/graphql` | Public | Execute the latest-version GraphQL handler. |
| GET | `/api/metrics[/{service}]` | Per-service metrics Bearer token | Return the configured service's metrics via the current-version alias. |
| GET | `/api/v1/server/metrics[/{service}]` | Per-service metrics Bearer token | Return metrics for the configured service. |
| GET | `/server/metrics[/{service}]` | Per-service metrics Bearer token | Return the configured service's metrics. |
| GET | `/metrics[/{service}]` | Per-service metrics Bearer token | Root metrics alias when enabled. |
| POST | `/api/v1/server/preferences` | None | Apply a partial guest-preference JSON update and return the resulting preferences. |
| POST | `/api/v1/server/mcp` | Owner or operator Bearer | Accept Streamable HTTP MCP JSON-RPC requests and return protocol responses. |
| POST | `/api/mcp` | Owner or operator Bearer | Current-version MCP alias. |
| POST | `/mcp` | Owner or operator Bearer | Root MCP alias. |
| POST | `/api/v1/server/preferences/tokens/primary` | None; strict rate limit | Issue a new owner and primary token; returns the token in `owner_token`, unless the caller already has a valid token. |
| GET | `/api/v1/server/preferences/tokens` | Primary Bearer | Return masked primary and linked-token metadata. |
| POST | `/api/v1/server/preferences/tokens/primary/reveal` | Primary Bearer | Reveal the primary token in the response after explicit request. |
| POST | `/api/v1/server/preferences/tokens/primary/rotate` | Primary Bearer | Rotate the primary token in place and return its new value. |
| DELETE | `/api/v1/server/preferences/tokens/primary` | Primary Bearer | Delete the owner, its tokens and its resources. |
| POST | `/api/v1/server/preferences/tokens` | Primary Bearer | Create a linked token from a name and expiry; return its value once in the response. |
| POST | `/api/v1/server/preferences/tokens/{id}/reveal` | Primary Bearer | Reveal the selected linked token. |
| POST | `/api/v1/server/preferences/tokens/{id}/rotate` | Primary Bearer | Rotate the selected linked token in place. |
| PATCH | `/api/v1/server/preferences/tokens/{id}` | Primary Bearer | Rename the selected linked token or change its expiry. |
| DELETE | `/api/v1/server/preferences/tokens/{id}` | Primary Bearer | Revoke the selected linked token. |

#### Browser and discovery routes

| Method | Path | Auth | Purpose and request/response summary |
|--------|------|------|--------------------------------------|
| GET | `/` | Public | Render the searchable timezone homepage. |
| GET | `/timezones` | Public | Render a paginated timezone list. |
| GET | `/timezones/search?q={query}` | Public | Render matching timezone records. |
| GET | `/timezones/offset/{offset}` | Public | Render timezone records matching an offset. |
| GET | `/timezones/abbr/{abbr}` | Public | Render timezone records matching an abbreviation. |
| GET | `/timezones/utc/{iana-id}` | Public | Render the record for an IANA identifier. |
| GET | `/timezones/value/{name}` | Public | Render the record matching a display name. |
| GET | `/timezones/random` | Public | Render a random timezone record. |
| GET | `/timezones/current/{iana-id}` | Public | Render current local time for an IANA identifier. |
| GET | `/stats` | Public | Render aggregate dataset statistics. |
| GET | `/count` | Public | Render the number of timezone records. |
| GET | `/server/about` | Public | Render project scope and service information. |
| GET | `/server/help` | Public | Render API usage and links to registered help pages. |
| GET | `/server/help/mcp` | Public | Explain MCP connection, authentication, tools and limits. |
| GET | `/server/privacy` | Public | Render privacy and data-handling details. |
| GET | `/server/contact` | Public | Render the abuse, DMCA, privacy or security report form; security mode is selected only by a valid `security_id`. |
| GET | `/server/dpo` | Public | Render data protection officer contact details. |
| GET | `/server/security` | Public | Render live disclosure contacts and instructions, with links to policy, acknowledgments, PGP key and reporting form. |
| GET | `/server/security/policy` | Public | Render the configured or built-in coordinated disclosure policy. |
| GET | `/server/security/thanks` | Public | Render researchers who opted into credit for disclosed reports. |
| GET | `/server/security/report/{tracking_id}` | One-use daily status token in URL | Show the report's researcher-visible status, comment and expected disclosure date; no API mirror. |
| GET | `/server/terms` | Public | Render terms for using the reference service. |
| GET | `/server/healthz` | Public | Render service status and public health details. |
| GET | `/server/healthz.txt` | Public | Return public service health as plain text. |
| GET | `/healthz` | Public; optional configured alias | Direct health-page alias when enabled by configuration. |
| GET | `/server/docs/swagger` | Public | Render interactive OpenAPI documentation with a no-JavaScript link. |
| GET | `/server/docs/graphql` | Public | Render the GraphQL workbench with a no-JavaScript endpoint link. |
| GET | `/server/preferences` | Public | Render guest preferences and links to registered preference pages. |
| GET | `/server/preferences/sync` | Public | Render preference export and import controls. |
| POST | `/server/preferences/sync` | CSRF token | Apply supported preference values and redirect with post/redirect/get. |
| GET | `/server/preferences/sync/import` | Public | Import preference values from a shared URL and redirect. |
| GET | `/server/preferences/tokens` | Public page; owner cookie for token list/actions | Render token creation and restore forms anonymously; show owner token controls only for a valid owner cookie. |
| POST | `/server/preferences` | CSRF token | Set the browser's theme preference and redirect to a same-origin page. |
| POST | `/server/preferences/token/create` | None; strict rate limit | Issue a primary owner token, set the owner cookie and render the new token. |
| POST | `/server/preferences/token/restore` | Pasted primary token; strict rate limit and CSRF | Validate a primary token and restore the browser owner cookie. |
| POST | `/server/preferences/token/reveal` | Owner cookie; primary token and CSRF | Reveal the primary token for one page load. |
| POST | `/server/preferences/token/rotate` | Owner cookie; primary token and CSRF | Rotate the primary token and update the owner cookie. |
| POST | `/server/preferences/token/revoke` | Owner cookie; primary token and CSRF | Delete the owner, its tokens and its resources after confirmation. |
| POST | `/server/preferences/token/forget` | CSRF token | Clear this browser's owner cookie and saved token copy without revoking it. |
| POST | `/server/preferences/tokens` | Owner cookie; primary token and CSRF | Create a linked token from a name and expiry. |
| POST | `/server/preferences/tokens/{id}/reveal` | Owner cookie; primary token and CSRF | Reveal a linked token for one page load. |
| POST | `/server/preferences/tokens/{id}/rotate` | Owner cookie; primary token and CSRF | Rotate the selected linked token in place. |
| POST | `/server/preferences/tokens/{id}/update` | Owner cookie; primary token and CSRF | Rename the selected linked token or change its expiry. |
| POST | `/server/preferences/tokens/{id}/revoke` | Owner cookie; primary token and CSRF | Revoke the selected linked token. |
| GET | `/favicon.ico` | Public | Serve the embedded favicon. |
| GET | `/manifest.json` | Public | Serve the PWA manifest. |
| GET | `/sw.js` | Public | Serve the application service worker. |
| GET | `/robots.txt` | Public | Serve configured crawler rules and an absolute sitemap URL. |
| GET | `/sitemap.xml` | Public | Serve the dynamically generated sitemap of allowed public HTML pages. |
| GET | `/feed.xml` | Public | Serve an Atom feed with public resource updates, or a valid empty feed. |
| GET | `/rss.xml` | Public | Serve an RSS feed with the same public resource updates. |
| GET | `/.well-known/llms.txt` | Public | Serve the agent-readable service and API summary. |
| GET | `/llms.txt` | Public | Alias for the agent-readable service summary. |
| GET | `/.well-known/security.txt` | Public | Serve security contact and reporting guidance. |
| GET | `/.well-known/pgp-key.asc` | Public; 404 until a key exists | Serve the configured public report-encryption key. |
| GET | `/static/*` | Public | Serve embedded CSS, JavaScript, images and vendored documentation assets. |
