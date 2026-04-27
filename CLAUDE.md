# pipectl

## Project Overview

`pipectl` is a Go CLI that runs YAML-defined data pipelines.

A pipeline is an ordered list of steps operating on a payload. Each step reads
from `ExecutionContext.Payload`, transforms it or performs an action, then
passes control to the next step.

Primary goals:

- composable pipeline definitions
- predictable step behavior
- strongly typed payload handling
- minimal reflection
- small focused packages

## Source Map

Current implementation lives under `internal/`:

- `internal/engine`: runtime execution (`ExecutableStep`, `Engine`, context)
- `internal/engine/steps/*`: concrete step implementations
- `internal/engine/payload`: payload formats and read/write logic
- `internal/pipeline/spec`: YAML schema/types + step registry/unmarshal
- `internal/pipeline/plan`: compiles spec steps into executable engine steps
- `internal/pipeline`: top-level pipeline run orchestration
- `cmd/pipectl`: CLI entrypoint

## Boundaries

Keep dependency flow directional and acyclic:

```text
cmd -> internal/pipeline -> internal/pipeline/{spec,plan}
internal/pipeline/plan -> internal/engine + internal/engine/steps/*
internal/engine/steps/* -> internal/engine/payload
```

Rules:

- `spec` parses/validates pipeline config; it does not execute steps.
- `engine` executes steps; it does not parse YAML.
- step packages should not import `internal/pipeline/*`.

## Step Contract

All executable steps implement:

```go
type ExecutableStep interface {
  Execute(ctx *ExecutionContext) error
  Supports(payload payload.Payload) bool
  Name() string
}
```

Step guidelines:

- single responsibility per step
- deterministic behavior unless external I/O is the purpose
- minimal runtime state
- explicit errors with useful context

## Playbooks

Task-level implementation details live in `playbooks/`.

- add a new step: `playbooks/add-new-step.md`
- add a payload format: `playbooks/add-payload-format.md`
- add or extend step options: `playbooks/add-step-options.md`

Use these as checklists during implementation and review.

## Testing Expectations

Minimum expectations for feature work:

- table-driven unit tests where applicable
- happy path and failure path coverage
- parser/validation tests when YAML config is changed

Run:

```bash
go test ./...
```

## Change Quality Bar

Before considering work done:

1. Architecture boundaries are preserved.
2. New config surfaces are validated close to parse time.
3. Tests are added/updated and pass locally.
4. Testdata pipelines and golden files are added or updated to reflect any feature that affects pipeline execution — not just step changes. Each payload format the step supports should have its own pipeline test case (e.g. a JSON and a CSV variant). CLI-level features (e.g. `--var`) also require a testdata pipeline and golden file wired into the appropriate `Test*Pipelines` function in `internal/pipeline/integration_test.go`; extend the `testCase` struct if the harness needs new fields (e.g. `vars`). Generate golden files with `go test -run <TestName> -update` from the repo root.
5. Documentation is updated to match any user-visible changes:
   - **CLI flags** (`cmd/pipectl/run.go`): update `README.md`, `website/docs/cli.md`, and `website/docs/getting-started.md`.
   - **Step implementation** (`internal/engine/steps/*/step.go`): update `README.md` (step table), `docs.md`, and `website/docs/steps/<stepname>.md`.
   - **Step registry** (`internal/pipeline/spec/unmarshal.go`): if a new step was added, also update `website/docs/steps/index.md`.
   - **Payload formats** (`internal/engine/payload/`): if a new format was added, update `website/docs/formats.md` (or equivalent format reference page).

## Non-Goals

Do not introduce:

- framework-heavy abstractions
- unnecessary dependency injection layers
- plugin systems beyond current project scope
