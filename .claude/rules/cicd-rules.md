# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never invoke Makefile targets from CI workflows.
- Never use mutable/unpinned third-party actions or leak secrets in logs/artifacts.
- Never skip test, lint, security, supply-chain, release, or provider-specific requirements.
- Never edit provider workflow files without validating their native syntax and pin rules.

## CRITICAL - ALWAYS DO

- Keep GitHub, GitLab, Gitea, Forgejo, and Jenkins workflows aligned with the required supported-provider set.
- Use explicit Docker toolchain commands in CI, including the same prescribed test and lint gates as local verification.
- Pin GitHub actions to verified full commit SHAs and maintain the documented SHA table.
- Apply least-privilege permissions, concurrency controls, artifact retention, and secret handling from AI.md.
- Validate each workflow with its provider's native validator before completion.

## Key Rules

- Workflow changes must retain required branch triggers, build matrix, tests, security scans, and release conditions.
- Keep dependency automation and community/security metadata synchronized with repository behavior.

For complete details, see AI.md PART 27.
