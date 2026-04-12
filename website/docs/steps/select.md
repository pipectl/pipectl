# select

Keeps only the specified fields, dropping everything else.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | list | Yes | Field names to keep |

## Example

```yaml
- select:
    fields: [first_name, last_name, email, plan]
```

## Notes

- Fields are kept in the order they appear in the `fields` list.
- For JSON and JSONL, fields missing from a record are silently omitted — no error is returned.
- For CSV, only the specified column names are kept. The column order in the output matches the `fields` list.
