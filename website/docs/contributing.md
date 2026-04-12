# Contributing

Contributions are welcome — bug reports, new steps, documentation improvements, or fixes.

## Project layout

```
pipectl/
├── cmd/pipectl/              CLI entrypoint (cobra)
├── internal/
│   ├── engine/               Runtime execution
│   │   ├── payload/          JSON, JSONL, CSV payload types
│   │   └── steps/            One package per step implementation
│   └── pipeline/
│       ├── spec/             YAML schema types + step registry
│       ├── plan/             Compiles spec steps into engine steps
│       └── integration_test.go
├── examples/                 Runnable example pipelines
├── playbooks/                Task-level implementation guides
└── website/                  VitePress documentation site
```

### Dependency boundaries

```
cmd → internal/pipeline → internal/pipeline/{spec,plan}
internal/pipeline/plan → internal/engine + internal/engine/steps/*
internal/engine/steps/* → internal/engine/payload
```

- `spec` parses and validates YAML — it does not execute steps.
- `engine` executes steps — it does not parse YAML.
- Step packages must not import `internal/pipeline/*`.

## Running tests

```bash
go test ./...
```

Regenerate golden files after changing step output:

```bash
go test -run TestStepPipelines -update
```

## Adding a new step

The detailed checklist lives in [`playbooks/add-new-step.md`](https://github.com/pipectl/pipectl/blob/main/playbooks/add-new-step.md). The short version:

1. Add spec type in `internal/pipeline/spec/`
2. Register it in `internal/pipeline/spec/unmarshal.go`
3. Add the engine step in `internal/engine/steps/<name>/`
4. Wire it in `internal/pipeline/plan/builder.go`
5. Add unit tests in the step package
6. Add integration test pipeline(s) and golden files in `internal/pipeline/testdata/`
7. Wire the test case into `TestStepPipelines`
8. Add a step page under `website/docs/steps/`

Each payload format a step supports (JSON, JSONL, CSV) should have its own integration test case.

## Adding a new payload format

See [`playbooks/add-payload-format.md`](https://github.com/pipectl/pipectl/blob/main/playbooks/add-payload-format.md).

## Quality bar

Before submitting a pull request:

- Architecture boundaries are preserved
- New config surfaces are validated close to parse time
- Tests are added or updated and pass locally
- Testdata pipelines and golden files are added or updated
- The relevant docs page under `website/docs/` is updated

## Opening a pull request

Target the `main` branch. Include a brief description of what the PR does and why.
