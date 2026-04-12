# Steps

A step is a single transformation applied to the pipeline payload. Steps run in the order they are listed. If any step fails, the pipeline stops.

## All steps

| Step | What it does | JSON | JSONL | CSV |
|------|-------------|:----:|:-----:|:---:|
| [`normalize`](./normalize) | Normalise string fields (lower, upper, trim, capitalize…) | ✓ | ✓ | ✓ |
| [`filter`](./filter) | Keep records matching a condition or nested `all`/`any` group | ✓ | ✓ | ✓ |
| [`redact`](./redact) | Replace field values with `mask`, `sha256`, or `REDACTED` | ✓ | ✓ | ✓ |
| [`cast`](./cast) | Convert field types — int, float, bool, string, time | ✓ | ✓ | |
| [`default`](./default) | Fill missing or empty fields with a default value | ✓ | ✓ | ✓ |
| [`select`](./select) | Keep only the specified fields | ✓ | ✓ | ✓ |
| [`rename`](./rename) | Rename fields | ✓ | ✓ | ✓ |
| [`sort`](./sort) | Sort records by a field, ascending or descending | ✓ | ✓ | ✓ |
| [`limit`](./limit) | Truncate to N records | ✓ | ✓ | ✓ |
| [`dedupe`](./dedupe) | Remove duplicate records by key fields | ✓ | ✓ | ✓ |
| [`convert`](./convert) | Convert payload format — json ↔ jsonl ↔ csv | ✓ | ✓ | ✓ |
| [`validate-json`](./validate-json) | Validate records against a JSON Schema | ✓ | ✓ | |
| [`assert`](./assert) | Assert record count or field existence | ✓ | ✓ | ✓ |
| [`http-transform`](./http-transform) | POST/PUT/PATCH payload to an HTTP endpoint, continue with the response | ✓ | ✓ | |
| [`log`](./log) | Print record count and samples to stdout | ✓ | ✓ | ✓ |
| [`count`](./count) | Print current record count to stdout | ✓ | ✓ | ✓ |

## Step syntax

Each step is a single-key object in the `steps` list:

```yaml
steps:
  - normalize:
      fields:
        email: lower
  - redact:
      fields: [password]
      strategy: mask
  - filter:
      field: country
      equals: AU
```

The key is the step type. The value is the step's configuration. Each list item must contain exactly one step type.
