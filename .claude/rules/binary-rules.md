# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use Go on the host or run a built binary on the host; build and execute inside the prescribed Docker or Incus environments.
- Never enable CGO or add an external runtime dependency beyond the specified optional Tor and I2P daemons.
- Never omit the companion client or invent unsupported flags, output formats, or configuration behavior.
- Never write binaries or runtime files into the source tree except configured build outputs.

## CRITICAL - ALWAYS DO

- Embed required assets and produce static server/client binaries for all eight OS/architecture combinations.
- Implement the server CLI, service, maintenance, status, update, path, and privilege behavior in PART 8.
- Keep `timezones-cli` in parity with the server API and implement its profile, output, language, paging, shell, TUI, GUI, and update behavior from PART 32.
- Preserve `NO_COLOR`, structured errors, safe token handling, and exact documented output contracts.

## Key Rules

- Server binary: `timezones`; required companion client: `timezones-cli`.
- CLI output formats are `table`, `json`, and `plain` where specified; use only flags defined in AI.md.
- Run all Go build, test, and binary execution steps through the prescribed container workflow.

For complete details, see AI.md PART 7, 8, 32.
