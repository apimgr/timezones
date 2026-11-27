# Configuration reference

This reference lists the current defaults from the server configuration
loader and CLI profile initializer. A value shown as `{}` is an empty mapping;
`[]` is an empty list. Values supplied at runtime or generated from the
deployment are marked as dynamic. Server configuration changes are read from
`server.yml`, sorted `yaml.d/*.yml` drop-ins, `secrets.yml`, and configuration
assets; the loader watches these inputs and reloads settings. The listener,
database, TLS, daemon, and Tor settings identified by AI.md as restart-required
take effect after restart.

## Server files and precedence

The configuration directory contains `server.yml` for node-local listener,
identity, database, cache, and TLS settings; `yaml.d/{area}.yml` files for
area settings; and `secrets.yml` for credentials. Drop-ins are merged in
filename order, then secrets are merged. The loader applies built-in defaults
for omitted settings. Environment variables only override the settings and
phases listed below; they are not a general `PROJECT_SECTION_KEY` override
mechanism. Invalid drop-in entries are warned about and skipped.
String values also support environment interpolation as `${NAME}` or
`${NAME:-fallback}`; an unset variable without a fallback expands to empty.

## Server default inventory

### Node and logging

| Setting | Default |
|---|---|
| `server.token` | `""` |
| `server.port` | `""` (selected at startup) |
| `server.fqdn` | `""` |
| `server.address` | `[::]` |
| `server.mode` | `production` |
| `server.api_version` | `v1` |
| `server.baseurl` | `/` |
| `server.pidfile` | `true` |
| `server.daemonize` | `false` |
| `server.logs.level` | `info` |
| `server.logs.access.format` | `apache` |

### Database, cache, and TLS

| Setting | Default |
|---|---|
| `server.database.driver` | `sqlite` |
| `server.database.url` | `""` |
| `server.database.token` | `""` |
| `server.cache.type` | `memory` |
| `server.cache.url` | `""` |
| `server.cache.host` / `port` / `username` / `password` / `db` | `localhost` / `6379` / `""` / `""` / `0` |
| `server.cache.tls` / `tls_skip_verify` | `false` / `true` |
| `server.cache.tls_ca` / `tls_cert` / `tls_key` | `""` / `""` / `""` |
| `server.cache.pool_size` / `min_idle` | `10` / `2` |
| `server.cache.timeout` | `5s` |
| `server.cache.prefix` | `{project_name}:` |
| `server.cache.ttl` | `1h` |
| `server.ssl.enabled` | `false` |
| `server.ssl.cert` / `key` | `""` / `""` |
| `server.ssl.min_version` | `TLS1.2` |
| `server.ssl.letsencrypt.enabled` | `false` |
| `server.ssl.letsencrypt.email` | `""` |
| `server.ssl.letsencrypt.challenge` | `http-01` |
| `server.ssl.letsencrypt.staging` | `false` |
| `server.ssl.letsencrypt.dns_provider` | `""` |
| `server.ssl.letsencrypt.dns_credentials` | `{}` |

### Metrics, diagnostics, limits, and HTTP behavior

| Setting | Default |
|---|---|
| `server.metrics.enabled` / `root.enabled` | `true` / `true` |
| `server.metrics.auth.allow_unauthenticated` | `false` |
| `server.metrics.auth.tokens.prometheus` / `grafana` / `loki` | `""` / `""` / `""` |
| `server.metrics.include_system` / `include_runtime` | `true` / `true` |
| `server.metrics.loki.max_entries` / `max_age` | `1000` / `1h` |
| `server.metrics.duration_buckets` | `[0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]` |
| `server.metrics.size_buckets` | `[100, 1000, 10000, 100000, 1000000, 10000000]` |
| `server.debug.pprof` / `log_queries` / `log_cache` | `true` / `true` / `true` |
| `server.debug.log_bodies` | `false` |
| `server.debug.max_body_log_size` | `10KB` |
| `server.debug.block_profile_rate` / `mutex_profile_fraction` | `1` / `1` |
| `server.debug.runtime_endpoints` | `true` |
| `server.limits.max_body_size` | `10MB` |
| `server.limits.read_timeout` / `write_timeout` / `idle_timeout` | `30s` / `30s` / `120s` |
| `server.compression.enabled` / `level` | `true` / `5` |
| `server.compression.types` | `[text/html, text/css, text/javascript, application/json, application/xml]` |
| `server.trusted_proxies.additional` | `[]` |
| `server.trusted_proxies.client_ip_header` | `auto` |
| `server.cors.allowed_origins` | `[*]` |
| `server.cors.allow_credentials` / `max_age` | `true` / `86400` |
| `server.healthz.root.enabled` | `false` |
| `server.rate_limit.read.requests` / `window` | `120` / `60` |
| `server.rate_limit.write.requests` / `window` | `10` / `60` |
| `server.rate_limit.health.requests` / `window` | `120` / `60` |
| `server.rate_limit.strict.requests` / `window` | `10` / `900` |
| `server.rate_limit.auth_failures.requests` / `window` | `5` / `900` |
| `server.rate_limit.global_burst` | `240` |

