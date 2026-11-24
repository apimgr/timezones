# Project SPEC

Project: timezones
Role: Efficient loader for AI.md

⚠️ **THIS FILE IS AUTO-LOADED EVERY CONVERSATION. FOLLOW IT EXACTLY.** ⚠️

`AI.md` is the full, read-only implementation specification. `IDEA.md` defines the product and canonical project variables. `SPEC.md`, when present, contains project-specific rule overrides and takes precedence over AI.md.

## Before ANY Code Change

1. Read the complete AI.md PART relevant to the change and follow its references.
2. Search the repository for reusable code and read every file before editing it.
3. Do not guess. Ask when a required product decision is missing or the spec is ambiguous.
4. Keep IDEA.md within the WHAT boundary; never change its project variables without user confirmation.
5. Verify the implementation and documentation against AI.md and its FINAL CHECKPOINT before reporting completion.

## CRITICAL - NEVER DO

- Never edit AI.md or silently override it in code, IDEA.md, or documentation.
- Never use Go on the host; build, test, and run binaries in the prescribed Docker or Incus environments.
- Never execute a built binary directly on the host or write test/runtime data into the repository.
- Never use CGO, client-side rendering, unsafe input handling, external cron, or unapproved dependencies.
- Never add stubs, TODO/FIXME/HACK comments, placeholder content, or undocumented behavior.
- Never use Makefile targets in CI/CD; workflows use explicit commands.

## CRITICAL - ALWAYS DO

- Implement every mandatory AI.md requirement; optional features are selected by IDEA.md or SPEC.md as specified.
- Use server-rendered Go templates, progressive enhancement, mobile-first accessible UI, and the required companion client.
- Build the static binaries for all eight OS/architecture combinations with the project Docker toolchain.
- Keep README.md, docs/, OpenAPI, GraphQL, and actual behavior synchronized.
- Run the required test and lint gates through the prescribed wrapper; report actual results only.
- Use `gitcommit --dir {project_dir} all` for commits and follow `.git/COMMIT_MESS` requirements.

## Rule Files

Read the focused rules before work: `.claude/rules/ai-rules.md`, `project-rules.md`, `config-rules.md`, `binary-rules.md`, `backend-rules.md`, `api-rules.md`, `frontend-rules.md`, `features-rules.md`, `service-rules.md`, `makefile-rules.md`, `docker-rules.md`, `cicd-rules.md`, and `testing-rules.md`.

## Binary Terminology

- Server: `timezones`
- Required companion client: `timezones-cli`

## Key Paths

- Source: `src/`
- Tests: `tests/` and Go tests beside source
- Docker: `docker/`
- Documentation: `docs/`
- Full source of truth: `AI.md`; current tasks: `TODO.AI.md`
