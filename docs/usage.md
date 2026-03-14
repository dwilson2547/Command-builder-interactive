# Usage Guide

## Starting the app

```bash
./command-builder
```

The search bar is focused immediately. Start typing to find a command.

---

## Search screen

### Basic search

Type any part of a command name or description:

```
openssl p12
```

Results appear as you type, scored by relevance.

### Search filters

Prefix your query with a `/` modifier to narrow results:

| Query              | Behaviour                                     |
|--------------------|-----------------------------------------------|
| `/default <terms>` | Search only the built-in default config       |
| `/all <terms>`     | Search all configs (same as no prefix)        |
| `/<name> <terms>`  | Search a specific config by name              |

Examples:
```
/default tar
/my-tools pg dump
```

### Config management commands

| Query                            | Action                            |
|----------------------------------|-----------------------------------|
| `/config`                        | Open the config manager screen    |
| `/import https://example.com/x.yaml` | Import a config from a URL  |

### Keyboard shortcuts

| Key        | Action                          |
|------------|---------------------------------|
| Type       | Update search query             |
| ↑ / ↓      | Navigate results                |
| PgUp/PgDn  | Jump 10 results                 |
| Enter      | Open form for selected result   |
| Ctrl+C     | Quit                            |

---

## Form screen

When you select a result, the form screen shows all the inputs for that command option.

- **Required** fields are highlighted in pink with `*`.
- The current input shows a rounded border in accent colour.
- Navigate between fields with **Tab** / **Shift+Tab** or **↑** / **↓**.

### Tab completion (file / dir fields)

When the focused field has type `file` or `dir`:

1. Start typing a path (relative or absolute).
2. Press **Tab** to complete to the longest common prefix.
3. If multiple matches exist, a completion list appears — press **Tab** / **↑** / **↓** to cycle through them.

### Command preview

Once all required fields are filled, the built command appears at the bottom of the form.  
Template placeholders `{{...}}` that are not yet filled show `<placeholder_name>`.

### Keyboard shortcuts

| Key          | Action                                        |
|--------------|-----------------------------------------------|
| Tab          | Next field / path completion                  |
| Shift+Tab    | Previous field                                |
| ↑ / ↓        | Previous/next field or completion item        |
| Enter        | Next field; **confirm** when all required filled |
| Esc          | Go back to search                             |
| Ctrl+C       | Quit                                          |

When you confirm, the built command is printed to **stdout** after the TUI exits, so you can pipe or redirect it:

```bash
./command-builder > /tmp/cmd.sh
```

---

## Config manager screen

Open with `/config` in the search bar.

Shows a list of all loaded configs with command counts and badges.

| Key | Action                            |
|-----|-----------------------------------|
| ↑/↓ | Navigate config list              |
| n   | Create a new empty config         |
| d   | Delete selected config (confirm)  |
| x   | Export selected config to a file  |
| i   | Import a config from a URL        |
| Esc | Return to search                  |
