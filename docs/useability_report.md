# Usability Report — Command Builder Interactive

**Date:** 2026-03-18  
**Version reviewed:** latest (`copilot/add-usability-report` branch)  
**Scope:** All six TUI screens, the search engine, config management, form input, starred commands, and settings.

---

## 1. Executive Summary

Command Builder Interactive is a capable and well-structured terminal tool with an impressive feature set. The search, form-building, sub-command picker, path completion, starring, and theming systems all work well. This report identifies **usability gaps** — things a user might expect but that are currently absent, unclear, or harder than necessary — and proposes actionable improvements ranked by impact.

---

## 2. Discovery & Navigation

### 2.1 No contextual help / `?` key

**Finding:** There is no in-app help available at any screen. A first-time user landing on the search screen has no way to discover available slash commands (`/config`, `/settings`, `/import`, `/s`) or keyboard shortcuts without reading the external documentation.

**Impact:** High — users who install the binary without reading the README will not know how to reach the config manager or settings screen.

**Suggestion:** Add a `?` key (or `/help` slash command) that opens a scrollable help overlay listing all slash commands and key bindings for the current screen. The overlay can re-use the same completion-style pop-up already used for sub-command pickers, so no new UI infrastructure is needed.

---

### 2.2 No visual hint that slash commands exist

**Finding:** The search bar placeholder (`Search commands…`) gives no indication that typing `/` triggers special commands. A user might type `/config` by accident and only then discover the feature.

**Impact:** Medium — the most powerful navigation features are entirely invisible.

**Suggestion:** Change the placeholder to `Search or /command…` and add a one-line footer hint such as `/ for commands · ? for help` that appears whenever the search bar is empty and unfocused.

---

### 2.3 No breadcrumb or screen title

**Finding:** When a user navigates from search → config manager → editor → inputs level, there is no consistent header showing where they are. Each screen uses a different style of title display, making it easy to lose context.

**Impact:** Medium — particularly noticeable in the three-level editor (Commands → Options → Inputs).

**Suggestion:** Add a persistent, single-line breadcrumb bar at the top of every screen, e.g.:

```
Command Builder  ›  Config Manager  ›  docker  ›  build-image  ›  Inputs
```

---

### 2.4 No way to jump directly to a specific config's commands from search

**Finding:** Users can filter by config using `/<config-name> <terms>`, but there is no quick way to browse *all* commands in a specific config without typing a filter term. An empty `/myconfig` search falls back to no results.

**Impact:** Low-Medium — advanced users who manage many configs will notice this.

**Suggestion:** Treat `/<config-name>` with no following terms as "show all commands in this config", sorted alphabetically.

---

## 3. Form Screen

### 3.1 Required-field errors are not surfaced on attempted submit

**Finding:** When a user presses Enter to confirm a form that still has unfilled required fields, nothing happens — the cursor simply moves to the next field. There is no error message, no toast, and no visual change indicating *why* the form didn't submit.

**Impact:** High — users think the Enter key is broken or that the command was silently accepted.

**Suggestion:** Display a brief, dismissible inline message such as `⚠ 2 required field(s) still empty` in the status bar for 2–3 seconds whenever the user attempts to confirm an incomplete form.

---

### 3.2 No way to clear / reset all form fields at once

**Finding:** There is no keyboard shortcut to clear all inputs and start fresh. A user who accidentally fills in the wrong values must navigate to each field and delete the text individually.

**Impact:** Medium — particularly painful for commands with many inputs.

**Suggestion:** Add a `Ctrl+R` (reset) shortcut that clears all text inputs, resets flags to their defaults, and refocuses the first field. Show a brief confirmation hint before clearing.

---

### 3.3 No copy-to-clipboard shortcut for partial / in-progress commands

**Finding:** The command preview at the bottom of the form updates live, but there is no way to copy it to the clipboard *before* confirming. The clipboard copy only happens on final confirmation.

**Impact:** Low-Medium — users who want to tweak the command manually in their terminal cannot easily grab the partial preview.

**Suggestion:** Add a `Ctrl+Y` (yank) shortcut that copies the current preview text to the clipboard without exiting the form, accompanied by a brief "Copied!" flash in the status bar.

---

### 3.4 Password fields have no visibility toggle

