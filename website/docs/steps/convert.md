# convert

Converts the payload to a different format.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `format` | string | Yes | Target format: `json`, `jsonl`, or `csv` |

## Example

```yaml
- convert:
    format: csv
```

## Conversion behaviour

| From | To | Notes |
|------|----|-------|
| CSV | JSON / JSONL | First row becomes field names. Dot-separated headers (e.g. `user.name`) become nested JSON objects. |
| JSON / JSONL | CSV | Records are flattened. Nested objects become dot-separated headers. Arrays and objects within field values are JSON-encoded as strings. |
| JSON | JSONL | Array elements become individual lines. |
| JSONL | JSON | Lines are collected into a JSON array. |

## Notes

- Converting to the same format the payload is already in is a no-op.
- `convert` only changes the in-memory format. It does not affect the `output.format` setting, which controls how the final payload is written.
