# Security and privacy

## Access model

Timezone reads and public information pages do not require authentication.
Owner and linked Bearer tokens protect token management; MCP and metrics use
their configured credentials. The CLI setup check at
`/api/v1/server/auth/check` validates a Bearer token and returns its kind.
There are no user accounts or password-based login. Send API credentials in
`Authorization: Bearer …`, and use HTTPS on clearnet deployments.

## Public status and security endpoints

Health is intentionally public at `/server/healthz`, `/server/healthz.txt`,
`/api/v1/server/healthz`, and `/api/v1/server/healthz.txt`, with equivalent
unversioned aliases under `/api/healthz`. It reports project/build/runtime
information, service status, enabled features, component and dataset checks,
and public-safe aggregate request counts. Overlay status is included when
enabled. It does not return tokens, database credentials, or filesystem paths.
Metrics are public only when configured to allow unauthenticated access;
otherwise use the corresponding service Bearer token.

The current public `/.well-known/` entries are:

| Path | Behavior |
|---|---|
| `/.well-known/security.txt` | Dynamic security contact and disclosure information |
| `/.well-known/pgp-key.asc` | Active security-report PGP public key when published; otherwise unavailable |
| `/.well-known/llms.txt` | Machine-readable endpoint and service description |

`/llms.txt` is an alias for the last entry. Unsupported well-known paths return
404. OAuth metadata, WebFinger, Android App Links, Apple association, and
MTA-STS are not enabled by this project.

## Request protections

The server limits request targets and bodies, adds CSP and browser security
headers, and trusts forwarded client-address headers only when the immediate
peer is a private or loopback address. Its general rate limiter currently
allows 240 requests per minute per client address; strict token operations
have a separate limit. This general limiter is fixed in code. The configured
`server.rate_limit.*` values are described in the
[configuration reference](configuration.md), but are not currently used to
set endpoint-class limits.

The server can enforce its configured IP/domain blocklist and GeoIP policy.
`server.security.allowlist` entries bypass those checks and rate limiting.
Blocklist source health and GeoIP database health appear on the public health
endpoint. See [configuration](configuration.md) for default blocklist settings
and the current configuration surface.

## Browser state and data

The web interface uses first-party theme, language, consent, CSRF, and
owner-token cookies. Owner-token values are stored as hashes in the token
database; raw tokens are returned only during issuance or explicit reveal, and
list responses contain metadata. Timezone lookups are not stored as
application records. The HTTP access logger records request information;
deployments may add proxy or hosting logs with separate retention settings.

## Vulnerability reporting

Use [GitHub private vulnerability reporting](https://github.com/apimgr/timezones/security/advisories/new)
for source-code vulnerabilities. For a deployed instance, see
`/server/security`, `/.well-known/security.txt`, and
`/server/contact?security_id=...`. The report form uses review and confirmation
steps. `/server/security/report/{tracking_id}` exposes report status and
researcher-facing comments, not submitted report contents. Reports use the
configured PGP public key when available and otherwise the server's encrypted
storage path. When SMTP is available, the researcher receives an acknowledgment
and later status changes send a notice containing the tracking ID and new status,
never report content. Without working SMTP, those messages are skipped. See
[SECURITY.md](https://github.com/apimgr/timezones/blob/main/.github/SECURITY.md)
for supported releases and response targets.
