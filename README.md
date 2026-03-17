# Command Builder Interactive

A terminal UI (TUI) for composing complex CLI commands through interactive forms.
Browse a searchable library of command templates, fill in the required fields, and
get a fully-assembled command printed to stdout — ready to run, pipe, or capture.

---

## Features

- 🔍 **Fuzzy search** across all command templates with relevance scoring
- 📝 **Interactive form** — fill in placeholders, preview the command live
- 📋 **Clipboard copy** — the built command is automatically copied to your clipboard on confirm
- ⚡ **Dynamic pickers** — inputs can run a shell command and show a live scrollable picker on Tab
- 🏷️ **Searchable tags** — find commands by alternative terms you might think of
- ⭐ **Starred commands** — save frequently-used commands with their values for quick re-use
- 🗂️ **Config manager** — add, edit, delete, export, and import config packs
- 🌐 **URL import** — pull shared command packs from any HTTP/HTTPS URL
- 📁 **Local file import** — load configs from local YAML files with tab-completion
- ✏️ **Built-in editor** — create and edit commands, including dynamic picker sources, without leaving the TUI
- ⚙️ **Settings** — customise the colour palette and toggle "run on enter" mode
- 🎨 **Themeable** — all colours configurable via ANSI 256 codes or hex values
- 🏷️ **Custom app name** — rename the app and optionally add a shell alias
- 🚀 **Run on Enter** — optionally execute the command directly from the TUI

---

## Installation

### Download a binary

