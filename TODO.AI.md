# Implementation backlog

This list records verified gaps between `AI.md`, `IDEA.md`, and the current implementation. Items are removed only after the resolving change is committed, per the AI.md TODO rules.

## Required product surface

- [x] Implement the IDEA-listed GraphQL endpoint and browser query workbench; document and align its schema with REST. (Implemented in the working tree; retain until committed.)
- [x] Verify the current-time API and rendered page. (Covered by `make test`; retain until committed.)
- [x] Provide a bundled Swagger UI at `/server/docs/swagger`. (Implemented in the working tree; retain until committed.)
- [ ] Complete content negotiation for browser, text browser, HTTP tool, and API clients across all public routes.
- [x] Add browser pages for list/search/random/filter/lookup timezone routes and dataset statistics/count. (Implemented in the working tree; retain until committed.)
- [ ] Complete API alignment for all routes and detailed health data under AI.md PARTS 9, 13, and 14. Standard success/error envelopes, list/search pagination and filtering, pagination overflow handling, forwarded-host allow-list enforcement, and the `/timezones.json` OpenAPI response schema are implemented in the working tree; retain those fixes until committed.

## Server baseline

- [ ] Replace the partial legacy configuration schema with the required file-only `server.yml` configuration, validation, defaults, reload behavior, and documented privacy/branding/page settings.
- [ ] Complete the AI.md server baseline for TLS, service lifecycle, logging, health details, scheduler, backup/restore, metrics, and optional Tor/I2P behavior.
- [ ] Complete web security, SEO metadata, cache/version invalidation, accessibility, and configured page-content behavior from AI.md PARTS 9, 11, 13, and 16. The static asset mount, service-worker registration, and offline navigation fallback are implemented and covered by tests; retain these changes until committed.
- [ ] Implement the required multilingual catalogs and parity checks, or reconcile the English-only product scope with the required AI.md I18N baseline.
- [ ] Provide the required server/client CLI parity, configuration, completion, and interactive client behavior in AI.md PARTS 8 and 32.

## Build, governance, and verification

- [ ] Complete the required container, Makefile, CI/release, supply-chain, community, and multi-provider files against AI.md PARTS 25–28 and GitHub requirements.
- [ ] Complete required Docker, Incus, endpoint-coverage, browser E2E, and coverage-gate verification. `make test` now passes; browser E2E and external-provider workflow verification remain outstanding.
- [ ] Verify and document timezone dataset provenance and redistribution terms in `LICENSE.md`.
- [ ] Complete browser interaction, responsive/accessibility, and OpenAPI UI verification. The built server passed a container HTTP smoke check, and headless Chromium loaded the homepage and its local assets; the full visual/accessibility review remains outstanding.
