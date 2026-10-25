---
name: no-git-command-suggestions
description: "User uses IntelliJ for all git operations, do not suggest git commands"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: ca0660c1-3c56-4d8c-9c90-d04e7667cbd8
  modified: 2026-10-25T04:44:31.518Z
---

Do not suggest git commands (commit, push, add, etc.) to the user.

**Why:** User manages all git operations through IntelliJ IDE and finds the suggestions unnecessary.

**How to apply:** When changes are ready to commit/push, simply say the changes are ready — don't include the git commands. This is about not telling the user to type commands themselves — it does not mean avoid git entirely. Commits/pushes stay in IntelliJ unless the user says otherwise.