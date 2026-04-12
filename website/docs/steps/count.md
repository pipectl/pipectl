# count

Prints the current record count to stdout. The payload passes through unchanged.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `message` | string | No | Text to print before the count |

## Example

```yaml
- count:
    message: Records ready for export
```

## Notes

- The payload is not modified — `count` is a read-only inspection step.
- Use `log` instead if you also want sample records printed alongside the count.
