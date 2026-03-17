# Changelog

All notable changes to Command Builder are documented here.

## [v1.37.0] - 2026-03-16

### Added

- **system command set** (`configs/system.yaml`) — 34 options across 5 command groups, covering all System abstraction items:
  - **services**: `list`, `status`, `start`, `stop`, `restart`, `reload`, `enable`, `disable`, `logs` — service name fields Tab-pick from live `systemctl list-units` output
  - **processes**: `kill-on-port`, `show-on-port`, `kill-by-pid`, `kill-high-cpu`, `kill-high-memory`, `kill-by-name` — PID picker auto-fills top 30 processes sorted by CPU; name picker lists running process names
  - **usb**: `list`, `list-verbose`, `tree`, `watch`
  - **network-devices**: `list-interfaces`, `interface-info`, `interface-stats`, `list-connections`, `wifi-list` — interface fields Tab-pick from `ip link show`
  - **network-info**: `ip-addresses`, `public-ip`, `routing-table`, `default-gateway`, `dns-config`, `arp-table`, `open-ports`

## [v1.36.0] - 2026-03-16

### Added

- **go command set** (`configs/go-tools.yaml`) — 29 options across 6 command groups:
  - **run**: `run` (race flag)
  - **build**: `build`, `build-all`, `install`, `clean`
  - **test**: `test`, `test-run`, `test-coverage`, `test-coverage-html`, `benchmark`, `test-count`
  - **modules**: `init`, `tidy`, `get`, `get-latest`, `download`, `vendor`, `graph`, `why`, `list-modules`
  - **quality**: `fmt`, `vet`, `generate`
  - **doc**: `doc`, `godoc-server`
  - **profiling**: `test-cpu-profile`, `test-mem-profile`, `pprof`
  - Package paths Tab-pick from `go list ./...`; module names Tab-pick from `go list -m all`.

## [v1.35.0] - 2026-03-16

### Added

- **sed command set** (`configs/sed.yaml`) — 24 options across 5 command groups:
  - **substitute**: `replace-first`, `replace-all`, `replace-nth`, `replace-in-range`, `replace-in-matching-lines`, `replace-extended` (ERE)
  - **lines**: `print-matching`, `print-range`, `delete-matching`, `delete-not-matching`, `delete-range`, `delete-empty`, `delete-comments`
  - **insert**: `insert-before`, `append-after`, `insert-at-line`
  - **inplace**: `replace-in-file`, `replace-in-file-backup`, `delete-lines-in-file`, `multi-replace-in-file`
  - **transform**: `trim-leading-whitespace`, `trim-trailing-whitespace`, `trim-whitespace`, `add-prefix`, `add-suffix`, `number-lines`, `double-space`

## [v1.34.0] - 2026-03-16

### Added

- **kubectl command set** (`configs/kubectl.yaml`) — 43 options across 8 command groups:
  - **pods**: `get`, `describe`, `logs`, `exec`, `port-forward`, `delete`, `top`
  - **deployments**: `get`, `describe`, `scale`, `restart`, `rollout-status`, `rollout-history`, `rollout-undo`, `set-image`
  - **services**: `get`, `describe`, `expose`
  - **config**: `get-contexts`, `current-context`, `use-context`, `set-namespace`, `view`
  - **manifests**: `apply`, `delete-manifest`, `diff`, `kustomize-apply`
  - **nodes**: `get`, `describe`, `top`, `cordon`, `uncordon`, `drain`
  - **secrets**: `get`, `describe`, `create-generic`, `create-from-file`, `create-docker-registry`, `delete`
  - **configmaps**: `get`, `describe`, `create-from-literal`, `create-from-file`, `delete`
  - **resources**: `get`, `describe`, `delete`, `label`, `annotate` (generic for any resource type)
  - Namespace, pod, deployment, service, node, secret, configmap and context fields Tab-pick from live cluster output throughout.

## [v1.33.0] - 2026-03-16

### Added

