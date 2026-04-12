# sort

Sorts records by a single field.

**Supported formats:** `json` (arrays only) `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `field` | string | Yes | Field name to sort by |
| `direction` | string | No | `asc` or `desc`. Defaults to `asc`. |

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

## Notes

- Records with a missing or `null` field value are always sorted last, regardless of direction.
- For CSV, empty string values are treated as missing and sorted last.
- When field values parse as numbers, numeric ordering is used. Otherwise, string ordering is used.
- JSON object payloads (single record) are not supported — convert to a JSON array or JSONL first.
