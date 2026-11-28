# Development

Go builds and tests are run through the Docker toolchain configured by the Makefile. Do not install or invoke Go directly on the local system.

```sh
make build
make test
./tests/run_tests.sh
./tests/e2e.sh
./tests/incus.sh
make release
```

`make build` creates debug binaries for the host and all eight distribution targets. `make test` runs Go vet, staticcheck, vulnerability and shell checks, unit tests, and the coverage threshold, then runs the basic binary checks. `make release` creates optimized production binaries, a source archive, and checksums without tagging or publishing. `./tests/run_tests.sh` builds binaries and checks their CLI behavior, API routes, rendered pages, and content types inside Docker. `./tests/docker.sh` reuses executable server and CLI binaries supplied via `BINDIR`, or builds them with the Docker Go toolchain otherwise, then runs them in an isolated Alpine container. It checks binary metadata and CLI basics, starts the server with `--debug`, verifies private permissions on generated secret storage, exercises operator, primary-owner, and linked-token authorization, confirms authenticated CLI timezone reads and token create/list/rename/revoke operations plus rejection of invalid credentials, checks debug-route authorization and redaction, and verifies Tor and I2P health fields and bootstrap state alongside public API and page routes. `TOR_WAIT` and `I2P_WAIT` set overlay startup ceilings (300 and 600 seconds by default); unavailable bootstrap is reported as skipped unless `REQUIRE_OVERLAY=1` is set. `./tests/e2e.sh` uses an isolated Docker network, the server runtime image, Chromium headless-shell, and the Go toolchain image to exercise server-rendered pages, browser interactions, JavaScript-disabled forms, localization, and responsive layouts. When Docker automatic address pools are exhausted, set `E2E_SUBNET` to an unused private CIDR (for example, `10.250.0.0/24`) so the script can create its isolated network explicitly. These Docker scripts require Docker, network access to pull the test images and Go modules, and available Docker resources.

`make docker` builds the current host platform locally, applies the required OCI metadata labels, and runs `tests/image.sh` against the resulting image. The image check confirms its required runtime commands and waits up to `TOR_WAIT` seconds (180 by default) for healthy Tor status. It reports an explicit skip when outbound networking is unavailable; set `REQUIRE_OVERLAY=1` to require Tor bootstrap.

`./tests/incus.sh [systemd|openrc|all]` runs the full application lifecycle in disposable Incus containers: Debian/Trixie with systemd and Alpine 3.24 with OpenRC. The default `all` mode runs both sequentially. It needs an initialized Incus host with the `images:` remote and network access for guest package installation; it also invokes `make build`, so Docker is required. `INCUS_CPUS`, `INCUS_MEM`, and `INCUS_CPU_PRIORITY` set container resource limits. `INCUS_OP_TIMEOUT`, `INCUS_INSTALL_TIMEOUT`, and `INCUS_RUN_TIMEOUT` set command, package-install, and full-run ceilings; `DOCKER_TIMEOUT` bounds `make build`. `REQUIRE_OVERLAY=1` makes the run require Tor and I2P checks, with `TOR_WAIT` and `I2P_WAIT` controlling their startup ceilings.

The server implementation is under `src/server/`; templates and static assets are embedded from `src/server/templates/` and `src/server/static/`. Timezone lookup logic lives in `src/timezones/`, and the dataset is in `src/data/`.

When changing routes or user-facing behavior, update the README, this documentation set, and the relevant server-rendered pages together. Do not include credentials or private keys in changes.
