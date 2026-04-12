# Playbook: Add a New Step

## Goal

Add a new pipeline step end-to-end: YAML parsing, planning, execution, and tests.

## Implementation Checklist

1. Add spec type in `internal/pipeline/spec`.
2. Register the step key in `stepRegistry` in `internal/pipeline/spec/unmarshal.go`.
3. Add planner mapping in `internal/pipeline/plan/builder.go`.
4. Implement executable step in `internal/engine/steps/<stepname>/step.go`.
5. Add unit tests for both spec and engine step behavior.
6. Add a pipeline YAML for the step in `internal/pipeline/testdata/pipelines/step/<stepname>.yaml` and a corresponding golden output file in `internal/pipeline/testdata/golden/step/<stepname>.<ext>`.
7. Update `docs.md` with info about the new step.
8. Add a step page at `website/docs/steps/<stepname>.md` following the structure of existing step pages: supported formats, options table, examples, notes.

## Files Usually Touched

- `internal/pipeline/spec/<stepname>.go`
- `internal/pipeline/spec/unmarshal.go`
- `internal/pipeline/plan/builder.go`
- `internal/engine/steps/<stepname>/step.go`
- `internal/engine/steps/<stepname>/step_test.go`
- `internal/pipeline/testdata/pipelines/step/<stepname>.yaml`
- `internal/pipeline/testdata/golden/step/<stepname>.<ext>`
- `docs.md`
- `website/docs/steps/<stepname>.md`

## Spec Layer Requirements

- Define a config struct with yaml tags.
- Implement `StepType() string`.
- Keep config flat unless nesting is required.
- Validate required fields as early as possible.

## Engine Layer Requirements

- Implement `Name()`, `Supports()`, and `Execute()`.
- Return explicit errors; avoid panics.
- Keep step behavior isolated and deterministic.

## Test Matrix

- valid config unmarshals correctly
- unknown/missing required options fail predictably
- `Supports()` truth table is correct
- `Execute()` success case
- `Execute()` failure cases and error messages

## Done Criteria

- `go test ./...` passes.
- pipeline can parse and execute the new step.
- step is wired in spec registry and planner switch.
- testdata pipeline and golden file exist for the step.
- `docs.md` has a complete entry: supported payloads, all options with accepted values and defaults, notes for non-obvious behavior, and at least one example.
- `website/docs/steps/<stepname>.md` exists with supported formats, options table, examples, and notes.
- Step is listed in the table in `website/docs/steps/index.md`.
