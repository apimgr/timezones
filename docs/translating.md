# Translating

The shared, embedded catalogs live in `src/common/i18n/locales/`. English,
Spanish, French, Arabic, German, and Italian are registered in
`src/common/i18n/i18n.go`; Arabic is right-to-left. The Go server and CLI use
these same files, and the server exposes each registered catalog at
`/locales/{language}.json`.

For a browser request, `?lang={code}` selects a supported language for that
response. It stores a one-year `lang` cookie only after the visitor accepts the
Preferences consent category; without that consent, the selection applies only
to the current response. On later requests the server uses the language cookie
only while Preferences consent remains accepted, then negotiates
`Accept-Language`, and finally falls back to English. An unsupported query
value selects English for that response and does not write a cookie.

Visitors can also choose a language from the Language section of
`/server/preferences`; the first choice is Automatic. The footer's active
language link opens that section. The six locale catalogs share the same keys
in server-rendered pages, API error/status messages, OpenAPI descriptions,
GraphQL type/field descriptions, and client interfaces. Arabic uses RTL layout.
JavaScript and server-rendered pages must keep the same language and direction.

To add a language, add a complete JSON catalog with the same message keys and
placeholder names as `en.json`, register its native name, text direction, and
CLDR plural categories, then update the language list in the README and this
guide. Keep protocol names, command examples, API paths, configuration keys,
and literal values unchanged. Review accessibility labels and RTL layout in a
browser before describing the translation as complete.
