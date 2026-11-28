# CLI reference

`timezones-cli` is the companion client for a Timezones API server. It supports
one-shot queries, interactive terminal and desktop interfaces, output
formatting, caching, and owner-token commands. With no query, it opens the
configured interactive display.

## Setup and authentication

The client reads its `cli.yml` profile from the platform's user configuration
directory. Select another profile with `--config`. The server URL can be set in
the profile, with `TIMEZONES_SERVER_PRIMARY`, or with `--server`. The built-in
official server URL is the final fallback when one is configured into the
release.

Supply an API token using `--token`, `--token-file`, `TIMEZONES_TOKEN`,
`TIMEZONES_AUTH_TOKEN`, or the profile's `auth.token` / `auth.token_file`
settings. Avoid command-line tokens in shared terminal sessions because they
may be visible in process listings or shell history. The client sends tokens
only in the `Authorization: Bearer` header.

On first run without a configured server, the client opens the detected GUI or
terminal setup screen when available. Enter a server URL and optionally a
token; the client checks public connectivity without sending the token, then
validates the token with the machine-only auth-check endpoint. A rejected token
is not saved. After a successful connection the screen shows the server name
and version, and the URL and accepted token are saved to the CLI profile.

```sh
timezones-cli --server https://YOUR_SERVER search America/New_York
timezones-cli --config ~/.config/timezones/cli.yml --token-file ~/.config/timezones/token list
timezones-cli --help
timezones-cli --version
```

For a local server started on port 8080, use `http://127.0.0.1:8080` as the
`--server` value.

## Commands

| Command | Purpose |
|---|---|
| `timezones-cli [query]` | Search using a display name, abbreviation, or IANA identifier |
| `list` | List records; accepts `--page`, `--utc-offset`, `--isdst`, and `--limit` |
| `search [query]` | Search timezone names and identifiers |
| `get [name]` | Look up a display name |
| `current [identifier]` | Get current local time for an IANA identifier |
| `utc [identifier]` | Look up an IANA identifier |
| `offset [hours]` | List records at a UTC offset |
| `abbr [abbreviation]` | List records matching an abbreviation |
| `value [name]` | Look up a timezone by display name |
| `random` | Get a random timezone record |
| `stats` | Show dataset statistics |
| `count` | Show the number of records |
| `token list` | List token metadata |
| `token issue` | Claim a primary owner token and save it to the CLI profile |
| `token create --name NAME [--expires DURATION]` | Create a linked token |
| `token reveal [id]` / `token rotate [id]` | Reveal or rotate the primary token, or the selected linked token |
| `token rename ID NAME` | Rename a linked token |
| `token revoke ID\|primary [--confirm TEXT]` | Revoke a linked token or confirmed primary token |

Use `--output table`, `json`, or `plain` to select the output format.
`--lang` selects a supported language. `--color auto|yes|no` controls
color output. `--limit` controls the result count for list and search commands.

```sh
timezones-cli --server https://YOUR_SERVER search --output json 'America/New_York'
timezones-cli --server https://YOUR_SERVER current America/New_York --output json
timezones-cli --server https://YOUR_SERVER list --utc-offset -5 --isdst false --page 2 --limit 50
timezones-cli --server https://YOUR_SERVER --lang es --output plain EST
```

Shell integration uses `--shell completions [SHELL]` to print a completion
script and `--shell init [SHELL]` to print an eval-ready command. If `SHELL` is
omitted, the client detects it from `$SHELL`. Completion scripts are available
for bash, zsh, fish, PowerShell, and PowerShell Core (`pwsh`). `sh`, `dash`, and
`ksh` have no programmable completion, so requesting completions for those
shells prints usage and exits with status 2; their `init` action prints a safe
no-op comment. Use `--shell --help` for detailed help.

The client profile includes server timeout/retry, auth, output,
theme, display, cache, logging, update, language, and default result-limit
settings. The profile is created with owner-only permissions on Unix-like
systems.

## Server status and overlay control

The server binary accepts `--status` to query health and print Tor/I2P
lifecycle values. Tor starts when its binary is found. I2P defaults to enabled
inside an application container and disabled elsewhere; use `--i2p` or
`--i2p=false` to override that default. `I2P_ENABLED=true|false` and
`server.i2p.enabled` provide the equivalent environment and configuration
controls. The server's help output lists these options.

## Server maintenance

The server binary also provides operator maintenance commands. Backups use the
configured backup password when present; encrypted restore and verification
prompt for the password separately from the server token. Keep the backup
password in `secrets.yml` and provide it interactively when restoring.

```sh
timezones --maintenance backup [filename] [--include-ssl] [--include-data]
timezones --maintenance restore /path/to/backup.tar.gz.enc
timezones --maintenance backup test [/path/to/backup.tar.gz.enc]
timezones --maintenance scheduler list
timezones --maintenance scheduler show backup_daily
timezones --maintenance scheduler run backup_daily
timezones --maintenance secret rotate encryption_key
```

At-rest key rotation requires root or the operator token, a typed confirmation,
and a second operator-token prompt when token authorization is used. Rotation
stages `secrets.yml.new`, keeps the previous key wrapped in the database for
the 30-day recovery window, and re-encrypts token, security-report, DNS, and
backup-password ciphertext before the server accepts requests after recovery.

Restore is destructive. It checks authorization and validates the archive,
manifest, extracted paths, and database before applying the backup. Scheduled
daily backups create a verified full archive and a cumulative incremental
when a compatible previous full archive is available. Hourly incrementals are
disabled by default. Retention lives under `server.backup.retention`; see
[Configuration](configuration.md#security-notifications-contact-and-language).
