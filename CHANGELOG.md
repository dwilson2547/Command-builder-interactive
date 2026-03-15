# Changelog

All notable changes to Command Builder are documented here.

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

### Added

- User guide with screenshots for navigation, adding commands, and using the tool.
- Variable auto-detection: when a template is added in the command editor, variables are automatically parsed and added as optional string inputs.

## [v1.9.0] - 2026-03-04

### Added

- Searchable tags for commands — users can add alias terms that surface results even when the real command name differs from how they think of it.

## [v1.8.0] - 2026-02-27

### Added

- GitHub Actions CI/CD pipelines: run tests, build, and lint on pull requests; publish tagged releases with arm64 and linux-x64 binaries when merged into main.

## [v1.7.0] - 2026-02-22

### Added

- Option to remove or edit the default config from within the application.

## [v1.6.0] - 2026-02-17

### Added

- `/settings` menu with customisable colour palette (primary, accent, success, warning, error, muted, text, selected-background).
- Colours persist between sessions.

## [v1.5.0] - 2026-02-12

### Added

- Application version displayed in the footer.

## [v1.4.0] - 2026-02-07

### Added

- Tab autocomplete for config file paths when loading local configs.
- Configs can now be loaded from local files, not just URLs.

## [v1.3.0] - 2026-02-02

### Added

- Project versioning: minor version incremented for every change.

## [v1.2.0] - 2026-01-28

### Added

- Copilot instructions file (`copilot-instructions.md`).

## [v1.0.0] - 2026-01-15

### Added

- Initial release: interactive TUI for building and searching CLI commands from YAML config files.
- Search screen with full-text fuzzy matching.
- Form screen for filling in command option inputs.
- Config manager for loading and managing configuration packs.
