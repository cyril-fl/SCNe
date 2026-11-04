---
name: i19n-audit
description: Find and fix hard-coded player-facing strings and unresolved translation keys, so all text goes through the i19n module. Use when adding or reviewing UI text, messages, or locale files, or when asked to check i18n/i19n coverage.
---

# i19n audit

All text the player reads goes through `app/shared/services/i19n`.
`i19n.Localize(id)` returns the id itself when the message is missing, so a
missing key shows up as raw text on screen rather than as an error.

## Run the audit

From the project root:

```sh
go run .agents/skills/i19n-audit/scripts/audit.go
```

It reports four lists. The first is reliable; the others need review.

| Section | Meaning | Action |
|---|---|---|
| Keys missing from a locale | Key not in every locale, often a typo in one | Fix or add the key |
| Keys used in code but absent | Rendered raw on screen | Add the key to all locales, unless it is a technical subject (see below) |
| Hard-coded text candidates | Heuristic, includes technical strings | Read each one in context |
| Raw errors | `logs.RawError`, `fmt.Errorf`, `errors.New` | Check whether the player sees it |

The audit cannot see keys built at runtime (`"stat."+name`) or names coming from
`data/`. For those, compare `data/` values with the locale files by hand, such as
oracle names against `oracles.<name>`.

## Fixing

- Add each key to every locale, in the themed folder that fits:
  `locales/<theme>/<locale>.yaml`, or `locales/<locale>.yaml` for global text.
- Keys are flat: `field.none:` then `other: "..."`. Reuse an existing key first.
- Use `i19n.Localize` for the player and `i19n.Unlocalize` for English logs.
  `logs.Error(id, data)` already handles both.
- String values in `data` are translated too, so pass keys there
  (`"Property": "property.name"`).
- Go always prints month names in English. Localize the layout
  (`journal.date:format`) and use a numeric one where English months would show.

## Leave these as they are

- **Before i19n is initialized:** `init()` in `config/`, and the errors in
  `i19n/config.go` and `initI19n/run.go`. `Localize` would return `{{id}}`.
- **Commands annotated `annotations.NODB`** (`version`, `config`, `try`): they skip
  `setup`, so i19n is not initialized unless the command does it itself.
- **Technical strings:** field names, collections, key bindings, colors, and
  formats with no words, such as `"%s | %s"`.
- **Subjects that only reach the logs** (`Silent*` calls, such as `"savestate"`).
- **The CLI help** in `config/.config.yaml`, which stays in English.

## Ask before changing

- **Stored values:** a translated string saved to the database freezes its language.
  Store a stable value and translate on display, as `campaign.DisplayAuthor` does.
- **`logs.RawError` to `logs.Error`:** this drops the `%w` chain, which breaks
  `errors.Is`.
- **Renaming a key used as a data name:** saved games may reference it.

## Verify

1. Run `go build ./...`.
2. Rerun the audit and check that the lists shrank as expected.
3. Run the tests with `$go-test-run`. A test calling a function that now uses
   `Localize` may need a translator; invoke `$go-test-create` before updating it.
