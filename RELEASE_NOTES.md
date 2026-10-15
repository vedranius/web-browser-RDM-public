## Web Remote Manager PRO v12.4.5 — One-click Medium and usable log search

### Added

- **Medium:** add the original trusted tool/config ZIPs and optional SSH ZIP, then **Prepare and run**. WRM starts Docker, prepares the original dependencies, validates configuration, synchronizes profiles and opens the unchanged tool. Approved imports survive setup failure and can be resumed without re-uploading.
- **Logs:** automatic discovery from manual roots, service log folders and Supervisor metadata; progressive server results, visible errors, operator/service filters and historical cross-file keyword/date search. Tiles show modification/event times, highlighted excerpts and open/follow actions. Search can be canceled or continued when a bound is reached.
- **Environments:** combined text/service/environment/server filters with sticky controls.

### Fixed

The whole Overview/Installations filter/check/action area stays at the top during scroll. Monitoring tab responses are isolated. Date-only searches and continuation-event timestamps are handled. Discovery avoids repeated per-file metadata round trips; host metric versus log prerequisites are explained.

### Upgrade

Retain the database/key and private runtime directories/volumes, then replace the executable/container. No schema migration. Windows setup may display the official Docker installer and vendor/OS prompts; required license choices, WSL/virtualization setup and restarts remain user actions. A Windows service needs a preinstalled Docker account; other Windows containers are never stopped or switched automatically. Linux without installed Docker/native prerequisites needs its administrator's initial package setup. Windows builds are cross-compiled, without a native Windows VM run.

Search is bounded per file/round and excerpts; incomplete results show Continue. Compressed/binary logs remain explicit unsupported text inputs. See [guide](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.5/docs/V12.4.5_GUIDE.md), [validation](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.5/docs/RELEASE_SUMMARY_v12.4.5.md) and [downloads](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.5).

### Included from v12.4.4 — Monitoring, logs and AI workflows

#### Added

- **Monitoring** brings Zabbix server metrics, history, agent detection and service Functionality into one workspace. Use passive agent/agent 2 over SSH or your configured network route, verified certificate TLS, PSK through the existing remote `zabbix_get`, or an existing HTTPS Zabbix server API for active-only agents. Unavailable metrics are unknown and diagnostics remain visible.
- **Logs** supports multiple roots/operators/services, historical segmented search and date/timezone/keyword filters, highlighted results, copying/downloads, terminal actions and up to 16 movable/resizable live windows with tiling/cascade/pause/rotation handling.
- **Client GitLab tokens and cache** let a computer with GitLab access supply bounded repository snapshots to WRM. The GitLab token stays in browser memory or the companion's local token file. The latest complete successful server cache survives client disconnect and WRM restart; refresh is atomic, sharing is explicit and reviewed deployment content stays pinned.
- **MCP onboarding and permissions:** CLI setup scripts, OAuth/cloud guidance, Cursor install link, `wrm` prompt and authorized server picker, native assistant `/wrm`, safe structured configuration inventory, visible Automatic within limits and editable existing grants. Grant edits retain bearer/expiry and terminate old sessions/approvals.
- **Terminal and navigation:** previewed/redacted selection-to-AI context, Open terminal here for POSIX/PowerShell, Copy/Paste/native-menu fallback, hostname groups/combined Functionality filters and folder/connection colors/icons.

#### Fixed