The two bucket settings are lists of numeric histogram boundaries. Rate-limit
window values are seconds. Empty token settings do not create credentials.
The loader supplies the blocklist defaults shown above. The allowlist defaults
to an empty list; entries use `{cidr, description}` objects.

### Overlay networks

| Setting | Loader default |
|---|---|
| `server.tor.binary` | `""` (auto-detect `tor`) |
| `server.tor.expose` / `contact_email` | `false` / `""` |
| `server.tor.use_network` | `false` |
| `server.tor.max_circuits` | `32` |
| `server.tor.startup_timeout` / `circuit_timeout` / `bootstrap_timeout` | `60` / `60` / `180` seconds |
| `server.tor.safe_logging` | `true` |
| `server.tor.max_streams_per_circuit` / `close_circuit_on_stream_limit` | `100` / `true` |
| `server.tor.bandwidth_rate` / `bandwidth_burst` / `max_monthly_bandwidth` | `1 MB` / `2 MB` / `100 GB` |
| `server.tor.num_intro_points` / `virtual_port` | `3` / `80` |
| `server.i2p.enabled` / `expose` | `false` / `false` |
| `server.i2p.binary` | `""` (auto-detect `i2pd`) |
| `server.i2p.sam_address` | `127.0.0.1:7656` |
| `server.i2p.virtual_port` | `80` |
| `server.i2p.inbound_length` / `outbound_length` | `3` / `3` |
| `server.i2p.inbound_quantity` / `outbound_quantity` | `5` / `5` |
| `server.i2p.signature_type` | `7` |
| `server.i2p.bootstrap_timeout` | `300` seconds |

I2P's loader default is disabled. At server startup, if no explicit value is
present in configuration or `I2P_ENABLED`, the application changes that
default to enabled inside application containers. `--i2p` takes precedence
over the environment and file setting. Tor has no enable switch: the hidden
service starts when a Tor binary is found. Both managers start asynchronously;
their status appears in health responses. Exposure flags control whether the
addresses are included in autodiscovery.

### Tracking and privacy

