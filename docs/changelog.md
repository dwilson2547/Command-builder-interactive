# Changelog

All notable changes to this project are documented here.

---

## [Unreleased]

### Bug Fixes

#### Search Input Not Responding to Keystrokes

**Commit:** `fa060c3`

**Files changed:**
- `internal/tui/search_screen.go`
- `docs/issues/search-not-responding.md` *(new)*

**Problem:**
Typing in the search box had no effect — the results list never filtered and
the cursor did not move.

**Root cause:**
`textinput.Focus()` was called inside `Init()`, which uses a value receiver.
The mutation set `focus = true` on a throwaway copy of the model; the real
`SearchModel` stored in `AppModel` kept `input.focus == false`.  Because
`charmbracelet/bubbles/textinput` silently drops all messages (including key
presses) when `focus == false`, every keystroke was discarded.

**Fix:**
- `NewSearchModel()` — call `ti.Focus()` on the local variable *before*
  embedding it in the struct, so the focused state is carried into the
  returned model.
- `Init()` — return `textinput.Blink` (cursor animation only); focus is now
  guaranteed at construction time.

See [`docs/issues/search-not-responding.md`](issues/search-not-responding.md)
for a full write-up.