General import recognizes mRemoteNG XML/PuTTY `.reg` ([issue #5](https://github.com/vedranius/web-browser-RDM-public/issues/5)); synthetic 4 MiB XML and UTF-16 registry exports were exercised in the browser. Provider header/request/idle timeouts are configurable for longer local AI requests. Wide replay fits phone/desktop widths. Search/filter bars remain at the top of supported scroll areas, and MCP wizard controls align correctly on narrow screens.

#### Upgrade and operation

Replace the executable/container and retain the database, encryption key and private runtime files. Migrations are additive. Functionality is now **Monitoring → Functionality**. Configure each server's installed Zabbix agent and log roots in Monitoring settings; automatic polling is opt-in. MCP stays disabled by default and existing explicit administrator policies remain in force.

Browser GitLab access needs the GitLab CORS/origin/local-network permissions; use the WRM companion when the browser cannot reach that API. The relay is for WRM repository API reads: native Medium Git clone/build/config Git operations and GitLab writes still need a direct route. Native Medium's Windows Linux-container prerequisite from v12.4.3b remains.

[Download v12.4.4](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.4) · [setup/limits](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.4/docs/V12.4.4_GUIDE.md) · [validation](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.4/docs/RELEASE_SUMMARY_v12.4.4.md). Windows packages are cross-compiled; no native Windows VM execution is claimed.

### Included from v12.4.3b — Windows Medium runtime fix

#### Fixed

Windows `.exe` installations run the original Medium deploy tool inside a managed Linux Docker runtime. Python, Git, OpenSSH and rsync are supplied by its pinned image; separate Windows installations are unnecessary. ZIP imports survive setup failure and can be prepared again without re-uploading. Private Linux volumes preserve deploy history, rollback state, permissions and config edits across restarts. Native Linux mode remains available.

The original tool remains unchanged. Runtime source/venv are read-only, SSH/config are private and builds use an unprivileged user. Runtime environment forwarding preserves Docker CLI context and keeps secret values out of command arguments. Only owned containers are stopped.

#### Windows upgrade

**Download the new Windows executable manually:** v12.4.3 cannot parse a `b` suffix; v12.4.3b supports it for later updates. Retain the DB/key. Start Docker Desktop with **Linux containers** on Windows 10/11, then use **Git → Medium → Prepare runtime / retry** for an existing tool import, or import your trusted tool ZIP first. Import deploy-config/SSH profile and start Medium. A browser client needs no deployment dependencies.

Docker Desktop is not supported on Windows Server; use a supported Linux WRM host/VM/container for Medium. The Docker engine must be local and accessible to the WRM account. Back up private `.medium` files and the managed Linux runtime/user Docker volumes.

#### Validation and packages

Real pinned-image setup, original comparison/build, two actual OpenSSH/rsync deployments, rollback and stop/start history persistence were exercised against a local fixture. The build user could not read private SSH/config. Browser layouts at 390/1100/1440/1920 px had no JavaScript errors. Windows is cross-compiled; no native Windows VM run is claimed. Full validation: [hotfix summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3b/docs/RELEASE_SUMMARY_v12.4.3b.md).

[Download v12.4.3b](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.3b): 16 platform packages plus checksums; Docker `ghcr.io/vedranius/wrm-pro:v12.4.3b` and `latest`. The remaining accepted work is planned for **v12.4.4** in ROADMAP, not included in this hotfix.

### Included from v12.4.3 — Original Medium deployment, config profiles and Bitwarden

#### Added

- **Medium** uses the original supplied deploy tool and deploy-config format unchanged: original screens, comparisons, file selection/exclusions, builds, native plans/confirmations, deployment progress, backups, history, rollback and config Git controls. Import the private tool/config ZIPs in Git → Medium. WRM handles authenticated access, private directories, encrypted runtime preferences and setup status. [Original-tool parity record](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/docs/MEDIUM_PARITY_v12.4.3.md).
- **Functionality config profiles** inspect uploaded reference settings across servers/operators and show searchable reference/actual values, differences, missing settings and masked secrets. INI/config/Python parsing is static; incomplete checks remain explicit.
- **Bitwarden** account sign-in/unlock and direct native SSH item references for SSH/SFTP. Private key retrieval happens on the server during authentication. Master passwords and unlock sessions are not persisted; per-user CLI profiles are encrypted between commands.
- **Folder order** persists after drag/drop, Move up/down and export/import, including nested folders.

#### Fixed and changed

- Git dialogs now stay visible above mobile drawers and keep wide tables within their scroll area.
- The bridge supports original same-origin no-referrer forms and rejects cross-site submissions; stop refuses active original-tool operations.
- Docker includes Python, git, OpenSSH, rsync and the official Bitwarden CLI. Standalone installations need these optional dependencies on the WRM host. Medium defaults to administrators; grant it only to trusted operators because its native build/config commands run as the WRM OS account.

#### Validation and packages

Original tool: **1,581 tests passed, 9 skipped**. Real Chromium exercised private archive imports and native comparison, deploy, history and rollback against test Git/SSH destinations. Config profiles, folder persistence and layouts at 390/1100/1440/1920 px were checked without JavaScript errors. Bitwarden mock-account flows and direct real SSH authentication were tested; a real cloud account was not used. Full WRM race/vet/format/UI checks and platform builds are detailed in the [release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/docs/RELEASE_SUMMARY_v12.4.3.md).

Release packages: 13 standalone targets, macOS universal and Linux amd64/arm64 AppImages, plus `SHA256SUMS.txt`. Docker: `ghcr.io/vedranius/wrm-pro:v12.4.3` and `latest`. The pipeline publishes only after the complete package set and Docker succeed.

#### Upgrading from v12.4.0

Replace the binary/AppImage or update Docker. Startup adds the config-profile table and folder order column. Preserve your database, encryption key and private adjacent `.medium`/`.bitwarden` directories. Tool and config archives are private imports, not public release assets. Git/SSH targets must be reachable from the runtime. The earlier client GitLab relay/replay, AI/navigation and icon/color requests remain accepted backlog, and were not included in this prioritized release.

**Downloads:** [v12.4.3 packages and checksums](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.3). [Changelog](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/CHANGELOG.md).

### Included from v12.4.0 — Service functionality, nested folders and remote previews

Inspect service functionality and configuration on your servers, distinguish identically named connections by client/site, and preview remote files directly in WRM.

#### Added

- **Git → Functionality** checks every Supervisor worker and failed systemd service on selected servers, as well as each catalog service's runtime units. It shows installation problems, disk/memory warnings, timestamps and bounded error logs. Optional HTTP/TCP health probes run from the selected server. Unavailable checks remain unverified; problem/recovery notifications are available.
- **Configuration** shows section, variable/function, value and description for each service's config/INI/environment files and JSON. Secrets are masked. Paths stay inside the installation; contents are read on demand with explicit limits and parse-gap warnings.
- **Nested folders** ([#57](https://github.com/vedranius/web-browser-RDM-public/issues/57)): recursive tree, parent selection, complete folder paths in connection pickers/windows/tabs, folder tab filtering and hierarchy-preserving imports. Folder shares include descendants; route defaults remain local to the immediate folder.
- **Remote previews** ([#55](https://github.com/vedranius/web-browser-RDM-public/issues/55)): video/audio streaming with seekable byte ranges over SFTP/FTP, image/PDF viewing and bounded text viewing. Select a path in a terminal and right-click to preview/download it; the file manager offers Preview too. Browser codec support determines playable media formats.

#### Fixed

- **AI entry points and diagnostics** ([discussion #56](https://github.com/vedranius/web-browser-RDM-public/discussions/56)): AI is visible in SSH/SFTP windows; unavailable permissions lead to an explanation. Settings → AI assistant and the top AI button provide setup/status access. Policy changes refresh other open browsers. A successful provider test is distinguished from account permission, and the terminal assistant requires no MCP setup.
- Missing installations stay visible with a warning. Supervisor programs in the same configuration file are associated with the correct service stanza.
- Docker Compose uses the current version consistently.

#### Packages and validation

Linux, Windows, macOS (including universal), Android/Termux, FreeBSD and OpenBSD standalone binaries; Linux amd64/arm64 AppImages; Docker `ghcr.io/vedranius/wrm-pro:v12.4.0` and `latest`. The release pipeline verifies all 16 downloadable packages and `SHA256SUMS.txt` before publication.

Validation details: [v12.4.0 release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.0/docs/RELEASE_SUMMARY_v12.4.0.md).

#### Upgrading from v12.3.6

Replace the binary/AppImage or update the container. Startup adds `folders.parent_id` and a service-health cache table without removing existing data. Previous binaries ignore the new fields; run the previous binary to downgrade. Existing catalogs need no new fields. Git runtime checks use the SSH account's read permissions and do not restart services.

**Downloads:** [v12.4.0 packages and checksums](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.0). See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.0/CHANGELOG.md) for the complete history.

### Included from v12.3.6

This release adds deployment review and server selection to Git Simple, gives AI tools the server's actual date, and improves Linux systemd service installation and failure diagnostics.

#### 🤖 AI assistant

- Ollama setup releases its action lock before reporting completion, so the next setup step can start immediately. This fixes the intermittent HTTP 409 exposed by the release CI run.

- `system_info` now reports the remote server's actual date and time instead of leaving the model to infer it from the kernel build string.
- The assistant is explicitly instructed not to invent tool calls, tool output, logs or internal `<tool_response>` wrappers.
- Enabling or disabling the assistant refreshes its buttons on already-open terminals without reconnecting. Shared connection roles remain enforced; legacy lowercase SSH protocols are accepted.
- The processing indicator now has its English and Croatian translation.
- The indicator remains visible until the model finishes. Recognized tool JSON is removed from the finalized chat message, and incomplete JSON is no longer repaired into executable calls or silently removed from answers.

#### ⎇ Git Simple

- Git Simple now has combined service, hostname, environment, status, branch and version/commit filters. Its deploy review shows changed files and their state per destination, lets you select successfully compared servers and exclude files, and shows when server state was compared. Deploys keep using WRM's fingerprinted plan, backup, lock and audit flow.
- Empty server selections are refused; confirmation is recomputed whenever the selection changes. Canceled or superseded comparisons cannot reopen a review. Rollback shows each selected server's own previous version. Simple mode keeps removed files; Advanced mode offers deletion after backup.
- Connection changes after comparison invalidate the plan, so a reviewed deployment cannot silently follow a changed SSH destination or login.

#### 🐧 Linux service

- Installation checks read/write access to SQLite and its WAL/SHM files, plus write/search access to its directory, before changing service settings. It points to `DB_PATH` when permissions prevent startup.
- If systemd cannot start WRM, installation reports the unit status and recent journal entries; read-only database errors include a specific permissions hint.
- Generated units limit rapid restart loops after a persistent failure. The install summary shows the database path the service will use.
- Service definitions handle spaces, literal dollars and percent characters in paths correctly. Reinstalling clears the unit's old failures and restarts the newly configured binary.

#### 📦 Packages and validation

- **Linux**: amd64, arm64, ARMv7, ARMv6, 386; **AppImage** for amd64 and arm64.
- **Windows**: amd64 and arm64 `.exe` files.
- **macOS**: Intel, Apple Silicon and universal binaries.
- **Android / Termux**: arm64; **FreeBSD**: amd64 and arm64; **OpenBSD**: amd64.
- Docker: `ghcr.io/vedranius/wrm-pro:v12.3.6` and `latest` (Linux amd64).
- `SHA256SUMS.txt` covers all 16 downloadable packages. GitHub also provides source archives.

AppImage quick start: `chmod +x wrm-pro-v12.3.6-linux-amd64.AppImage`, then run it and open **http://localhost:8080**. On ARM64 choose the arm64 AppImage. Without FUSE, add `--appimage-extract-and-run`. Data stays in `${XDG_DATA_HOME:-$HOME/.local/share}/wrm-pro/` unless `DB_PATH` is set. Replace the AppImage to upgrade; use the standalone binary for boot services.

The pre-push review covered race tests, UI regressions, all 13 cross-compilation targets, native Linux/Docker health checks, Git Simple browser deployment and rollback, and a real local Ollama tool call. Other OS binaries were compiled; they were not run on their respective OSes. See [the review report](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/docs/PRE_PUSH_REVIEW_v12.3.5.md).

#### ⬆️ Upgrading from v12.3.4

Replace the binary or use the in-app updater. No database migration is required.

For an existing Linux service, re-run `-install-service` with the new binary to regenerate its unit and apply the path and restart-limit fixes. This restarts that service.

**Downloads:** [all v12.3.6 files and `SHA256SUMS.txt`](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.3.6)

See [the release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/docs/RELEASE_SUMMARY_v12.3.6.md) for package and validation details.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/CHANGELOG.md) for the complete change history.
