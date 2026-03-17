# Desired Commandsets
List of command sets to build, while this application started with the goal of a focused, custom group of commands to be created by users for specific workflows, i'd like to make a number of more general command sets to allow for future expansion into chaining commands together (operations), so the following command sets should be built out fully and with careful consideration for final user interaction ie. auto-filling values where possible, providing optional flags on a more general form rather than building custom commands for each possible operation. command sets should be stored in {project_root}/configs, when a new config is built please update this file to reflect the commands in the command set so users can further update the list with desired functionality and it can be checked off when completed

# PS AUX command set
- [x] Common functions with ps aux
  - `list-all` — List all running processes (optional thread flag)
  - `search-by-name` — Find processes matching a name via grep
  - `top-by-cpu` — Show top N processes sorted by CPU usage
  - `top-by-memory` — Show top N processes sorted by memory usage
  - `filter-by-user` — Show all processes owned by a specific user (Tab picks from /etc/passwd)
  - `inspect-pid` — Inspect a specific PID with detailed columns (Tab picks from live process list sorted by CPU)
  - `process-tree` — Display processes as a parent-child tree (optional user filter)
  - `watch-cpu` — Continuously watch top CPU processes at a configurable interval
  - `watch-memory` — Continuously watch top memory processes at a configurable interval
  - `custom-columns` — Show selected ps columns with an optional sort flag
  - `count-by-name` — Count running instances of a named process
  - `processes-on-port` — Find the process listening on a specific port via ss

# System raw command set
- [x] mv commands
  - `move` — Move or rename a file/directory (interactive, no-clobber, backup, verbose flags)
  - `move-into-dir` — Move one or more files into a target directory
- [x] rm commands
  - `remove-file` — Remove one or more files (force, interactive, verbose flags)
  - `remove-recursive` — Remove a directory and all contents (force, interactive, verbose flags)
- [x] rmdir commands
  - `remove-empty-dir` — Remove empty directories, optionally with --parents
- [x] cp commands
  - `copy-file` — Copy a file (preserve, interactive, no-clobber, verbose flags)
  - `copy-recursive` — Copy a directory recursively (preserve, interactive, verbose flags)
- [x] rsync commands
  - `local-sync` — Sync a local directory to another local path (dry-run, delete, checksum flags)
  - `push-to-remote` — Sync local directory to a remote server over SSH (Tab-picks hosts from ~/.ssh/config)
  - `pull-from-remote` — Sync a remote directory to local over SSH (Tab-picks hosts from ~/.ssh/config)
- [x] sftp commands
  - `connect` — Open an interactive SFTP session (Tab-picks hosts from ~/.ssh/config)
  - `download-file` — Download a file from remote via SFTP batch mode
  - `upload-file` — Upload a file to remote via SFTP batch mode
- [x] ssh commands
  - `connect` — Open an interactive SSH session (X11 flag, Tab-picks hosts and identity keys)
  - `run-command` — Execute a command on a remote host
  - `local-port-forward` — Forward a local port to a service via SSH (-L tunnel)
  - `remote-port-forward` — Expose a local port on a remote server (-R reverse tunnel)
  - `copy-id` — Install public key on remote host for passwordless login