Grab the latest binary for your platform from the
[Releases](https://github.com/dwilson2547/command-builder/releases) page:

```bash
# Linux x86-64
chmod +x command-builder-linux-amd64
mv command-builder-linux-amd64 ~/.local/bin/command-builder

# Linux ARM64
chmod +x command-builder-linux-arm64
mv command-builder-linux-arm64 ~/.local/bin/command-builder

# macOS (Apple Silicon)
chmod +x command-builder-darwin-arm64
mv command-builder-darwin-arm64 ~/.local/bin/command-builder

# macOS (Intel)
chmod +x command-builder-darwin-amd64
mv command-builder-darwin-amd64 ~/.local/bin/command-builder

# Windows — download command-builder-windows-amd64.exe or command-builder-windows-arm64.exe
# and place it on your PATH.
```

### Build from source

Requires **Go 1.24+**.

```bash
git clone https://github.com/dwilson2547/command-builder
cd command-builder
./build.sh
./command-builder
```

---

## Quick start

```bash
./command-builder
```

The search bar is focused immediately — start typing to find a command. Press **Enter**
to open the form, fill in the fields, and press **Enter** again to confirm. The built
command is printed to stdout and copied to your clipboard.

```bash
# Capture the output in a variable
CMD=$(./command-builder)
echo "$CMD"

# Execute it immediately
eval "$(./command-builder)"

# Redirect to a script
./command-builder > /tmp/run.sh && bash /tmp/run.sh
```

---

## Screens

### Search screen

The home screen. Type to filter results in real time.

```
╭──────────────────────────────────────────────────────────────────────╮
│  > docker                                                            │
╰──────────────────────────────────────────────────────────────────────╯
  docker › build-image      Build Docker image from Dockerfile      default
  docker › run-container    Run a container with port mapping       default
  docker › exec-shell       Open a shell in a running container     default

  /config to manage configs · /settings for settings             v1.18.0
```

**Slash commands:**

| Query                          | Effect                                     |
|--------------------------------|--------------------------------------------|
| `/default <terms>`             | Search only the built-in config            |
| `/all <terms>`                 | Search all configs (default behaviour)     |
| `/<config-name> <terms>`       | Search one specific config by name         |
| `/s <terms>`                   | Show starred commands (filter by terms)    |
| `/config`                      | Open the Config Manager                    |
| `/import <url-or-path>`        | Import a config immediately                |
| `/settings`                    | Open the Settings screen                   |

**Keyboard shortcuts:**

| Key       | Action                              |
|-----------|-------------------------------------|
| Type      | Update search query                 |
| ↑ / ↓     | Navigate results                    |
| PgUp/PgDn | Jump 10 results                     |
| Enter     | Open form for selected result       |
| Ctrl+C    | Quit                                |

---

### Form screen

Fill in the template placeholders and watch the command assemble live.

```
  docker › build-image
  Build Docker image from Dockerfile
  ──────────────────────────────────────────────────────────────────────

  * Image name   ╭──────────────────────────────────────────╮
                 │ myapp                                    │
                 ╰──────────────────────────────────────────╯

  * Tag          ╭──────────────────────────────────────────╮
                 │ latest                                   │
                 ╰──────────────────────────────────────────╯

    Dockerfile     ./Dockerfile

    Context        .

  ──────────────────────────────────────────────────────────────────────
  docker build -t myapp:latest -f ./Dockerfile .

  Tab/↑↓ navigate · Enter confirm · Esc back                   v1.18.0
```

- Fields marked `*` are **required** — the command cannot be confirmed until all are filled.
- `file` and `dir` fields support **Tab path completion**.
- Inputs with a `sub_command` configured show a **Tab picker** — press Tab to run the command and select a value from the live results.
- `flag` inputs are toggled with **Space**.
- On confirm, the command is printed to stdout and **copied to clipboard**.
- Press `*` to **star** the current command with its filled values for quick re-use later.

**Keyboard shortcuts:**

| Key       | Action                                            |
|-----------|---------------------------------------------------|
| Tab       | Next field / path completion / open sub-cmd picker|
| Shift+Tab | Previous field                                    |
| ↑ / ↓     | Previous/next field, completion, or picker item   |
| Space     | Toggle a `flag` input on/off                      |
| Enter     | Next field; confirm when all required fields done |
| *         | Star this command with current values             |
| Esc       | Return to search                                  |
| Ctrl+C    | Quit                                              |

---

### Config manager

Open with `/config` in the search bar. Manage all loaded command packs.

| Key | Action                                        |
|-----|-----------------------------------------------|
| ↑/↓ | Navigate the config list                      |
| n   | Create a new empty config                     |
| i   | Import a config from a URL                    |
| f   | Import a config from a local file             |
| e   | Edit the selected config in the command editor|
| d   | Delete the selected config                    |
| x   | Export the selected config to a file          |
| u   | Pull an update from the config's source URL   |
| Esc | Return to search                              |

---

### Command editor

Select a config and press **e** in the Config Manager to edit it. Navigate three
levels: **Commands → Options → Inputs**. Pressing **n** or **e** opens an inline
form; save with **Ctrl+S**.

When you save an Option, any `{{variable}}` placeholders in the template that
don't yet have a matching Input are **automatically created** as optional string
inputs, saving you manual setup.

The Input form includes a **SubCommand** field — enter any shell command whose
stdout (parsed as CSV) will populate a live Tab picker for that input. When the
SubCommand field is focused, press **Enter** to preview the picker output before
saving.

---

### Stars

Press `*` on any completed form to star the command with its current values.
Access starred commands from the search bar with `/s` — navigate to one and press
**Enter** to reopen it pre-filled, then adjust and confirm as normal.

| Key | Action                                 |
|-----|----------------------------------------|
| `/s <terms>` | Show starred commands (inline list) |
| Enter        | Open pre-filled form for selected star |
| d            | Delete selected star                   |
| Esc          | Exit star mode / return to search      |

---

### Settings

Open with `/settings`. Customise colours and behaviour:

- **Colour palette** — change any colour with an ANSI 256 code or hex value (`#ff8700`).
  Press **r** to reset a single colour, **R** to reset the entire palette.
- **Run on Enter** — when enabled, the TUI executes the built command in your shell
  instead of printing it.
- **App name** — rename the application. Optionally adds a shell alias to `~/.bashrc`.

Settings are persisted to `~/.config/command-builder/settings.json`.

---

## Config file format

Configs are YAML files placed in `~/.config/command-builder/configs/` (or imported
via URL/file). The built-in `default` config is embedded in the binary.

```yaml
name: "my-tools"
description: "Personal toolbox"
version: "1.0.0"
commands:
  - name: "pg"
    description: "PostgreSQL helpers"
    options:
      - name: "dump"
        description: "Dump a database to a file"
        template: "pg_dump -h {{host}} -U {{user}} -d {{database}} -f {{output_file}}"
        tags: ["postgres", "backup", "export"]
        inputs:
          - name: "host"
            type: "string"
            description: "Database host"
            required: true
            default: "localhost"
          - name: "user"
            type: "string"
            description: "Database user"
            required: true
          - name: "database"
            type: "string"
            description: "Database name"
            required: true
          - name: "output_file"
            type: "file"
            description: "Output .sql file path"
            required: true
```

**Input types:**

| Type     | Behaviour                                                        |
|----------|------------------------------------------------------------------|
| `string` | Plain text field                                                 |
| `file`   | Text field with Tab path completion (files)                      |
| `dir`    | Text field with Tab path completion (dirs)                       |
| `flag`   | Boolean toggle; Space to toggle; omitted from command when off   |

Any input type can also have a `sub_command` — a shell command that populates a
live scrollable picker when the user presses Tab. See
[`docs/sub-command-completions.md`](docs/sub-command-completions.md) for details.

See [`docs/config-format.md`](docs/config-format.md) for the full schema reference.

---

## Sharing configs

Host any YAML config at a public URL and share it:

```
/import https://raw.githubusercontent.com/yourname/my-configs/main/tools.yaml
```

Team members can then press **u** in the Config Manager to pull updates. See
[`docs/plugins.md`](docs/plugins.md) for details.

---

## Data locations

| Path                                          | Contents                        |
|-----------------------------------------------|---------------------------------|
| `~/.config/command-builder/configs/`          | User config YAML files          |
| `~/.config/command-builder/settings.json`     | Colour palette and preferences  |
| `~/.config/command-builder/stars.json`        | Starred commands and saved values |

---

## Development

```bash
# Run tests
go test ./...

# Build
./build.sh

# Lint (installs golangci-lint if absent)
./lint.sh
```

See [`CHANGELOG.md`](CHANGELOG.md) for release history.
