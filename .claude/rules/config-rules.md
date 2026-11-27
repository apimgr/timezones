# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never generate JSON or `.yaml` config files; server config is runtime `server.yml`, with area drop-ins under `yaml.d/` and secrets in `secrets.yml`.
- Never silently ignore unknown configuration keys or accept invalid values without the specified warning/default behavior.
- Never store secrets in normal configuration, logs, docs, or generated source.
- Never let development mode disable sanitization or implicitly enable debug endpoints.

## CRITICAL - ALWAYS DO

- Apply precedence in order: CLI flags, environment, config files, defaults.
- Validate configuration at load, normalize supported booleans, preserve source/precedence metadata, and hot-reload only as specified.
- Keep branding, SEO, theme, and accent settings in `yaml.d/branding.yml`; never add them to generated `server.yml`.
- Default to production; keep debug explicitly gated by `--debug`, `DEBUG`, or the specified debug mode behavior.

## Key Rules

- Use exact `server.*`, `web.*`, `data.*`, and `cli.*` key paths from AI.md.
- Invalid values use their documented fallback and one warning names the affected key.
- Keep privacy, tracking, cookie, branding, and page-content behavior synchronized with rendered pages and docs.

For complete details, see AI.md PART 5, 6, 12.
