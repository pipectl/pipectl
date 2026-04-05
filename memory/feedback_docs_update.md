---
name: Update DOCS.md as part of feature work
description: Always update DOCS.md when completing a feature task, without waiting to be asked
type: feedback
---

Always update DOCS.md as part of completing a feature task — do not wait to be asked.

**Why:** CLAUDE.md item 4 of the Change Quality Bar explicitly requires it: "Public behavior is reflected in examples or docs when relevant." This was missed when adding JSON/JSONL support to the filter step.

**How to apply:** After implementing a step change (new payload support, new options, new step), update the corresponding section in DOCS.md before considering the task done.
