---
name: Update DOCS.md as part of feature work
description: Always update DOCS.md when completing any feature task — steps, options, CLI flags, config fields
type: feedback
---

Always update DOCS.md as part of completing a feature task — do not wait to be asked.

**Why:** CLAUDE.md item 4 of the Change Quality Bar explicitly requires it: "Public behavior is reflected in examples or docs when relevant." This was missed when adding JSON/JSONL support to the filter step, and again when adding the `--verbose` and `--dry-run` CLI flags.

**How to apply:** Any of the following changes require a DOCS.md update:
- New step added
- New option added to an existing step (including new operators, strategies, or payload support)
- New CLI flag added (e.g. `--verbose`, `--dry-run`, `-o`)
- New pipeline config field added (e.g. `input.delimiter`)
- Changed behavior or error conditions for existing features

Cover: the option name, accepted values, default, and a note about edge-case behavior where relevant.
