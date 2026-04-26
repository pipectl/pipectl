# TODO

## Additional steps

- `enrich` — add derived/computed fields using templates, eg: `"{{first_name}} {{last_name}}"`
- `map` — field-level numeric and string transforms (multiply, divide, add, subtract, round, to_lower, etc.)
- `http-request` — send payload to HTTP endpoint without replacing it (fire-and-forget style)
- `http-transform` — add CSV payload support

## Step enhancements

- `normalize` — support pipe-separated strategy chains, eg: `trim|lower|collapse-spaces`
- `filter` — document or add `on-missing` option for records missing the filter field (currently silently excluded, which may surprise users)

## CLI

### Additional CLI options

- `validate` - validate the pipeline is valid
- others?

### Paid feature

When a paid/cloud feature is included in a pipeline

```text
This step requires pipectl Cloud (secrets, scheduling, etc.)

Run with:
  pipectl run pipeline.yaml --cloud

Learn more: pipectl.dev/cloud
```