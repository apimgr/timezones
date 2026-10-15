# Timezones API implementation backlog

## Completed in this pass

- [x] Added a server-rendered timezone finder to the homepage. Search by
      display name, abbreviation, description, and IANA identifier works with
      a normal GET form and does not require JavaScript.
- [x] Added query length validation and global request target/body limits.
- [x] Added bounded, sharded per-client request rate limiting with `429`,
      `Retry-After`, and the shared API error shape.
- [x] Changed `--status` and the Docker health check to probe `/healthz` rather
      than reporting healthy unconditionally.
- [x] Return HTTP listener bind failures as a nonzero process exit and use the
      shared 30-second read/write and 120-second idle timeout defaults.
- [x] Added the IDEA-required `timezones-cli` lookup command and included its
      eight-platform release artifacts.
- [x] Use `casjaysdev/go:latest` for Docker and Makefile Go builds, add
      `-buildvcs=false` and `-trimpath`, and point local image builds at the
      committed Dockerfile path.
- [x] Replaced the invalid text-in-a-PNG favicon and missing PWA icon references
      with embedded SVG icons.
- [x] Removed unsupported admin-panel and credential instructions from the
      README. Recent GitHub CI, Docker Build, and Daily Build runs are green.

## Remaining IDEA.md product requirements

- [ ] Add tests for theme cookie rendering, all three cycle transitions, CSRF
      rejection, and safe return redirects.

- [ ] Implement `/graphql` with a real GraphQL schema and resolvers. It is
      advertised in IDEA.md but no GraphQL endpoint is present.
- [ ] Implement OpenAPI/Swagger JSON and the documented Swagger UI routes.
- [ ] Add the required `/server/about`, `/server/help`, `/server/healthz`,
      `/server/privacy`, and `/server/terms` pages. Only `/healthz` JSON is
      implemented; the preference POST is limited to the theme toggle.
- [ ] Reconcile the IDEA conflict between “current time in any timezone” and
      the non-goal of real-time DST calculation. The bundled records are static;
      using them to produce current local time would be wrong during DST changes.
- [ ] Identify, document, and attribute the embedded dataset's authoritative
      source, version, and redistribution license. Current repository content
      does not establish these details, so full-world completeness cannot be
      claimed safely.
- [ ] Add persistent/configurable rate-limit policy under `server.rate_limit.*`;
      the implemented default currently uses a fixed 240 requests per minute
      per client address.

## Remaining shared AI.md requirements

- [ ] Align the Docker Tor integration with AI.md PART 31: Tor is started by an
      entrypoint toggle and generic torrc today, instead of being app-managed
      with a dedicated loopback backend and the required isolation controls.
- [ ] Add GitLab, Gitea, Forgejo, and Jenkins equivalents with the same test,
      lint, security, release, and artifact gates. This repository currently
      has GitHub workflows only.
- [ ] Add CODEOWNERS, security policy, issue/PR templates, and branch-protection
      documentation or settings. Renovate configuration is already present.
- [ ] Complete the CLI contract documented by AI.md, including service/config
      lifecycle parity, configurable listen URL/base path, data/log directory
      behavior, and the required maintenance/update command semantics.
- [ ] Add integration coverage for container startup, health checks, request
      limits, PWA assets, API endpoints, and CLI behavior. The primary session
      ran `make test` successfully after the current changes; this does not
      replace those binary and endpoint integration checks.
