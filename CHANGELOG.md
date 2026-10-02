# Changelog

All notable changes to Web Remote Manager PRO. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Release notes with downloads are on the
[Releases page](https://github.com/vedranius/web-browser-RDM-public/releases).

## [10.8.0] — 2026-10-02 — quick connect, notes, network tools & installable app

### Added
- **Quick connect** (`quick.go`): a field at the top of the sidebar takes `user@host:port` or a URL-like target (`ssh://`, `sftp://`, `rdp://`, `vnc://`, `telnet://`, `https://`) and connects without the connection dialog:
  - a small dialog for the login (password, vault credential or stored key) and an optional jump host; `POST /api/connections/quick`;
  - the connection is a **temporary** connection of its owner (`connections.temp_until`) — same checks, host keys, recording and audit (`connection.quick`) — listed under *Quick connections*, kept 24 hours after its last use (open terminals keep it alive), at most 20 per user, not monitored, not exported;
  - 💾 / *Save as a connection* turns it into a normal connection (`PUT … {temporary: false}`).
- **Notes & runbook** per connection (`connections.notes`, up to 20 000 characters): a small markdown (headings, lists, bold, code, code blocks, links) rendered as escaped text; edited in the connection dialog or a notes window; 📝 in the sidebar, the terminal window bar and the context menu; owner only (not share members); kept by export / import and *Duplicate*; lists only carry `has_notes`.
- **Network tools** (`nettools.go`, 🧰 *Tools*, `POST /api/nettools`): port check (up to 100 ports and ranges; open / closed / no answer), ping, traceroute (`traceroute`, `tracepath` or `tracert`), DNS (addresses, CNAME, MX, TXT, PTR) and an HTTP/TLS check (status, server, redirect target, TLS version, certificate subject, issuer, expiry and days left, names, trust) — **from the WRM server or from one of the user's SSH connections through its jump hosts** (port and HTTP checks over SSH channels, the other tools with the server's commands); targets are validated as host names or addresses, one run per user at a time, audit event `nettool.run`, policy `network_tools` (off / admins / all, default all); context menu entries *Network tools from this server…* and *Check ports of this host*.
- **Installable app (PWA)** (`pwa.go`): `/manifest.webmanifest` and a service worker (`/sw.js`) that caches the app's own static files per version and shows an offline page; *Settings → General → Install as an app*.
- Tests: quick connect (targets, temporary connections, expiry with and without open terminals, limit, export), notes (save, limit, list flag, duplicate), network tools (ports and HTTP/TLS from WRM and through SSH, DNS, ping and traceroute on a server, validation, ownership, policy) and the PWA endpoints.

### Fixed
- The top bar overlapped its buttons at medium widths (and the version label never collapsed); it now gives way step by step.

## [10.7.0] — 2026-10-02 — out-of-band management, consoles & serial ports

### Added
- **Out-of-band management** (`bmc.go`): a BMC per connection (`connections.bmc`, JSON with the password encrypted, or a vault credential):
  - **Redfish** (iDRAC, iLO, XClarity, Supermicro, OpenBMC): power state, health, maker, model, serial number, BIOS, CPUs, memory, BMC firmware, power draw, inlet temperature; certificates pinned on first use like SSH host keys (shared with FTPS: `pinnedTLSConfig`);
  - **IPMI** with `ipmitool` (`chassis status`, `mc info`, `fru print`), the password in the environment (`-E`);
  - power actions on, graceful shutdown/restart, force off/restart, power cycle and NMI — only those the BMC offers — with confirmation and audit events `bmc.power` / `bmc.power_failed`;
  - through the connection's jump hosts: Redfish over an SSH channel of the last hop, `ipmitool` on the last hop (password over stdin).
