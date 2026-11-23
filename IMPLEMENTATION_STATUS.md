# Implementation status

`AI.md` is the implementation source of truth and `IDEA.md` defines the product scope. This status summarizes the current tree; it is not a claim of full compliance.

The service embeds timezone records and web assets, provides versioned REST lookup endpoints with response envelopes and pagination, a server-rendered homepage and server pages, theme preferences, and the `timezones-cli` binary. It also serves a GraphQL endpoint and explorer, a bundled Swagger UI, an OpenAPI JSON description, current-time views, and robots, sitemap, Atom, RSS, security.txt, and llms.txt discovery resources.

The implementation remains incomplete against the full AI.md baseline. Verified remaining work is tracked in [TODO.AI.md](TODO.AI.md). Major gaps include the broad server configuration and lifecycle baseline, authentication/MCP requirements pending resolution of an AI.md and IDEA.md conflict, complete content negotiation, localization, and optional platform integrations.

Validation completed in this pass: `make test` passed (including `go vet`, `govulncheck`, and Go tests); `golangci-lint run` passed with zero issues; `mkdocs build --strict` passed; and server and CLI binaries compiled in Docker. A container running the built server returned successful responses for the rendered feature pages, docs, static assets, sitemap, and feeds; headless Chromium loaded the homepage and its CSS, JavaScript, service worker, manifest, and icon. Full browser interaction, responsive/accessibility review, provider CI, and release workflows remain outstanding.
