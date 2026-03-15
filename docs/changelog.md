# Changelog

## [1.1.0] - 2026-03-15

### Maintenance

- **Added `.github/copilot-instructions.md`**
  Created Copilot instructions file documenting project overview, tech stack, architecture,
  code conventions, common workflows, and doc references to give GitHub Copilot consistent
  context across all sessions.

- **Centralized project versioning**
  Removed the hardcoded `const appVersion` from `search_screen.go`. Version is now defined
  in a single file (`internal/tui/version.go`) as exported `var AppVersion = "v1.1.0"` so
  it can be overridden at build time via `-ldflags`. Updated `build.sh` to automatically
  inject the version using `git describe --tags --always --dirty`, falling back to the
  hardcoded value when no git tag exists. All four TUI screens reference `AppVersion`.

### New Features

- **Edit existing configs from the Config Manager screen**
  Pressing `e` on any non-built-in config opens a new Edit screen that lets the
  user browse and mutate the full config hierarchy: Commands → Options → Inputs.
  Each level shows a navigable list (`↑`/`↓`, `Enter` to drill in, `Esc` to go
  back up). Within a level the user can create new items (`n`), edit the selected
  item (`e`), or delete it (`d`, with a `y`/`n` confirmation prompt). Forms are
  submitted with `Ctrl+S`. All changes are written to disk immediately via a new
  `Manager.UpdateConfig` method. Built-in (embedded) configs show an error when
  `e` is pressed rather than opening the editor.

- **Per-config source URL and one-key pull updates**
  `Config` now has an optional `source_url` YAML field. When a config is
  imported with `i` (import from URL) the URL is stored in the saved YAML file.
  Pressing `u` on any config that has a source URL prompts for confirmation
  (`yes`) and then re-fetches the YAML from that URL via a new
  `Manager.PullConfig` method, replacing the config's commands in place while
  preserving the local name and file path. Configs with a stored URL are marked
  with a `[url]` badge in the list. Configs without a source URL show a
  descriptive error instead of opening the prompt.

### Bug Fixes

- **Search input not responding to keystrokes**
  `Init()` used a value receiver, so `textinput.Focus()` mutated a temporary
  copy of the model that was immediately discarded. The stored model's input was
  never focused, causing all keystrokes to be silently dropped. Fixed by calling
  `ti.Focus()` directly in `NewSearchModel()` before the value is stored, and
  changing `Init()` to return `textinput.Blink` instead.

- **Search box disappearing when query returns no results**
  The view's `reserved` line count was `6`, but the actual chrome (title,
  rounded-border top, input, rounded-border bottom, hint, separator, status bar)
  occupies `7` lines. The resulting 1-line overflow caused the terminal renderer
  to clip the top of the view. On top of this, the "No results" message used
  `Padding(1,2)` (3 lines) and was written outside the bounded results area,
  adding further overflow when results were empty. Fixed by correcting `reserved`
  to `7` and rewriting the results/no-results block so the area always occupies
  exactly `visRows` lines.

- **Slash commands (`/config`, `/import`) breaking live search**
  The guard that prevented search from re-running during slash-command input was
  widened too broadly (blocking all `/`-prefixed queries), which broke the live-
  filter slash commands (`/default`, `/all`, `/<name>`). Reverted to only
  skipping `/config` and `/import`, which are the only commands handled solely
  on Enter.

- **Scroll position resetting every ~500ms**
  `textinput.Blink` fires a periodic tick message to toggle the cursor. This
  tick fell through to the text-input delegate block, which re-ran the search
  and unconditionally reset `selectedIdx` and `scrollTop` to `0` — snapping the
  list back to the top twice a second. Fixed by snapshotting the query value
  before delegating to the text input and only resetting scroll/selection when
  the query actually changes.

- **Search screen reverting to a fixed small width after navigating back**
  `NewSearchModel` hardcoded `ti.Width = 60` and left the model's `width` and
  `height` fields as zero. On first launch the initial `WindowSizeMsg` corrected
  the dimensions, but no `WindowSizeMsg` fires when navigating back from the
  form screen, so a fresh model created on back-navigation stayed at zero width
  for the rest of the session. Fixed by adding `w, h int` parameters to
  `NewSearchModel` and passing `a.width, a.height` from `AppModel` in the
  `backToSearchMsg` handler.