**Finding:** Inputs whose names contain "password" are auto-masked. There is currently no way to temporarily reveal the typed text to verify it.

**Impact:** Medium — users frequently mistype passwords in terminal forms and have no way to check without retyping.

**Suggestion:** Add a `Ctrl+P` (peek) toggle to switch between masked and visible mode on the currently focused password field, similar to the "show/hide password" button on web forms.

---

### 3.5 No indication of total field count or progress

**Finding:** On forms with many inputs (some options have 10+ fields), the user sees the current field but has no sense of overall progress through the form.

**Impact:** Low — but improves perceived completion and reduces anxiety.

**Suggestion:** Add a counter such as `Field 3 / 9` to the status bar whenever there are more than three fields.

---

### 3.6 Sub-command picker shows no loading feedback on slow commands

**Finding:** When a sub-command takes more than a fraction of a second to execute, the overlay appears blank until results arrive. There is a `loadingSubCmd` flag in the model but no spinner is rendered while it is true.

**Impact:** Medium — users assume the picker broke and press Esc, missing the feature.

**Suggestion:** Render a simple spinner (e.g. `⣷ Loading…`) inside the picker overlay while `loadingSubCmd` is true.

---

### 3.7 No sub-command timeout

**Finding:** The sub-command runner uses `exec.Command` with no timeout. A shell command that hangs (e.g. a network call that never returns) will freeze the TUI indefinitely.

**Impact:** High — this can make the application appear completely broken.

**Suggestion:** Wrap sub-command execution in a `context.WithTimeout` of 10 seconds. If the timeout is reached, display `⚠ Picker timed out` in the overlay.

---

## 4. Search Screen

### 4.1 No empty-state message when search returns no results

**Finding:** When a query returns zero matches, the results area is simply blank. There is no "No results found" message or suggestion to try a different term.

**Impact:** Medium — users may not know whether the search ran at all.

**Suggestion:** Display a centered message such as:

```
No commands matched "dockerr".
Try a shorter term or /all to search across all configs.
```

---

### 4.2 Search history not preserved

**Finding:** Every time the user presses Esc from a form and returns to search, the search bar is cleared. If they want to refine a query or try the next result, they must retype the entire search term.

**Impact:** Medium — very common interaction; affects daily use.

**Suggestion:** Restore the previous search query (and cursor position) when returning from the form screen via Esc. The model already receives a `backToSearchMsg`; it could carry the previous query string.

---

### 4.3 No keyboard shortcut to quickly open the Config Manager or Settings

**Finding:** Reaching the config manager requires typing `/config` + Enter. There is no direct hotkey (e.g. `Ctrl+,` for settings or `Ctrl+M` for config manager) from the search screen.

**Impact:** Low — slash commands work fine once discovered, but discoverability is low (see 2.1).

**Suggestion:** Map `Ctrl+,` to settings and `Ctrl+M` to config manager as alternative shortcuts, documented in the footer.

---

## 5. Config Manager & Editor

### 5.1 Delete confirmation requires typing the full config name

**Finding:** Deleting a config requires the user to type the exact config name to confirm. While safe, this is friction-heavy for configs with long or complex names.

**Impact:** Low — safety is good; friction is the trade-off.

**Suggestion:** Accept `yes` (case-insensitive) as an alternative confirmation in addition to the config name, while keeping the name-based confirmation as the primary method.

---

### 5.2 No undo / cancel when editing a command or option

**Finding:** The editor opens an inline form for editing commands, options, and inputs. There is no way to discard edits in progress — pressing Esc while an inline form is open submits or ignores changes without warning.

**Impact:** Medium — accidental edits that are saved cannot be easily undone.

**Suggestion:** When the user presses Esc on a dirty form (one where any field differs from the original value), display a `Discard changes? (y/n)` prompt before returning to the list.

---

### 5.3 No validation that template placeholders match input names

**Finding:** A user can write a template `docker run {{image_name}}` and name an input `image` — the placeholder will never be replaced and the literal `{{image_name}}` appears in the output. This failure is silent.

**Impact:** High — produces broken commands with no error message.

**Suggestion:** When saving an option, compare all `{{variable}}` placeholders in the template against the names of the defined inputs. Warn the user if any placeholder has no matching input:

```
⚠ Template contains {{image_name}} but no input named "image_name" was found.
```

---

### 5.4 No bulk import of multiple configs at once

**Finding:** The `/import` command and the `f` (file) action in the config manager each handle one file at a time. There is no way to point at a directory of YAML files and import them all.

**Impact:** Low — niche use case but useful for initial setup.

**Suggestion:** Allow `/import <directory-path>` (or a dedicated `F` shortcut) to scan a directory for `*.yaml` files and import them all, reporting a summary (`3 configs imported, 1 failed`).

---

## 6. Starred Commands Screen

### 6.1 No search / filter within the stars screen

**Finding:** The `/s <terms>` filter in the search bar is the only way to filter starred commands. Once you are inside the stars screen (reached from `/s`), there is no further filter input.

**Impact:** Low-Medium — users with many starred commands cannot search within the stars screen itself.

**Suggestion:** Add an inline filter bar at the top of the stars screen (similar to the search bar on the main screen) that filters the star list as the user types.

---

### 6.2 Starred commands have no edit capability

**Finding:** A star saves the option name and filled values, but there is no way to edit a star after saving it. To change a pre-filled value, the user must re-open the form, make changes, and re-star it (creating a duplicate), then delete the old one.

**Impact:** Medium — stars are intended for frequent reuse; inability to edit them is limiting.

**Suggestion:** Add an `e` shortcut on the stars screen that opens the pre-filled form with all values editable; saving re-stars (overwrites) the same entry rather than creating a new one.

---

### 6.3 Duplicate star names are not flagged

**Finding:** Two stars can share the same custom display name with no warning. In the `/s` list they are visually indistinguishable.

**Impact:** Low — confusing in lists with many stars.

**Suggestion:** When saving a star with a custom name that already exists, display a prompt: `A star named "my-build" already exists. Overwrite? (y/n)`.

---

## 7. Settings Screen

### 7.1 Color changes require knowing ANSI codes or hex values

**Finding:** The color settings accept ANSI 256 codes (0–255) or hex strings (`#rrggbb`). A user who doesn't know these formats has no visual reference and cannot experiment intuitively.

**Impact:** Medium — the color customisation feature exists but is inaccessible to non-technical users.

**Suggestion:** Add a small live preview swatch next to each color field that updates as the user types. If the terminal supports it, display the color name alongside the hex/ANSI code (e.g. `#5f87ff → Steel Blue`).

---

### 7.2 No way to export / share a theme

**Finding:** Themes are stored in `settings.json` locally. There is no way to share a theme with another user or save it as a named preset.

**Impact:** Low — nice-to-have for community sharing.

**Suggestion:** Add an `x` (export) shortcut on the settings screen that writes a `theme.json` file to the current directory (or a user-specified path), and an `i` (import) shortcut to apply a saved theme file.

---

### 7.3 No confirmation before applying potentially invisible color changes

**Finding:** If a user sets the foreground text color to the same value as the background, all text may become invisible. There is no validation or "apply and preview" step.

**Impact:** Low — recoverable by editing `settings.json` manually, but alarming for non-technical users.

**Suggestion:** After applying a color change, briefly flash a test string in the new color so the user can see the result immediately before committing.

---

## 8. Accessibility & General UX Polish

### 8.1 No mouse support

**Finding:** All navigation is keyboard-only. While this is intentional for a TUI, some users expect to be able to click on a list item to select it or scroll with a mouse wheel.

**Impact:** Low — TUI users typically prefer keyboards; still worth noting.

**Suggestion:** Consider adding optional Bubble Tea mouse mode (`tea.WithMouseCellMotion()`). A future setting could toggle it on/off.

---

### 8.2 Long command previews are not word-wrapped

**Finding:** The live command preview at the bottom of the form is rendered as a single line. Very long commands (common with `kubectl`, `docker`, `openssl`) overflow the terminal width and are truncated.

**Impact:** Medium — users cannot see the full command they are about to run.

**Suggestion:** Wrap the command preview at terminal width and allow the preview area to expand up to a configurable maximum number of lines (e.g. 4). If the command is still too long, add a scroll indicator.

---

### 8.3 No clipboard feedback when terminal does not support it

