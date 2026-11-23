# Implementation backlog

This list records verified gaps between `AI.md`, `IDEA.md`, and the current implementation. Items are removed only after the resolving change is committed, per the AI.md TODO rules.

## Required product surface

- [ ] Implement the IDEA-listed GraphQL endpoint and GraphiQL page; document and align its schema with REST.
- [ ] Verify the new current-time API and rendered page, then resolve this item in the follow-up commit.
- [ ] Provide a real Swagger UI at `/server/docs/swagger`; the current page only links to the OpenAPI JSON document.
- [ ] Complete content negotiation for browser, text browser, HTTP tool, and API clients across all public routes.
- [ ] Align remaining public API routes, response envelopes, error bodies, pagination/filtering, and health data with AI.md PARTS 9, 13, and 14.

## Server baseline

- [ ] Replace the partial legacy configuration schema with the required file-only `server.yml` configuration, validation, defaults, reload behavior, and documented privacy/branding/page settings.
- [ ] Complete the AI.md server baseline for TLS, service lifecycle, logging, health details, scheduler, backup/restore, metrics, and optional Tor/I2P behavior.
- [ ] Complete web security, SEO metadata, PWA offline behavior, cache/version invalidation, accessibility, and configured page-content behavior from AI.md PARTS 9, 11, 13, and 16.
- [ ] Implement the required multilingual catalogs and parity checks, or reconcile the English-only product scope with the required AI.md I18N baseline.
- [ ] Provide the required server/client CLI parity, configuration, completion, and interactive client behavior in AI.md PARTS 8 and 32.

## Build, governance, and verification

- [ ] Complete the required container, Makefile, CI/release, supply-chain, community, and multi-provider files against AI.md PARTS 25–28 and GitHub requirements.
- [ ] Add the required Docker, Incus, endpoint-coverage, browser E2E, and coverage-gate test workflows; tests were not run for this task.
- [ ] Verify and document timezone dataset provenance and redistribution terms in `LICENSE.md`.
- [ ] Verify the new pages, discovery endpoints, OpenAPI output, documentation links, and responsive/accessibility behavior in the primary session.
