---
name: Confirm doc changes with file:line references
description: After any task that includes documentation updates, explicitly list each file and line number changed so the user can verify without hunting
type: feedback
---
At the end of any task that touches documentation files, list every doc file that was changed with its path and approximate line number or section name.

**Why:** User checked for the --quiet doc update and couldn't find it, even though it had been written. Explicit file:line references let them verify quickly.

**How to apply:** In the end-of-task summary, include a bullet per doc file: path, line number or section, and what was added/changed. Don't just say "docs updated" — show exactly where.