# TODO

## MVP

Roughly in priority order.

### Steps

- `filter` — add JSON/JSONL support
- `filter` — add basic operators (`equals`, `not-equals`, `contains`, `starts-with`)
- `filter` — add multi-condition support (`all` / `any`)
- `select` — add JSON/JSONL support
- `sort` — order records by a field (asc/desc)
- `dedupe` — remove duplicate records by key field
- `cast` — add CSV support

### CLI

- `--dry-run` — validate the pipeline YAML without consuming stdin

### Error handling

Better error messages with step name, field name, and record index. eg:

```
[step 4: filter] field 'country' not found in record
```

---

## Backlog

Lower priority ideas for after MVP.

### Steps

- `enrich` — add derived/computed fields using templates, eg: `"{{first_name}} {{last_name}}"`
- `map` — field-level numeric and string transforms (multiply, divide, add, subtract, round, to_lower, etc.)
- `mask` — partial redaction (eg: expose last 4 chars of credit card)
- `http-request` — send payload to HTTP endpoint without replacing it (fire-and-forget style)
- `http-transform` — add CSV payload support

### Step enhancements

- `normalize` — support pipe-separated strategy chains, eg: `trim|lower|collapse-spaces`
- `filter` — combined `all`/`any` nesting (see design note below)

### Payload / format

- CSV — configurable delimiter (eg: `delimiter: ";"`)
- JSONL — stricter validation (one object per line, reject arrays)

### Logging

- Replace `fmt.Printf` calls with a proper logger

---

## Notes

### Filter multi-condition design

For combined `all`/`any` nesting, model the step representation like this:

```go
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

Example YAML:

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
