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

- [ ] `-q / --quiet` — suppress all output except final payload; scripting-friendly complement to `--verbose`
- [ ] `--output-format FORMAT` — override `output.format` from YAML at the CLI without editing the file
- [ ] `--timing` — print per-step table (duration, records in/out) to stderr after execution
- [ ] `--var KEY=VALUE` (repeatable) — substitute `${VAR}` tokens in pipeline YAML at parse time; enables reusable pipelines across environments
- [ ] `--from-step N` — skip steps 1–(N-1), start at step N using `--input` as the snapshot; useful for debugging expensive pipelines

### Paid feature

When a paid/cloud feature is included in a pipeline

```text
This step requires pipectl Cloud (secrets, scheduling, etc.)

Run with:
  pipectl run pipeline.yaml --cloud

Learn more: pipectl.dev/cloud
```