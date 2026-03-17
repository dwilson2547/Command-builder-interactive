# Config File Format

This guide explains how to write a Command Builder config file from scratch.
By the end you will be able to build configs with static text, optional flags,
file/directory pickers, and dynamic value pickers powered by live shell commands.

---

## Quick start

A config file is a YAML file placed in `~/.config/command-builder/configs/`.
The filename becomes the config's disk identity; the `name` field inside the
file is what appears in the UI.

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

Save this as `~/.config/command-builder/configs/my-tools.yaml`, reopen the
app, and search for **pg dump** to try it immediately.

---

## Structure overview

```
Config
 └── commands[]          ← groups related operations (e.g. "docker", "git")
      └── options[]      ← individual commands (e.g. "build-image", "clone")
           ├── template  ← the shell command with {{placeholder}} variables
           └── inputs[]  ← one input per placeholder
```

---

## Top-level fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✓ | Unique identifier used in `/name` search filters |
| `description` | | One-line summary shown in the config manager |
| `version` | | Semver string (informational only) |
| `source_url` | | URL the config was fetched from; enables `/config` → Refresh |
| `commands` | ✓ | List of Command objects |

---

## Command

A command groups related options under a common tool name (e.g. `docker`,
`git`, `openssl`).

```yaml
commands:
  - name: "ffmpeg"
    description: "Audio and video converter"
    options:
      - ...
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✓ | Short slug — no spaces recommended |
| `description` | | Shown in search results |
| `options` | ✓ | List of Option objects |

---

## Option

An option is one specific use-case of a command.

```yaml
options:
  - name: "convert-video"
    description: "Re-encode a video file"
    template: "ffmpeg -i {{input}} -c:v libx264 -crf {{crf}} {{output}}"
    tags: ["encode", "h264", "mp4"]
    inputs:
      - ...
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✓ | Short slug shown in the results list |
| `description` | | Longer description shown under the name |
| `template` | ✓ | Shell command with `{{placeholder}}` variables |
| `tags` | | Extra search keywords (aliases, abbreviations) |
| `inputs` | | List of Input objects (empty list = no form, just run the template) |

### Template placeholders

Every `{{name}}` in the template must have a matching input with the same
`name`. The app replaces each placeholder with the value the user enters:

- **Required input, filled** → replaced with the entered value
- **Required input, empty** → replaced with `<name>` (shown in red in preview)
- **Optional input, empty** → placeholder and any surrounding whitespace removed
- **Flag input, toggled on** → replaced with the `default` string
- **Flag input, toggled off** → removed entirely

---

## Input

```yaml
inputs:
  - name: "output_file"   # must match the {{placeholder}} name exactly
    type: "file"          # "string" | "file" | "dir" | "flag"
    description: "Output path"
    required: true
    default: ""
    sub_command: ""
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✓ | Matches `{{name}}` in the template |
| `type` | ✓ | Controls the input widget — see table below |
| `description` | | Shown as the input field's label/placeholder |
| `required` | | If `true`, the command cannot be confirmed until this is filled |
| `default` | | Pre-filled value when the form opens |
| `sub_command` | | Shell command that populates a dynamic picker on **Tab** |

### Input types

| Type | Widget | Tab key |
|------|--------|---------|
| `string` | Plain text field | No action |
| `file` | Text field | Path completion (files + dirs) |
| `dir` | Text field | Path completion (dirs only) |
| `flag` | Toggle checkbox | No action — use **Space** to toggle |

---

## Example 1 — required string and file inputs

The simplest form: a few text fields, one of which is a file path.

```yaml
name: "dev-tools"
description: "Developer utilities"
version: "1.0.0"
commands:
  - name: "openssl"
    description: "Certificate and key utilities"
    options:
      - name: "gen-rsa-key"
        description: "Generate an RSA private key"
        template: "openssl genrsa -out {{output_file}} {{key_size}}"
        inputs:
          - name: "output_file"
            type: "file"
            description: "Path to write the .pem key"
            required: true
            default: "private.pem"
          - name: "key_size"
            type: "string"
            description: "Key size in bits"
            required: true
            default: "4096"
```

**What the user sees:**
- Two text fields: *output_file* (pre-filled `private.pem`) and *key_size* (pre-filled `4096`)
- Tab on *output_file* opens a filesystem browser
- The live preview updates as they type: `openssl genrsa -out private.pem 4096`

---

## Example 2 — optional inputs

Optional inputs are omitted from the command entirely when left blank.
Use them for flags and arguments the user may or may not want to include.

```yaml
      - name: "curl-request"
        description: "Make an HTTP request with curl"
        template: "curl -X {{method}} {{headers}} {{data}} {{url}}"
        inputs:
          - name: "method"
            type: "string"
            description: "HTTP method"
            required: true
            default: "GET"
          - name: "url"
            type: "string"
            description: "Target URL"
            required: true
          - name: "headers"
            type: "string"
            description: "Extra headers, e.g. -H 'Authorization: Bearer token'"
            required: false   # ← omitted when blank
          - name: "data"
            type: "string"
            description: "Request body, e.g. -d '{\"key\":\"value\"}'"
            required: false   # ← omitted when blank
