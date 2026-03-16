# TODO

### Payload

- TODO
    - Support JSONL
    - Support arrays of JSON objects

## Steps

### Convert

Convert the payload to a different format.

```yaml
- convert:
    format: json | jsonl | csv
```

Future enhancements:

```yaml
- convert:
    format: json
    pretty: true
```

```yaml
- convert:
    format: csv
    delimiter: ";"
```

Conversions:

| In    | Out   |
|-------|-------|
| CSV   | JSON  |
| CSV   | JSONL |
| JSON  | JSONL |
| JSONL | JSON  |
| JSON  | CSV   |
| JSONL | CSV   |

### Select

- TODO
    - Add support for JSON

### Filter

- Add support for JSON payloads
- Add support for multiple conditions.

AND example:

```yaml
- filter:
    all:
      - field: country
        equals: AU
      - field: status
        equals: active
```

OR example:

```yaml
- filter:
    any:
      - field: country
        equals: AU
      - field: country
        equals: NZ
```

Combination:

```yaml
- filter:
    all:
      - field: age
        greater_than: 18
      - any:
          - field: country
            equals: AU
          - field: country
            equals: NZ
```

Note: if doing the combined above, model the step representation like this:

```go
package steps

type ConditionGroup struct {
	All  []ConditionGroup
	Any  []ConditionGroup
	Rule *Rule
}

type Rule struct {
	Field string
	Value interface{}
}
```

### HTTP Transform

- Add support for posting CSV payloads

### HTTP request

- Add a separate step for HTTP requests
- Does NOT transform the payload (the same payload is passed through)
- Sends the payload to the HTTP endpoint
- Fails on non 200 responses

### Map

Transform a field.

Note: Some overlap with `normalize`.

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

- TODO
    - Which operations are supported? eg:
        - `to_lower`
        - `to_upper`
        - `multiply_by`
        - `divide_by`
        - `add`
        - `subtract`
        - `round`
        - `floor`
        - `ceil`

### Cast

Convert types.

```yaml
- cast:
    field: age
    type: int
```

- TODO
    - Which casts are supported?

### Mask

Different from redact.

```yaml
- mask:
    field: credit_card
    strategy: last4
```

- TODO
    - Add support for CSV
    - Add support for JSON

### Dedupe

Remove duplicates

```yaml
- dedupe:
    by: email
```

### Enrich

Add derived fields

```yaml
- enrich:
    field: full_name
    value: "{{first_name}} {{last_name}}"
```