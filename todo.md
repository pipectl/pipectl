# TODO

## Refactoring

- **Move `sort` direction default to `UnmarshalYAML`** — `SortStep.Validate()` currently mutates `s.Direction = "asc"`; defaults belong in `UnmarshalYAML`, not the validator (same pattern `dedupe` already follows)
- **Move `log` sample default to spec** — the default `Sample=10` is set in `plan/builder.go:117`; it should live in `spec/log.go` `UnmarshalYAML` alongside the struct definition
- **Remove dead `String()` methods on spec step types** — every spec step type implements `String()` but it is never called; delete them

## Step enhancements

- `filter` — add `on-missing: exclude|include|error` option for records missing the filter field (currently silently excluded, which surprises users); default to `exclude` for backwards compatibility
- `filter` — add `case-sensitive: false` option for string operators (`equals`, `not-equals`, `contains`, `starts-with`, `ends-with`)
- `filter` — add `matches` operator for regex matching
- `sort` — multi-field sort: allow an ordered list of `field`/`direction` pairs as compound sort keys
- `sort` — add `nulls: first|last` option; currently nulls are always last with no way to override
- `cast` — add `on-error: fail|skip|default` strategy so a single unparseable value doesn't abort the whole pipeline
- `assert` — add value assertions (`field-equals`, `field-contains`, `field-matches`) alongside the existing `field-exists` and record-count checks
- `limit` — add `offset` option for "skip N, take M" patterns
- `redact` — support nested (non-top-level) fields; currently silently ignores them (code TODO at `redact/step.go:70`)
- `redact` — support non-string field types; currently silently skips them (code TODO at `redact/step.go:71`)

## Test coverage

- Add `http-request` integration test pipeline + golden file — no testdata pipeline exists; wire into `TestStepPipelines`
- Add `http-transform` integration test pipeline + golden file — same gap; mock HTTP server needed
- Add custom CSV delimiter integration test — `formats.md` documents the `delimiter:` option but no testdata pipeline exercises it
- Add `validate-json` integration test pipeline + golden file — only unit tests exist today
- Add cast nested-fields integration test — dot-notation (`a.b`) and array-index (`a[0]`) paths are documented but not covered end-to-end
- Add `count: message:` integration test — the `message` option has no testdata pipeline or golden file
- Add zero-record JSONL integration test — verify empty-input end-to-end behaviour
- Add `http-request` spec unmarshal tests — `unmarshal_test.go` covers every step except `http-request`
- Expand plan builder tests — `plan/builder_test.go` has no coverage for `filter`, `sort`, `dedupe`, `select`, `normalize`, `redact`, or `validate-json` plan compilation
- Add CLI flag tests — `cmd/pipectl/run_test.go` only tests `--output`; add coverage for `--verbose`, `--quiet`, `--dry-run`, `--timing`, and `--var`

## Documentation

- Clarify `http-transform` CSV support — spec allows `expect-format: csv` but step matrix shows ✗; decide if this is a supported path or a spec bug and fix whichever is wrong
- Add advanced examples — nested `all`/`any` filters, `--var` with multiple vars, `http-transform` chained with format conversion
- **Standardise step error-message phrasing** — steps use at least three different phrasings for the same class of error ("must be", "is required", "requires", "invalid"); pick one pattern and apply it consistently across all `spec/*.go` `Validate()` methods
- `redact` — add a concrete output example for the `sha256` strategy (lowercase hex digest) to `website/docs/steps/redact.md` and `DOCS.md`; no example exists today
- `dedupe` — add "(default: `true`)" to the `case-sensitive` option row in `website/docs/steps/dedupe.md`
- `filter` website docs — add a note that `greater-than` / `less-than` fail the pipeline if the field value is non-numeric (documented in `DOCS.md` but absent from `website/docs/steps/filter.md`)
- `cast` docs — document what happens when a value cannot be parsed as the target type (pipeline fails; no skip — cross-reference the planned `on-error` enhancement)

## Additional steps

- `enrich` — add derived/computed fields using templates, eg: `"{{first_name}} {{last_name}}"`
- `map` — field-level numeric and string transforms (multiply, divide, add, subtract, round, to_lower, etc.) — lower priority; most use cases covered by `enrich` once that exists

## CLI

### Additional CLI options

- [ ] `--output-format FORMAT` — override `output.format` from YAML at the CLI without editing the file
- [ ] `--from-step N` — skip steps 1–(N-1), start at step N using `--input` as the snapshot; useful for debugging expensive pipelines (lower priority — `--input` with a mid-pipeline snapshot covers most cases today)
