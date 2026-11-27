# Testing Rules (PART 28, 29, 30)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never run Go, build, tests, lint, formatters, or binaries on the host; use Docker or the prescribed Incus VM/container workflow.
- Never write test/runtime artifacts into the repository or use unscoped temp paths.
- Never treat a passing unit suite as proof that browser, service, overlay, CI, or end-to-end requirements are complete.
- Never ship stale docs, missing locale strings, mismatched catalogs, or untranslated user-visible page text.

## CRITICAL - ALWAYS DO

- Run test and direct lint commands through `test-lint-run.sh` and report actual output.
- Keep unit, Docker binary/API integration, basic container, Incus systemd/OpenRC, browser E2E, and CI workflow checks distinct and complete.
- Build ReadTheDocs with the specified strict MkDocs configuration and keep README/docs/API/OpenAPI synchronized.
- Maintain locale parity, language negotiation, RTL layout, accessible text, and progressive enhancement.
- Verify changed behavior against the FINAL CHECKPOINT and required client classes/endpoints.

## Key Rules

- Use project-scoped temp directories for all logs, caches, browser profiles, coverage, and container volumes.
- Primary agent alone runs verification and commit gates; report commands and results accurately.

For complete details, see AI.md PART 28, 29, 30.
