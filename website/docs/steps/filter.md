# filter

Keeps only records that match a condition. Records that do not match are dropped.

**Supported formats:** `json` `jsonl` `csv`

## Single condition

Specify `field` and exactly one operator:

| Operator | Description |
|----------|-------------|
| `equals` | Field value equals the given string (numeric comparison when both values are numbers) |
| `not-equals` | Field value does not equal the given string |
| `contains` | Field value contains the given substring |
| `starts-with` | Field value starts with the given string |
| `ends-with` | Field value ends with the given string |
| `matches` | Field value matches the given regular expression (RE2 syntax — see [Regex syntax](#regex-syntax)) |
| `greater-than` | Field value is numerically greater than the given number |
| `less-than` | Field value is numerically less than the given number |

```yaml
- filter:
    field: country
    equals: AU
```

```yaml
- filter:
    field: status
    not-equals: inactive
```

```yaml
- filter:
    field: age
    greater-than: "18"
```

```yaml
- filter:
    field: email
    contains: "@example.com"
```

```yaml
- filter:
    field: email
    matches: '^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$'
```

## Multi-condition: `all` (AND)

Keep records that match **every** listed condition:

```yaml
- filter:
    all:
      - field: status
        equals: active
      - field: age
        greater-than: "18"
```

## Multi-condition: `any` (OR)

Keep records that match **at least one** listed condition:

```yaml
- filter:
    any:
      - field: country
        equals: AU
      - field: country
        equals: NZ
```

## Nested groups

`all` and `any` can be nested arbitrarily deep:

```yaml
- filter:
    all:
      - field: age
        greater-than: "18"
      - any:
          - field: country
            equals: AU
          - field: country
            equals: NZ
```

## Missing fields

By default, records missing the field a rule tests are treated as non-matching and silently excluded — no error or warning is raised at the default log level. Run with `--verbose` (`-v`) to see the count of excluded records. Set `on-missing` on the step to change this:

| Value | Behavior |
|-------|----------|
| `exclude` (default) | Records missing the field are excluded — the original behavior. |
| `include` | Records missing the field are treated as matching. |
| `error` | The pipeline fails with an error naming the missing field. |

```yaml
- filter:
    field: status
    equals: active
    on-missing: include # keep records that don't have `status` at all
```

`on-missing` is a single step-level setting — it's not configurable per rule, and it applies uniformly to every leaf rule evaluated by the step, including ones nested inside `all`/`any` groups.

## Case sensitivity

By default, `equals`, `not-equals`, `contains`, `starts-with`, `ends-with`, and `matches` compare values case-sensitively. Set `case-sensitive: false` on the step to fold case for these operators:

```yaml
- filter:
    field: email
    equals: "alice@example.com"
    case-sensitive: false # matches ALICE@EXAMPLE.COM too
```

For `matches`, case folding is applied by prefixing the pattern with the inline flag `(?i)` before compiling — equivalent to writing `(?i)` at the start of the pattern yourself.

Like `on-missing`, `case-sensitive` is a single step-level setting — it's not configurable per rule, and it applies uniformly to every leaf rule, including ones nested inside `all`/`any` groups.

## Regex syntax

`matches` compiles the pattern using Go's standard library [`regexp`](https://pkg.go.dev/regexp) package, which implements **RE2** syntax — described in full at [pkg.go.dev/regexp/syntax](https://pkg.go.dev/regexp/syntax). This is **not** the same flavor as PCRE (Perl/JavaScript/Python-style regex): RE2 guarantees linear-time matching (a user-supplied pattern can't cause a catastrophic-backtracking hang), but it does not support **backreferences** (e.g. `\1`) or **lookahead/lookbehind** (e.g. `(?=...)`, `(?<=...)`). Users coming from PCRE-flavored tools may need to rework patterns that rely on those features.

Common syntax that *is* supported:

| Pattern | What it does |
|---------|--------------|
| `^\d{3}-\d{4}$` | Anchors + digit character class + quantifier (e.g. a phone extension field) |
| `(?i)^error` | Inline case-insensitive flag (equivalent to `case-sensitive: false`) |
| `^(foo\|bar\|baz)$` | Alternation |
| `^[A-Z][a-z]+ [A-Z][a-z]+$` | Character classes + quantifiers (e.g. a "First Last" name shape) |

An invalid pattern is rejected when the pipeline is validated (at load time), not at runtime — the same as an invalid `greater-than`/`less-than` value.

## Notes

- For JSON and JSONL, non-string field values are coerced to strings before comparison.
- `greater-than` and `less-than` require the field value to be parseable as a number. Records with non-numeric values will cause the step to fail.
- `all` and `any` cannot be combined at the same nesting level.
- Group conditions (`all`, `any`) and flat rule fields (`field`, `equals`, etc.) cannot be mixed on the same step.
- `case-sensitive: false` only affects `equals`, `not-equals`, `contains`, `starts-with`, `ends-with`, and `matches` — it has no effect on `greater-than`/`less-than`. Defaults to `true`.
