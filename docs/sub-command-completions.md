# Sub-Command Completions

Sub-command completions let any config input populate a dynamic, interactive picker
by running a shell command and parsing its output. This is useful for inputs whose
valid values are runtime-dependent — like Docker container names, running processes,
or database tables.

## How It Works

Add a `sub_command` field to any `Input` in your config YAML. When the user focuses
that input and presses **Tab**, the command runs, its CSV output is parsed, and a
scrollable picker appears below the input field.

## YAML Syntax

```yaml
- name: "container"
  type: "string"
  description: "Container name or ID"
  required: true
  sub_command: "docker ps --format '{{.Names}},{{.Image}}'"
```

The `sub_command` value is passed directly to `sh -c`, so shell quoting and
substitution work as expected.

## CSV Output Format

Each line of the command's stdout becomes one item in the picker:

```
column_0,column_1
```

| Column | Role | Required |
|--------|------|----------|
| 0 | **Value** — inserted into the input field when selected | Yes |
| 1 | **Detail** — displayed alongside the value for context; not inserted | No |

Lines are split on the **first comma only**, so a value containing commas is safe
as long as the detail (if any) is everything after that first comma.

Empty lines in the output are skipped.

### Example output for `docker ps --format '{{.Names}},{{.Image}}'`

```
my-nginx-container,nginx:latest
redis-1,redis:7-alpine
postgres-dev,postgres:15
```

The picker displays:

```
  my-nginx-container    nginx:latest
  redis-1               redis:7-alpine
❯ postgres-dev          postgres:15
```

Selecting `postgres-dev` inserts `postgres-dev` into the input field.

## Key Bindings

| Key | Action |
|-----|--------|
| **Tab** (picker hidden) | Run the sub-command and open the picker |
| **Tab** (picker open) | Close the picker |
| **Up / Down** | Navigate items in the picker |
| **Enter** | Select the highlighted item and fill the input |
| **Esc** | Close the picker without selecting |

When focused on an input that has a `sub_command`, the status bar shows a
`Tab: pick value` hint as a reminder.

## Error Handling

- If the command fails (non-zero exit) or produces no output, the picker displays
  a brief error message instead of an empty list.
- The input field remains editable — you can still type a value manually.

## Notes

- `sub_command` works with any input `type` (`string`, `file`, `dir`).
- The `{{...}}` placeholders in a `sub_command` value are **not** processed by the
  command-builder template engine — they are passed verbatim to the shell.
- Opening the sub-command picker dismisses any active file/directory completion
  overlay, and vice versa.
- Sub-commands run asynchronously so the UI stays responsive while the command
  executes.
