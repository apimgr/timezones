# Translating

The shared, embedded catalogs live in `src/common/i18n/locales/`. English,
Spanish, French, Arabic, German, and Italian are registered in
`src/common/i18n/i18n.go`; Arabic is right-to-left. The Go server and CLI use
these same files, and the server exposes each registered catalog at
`/locales/{language}.json`.

For a browser request, `?lang={code}` selects a supported language and stores a
one-year `lang` cookie. Without that query parameter, the server checks the
cookie, then the `Accept-Language` header, and finally uses English. Unsupported
codes fall back to English. JavaScript and server-rendered pages must keep the
same language and direction.

To add a language, add a complete JSON catalog with the same message keys and
placeholder names as `en.json`, register its native name, text direction, and
CLDR plural categories, then update the language list in the README and this
guide. Keep protocol names, command examples, API paths, configuration keys,
and literal values unchanged. Review accessibility labels and RTL layout in a
browser before describing the translation as complete.
