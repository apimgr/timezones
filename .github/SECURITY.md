# Security Policy

## Supported versions

Only the latest published release is supported and receives security fixes. Older releases are unsupported; fixes are not backported.

## Reporting a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/apimgr/timezones/security/advisories/new) as the primary reporting channel. Do not file vulnerabilities in public issues. On repository mirrors without private vulnerability reporting, report through the GitHub origin repository.

Maintainers will acknowledge reports within 48 hours and coordinate remediation and disclosure privately with the reporter. Critical issues target a fix within 14 calendar days. Other issues are prioritized by severity and may take longer. Public disclosure follows the coordinated-disclosure policy after a fix is available, normally within 90 days of the report unless the reporter and maintainers agree otherwise.

For deployed-instance reporting details, see the `/.well-known/security.txt` endpoint provided by the deployment and the [security documentation](https://apimgr-timezones.readthedocs.io/en/latest/security/).

The deployment also provides a security overview at `/server/security` and a deployed-instance reporting form at `/server/contact?security_id=...`. Use the form's security-report option and complete its review and confirmation steps. Do not include vulnerability details in a public issue or ordinary contact request.
