# dedupe

Removes duplicate records, keeping the first occurrence of each unique key.

**Supported formats:** `json` (arrays only) `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | list | Yes | One or more field names that together form the uniqueness key |
| `case-sensitive` | boolean | No | When `false`, values are compared case-insensitively. Defaults to `true`. |

## Examples

Deduplicate by email address:

```yaml
- dedupe:
    fields: [email]
```

Deduplicate by full name, case-insensitively:

```yaml
- dedupe:
    fields: [first_name, last_name]
    case-sensitive: false
```

## Notes

- Records are compared by the combined values of all listed fields. If every listed field matches a previously-seen record, the record is dropped.
- Field order in `fields` does not affect which records are considered duplicates — only the combination of values matters.
- If a field is missing from a record, its value is treated as an empty string for comparison purposes.
- For CSV, the header row is always preserved.
- JSON object payloads (single record) are not supported — only JSON arrays.
