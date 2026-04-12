# limit

Truncates the payload to at most N records.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `count` | integer | Yes | Maximum number of records to keep. Must be ≥ 1. |

## Example

```yaml
- limit:
    count: 100
```

## Notes

- If the payload already has fewer records than `count`, it passes through unchanged.
- For CSV, the header row is always preserved.
- Useful for sampling large inputs, capping output size before an `http-transform`, or testing a pipeline end-to-end with real data.