- **maven command set** (`configs/maven.yaml`) — 30 options across 7 command groups:
  - **lifecycle**: `clean`, `compile`, `test`, `package`, `verify`, `install`, `deploy`, `clean-install` — all with profile, skip-tests, threads, offline and quiet flags
  - **test**: `run-class`, `run-method`, `run-pattern` — targeted Surefire execution
  - **dependency**: `tree`, `analyze`, `copy-dependencies`, `get`, `resolve-sources`
  - **multimodule**: `build-module` (with --also-make), `resume-from` — module fields Tab-pick from nested pom.xml paths
  - **versions**: `check-dependency-updates`, `check-plugin-updates`, `set-version`, `revert`, `commit`
  - **archetype**: `generate-interactive`, `generate` (non-interactive with full coordinate inputs)
  - **release**: `prepare`, `perform`, `rollback`, `clean`
  - **help**: `effective-pom`, `effective-settings`, `describe-plugin`, `active-profiles`
  - Profile fields Tab-pick from `pom.xml` `<id>` entries throughout.

## [v1.32.0] - 2026-03-16

### Added

- **npm command set** (`configs/npm.yaml`) — 25 options across 7 command groups:
  - **packages**: `install`, `install-all`, `ci`, `uninstall`, `update`, `dedupe`
  - **scripts**: `run` (Tab-picks scripts from package.json), `start`, `test`, `build`, `exec` (npx)
  - **inspect**: `list`, `outdated`, `info`, `why`
  - **audit**: `audit`, `audit-fix`
  - **init**: `init`, `init-template`
  - **publish**: `pack`, `publish`, `deprecate`, `dist-tag-add`
  - **cache**: `verify`, `clean`
  - **config**: `list`, `set`, `set-registry`, `set-scope-registry`
  - `uninstall`, `update` and `why` fields Tab-pick from package.json dependencies.

## [v1.31.0] - 2026-03-16

### Added

- **pip command set** (`configs/pip.yaml`) — 22 options across 6 command groups:
  - **packages**: `install`, `install-from-file`, `install-editable`, `uninstall` (Tab-picks from installed list), `upgrade`, `upgrade-all`
  - **inspect**: `list` (outdated/uptodate/json flags), `show` (files flag), `check`, `outdated`
  - **requirements**: `freeze` (to file), `freeze-to-stdout`
  - **index**: `search` (index versions), `download` (offline/platform flags)
  - **cache**: `info`, `list`, `purge`, `remove`
  - **config**: `list`, `set`, `unset`, `set-index`, `set-trusted-host`
  - Uninstall and upgrade fields Tab-pick from `pip list` output.

## [v1.30.0] - 2026-03-16

### Added

- **git command set** (`configs/git.yaml`) — 38 options across 9 command groups covering the full everyday Git workflow:
  - **repository**: `init`, `clone`, `status`, `log`
  - **branch**: `list`, `create`, `switch`, `rename`, `delete`, `set-upstream`
  - **staging**: `add` (with patch mode flag), `restore` (unstage or discard)
  - **commit**: `commit`, `amend`, `fixup`
  - **remote**: `list`, `add`, `remove`, `fetch`, `pull`, `push`
  - **diff**: `working-tree`, `staged`, `between-refs`
  - **merge**: `merge`, `rebase`, `cherry-pick`
  - **stash**: `save`, `list`, `pop`, `apply`, `drop`, `show`
  - **tag**: `list`, `create`, `create-annotated`, `delete`, `push`, `delete-remote`
  - **reset**: `soft`, `mixed`, `hard`, `revert`
  - **config**: `set-identity`, `set-default-branch`, `list`, `set-editor`, `set-alias`
  - Branch, commit, stash, tag and remote fields Tab-pick from live `git` output throughout.

## [v1.29.0] - 2026-03-16

### Added

- **conda command set** (`configs/conda.yaml`) — 19 options across 3 command groups. Environment name fields Tab-pick from live `conda env list` output throughout:
  - **env**: `list`, `create`, `create-from-file`, `clone`, `activate`, `deactivate`, `remove`, `export`, `export-explicit`
  - **packages**: `list`, `install`, `install-from-file`, `update`, `remove-package`, `search`
  - **conda**: `info`, `update-conda`, `clean`, `config-show`, `add-channel`

