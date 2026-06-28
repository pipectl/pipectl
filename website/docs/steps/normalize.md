# normalize

Normalises string fields in the payload.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | map | Yes | Map of field name to normalisation strategy |

### Strategies

| Strategy | Effect |
|----------|--------|
| `lower` | Convert to lowercase |
| `upper` | Convert to uppercase |
| `trim` | Remove leading and trailing whitespace |
| `trim-left` | Remove leading whitespace only |
| `trim-right` | Remove trailing whitespace only |
| `collapse-spaces` | Replace runs of whitespace with a single space |
| `capitalize` | Capitalise the first letter of each word |

## Example

```yaml
- normalize:
    fields:
      email: trim|lower
      first_name: trim|capitalize
      last_name: trim|capitalize
      country: upper
      description: trim|collapse-spaces
```

## Notes

- Fields must exist in the payload. An error is returned if a configured field is missing.
- Only string values are normalised. Non-string fields that exist are left unchanged.
- Multiple strategies can be applied to a single field by separating them with `|` (e.g. `trim|lower`). Strategies are applied left-to-right.
