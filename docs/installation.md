# Installation

## Docker Compose

The production image listens on port 80 and includes Tor and I2P packages. The
provided Compose setup starts the API with Valkey, stores configuration and
data under `docker/volumes/`, and binds the host port to the Docker bridge
gateway.

The container stores configuration in `/config/timezones`, application data
and cache in `/data/timezones`, logs in `/data/log/timezones`, SQLite data in
`/data/db/sqlite`, and backups in `/data/backups/timezones`. These paths map to
the Compose configuration and data mounts.

```bash
docker compose -f docker/docker-compose.yml up -d
```

The container starts a managed Tor hidden service when the `tor` binary is
available. I2P defaults to enabled in application containers. Override it with
`I2P_ENABLED=true|false` or `--i2p[=true|false]`. Discovery publishes either
overlay address only when its `expose` setting is enabled.

## Binary

Download the matching server release binary for your operating system and
architecture, then run it in the foreground:

```sh
# Start a local test server on loopback.
./timezones --address 127.0.0.1 --port 8080
```

The server and `timezones-cli` are separate binaries. The server embeds the
timezone dataset and web assets. Optional desktop builds are published as
`timezones-cli-gui-{os}-{arch}` for Linux, macOS, and Windows. These builds use
the platform's dynamic loader and are intended for desktop sessions; standard
CLI builds remain static and portable.

## First start

On first start, the server creates a default `server.yml` in its configuration
directory. Use `CONFIG_DIR` or `--config` to select a different configuration
directory, and `DATA_DIR` or `--data` to select runtime storage. `PORT` or
`--port` selects the listener port; `LISTEN` or `--address` selects its
address. With no port set, the server chooses an unused port on a host and
uses port 80 in an application container. The server runs in the foreground by
default.

## systemd service

On a Linux host managed by systemd, install and start the service:

```bash
sudo timezones --service --install
```

The installer detects the platform service manager, enables the service, and
starts it. To remove it, run `sudo timezones --service --uninstall` and confirm
the prompt; this deletes application data and removes the Linux service user or
group only when the application recorded that it created them. The server binary
remains and must be deleted manually. See [Configuration](configuration.md)
for all runtime settings.
