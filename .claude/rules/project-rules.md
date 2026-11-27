# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never add `cmd/`, `internal/`, or `pkg/`; project Go source belongs under `src/`.
- Never put generated config, runtime state, test output, logs, caches, or secrets in the repository.
- Never add a dependency without checking the approved library list and recording required license attribution.
- Never change `IDEA.md` project variables without user confirmation.

## CRITICAL - ALWAYS DO

- Keep the required `src/`, `docker/`, `docs/`, and `tests/` structure and package responsibilities.
- Resolve paths from the frozen `internal_org` and `internal_name` project variables.
- Support Linux, macOS, Windows, BSD, user, root, and container path conventions defined in PART 4.
- Keep `LICENSE.md` and third-party attribution synchronized with embedded data and dependencies.

## Key Rules

- Project name and binary terms come from IDEA.md; server binary is `timezones`, companion client is `timezones-cli`.
- Keep source and documentation in the repository; execution artifacts go to a project-scoped OS temp directory.
- Use the prescribed OS-specific config, data, cache, and log roots; container paths are `/config/` and `/data/`.

For complete details, see AI.md PART 2, 3, 4.
