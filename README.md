# Timezones API

Timezones is a static Go service with an embedded world timezone reference dataset, a versioned REST API, and a server rendered lookup page. The API and web interface are read-only and public. The timezone offsets and DST flags describe the bundled static dataset; they are not calculated for the current instant.

## Features

- Search the bundled records by display name, abbreviation, description, or IANA identifier.
- Filter by UTC offset and abbreviation, and look up one IANA identifier or display name.
- View dataset counts and retrieve records as JSON or plain text.
- Browse records from the homepage with a regular HTML form; JavaScript is not required for search.
- Use the standalone `timezones-cli` lookup binary.
- Choose a dark, light, or system-matched theme. Theme preference is saved in a server-readable cookie and the control also works without JavaScript.
- Run the HTTP service or the multi-architecture container image.

## Run the server

```sh
./timezones --address 127.0.0.1 --port 8080
```

Open `http://127.0.0.1:8080/` to search the dataset. The server creates a default `server.yml` in its configuration directory. Set `CONFIG_DIR`, `PORT`, or `ADDRESS` to override the corresponding location or listener setting.

## CLI lookup

```sh
timezones-cli America/New_York
timezones-cli --json EST
```

## REST API

| Method and path | Purpose |
|---|---|
| `GET /api/v1/timezones` | Return every record in the API response envelope |
| `GET /api/v1/timezones.json` | Return the embedded JSON dataset |
| `GET /api/v1/timezones/search?q=...` | Search by name, description, abbreviation, or IANA ID |
| `GET /api/v1/timezones/offset/{offset}` | Match the dataset's UTC offset in hours |
| `GET /api/v1/timezones/abbr/{abbr}` | Match an abbreviation |
| `GET /api/v1/timezones/utc/{iana-id}` | Look up an IANA identifier, including identifiers with slashes |
| `GET /api/v1/timezones/value/{name}` | Look up a display name |
| `GET /api/v1/timezones/random` | Return one random record |
| `GET /api/v1/stats` | Return dataset counts |
| `GET /healthz` | Check whether the HTTP service is responding |

A plain text variant is available for list, search, random, statistics, and count endpoints using the `.txt` suffix. Search query length is limited to 128 bytes. Requests are rate limited by client address.

Example:

```sh
curl 'http://127.0.0.1:8080/api/v1/timezones/search?q=America/New_York'
```

## Container

The production image listens on port 80. The container entrypoint starts Tor by default; set `ENABLE_TOR=false` to disable it. Bind persistent configuration and data directories as needed.

```sh
docker run --rm -p 8080:80 -e ENABLE_TOR=false ghcr.io/apimgr/timezones:latest
```

## Build and development

The project embeds the timezone dataset and web assets into the server binary. GitHub Actions builds the server and CLI release binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64. Docker images target amd64 and arm64.

## License

Project code is MIT licensed. The embedded dataset's source and redistribution terms still need a complete per-dataset attribution and provenance manifest; see [TODO.AI.md](TODO.AI.md).
