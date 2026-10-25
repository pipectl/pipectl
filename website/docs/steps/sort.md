# sort

Sorts records by a single field.

**Supported formats:** `json` (arrays only) `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `field` | string | Yes | Field name to sort by |
| `direction` | string | No | `asc` or `desc`. Defaults to `asc`. |
| `nulls` | string | No | `first` or `last`. Where records with a missing, `null`, or empty (CSV) value are placed. Defaults to `last`. |

## Examples

Sort by last name ascending (default):

```yaml
- sort:
    field: last_name
```

Sort by date descending:

```yaml
- sort:
    field: created_at
    direction: desc
```

Put records with no `shipped_at` value first:

```yaml
- sort:
    field: shipped_at
    nulls: first
```

## Notes

- Records with a missing or `null` field value are sorted last by default. Set `nulls: first` to place them first instead. Null placement is independent of `direction`.
- For CSV, empty string values are treated as missing.
- Records with missing values keep their original relative order.
- When field values parse as numbers, numeric ordering is used. Otherwise, string ordering is used.
- JSON object payloads (single record) are not supported — convert to a JSON array or JSONL first.
