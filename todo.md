# Todo

- [x] create copilot-instructions.md
- [x] project versioning, increment minor version for all changes unless otherwise specified
- [x] Update configs so that they can be loaded from local files, not just urls. include tab autocomplete for config file selection
- [x] Display application version in the footer
- [x] add /settings menu for global settings within the application. for starters, allow the users to set custom colors for the application
- [x] allow the user to remove/edit the default config
- [x] github actions ci cd pipelines, run tests, build, and lint on pull request, when merged into main increment the application version and release the application. build binaries for arm64 and linux x64 for the release and create a tag of the release version as well
- [x] add searchable tags to the commands so that users can add terms they might think of instead of the real command and they still get search results
- [x] update command editor so that when a user adds a template, the variables are detected and added as optional string inputs by default
- [x] Create user guide with screenshots for how to add commands, navigate and use the tool, etc
- [x] add option to change application name in settings menu, default is command-builder. I'd like to prompt the user and request to add an alias to their bashrc to rename command-builder to their provided name. if they say yes, add the alias to the bashrc and show the new name in the application header, if they say no only reflect the change in the application header
- [x] Version is now displayed in header and footer, let's keep it in the footer and remove it from the header
- [x] add option in settings to run command on enter
- [x] id like to add optional flags to the application, for example with tar i'd like to have the option to gzip and when enabled. if an argument is not required it should not be rendered unless provided
- [x] when user hits enter on build command, copy it to the clipboard as well as printing it to console assuming run on enter is disabled
- [x] update project readme
- [ ] expand configs with additional tools, create separate libraries that can be imported


# Bugs

- [x] fix footer styling, this is the current:
    ```
    find › by-type  Find by type (f=file, d=dir, l=symlink) [default]
    53 result(s)                                                     Ctrl+C 
    quit                                                                                                     
    ```