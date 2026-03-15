# Command Builder

A fast, interactive terminal command builder powered by YAML config files.

## Features

- **Instant search** — start typing immediately to fuzzy-search all known commands
- **Form-based input** — fill in placeholders with a guided form; required fields are highlighted
- **Tab completion** — auto-complete file and directory paths
- **Config filters** — narrow searches with `/default`, `/all`, or `/<configname>`
- **Plugin system** — import shared command sets from any URL
- **Config management** — create, delete, export, and import configs without leaving the TUI

## Quick start

```bash
# Build
go build -o command-builder .

# Run
./command-builder
```

Type to search. Press **Enter** on a result to fill in its form. Press **Enter** again (once all required fields are filled) to print the built command to stdout.

## Usage guide

See [docs/usage.md](usage.md).

## Config format

See [docs/config-format.md](config-format.md).

## Plugin / sharing system

See [docs/plugins.md](plugins.md).
