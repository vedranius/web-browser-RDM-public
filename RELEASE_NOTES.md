<p align="center"><img src="https://raw.githubusercontent.com/vedranius/web-browser-RDM-public/v10.2.0-mimo/docs/brand/png/lockup/wrm-lockup-on-dark-664w.png" alt="WRM PRO — Web Remote Manager" width="332"></p>

## Web Remote Manager PRO v10.2.0-mimo — jump hosts, SSH tunnels & mRemoteNG import

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/vedranius)

Browser-based remote server management: SSH terminal, SFTP/FTP/FTPS file manager, server-to-server transfer, sharing, real-time collaboration, **voice calls**, audit trail and session recording. One self-contained binary (or container) with an embedded web UI.

v10.2 brings what admins in data centers use every day in mRemoteNG, PuTTY and `ssh -J`: **servers behind bastions**, **SSH tunnels** to internal web interfaces and services, and a one-click **import of mRemoteNG and OpenSSH configurations**.

### 🪜 Jump hosts (bastions)

- Every connection can go **through another saved SSH connection**: choose *Jump host (connect via)* in the connection. Like `ssh -J` / `ProxyJump` and the *SSH tunnel* setting of mRemoteNG.
- **Chains** up to 5 hops (`vpn-gw → bastion-dc1 → app-01`); loops are refused.
- **Everything works through it:** terminal, file manager, FTP/FTPS, search, editor, server-to-server transfer, *Test connection*, tunnels and web interfaces.
- Every hop uses **its own credentials** and its **host key is verified**. No port is opened on the jump host (`direct-tcpip`, the bastion only needs `AllowTcpForwarding yes`).
- The route is shown everywhere: **⤳** in the sidebar, *Test connection* (`… via bastion-dc1`), the terminal, session history and the audit log.

### 🔀 SSH tunnels (port forwarding)

- **Local (-L):** a port on the WRM machine → a host:port behind the server, e.g. `127.0.0.1:2443 → 10.0.0.5:443` for the web UI of a switch, iDRAC or iLO.
- **SOCKS proxy (-D):** one proxy for your browser or tools that reaches the whole network behind a server.
- **Remote (-R):** a port on the server → a host:port on the WRM side (administrators, policy).
- Configure them in the connection (**🔀 SSH tunnels**, with **templates**: web HTTPS/HTTP, SSH, RDP, VNC, PostgreSQL, MySQL, iDRAC/iLO/IPMI, SOCKS) and use the new **🔀 Tunnels** panel: state, listen → target, connections, traffic, **🌐 Open**, **📋 Copy address**, **▶ Start / ■ Stop**.
- **Start modes:** manual, **with the terminal** (starts with the first terminal and stops after the last one), **always** (starts with WRM).
- **Reliable:** each tunnel keeps its own SSH connection through the jump hosts, checks it every 30 s and **reconnects automatically** while its port stays open.
- **Safe defaults:** tunnels listen on `127.0.0.1`; network addresses and remote forwards are for administrators; a tunnel can be turned off or limited to administrators; every start, stop and error is **audited** with the traffic. Administrators see and stop all tunnels in *Admin panel → Overview*.

### 🌐 Web interface connections

- New protocols **Web interface (HTTPS / HTTP)** with host:port and path: save the web UIs of your iDRACs, iLOs, switches, firewalls and NAS next to your servers.
- **Double-click** opens it in a new tab — directly, or, **behind a jump host, through an automatic temporary tunnel** (closed after 30 minutes without traffic).

### 📥 Import from mRemoteNG and OpenSSH

- **mRemoteNG `confCons.xml`:** folders, SSH and web connections **with their passwords** (both encryption formats, master password, full-file encryption), inherited settings, and **SSH tunnel → jump host**. RDP/VNC/Telnet are listed as skipped; duplicates are skipped.
- **OpenSSH `~/.ssh/config`:** hosts, users, ports, `Host *` defaults, **`ProxyJump`** / `ProxyCommand ssh -W`, **`LocalForward` / `RemoteForward` / `DynamicForward`** and `IdentityFile`.
- *Settings → Data* shows a summary of what was imported, skipped and what to check. Export/import of WRM's own format now includes jump hosts and tunnels.

