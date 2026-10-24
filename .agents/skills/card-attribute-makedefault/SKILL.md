---
name: card-attribute-makedefault
description: Register card attribute values (rarity, variety, category, slot...) read from loaded card files with MakeDefault, or add MakeDefault to a new attribute package.
---

# Card attribute MakeDefault

Attribute packages live in `app/cmdrun/domain/card/attributes/<attr>/caching.go`.
`MakeDefault(value)` adds `value` to the attribute's runtime list of valid
values, then returns it. After that, `Validate()` accepts it.

## Rule

A value read from a loaded card file (YAML template) is trusted: register it.
A value typed by the user is not: use `AssertWithDefault`, never `MakeDefault`.

## Where to call it

In the card package's `loadFromFile`, right after `LoadFromFile` succeeds. Call
it only on non-empty values, because `MakeDefault("")` would make the empty
string valid for the whole game:

```go
// Values from a loaded card file are trusted: register them as valid.
if params.Variety != "" {
	variety.MakeDefault(string(params.Variety))
}
```

Call it once, at load time. In `newFromYaml`, convert the value with the type
(`variety.Variety(config.Variety)`), not with `MakeDefault`.

Existing examples: `card/characters/caching.go`, `card/locations/caching.go`,
`card/objects/equipment/caching.go`.

## Adding MakeDefault to an attribute that lacks it

Copy the `variety` version and use the package's own types:

```go
func MakeDefault(defaultValue string) Variety {
	_, _ = ListWithError() // load the file first so Update does not get overwritten
	cachedConfig.Update(func(config *yamlConfig) {
		config.addValue(Variety(defaultValue))
	})
	return Variety(defaultValue)
}
```

`addValue` must already skip duplicates. Check this in the package before you
reuse it. Add a test in the package's `caching_test.go` (see
`variety/caching_test.go`) that checks `Validate()` returns true after
`MakeDefault`.
