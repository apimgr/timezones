# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never expose credentials, private paths, raw tokens, report contents, or unredacted errors in public responses or logs.
- Never process untrusted HTML, Markdown, archives, or files without the required sanitization and path limits.
- Never bypass configured rate limits, CSRF/origin checks, authenticated blob access, or private `no-store` controls.
- Never make Tor or I2P a required dependency; they remain independently managed optional runtime daemons.

## CRITICAL - ALWAYS DO

- Use stable JSON error envelopes, appropriate cache headers, ETags, and version-change cache purging.
- Use SQLite by default, support configured libsql/Turso, migrate schema automatically, pool connections, and use prepared statements.
- Apply security headers, HSTS when TLS is active, audit logging, redaction, abuse controls, and content policies from AI.md.
- Keep public health data safe and aggregate; overlay hostnames are exposed only when the matching `server.*.expose` opt-in is true.
- Follow the Tor-required and I2P runtime-default rules in PART 31 exactly.

## Key Rules

- Sanitized Markdown uses the approved Goldmark and Bluemonday pipeline; operator page content uses the one `operatorHTML` renderer.
- Enforce request size, timeout, body, archive, and filesystem boundaries before expensive or persistent work.
- Keep database, cache, backup, security-report, and overlay failures in their documented lifecycle states.

For complete details, see AI.md PART 9, 10, 11, 31.