### ⚙️ New policies

`tunnels_enabled` (`TUNNELS_ENABLED`, on), `tunnel_users` (`all` / `admins`), `tunnel_bind_any` (off), `tunnel_remote_forward` (`off` / `admins` / `all`, default `admins`), `tunnel_idle_minutes` (0). All in *Admin panel → Security policies → SSH tunnels*, or forced with `WRM_<KEY>`.

### ⬆️ Upgrading from v10.1

Replace the binary. The database gets a new table (`connection_tunnels`) and two new columns (`connections.jump_conn_id`, `connections.web_path`); nothing else changes, and the previous binary still starts on the upgraded database.

- **Tunnels are available to all users by default**, on `127.0.0.1` of the WRM machine only. If WRM runs on a shared server, consider `tunnel_users=admins` (local users of that server can reach `127.0.0.1` ports), or `WRM_TUNNELS_ENABLED=0`.
- Where is a tunnel port? On the computer **where WRM runs**. That is usually your own PC. If WRM runs on a server, see *SSH tunnels → Where is the tunnel port?* in the README.

Coming from v10.0 or v9? Read the upgrade notes below as well. See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/CHANGELOG.md) for the full history.

---

## Included from v10.1.0 — audit trail & session recording

#### 🎬 Session recording & replay

- **Every SSH terminal session is recorded** in the open **asciinema** format (asciicast v2), compressed, with its SHA-256 checksum stored in the database.
- **▶ Replay in the browser** (play/pause, seek, 0.5–16× speed, *skip idle time*), or **⬇ download the `.cast`** and play it with `asciinema play`.
- **Transparent for users:** the terminal shows *“This session is recorded”* and a **● REC** badge.
- **No passwords in recordings:** servers do not echo passwords, so they never appear in the output. Keystroke recording is **off** by default. When an administrator turns it on, typing at password, passphrase and PIN prompts is masked.
- **Never slows the terminal:** the recording is written in the background. A per-session size limit stops the recording, not the terminal.
- **Who can watch:** administrators all; users their own sessions and sessions on **their** connections (e.g. guests of a share). Watching or downloading a recording is itself audited.
- **Where to find it:** *Admin panel → Sessions & recordings* (everything), *Settings → Session history* (your own).

#### 🧾 Audit trail

