# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never call Makefile targets from CI/CD workflows; CI must use explicit commands.
- Never run Go build or test targets with a host Go toolchain.
- Never let build, clean, release, or test commands write artifacts into unapproved runtime paths.
- Never make release targets perform destructive tag/release operations outside their specified workflow.

## CRITICAL - ALWAYS DO

- Make targets must use the project Docker Go toolchain, static flags, project identity, and configured cache/output locations.
- Keep the six required build, test, release, Docker, cleanup, and help targets consistent with AI.md; do not add targets.
- Build all required server/client platform combinations and optional GUI variants exactly as specified.
- Keep the 60% coverage threshold and prescribed vet/vulnerability/test checks in the test target.

## Key Rules

- Use explicit build metadata and strip settings from the Makefile spec.
- Use task-scoped `BINDIR`, `RELDIR`, `GO_CACHE`, `GO_BUILD`, and temp paths during verification.

For complete details, see AI.md PART 25.
