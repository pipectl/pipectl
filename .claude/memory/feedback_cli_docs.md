---
name: Update docs when CLI flags change
description: Whenever cmd/pipectl/run.go is modified to add/change/remove flags, update README.md, website/docs/cli.md, and website/docs/getting-started.md
type: feedback
---
When CLI flags in `cmd/pipectl/run.go` are added, changed, or removed, always update the documentation:
- `README.md` — the Usage/Flags block
- `website/docs/cli.md` — the Flags table, Examples section, and Notes
- `website/docs/getting-started.md` — any examples that show the relevant flag in use

**Why:** User explicitly asked for this after noticing the `--input` flag was implemented but not reflected in any docs.

**How to apply:** Treat doc updates as part of the same task as the flag change — not an afterthought.