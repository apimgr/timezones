# Implementation backlog

This list records verified gaps between `AI.md`, `IDEA.md`, and the current implementation. Items are removed only after the resolving change is committed, per the AI.md TODO rules.

## Required product surface

- [x] Implement the IDEA-listed GraphQL endpoint and browser query workbench; document and align its schema with REST. (Implemented in the working tree; retain until committed.)
- [x] Verify the current-time API and rendered page. (Covered by `make test`; retain until committed.)
- [x] Provide a bundled Swagger UI at `/server/docs/swagger`. (Implemented in the working tree; retain until committed.)
- [x] Add production, development, and automated-test Compose files under `docker/`; remove the obsolete root Compose files. (Implemented in the working tree; retain until committed.)
- [x] Add the Docker binary/API integration runner and make its build output mount usable with a project-selected `BINDIR`. (Implemented in the working tree; the runner still needs full AI.md route/auth/overlay coverage; retain until committed.)
- [x] Align Makefile build metadata and enforce the 60% Go coverage floor. (Implemented in the working tree; the latest recorded `make test` result was 70.1%; rerun after the current edits and retain until committed.)
- [x] Add the GitHub community policy files required for this public repository. (Implemented in the working tree; retain until committed.)
- [ ] Complete content negotiation for browser, text browser, HTTP tool, and API clients across all public routes.
- [x] Add browser pages for list/search/random/filter/lookup timezone routes and dataset statistics/count. (Implemented in the working tree; retain until committed.)
- [ ] Complete API alignment for all routes and detailed health data under AI.md PARTS 9, 13, and 14. Standard success/error envelopes, list/search pagination and filtering, pagination overflow handling, forwarded-host allow-list enforcement, `/timezones.json` OpenAPI schema, and browser/text/API content negotiation are implemented in the working tree; retain those fixes until committed.
- [ ] Implement operator and owner tokens, token-management browser/API routes, protected metrics and MCP endpoints, SQLite persistence, and their full security and lifecycle requirements from AI.md PARTS 5, 8, 10–12, 14, 20, 28, and 32.

## Server baseline

- [ ] Replace the partial legacy configuration schema with the required file-only `server.yml` configuration, validation, defaults, reload behavior, and documented privacy/branding/page settings.
- [ ] Complete the AI.md server baseline for TLS, service lifecycle, logging, health details, scheduler, backup/restore, metrics, and optional Tor/I2P behavior.
- [ ] Complete web security, SEO metadata, cache/version invalidation, accessibility, and configured page-content behavior from AI.md PARTS 9, 11, 13, and 16. The static asset mount, service-worker registration, offline navigation fallback, HTTP-tool text rendering, and no-JavaScript text-browser response are implemented and covered by tests; retain these changes until committed.
- [ ] Complete every required server-rendered page and synchronize page content, README, and all current documentation with implemented behavior and AI.md; verify all documentation links and route inventories.
- [ ] Implement the AI.md PART 11 disclosure and contact system: contact and DPO pages/API, CAPTCHA and review/confirmation flow, encrypted security submissions and researcher status links, operator maintenance commands, PGP lifecycle, hall of fame, security policy, security.txt and public key routes. The route inventory in IDEA.md now records the required public surfaces; none of these handlers or persistence workflows is implemented.
- [ ] Implement the required multilingual catalogs and parity checks, or reconcile the English-only product scope with the required AI.md I18N baseline.
- [ ] Provide the required server/client CLI parity, configuration, completion, and interactive client behavior in AI.md PARTS 8 and 32.

## Build, governance, and verification

- [ ] Complete the required container, Makefile, CI/release, supply-chain, community, and multi-provider files against AI.md PARTS 25–28 and GitHub requirements.
- [ ] Complete required Docker, Incus, endpoint-coverage, browser E2E, and provider-workflow verification. `make test` passes the coverage gate; browser E2E and provider workflow verification remain outstanding.
- [ ] Implement `tests/incus.sh` full-app service lifecycle and `tests/e2e.sh` browser E2E suites, then verify them in isolated environments.
- [ ] Verify and document timezone dataset provenance and redistribution terms in `LICENSE.md`.
- [ ] Complete browser interaction, responsive/accessibility, and OpenAPI UI verification. The built server passed a container HTTP smoke check, and headless Chromium loaded the homepage and its local assets; the full visual/accessibility review remains outstanding.
