# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never assume a service manager, init system, or privilege-escalation tool exists.
- Never perform broad host cleanup or change global host configuration to work around a project issue.
- Never leave stale pid files, service files, privilege state, or orphan processes after lifecycle operations.
- Never run the daemon as root when the specified privilege drop is available.

## CRITICAL - ALWAYS DO

- Detect the supported service manager and operating environment using the prescribed precedence.
- Make install, uninstall, start, stop, restart, reload, enable, disable, and status operations idempotent and bounded.
- Keep privilege escalation explicit, narrow, validated, and auditable.
- Run service lifecycle verification in the required disposable Incus OS environments.

## Key Rules

- Keep service files and user/group behavior platform-specific and aligned with `--service` CLI semantics.
- Clean up only resources created for the task; preserve unrelated containers, services, and user changes.

For complete details, see AI.md PART 23, 24.
