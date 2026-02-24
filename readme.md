# Pipectl

## Steps

### Filter

Drop records based on conditions.

```yaml
- filter:
    field: status
    equals: active
```

### Select

Keep only certain fields.

```yaml
- select:
    fields: [id, email, created_at]
```

### Rename

Rename fields.

```yaml
- rename:
    mappings:
      firstName: first_name
      lastName: last_name
```

### Map

Transform a field.

```yaml
- map:
    field: email
    to_lower: true
```

```yaml
- map:
    field: price
    multiply_by: 1.1
```

### Default

Set missing values.

```yaml
- default:
    field: country
    value: AU
```

### Cast

Convert types.

```yaml
- cast:
    field: age
    type: int
```

### Validate schema

You already have validate_json — this is broader

```yaml
- validate_schema:
    required: [id, email]
    types:
      age: int
      email: string
```

### Dedupe

Remove duplicates

```yaml
- dedupe:
    by: email
```

### Mask

Different from redact.

```yaml
- mask:
    field: credit_card
    strategy: last4
```

### Enrich

Add derived fields

```yaml
- enrich:
    field: full_name
    value: "{{first_name}} {{last_name}}"
```