- **Serial consoles** in terminal windows (`console.go`): over the BMC's SSH with the vendor's console command (iDRAC `console com2`, iLO `vsp`, …, suggested from the BMC maker), or IPMI **Serial-over-LAN** on the jump host or on the WRM server in a pseudo terminal; owner only, audited and recorded (`protocol` console / sol); saved in workspace sessions.
- **Serial ports** of the WRM server: protocol `SERIAL` (device name, speed, data bits, parity, stop bits, flow control), raw terminal, *Test connection*; policy `serial_ports` (off / admins / all, default admins), Linux.
- Policies `bmc_enabled` (`BMC_ENABLED`) and `serial_ports`; ipmitool in the admin overview; audit filters for `bmc.` and `inventory.`.
- Tests: Redfish against a fake BMC (status, actions, pinning, changed certificate, wrong password, jump host, credential), IPMI with a fake ipmitool (locally and on a jump host), Serial-over-LAN (PTY and jump host), the BMC SSH console and a serial port on a pseudo terminal.

### Changed
- The terminal WebSocket runs on terminal backends (SSH shell, BMC console, SOL, serial port) instead of an SSH session only.

### Fixed
- Saved workspace sessions did not include remote desktop windows.

## [10.6.0] — 2026-10-02 — tags & inventory import (CSV, Excel, NetBox)

### Added
- **Tags** on connections (`connections.tags`): free tags and tags with a key (`env:`, `site:`, `rack:`, `role:`, `tenant:`, `platform:`, `cluster:`), normalised to lower case, at most 30:
  - a tags field with suggestions in the connection dialog, bulk tagging (`POST /api/connections/bulk` action `tag`, add / remove);
  - a tag filter bar in the sidebar (several tags = all must match), tag search (`#tag`, several words);
  - environment badges (PROD / staging / test / dev colours) and a red edge for production connections;
  - kept by export / import and *Duplicate*.
