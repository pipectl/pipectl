# TODO

## Additional steps

- `enrich` — add derived/computed fields using templates, eg: `"{{first_name}} {{last_name}}"`
- `map` — field-level numeric and string transforms (multiply, divide, add, subtract, round, to_lower, etc.)
- `http-request` — send payload to HTTP endpoint without replacing it (fire-and-forget style)
s
## Step enhancements

- `normalize` — support pipe-separated strategy chains, eg: `trim|lower|collapse-spaces`
- `filter` — document or add `on-missing` option for records missing the filter field (currently silently excluded, which may surprise users)

## Documentation

- The CLI says "pipectl executes declarative pipelines defined in YAML or JSON." - JSON is not supported yet, remove it
- Website docs - should steps be alphabetical? Or some other order?

## CLI

### Additional CLI options

- [ ] `--output-format FORMAT` — override `output.format` from YAML at the CLI without editing the file
- [ ] `--from-step N` — skip steps 1–(N-1), start at step N using `--input` as the snapshot; useful for debugging expensive pipelines

### Paid feature

When a paid/cloud feature is included in a pipeline

```text
This step requires pipectl Cloud (secrets, scheduling, etc.)

Run with:
  pipectl run pipeline.yaml --cloud

Learn more: pipectl.dev/cloud
```