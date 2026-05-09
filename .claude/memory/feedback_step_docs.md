---
name: Update README and sidebar when steps change
description: Adding or changing a step requires updates to README.md, website/docs/steps/index.md, and website/.vitepress/config.mts sidebar
type: feedback
---
When adding a new step or changing a step's user-visible behavior, update all of:
- `README.md` — step table row
- `website/docs/steps/index.md` — "All steps" table row
- `website/docs/steps/<stepname>.md` — full step documentation
- `website/.vitepress/config.mts` — sidebar `items` list under the Steps section (alphabetical order)

**Why:** The VitePress sidebar is manually maintained in `config.mts` and is separate from the steps index table. Adding a step doc file and updating the index table is not enough — the sidebar entry must also be added or the step won't appear in the left-hand navigation.

**How to apply:** Any time a new step is added, treat `website/.vitepress/config.mts` as a required update alongside the other doc files. The sidebar items list is under `themeConfig.sidebar`, in the `Steps` section.