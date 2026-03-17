# Issue: Search Input Not Responding to Keystrokes

## Symptom

Launching the application and typing in the search box produces no visible
output.  The results list remains unchanged regardless of what is typed.

## Root Cause

The bug lives in `internal/tui/search_screen.go`.

### Background – Bubble Tea value semantics

Bubble Tea models are ordinary Go structs passed by value.  `Init()`,
`Update()`, and `View()` all use **value receivers**, meaning every call
receives a *copy* of the model.  Any state mutation made inside `Init()` is
silently discarded because `Init()` only returns a `tea.Cmd`, not the updated
model.

### Why the text input was never focused

`charmbracelet/bubbles/textinput` ignores all keyboard events when its
internal `focus` field is `false` (the zero/default value):

```go
// textinput/textinput.go – line 556
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    if !m.focus {
        return m, nil   // drops every message, including key presses
    }
    …
}
```

To make the input accept keystrokes, `Focus()` must be called and its result
persisted.  `Focus()` has a **pointer receiver** (`func (m *Model) Focus()`),
so it mutates the model in place.

The original `Init()` implementation was:

```go
// BEFORE (broken)
func (m SearchModel) Init() tea.Cmd {
    return m.input.Focus()   // m is a copy; mutation is thrown away
}
```

Because `m` is a copy, `m.input.Focus()` sets `focus = true` only on that
throwaway copy.  The real `SearchModel` stored inside `AppModel` keeps
`input.focus == false` forever, so every keystroke is silently dropped by
the text input's `Update` method.

## Fix

**File:** `internal/tui/search_screen.go`

### 1 – Focus the input during construction

`NewSearchModel` creates the `textinput.Model` as a local variable before
embedding it in the struct.  Calling `Focus()` on this local variable
correctly mutates it in place (pointer receiver auto-addressing works
here), and the focused state is then copied into the returned `SearchModel`:

```go
// AFTER (fixed)
func NewSearchModel(mgr *config.Manager) SearchModel {
    ti := textinput.New()
    ti.Placeholder = "Search commands… (e.g. 'openssl print p12')"
    ti.Width = 60
    ti.Focus()   // sets focus = true on the local variable before copying

    m := SearchModel{mgr: mgr, input: ti}
    m.results = runSearch("", mgr)
    return m
}
```

### 2 – Simplify `Init()` to only start cursor blinking

Because focus is now guaranteed at construction time, `Init()` only needs to
return the blink command so the cursor animates:

```go
// AFTER (fixed)
func (m SearchModel) Init() tea.Cmd {
    return textinput.Blink
}
```

This removes the misleading and ineffective `m.input.Focus()` call from a
value-receiver method.

## Why the fix works

| Step | Before fix | After fix |
|---|---|---|
| `NewSearchModel` returns | `input.focus == false` | `input.focus == true` |
| Bubbletea stores `AppModel` | `search.input.focus == false` | `search.input.focus == true` |
| User presses a key | `textinput.Update` returns early (`!m.focus`) | `textinput.Update` processes the key, updates value |
| `runSearch` is called | always with `""` | with the current typed query |
| Results list | always shows all results | filters in real time |

## Affected Versions

All versions prior to this fix.