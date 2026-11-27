# Integrations and discovery

The application has no third-party identity-provider or federation dependency. Its MCP endpoint uses the server's owner and linked Bearer tokens. When a Tor binary is available, the server starts its own managed Tor hidden service. I2P uses a managed `i2pd` process when its binary is available and otherwise connects to the configured SAM bridge. Both use separate loopback backends, are stopped with the server, and report lifecycle state through health. `/api/autodiscover` includes an overlay address only when its corresponding `expose` setting is enabled.

Tor starts automatically when its binary is found; there is no enable switch.
Its binary path and operational settings are under `server.tor.*`. I2P is
disabled by default outside application containers and enabled by default
inside them. Set `server.i2p.enabled`, `I2P_ENABLED`, or `--i2p[=true|false]`
to control it. Its provider settings are under `server.i2p.*`. The Tor `.onion`
and I2P `.b32.i2p` services use HTTP as required by their overlay protocols;
they do not use TLS or HTTPS redirects.

Tor maps the hidden service to a dedicated local backend and uses PROXY
protocol metadata for circuit identity. I2P uses a separate plain loopback
backend and stores its destination key under the data directory. The server
does not expose the Tor/I2P service status or addresses through discovery
unless the relevant exposure setting is enabled; health reports lifecycle
status. `server.tor.use_network` enables the managed outbound SOCKS client, but
the current application does not route its own outbound requests through that
client.

## Public endpoints

The browser interface serves `/`, `/server/about`, `/server/help`,
`/server/help/mcp`, `/server/privacy`, `/server/terms`, `/server/dpo`,
`/server/security`, `/server/security/policy`, `/server/security/thanks`,
`/server/security/report` and `/server/security/report/{tracking_id}`,
`/server/contact`, `/server/preferences`, `/server/preferences/sync`,
`/server/preferences/sync/import`, `/server/preferences/tokens`,
`/server/docs/swagger`, and `/server/docs/graphql`. Timezone pages are served
at `/timezones` and `/timezones/search`, `/timezones/random`,
`/timezones/offset/{offset}`, `/timezones/abbr/{abbr}`, `/timezones/utc/{id}`,
`/timezones/value/{name}`, and `/timezones/current/{id}`, with `/stats` and
`/count` for dataset summaries. The contact page accepts security reports as
well as general contact submissions.

- REST API: `/api/v1/`
- OpenAPI JSON: `/api/v1/server/swagger`
- OpenAPI alias: `/api/swagger`
- GraphQL API (GET and POST) and browser explorer: `/api/v1/server/graphql` and `/server/docs/graphql`
- GraphQL GET and POST alias: `/api/graphql`
- Health: `/api/v1/server/healthz` and the server-rendered `/server/healthz`
- Health alias: `/api/healthz`
- Metrics aliases: `/server/metrics[/{service}]`, `/metrics[/{service}]`, and `/api/metrics[/{service}]`
- Health plain-text aliases: `/api/healthz.txt` and `/api/v1/server/healthz.txt`
- Robots rules: `/robots.txt`
- Sitemap: `/sitemap.xml`
- Atom and RSS feeds: `/feed.xml` and `/rss.xml` (valid feeds with no entries because the embedded dataset has no update history)
- Agent-readable endpoint summary: `/.well-known/llms.txt` and `/llms.txt`
- Machine-only Bearer-token validation for CLI setup: `/api/v1/server/auth/check`
- Security overview and reporting: `/server/security`, `/server/contact`, and `/server/security/report/{tracking_id}`
- Security disclosure policy and researcher acknowledgments: `/server/security/policy` and `/server/security/thanks`
- Machine-readable security contact: `/.well-known/security.txt`
- Project OpenPGP public key when configured: `/.well-known/pgp-key.asc`
- MCP Streamable HTTP aliases: `/mcp`, `/api/mcp`, and `/api/v1/server/mcp` (Bearer token required)
- CLI release discovery and verified downloads: `/api/autodiscover` and `/cli/binaries/{binary}`
- PWA manifest and service worker: `/manifest.json` and `/sw.js`
- Locale catalogs: `/locales/{language}.json`

The dataset and static web assets are embedded in the server binary. CORS
origins default to `*`; see [configuration](configuration.md) for the CORS
settings.

The browser interface is installable as a PWA through its web manifest. Its service worker caches the application shell and falls back to the cached homepage for offline navigation; API data still requires a network connection.