```

If the user fills in only `method` (GET) and `url` (https://example.com), the
preview collapses the gaps automatically:

```
curl -X GET https://example.com
```

If they also fill `headers` with `-H 'Accept: application/json'`, the output becomes:

```
curl -X GET -H 'Accept: application/json' https://example.com
```

> **Tip:** Put optional inputs *before* required ones in the template only if
> the CLI tool accepts them in that position. The whitespace collapsing handles
> the gaps, but order still matters for the final command.

---

## Example 3 — boolean flags

Use `type: "flag"` for options that are either present or absent in the
command (like `-v`, `--dry-run`, `--force`).

The `default` field holds the **string inserted when the flag is on**.
The placeholder is removed entirely when the flag is off.

```yaml
      - name: "rsync-deploy"
        description: "Sync a local directory to a remote server"
        template: "rsync -av{{dry_run}}{{delete}} {{source}}/ {{user}}@{{host}}:{{remote_path}}"
        inputs:
          - name: "dry_run"
            type: "flag"
            description: "Dry run — show changes without applying (--dry-run)"
            required: false
            default: " --dry-run"   # ← note the leading space
          - name: "delete"
            type: "flag"
            description: "Delete files on remote that no longer exist locally (--delete)"
            required: false
            default: " --delete"
          - name: "source"
            type: "dir"
            description: "Local directory to sync"
            required: true
          - name: "user"
            type: "string"
            description: "Remote SSH user"
            required: true
          - name: "host"
            type: "string"
            description: "Remote hostname or IP"
            required: true
          - name: "remote_path"
            type: "string"
            description: "Destination path on the remote server"
            required: true
```

**How the flags compose:**

| dry_run | delete | Rendered fragment |
|---------|--------|-------------------|
| off | off | `rsync -av /src/ user@host:/dest` |
| on | off | `rsync -av --dry-run /src/ user@host:/dest` |
| off | on | `rsync -av --delete /src/ user@host:/dest` |
| on | on | `rsync -av --dry-run --delete /src/ user@host:/dest` |

> **Convention for flag `default` values:**
> - For long flags (`--dry-run`), include a leading space: `" --dry-run"`
> - For short single-char flags embedded mid-word (like `tar -czvf`), omit
>   the space so they concatenate directly: `default: "z"`

### Embedded short flags

Some tools use combined short flags (e.g. `tar -czvf`). Use flags without
leading spaces and place the placeholders immediately adjacent in the template:

```yaml
      - name: "create-archive"
        description: "Create a compressed tar archive"
        template: "tar -cv{{gzip}}{{bzip2}}{{xz}}f {{output_file}} {{source}}"
        inputs:
          - name: "gzip"
            type: "flag"
            description: "gzip compression (-z)"
            required: false
            default: "z"
          - name: "bzip2"
            type: "flag"
            description: "bzip2 compression (-j)"
            required: false
            default: "j"
          - name: "xz"
            type: "flag"
            description: "xz compression (-J)"
            required: false
            default: "J"
          - name: "output_file"
            type: "file"
            description: "Output archive file"
            required: true
          - name: "source"
            type: "dir"
            description: "Source directory or file"
            required: true
```

If the user toggles **gzip on** and leaves bzip2/xz off:
`tar -cvzf archive.tar.gz ./project`

If the user toggles nothing (plain tar):
`tar -cvf archive.tar ./project`

---

## Example 4 — dynamic value picker with `sub_command`

Any input can include a `sub_command` field containing a shell command.
When the user presses **Tab** on that input, the app runs the command, parses
its stdout as CSV, and shows a scrollable picker.

```yaml
      - name: "exec-container"
        description: "Open a shell inside a running Docker container"
        template: "docker exec -it {{container}} {{shell}}"
        inputs:
          - name: "container"
            type: "string"
            description: "Running container name"
            required: true
            sub_command: "docker ps --format '{{.Names}},{{.Image}}'"
            #                        ↑ col 0 = value    ↑ col 1 = detail (display only)
          - name: "shell"
            type: "string"
            description: "Shell to open"
            required: false
            default: "/bin/bash"
```

When the user presses Tab on the *container* field, the app runs
`docker ps --format '{{.Names}},{{.Image}}'`, which outputs something like:

```
web-app,nginx:1.25
db,postgres:16
cache,redis:7
```

The picker shows:

```
▶ web-app       nginx:1.25
  db            postgres:16
  cache         redis:7
