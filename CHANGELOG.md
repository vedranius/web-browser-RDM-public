# Changelog

All notable changes to Web Remote Manager PRO. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Release notes with downloads are on the
[Releases page](https://github.com/vedranius/web-browser-RDM-public/releases).

## [10.1.0-mimo] — 2026-10-01 — audit trail & session recording

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

## [10.0.1-mimo] — 2026-09-28

### Fixed
- "404 page not found" at `/static/`, the v9 address of the app. It now redirects to `/`.
- Upgrading a database from an earlier build with an incompatible `audit_log` / `collab_messages` table printed `Index warning: no such column: ts / share_id`, and the audit log and chat history could not be written. Such tables are now kept as `<table>_old_<time>` and created again.

### Added
- The startup log shows the address to open in the browser.

## [10.0.0-mimo] — 2026-09-28

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

## [9.10.2-mimo]

### Security
- In a collaboration session any participant could give themselves keyboard control of another participant's shared terminal. Now only the sharer can grant control, keystrokes go only to that sharer, and control ends when sharing stops.

### Fixed
- The *Control* button appears next to participants while you share a terminal; late joiners see *Watch* at once; duplicate names get a suffix.

## [9.10.1-mimo]

### Added
- Saved sessions restore all windows (connection, mode, folder, position, snap, minimized, focus).
- Reconnect button on SSH windows (Ctrl+Shift+R).
- Top/bottom half snap targets, drag-to-edge snapping, snap layouts menu.

### Fixed
- Windows stay inside the screen when the zoom, resolution or panel sizes change.

## [9.10.0-mimo]

### Added
- Terminal: reconnect and automatic reconnect, connection state, flow control, SSH keepalive, search in output, clickable URLs, configurable scrollback.
- File manager: instant filter and recursive name/content search, breadcrumb path, multi-select, keyboard shortcuts, folder and drag & drop uploads with progress, built-in text editor, pooled SFTP connections.
- Connections: test connection, duplicate, FTPS. Modern dark responsive UI with an accent color; mobile drawers.
- Security: login rate limiting, WebSocket origin check, optional HTTPS.

### Fixed
- Disconnects on large outputs with multi-byte characters (binary frames).
- Full-screen programs drawn at the wrong size (the PTY size follows the window).
- Folder transfer when WRM runs on Windows; several smaller issues.

[10.1.0-mimo]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.1.0-mimo
[10.0.1-mimo]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.0.1-mimo
[10.0.0-mimo]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.0.0-mimo
[9.10.2-mimo]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v9.10.2-mimo
