# Contributing

## Local setup and checks

Go builds and tests use the Docker toolchain configured by the Makefile. Do not install Go on the host for this project.

```sh
make dev
make test
./tests/run_tests.sh
docker run --rm --name "timezones-docs-$(tr -dc 'a-z0-9' </dev/urandom | head -c8)" -v "$PWD:/docs" -w /docs python:3.12-slim sh -c 'pip install -r docs/requirements.txt && mkdocs build --strict -d /tmp/timezones-docs/site'
```

The integration script requires Docker and exercises the server and CLI binaries. MkDocs requirements are listed in `docs/requirements.txt`.

## Changes and pull requests

- Work on a branch and open a pull request against the repository's default branch.
- Include tests for behavior changes and update the README and documentation when routes or user-facing behavior change.
- Describe the change, its reason, test evidence, configuration or documentation updates, and any breaking changes.
- State any security or privacy impact.
- Do not add placeholders, stubs, or unfinished TODO behavior.

## Security reports

Do not report vulnerabilities in public issues. Use the private reporting path in [SECURITY.md](SECURITY.md).
