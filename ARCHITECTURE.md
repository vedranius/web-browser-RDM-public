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
- **Dialing and jump hosts** (`jump.go`): every SSH/SFTP operation (terminal, SFTP pool, search, transfer, connection test, tunnels) gets its client from `dialSSH(c, onNewKey)`. When `connections.jump_conn_id` is set, `jumpChain` resolves the chain (same owner, SSH only, at most 5 hops, no loops), the first hop is dialed with `ssh.Dial`, each further hop with `prev.Dial("tcp", host)` (a `direct-tcpip` channel) plus `ssh.NewClientConn`. Each hop uses its own credentials and host key check. When the target client closes, the hops close too. FTP/FTPS behind a jump host use `ftp.DialWithDialFunc` with the last hop's `Dial`. `validateJump` checks a choice when a connection is saved.
- **Tunnels** (`tunnels.go`): definitions in `connection_tunnels` (kind `local` / `remote` / `dynamic`, listen host:port, target host:port, open scheme/path, start mode `manual` / `connect` / `always`). The `tunnelManager` keeps running tunnels by key (`t<id>`, or `w<conn>u<user>` for temporary web tunnels). A running tunnel has its own SSH client (via `dialSSH`), a listener (local: `net.Listen` on the WRM machine → `client.Dial`; dynamic: a small SOCKS5 server (CONNECT, no auth) → `client.Dial`; remote: `client.Listen` (`tcpip-forward`) → `net.Dial` on the WRM side), per-tunnel counters, a 30 s keepalive and automatic reconnect with backoff while the listener stays open. `connect` tunnels are reference-counted by the terminals of the connection (`holdConn` / `releaseConn` with a grace period); `always` tunnels start after the database at boot. Changes are pushed to the owner's browsers through `/ws/events` (`tunnels_changed`). Policies are checked when a tunnel is saved, started, and again when a policy changes (`applyPolicy`).
- **Remote desktop** (`guac.go`): RDP, VNC and Telnet go through guacd. `/ws/desktop` authorizes the connection like a terminal, opens a temporary tunnel through the jump hosts when needed, connects to `guacd_address`, performs the handshake (`select` → `args` → `size`/`audio`/`video`/`image`/`timezone`/`connect` → `ready`) with the stored parameters (`connections.options`, JSON) and then relays complete instructions: guacd → browser in batches (recorded to `recordings/…/<uid>.guac.gz`), browser → guacd only allow-listed opcodes. Sessions use `terminal_sessions` (`protocol` rdp/vnc/telnet) and `session_recordings` (`format` `guacamole+gzip`); the browser replays them with `Guacamole.SessionRecording`.
- **Terminal backends** (`console.go`): the terminal WebSocket talks to a `termIO` (stdin, stdout, start, resize, wait, keepalive, close): an SSH shell, the serial console over the BMC's SSH (the console command is typed once the BMC prompt is quiet), IPMI Serial-over-LAN (`ipmitool sol activate` on the last jump host — the password goes over stdin after the wrapper turned echo off — or locally in a pseudo terminal, `tty_linux.go`), or a serial port of the WRM server (raw termios, `tty_linux.go`).
- **Out-of-band management** (`bmc.go`): `connections.bmc` holds a JSON BMC setting. Redfish requests use an `http.Transport` whose `DialContext` goes through the last jump hop when needed and whose TLS check pins self-signed certificates (`pinnedTLSConfig`, shared with FTPS, stored in `known_hosts` as `bmc://host:port`). IPMI runs `ipmitool` with `-E`, locally or on the last hop. Power actions map to Redfish `ResetType`s / `ipmitool chassis power`.
- **Quick connect & notes** (`quick.go`): a quick connection is an ordinary row in `connections` with `temp_until` set (RFC 3339, 24 hours after the last use). It goes through `normalizeConnection`, `validateJump` and `checkAuthRefs` like a saved one, is created with `monitor=0`, is excluded from export, and `touchQuick` moves `temp_until` whenever a terminal or desktop opens it. `runQuickCleanup` deletes expired rows every 10 minutes unless a terminal of the connection is open; `PUT … {temporary: false}` clears `temp_until`. `connections.notes` holds the owner's notes; lists only return `has_notes`, `GET /api/connections/{id}` returns the text, the browser renders a small markdown after escaping everything.
- **Network tools** (`nettools.go`): `POST /api/nettools` validates the target (host name / address characters only), takes a per-user lock and runs the tool either on the WRM server (`net.Dial`, `exec.CommandContext` with arguments, no shell, `net.Resolver`, `http.Client`) or from a source connection of the user through `dialSSH`: port and HTTP checks use `client.Dial` (`direct-tcpip`) so nothing runs on the server, ping / traceroute / DNS run a fixed POSIX script with the shell-quoted target (`runRemoteScript`). The HTTP check does not follow redirects and verifies the certificate separately against the system roots to report trust without failing.
- **Installable app** (`pwa.go`): `/manifest.webmanifest` and `/sw.js`. The service worker cache is named after `AppVersion` (`wrm-<version>`), precaches the vendor scripts and brand files, serves `/static/` cache-first, never touches `/api/` or WebSockets, and falls back to an offline page for navigations when the server is unreachable.
- **Tags & inventory** (`inventory.go`): `connections.tags` holds normalised, comma-joined tags. CSV files are decoded (UTF-8 / UTF-16 / Windows-1250) and split with `encoding/csv` after separator detection; `.xlsx` workbooks are read with `archive/zip` and `encoding/xml` (workbook, relationships, shared strings, sheet). Rows become `importItem`s through a field → column mapping and go through the same `importer.commit` as the mRemoteNG / OpenSSH imports. NetBox objects are fetched page by page (`/api/dcim/devices/`, `/api/virtualization/virtual-machines/`), get `ext_id` `netbox:<host>:<kind>:<id>` and `ext_tags` (the tags NetBox controls), so a later import updates them instead of duplicating; `inventory_sources` keeps the address and the encrypted token per user.
- **SSH key store** (`keys.go`): table `ssh_keys` (owner, type, public line, encrypted private key, fingerprint) and `ssh_key_deployments` (key ↔ connection). Connections with `auth_method` `KEY_REF` point to a key (`connections.key_id`). Deploy / revoke / "who has access" log in with `dialSSH` and run small POSIX `sh` + `awk` scripts (`runRemoteScript`) that read the key or its base64 blob from stdin and change `~/.ssh/authorized_keys` in place.
- **Credentials vault** (`credentials.go`): tables `credentials` (owner, user name, encrypted password and pending password, optional `key_id`, host patterns, `shared_all`, rotation interval and status) and `credential_grants`. Connections with `auth_method` `CREDENTIAL` point to one (`connections.credential_id`). `loadConnection` calls `resolveConnectionAuth`, which replaces `KEY_REF` / `CREDENTIAL` by `KEY` / `PASSWORD` with the decrypted secret after checking ownership, grants and host patterns (for grantees also the jump hosts); `dialSSH` resolves every hop the same way, and a failure is kept in `authErr` and reported at login. Rotation (`startRotation` / `runRotation`) runs in the background: pre-flight logins, `passwd` over a PTY per user@host, verification, rollback; the browser polls `/api/credentials/{id}/rotation`.
- **Snippets** (`snippets.go`): table `snippets` (owner, scope all/folder/connection, `auto_run`, `shared`). Variables are expanded per terminal (browser for manual use, server for run on connect). `ssh_ws.go` types run-on-connect snippets into stdin once the first output has been quiet for 400 ms (at most 5 s after the shell started).
- **Broadcast input** is done by the browser (`termInput` sends the same keystrokes to every terminal of the group); each terminal WebSocket receives a `broadcast` control message so the server can audit it per terminal session and refuse it when `broadcast_enabled` is off.
- **Live status** (`status.go`): `statusMonitor.round()` loads all monitored connections, builds one target per host:port (or per jump host and host:port when `status_jump_checks` is on), probes them with 32 workers (TCP connect, SSH/FTP greeting, one retry), keeps state and "since" per target, derives the state of connections behind jump hosts from the jump host, and sends `status_changed` to owners whose connections changed.
- **Web interface connections**: protocol `HTTP` / `HTTPS` with `web_path`. `POST /api/connections/{id}/open-web` returns the URL directly, or — behind a jump host — opens a temporary local tunnel over the last hop (30 minutes idle limit) and returns its URL.
- **Import** (`importers.go`): mRemoteNG `confCons.xml` (AES-GCM with PBKDF2-HMAC-SHA1 key / legacy AES-CBC with MD5 key; `Protected` check of the master password; full-file encryption; inheritance; `SSHTunnelConnectionName` → jump host) and OpenSSH config (`Host`, `Host *`, `ProxyJump`, `ProxyCommand ssh -W`, `*Forward`, `IdentityFile`). Both go through one `importer` that skips duplicates, links jump hosts by name and validates tunnels.
- **Server-to-server SFTP**: `transferRemoteHandler` in `main.go` opens SFTP to both servers from WRM and copies through a 1 MiB buffer, streaming NDJSON progress to the browser.

