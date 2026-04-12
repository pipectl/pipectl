# rename

Renames fields in the payload.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | map | Yes | Map of current field name → new field name |

## Example

```yaml
- rename:
    fields:
      first_name: firstName
      last_name: lastName
      credit_card: creditCard
```

## Notes

- For CSV, the rename applies to column headers.
- Fields that are listed but do not exist in the payload are silently ignored.
