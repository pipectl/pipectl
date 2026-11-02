# sort

Sorts records by one or more fields.

**Supported formats:** `json` (arrays only) `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `field` | string | Yes, unless `keys` is set | Field name to sort by |
| `direction` | string | No | `asc` or `desc`. Defaults to `asc`. Not allowed with `keys`. |
| `nulls` | string | No | `first` or `last`. Where records with a missing, `null`, or empty (CSV) value are placed. Defaults to `last`. With `keys`, this is the default for any key that doesn't set its own `nulls`. |
| `keys` | list | Yes, unless `field` is set | Ordered list of sort keys for a compound sort. Each entry takes `field` (required), `direction` (default `asc`), and `nulls` (default: the step-level `nulls`). Cannot be combined with `field` or `direction`. |

## Examples

Sort by last name ascending (default):

```yaml
- sort:
    field: last_name
```

Sort by date descending:

```yaml
- sort:
    field: created_at
    direction: desc
```

Put records with no `shipped_at` value first:

```yaml
- sort:
    field: shipped_at
    nulls: first
```

Sort by country, then by age (oldest first) within each country:

```yaml
- sort:
    keys:
      - field: country
      - field: age
        direction: desc
```

Compound sort where every key puts nulls first except `age`:

```yaml
- sort:
    nulls: first
    keys:
      - field: region
      - field: age
        nulls: last
```

## Notes

- With `keys`, records are compared on the first key; later keys are only used to break ties left by earlier ones.

- Records with a missing or `null` field value are sorted last by default. Set `nulls: first` to place them first instead. Null placement is independent of `direction`.
- For CSV, empty string values are treated as missing.
- The sort is stable: records that tie on every key (including records with missing values) keep their original relative order.
- When field values parse as numbers, numeric ordering is used. Otherwise, string ordering is used.
- JSON object payloads (single record) are not supported — convert to a JSON array or JSONL first.
