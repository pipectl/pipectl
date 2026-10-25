# redact

Replaces the values of selected fields with a redacted form.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | list | Yes | Field paths to redact. Nested JSON fields use dot notation (`user.email`) and array indexing (`contacts[0].phone`). |
| `strategy` | string | No | Redaction strategy. Defaults to `REDACTED` if omitted. |

### Strategies

| Strategy | Output |
|----------|--------|
| `mask` | Replaces each character with `*` (e.g. `secret` → `******`) |
| `sha256` | Replaces the value with its SHA-256 hex digest |
| `partial-last` or `partial-last:N` | Reveals the last N characters (default 4), masks the rest with `*` (e.g. `partial-last:4` on `1234-5678-9012-3456` → `***************3456`) |
| `partial-first` or `partial-first:N` | Reveals the first N characters (default 4), masks the rest with `*` (e.g. `partial-first:4` on `1234-5678-9012-3456` → `1234***************`) |
| `nullify` | Replaces the value with `null`. Works on any type, so numbers, bools, objects, and arrays can be redacted without changing them to strings. |
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

Nested fields — paths that contain `[` must be quoted inside a YAML flow list:

```yaml
- redact:
    fields: [user.email, "contacts[0].phone"]
    strategy: mask
```

Redact a numeric salary without changing its type (`85000` → `null`):

```yaml
- redact:
    fields: [salary]
    strategy: nullify
```

Mask a number by casting it to a string first (`12345678` → `"****5678"`):

```yaml
- cast:
    fields:
      account_number:
        type: string
- redact:
    fields: [account_number]
    strategy: partial-last:4
```

## Notes

- For JSON and JSONL, every strategy except `nullify` only applies to string values. If a field holds a number, bool, object, or array, the step fails rather than silently changing its type — `cast` it to string first, or use `nullify`.
- Values that are already `null` are left unchanged, whatever the strategy.
- A field path that doesn't exist in a record fails the step.
- For CSV, columns are matched by exact header name (e.g. a flattened `user.email` header) and every value is redacted as a string. `nullify` writes an empty value, since CSV has no null.
- `sha256` produces a stable deterministic hash — the same input always produces the same output.