## Stored connections and credentials

- `connections` (host, port in `host`, username, auth method, `password`, `private_key`, `key_path`, `key_id`, `credential_id`, folder, owner `user_id`, `jump_conn_id`, `web_path`, `tags`, `ext_id` / `ext_tags` for inventory sync), `connection_tunnels` and `folders`. This is the "servers" inventory of the plan. Keys of the key store (`ssh_keys`) and vault credentials (`credentials`, `credential_grants`) are shared by many connections.
- Passwords, private keys and TOTP secrets are encrypted with **AES-256-GCM** before they are stored (`security.go`). The key comes from `ENCRYPTION_KEY`, `ENCRYPTION_KEY_FILE` or a random key file `<DB_PATH>.key` (mode `0600`), never from the database or git.
- Secrets are never returned by the API: the UI only "sets/replaces" them. They are decrypted in memory when connecting. The only exceptions are explicit, audited owner actions after re-authentication (config export with secrets, private key export, *Show* of a vault password), all behind the `allow_secret_export` policy.
- Server-side key files (`key_path`, `~/.ssh`) are limited to administrators by default.

## Canvas and collaboration

- The multi-window canvas (windows, tabs, snapping) is **client-side state**. Saved workspaces are stored on the server (`sessions` table, `/api/sessions`).
- **Sharing** (`shares.go`): `share_links`, `share_items` (connections/folders), `share_members` and `share_participants`. Five roles (observer, viewer, operator, moderator, owner) with permissions, enforced on every request by `authorizeConnection(r, connID, perm)`.
- **Realtime** (`collab.go`): `/ws/share/{token}` rooms. The server is the authority for identities, roles, presence, chat, terminal sharing (the sharer's browser streams output to watchers), keyboard control grants and voice signalling. `/ws/events` carries per-user notifications.

## Tests and migrations

- `go test ./...`:
  - `security_test.go`: TOTP, recovery codes, roles, share access, CSRF, encryption and key migration, settings;
  - `upgrade_test.go`: upgrading databases from earlier builds, old URLs;
  - `integration_test.go`: an in-process SSH + SFTP server (also `direct-tcpip` and `tcpip-forward`). It covers connect → audit entries → recording, file transfers with checksums, and the append-only audit trail;
  - `tunnels_test.go`: two-hop jump chains (terminal, SFTP, connection test, loop/ownership checks), local / SOCKS5 / remote tunnels with traffic counters, ownership, policies, audit, start modes and web interfaces behind a jump host;
  - `snippets_test.go`: snippets API, sharing and validation, variables, run on connect (order and scope), export/import, broadcast audit and policy;
  - `guac_test.go`: Guacamole codec and the desktop relay against a fake guacd (handshake parameters, input allow-list, pings, errors, policy, recording, VNC through a jump host);
  - `quick_test.go`: quick connect (targets and schemes, temporary flag, expiry with and without an open terminal, per-user limit, export), notes (save, length limit, list flag, duplicate), network tools (ports and HTTP/TLS from WRM and through an SSH connection, DNS, ping and traceroute in a real `sh`, validation, ownership, policy), manifest and service worker;
  - `bmc_test.go`: Redfish against a fake BMC (status, actions, pinning, jump host, vault credential), IPMI with a fake ipmitool (locally and on a jump host), Serial-over-LAN, the BMC SSH console, a serial port on a pseudo terminal;
  - `inventory_test.go`: tags (normalisation, bulk, export/import), CSV separators and encodings, an `.xlsx` import (mapping, folders, tags, jump hosts, credential login), NetBox import and sync against a fake NetBox (pagination, token, filters, update, missing objects);
  - `keys_test.go`: key store (generate, import with passphrase, public keys, encryption at rest), deploy / revoke / who has access against a fake SSH server that runs commands in a real `sh`, vault sharing and host / jump-host restrictions, rotation (check, success, rollback, incomplete rollback) with a simulated `passwd`;
  - `status_test.go`: live status (up/down, SSH banner, jump hosts, checks through jump hosts, check now, opt-out, FTP greeting);
  - `importers_test.go`: mRemoteNG files (both encryption formats, master password, inheritance, folders, jump hosts) and OpenSSH config.
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
| `tunnels_enabled` | on | `WRM_TUNNELS_ENABLED` or `TUNNELS_ENABLED` |
| `broadcast_enabled` | on | `WRM_BROADCAST_ENABLED` or `BROADCAST_ENABLED` |
| `status_enabled` / `status_interval_seconds` / `status_jump_checks` | on / 60 / off | `WRM_STATUS_ENABLED` or `STATUS_ENABLED`, `WRM_STATUS_INTERVAL_SECONDS`, `WRM_STATUS_JUMP_CHECKS` |
| `tunnel_users` / `tunnel_bind_any` / `tunnel_remote_forward` / `tunnel_idle_minutes` | all / off / admins / 0 | `WRM_TUNNEL_USERS`, `WRM_TUNNEL_BIND_ANY`, `WRM_TUNNEL_REMOTE_FORWARD`, `WRM_TUNNEL_IDLE_MINUTES` |
| `bmc_enabled` / `serial_ports` | on / admins | `WRM_BMC_ENABLED` or `BMC_ENABLED`, `WRM_SERIAL_PORTS` |
| `network_tools` | all | `WRM_NETWORK_TOOLS` |

Later phases add their flags (`RBAC_ENABLED`, `OIDC_ENABLED`, `BROADCAST_ENABLED`, `FLEET_ENABLED`, …) the same way. With all new flags off, the existing flow (connect → work) is unchanged.

## Extension plan: status and touch points

| Phase | Already in WRM (v10.0) | Added / to add | Code it touches |
|---|---|---|---|
| **1. Audit + recording** | `audit_log` of sign-ins, admin actions, connections, shares, terminals, file changes; CSV export | **v10.1:** `terminal_sessions`, `session_recordings` (asciicast v2, replay), `file_transfers` (size + SHA-256 for upload/download/ZIP/server-to-server), audit entries linked to connection and session, hash chain + append-only triggers, filters, `/api/recordings`, `/api/admin/transfers`, `/healthz` | `audit.go`, `recording.go` (new), `ssh_ws.go`, `files.go`, `main.go` (schema, transfer), `settings.go`, `index.html` |
| **2. RBAC** | Owner-only connections, share roles with permissions, server-side checks (`authorizeConnection`) | Groups, central inventory with `server_access` grants (user/group → role), deny-by-default, one `authorize(user, action, server)` guard, `access_denied` audit, access matrix UI | `shares.go` → new `rbac.go`, connection API in `main.go`, `users.go`, admin UI |
| **3. Credentials, 2FA, SSO** | AES-256-GCM at rest, key from env/file, secrets never returned, TOTP + recovery codes, 2FA policy | **v10.5:** per-user SSH key store with deploy / revoke / who has access (`keys.go`), shared credentials vault with grants, host lists and password rotation (`credentials.go`). **To add:** encryption key IDs and rotation, refusing to start when encrypted secrets exist but the key is missing, SSH CA with short-lived certificates, OIDC SSO (Authorization Code + PKCE) with group mapping | `security.go`, `users.go`, `keys.go`, `credentials.go`, `helpers.go` (`loadConnection`), `jump.go`, new `oidc.go` |
| **4. Collab: broadcast + live join** | Live rooms with presence, roles, terminal sharing, request/grant/revoke control (read-only / read-write / owner) | **v10.3:** broadcast input to several terminals (confirmation, destructive-command warning, `terminal.broadcast` audit per session, policy). **To add:** live join of other users' sessions | `ssh_ws.go`, `settings.go`, `index.html` |
| **5. Fleet operations** | — | **v10.2:** jump hosts (chains, all operations, `terminal_sessions.jump_path`), SSH tunnels (local / remote / SOCKS5, start modes, reconnect, panel, policies, audit), web interface connections, import from mRemoteNG and OpenSSH config. **v10.3:** live up/down status (`status.go`). **v10.6:** tags and environments, inventory import from CSV / Excel and NetBox with sync (`inventory.go`). **To add:** GitLab / other inventories | `jump.go`, `tunnels.go`, `importers.go` (new), `main.go` (schema, connections API), `ssh_ws.go`, `helpers.go`, `admin.go`, `users.go`, `files.go`, `conntest.go`, `settings.go`, `index.html` |
| **6. Quick wins** | Inline SFTP editor, quick connection search, responsive/mobile UI. **v10.4:** RDP / VNC / Telnet through guacd (`guac.go`) | **v10.3:** snippets and run on connect (`snippets.go`). **v10.8:** quick connect, connection notes / runbooks (`quick.go`), network tools (`nettools.go`), installable PWA (`pwa.go`). **To add:** tmux-backed sessions | `index.html`, new tables |

### Decisions where the plan's example names collide with existing code

- **`terminal_sessions`** instead of `sessions`, because `sessions` already stores saved workspaces.
- **`/api/recordings`** instead of `/api/sessions`, because `/api/sessions` is the workspace API.
- The plan's `servers` are the existing **`connections`** (each with an owner). Phase 2 adds grants on top of them.
- `audit_log` keeps its existing column names (`action` = the plan's `event_type`, `details` = `detail_json`, `ip` = `client_ip`) and gains `conn_id` (server), `session_id`, `prev_hash` and `hash`.
- Recordings are **gzip-compressed asciicast v2** files on disk (`WRM_RECORDINGS_DIR`, default `recordings/` next to the database). The checksum is the SHA-256 of the stored file. The API serves the plain `.cast`.
