# Playbook: Add a Payload Format

## Goal

Introduce a new payload format (for example `jsonl`) with predictable read/write behavior.

## Implementation Checklist

1. Add payload type in `internal/engine/payload` implementing `Type() string`.
2. Extend `Read()` in `internal/engine/payload/payload.go`.
3. Extend `Write()` in `internal/engine/payload/payload.go`.
4. Ensure existing steps handle or reject the format via `Supports()`.
5. Add/extend tests for parse, serialize, and failure behavior.

## Files Usually Touched

- `internal/engine/payload/<format>.go`
- `internal/engine/payload/<format>_test.go`
- `internal/engine/payload/payload.go`
- `internal/engine/steps/*/step.go` (if support rules must change)

## Design Rules

- Keep payload structure explicit and minimal.
- Avoid hidden conversion behavior between unrelated formats.
- Prefer explicit "unsupported format" errors.
- Do not panic on malformed input.

## Test Matrix

- valid input parses into expected payload shape
- malformed input returns an error
- writer emits expected output format
- unsupported conversion paths fail clearly

## Done Criteria

- read/write path for the new format is covered by tests
- all affected step `Supports()` behavior is intentional and tested
- `go test ./...` passes
