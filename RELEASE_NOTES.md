<p align="center"><img src="https://raw.githubusercontent.com/vedranius/web-browser-RDM-public/v12.4.0/docs/brand/png/lockup/wrm-lockup-on-dark-664w.png" alt="WRM PRO — Web Remote Manager" width="332"></p>

## Web Remote Manager PRO v12.4.0 — Service functionality, nested folders and remote previews

Inspect service functionality and configuration on your servers, distinguish identically named connections by client/site, and preview remote files directly in WRM.

### Added

- **Git → Functionality** checks every Supervisor worker and failed systemd service on selected servers, as well as each catalog service's runtime units. It shows installation problems, disk/memory warnings, timestamps and bounded error logs. Optional HTTP/TCP health probes run from the selected server. Unavailable checks remain unverified; problem/recovery notifications are available.
- **Configuration** shows section, variable/function, value and description for each service's config/INI/environment files and JSON. Secrets are masked. Paths stay inside the installation; contents are read on demand with explicit limits and parse-gap warnings.
- **Nested folders** ([#57](https://github.com/vedranius/web-browser-RDM-public/issues/57)): recursive tree, parent selection, complete folder paths in connection pickers/windows/tabs, folder tab filtering and hierarchy-preserving imports. Folder shares include descendants; route defaults remain local to the immediate folder.
- **Remote previews** ([#55](https://github.com/vedranius/web-browser-RDM-public/issues/55)): video/audio streaming with seekable byte ranges over SFTP/FTP, image/PDF viewing and bounded text viewing. Select a path in a terminal and right-click to preview/download it; the file manager offers Preview too. Browser codec support determines playable media formats.

### Fixed

- **AI entry points and diagnostics** ([discussion #56](https://github.com/vedranius/web-browser-RDM-public/discussions/56)): AI is visible in SSH/SFTP windows; unavailable permissions lead to an explanation. Settings → AI assistant and the top AI button provide setup/status access. Policy changes refresh other open browsers. A successful provider test is distinguished from account permission, and the terminal assistant requires no MCP setup.
- Missing installations stay visible with a warning. Supervisor programs in the same configuration file are associated with the correct service stanza.
- Docker Compose uses the current version consistently.

### Packages and validation

Linux, Windows, macOS (including universal), Android/Termux, FreeBSD and OpenBSD standalone binaries; Linux amd64/arm64 AppImages; Docker `ghcr.io/vedranius/wrm-pro:v12.4.0` and `latest`. The release pipeline verifies all 16 downloadable packages and `SHA256SUMS.txt` before publication.

Validation details: [v12.4.0 release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.0/docs/RELEASE_SUMMARY_v12.4.0.md).

### Upgrading from v12.3.6

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
