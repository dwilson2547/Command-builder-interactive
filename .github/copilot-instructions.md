# GitHub Copilot Instructions — Command Builder Interactive

## Project Overview

**Command Builder Interactive** is a terminal UI (TUI) application written in Go.
It lets users browse a searchable library of CLI command templates, fill in required
inputs via an interactive form, and receive the fully-assembled command as output.
Configs are YAML files that define commands, options, input placeholders, and templates.

---

## Changelog

- All notable changes must be documented in `CHANGELOG.md`.
- Follow [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.
- The project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
- **Increment the minor version** for every change unless the user specifies otherwise.
- Version is defined in `internal/tui/version.go` as `var AppVersion = "vX.Y.Z"`.
  At build time `build.sh` overrides it via `-ldflags "-X …tui.AppVersion=<tag>"` using `git describe`.

## README

- Location: `README.md` in the project root.
- Keep it updated whenever user-facing behaviour changes.

## Todo

- Location: `todo.md` in the project root.
- Check off items as they are completed (`- [x]`).
- After completing an item, also add a corresponding entry to `CHANGELOG.md`.

---

## Tech Stack

| Layer | Library / Tool |
|-------|----------------|
| Language | Go 1.24+ |
| TUI framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) (Elm-architecture) |
| TUI components | [Bubbles](https://github.com/charmbracelet/bubbles) (textinput, etc.) |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Config format | YAML (`gopkg.in/yaml.v3`) |
| Build | `build.sh` → produces `./command-builder` binary |

---

## Project Structure

```
command-builder-interactive/
├── main.go                    # Entry point; embeds default config; starts Bubbletea
├── go.mod / go.sum
├── configs/
│   └── default.yaml           # Embedded default config (built into binary via go:embed)
├── internal/
│   ├── config/
│   │   ├── types.go           # Config, Command, Option, Input structs
│   │   ├── loader.go          # YAML load/save helpers
│   │   ├── manager.go         # Manager: in-memory registry + disk persistence
│   │   └── config_test.go
│   ├── tui/
│   │   ├── app.go             # Root AppModel, screen routing, inter-screen messages
│   │   ├── search_screen.go   # Search / main screen
│   │   ├── form_screen.go     # Input form screen
│   │   ├── config_screen.go   # Config manager screen
│   │   ├── edit_screen.go     # Config editor screen
│   │   └── styles.go          # All Lipgloss styles (single source of truth)
│   ├── search/
│   │   ├── search.go          # Fuzzy/scored search logic
│   │   └── search_test.go
│   └── plugin/
│       └── plugin.go          # Thin wrapper: URL import via Manager
├── CHANGELOG.md
├── docs/
│   ├── config-format.md
│   ├── plugins.md
│   ├── usage.md
│   └── README.md
└── todo.md
```

---

## Architecture

### Bubbletea (Elm) Pattern

Every screen is a struct that implements `tea.Model` (`Init`, `Update`, `View`).
`AppModel` in `tui/app.go` owns all screen models and routes messages between them.

- Use **pointer receivers on Manager methods** but **value receivers on tea.Model** where Bubbletea expects it; be careful about where mutations happen.
- Inter-screen navigation uses typed message structs (e.g. `selectOptionMsg`, `backToSearchMsg`, `goToConfigMsg`).
- **Never mutate shared state inside `View()`** — it is called frequently and must remain pure.

### Config System

- `Config` → `[]Command` → `[]Option` → `[]Input`
- `Manager` is the only authorised way to add, update, delete, pull/refresh, or list configs.
- The embedded `default` config (compiled via `//go:embed`) has no `FilePath`; treat it as read-only.
- User configs live in `~/.config/command-builder/configs/` as individual `<name>.yaml` files.
- A config may store a `source_url` — if present, it can be refreshed with `Manager.PullConfig`.

### Input Types

| Type | Behaviour |
|------|-----------|
| `string` | Free-text input |
| `file` | File path (tab-autocomplete planned) |
| `dir` | Directory path |
| `flag` | Boolean toggle |

### Search

- Lives in `internal/search/`.
- Scores results across all loaded configs.
- Supports `/filter` prefixes: `/default`, `/all`, `/<config-name>` for scoping.
- Slash commands `/config` and `/import` are handled only on Enter (not as live filters).

---

## Code Conventions

### Go Style
- Standard Go formatting (`gofmt`/`goimports`).
- Package names are short, lowercase, no underscores.
- Exported types/functions have godoc comments.
- Error handling: return `error`; do not `panic` in library code.
- Avoid global mutable state outside `Manager`.

### TUI / Lipgloss
- **All styles are defined in `internal/tui/styles.go`** — never create ad-hoc inline styles elsewhere.
- Colors are named constants (`colorPrimary`, `colorAccent`, etc.) — use them instead of raw color codes.
- Terminal width/height are passed down from `AppModel` via `tea.WindowSizeMsg`; every screen must handle this message and resize accordingly.

### YAML Config Files
- Follow the schema in `docs/config-format.md`.
- Template placeholders use `{{variable_name}}` syntax.
- `required: true` fields must be filled before the form can be confirmed.

### Tests
- Unit tests live alongside source files (`_test.go`).
- Use the standard `testing` package; no third-party test framework.
- Run tests with `go test ./...`.

---

## Common Workflows

### Adding a new TUI screen
1. Create `internal/tui/<screen_name>_screen.go` implementing `tea.Model`.
2. Add a `screen<Name>` constant in `app.go`.
3. Add the model field to `AppModel` and route messages in `AppModel.Update`.
4. Define any new inter-screen message types near the top of `app.go`.
5. Add corresponding styles to `styles.go`.

### Adding a new Manager method
1. Add the method to `internal/config/manager.go`.
2. Write a test in `internal/config/config_test.go`.
3. Wire up any TUI-side message/command in the relevant screen file.

### Adding commands to the default config
- Edit `configs/default.yaml` following the schema in `docs/config-format.md`.
- The file is embedded at compile time; a rebuild is required to pick up changes.

### Building
```bash
./build.sh
```
Produces the `./command-builder` binary.

### Linting
Run the linter after making any code changes:
```bash
./lint.sh
```
Installs `golangci-lint` automatically if not present. Fix all reported issues before committing.

---

## Docs to Consult

| Question | Document |
|----------|----------|
| Config YAML schema | `docs/config-format.md` |
| Keyboard shortcuts & search filters | `docs/usage.md` |
| URL import / sharing | `docs/plugins.md` |
| Past changes | `CHANGELOG.md` |
