# Implementation status

`AI.md` is the implementation source of truth and `IDEA.md` defines the product scope. This status summarizes the current tree; it is not a claim of full compliance.

The service embeds timezone records and web assets, provides versioned REST lookup endpoints, a server-rendered homepage, theme preferences, and the `timezones-cli` binary. It also serves public about, help, health, privacy, terms, and API documentation pages, paired JSON page endpoints, an OpenAPI JSON description, and basic robots, sitemap, security.txt, and llms.txt discovery resources.

The implementation remains incomplete against the full AI.md baseline. Verified remaining work is tracked in [TODO.AI.md](TODO.AI.md). In particular, GraphQL, full CLI and server baseline requirements, complete content negotiation, localization, and integration verification are still outstanding. The current-time API and page were added in this pass and still need execution-level verification.

No build, test, lint, formatter, or workflow validation was run during the implementation pass.
