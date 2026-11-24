# Implementation status

`AI.md` is the implementation source of truth and `IDEA.md` defines the product scope. This status summarizes the current tree; it is not a claim of full compliance.

The service embeds timezone records and web assets, provides versioned REST lookup endpoints with response envelopes and pagination, a server-rendered homepage and server pages, theme preferences, and the `timezones-cli` binary. It also serves a GraphQL endpoint and explorer, a bundled Swagger UI, an OpenAPI JSON description, current-time views, and robots, sitemap, Atom, RSS, security.txt, and llms.txt discovery resources.

The implementation remains incomplete against the full AI.md baseline. Verified remaining work is tracked in [TODO.AI.md](TODO.AI.md). Major gaps include the broad server configuration and lifecycle baseline, authentication/MCP, the required server/client CLI behavior, localization, security disclosure/contact workflows, and platform integrations. The required security and contact handlers and workflows remain unimplemented. IDEA.md follows AI.md's public-read and protected-operation access model.

Two inputs are awaiting the coordinating session: the exact stable `owner_token` cookie/localStorage name required by AI.md before token UI implementation, and the project's supported-release policy for `.github/SECURITY.md`. The timezone dataset's provenance and redistribution terms are now verified against an exact source-file match at the pinned upstream commit and recorded in `LICENSE.md`.

Previously recorded verification, before the current uncommitted edits: `make test` passed with 70.1% coverage, `golangci-lint run` passed, `mkdocs build --strict` passed, and Docker smoke checks confirmed the server, CLI, rendered pages, static assets, sitemap, feeds, and homepage asset loading. These results do not verify the current edits. Full browser interaction, responsive/accessibility review, provider CI, release workflows, and the remaining requirements in TODO.AI.md are outstanding.
