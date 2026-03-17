# User Guide

A step-by-step guide to using **Command Builder Interactive** — a terminal UI
for composing complex CLI commands from interactive forms.

---

## Table of Contents

1. [Installation & launching](#1-installation--launching)
2. [The search screen](#2-the-search-screen)
3. [Building a command (form screen)](#3-building-a-command-form-screen)
4. [Config manager](#4-config-manager)
5. [Adding your own commands (editor)](#5-adding-your-own-commands-editor)
6. [Settings & colour customisation](#6-settings--colour-customisation)
7. [Stars](#stars)
8. [Tips & tricks](#7-tips--tricks)

---

## 1. Installation & launching

Download the binary for your platform from the
[Releases](https://github.com/dwilson2547/command-builder/releases) page and
make it executable:

```bash
chmod +x command-builder-linux-amd64
./command-builder-linux-amd64
```

Or for macOS:

```bash
# Apple Silicon
chmod +x command-builder-darwin-arm64
./command-builder-darwin-arm64

# Intel Mac
chmod +x command-builder-darwin-amd64
./command-builder-darwin-amd64
```

Windows users: download `command-builder-windows-amd64.exe` (or `arm64`), place it on your `PATH`, and run it from any terminal.

Or build from source:

```bash
git clone https://github.com/dwilson2547/command-builder
cd command-builder
./build.sh
./command-builder
```

The TUI opens immediately in your terminal.

---

## 2. The search screen

The search screen is the home screen. The text input is focused as soon as the
app launches — start typing to search.

```
╭──────────────────────────────────────────────────────────────────────╮
│  > tar                                                               │
╰──────────────────────────────────────────────────────────────────────╯
  tar › create-archive      Create a compressed tar archive         default
  tar › extract-archive     Extract a tar archive                   default
  tar › list-contents       List contents of an archive             default

  /config to manage configs · /settings for settings                 v1.9.0
```

### Basic search

Type any part of the command name, description, or tags:

```
pg dump
openssl certificate
docker run
```

Results are ranked by relevance and updated as you type.

### Searching by tags

Commands can have tags — alternative terms you might think of. For example,
searching `pfx` or `keystore` can surface the `openssl › print-p12` command
even though those words don't appear in its name. Tags appear in square
brackets in the editor but are invisible to the user at search time; they just
work.

### Search filters (slash prefixes)

Prefix your query with a `/` modifier to narrow the search scope:

| Query                        | Effect                                     |
|------------------------------|--------------------------------------------|
| `/default <terms>`           | Search only the built-in default config    |
| `/all <terms>`               | Search all configs (same as no prefix)     |
| `/<config-name> <terms>`     | Search one specific config by name         |
| `/s <terms>`                 | Show starred commands (filter by terms)    |
| `/config`                    | Open the Config Manager screen             |
| `/import <url-or-path>`      | Import a config immediately                |
| `/settings`                  | Open the Settings screen                   |

### Keyboard shortcuts — search screen

| Key        | Action                                     |
|------------|--------------------------------------------|
| Type       | Update search query                        |
| ↑ / ↓      | Navigate results                           |
| PgUp/PgDn  | Jump 10 results at a time                  |
| Enter      | Open the form for the selected result      |
| Ctrl+C     | Quit                                       |

---

## 3. Building a command (form screen)

Pressing **Enter** on a search result opens the form screen. Each placeholder
in the command template is shown as an input field.

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

  Tab/↑↓ navigate · Enter confirm · Esc back                   v1.9.0
```

- Fields marked with `*` are **required** — the command cannot be confirmed
  until all required fields have a value.
- The built command preview at the bottom updates in real time. Unfilled
  optional placeholders appear as `<placeholder_name>`.
- The focused field has a rounded accent-colour border.
- `flag` type inputs are toggled with **Space** (no text entry needed).

### Dynamic value pickers

Some inputs have a `sub_command` configured, shown with a `Tab: pick value`
hint in the status bar. Press **Tab** on such a field to run the command and
open a live scrollable picker:

```
  Container name  ╭──────────────────────────────────────────╮
                  │                                          │
                  ╰──────────────────────────────────────────╯
  ╭──────────────────────────────────────────────────────╮
  │ ▶ web-app       nginx:1.25                           │
  │   db            postgres:16                          │
  │   cache         redis:7                              │
  ╰──────────────────────────────────────────────────────╯
  ↑↓: navigate  Enter: select  Esc: close
```

Use **↑**/**↓** to highlight an entry, **Enter** to select it, and **Esc** or
**Tab** to dismiss without selecting.

### File and directory fields

Fields typed `file` or `dir` support **Tab completion**:

1. Start typing a path, e.g. `./out`
2. Press **Tab** — the path expands to the longest common match.
3. If multiple matches exist, a list appears. Keep pressing **Tab** or use
   **↑**/**↓** to cycle through them.

```
    Output file  ╭──────────────────────────────────────────╮
                 │ ./output/                                │
                 ╰──────────────────────────────────────────╯
                   ./output/backup.sql
                 › ./output/dump.sql
                   ./output/schema.sql
```

### Confirming the command

Once all required fields are filled, press **Enter** on the last field (or any
field if the rest are already valid). The TUI exits and prints the built
command to **stdout**, so you can capture it:

```bash
CMD=$(./command-builder)
echo "$CMD"       # docker build -t myapp:latest -f ./Dockerfile .
```

Or redirect it to a script:

```bash
./command-builder > /tmp/run.sh && bash /tmp/run.sh
```

### Starring a command

Once all required fields are filled, press **`*`** to star the command. You can
optionally enter a custom name for the star, or leave blank to use the default
`command › option` label. Starred commands are saved with all their current
values so you can quickly re-run them later.

Access stars with `/s` in the search bar.

### Keyboard shortcuts — form screen

| Key          | Action                                            |
|--------------|---------------------------------------------------|
| Tab          | Next field / advance path completion / open picker|
| Shift+Tab    | Previous field                                    |
| ↑ / ↓        | Previous/next field, completion, or picker item   |
| Space        | Toggle a `flag` input on/off                      |
| Enter        | Next field; confirm when all required fields done |
| *            | Star this command with current values             |
| Esc          | Return to search                                  |
| Ctrl+C       | Quit                                              |

---

## 4. Config manager

Type `/config` in the search bar and press **Enter** (or just type `/config`
and hit Enter) to open the Config Manager.

```
  Configs
  ──────────────────────────────────────────────────────────────────────
  › default        Built-in commands (openssl, tar, docker, …)  [built-in]  12 cmds
    my-tools       Personal toolbox                                           3 cmds
    work-scripts   Company scripts            [url]                           7 cmds

  ──────────────────────────────────────────────────────────────────────
  n new · i import URL · f import file · d delete · x export · e edit · u pull

                                                                      v1.9.0
```

### Importing a config from a URL

1. Press **i**.
2. Enter the URL of a remote YAML config file.
3. Press **Enter** — the config is fetched and added to
   `~/.config/command-builder/configs/`.

After importing, configs with a stored source URL show a `[url]` badge.
Press **u** on them to re-fetch the latest version.

### Importing a config from a local file

1. Press **f**.
2. Type the path to your YAML file. Press **Tab** to autocomplete paths.
3. Press **Enter** — the file is copied into the managed config directory.

You can also import from the search screen directly:

```
/import ~/my-tools.yaml
/import https://example.com/tools.yaml
```

### Creating a new empty config

Press **n**, type a name, and press **Enter**. An empty config is created and
you can immediately press **e** to open the command editor.

### Exporting a config

Press **x**, enter a destination path (Tab completes paths), and press
**Enter**. The config is written as a YAML file you can share or back up.

### Updating a config from its source URL

Press **u** on any config that shows a `[url]` badge. Type `yes` to confirm.
The commands are replaced with the latest version from the URL while the
local name and file path are preserved.

### Deleting a config

Press **d** and confirm with **y**. The file is removed from disk.

> **Note:** Deleting the built-in default config hides it for the current
> session. Delete the tombstone file
> `~/.config/command-builder/configs/.default-hidden` to restore it.

### Config manager keyboard shortcuts

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

## 5. Adding your own commands (editor)

Select a config in the Config Manager and press **e** to open the command
editor. The editor has three levels you can navigate: **Commands → Options →
Inputs**.

### Level 1 — Commands

```
  Editing: my-tools
  ──────────────────────────────────────────────────────────────────────
  Commands
  › docker          Docker container platform
    pg              PostgreSQL helpers
    git             Git shortcuts

  ──────────────────────────────────────────────────────────────────────
  Enter drill-in · n new · e edit · d delete · Esc back              v1.9.0
```

- Press **n** to create a new command (fill in name and description, submit
  with **Ctrl+S**).
- Press **e** on an existing command to rename/re-describe it.
- Press **Enter** to drill into a command's options list.

### Level 2 — Options

```
  Editing: my-tools › docker
  ──────────────────────────────────────────────────────────────────────
  Options
  › build-image    docker build -t {{image_name}}:{{tag}} …   [dockerfile, image]
    run-container  docker run --name {{name}} …

  ──────────────────────────────────────────────────────────────────────
  Enter drill-in · n new · e edit · d delete · Esc back              v1.9.0
```

Tags (if set) appear in square brackets after the template preview.

Press **n** or **e** to open the option form:

```
  ── New Option ──────────────────────────────────────────────────────────
  Name         ╭──────────────────────╮
               │ push-image           │
               ╰──────────────────────╯
  Description  ╭──────────────────────╮
               │ Push image to registry│
               ╰──────────────────────╯
  Template     ╭────────────────────────────────────────────────────────╮
               │ docker push {{registry}}/{{image_name}}:{{tag}}        │
               ╰────────────────────────────────────────────────────────╯
  Tags         ╭──────────────────────╮
               │ upload, publish      │
               ╰──────────────────────╯
  ── Ctrl+S to save · Esc to cancel ─────────────────────────────────────
```

- **Template**: Use `{{variable_name}}` placeholders — these become form
  fields when a user selects the option.
- **Tags**: Comma-separated alternative search terms. Users can find this
  option by typing these words.
- When you save with **Ctrl+S**, any `{{variable_name}}` placeholders in the
  template that don't already have a matching Input are **automatically
  created** as optional string inputs. You can then drill in to refine them.

### Level 3 — Inputs

After saving an option, press **Enter** on it to see its inputs list:

```
  Editing: my-tools › docker › push-image
  ──────────────────────────────────────────────────────────────────────
  Inputs
  › registry      string   Container registry URL
    image_name    string   Image name
    tag           string   Image tag

  ──────────────────────────────────────────────────────────────────────
  n new · e edit · d delete · Esc back                               v1.9.0
```

Press **e** on an input to edit it, or **n** to create one manually:

```
  ── Edit Input ──────────────────────────────────────────────────────────
  Name         ╭────────────────────╮
               │ registry           │
               ╰────────────────────╯
  Type         ╭────────────────────╮
               │ string             │    (string | file | dir | flag)
               ╰────────────────────╯
  Description  ╭────────────────────╮
               │ Container registry │
               ╰────────────────────╯
  Required     ╭────────────────────╮
               │ true               │    (true | false)
               ╰────────────────────╯
  Default      ╭────────────────────╮
               │ docker.io          │
               ╰────────────────────╯
  SubCommand   ╭────────────────────────────────────────────────────────╮
               │ docker ps --format '{{.Names}},{{.Image}}'             │
               ╰────────────────────────────────────────────────────────╯
               Enter to preview
  ── Ctrl+S to save · Esc to cancel ─────────────────────────────────────
```

**Input types:**

| Type     | Behaviour in the form screen                    |
|----------|-------------------------------------------------|
| `string` | Plain text field                                |
| `file`   | Text field with Tab path completion (files)     |
| `dir`    | Text field with Tab path completion (dirs only) |
| `flag`   | Boolean flag — Space to toggle; `default` value is inserted when on |

**SubCommand field:**

Enter any shell command whose stdout (CSV format) will populate a live picker when
the user presses Tab on this input in the form screen. While the SubCommand field
is focused in the editor, press **Enter** to preview the picker output immediately
without leaving the editor. Use **↑**/**↓** to browse the results and **Esc** to
close the preview.

Inputs with a SubCommand configured are marked with a `⚡` in the inputs list.

All changes are written to disk immediately on each **Ctrl+S** save — no
separate "commit" step is needed.

---

## 6. Settings & colour customisation

Type `/settings` in the search bar and press **Enter** to open the Settings
screen.

```
  Settings
  ──────────────────────────────────────────────────────────────────────
  › Primary      #5f87ff   ██  Commands, borders, title highlights
    Accent       #87d7ff   ██  Options, focused inputs, required fields
    Success      #87ff87   ██  Command previews, success messages
    Warning      #ffd787   ██  Completion overlays, warnings
    Error        #ff5f5f   ██  Error messages and banners
    Muted        #6c6c6c   ██  Descriptions, hints, separators
    Text         #d0d0d0   ██  Normal result rows and labels
    Selected BG  #262626   ██  Background of the selected row

  ──────────────────────────────────────────────────────────────────────
  ↑↓ navigate · e/Enter edit · r reset colour · R reset all · Esc back

                                                                      v1.9.0
```

### Changing a colour

1. Navigate to the colour you want to change with **↑**/**↓**.
2. Press **e** or **Enter** to enter edit mode.
3. Type a new value — either:
   - An ANSI 256-colour code: `208`
   - A CSS hex value: `#ff8700`
4. Press **Enter** to confirm. The new colour is applied immediately across
   all screens.

### Resetting colours

- Press **r** to reset the currently selected colour to its default.
- Press **R** to reset the entire palette to built-in defaults.

### Run on Enter

Toggle **Run on Enter** in the settings screen to execute the built command
directly in your current shell instead of printing it to stdout. When enabled,
the status bar shows `run command & quit` instead of `copy command & quit`.

### App name

Set a custom **App name** to personalise the header. You can also have the
settings screen append a shell alias to `~/.bashrc` so you can launch the app
by your chosen name.

Colour settings are saved to `~/.config/command-builder/settings.json` and
loaded automatically at startup.

---

## Stars

Starred commands let you save a command with all its current input values for
quick re-use.

### Starring a command

1. Fill in the form for any command.
2. Press **`*`** — optionally enter a custom name, or press **Enter** to use
   the default `command › option` label.
3. The star is saved to `~/.config/command-builder/stars.json`.

### Using starred commands

Type `/s` in the search bar (or `/s <term>` to filter):

```
╭──────────────────────────────────────────────────────────────────────╮
│  > /s docker                                                         │
╰──────────────────────────────────────────────────────────────────────╯
  ★ docker › build-image   myapp:latest                          ⭐ star
  ★ docker › exec-shell    web-app  /bin/bash                    ⭐ star
```

Press **Enter** on a starred entry to reopen its form pre-filled with the
saved values. Adjust any fields and confirm to run or copy the command.

Press **d** to delete a star.

| Key         | Action                               |
|-------------|--------------------------------------|
| `/s <term>` | Show/filter starred commands         |
| ↑ / ↓       | Navigate stars                       |
| Enter       | Open pre-filled form                 |
| d           | Delete selected star                 |
| Esc         | Exit star mode / return to search    |

---

## 7. Tips & tricks

### Pipe the output into a script

```bash
# Build a docker run command and execute it immediately
eval "$(./command-builder)"
```

### Use a custom config for a team

1. Write a YAML config file (see [config-format.md](config-format.md)).
2. Host it on a web server or internal file share.
3. Each team member runs `/import https://your-server/tools.yaml` once.
4. Everyone can pull updates with the `u` key in the Config Manager.

### Quick import from the search bar

```
/import ~/Downloads/new-tools.yaml
```

No need to open the Config Manager — just type and press **Enter**.

### Back up your configs

All user configs live in `~/.config/command-builder/configs/`. Copy this
directory to back up or transfer your setup to another machine.

### Restore the built-in default config

If you deleted the default config and want it back:

```bash
rm ~/.config/command-builder/configs/.default-hidden
```

Restart the app and the built-in default will reappear.
