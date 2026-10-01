# Architecture

This document describes how Web Remote Manager PRO (WRM) is built. It is the result of
"Step 0 — repo audit" of the extension plan (WRM as a lightweight PAM for teams), and it maps
every phase of that plan to the code that it touches.

## Stack

| Part | What it is |
|---|---|
| Backend | **Go** (version in `remote-manager/go.mod`), standard library `net/http`, no web framework. One self-contained binary. |
| Libraries | `golang.org/x/crypto/ssh` (SSH client), `github.com/pkg/sftp`, `github.com/jlaffaye/ftp`, `github.com/gorilla/websocket`, `github.com/pion/turn/v4` (TURN relay for voice), `modernc.org/sqlite` (pure-Go SQLite, no CGO), `rsc.io/qr` |
| Frontend | One page, `remote-manager/static/index.html`: plain JavaScript, no framework, no build step. `xterm.js` and fonts are vendored in `static/vendor/`. Everything under `static/` is embedded into the binary with `go:embed`. |
| Database | SQLite in WAL mode (`DB_PATH`, default `./remote_manager.db`), mode `0600`. |
| Build / run | `go build`; CI (`.github/workflows/build.yml`) runs gofmt, vet, tests with the race detector, a JS syntax check, cross-compiles 13 targets and publishes releases on tags. `Dockerfile` + `docker-compose.yml` for container deployments. |

## Auth model

- Local accounts in `users` (bcrypt cost 12). The first account is the administrator; self-registration is closed by default.
- Sign-in sessions in `auth_sessions`: a random 256-bit token in the `wrm_session` cookie (HttpOnly, SameSite, Secure over HTTPS). Only the SHA-256 of the token is stored. Sessions expire after idle and absolute timeouts (policies).
- Two-factor authentication: TOTP (RFC 6238) with single-use recovery codes; a policy can require it for admins or everyone (`users.go`, `totp.go`).
- Roles today: `users.is_admin` (administrator or user) plus per-share roles (see Sharing). Policies are in `app_settings` and can be forced with `WRM_<KEY>` environment variables (`settings.go`).
- CSRF: every non-GET `/api` request needs `X-WRM-Request: 1` and a same-origin `Origin`; WebSockets check `Origin` (`security.go`).

## SSH / connection layer

- **Terminal**: `GET /ws/ssh?id=…` (`ssh_ws.go`). WRM opens the SSH connection itself (`x/crypto/ssh`), requests a PTY and starts a shell. Two pump goroutines read stdout/stderr and send binary WebSocket frames; a writer goroutine forwards keystrokes. The **PTY output stream lives here**. This is where session recording taps in (`termAudit.output/input/resize`).
- **Files**: `files.go` (list, download, ZIP, streamed multipart upload, mkdir, rename, delete) over a small pool of SFTP clients (`sftp_pool.go`) or FTP/FTPS; `search.go` (name and content search); `conntest.go` (test a connection).
- **Host keys**: `hostkeys.go` verifies SSH host keys (trust on first use or strict) and pins FTPS certificates (`known_hosts` table).
- **Server-to-server SFTP**: `transferRemoteHandler` in `main.go` opens SFTP to both servers from WRM and copies through a 1 MiB buffer, streaming NDJSON progress to the browser.

## Stored connections and credentials

- `connections` (host, port in `host`, username, auth method, `password`, `private_key`, `key_path`, folder, owner `user_id`) and `folders`. This is the "servers" inventory of the plan.
- Passwords, private keys and TOTP secrets are encrypted with **AES-256-GCM** before they are stored (`security.go`). The key comes from `ENCRYPTION_KEY`, `ENCRYPTION_KEY_FILE` or a random key file `<DB_PATH>.key` (mode `0600`), never from the database or git.
- Secrets are never returned by the API: the UI only "sets/replaces" them. They are decrypted in memory when connecting.
- Server-side key files (`key_path`, `~/.ssh`) are limited to administrators by default.

## Canvas and collaboration

