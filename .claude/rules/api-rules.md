# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never keep obsolete project endpoints through redirects, aliases, compatibility shims, or deprecation periods.
- Never implement a route separately from the shared route registry or allow OpenAPI, GraphQL, MCP, and serializers to drift from it.
- Never leak hostnames, proxy input, errors, credentials, or unbounded pagination data.
- Never bypass certificate verification or weaken TLS/security defaults.

## CRITICAL - ALWAYS DO

- Define each API route once with its method, auth, visibility, request/response model, and primary fields.
- Use exact `/api/{api_version}/...` paths, canonical success/error envelopes, stable error codes, and documented pagination.
- Keep REST, GraphQL, MCP, OpenAPI, Swagger, text output, and route docs aligned with the same data contract.
- Expose health/version data according to the privacy tiers and configured hostname exposure switches.
- Follow the TLS certificate, ACME, renewal, and HTTP challenge rules in PART 15.

## Key Rules

- Health is bare JSON when specified; GraphQL and MCP keep their protocol-native envelopes.
- Validate host/proxy input against configured allowlists and sanitize forwarded URLs.
- Test route metadata and every registered method/client representation.

For complete details, see AI.md PART 13, 14, 15.
