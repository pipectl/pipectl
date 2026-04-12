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

## Notes

- Records missing the specified field are always excluded.
- For JSON and JSONL, non-string field values are coerced to strings before comparison.
- `greater-than` and `less-than` require the field value to be parseable as a number. Records with non-numeric values will cause the step to fail.
- `all` and `any` cannot be combined at the same nesting level.
- Group conditions (`all`, `any`) and flat rule fields (`field`, `equals`, etc.) cannot be mixed on the same step.