- The multi-window canvas (windows, tabs, snapping) is **client-side state**. Saved workspaces are stored on the server (`sessions` table, `/api/sessions`).
- **Sharing** (`shares.go`): `share_links`, `share_items` (connections/folders), `share_members` and `share_participants`. Five roles (observer, viewer, operator, moderator, owner) with permissions, enforced on every request by `authorizeConnection(r, connID, perm)`.
- **Realtime** (`collab.go`): `/ws/share/{token}` rooms. The server is the authority for identities, roles, presence, chat, terminal sharing (the sharer's browser streams output to watchers), keyboard control grants and voice signalling. `/ws/events` carries per-user notifications.

## Tests and migrations

- `go test ./...`:
  - `security_test.go`: TOTP, recovery codes, roles, share access, CSRF, encryption and key migration, settings;
  - `upgrade_test.go`: upgrading databases from earlier builds, old URLs;
  - `integration_test.go`: an in-process SSH + SFTP server. It covers connect → audit entries → recording, file transfers with checksums, and the append-only audit trail.
- **Migrations** run at startup in `initDB` (`main.go`), and they are idempotent: `CREATE TABLE IF NOT EXISTS`, `ALTER TABLE … ADD COLUMN`, indexes and triggers. `ensureTable` keeps a table left by an earlier build with an incompatible layout as `<name>_old_<time>` and creates it again.
- Migrations are **additive only**: new tables, columns, indexes and triggers, and no data is removed or rewritten. A previous binary ignores the new objects, so a downgrade is "run the previous binary". This is the reversible ("down") path, and it loses no data.

## Feature flags

Feature flags are settings in `app_settings`. An administrator changes them in *Admin panel → Security policies*, and an environment variable can force them (`WRM_<KEY>`, some also accept the plan's names):

| Flag | Default | Environment |
|---|---|---|
| `audit_enabled` | on | `WRM_AUDIT_ENABLED` or `AUDIT_ENABLED` |
| `session_recording` | on | `WRM_SESSION_RECORDING` or `SESSION_RECORDING_ENABLED` |
| `session_recording_input` | off | `WRM_SESSION_RECORDING_INPUT` |
| `recording_retention_days` / `recording_max_mb` | 90 / 100 | `WRM_RECORDING_RETENTION_DAYS` / `WRM_RECORDING_MAX_MB` |

Later phases add their flags (`RBAC_ENABLED`, `OIDC_ENABLED`, `BROADCAST_ENABLED`, `FLEET_ENABLED`, …) the same way. With all new flags off, the existing flow (connect → work) is unchanged.

## Extension plan: status and touch points

| Phase | Already in WRM (v10.0) | Added / to add | Code it touches |
|---|---|---|---|
| **1. Audit + recording** | `audit_log` of sign-ins, admin actions, connections, shares, terminals, file changes; CSV export | **v10.1:** `terminal_sessions`, `session_recordings` (asciicast v2, replay), `file_transfers` (size + SHA-256 for upload/download/ZIP/server-to-server), audit entries linked to connection and session, hash chain + append-only triggers, filters, `/api/recordings`, `/api/admin/transfers`, `/healthz` | `audit.go`, `recording.go` (new), `ssh_ws.go`, `files.go`, `main.go` (schema, transfer), `settings.go`, `index.html` |
| **2. RBAC** | Owner-only connections, share roles with permissions, server-side checks (`authorizeConnection`) | Groups, central inventory with `server_access` grants (user/group → role), deny-by-default, one `authorize(user, action, server)` guard, `access_denied` audit, access matrix UI | `shares.go` → new `rbac.go`, connection API in `main.go`, `users.go`, admin UI |
| **3. Credentials, 2FA, SSO** | AES-256-GCM at rest, key from env/file, secrets never returned, TOTP + recovery codes, 2FA policy | Key IDs and key rotation, refusing to start when encrypted secrets exist but the key is missing, per-user SSH keys, SSH CA with short-lived certificates, OIDC SSO (Authorization Code + PKCE) with group mapping | `security.go`, `users.go`, new `oidc.go`, `buildAuthMethods` |
| **4. Collab: broadcast + live join** | Live rooms with presence, roles, terminal sharing, request/grant/revoke control (read-only / read-write / owner) | Broadcast input to several sessions (with confirmation, destructive-command warning, `broadcast_cmd` audit) | `collab.go`, `ssh_ws.go`, `index.html` |
| **5. Fleet operations** | — | Live up/down status, tags/environments, inventory import (SSH config, GitLab…), jump hosts (`jump_path` already in `terminal_sessions`), tunnels | new `fleet.go`, connections schema, SSH dialing |
| **6. Quick wins** | Inline SFTP editor, quick connection search, responsive/mobile UI | Snippets / run-on-connect, tmux-backed sessions, session notes, installable PWA | `index.html`, new tables |

### Decisions where the plan's example names collide with existing code

- **`terminal_sessions`** instead of `sessions`, because `sessions` already stores saved workspaces.
- **`/api/recordings`** instead of `/api/sessions`, because `/api/sessions` is the workspace API.
- The plan's `servers` are the existing **`connections`** (each with an owner). Phase 2 adds grants on top of them.
- `audit_log` keeps its existing column names (`action` = the plan's `event_type`, `details` = `detail_json`, `ip` = `client_ip`) and gains `conn_id` (server), `session_id`, `prev_hash` and `hash`.
- Recordings are **gzip-compressed asciicast v2** files on disk (`WRM_RECORDINGS_DIR`, default `recordings/` next to the database). The checksum is the SHA-256 of the stored file. The API serves the plain `.cast`.
