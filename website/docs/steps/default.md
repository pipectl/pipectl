# default

Fills missing or empty fields with a default value.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | map | Yes | Map of field name to default value |

## Example

```yaml
- default:
    fields:
      country: AU
      currency: AUD
      source: import
```

## Notes

- **JSON / JSONL:** defaults are applied only when the field does not already exist in a record. Existing fields — including those with `null` values — are not overwritten.
- **CSV:** a missing column is added to the header and populated for all rows. An existing column is only filled where the cell value is empty.
