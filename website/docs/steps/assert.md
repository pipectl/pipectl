# assert

Checks record-count, field-existence, and field-value conditions. The pipeline fails if any assertion is not met.

**Supported formats:** `json` `jsonl` `csv`

## Options

At least one option is required.

| Option | Type | Description |
|--------|------|-------------|
| `min-records` | integer | Payload must have at least this many records. Must be ≥ 0. |
| `max-records` | integer | Payload must have at most this many records. Must be ≥ 0. |
| `records-equal` | integer | Payload must have exactly this many records. Must be ≥ 0. |
| `field-exists` | string | The named field must exist in the payload. |
| `field-equals` | object (`field`, `value`) | Every record's `field` must equal `value`. |
| `field-contains` | object (`field`, `value`) | Every record's `field` must contain `value` as a substring. |
| `field-matches` | object (`field`, `value`) | Every record's `field` must match the regular expression in `value`. |
| `case-sensitive` | boolean | Whether `field-equals`, `field-contains`, and `field-matches` compare case-sensitively. Defaults to `true`. |

## Examples

Assert a minimum and maximum:

```yaml
- assert:
    min-records: 1
    max-records: 10000
```

Assert an exact count:

```yaml
- assert:
    records-equal: 5
```

Assert a field exists:

```yaml
- assert:
    field-exists: email
```

Assert every record's field equals a value:

```yaml
- assert:
    field-equals:
      field: status
      value: active
```

Assert every record's field contains a substring:

```yaml
- assert:
    field-contains:
      field: email
      value: "@"
```

Assert every record's field matches a regular expression:

```yaml
- assert:
    field-matches:
      field: email
      value: '^[^@]+@[^@]+\.[^@]+$'
```

Case-insensitive value checks:

```yaml
- assert:
    field-equals:
      field: status
      value: ACTIVE
    case-sensitive: false
```

Combined:

```yaml
- assert:
    min-records: 10
    max-records: 1000
    field-exists: email
    field-equals:
      field: status
      value: active
```

## Notes

- `min-records` must be ≤ `max-records` when both are set.
- `records-equal` must fit within the `min-records` / `max-records` bounds when they are also set.
- For CSV, `field-exists` checks the header row.
- For JSON and JSONL, `field-exists` passes if any record contains the field.
- `field-equals`, `field-contains`, and `field-matches` check every record, not just any record — this differs from `field-exists`, which only checks the schema/header. The assertion fails on the first record where the field is missing, or where the value fails the check.
- `field-matches` uses Go's RE2 regular expression engine — see the [filter step's regex syntax notes](./filter#regex-syntax) for supported syntax and caveats.
- `case-sensitive` defaults to `true` and only affects `field-equals`, `field-contains`, and `field-matches` (it has no effect on `field-exists`, which compares field names, not values).