| Setting | Default |
|---|---|
| `server.tracking.type` / `id` / `url` | `""` / `""` / `""` |
| `server.privacy.data.sold` / `stored_on_server` | `false` / `true` |
| `server.privacy.data.sharing[0].condition` / `.when` / `.data` | `analytics` / `Tracking configured (server.tracking.type set) AND user consents` / `Anonymized: page views, browser type, country` |
| `server.privacy.data.sharing[1].condition` / `.when` / `.data` | `email` / `SMTP configured for sending emails` / `Email address, message content` |
| `server.privacy.data.sharing[2].condition` / `.when` / `.data` | `user_initiated` / `User explicitly shares content (social buttons, exports)` / `Whatever user chooses to share` |
| `server.privacy.retention.period` | `""` |
| `server.privacy.retention.export_available` / `deletion_available` | `true` / `true` |
| `server.privacy.consent.show_until_acknowledged` | `true` |
| `server.privacy.consent.message` | `In accordance with the EU GDPR law this message is being displayed. We use cookies for essential site functionality and, with your consent, for preferences and analytics. Your data is stored on our servers and is never sold.` |
| `server.privacy.consent.message_if_sold` | `In accordance with the EU GDPR law this message is being displayed. We use cookies for essential site functionality and, with your consent, for preferences and analytics. Your data may be shared with or sold to third parties as described in our Privacy Policy.` |
| `server.privacy.consent.policy.text` / `url` / `version` | `Privacy Policy` / `/server/privacy` / `1` |
| `server.privacy.consent.buttons.decline` / `accept` | `Decline` / `I Agree` |
| `server.privacy.consent.position` | `bottom` |
| `server.privacy.consent.show_preferences` | `true` |
| `server.privacy.consent.preferences_text` | `Manage Preferences` |
| `server.privacy.cookies.essential.enabled` | `true` |
| `server.privacy.cookies.essential.description` | `Required for the site to function. Includes the owner cookie, security tokens (CSRF) and the consent cookie itself. The theme and language preferences are NOT essential: they are the separate consent-controlled category below. These cookies are strictly necessary and cannot be disabled.` |
| `server.privacy.cookies.preferences.enabled` | `true` |
| `server.privacy.cookies.preferences.description` | `Remember your settings such as theme (dark/light), language, and UI preferences. Without it nothing is saved and the defaults apply on each visit.` |
| `server.privacy.cookies.analytics.enabled` | `true` |
| `server.privacy.cookies.analytics.description` | `Help us understand how visitors use our site to improve the experience.` |
| `server.privacy.cookies.analytics.description_suffix_not_sold` | `Analytics data is anonymized and never sold.` |
| `server.privacy.cookies.analytics.description_suffix_sold` | `Analytics data may be shared with third parties.` |
| `server.privacy.third_party.services` | `[]` |

The default `data.sharing` records are: analytics only when tracking is set
and the visitor consents (anonymized page views, browser type, country); email
when SMTP is configured (address and message); and user-initiated sharing
(content the user chooses to share).

### Security, notifications, contact, and language

| Setting | Default |
|---|---|
| `server.security.encryption_key` / `encryption_key_version` | `""` / `0` |
| `server.security.blocklists.enabled` | `true` |
| `server.security.blocklists.auto_update` | `true` |
| `server.security.blocklists.action` | `reject` |
| `server.security.blocklists.reject_message` | `""` |
| `server.security.blocklists.log_blocked` | `true` |
| `server.security.blocklists.extra_sources` | `[]` |
| `server.security.allowlist` | `[]` (list of `{cidr, description}` entries) |
| `server.notifications.email.smtp.host` | `""` |
| `server.notifications.email.smtp.port` | `587` |
| `server.notifications.email.smtp.username` / `password` | `""` / `""` |
| `server.notifications.email.smtp.from_name` / `from_email` | `""` / `""` |
| `server.notifications.email.smtp.tls` | `auto` |
| `server.pages.contact.captcha_secret_key` | `""` |
| `server.pages.help.mcp` | Project-specific Markdown for the MCP help page; edit through `yaml.d/pages.yml` |
| `server.backup.encryption_password` | `""` (store in `secrets.yml`) |
| `server.backup.retention.max_backups` | `1` daily full backups |
| `server.backup.retention.keep_weekly` / `keep_monthly` / `keep_yearly` | `0` each |
| `server.backup.retention.max_total_size` | `10%` of the backup volume; `0` disables the cap |
| `server.contact.admin.email` | `""` |
| `server.contact.admin.webhooks.telegram` / `discord` / `slack` / `generic` | `""` each |
| `server.contact.security.email` | `security@{fqdn}` |
| `server.contact.security.webhooks.telegram` / `discord` / `slack` / `generic` | `""` each |
| `server.contact.abuse.email` | `""` |
| `server.contact.abuse.webhooks.telegram` / `discord` / `slack` / `generic` | `""` each |
| `server.contact.dpo.name` / `email` | `""` / `""` |
| `server.contact.general.email` | `""` |
| `server.contact.general.webhooks.telegram` / `discord` / `slack` / `generic` | `""` each |
| `server.i18n.default_language` | `en` |
| `server.i18n.available_languages` | `[en, es, fr, de, it, ar]` |

