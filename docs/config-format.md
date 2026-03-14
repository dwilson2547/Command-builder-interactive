# Config File Format

Command Builder loads YAML config files from `~/.config/command-builder/configs/`.  
The built-in `default` config is embedded in the binary.

## Top-level fields

```yaml
name: "my-tools"          # unique identifier used in search filters
description: "My tools"   # shown in the config manager
version: "1.0.0"          # semver string (informational)
commands:                  # list of Command objects (see below)
  - ...
```

## Command

```yaml
- name: "docker"           # short slug (no spaces recommended)
  description: "Docker container platform"
  options:
    - ...
```

## Option

```yaml
- name: "build-image"
  description: "Build Docker image from Dockerfile"
  template: "docker build -t {{image_name}}:{{tag}} -f {{dockerfile}} {{context}}"
  inputs:
    - ...
```

The `template` field uses `{{input_name}}` placeholders that are replaced with the
values the user types in the form.

## Input

```yaml
- name: "image_name"
  type: "string"          # "string" | "file" | "dir" | "flag"
  description: "Image name"
  required: true
  default: ""             # pre-filled default value (optional)
```

### Input types

| Type     | Behaviour                                      |
|----------|------------------------------------------------|
| `string` | Plain text field                               |
| `file`   | Text field with **Tab** path completion        |
| `dir`    | Text field with **Tab** directory completion   |
| `flag`   | Boolean flag; leave blank to omit              |

## Full example

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
