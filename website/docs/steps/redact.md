# redact

Replaces the values of selected fields with a redacted form.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | list | Yes | Field names to redact |
| `strategy` | string | No | Redaction strategy. Defaults to `REDACTED` if omitted. |

### Strategies

| Strategy | Output |
|----------|--------|
| `mask` | Replaces each character with `*` (e.g. `secret` → `******`) |
| `sha256` | Replaces the value with its SHA-256 hex digest |
| *(omitted)* | Replaces the value with the string `REDACTED` |

## Examples

Mask credit card numbers:

```yaml
- redact:
    fields: [credit_card]
    strategy: mask
```

Hash email addresses for analytics (stable, non-reversible):

```yaml
- redact:
    fields: [email, phone]
    strategy: sha256
```

Simple redaction with default placeholder:

```yaml
- redact:
    fields: [password, api_key]
```

## Notes

- For JSON and JSONL, only top-level string fields are redacted. Non-string values are left unchanged.
- For CSV, all field values in the specified columns are redacted regardless of type.
- `sha256` produces a stable deterministic hash — the same input always produces the same output.
