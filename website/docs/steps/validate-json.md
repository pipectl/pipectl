# validate-json

Validates JSON or JSONL records against a JSON Schema. The pipeline fails immediately if any record does not conform.

**Supported formats:** `json` `jsonl`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `schema` | string | Yes | Path to a JSON Schema file, or an inline JSON Schema string |

## Examples

Reference a schema file:

```yaml
- validate-json:
    schema: ./schemas/customer.json
```

Inline schema:

```yaml
- validate-json:
    schema: |
      {
        "type": "object",
        "required": ["email", "name"],
        "properties": {
          "email": { "type": "string", "format": "email" },
          "name":  { "type": "string" }
        }
      }
```

## Notes

- Validation errors include the field path and the constraint that failed, making them easy to act on.
- Use `validate-json` early in a pipeline to fail fast before expensive steps like `http-transform`.
- Schema file paths are resolved relative to the current working directory.
