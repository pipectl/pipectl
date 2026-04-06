# Playbook: Add or Extend Step Options

## Goal

Add new options to an existing step without breaking existing pipelines.

## Implementation Checklist

1. Extend the spec struct in `internal/pipeline/spec/<step>.go`.
2. Add validation/defaulting behavior in spec layer as needed.
3. Wire new fields through planner in `internal/pipeline/plan/builder.go`.
4. Update engine step in `internal/engine/steps/<step>/step.go`.
5. Add tests for backward compatibility and new behavior.
6. Update `docs.md`: add the new option to the step's Options table, update Notes if behavior is non-obvious, and add an example if a new operator or strategy is introduced.

## Compatibility Rules

- Existing valid YAML should continue to work unless a breaking change is intentional.
- New options should be optional by default where possible.
- Breaking changes require explicit documentation and test updates.

## Validation Rules

- Validate required combinations early (spec layer preferred).
- Return actionable errors that identify the offending option.
- Keep option schema flat and readable unless complexity is unavoidable.

## Test Matrix

- previous config shape still works
- new option changes behavior as intended
- invalid option values fail with clear errors
- planner correctly maps new fields from spec to step implementation

## Done Criteria

- step supports old and new configs as designed
- behavior and validation are covered by unit tests
- `go test ./...` passes
- `docs.md` reflects the new option, its accepted values, default, and any edge-case behavior