### Scheduled tasks

| Task setting | Default |
|---|---|
| `server.scheduler.tasks.geoip_update.enabled` / `schedule` / `retry_on_fail` / `retry_delay` | `true` / `0 3 * * 0` / `true` / `1h` |
| `server.scheduler.tasks.blocklist_update.enabled` / `schedule` / `retry_on_fail` / `retry_delay` | `true` / `0 4 * * *` / `true` / `1h` |
| `server.scheduler.tasks.cve_update.enabled` / `schedule` / `retry_on_fail` / `retry_delay` | `true` / `0 5 * * *` / `true` / `1h` |
| `server.scheduler.tasks.log_rotation.enabled` / `schedule` | `true` / `0 0 * * *` |
| `server.scheduler.tasks.token_cleanup.enabled` / `schedule` | `true` / `*/15 * * * *` |
| `server.scheduler.tasks.backup_daily.enabled` / `schedule` | `true` / `0 2 * * *` |
| `server.scheduler.tasks.backup_hourly.enabled` / `schedule` | `false` / `@hourly` |
| `server.scheduler.tasks.ssl_renewal.enabled` / `schedule` | `true` / `0 3 * * *` |
| `server.scheduler.tasks.health_check.enabled` / `schedule` | `true` / `*/5 * * * *` |
| `server.scheduler.tasks.tor_health.enabled` / `schedule` | `true` / `*/10 * * * *` |

Schedules use cron expressions.

### Maintenance defaults

| Setting | Default |
|---|---|
| `server.maintenance.self_healing.retry_interval` / `max_attempts` | `30s` / `0` (unlimited retries) |
| `server.maintenance.cleanup.disk_threshold` / `log_retention_days` | `90%` / `7` days |
| `server.maintenance.notify.on_enter` / `on_exit` | `true` / `true` |

### Branding, search metadata, browser, and data provider

| Setting | Default |
|---|---|
| `server.branding.title` | `{project_name}` |
| `server.branding.tagline` / `description` / `favicon` / `logo` | `""` each |
| `server.seo.keywords` | `[]` |
| `server.seo.author` / `og_image` / `twitter_handle` | `""` each |
| `server.seo.title_separator` | ` - ` |
| `server.seo.meta_tags` | `{}` |
| `server.seo.canonical` | `""` |
| `server.seo.noindex` | `false` |
| `server.seo.opengraph.enabled` / `type` / `site_name` | `true` / `website` / `""` |
| `server.seo.twitter.enabled` / `card` / `creator` | `true` / `summary_large_image` / `""` |
| `web.ui.theme` / `accent` | `dark` / `""` |
| `web.robots.allow` / `deny` | `[/, /api]` / `[/debug]` |
| `data.cve.api_key` | `""` |

## CLI profile defaults

`timezones-cli` creates `cli.yml` under the user configuration directory.
An explicitly selected `--config` path is used instead. Empty server fields
such as `server.primary` remain empty unless provided by a profile, environment,
flag, or release build.

| Setting | Default |
|---|---|
| `server.api_version` | `v1` |
| `server.timeout` | `30s` |
| `server.retry` / `server.retry_delay` | `3` / `1s` |
| `auth.token` / `auth.token_file` | `""` / `""` |
| `output.format` / `output.color` / `output.pager` | `table` / `auto` / `auto` |
| `output.quiet` / `output.verbose` | `false` / `false` |
| `tui.enabled` / `tui.theme` / `tui.mouse` / `tui.unicode` | `true` / `dark` / `true` / `true` |
| `display.mode` | `auto` |
| `logging.level` | `warn` |
| `logging.file` | `""` (uses `{user_log_dir}/cli.log`) |
| `logging.max_size` / `logging.max_files` | `10MB` / `5` |
| `cache.enabled` / `cache.ttl` / `cache.max_size` | `true` / `5m` / `100MB` |
| `debug` | `false` |
| `update.auto` / `update.check_interval` / `update.channel` | `false` / `per_invocation` / `stable` |
| `defaults.lang` / `defaults.output` / `defaults.limit` | `auto` / `table` / `20` |