## [v1.28.0] - 2026-03-16

### Added

- **openssl command set** (`configs/openssl.yaml`) — 22 options across 4 command groups:
  - **inspect**: `inspect-cert`, `inspect-cert-dates`, `inspect-csr`, `inspect-pkcs12`, `inspect-private-key`, `check-remote-cert`, `check-remote-cert-dates`, `verify-cert-chain`, `check-key-cert-match`
  - **generate**: `rsa-key`, `ec-key`, `csr`, `key-and-csr`, `self-signed-cert`, `dh-params`
  - **convert**: `pem-to-p12`, `p12-to-pem`, `der-to-pem`, `pem-to-der`, `extract-public-key`, `remove-key-passphrase`
  - **digest**: `hash-file`, `encrypt-file`, `decrypt-file`

## [v1.27.0] - 2026-03-16

### Added

- **keytool command set** (`configs/keytool.yaml`) — 10 options covering Java keystore and certificate management:
  `list` (verbose flag), `generate-keypair`, `generate-csr`, `import-cert`, `export-cert` (PEM/DER flag),
  `import-keystore` (PKCS12 ↔ JKS conversion), `delete-entry`, `change-alias`, `change-storepass`, `print-cert`.

## [v1.26.0] - 2026-03-16

### Added

- **system-raw command set** (`configs/system-raw.yaml`) — 17 options across 7 commands covering core filesystem and remote-access operations:
  - **mv**: `move`, `move-into-dir`
  - **rm**: `remove-file`, `remove-recursive`
  - **rmdir**: `remove-empty-dir`
  - **cp**: `copy-file`, `copy-recursive`
  - **rsync**: `local-sync`, `push-to-remote`, `pull-from-remote` (Tab-picks SSH hosts from `~/.ssh/config`)
  - **sftp**: `connect`, `download-file`, `upload-file`
  - **ssh**: `connect`, `run-command`, `local-port-forward`, `remote-port-forward`, `copy-id` (Tab-picks hosts and identity keys)
  - **ssh-keygen**: `generate`, `change-passphrase`, `show-fingerprint`, `scan-host-key`, `add-to-known-hosts`

## [v1.25.0] - 2026-03-16

### Added

- **ps-aux command set** (`configs/ps-aux.yaml`) — 12 options covering common `ps aux` workflows:
  `list-all`, `search-by-name`, `top-by-cpu`, `top-by-memory`, `filter-by-user` (Tab-picks users from `/etc/passwd`),
  `inspect-pid` (Tab-picks live processes sorted by CPU), `process-tree`, `watch-cpu`, `watch-memory`,
  `custom-columns`, `count-by-name`, `processes-on-port`.

## [v1.24.0] - 2026-03-16

### Added

- **Custom star names** — when pressing `*` to star a command, a name prompt now appears inline in the form. Type a custom name and press **Enter** to save it, or press **Enter** with an empty field to keep the default `command › option` label. Press **Esc** to cancel without starring.
  - Custom names are shown in place of the default label in the `/s` starred-commands list and on the dedicated Stars screen.
  - Searching with `/s <term>` now filters stars by their custom name, command name, and option name.

## [v1.23.0] - 2026-03-16

### Added

- **Starred commands** — users can now save any command with its current input values for quick re-use.
  - Press **`*`** in any command form to star it. The command, option name, all input values, and flag states are saved to `~/.config/command-builder/stars.json`.
  - Type **`/s`** in the main search and press **Enter** to open the Starred Commands screen, which lists all saved stars.
  - Press **Enter** on a star to re-open its form pre-filled with the saved values, ready to review or run.
  - Press **`d`** on a star to delete it permanently.

## [v1.22.0] - 2026-03-16

### Added

- **Visual colour picker** — pressing `e` or `Enter` on any colour entry in `/settings` now opens an interactive grid of all 256 ANSI colours instead of a plain text box. The cursor (◆) moves with the arrow keys and the selected ANSI number is shown in a preview below the grid. Pressing `Enter` confirms the selection, `Esc` cancels, and `t` drops into the existing free-text input for entering hex (`#rrggbb`) or other values directly.

