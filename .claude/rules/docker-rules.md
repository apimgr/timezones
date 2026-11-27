# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never put runtime volumes, generated config, secrets, or test state in the repository.
- Never run Compose workflows from the source tree when the prescribed temp-directory workflow applies.
- Never use an unpinned or privileged runtime workaround contrary to the project Docker specification.
- Never rely on an init script for application initialization that belongs in the binary.

## CRITICAL - ALWAYS DO

- Use the specified Docker base/toolchain images, static builds, embedded assets, and minimal runtime layout.
- Provide production, development, and test Compose definitions under `docker/` with correct config/data/log/cache/backup mounts.
- Run as the required non-root identity with the specified permissions and signal handling.
- Keep health checks, overlay daemon handling, entrypoint behavior, image metadata, and resource settings aligned with PART 26.

## Key Rules

- Runtime configuration and data belong in `/config/` and `/data/`; test volumes live in project-scoped OS temp directories.
- Do not use Docker image execution as a replacement for the required compiled-binary integration suite.

For complete details, see AI.md PART 26.