- [x] ssh-keygen commands
  - `generate` — Generate a new SSH key pair (type, bits, comment, output path)
  - `change-passphrase` — Change/remove passphrase of an existing key (Tab-picks keys from ~/.ssh)
  - `show-fingerprint` — Display fingerprint of a public key (Tab-picks ~/.ssh/*.pub)
  - `scan-host-key` — Retrieve remote host's public key fingerprint via ssh-keyscan
  - `add-to-known-hosts` — Scan host and append key to ~/.ssh/known_hosts

# Keytool command set
- [x] Keytool commands and usages
  - `list` — List all keystore entries (optional verbose flag)
  - `generate-keypair` — Generate a key pair and self-signed certificate (RSA/EC/DSA, custom dname, validity)
  - `generate-csr` — Generate a Certificate Signing Request for an existing key pair
  - `import-cert` — Import a certificate or CA reply into a keystore (noprompt flag)
  - `export-cert` — Export a certificate to a file (PEM/DER via -rfc flag)
  - `import-keystore` — Convert/migrate between keystore formats (e.g. PKCS12 ↔ JKS)
  - `delete-entry` — Delete an entry from a keystore by alias
  - `change-alias` — Rename an alias within a keystore
  - `change-storepass` — Change a keystore's password
  - `print-cert` — Display certificate file contents without a keystore

# Openssl Command set
- [x] Openssl commands and usages
  - **inspect**: `inspect-cert`, `inspect-cert-dates`, `inspect-csr`, `inspect-pkcs12`, `inspect-private-key`, `check-remote-cert`, `check-remote-cert-dates`, `verify-cert-chain`, `check-key-cert-match`
  - **generate**: `rsa-key`, `ec-key`, `csr`, `key-and-csr`, `self-signed-cert`, `dh-params`
  - **convert**: `pem-to-p12`, `p12-to-pem`, `der-to-pem`, `pem-to-der`, `extract-public-key`, `remove-key-passphrase`
  - **digest**: `hash-file`, `encrypt-file`, `decrypt-file`

# Conda command set
- [x] Conda commands and usages
  - **env**: `list`, `create`, `create-from-file`, `clone`, `activate`, `deactivate`, `remove`, `export`, `export-explicit` (Tab-picks env names from `conda env list` throughout)
  - **packages**: `list`, `install`, `install-from-file`, `update`, `remove-package`, `search`
  - **conda**: `info`, `update-conda`, `clean`, `config-show`, `add-channel`

# Git command set
- [x] git commands and usages
  - **repository**: `init`, `clone`, `status`, `log`
  - **branch**: `list`, `create`, `switch`, `rename`, `delete`, `set-upstream`
  - **staging**: `add` (patch flag), `restore` (unstage or discard)
  - **commit**: `commit`, `amend`, `fixup`
  - **remote**: `list`, `add`, `remove`, `fetch`, `pull`, `push`
  - **diff**: `working-tree`, `staged`, `between-refs`
  - **merge**: `merge`, `rebase`, `cherry-pick`
  - **stash**: `save`, `list`, `pop`, `apply`, `drop`, `show`
  - **tag**: `list`, `create`, `create-annotated`, `delete`, `push`, `delete-remote`
  - **reset**: `soft`, `mixed`, `hard`, `revert`
  - **config**: `set-identity`, `set-default-branch`, `list`, `set-editor`, `set-alias`

# Pip command set
- [x] Pip commands and usages
  - **packages**: `install`, `install-from-file`, `install-editable`, `uninstall`, `upgrade`, `upgrade-all`
  - **inspect**: `list`, `show`, `check`, `outdated`
  - **requirements**: `freeze`, `freeze-to-stdout`
  - **index**: `search`, `download`
  - **cache**: `info`, `list`, `purge`, `remove`
  - **config**: `list`, `set`, `unset`, `set-index`, `set-trusted-host`

# npm command set
- [x] npm commands and usages
  - **packages**: `install`, `install-all`, `ci`, `uninstall`, `update`, `dedupe`
  - **scripts**: `run` (Tab-picks scripts from package.json), `start`, `test`, `build`, `exec` (npx)
  - **inspect**: `list`, `outdated`, `info`, `why`
  - **audit**: `audit`, `audit-fix`
  - **init**: `init`, `init-template`
  - **publish**: `pack`, `publish`, `deprecate`, `dist-tag-add`
  - **cache**: `verify`, `clean`
  - **config**: `list`, `set`, `set-registry`, `set-scope-registry`

# maven command set
- [x] maven commands and usages
  - **lifecycle**: `clean`, `compile`, `test`, `package`, `verify`, `install`, `deploy`, `clean-install`
  - **test**: `run-class`, `run-method`, `run-pattern`
  - **dependency**: `tree`, `analyze`, `copy-dependencies`, `get`, `resolve-sources`
  - **multimodule**: `build-module`, `resume-from`
  - **versions**: `check-dependency-updates`, `check-plugin-updates`, `set-version`, `revert`, `commit`
  - **archetype**: `generate-interactive`, `generate`
  - **release**: `prepare`, `perform`, `rollback`, `clean`
  - **help**: `effective-pom`, `effective-settings`, `describe-plugin`, `active-profiles`

# kubectl command set
- [x] kubectl commands and usages
  - **pods**: `get`, `describe`, `logs`, `exec`, `port-forward`, `delete`, `top`
  - **deployments**: `get`, `describe`, `scale`, `restart`, `rollout-status`, `rollout-history`, `rollout-undo`, `set-image`
  - **services**: `get`, `describe`, `expose`
  - **config**: `get-contexts`, `current-context`, `use-context`, `set-namespace`, `view`
  - **manifests**: `apply`, `delete-manifest`, `diff`, `kustomize-apply`
  - **nodes**: `get`, `describe`, `top`, `cordon`, `uncordon`, `drain`
  - **secrets**: `get`, `describe`, `create-generic`, `create-from-file`, `create-docker-registry`, `delete`
  - **configmaps**: `get`, `describe`, `create-from-literal`, `create-from-file`, `delete`
  - **resources**: `get`, `describe`, `delete`, `label`, `annotate`

# Sed command set
- [x] sed commands and usages
  - **substitute**: `replace-first`, `replace-all`, `replace-nth`, `replace-in-range`, `replace-in-matching-lines`, `replace-extended`
  - **lines**: `print-matching`, `print-range`, `delete-matching`, `delete-not-matching`, `delete-range`, `delete-empty`, `delete-comments`
  - **insert**: `insert-before`, `append-after`, `insert-at-line`
  - **inplace**: `replace-in-file`, `replace-in-file-backup`, `delete-lines-in-file`, `multi-replace-in-file`
  - **transform**: `trim-leading-whitespace`, `trim-trailing-whitespace`, `trim-whitespace`, `add-prefix`, `add-suffix`, `number-lines`, `double-space`

# Go command set
- [x] command set for golang
  - **run**: `run`
  - **build**: `build`, `build-all`, `install`, `clean`
  - **test**: `test`, `test-run`, `test-coverage`, `test-coverage-html`, `benchmark`, `test-count`
  - **modules**: `init`, `tidy`, `get`, `get-latest`, `download`, `vendor`, `graph`, `why`, `list-modules`
  - **quality**: `fmt`, `vet`, `generate`
  - **doc**: `doc`, `godoc-server`
  - **profiling**: `test-cpu-profile`, `test-mem-profile`, `pprof`

# System abstraction Command set
- [x] restarting services, auto populate service list
- [x] kill processes on port
- [x] kill process by id with auto fill of high ram or cpu processes
- [x] list usb devices
- [x] list network devices
- [x] get ip address, default gateway, other network information
  - **services**: `list`, `status`, `start`, `stop`, `restart`, `reload`, `enable`, `disable`, `logs`
  - **processes**: `kill-on-port`, `show-on-port`, `kill-by-pid`, `kill-high-cpu`, `kill-high-memory`, `kill-by-name`
  - **usb**: `list`, `list-verbose`, `tree`, `watch`
  - **network-devices**: `list-interfaces`, `interface-info`, `interface-stats`, `list-connections`, `wifi-list`
  - **network-info**: `ip-addresses`, `public-ip`, `routing-table`, `default-gateway`, `dns-config`, `arp-table`, `open-ports`

# Docker command set
- [x] Docker container and image management
  - **containers**: `list`, `list-all`, `run`, `run-interactive`, `exec`, `logs`, `stop`, `start`, `restart`, `remove`, `inspect`, `stats`, `top`, `copy`
  - **images**: `list`, `pull`, `build`, `tag`, `push`, `remove`, `inspect`, `prune`
  - **volumes**: `list`, `create`, `inspect`, `remove`, `prune`
  - **networks**: `list`, `create`, `inspect`, `connect`, `disconnect`, `remove`
  - **system**: `info`, `prune-all`, `df`

# Docker Compose command set
- [x] Docker Compose multi-container orchestration
  - **lifecycle**: `up`, `up-detached`, `down`, `start`, `stop`, `restart`
  - **inspect**: `ps`, `logs`, `top`, `config`
  - **build**: `build`, `pull`
  - **exec**: `exec`, `run`

# curl command set
- [x] HTTP requests and file transfers
  - **request**: `get`, `post-json`, `post-form`, `put`, `patch`, `delete`, `head`
  - **auth**: `bearer-token`, `basic-auth`
  - **download**: `download-file`, `download-resume`, `follow-redirects`
  - **inspect**: `verbose`, `response-headers-only`, `time-request`

# jq command set
- [x] JSON processing and querying
  - **query**: `pretty-print`, `select-field`, `select-nested`, `filter-array`, `select-where`, `keys`, `values`
  - **transform**: `map`, `to-csv`, `to-entries`, `from-entries`, `del-field`, `add-field`
  - **reduce**: `length`, `first`, `last`, `unique`, `group-by`, `sort-by`
  - **file**: `from-file`, `slurp-multiple`

# awk command set
- [x] Pattern scanning and text processing
  - **print**: `print-field`, `print-range`, `print-matching`, `print-not-matching`, `print-last-field`, `print-line-count`
  - **transform**: `add-prefix`, `reorder-fields`, `custom-delimiter`, `sum-column`, `average-column`
  - **filter**: `match-pattern`, `numeric-compare`, `between-lines`

# find command set
- [x] File search and bulk operations
  - **search**: `by-name`, `by-extension`, `by-type`, `by-size`, `by-modified-time`, `by-owner`, `by-permissions`
  - **action**: `delete-matching`, `exec-on-matching`, `copy-matching`, `print-with-size`

# tar command set
- [x] Archive creation and extraction
  - **create**: `create-gzip`, `create-bzip2`, `create-xz`, `create-uncompressed`
  - **extract**: `extract`, `extract-to-dir`, `extract-single-file`
  - **inspect**: `list-contents`, `test-archive`

# grep command set
- [x] Pattern search in files and streams
  - **search**: `search-file`, `search-recursive`, `search-case-insensitive`, `search-word`, `invert-match`, `count-matches`
  - **context**: `before-context`, `after-context`, `surrounding-context`
  - **output**: `filenames-only`, `line-numbers`, `with-color`
  - **extended**: `extended-regex`, `fixed-string`, `perl-regex`

# disk command set
- [x] Disk usage and filesystem inspection
  - **usage**: `disk-free`, `disk-free-human`, `inode-usage`, `dir-size`, `dir-size-sorted`, `find-large-files`
  - **filesystem**: `list-mounts`, `mount`, `unmount`, `remount-readonly`
  - **io**: `io-stats` (iostat), `io-per-process` (iotop)

# permissions command set
- [x] File and directory permission management
  - **chmod**: `set-permissions`, `add-execute`, `remove-write`, `recursive-chmod`, `set-from-reference`
  - **chown**: `change-owner`, `change-group`, `change-owner-and-group`, `recursive-chown`
  - **acl**: `get-acl`, `set-acl`, `remove-acl`

# terraform command set
- [x] Infrastructure as Code with Terraform
  - **workspace**: `list`, `new`, `select`, `delete`
  - **core**: `init`, `validate`, `plan`, `apply`, `destroy`, `refresh`
  - **state**: `list`, `show`, `move`, `remove`, `pull`, `push`
  - **inspect**: `output`, `show`, `graph`, `providers`

# helm command set
- [x] Kubernetes package manager
  - **repo**: `list`, `add`, `update`, `remove`
  - **chart**: `search`, `inspect`, `pull`
  - **release**: `list`, `install`, `upgrade`, `rollback`, `uninstall`, `status`, `history`
  - **inspect**: `get-values`, `get-manifest`, `diff`

# gpg command set
- [x] GNU Privacy Guard encryption and signing
  - **keys**: `list-public`, `list-secret`, `import`, `export`, `delete-public`, `delete-secret`, `generate`, `search-keyserver`
  - **encrypt**: `encrypt-file`, `encrypt-symmetric`, `decrypt-file`
  - **sign**: `sign-file`, `clearsign`, `verify`, `sign-and-encrypt`
  - **keyring**: `fingerprint`, `edit-key`, `export-ownertrust`, `import-ownertrust`

# crontab command set
- [x] Scheduled task management
  - **manage**: `list`, `edit`, `remove-all`, `install-from-file`
  - **system**: `list-system-crons`, `validate-expression`

# ffmpeg command set
- [x] Audio and video processing
  - **convert**: `video-to-mp4`, `video-to-gif`, `audio-to-mp3`, `change-container`
  - **inspect**: `probe`, `probe-json`
  - **edit**: `trim`, `extract-audio`, `extract-frames`, `concat`, `scale`, `compress-video`
  - **stream**: `record-screen`, `stream-to-rtmp`