- **Terminal sessions:** who, from which IP, which server and remote user, start, duration, and how the session ended (closed, connection lost, ended by an administrator, interrupted by a restart).
- **File transfers:** every upload, download (also files opened in the editor and every file of a ZIP download) and server-to-server copy, with source, destination, size and **SHA-256**. See *Admin panel → File transfers*, with CSV export.
- **Append-only, tamper-evident audit log.** Database triggers refuse to change audit entries, transfers and recordings, and refuse to delete anything younger than 7 days. Only the retention job removes old data. Every entry contains the hash of the previous one, and **🔏 Verify integrity** checks the whole chain.
- **Better audit search:** date range, events linked to the connection and the terminal session (click **▶ #id** to replay), recording events, and secrets in event details are always redacted.
- **Feature flags:** `audit_enabled` (`AUDIT_ENABLED`) and `session_recording` (`SESSION_RECORDING_ENABLED`) are on by default. Also new: `session_recording_input`, `recording_retention_days` (90) and `recording_max_mb` (100). Administrator actions are always recorded.

#### 🐳 Operations

- **Docker:** `Dockerfile` (static binary, unprivileged user, data volume `/data`, health check) and **`docker compose up -d`**. Release images are published to `ghcr.io/vedranius/wrm-pro`.
- **`/healthz`** for load balancers and monitoring; **`wrm -healthcheck`** for container health checks.
- Recordings location: `WRM_RECORDINGS_DIR` (default `recordings/` next to the database).

#### 🎨 New logo

The WRM PRO logo is now used for the favicon, the app icons, the sign-in screen, the top bar, the README and the release. A **WRM Orange** accent color is available in *Settings → General*. The complete logo kit (SVG and PNG) is in `docs/brand/`.

#### 📐 Architecture & tests

- [`ARCHITECTURE.md`](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/ARCHITECTURE.md) describes how WRM is built. It also maps the extension plan (RBAC, SSO, broadcast input, fleet status, …) to the code.
- A new **integration test** runs a real SSH + SFTP server inside the test. It checks connect → audit entries → recording (replayable, no passwords), transfer checksums and the append-only audit trail, and it runs in CI with the race detector.

#### ⬆️ Upgrading from v10.0

Replace the binary. The database is only **extended**: new tables, columns and triggers are added, and nothing is changed or removed. The previous binary still runs on the upgraded database.

- **Session recording is on by default.** Tell your users about it (in many countries this is required), or turn it off in *Admin panel → Security policies → Session recording*.
- Recordings use disk space in `recordings/` next to the database. A session usually takes a few KB to a few MB. Recordings are kept for 90 days, with at most 100 MB per session.
- Audit entries written before the upgrade have no hash. *Verify integrity* checks the chain from the first new entry on.

---

## Included from v10.0.x

### 🔧 v10.0.1-mimo

- `/static/` (the v9 address) forwards to the app instead of returning 404.
- A database from an earlier build with an incompatible `audit_log` / `collab_messages` table is upgraded at startup. The old table is kept as `<table>_old_<time>`.
- The server log shows the address to open in the browser.

### v10.0.0-mimo

#### 🎙 Voice calls in collaboration rooms (back, rebuilt)

Everyone in a share room can talk. It works like Discord or Jitsi:

- **🎙 Join voice** in the collaboration bar; **mute** (Ctrl+Shift+M), **deafen** (Ctrl+Shift+D), **leave**.
- **Push-to-talk** with a key of your choice (e.g. F8), or voice activity.
- **Microphone and speaker selection**, live input level meter, test sound, **noise suppression**, **echo cancellation**, **automatic gain control**.
- **Speaking indicators** (green ring), **per-person volume**, join/leave sounds, connection quality per person (P2P or relay, round-trip time).
- **Listen-only**: without a microphone, or without permission to use it, you can still join and listen.
- **Stable on real networks**: WebRTC with "perfect negotiation", automatic ICE restarts when the network changes, and automatic rejoin after reconnects.
- **Built-in TURN relay** (UDP and TCP, port 3478): calls work behind NAT and strict firewalls with no extra software. It accepts only short-lived credentials issued to room participants. By default it refuses to relay to private or loopback networks.
- **End-to-end encrypted** audio (DTLS-SRTP). The relay cannot decrypt it.
- Moderators can **mute someone for everyone**. Admins can turn voice on or off, limit participants per call (default 12), and add external STUN/TURN servers.

> Browsers allow the microphone only on **HTTPS** pages (or `localhost`). The quickest way to get HTTPS is `HTTPS_SELF_SIGNED=1`. For calls across the internet, open **UDP+TCP 3478** and the **UDP relay range** (default `49152-65535`) in the firewall. Behind NAT, set `turn_public_ip`.

#### 👥 Users, members & roles

- **Admin panel → Users**: create users (a temporary password is generated if you leave it empty), set the display name, make or remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete. The last active administrator cannot be removed.
- **Shares have three access modes**: 👥 *Members only* (default), 🏢 *Everyone signed in*, 🌐 *Anyone with the link* (guests; admins can disable this mode).
- **Members** with an individual role each, plus a role for everyone else. There are **five roles**:

  | | Observer | Viewer | Operator | Moderator | Owner |
  |---|:-:|:-:|:-:|:-:|:-:|
  | Chat & voice | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Watch shared terminals | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Browse & download files | — | ✔ | ✔ | ✔ | ✔ |
  | Open terminals, change files, transfer | — | — | ✔ | ✔ | ✔ |
  | Share own terminal, receive keyboard control | — | — | ✔ | ✔ | ✔ |
  | Change roles, mute, remove, ban people | — | — | — | ✔ | ✔ |
  | Share settings & members | — | — | — | — | ✔ |

  The **server enforces roles** on every file operation, terminal and room action.
- **People** view per share: everyone who ever opened it, with last seen and IP. From there you can change a person's role or **ban** them.
- **Share expiry** (1 hour to 30 days, or a custom date), **pause/resume**, **new link** (the old link stops working immediately).
- **Revocation takes effect at once.** Deleting, pausing or rotating a share, removing a member, lowering a role, banning, or disabling a user closes their open terminals and room connections.

#### 💬 Collaboration

- A **new collaboration bar** shows the share, your role, avatars of the people online, and the voice controls.
- **👥 People panel**: who is in the call and who is online, with roles, mute and deafen state, shared terminals, raised hands, connection quality and volume. Moderators have a **⋯** menu for *change role*, *mute*, *stop sharing*, *lower hand*, *remove* and *ban*.
- **Chat history** is kept, so people who join later see it. The chat has timestamps, clickable links, an unread badge and notifications. **File exchange** up to the policy limit: images get a preview, and other files are always downloaded, never opened in the page.
- **✋ Raise hand.**
- **Terminal sharing**: you can share several terminals at once. Viewers get a **snapshot of the current screen** (with colors) and then the live output. Terminal data is sent only to people who watch it.
- **Remote keyboard control**: *Request control* → *Allow / Deny*. Keystrokes go only to the granted terminal, and only while it is shared.
- **Identities come from the server.** Signed-in users appear under their account name. Guests pick a name and are marked *guest*.
- The room **reconnects automatically** and restores watching, sharing and the voice call.
- **Per-participant send queues**: one slow client can no longer stall a room.

#### 🔐 Security hardening

- **Two-factor authentication (TOTP)** with QR code and single-use **recovery codes**. The policy can **require 2FA** for admins or for everyone.
- **Self-registration is closed by default.** Only the first account (the administrator) registers itself.
- **Accounts lock** after repeated wrong passwords (policy). Per-IP rate limiting. Unknown users get constant-time responses. After an admin reset, the user **must change the password**.
- **Sessions**: 256-bit tokens, **stored hashed**. Idle and maximum lifetime (policy). A list of **signed-in devices** with remote sign-out.
- **Stored secrets** (connection passwords, private keys, 2FA secrets) are encrypted with **AES-256-GCM** and a **random key per installation** (`remote_manager.db.key`, or your own `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`).
- **Secrets are never sent to the browser**, not even to their owner. The *Edit connection* dialog shows "saved and encrypted — leave empty to keep it". *Export* leaves secrets out. *Export with passwords & keys* asks for your password again.
- **SSH host key verification** (trust on first use, or strict). A changed key is refused and shown with its fingerprint. **FTPS certificates** are validated, or pinned on first use.
- **Server-side key files** (*Key file* / *Auto (~/.ssh)*) are limited to administrators by default.
- **CSRF protection** on every state-changing API call. **Content-Security-Policy**, `X-Frame-Options`, `nosniff`, `Referrer-Policy: no-referrer`, `Permissions-Policy`, and HSTS on HTTPS. Request-size limits. No directory listings.
- **All assets are served locally** (xterm.js and fonts are embedded). There is no CDN dependency at runtime.
- **Share passwords** are held in signed cookies, which become invalid when the password changes. Guests get their own identity cookie. Share connection lists stay hidden until the share password is entered.
- **HTTPS made easy**: `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE`, or `HTTPS_SELF_SIGNED=1`. Behind a reverse proxy, `WRM_TRUST_PROXY=1` gives correct client IPs and secure cookies.
- **Audit log** of sign-ins (failed ones too), account and policy changes, connections, shares, room joins and moderation, terminals, file changes, host keys, exports and imports. It is searchable in the admin panel, exportable as **CSV**, and written as `AUDIT …` lines to the server log (journald/syslog → SIEM).
- **Recovery from the command line**: `wrm -reset-password USER [-reset-2fa]`.

#### 🛡 Admin panel

*Settings → 🛡 Admin panel* has these tabs:
- **Overview**: version, uptime, users, HTTPS, relay status, key source, **security warnings**, **open terminals** (with *End*), live rooms.
- **Users**
- **Shares**: every share on the server.
- **Security policies**
- **Voice & network**
- **Host keys**
- **Audit log**

Any policy can also be **forced by an environment variable** (`WRM_<KEY>`, e.g. `WRM_REQUIRE_2FA=all`). A forced policy is shown as locked in the panel.

#### ✨ Other improvements

- **Settings** are reorganised into tabs: General, Terminal, Security, Voice & audio, Data.
- **Connections**: SSH keyboard-interactive authentication (servers that ask for the password that way). *Duplicate* copies the stored secrets on the server, so they never pass through the browser. Deleting a connection closes its open terminals.
- **Uploads** respect the server-side size limit (`max_upload_mb`), and partial files are removed.
- Translations (English / Hrvatski) cover the whole new UI.
- A **unit test suite** (TOTP, recovery codes, roles, share access, CSRF, encryption and migration, signed values, settings, input validation) runs in CI with the race detector.

---

#### ⬆️ Upgrading from v9

Replace the binary and start it with the same database, and the same `ENCRYPTION_KEY` if you set one. Everything is migrated automatically.

- **Stored secrets are re-encrypted.** If `ENCRYPTION_KEY` was not set, v9 used a public default key. v10 generates a random key file, `remote_manager.db.key`, next to the database and re-encrypts all secrets with it. **Back that file up separately from the database.** If you did set `ENCRYPTION_KEY`, nothing changes.
- **Self-registration is now closed.** Existing accounts keep working. To reopen it: *Admin panel → Security policies*.
- **SSH host keys are now verified.** The first connection to each server remembers its key.
- **Key file / Auto (~/.ssh)** authentication is now admin-only, unless you turn on the policy *Allow server key files for all users*.
- **Existing shares keep their behaviour**: anyone with the link can open them, with the role *Operator*. Edit a share to switch it to members only or to change roles.
- **Behind a reverse proxy**, state-changing API calls now require the `Origin` to match the `Host`. Pass the `Host` header through (`proxy_set_header Host $host;`) or set `WRM_ALLOWED_ORIGINS`. Also set `WRM_TRUST_PROXY=1`.
- For voice calls, open **UDP+TCP 3478** and the relay port range in the firewall (see above).

---

---

### 📦 Downloads

| Platform | Architecture | Binary |
|---|---|---|
| Linux | x86-64 | `wrm-pro-v10.2.0-mimo-linux-amd64` |
| Linux | arm64 | `wrm-pro-v10.2.0-mimo-linux-arm64` |
| Linux | ARMv7 (Raspberry Pi) | `wrm-pro-v10.2.0-mimo-linux-armv7` |
| Linux | ARMv6 | `wrm-pro-v10.2.0-mimo-linux-armv6` |
| Linux | 32-bit | `wrm-pro-v10.2.0-mimo-linux-386` |
| Windows | x86-64 | `wrm-pro-v10.2.0-mimo-windows-amd64.exe` |
| Windows | arm64 | `wrm-pro-v10.2.0-mimo-windows-arm64.exe` |
| macOS | Intel | `wrm-pro-v10.2.0-mimo-darwin-amd64` |
| macOS | Apple Silicon | `wrm-pro-v10.2.0-mimo-darwin-arm64` |
| macOS | Universal | `wrm-pro-v10.2.0-mimo-darwin-universal` |
| Android | arm64 (Termux) | `wrm-pro-v10.2.0-mimo-android-arm64` |
| FreeBSD | x86-64 | `wrm-pro-v10.2.0-mimo-freebsd-amd64` |
| FreeBSD | arm64 | `wrm-pro-v10.2.0-mimo-freebsd-arm64` |
| OpenBSD | x86-64 | `wrm-pro-v10.2.0-mimo-openbsd-amd64` |

Verify integrity with `SHA256SUMS.txt`. The Android build has no built-in TURN relay; configure an external TURN server there if you need one.

### 🚀 Quick start

**Linux / macOS**
```bash
chmod +x wrm-pro-v10.2.0-mimo-linux-amd64
HTTPS_SELF_SIGNED=1 ./wrm-pro-v10.2.0-mimo-linux-amd64
# open https://<server>:8080 — create the administrator account (the first account)
```
On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

**Windows**: double-click `wrm-pro-v10.2.0-mimo-windows-amd64.exe`, or in PowerShell:
```powershell
$env:HTTPS_SELF_SIGNED=1; .\wrm-pro-v10.2.0-mimo-windows-amd64.exe
```

**Android (Termux)**
```bash
pkg install wget
wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v10.2.0-mimo/wrm-pro-v10.2.0-mimo-android-arm64
chmod +x wrm-pro-v10.2.0-mimo-android-arm64 && ./wrm-pro-v10.2.0-mimo-android-arm64
```

**Docker**
```bash
docker run -d --name wrm -p 8080:8080 -v wrm-data:/data -e HTTPS_SELF_SIGNED=1 ghcr.io/vedranius/wrm-pro:v10.2.0-mimo
# or, from the source tree:  docker compose up -d
```

On the first start WRM creates `remote_manager.db` (the database) and `remote_manager.db.key` (the encryption key) next to each other, both with mode `0600`. **Back up the key file.**

### ⚙️ Most important settings

```
PORT=8080 / LISTEN_ADDR=:8080     listening address
DB_PATH=./remote_manager.db       database
ENCRYPTION_KEY / ENCRYPTION_KEY_FILE   your own key for stored secrets (else <DB_PATH>.key)
HTTPS_CERT_FILE + HTTPS_KEY_FILE  HTTPS with your certificate
HTTPS_SELF_SIGNED=1               HTTPS with a generated self-signed certificate
WRM_TRUST_PROXY=1                 behind a reverse proxy (real client IPs, secure cookies)
WRM_ALLOWED_ORIGINS=host          extra allowed Origin hosts
WRM_<POLICY>=value                force a policy, e.g. WRM_REQUIRE_2FA=all, WRM_TURN_PUBLIC_IP=203.0.113.10
AUDIT_ENABLED / SESSION_RECORDING_ENABLED   audit log / session recording on (1) or off (0)
TUNNELS_ENABLED=0                 turn SSH tunnels off (WRM_TUNNEL_USERS=admins: administrators only)
WRM_RECORDINGS_DIR=/path          where session recordings are stored (default: recordings/ next to the database)
```

For the full documentation (Docker, reverse proxy, systemd, firewall, API), see the [README](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/README.md). To report a vulnerability, see [SECURITY.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/SECURITY.md).

---

### 📜 License

Web Remote Manager PRO is **source-available** under the [PolyForm Noncommercial License 1.0.0 or the PolyForm Internal Use License 1.0.0](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/LICENSE). It is free for personal, educational, non-profit and other noncommercial use, and **free for companies that use it as a work tool**, including paid work for their customers. **Offering WRM as a hosted service, charging for its use, reselling or bundling it requires a commercial license**; see [COMMERCIAL-LICENSE.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/COMMERCIAL-LICENSE.md). Contributions are welcome; see [CONTRIBUTING.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.2.0-mimo/CONTRIBUTING.md).

☕ **Like WRM?** Support its development on **[Ko-fi](https://ko-fi.com/vedranius)**. Thank you!
