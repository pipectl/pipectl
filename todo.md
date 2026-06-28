# TODO

## Additional steps

- `enrich` — add derived/computed fields using templates, eg: `"{{first_name}} {{last_name}}"`
- `map` — field-level numeric and string transforms (multiply, divide, add, subtract, round, to_lower, etc.) — lower priority; most use cases covered by `enrich` once that exists

## Step enhancements

- `filter` — add `on-missing: exclude|include|error` option for records missing the filter field (currently silently excluded, which surprises users); default to `exclude` for backwards compatibility

## Documentation

- Document `filter` silent exclusion behaviour — "missing fields are treated as non-matching" is in the code comment but not in the user-facing step docs
- Clarify `http-transform` CSV support — spec allows `expect-format: csv` but step matrix shows ✗; decide if this is a supported path or a spec bug and fix whichever is wrong
- Add advanced examples — nested `all`/`any` filters, `--var` with multiple vars, `http-transform` chained with format conversion

## CLI

### Additional CLI options

- [ ] `--output-format FORMAT` — override `output.format` from YAML at the CLI without editing the file
- [ ] `--from-step N` — skip steps 1–(N-1), start at step N using `--input` as the snapshot; useful for debugging expensive pipelines (lower priority — `--input` with a mid-pipeline snapshot covers most cases today)