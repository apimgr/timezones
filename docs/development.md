# Development

Go builds and tests are run through the Docker toolchain configured by the Makefile. Do not install or invoke Go directly on the local system.

```sh
make dev
make test
./tests/run_tests.sh
./tests/e2e.sh
make build
./tests/incus.sh
```

`make dev` builds the current platform into a temporary directory. `make test` runs Go vet, vulnerability checks, unit tests, and the coverage threshold. `./tests/run_tests.sh` builds the binaries and checks their CLI behavior, API routes, rendered pages, and content types inside Docker. `./tests/docker.sh` reuses executable server and CLI binaries supplied via `BINDIR`, or builds them with the Docker Go toolchain otherwise, then runs them in an isolated Alpine container. It checks binary metadata and CLI basics, starts the server with `--debug`, verifies private permissions on generated secret storage, exercises operator, primary-owner, and linked-token authorization, confirms authenticated CLI timezone reads and token create/list/rename/revoke operations plus rejection of invalid credentials, checks debug-route authorization and redaction, and verifies Tor and I2P health fields and bootstrap state alongside public API and page routes. `TOR_WAIT` and `I2P_WAIT` set overlay startup ceilings (300 and 600 seconds by default); unavailable bootstrap is reported as skipped unless `REQUIRE_OVERLAY=1` is set. `./tests/e2e.sh` uses an isolated Docker network, the server runtime image, Chromium headless-shell, and the Go toolchain image to exercise server-rendered pages, browser interactions, JavaScript-disabled forms, localization, and responsive layouts. These Docker scripts require Docker, network access to pull the test images and Go modules, and available Docker resources. `make build` builds release binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64.

`./tests/incus.sh [systemd|openrc|all]` runs the full application lifecycle in disposable Incus containers: Debian/Trixie with systemd and Alpine 3.24 with OpenRC. The default `all` mode runs both sequentially. It needs an initialized Incus host with the `images:` remote and network access for guest package installation; it also invokes `make build`, so Docker is required. `INCUS_CPUS`, `INCUS_MEM`, and `INCUS_CPU_PRIORITY` set container resource limits. `INCUS_OP_TIMEOUT`, `INCUS_INSTALL_TIMEOUT`, and `INCUS_RUN_TIMEOUT` set command, package-install, and full-run ceilings; `DOCKER_TIMEOUT` bounds `make build`. `REQUIRE_OVERLAY=1` makes the run require Tor and I2P checks, with `TOR_WAIT` and `I2P_WAIT` controlling their startup ceilings.

The server implementation is under `src/server/`; templates and static assets are embedded from `src/server/templates/` and `src/server/static/`. Timezone lookup logic lives in `src/timezones/`, and the dataset is in `src/data/`.

When changing routes or user-facing behavior, update the README, this documentation set, and the relevant server-rendered pages together. Do not include credentials or private keys in changes.