## [v1.21.0] - 2026-03-16

### Fixed

- **Settings edit input misaligned** — the bordered text input shown when editing a colour or other setting had its left/bottom border edges flush with the left margin while the top border was indented. The cause was prepending `"  "` via string concatenation, which only applied to the first rendered line. Fixed by wrapping the rendered block with `lipgloss.NewStyle().MarginLeft(2)` so all border lines are uniformly indented.

## [v1.20.0] - 2026-03-16

### Added

- **Sub-command input completions** — any config `Input` can now specify a `sub_command` field containing a shell command. When the user focuses that input and presses **Tab**, the command runs asynchronously and its CSV output populates a scrollable picker overlay. Column 0 of each line is the value inserted into the field; optional column 1 is a display-only detail (e.g. image name for Docker containers). Navigate with **↑↓**, confirm with **Enter**, or dismiss with **Esc** or **Tab**. A `Tab: pick value` hint appears in the status area when the focused input has a sub-command defined.
- **Docker container picker** — the `container` input on `docker › exec-container`, `docker › inspect`, and `docker › logs` now use `docker ps --format '{{.Names}},{{.Image}}'` to populate a live list of running containers, showing the image alongside each name for easy identification.
- New documentation: `docs/sub-command-completions.md` — feature guide with YAML syntax, CSV format contract, and key binding reference.

## [v1.19.0] - 2026-03-16

### Changed

- **README rewrite** — replaced the placeholder README with a comprehensive guide covering installation, all screens and keyboard shortcuts, config file format, sharing/import, data locations, and development workflow.



### Added

- **Clipboard copy on command build** — when "run on enter" is disabled, the assembled command is now automatically copied to the clipboard in addition to being printed to stdout. A warning is printed to stderr if the clipboard is unavailable.

## [v1.17.0] - 2026-03-16

### Fixed

- **Lint cleanups** — removed deprecated lipgloss `Copy()` usage and unused UI fields/variables so `golangci-lint` passes cleanly.

## [v1.16.0] - 2026-03-15

### Fixed

