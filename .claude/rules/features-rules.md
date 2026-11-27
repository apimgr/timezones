# Feature Rules (PART 17, 18, 19, 20, 21, 22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never silently discard email, scheduler, GeoIP, metrics, backup, restore, or update failures.
- Never make optional external data sources or notification providers startup-critical.
- Never expose private backup paths, security reports, tokens, or notification secrets.
- Never run external cron or update binaries without the configured integrity and permission checks.

## CRITICAL - ALWAYS DO

- Follow SMTP/provider configuration, safe message handling, notification routing, and report privacy rules.
- Keep scheduled tasks bounded, lock-protected, observable, and under the documented lifecycle.
- Treat GeoIP and other downloaded datasets as non-critical, validate updates, and preserve attribution.
- Protect metrics with their configured credentials and expose only the documented formats.
- Verify backup archives, enforce size/path limits, protect encryption secrets, and make restore explicit and safe.
- Verify update channel, signature/checksum, platform, and rollback/error behavior as specified.

## Key Rules

- Use exact configuration paths and defaults from the corresponding PART.
- Every feature needs operational status, documented failure behavior, and tests for its required failure paths.

For complete details, see AI.md PART 17, 18, 19, 20, 21, 22.
