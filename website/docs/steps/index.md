# Steps

A step is a single transformation applied to the pipeline payload. Steps run in the order they are listed. If any step fails, the pipeline stops.

## All steps

<div class="steps-table">

| Step | What it does | JSON | JSONL | CSV |
|------|-------------|:----:|:-----:|:---:|
| [`assert`](./assert) | Assert record count or field existence | ✓ | ✓ | ✓ |
| [`cast`](./cast) | Convert field types — int, float, bool, string, time | ✓ | ✓ | |
| [`convert`](./convert) | Convert payload format — json ↔ jsonl ↔ csv | ✓ | ✓ | ✓ |
| [`count`](./count) | Print current record count to stdout | ✓ | ✓ | ✓ |
| [`dedupe`](./dedupe) | Remove duplicate records by key fields | ✓ | ✓ | ✓ |
| [`default`](./default) | Fill missing or empty fields with a default value | ✓ | ✓ | ✓ |
| [`filter`](./filter) | Keep records matching a condition or nested `all`/`any` group | ✓ | ✓ | ✓ |
| [`http-request`](./http-request) | Send payload to an HTTP endpoint as a side effect, continue with the same payload | ✓ | ✓ | ✓ |
| [`http-transform`](./http-transform) | POST/PUT/PATCH payload to an HTTP endpoint, continue with the response | ✓ | ✓ | |
| [`limit`](./limit) | Truncate to N records | ✓ | ✓ | ✓ |
| [`log`](./log) | Print record count and samples to stdout | ✓ | ✓ | ✓ |
| [`normalize`](./normalize) | Normalise string fields (lower, upper, trim, capitalize…) | ✓ | ✓ | ✓ |
| [`redact`](./redact) | Replace field values with `mask`, `sha256`, or `REDACTED` | ✓ | ✓ | ✓ |
| [`rename`](./rename) | Rename fields | ✓ | ✓ | ✓ |
| [`select`](./select) | Keep only the specified fields | ✓ | ✓ | ✓ |
| [`sort`](./sort) | Sort records by a field, ascending or descending | ✓ | ✓ | ✓ |
| [`validate-json`](./validate-json) | Validate records against a JSON Schema | ✓ | ✓ | |

</div>

<style>
.steps-table table th:nth-child(n+3),
.steps-table table td:nth-child(n+3) {
  width: 1%;
  white-space: nowrap;
}
</style>

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
