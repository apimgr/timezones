# Frontend Rules (PART 16)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use client-side rendering as the primary page implementation or require JavaScript for core navigation, forms, or content.
- Never render unsanitized operator HTML/Markdown or link operator content to `/server/contact`.
- Never omit the page shell, one main landmark, one page `h1`, skip links, secure POST forms, or the required footer.
- Never place Help, About, Privacy, Terms, Preferences, or API docs in the application-section navigation.

## CRITICAL - ALWAYS DO

- Use server-rendered Go templates, mobile-first accessible HTML, and progressive enhancement.
- Keep skip links first-focusable, the header in its required three zones, current-nav state paired with `aria-current`, and breadcrumbs where required.
- Implement all required public pages, registries, error/health pages, contact/security pages, preferences, and data views.
- Render operator branding through `operatorHTML`; use Goldmark plus Bluemonday, language fallback, placeholder substitution, and the 256 KiB limit.
- Keep the shared header/footer, responsive CSS, no-JS fallbacks, and page-conformance tests synchronized.

## Key Rules

- About and Privacy extend built-in content; Terms overrides it; custom Footer branding sits above the always-shown application footer.
- Header/footer allow strict sanitized HTML only; document content preserves the documented safe table/code elements.
- Follow the `server.branding.content` behavior and heading-level shift in AI.md exactly.

For complete details, see AI.md PART 16.
