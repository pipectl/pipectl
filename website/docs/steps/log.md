# log

Prints a message, the current record count, and sample records to stdout. The payload passes through unchanged.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `message` | string | No | Text to print before the count and samples |
| `count` | boolean | No | Whether to print the record count. Defaults to `true`. |
| `sample` | integer | No | Number of sample records to print. Defaults to `10`. Set to `0` to disable. |

## Example

```yaml
- log:
    message: After normalization
    count: true
    sample: 5
```

## Notes

- `sample: 0` (or a negative value) disables sample output.
- Useful as a debugging checkpoint between steps.
- The payload is not modified — `log` is a read-only inspection step.
