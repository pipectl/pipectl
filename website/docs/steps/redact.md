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
| `partial-last` or `partial-last:N` | Reveals the last N characters (default 4), masks the rest with `*` (e.g. `partial-last:4` on `1234-5678-9012-3456` → `***************3456`) |
| `partial-first` or `partial-first:N` | Reveals the first N characters (default 4), masks the rest with `*` (e.g. `partial-first:4` on `1234-5678-9012-3456` → `1234***************`) |
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

PCI-DSS-style credit card display — last 4 digits only:

```yaml
- redact:
    fields: [credit_card]
    strategy: partial-last
```

Or with an explicit count:

```yaml
- redact:
    fields: [credit_card]
    strategy: partial-last:4
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