- **Tab completion navigation broken by cursor blink** — pressing Down/Up to browse path completions would silently replace the completion list after the first keypress. The live-completion recompute in the input delegate ran on every message (including the textinput's periodic cursor-blink event). After `SetValue` placed a selected path (e.g. `./build/`) into the field, the blink message triggered a recompute using that new value, overwriting the original `./` listing with the subdirectory's contents. Subsequent Down presses appeared to do nothing (list was gone or different), and with all required fields now filled, the next accidental Enter would submit. Fix: the live recompute is now skipped when `showCompletion` is already `true` so the browsable list is preserved for the full navigation session.

## [v1.15.0] - 2026-03-15

### Changed

- **Flag input visual improvements** — enabled/disabled flag states are now immediately obvious in the form screen.
  - Flag rows render with a green rounded border and a `✓` prefix when **on**, and a muted border with blank prefix when **off**.
  - The focused flag row shows `(Space: toggle)` inline so the interaction is self-documenting.

## [v1.14.0] - 2026-03-15

### Added

- **Flag-type inputs** — command options can now declare inputs with `type: flag`.
  - Rendered as a `[ ]` / `[x]` toggle in the form screen instead of a text field.
  - Press **Space** to enable or disable a flag when it is focused.
  - When a flag is toggled on, its `default` value (e.g. `z`, `-z`) is inserted into the template; when off it is omitted entirely.
- **Optional input omission** — non-required inputs that have no value are now removed from the built command rather than showing an `<input_name>` placeholder.
  - Consecutive whitespace left by omitted arguments is collapsed automatically.
- **`tar create` example** in `default.yaml` demonstrating flag inputs for gzip (`-z`), bzip2 (`-j`), and xz (`-J`) compression, producing a clean command regardless of which flags are enabled.

## [v1.13.0] - 2026-03-15

### Added

- **Run on Enter** setting in `/settings` → General section.
  - When enabled, confirming a built command executes it directly in the user's `$SHELL` (via `shell -c "<cmd>"`) instead of printing it to stdout.
  - When disabled (default), the existing print-to-stdout behaviour is preserved.
  - The form screen footer hint updates to show **"run command & quit"** vs **"copy command & quit"** depending on the setting.
  - The toggle is saved to `~/.config/command-builder/settings.json` and persists between sessions.

## [v1.12.0] - 2026-03-15

### Changed

- Application version is now displayed only in the footer; removed from the header across all screens (search, config manager, settings).

## [v1.11.0] - 2026-03-15

### Added

- **Custom application name** — users can now set a custom name for the application in the `/settings` menu.
  - The new **App Name** field appears in a "General" section at the top of the settings screen.
  - The chosen name is reflected immediately in the header across all screens (search, form, config manager, edit, settings).
  - After saving a new name, the user is prompted:  
    _"Add alias to ~/.bashrc?"_ — pressing **y** appends `alias <name>='command-builder'` to `~/.bashrc` so the app is callable by the new name; pressing **n** or **Esc** skips the alias step.
  - The name persists between sessions via `~/.config/command-builder/settings.json`.
  - Default value: `Command Builder`.
- Settings screen now shows a **General** section (app name) above the existing **Colour Palette** section.

## [v1.10.0] - 2026-03-09

### Documentation

- **User guide** (`docs/user-guide.md`)
  Added a comprehensive end-user guide covering all major workflows:
  - Installation and launching (pre-built binary and build-from-source).
  - Search screen: basic search, tag-based search, slash-prefix filters
    (`/default`, `/all`, `/<name>`, `/config`, `/import`, `/settings`).
  - Form screen: filling required and optional fields, Tab path completion
    for `file`/`dir` inputs, live command preview, confirming and capturing
    output via stdout.
  - Config Manager: importing from a URL or local file (with Tab completion),
    creating empty configs, updating from a source URL, exporting, and
    deleting (including restoring the hidden built-in default).
  - Command editor: all three levels (Commands → Options → Inputs), creating
    and editing options with template placeholders and tags, auto-generation
    of inputs from `{{variable}}` placeholders on save, and all four input
    types (`string`, `file`, `dir`, `flag`).
  - Settings screen: changing individual colours (ANSI or hex), resetting one
    or all colours, and the persistent `settings.json` storage location.
  - Tips & tricks section: eval piping, team config sharing via URL, quick
    `/import` from the search bar, config directory backup, and restoring the
    built-in default.
  ASCII-art mock screenshots are included throughout to illustrate each screen.

### Added

- **Auto-detect template variables as inputs**
  When saving a new or edited Option in the command editor, the template string
  is scanned for `{{varName}}` placeholders. Any variable that does not already
  have a corresponding Input entry is automatically added as an optional
  `string` input, saving the user from manually re-entering every placeholder
  as an Input. Variables are deduplicated — repeated occurrences of the same
  placeholder produce a single Input. The generated inputs can be further
  refined (type, description, required flag, default value) by drilling into
  the Inputs level with `Enter` and pressing `e` to edit.

## [v1.9.0] - 2026-03-04

### Added

- **Searchable tags on command options**
  Each `Option` in a config now supports an optional `tags` list — a set of
  alternate terms or aliases a user might think of instead of the command's
  real name. Tags are matched at search time with the same priority as option
  names (exact match scores 80, prefix match 50, substring match 25), so
  relevant options surface even when the query doesn't match the name or
  description directly.

  **YAML format:**
  ```yaml
  options:
    - name: "print-p12"
      description: "Print P12 keystore content"
      template: "openssl pkcs12 -info -in {{input_file}} -passin pass:{{password}}"
      tags: ["pfx", "certificate", "inspect", "keystore"]
      inputs: ...
  ```

- **Tag editing in the command editor**
  The Option edit form (`e` / `n` while at the Options level of the editor
  screen) now includes a **Tags** field. Enter tags as a comma-separated list;
  they are split, trimmed and stored when the form is saved with `Ctrl+S`.
  Existing tags are pre-populated when editing an option.

- **Tag display in the option list**
  When an option has tags, they are shown in square brackets after the
  template column in the Options browse list for quick reference.

- **Example tags in `configs/default.yaml`**
  Seven options across `openssl`, `tar`, `grep`, and `docker` commands have
  been annotated with example tags to demonstrate the feature out of the box.

## [v1.8.0] - 2026-02-27

### Added

- **GitHub Actions CI workflow** (`.github/workflows/ci.yml`)
  Runs automatically on every pull request targeting `main`.
  - **Lint** — `golangci-lint` with a 5-minute timeout.
  - **Test** — `go test -v -race -count=1 ./...`
  - **Build** — cross-compiles for `linux/amd64` and `linux/arm64` to confirm
    both targets build cleanly. Build job depends on lint + test passing.

- **GitHub Actions release workflow** (`.github/workflows/release.yml`)
  Runs automatically on every push to `main` (skips its own version-bump
  commit via the `chore: bump version` message guard).
  1. Runs the test suite.
  2. Reads the current version from `internal/tui/version.go`, increments the
     minor component, and writes the new version back.
  3. Commits the change as `chore: bump version to vX.Y.Z`, creates and pushes
     a matching git tag.
  4. Cross-compiles `command-builder-linux-amd64` and
     `command-builder-linux-arm64` with the new version baked in via
     `-ldflags`.
  5. Extracts the matching section from `CHANGELOG.md` and creates a
     GitHub Release with both binaries attached.

### Notes

- From this version onward, version bumps in `internal/tui/version.go` are
  automated by the release workflow on every merge to `main`. Manual bumps
  are still applied locally for the initial commit of each feature.
- The `concurrency: release` group ensures only one release can run at a time.
- `GITHUB_TOKEN` is used for all git operations; no extra secrets are required
  unless the repo has branch-protection rules that block bot pushes.
- The `softprops/action-gh-release@v2` action is used for creating releases.

## [v1.7.0] - 2026-02-22

### Added

- **Edit the built-in default config**
  Pressing `e` on the built-in default in the Config Manager now works.
  On first edit the embedded default is "promoted" to a real file at
  `~/.config/command-builder/configs/default.yaml` before the editor
  opens. Subsequent launches load that file instead of the embedded
  version, so edits are preserved across restarts.

- **Delete the built-in default config**
  Pressing `d` on the built-in default and confirming now removes it
  from the current session. A tombstone marker file
  (`~/.config/command-builder/configs/.default-hidden`) is written so
  the embedded default is not reloaded on the next launch. The tombstone
  can be deleted manually to restore the embedded default.

- **`Manager.PromoteDefaultConfig`**
  New method in `internal/config` that saves an embedded (file-path-less)
  config to the user's config directory and updates its `FilePath` in place.

### Changed

- `Manager.DeleteConfig` now writes the tombstone file when deleting the
  embedded built-in config instead of silently succeeding without any
  persistent effect.
- `Manager.NewManager` startup order changed: user configs (including a
  user-saved `default.yaml`) are loaded first, and the embedded default is
  only loaded when no user-backed version exists and the tombstone is absent.
  Previously user configs named `"default"` were silently skipped.

## [v1.6.0] - 2026-02-17

### Added

- **`/settings` menu for global application settings**
  Typing `/settings` in the search bar (or pressing Enter while the query
  starts with `/settings`) opens a new full-screen Settings panel. It is
  integrated into the same screen-routing system as the Config Manager and
  Command Editor.

- **Customisable colour palette**
  The Settings screen exposes all nine theme colours (Primary, Accent,
  Success, Warning, Error, Muted, Text, Selected BG) as editable entries.
  Each row shows a coloured swatch, the current value, and a short
  description of where that colour is used.

  - `↑`/`↓` — navigate entries
  - `e` or `Enter` — enter edit mode for the selected colour; accepts ANSI
    terminal codes (`0`–`255`) or CSS hex values (`#rrggbb`)
  - `r` — reset the selected colour to its built-in default
  - `R` — reset the entire palette to built-in defaults
  - `Esc` — return to the search screen

- **Persistent colour settings**
  Colour choices are saved to `~/.config/command-builder/settings.json` and
  loaded automatically at startup. The palette is applied immediately on every
  colour change — no restart required.

- **`ApplyTheme` API in `internal/tui`**
  `styles.go` is refactored so that every Lipgloss style variable is
  reassigned by `ApplyTheme(config.AppSettings)`. Previously styles were
  plain `var` initialisers run once at package init; they are now rebuilt
  whenever the user changes a colour, ensuring every screen reflects the
  updated palette on its very next render.

- **`config.AppSettings` type**
  `internal/config/settings.go` introduces `AppSettings`, `DefaultSettings()`,
  `LoadSettings()`, and `SaveSettings()` — a self-contained layer for
  persisting non-config user preferences independently of the YAML config
  files.

## [v1.5.0] - 2026-02-12

### Added

- **Application version displayed in every screen's footer**
  All four screens (search, form, config manager, editor) now show `AppVersion`
  right-aligned in the status bar, rendered with the muted title-version style.
  The version badge is produced by a new shared `footerVersion()` helper.

### Fixed

- **Footer wrapping on all screens**
  `StyleStatus` has `Padding(0,1)`, meaning Lipgloss's `Width(w)` sets the
  *content* area to `w` characters and then adds 1 space of padding on each
  side — yielding a total of `w+2` columns, which wrapped on any terminal at
  exactly the content width. Fixed by a new `renderFooter(w, left, right)`
  helper in `internal/tui/footer.go` that targets `Width(w-2)` and uses
  `lipgloss.Width()` (ANSI-aware) for the gap calculation, ensuring the status
  bar fits exactly `w` columns on every screen.

## [v1.4.0] - 2026-02-07

### Added

- **Tab autocomplete for `/import` in the search screen**
  While typing `/import <path>` in the main search bar, pressing `Tab` now triggers
  the same glob-based path completion used by the Config Manager's file import prompt:
  - Single match → completed immediately.
  - Multiple matches → first `Tab` fills the longest common prefix; subsequent presses
    cycle through all matches, with the current selection highlighted.
  - Up to 8 completions are shown inline in the results area while in `/import` mode.
  - The hint bar updates to `Enter to import · Tab to autocomplete path` while
    a `/import` prefix is active, and resets when the user clears the query.
  - Completions reset automatically when the path is edited manually.

- **Import configs from local files**
  Added `Manager.ImportConfigFromFile(path)` which reads a local YAML file, handles
  `~` expansion, resolves to an absolute path, and copies the config into the managed
  config directory — mirroring the behaviour of the existing URL import.
  A leading `~` in the path is expanded to the user's home directory.

- **`f` key in Config Manager to import from a local file**
  Pressing `f` on the Config Manager screen opens a new "Import from local file" prompt.
  The status bar now shows both `i import URL` and `f import file` hints.

- **Tab autocomplete for the file import path**
  While typing a file path in the import-file prompt, pressing `Tab` triggers
  glob-based path completion:
  - Single match → path is completed immediately.
  - Multiple matches → first `Tab` fills the longest common prefix (bash-style);
    subsequent `Tab` presses cycle through all matches.
  - Up to five completions are shown as an inline list below the input; the currently
    selected entry is highlighted.
  - Completions reset automatically when the path is edited manually.

- **`/import` slash command now accepts local file paths**
  On the search screen, `/import <value>` previously only accepted URLs. It now
  auto-detects the argument: values starting with `http://` or `https://` are fetched
  as URLs; everything else is treated as a local file path and imported via
  `ImportConfigFromFile`. The hint text is updated to reflect this.

## [v1.3.0] - 2026-02-02

### Maintenance

- Project versioning established: minor version incremented for every change.

## [v1.2.0] - 2026-01-28

### Added

- **Added `.github/copilot-instructions.md`**
  Created Copilot instructions file documenting project overview, tech stack, architecture,
  code conventions, common workflows, and doc references to give GitHub Copilot consistent
  context across all sessions.

## [v1.1.0]

### Added

- **Edit existing configs from the Config Manager screen**
  Pressing `e` on any non-built-in config opens a new Edit screen that lets the
  user browse and mutate the full config hierarchy: Commands → Options → Inputs.
  Each level shows a navigable list (`↑`/`↓`, `Enter` to drill in, `Esc` to go
  back up). Within a level the user can create new items (`n`), edit the selected
  item (`e`), or delete it (`d`, with a `y`/`n` confirmation prompt). Forms are
  submitted with `Ctrl+S`. All changes are written to disk immediately via a new
  `Manager.UpdateConfig` method. Built-in (embedded) configs show an error when
  `e` is pressed rather than opening the editor.

- **Per-config source URL and one-key pull updates**
  `Config` now has an optional `source_url` YAML field. When a config is
  imported with `i` (import from URL) the URL is stored in the saved YAML file.
  Pressing `u` on any config that has a source URL prompts for confirmation
  (`yes`) and then re-fetches the YAML from that URL via a new
  `Manager.PullConfig` method, replacing the config's commands in place while
  preserving the local name and file path. Configs with a stored URL are marked
  with a `[url]` badge in the list. Configs without a source URL show a
  descriptive error instead of opening the prompt.

### Maintenance

- **Centralized project versioning**
  Removed the hardcoded `const appVersion` from `search_screen.go`. Version is now defined
  in a single file (`internal/tui/version.go`) as exported `var AppVersion = "v1.1.0"` so
  it can be overridden at build time via `-ldflags`. Updated `build.sh` to automatically
  inject the version using `git describe --tags --always --dirty`, falling back to the
  hardcoded value when no git tag exists. All four TUI screens reference `AppVersion`.

### Fixed

- **Search input not responding to keystrokes**
  `Init()` used a value receiver, so `textinput.Focus()` mutated a temporary
  copy of the model that was immediately discarded. The stored model's input was
  never focused, causing all keystrokes to be silently dropped. Fixed by calling
  `ti.Focus()` directly in `NewSearchModel()` before the value is stored, and
  changing `Init()` to return `textinput.Blink` instead.

- **Search box disappearing when query returns no results**
  The view's `reserved` line count was `6`, but the actual chrome (title,
  rounded-border top, input, rounded-border bottom, hint, separator, status bar)
  occupies `7` lines. The resulting 1-line overflow caused the terminal renderer
  to clip the top of the view. On top of this, the "No results" message used
  `Padding(1,2)` (3 lines) and was written outside the bounded results area,
  adding further overflow when results were empty. Fixed by correcting `reserved`
  to `7` and rewriting the results/no-results block so the area always occupies
  exactly `visRows` lines.

