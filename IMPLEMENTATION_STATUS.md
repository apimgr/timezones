# Implementation status

The active product specification is [IDEA.md](IDEA.md), and the implementation
rules are in [AI.md](AI.md). This file replaces an older status report that
referred to a removed SPEC.md and described features the current IDEA.md does
not require.

The current service embeds timezone JSON and web assets, provides REST lookup
routes, a server-rendered search page, static PWA assets, a health endpoint,
client-side rate protection, and the `timezones-cli` lookup binary. GitHub
workflows build the server and CLI releases and the multi-architecture Docker
image.

Known unfinished requirements and verification limits are tracked in
[TODO.AI.md](TODO.AI.md). This status is descriptive only; no build, test, lint,
format, or workflow validation was performed for the current implementation
pass.