The CLI creates config, data, cache, and log directories with mode `0700` and
the profile and log file with mode `0600` on Unix-like systems.
The CLI checks the selected server's published release on each invocation.
`update.auto: true` installs an available update without prompting only for
noninteractive runs; interactive runs prompt, while `--update yes` installs
without prompting and `--update check` only reports availability. The server
currently publishes one CLI release channel, so `update.channel` must remain
`stable` and `update.check_interval` must remain `per_invocation`.

## Environment variables

### Server

| Variable | Effect and scope |
|---|---|
| `MODE` | Overrides `server.mode` on configuration load. |
| `SMTP_HOST` / `SMTP_USERNAME` / `SMTP_PASSWORD` | Overrides corresponding SMTP fields on configuration load. |
| `SMTP_FROM_NAME` / `SMTP_FROM_EMAIL` | Overrides `server.notifications.email.from.name` / `server.notifications.email.from.email` on configuration load. |
| `SMTP_TLS` | Overrides `server.notifications.email.smtp.tls` on configuration load. |
| `SMTP_PORT` | Overrides SMTP port only when it parses as an integer. |
| `PORT` / `LISTEN` | Seeds `server.port` / `server.address` during first-run setup; also supplies the listener when CLI flags have not supplied values. |
| `I2P_ENABLED` | Explicit I2P enable/disable setting; a valid `--i2p` flag overrides it. |
| `APPLICATION_NAME` / `APPLICATION_TAGLINE` | On first run seed `server.branding.title` / `server.branding.description`. |
| `DATABASE_URL` / `CACHE_URL` | On first run seed connection URLs if not already configured. Existing saved values win. |
| `DATABASE_DRIVER` | First-run driver selection when applicable; URL scheme takes precedence when it identifies the driver. |
| `DOMAIN` | Public FQDN and allowed-domain input used by URL and contact handling. |
| `CONFIG_DIR` / `DATA_DIR` / `CACHE_DIR` / `LOG_DIR` / `BACKUP_DIR` | Override corresponding application directories. |
| `DATABASE_DIR` | Override the database directory when resolving the database location. |
| `XDG_CONFIG_HOME` / `XDG_DATA_HOME` / `XDG_CACHE_HOME` | Unix user-directory bases when the corresponding explicit directory is not set. |
| `APPDATA` / `LOCALAPPDATA` / `PROGRAMDATA` / `USERPROFILE` / `SystemDrive` | Windows application directory bases. |
| `HOSTNAME` | Candidate public host during URL resolution. |
| `INVOCATION_ID` / `RC_SVCNAME` / `KUBERNETES_SERVICE_HOST` / `container` | Runtime/service/container detection; not application settings. |

The container entrypoint defaults these directory variables to
`CONFIG_DIR=/config/timezones`, `DATA_DIR=/data/timezones`,
`CACHE_DIR=/data/timezones/cache`, `LOG_DIR=/data/log/timezones`,
`DATABASE_DIR=/data/db/sqlite`, and `BACKUP_DIR=/data/backups/timezones`.
Explicit environment values override these defaults.

`PORT`, `LISTEN`, and server options also have command-line equivalents. The
server does not support a generic environment variable for every YAML leaf.

### CLI

The CLI uses Viper automatic environment mapping with prefix `TIMEZONES_` and
`.` to `_` key conversion, so profile keys can be overridden using their
corresponding upper-case names (for example, `TIMEZONES_SERVER_TIMEOUT`).
Explicitly read variables are:

| Variable | Effect |
|---|---|
| `TIMEZONES_SERVER_PRIMARY` | Primary server URL |
| `TIMEZONES_TOKEN` / `TIMEZONES_AUTH_TOKEN` | API token, checked in that order |
| `TIMEZONES_LANG` | Output language |
| `TIMEZONES_FORMAT` | Output format |
| `LANG` | Locale fallback when CLI language is `auto` |
| `NO_COLOR` / `TERM` | Disable automatic terminal color (`NO_COLOR` set or `TERM=dumb`) |
| `COLORFGBG` | Terminal background-color detection for automatic theme selection |
| `SHELL` | Shell selection for shell integration output |
| `SSH_CLIENT` / `SSH_TTY` / `SSH_CONNECTION` / `MOSH` | Detect remote sessions for GUI availability |
| `DISPLAY` / `WAYLAND_DISPLAY` / `OS` / `TERM_PROGRAM` / `__CFBundleIdentifier` / `XPC_SERVICE_NAME` | Detect a local desktop for GUI availability |
| `XDG_CONFIG_HOME` | Unix user configuration directory base |
| `APPDATA` / `LOCALAPPDATA` | Windows configuration and data/cache/log directory bases |

## Command-line flags

### Server binary

| Flag | Effect |
|---|---|
| `--port` | Listener port |
| `--address` | Listener address |
| `--config` | Configuration directory |
| `--data` | Data directory |
| `--cache` | Cache directory |
| `--log` | Log directory |
| `--backup` | Backup directory |
| `--mode` | Runtime mode |
| `--i2p[=BOOL]` | Override whether to start I2P; default is on in an application container when not explicitly configured |
| `--version` | Print build version information |
| `--status` | Check the configured health endpoint |
| `--help` | Print usage |
| `--service {--install,--uninstall,--disable,start,stop,restart,reload,--help}` | Service manager operation |
| `--maintenance on` / `--maintenance off` | Enable or disable runtime maintenance mode |
| `--maintenance token list` / `revoke PREFIX [--yes]` | List token metadata or revoke a token |
| `--maintenance security list [--status STATUS]` / `show ID` / `export ID [--output FILE]` / `update ID [options]` | Manage security reports; with working SMTP, a changed status sends the researcher a notice without report content |
| `--maintenance pgp generate` / `rotate` / `publish` / `export` / `import` / `delete` | Manage the PGP key |
| `--maintenance scheduler list` / `show` / `run` / `enable` / `disable` / `history` | Manage scheduler tasks |
| `--update [check\|yes\|branch stable\|beta\|devel]` | Check or apply an update |

### `timezones-cli`

| Flag | Effect |
|---|---|
| `--config` | CLI profile path |
| `--help` | Show CLI help (Cobra built-in) |
| `--version` | Show CLI version (Cobra built-in) |
| `--server` | Server URL |
| `--token` / `--token-file` | Token value or token file |
| `--lang` | Output language |
| `--output {table,json,plain}` | Output encoding |
| `--color {auto,yes,no}` | Color selection |
| `--limit` | Maximum result count |
| `--debug` | Enable client debug output |
| `--shell` | Print shell integration or completion output |
| `--update {check,yes}` | Check for or install a CLI update |
| `list --page` | List page number; default `1` |
| `list --utc-offset` | Filter by UTC offset |
| `list --isdst {true,false}` | Filter daylight-saving records |
| `token create --name` / `--expires` | Linked-token name and optional expiry (`7d`, `30d`, `90d`, `1y`, `never`) |
| `token revoke --confirm` | Confirmation text for revoking the primary token |

## AI.md configuration gaps

The server currently supports the Tor and I2P configuration tables above.
The following AI.md-required settings or behaviors are not present in the
loader default tree or were not found in the current implementation:

- AI.md PART 11's `server.security.blocked_ips` and
  `server.security.abuse_detection.*` settings are absent. The current code
  supports the separate `server.security.allowlist` and source-based
  `server.security.blocklists` settings documented above.
- AI.md PART 12's general rate-limit class and counter behavior is not
  implemented as specified. The loader defines `server.rate_limit.*` defaults,
  but the general request limiter remains a fixed 240 requests per minute per
  client address; it does not select limits from those settings or key the
  general quota by Bearer token.
- AI.md configuration areas such as `server.csrf`, `server.url_detection`,
  `server.compliance`, `server.tokens`, `server.restore`, and top-level `cli`
  do not have built-in leaf defaults in the current loader. Although the
  loader preserves some of these sections in the merged configuration, source
  usage was not confirmed for every AI.md-required setting. Do not assume an
  undocumented leaf is consumed just because a YAML file accepts it.