- **Slash commands (`/config`, `/import`) breaking live search**
  The guard that prevented search from re-running during slash-command input was
  widened too broadly (blocking all `/`-prefixed queries), which broke the live-
  filter slash commands (`/default`, `/all`, `/<name>`). Reverted to only
  skipping `/config` and `/import`, which are the only commands handled solely
  on Enter.

- **Scroll position resetting every ~500ms**
  `textinput.Blink` fires a periodic tick message to toggle the cursor. This
  tick fell through to the text-input delegate block, which re-ran the search
  and unconditionally reset `selectedIdx` and `scrollTop` to `0` — snapping the
  list back to the top twice a second. Fixed by snapshotting the query value
  before delegating to the text input and only resetting scroll/selection when
  the query actually changes.

- **Search screen reverting to a fixed small width after navigating back**
  `NewSearchModel` hardcoded `ti.Width = 60` and left the model's `width` and
  `height` fields as zero. On first launch the initial `WindowSizeMsg` corrected
  the dimensions, but no `WindowSizeMsg` fires when navigating back from the
  form screen, so a fresh model created on back-navigation stayed at zero width
  for the rest of the session. Fixed by adding `w, h int` parameters to
  `NewSearchModel` and passing `a.width, a.height` from `AppModel` in the
  `backToSearchMsg` handler.

## [v1.0.0] - 2026-01-15

### Added

- Initial release: interactive TUI for building and searching CLI commands from YAML config files.
- Search screen with full-text fuzzy matching.
- Form screen for filling in command option inputs.
- Config manager for loading and managing configuration packs.