- **Inventory import** (`inventory.go`, *Settings → Data → Import inventory*):
  - **CSV** with automatic separator (`,` `;` tab `|`, Excel's `sep=` line) and encoding detection (UTF-8, UTF-16, Windows-1250), and **Excel .xlsx** (read with archive/zip + encoding/xml, sheet selection, shared and inline strings);
  - a column mapping guessed from English and Croatian header names, checked and changed in a preview;
  - host fields with ports, users, CIDR suffixes, IPv6 and URLs; Windows platforms become RDP; environment / site / rack / role / platform columns become tags; jump hosts by name;
  - **NetBox** devices and virtual machines through the REST API (filters, all pages, primary IP or DNS name, environment custom field), with **sync**: `connections.ext_id` remembers the NetBox object, a later import updates address, NetBox tags and folder (`ext_tags`) and lists objects that disappeared; address and token can be remembered per user (token encrypted, table `inventory_sources`);
  - for both: the login of the imported connections (none, a **vault credential** or a **stored SSH key**), default user and protocol, folders by column / site / role / tenant / one folder, extra tags.
- Audit events `inventory.imported`, `inventory.netbox_imported`, `connection.tagged`.
- Tests: tags (normalisation, bulk, export/import), CSV encodings and separators, an XLSX import, NetBox import and sync against a fake NetBox with pagination.

### Fixed
- Texts that still said mRemoteNG RDP/VNC/Telnet connections are skipped (they are imported since 10.4).

## [10.5.0] — 2026-10-02 — SSH keys & credentials vault

### Added
- **SSH key store** (🔑 *Keys* in the top bar, `keys.go`):
  - generate ED25519, RSA 2048/3072/4096 or ECDSA 256/384/521 keys, or import OpenSSH/PEM private keys (with passphrase) and public keys; private keys are encrypted at rest;
  - new auth method **SSH key from the key store** (`KEY_REF`, `connections.key_id`); last use is recorded;
  - **deploy** a key to many SSH connections at once like `ssh-copy-id` (through jump hosts, idempotent, keeps other lines, restores SELinux contexts) and optionally switch them to the key after a test login;
  - **revoke** a key from many servers (never the key WRM logs in with);
  - **who has access**: read a connection's `~/.ssh/authorized_keys`, match the keys of the key store, mark WRM's login key, remove single entries;
  - public key copy/download, rename, delete (refused while in use), export of the private key with re-authentication and an optional new passphrase.
- **Credentials vault** (`credentials.go`):
  - a credential holds a user name with a password and/or a stored key; new auth method **Credential from the vault** (`CREDENTIAL`, `connections.credential_id`) for SSH, SFTP, RDP, VNC, Telnet and FTP;
  - sharing with selected users (administrators: everybody) without revealing the secret; host lists (glob patterns, CIDR) limit where it can be used — for grantees also the jump hosts;
  - **password rotation** on all SSH servers that use it: pre-flight login check, `passwd` over a PTY, verification of the new login, all-or-nothing with rollback, a *pending* password if a rollback fails, live progress per server; a check-only mode;
  - rotation reminders (badge, admin overview warning), owner-only *Show* with re-authentication, usage list.
- Audit events `ssh_key.*` and `credential.*`; audit filter entries for both; admin overview card with key, deployment and credential counts.
- Bulk action **🔑 Key** for selected connections; context menu entries *Who has access…* and *Deploy SSH key…*.
- Export/import reference keys and credentials by name (`key_ref`, `credential_ref`) instead of copying them.
- Tests with a fake SSH server that runs scripts in `sh` and simulates `passwd`: key store, deploy/revoke/who has access, vault sharing, host and jump-host restrictions, rotation with success, rollback and incomplete rollback.

### Changed
- Deleting an account also deletes its snippets, tunnel definitions, keys and credentials (and its grants).
- RDP, VNC, Telnet and FTP connections offer only password and vault logins in the connection dialog.

## [10.4.0] — 2026-10-02 — RDP, VNC & Telnet in the browser

### Added
- **Remote desktop** connections — **RDP**, **VNC** and **Telnet** — in WRM windows, through guacd (Apache Guacamole proxy daemon):
  - WRM performs the guacd handshake with the stored credentials and options and relays the Guacamole protocol over a WebSocket (`/ws/desktop`, subprotocol `guacamole`);
  - browser → guacd input is limited to an allow-list of instructions; internal pings are answered by WRM;
  - guacamole-common-js 1.5.0 (Apache 2.0) is bundled and loaded on demand.
- **Desktop windows:**
  - scale-to-fit or 1:1 with scroll bars, full screen, dynamic resolution for RDP (display update);
  - keyboard and mouse (also touch), **Ctrl+Alt+Del**, **send to the remote clipboard**, **type text as keystrokes**;
  - remote clipboard into WRM's clipboard panel, sound from RDP, prompts for credentials the server requires, reconnect;
  - tabs with connection state; restored by saved workspace sessions.
- **Connection options:**
  - RDP: domain, security (NLA/TLS/RDP/Hyper-V), keyboard layout, resize behaviour, colors, start program, certificate, console session, sound, clipboard, wallpaper, **RD Gateway**;
  - VNC: colors, pointer, view only, clipboard;
  - Telnet: font size, colors, login/password prompt patterns.
- **Behind jump hosts:** a temporary tunnel through the jump hosts is opened for guacd for the duration of the session.
- **Audit and recording:**
  - desktop sessions in *Sessions & recordings* (`protocol` rdp/vnc/telnet) and audit events `desktop.open`, `desktop.close`, `desktop.error`;
  - the screen stream is recorded (gzip, SHA-256) and replayed in the browser with play/pause/seek, or downloaded as `.guac`.
- **Admin overview** shows whether guacd answers and its version; *Test connection* checks the port and guacd.
- **Policies** `desktop_enabled` (`REMOTE_DESKTOP_ENABLED`), `guacd_address` (`GUACD_ADDRESS`), `desktop_tunnel_bind`.
- `docker-compose.yml` contains a `guacd` service sharing WRM's network.
- **mRemoteNG import** brings RDP (domain, console, colors, RD Gateway, sound), VNC (view only) and Telnet connections instead of skipping them.
- **Live status** checks RDP/VNC/Telnet ports too.
- **Tests:**
  - Guacamole codec;
  - desktop relay against a fake guacd (handshake parameters, allow-list, ping, wrong password, policy, recording and download);
  - VNC through a jump host.

### Fixed
- SQLite `busy_timeout` is now set on every pooled database connection (it reached only one before), avoiding rare "database is locked" errors under concurrent load.

## [10.3.0] — 2026-10-02 — snippets, broadcast input & live status

Daily terminal work for data center admins (extension plan phases 4–6, first part). Release names no longer carry a suffix: this release is `v10.3.0`.

### Added
- **Snippets** — saved commands:
  - in the right panel (⚡ tab, grouped, searchable) and in a quick picker in every terminal (**Ctrl+Shift+Space** or ⚡ in the title bar);
  - click runs in the focused terminal, Shift+click inserts without Enter;
  - scope: all connections, one folder, or one connection; multi-line commands;
  - variables `{{host}} {{port}} {{user}} {{name}} {{folder}} {{wrm_user}} {{date}} {{time}}` and prompts `{{?Label}}` / `{{?Label=default}}` with a preview;
  - shared snippets for all users (administrators); example set; export/import; copied with *Duplicate*.
- **Run on connect:** snippets typed into the shell by the server once a terminal of a matching connection is connected (also after reconnects), e.g. `sudo -i`, `cd /srv/app`, `tmux attach`. Created from the connection or folder context menu; audited as `terminal.auto_run`.
- **Broadcast input:**
  - type into several terminals at once (📣 in the top bar, per-window toggle, orange frame and status bar);
  - dangerous commands (rm -rf, shutdown, mkfs, dd, iptables -F, DROP TABLE, kubectl delete, …) and multi-line pastes ask first, with Cancel as the default;
  - snippets run on every terminal of the broadcast with their own variables;
  - every start and stop is audited per terminal (`terminal.broadcast`); policy `broadcast_enabled` (`BROADCAST_ENABLED`).
- **Live up/down status:**
  - a background monitor checks every saved connection (TCP connect; SSH and FTP read the greeting and close politely, like Nagios check_ssh) once per host:port and round, retrying before reporting down;
  - status dots in the sidebar with latency, banner and "since"; down counter per folder; "show only down" filter;
  - notifications (also desktop notifications in the background) when a server goes down or comes back;
  - *Check status now* for a connection or a whole folder, also through jump hosts;
  - connections behind jump hosts show their jump host's state, or are checked through it (`status_jump_checks`);
  - per-connection opt-out (*Monitor up/down status*); policies `status_enabled` (`STATUS_ENABLED`), `status_interval_seconds`, `status_jump_checks`; API `GET /api/status`, `POST /api/status/check`; event `status_changed`.
- **Tests** for snippets (API, sharing, validation, variables, run on connect, export/import), broadcast audit and the status monitor (up/down, banners, jump hosts, check now, FTP greeting).

### Changed
- Release tags, binaries and the Docker image no longer have a suffix after the version number (`v10.3.0`, `wrm-pro-v10.3.0-linux-amd64`). CI names the release files after the tag.
- Confirmation dialogs for dangerous actions focus *Cancel*, so Enter cannot confirm them by accident.
- The "terminal connected" marker in the sidebar is blue (green is now the live status).

### Fixed
- Deleting several connections at once also stops their tunnels, removes their tunnel definitions and detaches connections that used them as jump host (single deletes already did).

## [10.2.0] — 2026-10-02 — jump hosts, SSH tunnels & mRemoteNG import

Phase 5 (fleet operations) of the extension plan, first part: what mRemoteNG, PuTTY and `ssh -J` users need to manage servers behind bastions (see [ARCHITECTURE.md](ARCHITECTURE.md)).

### Added
- **Jump hosts** (ProxyJump):
  - any connection can go through another saved SSH connection, in chains of up to 5 hops; loops, foreign and non-SSH jump hosts are refused;
  - used by terminals, the file manager (SFTP and FTP/FTPS), search, editor, server-to-server transfer, *Test connection*, tunnels and web interfaces;
  - every hop authenticates with its own credentials and its host key is verified;
  - the route is shown in the sidebar (⤳), connection test, terminal, session history (`jump_path`) and audit log.
- **SSH tunnels** (port forwarding) per connection:
  - local (`-L`), remote (`-R`) and dynamic SOCKS5 (`-D`) tunnels, through the jump hosts too;
  - listen address and port (0 = automatic), target, *Open as* http/https + path, start mode *manual*, *with terminal* or *always* (also after a restart);
  - templates: web interface (HTTPS/HTTP), SSH, RDP, VNC, PostgreSQL, MySQL/MariaDB, iDRAC/iLO/IPMI, SOCKS;
  - a hint in plain words for every tunnel;
  - each tunnel keeps its own SSH connection with a 30 s keepalive and automatic reconnect (backoff 2–60 s) while its port stays open;
  - statistics: open and total connections, bytes up/down;
  - **🔀 Tunnels** panel (top bar, with a counter of running tunnels): state, listen → target, traffic, *Open*, *Copy address*, *Start/Stop*; administrators can show and stop the tunnels of all users;
  - connection context menu: *Tunnels*, *Start / Stop all tunnels*, *Add tunnel…*;
  - live updates through `/ws/events` (`tunnels_changed`).
- **Web interface connections** (protocols `HTTPS` / `HTTP`, with a path): a double-click opens the page in a new tab, directly or, behind a jump host, through an automatic temporary tunnel (30 minutes idle limit). *Test connection* checks that the port answers.
- **Import from mRemoteNG** (`confCons.xml`):
  - AES-GCM and legacy AES-CBC passwords, master password (checked with `Protected`, asked for when needed), full-file encryption;
  - inheritance of user name, password, port, domain and SSH tunnel;
  - containers become folders; `SSHTunnelConnectionName` becomes the jump host;
  - unsupported protocols and duplicates are reported as skipped.
- **Import from OpenSSH config:** `Host` blocks, `Host *` and wildcard defaults, `HostName`/`User`/`Port`, `ProxyJump` (implicit jump hosts are created), `ProxyCommand ssh -W`, `LocalForward`/`RemoteForward`/`DynamicForward` (tunnels that start with the terminal), `IdentityFile` (when server key files are allowed).
- **Policies:** `tunnels_enabled` (`TUNNELS_ENABLED`), `tunnel_users`, `tunnel_bind_any`, `tunnel_remote_forward`, `tunnel_idle_minutes`. Running tunnels that a policy change no longer allows are stopped.
- **Audit events:** `tunnel.configured`, `tunnel.start`, `tunnel.stop` (with traffic), `tunnel.error`, `tunnel.reconnected`, `web.open`; connection changes record the jump host.
- **Admin panel → Overview:** active SSH tunnels with *Stop*. **Admin status API** includes the tunnels.
- **API:** `GET/PUT /api/connections/{id}/tunnels`, `GET /api/tunnels[?all=1]`, `POST /api/tunnels/{key}/start|stop`, `POST /api/connections/{id}/open-web`, `POST /api/config/import/mremoteng`, `POST /api/config/import/sshconfig`. Connections have `jump_id`, `web_path`, `route` and `tunnels`.
- **Tests:**
  - the in-process SSH test server also handles `direct-tcpip` and `tcpip-forward`;
  - new tests for two-hop jump chains (terminal, SFTP, connection test, validation), local/SOCKS/remote tunnels (traffic, ownership, policies, audit, start modes) and web interfaces;
  - new tests for mRemoteNG import (both encryption formats, master password, inheritance) and OpenSSH config import.

### Changed
- Export and import of WRM's own format include jump hosts and tunnels (IDs are remapped).
- *Duplicate* copies the jump host, path and tunnels of a connection. Deleting a connection stops its tunnels and makes connections that used it as jump host direct.
- Editing a connection restarts its running tunnels with the new settings.
- The connection dialog is wider and fully translated (labels were partly English in Croatian).

### Fixed
- The connection dialog no longer moves the focus back to *Name* when you already started typing in another field.

## [10.1.0] — 2026-10-01 — audit trail & session recording

First phase of the extension plan: WRM as a lightweight PAM for teams (see [ARCHITECTURE.md](ARCHITECTURE.md)).

### Added
- **Session recording** of every SSH terminal in asciinema format (asciicast v2, gzip, SHA-256 in the database). Recording is written by a background goroutine through a bounded queue, so the terminal is never slowed down. There is a size limit per session (`recording_max_mb`) and a retention period (`recording_retention_days`).
- **Replay in the browser** (play/pause, seek, 0.5–16×, skip idle time) and `.cast` download. Viewing and downloading are audited.
- **Optional keystroke recording** (`session_recording_input`, off by default). Typing at password, passphrase and PIN prompts is masked.
- **● REC** badge and a *“This session is recorded”* notice in the terminal.
- **Terminal sessions log** (`terminal_sessions`): user, IP, server, remote user, start/end, duration, status (closed, failed, connection lost, ended by an administrator, interrupted by a restart), exit code, bytes. Sessions left open by a crash are closed at the next start, and their partial recordings are kept.
- **File transfer log** (`file_transfers`): every upload, download (also editor opens and each file of a ZIP download) and server-to-server copy, with source, destination, size, SHA-256 and status.
- **Admin panel:** *Sessions & recordings* and *File transfers* tabs with filters and CSV export. The *Audit log* tab gained date range filters, recording events, links to sessions and **Verify integrity**.
- **Settings → Session history:** your own sessions and sessions on your connections.
- **Audit integrity:**
  - audit entries now reference the connection and the terminal session;
  - every entry carries a SHA-256 hash chain, checked by `GET /api/admin/audit/verify`;
  - database triggers make `audit_log`, `file_transfers` and `session_recordings` append-only (no UPDATE; no DELETE of entries younger than 7 days);
  - secret-looking keys in audit details are redacted.
- **Feature flags:** `audit_enabled` (`AUDIT_ENABLED`), `session_recording` (`SESSION_RECORDING_ENABLED`), `session_recording_input`, `recording_retention_days`, `recording_max_mb`. Environment aliases are supported for the plan's flag names.
- **API:** `GET /api/recordings`, `/api/recordings/{id}`, `/api/recordings/{id}/cast`, `GET /api/admin/transfers`, `GET /api/admin/audit/verify`, `GET /healthz`.
- **Operations:**
  - `Dockerfile` (unprivileged user, `/data` volume, health check) and `docker-compose.yml`;
  - CI builds and smoke-tests the image and publishes release images to `ghcr.io/vedranius/wrm-pro`;
  - `wrm -healthcheck`;
  - `WRM_RECORDINGS_DIR`.
- **Brand:** new WRM PRO logo for the favicon, app icons, sign-in screen, top bar, sidebar, README, release notes and GitHub social preview. A *WRM Orange* accent color. The logo kit is in `docs/brand/`.
- **Documentation:**
  - `ARCHITECTURE.md` (stack, auth, SSH layer, credentials, collaboration, migrations, extension plan status and touch points);
  - this changelog;
  - README, SECURITY and CONTRIBUTING updated.
- **Tests:**
  - an integration test with an in-process SSH + SFTP server (connect → audit entries → replayable recording without passwords; transfer checksums; access to recordings);
  - unit tests for the recorder (UTF-8 split across chunks, masking, resize, size limit) and for the append-only, hash-chained audit log.

### Changed
- Downloads by connection owners are now audited too (before, only downloads through shares were).
- Server-to-server transfer reports a read error as a failed file instead of a partial "ok".
- The audit log CSV export has new columns: connection, session and hash.

## [10.0.1] — 2026-09-28

### Fixed
- "404 page not found" at `/static/`, the v9 address of the app. It now redirects to `/`.
- Upgrading a database from an earlier build with an incompatible `audit_log` / `collab_messages` table printed `Index warning: no such column: ts / share_id`, and the audit log and chat history could not be written. Such tables are now kept as `<table>_old_<time>` and created again.

### Added
- The startup log shows the address to open in the browser.

## [10.0.0] — 2026-09-28

### Added
- **Voice calls** in collaboration rooms (WebRTC):
  - mute, deafen, push-to-talk, device selection, noise suppression, echo cancellation and gain control;
  - speaking indicators, per-person volume, connection quality;
  - a built-in TURN relay (UDP/TCP 3478) with short-lived credentials.
- **Users and roles:**
  - admin user management (create, disable, reset password/2FA, sign out everywhere, unlock, delete);
  - share access modes (members only / signed-in users / anyone with the link);
  - five roles enforced by the server (observer, viewer, operator, moderator, owner);
  - members, people history, bans, expiry, link rotation, immediate revocation.
- **Collaboration:** collaboration bar, people panel with moderation, chat history, file exchange, raise hand, several shared terminals with screen snapshot, remote keyboard control.
- **Security:**
  - TOTP 2FA with recovery codes and a 2FA policy; closed self-registration; account lockout; hashed session tokens with idle/maximum lifetime; device list;
  - a random AES-256-GCM key per installation (migrated from the v9 default key); secrets never sent to the browser; secret export only after re-entering the password;
  - SSH host key verification (trust on first use, or strict) and FTPS certificate pinning;
  - CSRF protection, CSP and security headers, request limits, local assets (no CDN), self-signed HTTPS option, trusted proxy support;
  - an audit log with CSV export.
- **Admin panel:** overview with warnings, users, shares, policies, voice & network, host keys, audit log. Policies can be forced with `WRM_<KEY>` environment variables.
- `wrm -version`, `wrm -reset-password USER [-reset-2fa]`; unit tests in CI with the race detector; `SECURITY.md`.

## [9.10.2]

### Security
- In a collaboration session any participant could give themselves keyboard control of another participant's shared terminal. Now only the sharer can grant control, keystrokes go only to that sharer, and control ends when sharing stops.

### Fixed
- The *Control* button appears next to participants while you share a terminal; late joiners see *Watch* at once; duplicate names get a suffix.

## [9.10.1]

### Added
- Saved sessions restore all windows (connection, mode, folder, position, snap, minimized, focus).
- Reconnect button on SSH windows (Ctrl+Shift+R).
- Top/bottom half snap targets, drag-to-edge snapping, snap layouts menu.

### Fixed
- Windows stay inside the screen when the zoom, resolution or panel sizes change.

## [9.10.0]

### Added
- Terminal: reconnect and automatic reconnect, connection state, flow control, SSH keepalive, search in output, clickable URLs, configurable scrollback.
- File manager: instant filter and recursive name/content search, breadcrumb path, multi-select, keyboard shortcuts, folder and drag & drop uploads with progress, built-in text editor, pooled SFTP connections.
- Connections: test connection, duplicate, FTPS. Modern dark responsive UI with an accent color; mobile drawers.
- Security: login rate limiting, WebSocket origin check, optional HTTPS.

### Fixed
- Disconnects on large outputs with multi-byte characters (binary frames).
- Full-screen programs drawn at the wrong size (the PTY size follows the window).
- Folder transfer when WRM runs on Windows; several smaller issues.

[10.7.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.7.0
[10.6.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.6.0
[10.5.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.5.0
[10.4.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.4.0
[10.3.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.3.0
[10.2.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.2.0&expanded=true
[10.1.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.1.0&expanded=true
[10.0.1]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.0.1&expanded=true
[10.0.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.0.0&expanded=true
[9.10.2]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v9.10.2&expanded=true