```

Selecting `web-app` fills the input; `nginx:1.25` is display-only.

**CSV format rules:**
- One entry per line
- Column 0 → value inserted into the field
- Column 1 (optional) → display detail shown in the picker but not inserted
- Lines are split on the **first comma only**, so values may contain commas
- Blank lines are ignored

**Picker key bindings:**

| Key | Action |
|-----|--------|
| **Tab** (field focused) | Run command and open picker |
| **Up / Down** | Navigate items |
| **Enter** | Select item and fill the input |
| **Esc** or **Tab** (picker open) | Close without selecting |

> **Any shell command works.** You are not limited to Docker.
> Examples: `kubectl get pods -o name`, `git branch`, `ls *.sql`, or a custom
> script that queries a database and prints `id,description` lines.

---

## Example 5 — putting it all together

This config combines required fields, optional fields, boolean flags, a
`dir` input with Tab completion, and a `sub_command` picker.

```yaml
name: "k8s-tools"
description: "Kubernetes helpers"
version: "1.0.0"
commands:
  - name: "kubectl"
    description: "Kubernetes CLI helpers"
    options:
      - name: "exec"
        description: "Open a shell in a running pod"
        template: "kubectl exec -it {{pod}} {{namespace}} {{container}} -- {{shell}}"
        tags: ["shell", "bash", "debug", "pod"]
        inputs:
          - name: "pod"
            type: "string"
            description: "Pod name"
            required: true
            sub_command: "kubectl get pods --no-headers -o custom-columns=':metadata.name,:status.phase'"

          - name: "namespace"
            type: "string"
            description: "Namespace flag, e.g. -n production"
            required: false
            # Optional: leave blank to use the current context's default namespace

          - name: "container"
            type: "string"
            description: "Container flag, e.g. -c sidecar (for multi-container pods)"
            required: false

          - name: "shell"
            type: "string"
            description: "Shell binary"
            required: false
            default: "/bin/sh"

      - name: "port-forward"
        description: "Forward a local port to a pod port"
        template: "kubectl port-forward {{pod}} {{local_port}}:{{pod_port}} {{namespace}}"
        tags: ["proxy", "tunnel", "forward"]
        inputs:
          - name: "pod"
            type: "string"
            description: "Pod name"
            required: true
            sub_command: "kubectl get pods --no-headers -o custom-columns=':metadata.name,:status.phase'"

          - name: "local_port"
            type: "string"
            description: "Local port"
            required: true
            default: "8080"

          - name: "pod_port"
            type: "string"
            description: "Pod port"
            required: true
            default: "80"

          - name: "namespace"
            type: "string"
            description: "Namespace flag, e.g. -n production"
            required: false

      - name: "apply"
        description: "Apply a manifest file"
        template: "kubectl apply{{dry_run}}{{server_side}} -f {{manifest}}"
        tags: ["deploy", "create", "update"]
        inputs:
          - name: "dry_run"
            type: "flag"
            description: "Dry run (--dry-run=client)"
            required: false
            default: " --dry-run=client"

          - name: "server_side"
            type: "flag"
            description: "Server-side apply (--server-side)"
            required: false
            default: " --server-side"

          - name: "manifest"
            type: "file"
            description: "Path to the YAML manifest"
            required: true
```

---

## Tips and common mistakes

### Placeholder names must match exactly

The `name` field of an input and the `{{name}}` in the template are
case-sensitive and must be identical.

```yaml
# ✓ Correct
template: "ssh {{user}}@{{host}}"
inputs:
  - name: "user"   # matches {{user}}
  - name: "host"   # matches {{host}}

# ✗ Wrong — mismatch causes the placeholder to remain literally in the output
template: "ssh {{User}}@{{Host}}"
inputs:
  - name: "user"
  - name: "host"
```

### Optional inputs leave gaps — use spacing carefully

When an optional input is empty its placeholder is removed, along with
adjacent whitespace. Design your templates so the command still reads
correctly when any subset of optional inputs is empty.

```yaml
# ✓ Good — each optional flag carries its own spacing inside the default value
template: "docker run{{ports}}{{env_file}} {{image}}"
inputs:
  - name: "ports"
    type: "flag"
    default: " -p 8080:80"   # leading space included in the value
  - name: "env_file"
    type: "flag"
    default: " --env-file .env"

# ✗ Risky — if 'ports' is empty you get "docker run --env-file .env image"
# which is fine here, but double-check your specific tool accepts that order
```

### Tags make commands easier to find

Tags are extra search keywords. Use abbreviations, synonyms, and common
misspellings that users might type.

```yaml
tags: ["k8s", "kube", "pod", "shell", "debug", "bash"]
```

### Test incrementally

Start with the template string alone (no inputs) and confirm it looks right
in the preview, then add inputs one by one.

