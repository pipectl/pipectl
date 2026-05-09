---
name: Update README when step implementations change
description: Edits to internal/engine/steps/*/step.go must include README.md step table updates
type: feedback
---
When a step implementation changes (new strategies, new behavior), update `README.md` (the step table) in addition to `website/docs/steps/<stepname>.md`.

**Why:** The README step table describes each step's strategies/capabilities and is user-facing. It's easy to forget since CLAUDE.md originally only listed `docs.md` and `website/docs/steps/` for step changes.

**How to apply:** For any change to `internal/engine/steps/*/step.go` that adds or changes user-visible behavior, update both:
- `README.md` — find the `| \`<stepname>\`` row in the step table and update the description
- `website/docs/steps/<stepname>.md` — full strategy/option documentation