**Finding:** The app copies the built command to the clipboard on confirm. If the terminal does not support clipboard access (e.g. SSH session without X forwarding), the copy silently fails.

**Impact:** Low — the command is still printed to stdout, so the user isn't blocked; they just don't know the copy failed.

**Suggestion:** Catch the clipboard error and display `Clipboard not available — command printed to stdout` in the status bar.

---

### 8.4 Quit confirmation is absent

**Finding:** `Ctrl+C` exits the application immediately with no confirmation, even if the user has unsaved changes in the config editor or a partially completed form.

**Impact:** Low-Medium — accidental Ctrl+C is a common terminal muscle-memory mistake.

**Suggestion:** If there are unsaved changes in the editor (`dirty` flag) or an in-progress form, display a one-line `Press Ctrl+C again to quit` message. A second Ctrl+C within 2 seconds exits. Otherwise proceed immediately.

---

### 8.5 Footer hints are not always contextual

**Finding:** The footer shows a fixed set of hints (`Ctrl+C quit`, etc.) that do not change based on the currently focused element or active mode. For example, when a sub-command picker is open, the footer still shows generic search hints rather than picker-specific navigation hints.

**Impact:** Low — experienced users will know, but the footer is the primary discoverability surface for shortcuts.

**Suggestion:** Make footer hints context-sensitive: update them whenever the active mode changes (normal form, picker open, completion list open, star naming mode).

---

## 9. Summary Table

| # | Area | Finding | Impact | Effort |
|---|------|---------|--------|--------|
| 2.1 | Navigation | No in-app help (`?` key) | High | Low |
| 2.2 | Navigation | Slash commands not discoverable | Medium | Low |
| 2.3 | Navigation | No breadcrumb / screen title | Medium | Medium |
| 2.4 | Navigation | `/config-name` with no terms returns empty | Low | Low |
| 3.1 | Form | Required-field error not shown on submit attempt | High | Low |
| 3.2 | Form | No reset all fields shortcut | Medium | Low |
| 3.3 | Form | No mid-form clipboard copy | Low | Low |
| 3.4 | Form | No password visibility toggle | Medium | Low |
| 3.5 | Form | No field progress indicator | Low | Low |
| 3.6 | Form | Sub-command picker shows no loading spinner | Medium | Low |
| 3.7 | Form | Sub-command has no execution timeout | High | Low |
| 4.1 | Search | No empty-state message | Medium | Low |
| 4.2 | Search | Search query cleared on back navigation | Medium | Low |
| 4.3 | Search | No direct hotkeys for Config/Settings | Low | Low |
| 5.1 | Config | Delete confirmation is friction-heavy | Low | Low |
| 5.2 | Config | No undo/cancel on edit | Medium | Medium |
| 5.3 | Config | Mismatched placeholders fail silently | High | Medium |
| 5.4 | Config | No bulk directory import | Low | Medium |
| 6.1 | Stars | No filter inside stars screen | Low | Low |
| 6.2 | Stars | Stars cannot be edited in place | Medium | Medium |
| 6.3 | Stars | Duplicate star names not flagged | Low | Low |
| 7.1 | Settings | Color input has no visual preview | Medium | Medium |
| 7.2 | Settings | No theme export/import | Low | Medium |
| 7.3 | Settings | No validation against invisible colors | Low | Low |
| 8.1 | General | No mouse support | Low | High |
| 8.2 | General | Long command previews truncated | Medium | Low |
| 8.3 | General | Silent clipboard failure | Low | Low |
| 8.4 | General | No quit confirmation with unsaved changes | Low | Low |
| 8.5 | General | Footer hints are not context-sensitive | Low | Medium |

---

## 10. Top Priorities

Based on impact and implementation effort, the following five items are recommended as the highest-value improvements:

1. **Sub-command timeout** (3.7) — prevents TUI from freezing; small code change with large reliability gain.
2. **Required-field error on submit attempt** (3.1) — most common source of user confusion; one status-bar message.
3. **Template placeholder validation** (5.3) — silent failure produces broken commands; catch-and-warn at save time.
4. **In-app help overlay** (2.1) — dramatically improves discoverability for new users.
5. **Restore search query on back navigation** (4.2) — affects every user on every session; minimal code change.
