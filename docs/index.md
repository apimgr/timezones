# Timezones API documentation

Timezones API serves a bundled, public timezone reference dataset through versioned REST and read-only GraphQL APIs, a server-rendered web interface, owner and linked bearer tokens, authenticated MCP, and the `timezones-cli` companion client.

The records include display names, abbreviations, UTC offsets, daylight-saving flags, labels, and IANA identifiers. Record offsets and flags are static data. The API separately calculates current local time from an IANA identifier, but does not provide future daylight-saving transition schedules.

## Start here

- [Install the server](installation.md)
- [Configure the server](configuration.md)
- [Use the REST API](api.md)
- [Use the CLI](cli.md)
- [Security and privacy](security.md)
- [Integrations and discovery](integrations.md)
- [Development](development.md)
- [Translation guide](translating.md)

The application has no user accounts, password login, or third-party identity provider. Tor hidden-service and optional I2P eepsite behavior is described in [Integrations](integrations.md).
