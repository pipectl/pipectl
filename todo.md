# TODO

## MVP

Roughly in priority order.

### Steps

- `dedupe` — remove duplicate records by key field

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
- `filter` — document or add `on-missing` option for records missing the filter field (currently silently excluded, which may surprise users)

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
