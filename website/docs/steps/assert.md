# assert

Checks record-count and field-existence conditions. The pipeline fails if any assertion is not met.

**Supported formats:** `json` `jsonl` `csv`

## Options

At least one option is required.

| Option | Type | Description |
|--------|------|-------------|
| `min-records` | integer | Payload must have at least this many records. Must be ≥ 0. |
| `max-records` | integer | Payload must have at most this many records. Must be ≥ 0. |
| `records-equal` | integer | Payload must have exactly this many records. Must be ≥ 0. |
| `field-exists` | string | The named field must exist in the payload. |

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

Combined:

```yaml
- assert:
    min-records: 10
    max-records: 1000
    field-exists: email
```

## Notes

- `min-records` must be ≤ `max-records` when both are set.
- `records-equal` must fit within the `min-records` / `max-records` bounds when they are also set.
- For CSV, `field-exists` checks the header row.
- For JSON and JSONL, `field-exists` passes if any record contains the field.
