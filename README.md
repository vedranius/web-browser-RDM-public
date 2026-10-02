<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/brand/svg/wrm-lockup-on-dark.svg">
    <img src="docs/brand/svg/wrm-lockup-on-light.svg" alt="WRM PRO — Web Remote Manager" width="420">
  </picture>
</p>

# Web Remote Manager PRO (WRM)

**A remote server manager that runs in any web browser.** SSH terminal with **snippets**, **broadcast input** and **live up/down status**, SFTP / FTP / FTPS file manager, **jump hosts** (bastions, chains), **SSH tunnels** (local, remote, SOCKS), **web interfaces** behind jump hosts, **import from mRemoteNG** and `~/.ssh/config`, server-to-server transfers, saved workspaces, sharing with roles, real-time collaboration with **voice calls**, and enterprise security (2FA, policies, a tamper-evident **audit log**, **session recording** with replay, file transfer log): one self-hosted binary (or container) for your PC, server or company.

**Current version: v10.3.0** · [Download](https://github.com/vedranius/web-browser-RDM-public/releases/latest) · [Release notes](RELEASE_NOTES.md) · [Changelog](CHANGELOG.md) · [Security](SECURITY.md) · [Architecture](ARCHITECTURE.md)

---

## 📜 License: free for personal use and as a work tool in companies

WRM is **source-available**. You may use it under the **[PolyForm Noncommercial License 1.0.0](LICENSE)** or the **[PolyForm Internal Use License 1.0.0](LICENSE)**, whichever fits.

| | |
|---|---|
| ✅ **Free for personal & noncommercial use** | Personal use, home labs, learning, hobby projects, non-profits, schools and universities, public research, government. You may **use, fork, modify and share** it, including modified versions. Keep the [LICENSE](LICENSE) and its `Required Notice:` lines. |
| ✅ **Free for companies as a work tool** | Companies of any size may run WRM on their own servers and use it for their work, **including paid work for their customers** (e.g. managing customers' servers), and may modify it for internal use. |
| 💼 **Needs a commercial license** | Offering WRM to others as a hosted service / SaaS, **charging anyone for access to or use of WRM**, reselling, renting or white-labelling it, bundling it in a product or appliance, or using its code in a product you offer to others. Contact the owner first. See **[COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)**. |
| 🤝 **Contributions welcome** | Ideas, bug reports, translations and pull requests. See **[CONTRIBUTING.md](CONTRIBUTING.md)**. |

Copyright © 2026 vedranius. Third-party libraries keep their own licenses.

## ☕ Support the project

WRM is built in my spare time. If it saves you time, you can buy me a coffee:

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/vedranius)

**https://ko-fi.com/vedranius**. Thank you! ❤️

---

## Table of contents

1. [What WRM is](#what-wrm-is)
2. [Quick start](#quick-start)
3. [Upgrading from v10.2](#upgrading-from-v102) · [from v10.1](#upgrading-from-v101) · [from v10.0](#upgrading-from-v100) · [from v9](#upgrading-from-v9)
4. [How it works](#how-it-works)
5. [Features in detail](#features-in-detail)
   - [Accounts, sign-in & two-factor authentication](#accounts-sign-in--two-factor-authentication)
   - [Connections & folders](#connections--folders)
   - [Jump hosts (bastions)](#jump-hosts-bastions)
   - [SSH tunnels (port forwarding)](#ssh-tunnels-port-forwarding)
   - [Web interfaces (HTTP / HTTPS connections)](#web-interfaces-http--https-connections)
   - [Import from mRemoteNG and OpenSSH](#import-from-mremoteng-and-openssh)
   - [Windows, tabs & snapping](#windows-tabs--snapping)
   - [SSH terminal](#ssh-terminal)
   - [Snippets & run on connect](#snippets--run-on-connect)
   - [Broadcast input](#broadcast-input)
   - [Live up/down status](#live-updown-status)
   - [File manager (SFTP / FTP / FTPS)](#file-manager-sftp--ftp--ftps)
   - [Search in files](#search-in-files)
   - [Text editor](#text-editor)
   - [Server-to-server transfer](#server-to-server-transfer)
   - [Workspace sessions](#workspace-sessions)
   - [Clipboard panel](#clipboard-panel)
   - [Sharing: members, roles & links](#sharing-members-roles--links)
   - [Real-time collaboration](#real-time-collaboration)
   - [Voice calls](#voice-calls)
   - [Settings](#settings)
   - [Audit log, session recording & file transfers](#audit-log-session-recording--file-transfers)
   - [Admin panel](#admin-panel)
   - [Mobile & responsive UI](#mobile--responsive-ui)
6. [Keyboard shortcuts](#keyboard-shortcuts)
7. [Configuration](#configuration)
8. [Docker, service, reverse proxy & firewall](#docker-service-reverse-proxy--firewall)
9. [Security model](#security-model)
10. [Limitations](#limitations)
11. [API reference](#api-reference)
12. [Building from source & releases](#building-from-source--releases)
13. [Troubleshooting](#troubleshooting)
14. [Contributing](#contributing)

---

## What WRM is

WRM is a **single executable** with a built-in web server and a built-in web app. You start it on any machine: your PC, a Raspberry Pi, a VPS or a company server. Then you open it in a browser and manage your remote servers from there:

- **SSH terminals** in the browser (xterm.js, 256 colors, full-screen apps like `nano`, `vim`, `htop`, `mc`), with **snippets** (saved commands with variables, one click or Ctrl+Shift+Space), **run on connect** (`sudo -i`, `cd /srv/app`…) and **broadcast input** (type into many terminals at once, with a safety check for dangerous commands).
- **Live up/down status** of every saved connection in the sidebar, with latency, server banner and a notification when a server goes down.
- **File manager** for **SFTP** (over SSH), **FTP** and **FTPS**: browse, upload, download, rename, delete, edit, search.
- **Jump hosts** (like `ssh -J`): reach servers behind a bastion, also in chains — for terminals, files, transfers and tunnels.
- **SSH tunnels** (like `ssh -L / -R / -D`, PuTTY, mRemoteNG): reach the web interface of a switch, an iDRAC/iLO or a database behind a server; a SOCKS proxy into a whole management network. **Web interface connections** open such pages with one double-click.
- **Import** your connections from **mRemoteNG** (with passwords, folders and SSH tunnels) and from **OpenSSH** `~/.ssh/config`.
- **Server-to-server copy** between two SSH servers, without downloading to your computer first.
- **Workspaces**: many terminal and file windows side by side, tabs, snapping, saved sessions.
- **Sharing with roles**: give colleagues or guests access to some connections — as *Observer*, *Viewer*, *Operator* or *Moderator* — without ever revealing the passwords.
- **Live collaboration**: chat with history, file exchange, **voice calls** (like Discord/Jitsi), terminal sharing with snapshots, remote keyboard control on request, raised hands and moderation.
- **Enterprise security**: two-factor authentication, closed registration, admin-managed accounts, password & session policies, account lockout, audit log (with CSV export and SIEM-friendly log lines), SSH host-key verification, encrypted secrets, CSRF protection and strict security headers.

Everything is stored in one local **SQLite** file. There is no external database, no Docker requirement, no CDN and no cloud account: the binary also serves all scripts and fonts, so WRM works in closed networks.

**Supported platforms (prebuilt binaries):** Linux (x86-64, arm64, ARMv7/Raspberry Pi, ARMv6, 32-bit), Windows (x86-64, arm64), macOS (Intel, Apple Silicon, Universal), Android arm64 (Termux), FreeBSD (x86-64, arm64), OpenBSD (x86-64).

---

## Quick start

1. Download the binary for your system from **[Releases](https://github.com/vedranius/web-browser-RDM-public/releases/latest)** and optionally verify it with `SHA256SUMS.txt`.
2. Run it:

   **Linux / macOS / FreeBSD / OpenBSD**
   ```bash
   chmod +x wrm-pro-v10.3.0-linux-amd64
   ./wrm-pro-v10.3.0-linux-amd64
   ```
   On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

   **Windows** (PowerShell), or just double-click the `.exe`:
   ```powershell
   .\wrm-pro-v10.3.0-windows-amd64.exe
   ```

   **Android (Termux)**
   ```bash
   pkg install wget
   wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v10.3.0/wrm-pro-v10.3.0-android-arm64
   chmod +x wrm-pro-v10.3.0-android-arm64 && ./wrm-pro-v10.3.0-android-arm64
   ```

   **Docker**
   ```bash
   git clone https://github.com/vedranius/web-browser-RDM-public.git && cd web-browser-RDM-public
   docker compose up -d          # https://<host>:8080 (self-signed certificate)
   ```
   See [Docker](#docker) for details.
3. Open **http://localhost:8080** (or `http://<server-ip>:8080`). For voice calls from other computers use HTTPS — the quickest way is `HTTPS_SELF_SIGNED=1` (see [Configuration](#configuration)).
4. **Create the administrator account** (the first account). After that, self-registration is **closed**: the administrator creates accounts in *Admin panel → Users* (or opens registration in *Security policies*).
5. Click **+ Connection**, enter host, user and password or key, then double-click the connection to open it.

On the first start WRM creates, next to the database:

| File | What it is |
|---|---|
| `remote_manager.db` | SQLite database (users, connections, shares, audit log). Mode `0600`. |
| `remote_manager.db.key` | Random 256-bit key that encrypts stored passwords, private keys and 2FA secrets. Mode `0600`. **Back it up separately** — without it stored secrets cannot be decrypted. Or set your own key with `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`. |
| `recordings/` | Terminal session recordings (`YYYY/MM/<id>.cast.gz`, mode `0600`), created when the first session is recorded. Location: `WRM_RECORDINGS_DIR`. |

Locked out? `./wrm-pro-… -reset-password admin` prints a new temporary password (add `-reset-2fa` to also turn off two-factor authentication).

---

## Upgrading from v10.2

Replace the binary. The database gets a new table (`snippets`) and one new column (`connections.monitor`); the previous binary still starts on it.

- **File names no longer carry a suffix:** releases are `v10.3.0`, binaries `wrm-pro-v10.3.0-linux-amd64` etc. Update scripts and systemd units that use the file name.
- **Live status is on by default:** WRM opens a TCP connection to every saved host:port once a minute (SSH and FTP: it reads the greeting and closes politely, like Nagios `check_ssh`). Servers see it in their logs. Change the interval or turn it off in *Admin panel → Security policies → Live status*, or untick *Monitor up/down status* in a connection. See [Live up/down status](#live-updown-status) about fail2ban.
- **Broadcast input** is available to everybody; turn it off with `broadcast_enabled` if you do not want it. Every start and stop is in the audit log.

## Upgrading from v10.1

Replace the binary. The database gets a new table (`connection_tunnels`) and two new columns on `connections` (`jump_conn_id`, `web_path`); nothing else changes. To go back, run the v10.1 binary on the same database; it does not know these features (connections behind a jump host then try to connect directly, web interface connections cannot be opened, tunnels do not run).

- New: [jump hosts](#jump-hosts-bastions), [SSH tunnels](#ssh-tunnels-port-forwarding), [web interface connections](#web-interfaces-http--https-connections) and [import from mRemoteNG / OpenSSH](#import-from-mremoteng-and-openssh).
- **Tunnels are on by default for all users**, listening on `127.0.0.1` of the WRM machine only. Remote forwards (`-R`) and listening on network addresses are limited to administrators. Review *Admin panel → Security policies → SSH tunnels* — e.g. turn tunnels off (`WRM_TUNNELS_ENABLED=0`) or limit them to administrators if WRM runs on a shared server.
- Jump hosts need `AllowTcpForwarding yes` on the bastion (the OpenSSH default), exactly like `ssh -J`.

## Upgrading from v10.0

Replace the binary. The database is extended automatically (new tables, columns and triggers only — nothing is changed or removed; the previous binary still runs on the upgraded database).

- **Terminal sessions are now recorded by default** (`session_recording`). Users see *“This session is recorded”* and a **● REC** badge on the terminal. Inform your users (in many countries this is required), or turn it off in *Admin panel → Security policies → Session recording*. Keystrokes are **not** recorded unless you turn on *Record keystrokes*.
- Recordings use disk space next to the database (`recordings/`, typically a few KB to a few MB per session; limit per session `recording_max_mb`, kept `recording_retention_days` = 90 days). Back up or exclude that folder as you need.
- The audit log is now **append-only and hash-chained**. *Admin panel → Audit log → Verify integrity* checks it. Entries written before the upgrade have no hash and are skipped by the check.
- Every upload, download (also files opened in the editor and files in ZIP downloads) and server-to-server copy is logged in *Admin panel → File transfers* with size and SHA-256.
- New: `/healthz` for load balancers and monitoring, `wrm -healthcheck` for container health checks, a `Dockerfile` and `docker-compose.yml`.

## Upgrading from v9

Replace the binary and start it with the same database and the same `ENCRYPTION_KEY` (if you set one). Everything is migrated automatically. Things to know:

- **Stored secrets are re-encrypted.** Without `ENCRYPTION_KEY`, v9 used a public default key. v10 generates a random key file (`remote_manager.db.key`) and re-encrypts all secrets with it on the first start. Back that file up. (If you did set `ENCRYPTION_KEY`, nothing changes.)
- **Self-registration is closed by default.** Existing accounts keep working. Open it again in *Admin panel → Security policies* if you want it.
- **Passwords and keys are no longer sent to the browser.** The *Edit connection* dialog shows “saved and encrypted — leave empty to keep it”. *Export* is without secrets; *Export with passwords & keys* asks for your password.
- **SSH host keys are verified** (trust on first use). The first connection to each server remembers its key; a changed key is refused until you accept it.
- **Key file / Auto (~/.ssh)** authentication (keys stored on the WRM server) is limited to administrators unless the policy *Allow server key files for all users* is on.
- **Existing shares keep their behaviour** (anyone with the link, role *Operator*). Edit them to use the new access modes and roles.
- **Behind a reverse proxy:** state-changing API calls now require the `Origin` to match the `Host` (like WebSockets already did). Pass the `Host` header, or set `WRM_ALLOWED_ORIGINS`. Set `WRM_TRUST_PROXY=1` so rate limits and the audit log see real client IPs.
- The per-IP sign-in limit is now 20 failures / 10 minutes; accounts additionally lock for 15 minutes after 10 wrong passwords (policy).
- If the database comes from an earlier build that already had tables named like the new v10 tables (`audit_log`, `collab_messages`, …) but with a different layout, WRM keeps them as `<table>_old_<time>` and creates new ones (logged as *Upgraded table …*). Nothing is deleted.
- The app is served at `/`. Old bookmarks to `/static/` are forwarded there.

---

## How it works

```
 Browser (any device)                      WRM server (one binary)                     Your servers
┌──────────────────────┐   HTTPS/HTTP   ┌──────────────────────────────┐   SSH / SFTP   ┌───────────┐
│ Web UI (embedded     │ ─────────────▶ │ REST API  (/api/...)         │ ─────────────▶ │ Linux/BSD │
│ single-page app,     │                │ WebSockets:                  │   FTP / FTPS   │ servers,  │
│ xterm.js terminal)   │ ◀───────────── │  /ws/ssh     terminal        │ ─────────────▶ │ NAS, ...  │
│                      │   WebSockets   │  /ws/share/  collaboration   │                └───────────┘
└──────────┬───────────┘                │  /ws/events  live updates    │
           │  voice (WebRTC, DTLS-SRTP) │ TURN relay  UDP/TCP 3478     │
           └────────────────────────────┤ SQLite: users, connections,  │
             peer-to-peer or via relay  │ shares, audit log, settings  │
                                        └──────────────────────────────┘
```

- **The WRM server makes the connections.** Your browser talks only to WRM, and WRM talks to your servers with SSH (`golang.org/x/crypto/ssh`), SFTP (`pkg/sftp`) and FTP (`jlaffaye/ftp`). Your servers only need to be reachable **from the WRM machine**.
- **Terminal:** every terminal window is a WebSocket (`/ws/ssh`) bound to one SSH session with a PTY. Keystrokes and output go as **binary frames**, so large outputs and multi-byte UTF-8 are never corrupted. The PTY size follows the window.
- **Flow control:** when a command prints faster than the browser can draw, the browser asks the server to pause and resume. Nothing is lost.
- **File operations** use a small **pool of SFTP connections** per saved connection (reused for 5 minutes and health-checked). Uploads are **streamed**; downloads and ZIP archives are streamed to the browser.
- **Secrets at rest** (connection passwords, private keys, 2FA secrets) are **AES-256-GCM encrypted**. Session tokens are stored only as SHA-256 hashes.
- **Collaboration:** `/ws/share/{token}` is the room of a share. The server decides identities and roles and relays chat, files, terminal streams (only to people watching them), remote-control keystrokes (only to the person who granted control) and voice signalling.
- **Jump hosts and tunnels:** WRM logs in to the jump host, opens an SSH `direct-tcpip` channel through it and runs a second, separately encrypted SSH session inside it (like `ssh -J`). Tunnels listen on the **WRM machine** (local / SOCKS) or on the SSH server (remote) and forward every connection through SSH. Tunnel ports are opened by WRM, not by the browser.
- **Voice:** audio flows **directly between browsers** (WebRTC, end-to-end encrypted with DTLS-SRTP). When browsers cannot reach each other (NAT, firewalls), the **built-in TURN relay** carries the still-encrypted audio.

---

## Features in detail

### Accounts, sign-in & two-factor authentication

- The **first account becomes administrator**. Afterwards **self-registration is closed** (policy); administrators create accounts, optionally with a generated temporary password that must be changed at the first sign-in.
- Passwords are stored as **bcrypt** hashes (cost 12). Minimum length is a policy (default 8).
- **Two-factor authentication (TOTP)** with any authenticator app (Google/Microsoft Authenticator, 1Password, Bitwarden, FreeOTP…): *Settings → Security → Two-factor authentication*. A QR code is shown; **10 single-use recovery codes** are generated. Codes cannot be replayed. Administrators can **require 2FA** for administrators or for everyone — users without it must set it up at their next sign-in.
- **Brute-force protection:** 20 failed sign-ins from one IP within 10 minutes lock that IP for 5 minutes; an **account** locks for 15 minutes after 10 wrong passwords (policy). Failed sign-ins take the same time whether the user exists or not.
- **Sessions:** HttpOnly, SameSite cookie (`Secure` on HTTPS); the database only stores a hash of the token. Sessions end after inactivity (default 7 days) and after a maximum age (default 30 days) — both policies. *Settings → Security* lists your **signed-in devices** and signs out one or all others. Changing your password signs out every other device.
- If your login expires while terminals or transfers are running, WRM shows a message instead of reloading the page.

### Connections & folders

Left sidebar ("PRO MANAGER"):

- **+ Connection** creates a connection:
  - **Protocol:** `SSH` (terminal, plus files over SFTP), `SFTP` (opens the file manager by default), `FTP`, `FTPS` (FTP with explicit TLS; file manager only), **🌐 Web interface** `HTTPS` / `HTTP` (see [Web interfaces](#web-interfaces-http--https-connections)).
  - **Host : Port** (default port 22 for SSH/SFTP, 21 for FTP/FTPS, 443/80 for web interfaces), **Username**, **Folder**.
  - **Jump host (connect via):** another saved SSH connection to go through (see [Jump hosts](#jump-hosts-bastions)).
  - **🔀 SSH tunnels:** port forwards of this connection (see [SSH tunnels](#ssh-tunnels-port-forwarding)).
  - **Authentication:** *Password*, *Private key (paste)*, *Key file (path on the WRM server)*, or *Auto* (tries `~/.ssh/id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa` of the WRM server user). *Key file* and *Auto* use keys of the **WRM server itself** and are therefore limited to administrators (policy). Password authentication also answers keyboard-interactive prompts.
  - **🔌 Test connection** logs in once and shows the result, latency and — for a new server — its host key fingerprint.
- **Secrets never leave the server.** Passwords and private keys are write-only: the edit dialog shows “saved and encrypted — leave empty to keep it”, with an option to remove the saved secret. *Duplicate* copies the connection on the server.
- **Host key verification:** the first connection to a server remembers its SSH host key (and, for FTPS, the certificate if it is not signed by a public CA). A **changed key is refused** with a clear warning and both fingerprints; the owner of the connection (or an administrator) can accept the new key after checking it. Policy: *trust on first use* (default), *strict* (new hosts must be approved) or *off*.
- **Double-click** a connection to open it: SSH opens a terminal; SFTP/FTP/FTPS open the file manager. On touch devices a single tap opens it.
- **Right-click** a connection: *Open as Terminal / File Manager* (in a new or an existing window), *Tunnels*, *Start / Stop all tunnels*, *Add tunnel…*, *Edit*, *Duplicate*, *Share*, *Delete*. Web interface connections: *Open web interface*.
- **📁 Folder** creates folders. **Drag** connections onto a folder to move them, or onto *"↓ Drop here"* to move them back to the root. Right-click a folder to *Share* or *Delete* it (its connections move to the root).
- **Search box** filters by name, host, username or protocol. **Multi-select** with **Ctrl/Cmd + click**, then use **🤝 Share** or **Delete**.
- The **dot in front of a connection** is its [live status](#live-updown-status): 🟢 up, 🟡 up but slow (> 300 ms), 🔴 down, a ring = behind a jump host (colored by the jump host's state). A **blue dot after the name** means one of its terminals is connected. **⤳** means it goes through a jump host (hover shows the route), **🔀** that it has tunnels (colored while one runs; click opens them), **WEB** marks web interfaces.
- **Shared with me** lists shares you are a member of (and shares published for all users), with your **role**; ↗ opens the collaboration room.

### Jump hosts (bastions)

Servers in data centers and customer networks are often reachable only through a bastion ("jump host"). WRM goes through it like `ssh -J` / `ProxyJump` and like the *SSH tunnel* setting of mRemoteNG:

1. Create the bastion as a normal SSH connection, e.g. `bastion-dc1` → `203.0.113.10:22`.
2. Create the internal server, e.g. `app-01` → `10.0.0.21:22`, and choose **Jump host (connect via) → bastion-dc1**.
3. Double-click `app-01`. The terminal shows `Connecting to 10.0.0.21:22 via bastion-dc1…`.

```
 WRM ──SSH──► bastion-dc1 ──SSH inside the first connection──► app-01 (10.0.0.21)
```

- **Everything works through it:** terminal, file manager (SFTP), FTP/FTPS (the FTP connection is made from the jump host), search, editor, server-to-server transfer, *Test connection*, tunnels and web interfaces.
- **Chains:** a jump host can have its own jump host — up to 5 hops, e.g. `vpn-gw → bastion-dc1 → app-01`. Loops are refused when you save.
- **Every hop logs in with its own credentials** (password, key, key file) and **its host key is verified** (trust on first use / strict), like a direct connection. WRM opens no port on the jump host; it uses a `direct-tcpip` channel, so the bastion only needs `AllowTcpForwarding yes` (the OpenSSH default).
- The route is shown in the sidebar (**⤳**), in *Test connection* (`… via bastion-dc1`), in the terminal, in *Sessions & recordings* and in the audit log.
- Jump hosts are your own saved SSH connections. **Shares** of a connection behind a jump host work too — members never see the credentials of either.
- Deleting a jump host makes the connections that used it direct again.

### SSH tunnels (port forwarding)

Like `ssh -L / -R / -D`, the tunnels of PuTTY and the SSH tunnel feature of mRemoteNG. You configure tunnels in a connection (**✏ Edit → 🔀 SSH tunnels**, or right-click → *Add tunnel…*) and start and stop them in the **🔀 Tunnels** panel in the top bar.

| Type | What it does | Example |
|---|---|---|
| **Local (-L)** | A port on the **WRM machine** leads, through the SSH connection, to a host:port **as seen from that server**. | Listen `127.0.0.1:2443` → target `10.0.0.5:443`: open `https://127.0.0.1:2443` to reach the web UI of a switch behind the server |
| **SOCKS proxy (-D)** | A SOCKS5 proxy on the WRM machine. Every destination is reached through the server. | `127.0.0.1:1080`: set it in the browser (e.g. FoxyProxy, *SOCKS v5, proxy DNS*) and browse the whole management network, or `curl --socks5-hostname 127.0.0.1:1080 http://10.0.0.5` |
| **Remote (-R)** | A port **on the SSH server** leads back to a host:port reachable **from WRM**. | `127.0.0.1:9000` on the server → `10.10.0.2:80`: the server downloads from a mirror on your side |

**Fields of a tunnel**

- **Listen:** address and port. `127.0.0.1` (default) means only programs on the WRM machine can use the tunnel. Leave the port empty for a free port chosen automatically (shown while running). Local and SOCKS ports must be 1024–65535. A network address (`0.0.0.0`, a LAN IP) makes the tunnel reachable by **anyone who can reach that port, without signing in to WRM** — administrators only, unless the policy allows it.
- **Target:** host:port, as seen from the server (local) or from WRM (remote). Not used by SOCKS.
- **Open as** `http`/`https` + path: adds a **🌐 Open** button that opens the tunnel in a new browser tab.
- **Start:** *Manual* (▶ in the Tunnels panel); *With terminal* (starts with the first terminal of the connection and stops about 15 s after the last one closes); *Always* (starts with WRM and keeps running, also after a restart).
- **Templates:** web interface (HTTPS 443 / HTTP 80), SSH, RDP and VNC to an internal host, PostgreSQL, MySQL/MariaDB, iDRAC / iLO / IPMI web, SOCKS proxy.

The hint under each tunnel says in plain words what it will do, e.g. *Connections to 127.0.0.1:2443 on the WRM machine go through Site1 to 10.0.0.5:443.*

**The 🔀 Tunnels panel** lists your tunnels by connection: state (🟢 running, 🟡 reconnecting, 🔴 error), listen → target, open / total connections, traffic ↑↓, since when and who started it, with **🌐 Open**, **📋 Copy address** (for SOCKS as `socks5h://…`), **▶ Start** and **■ Stop**. The number on the button counts running tunnels. Administrators can show **all users'** tunnels and stop them (also in *Admin panel → Overview*).

**Reliability:** every running tunnel has its own SSH connection (through the jump hosts), checks it every 30 seconds and **reconnects automatically** (after 2, 5, 10, 30, 60 s…) while its port stays open. Changing a connection restarts its running tunnels; deleting it stops them. Each start, stop (with the bytes transferred) and error is in the **audit log**.

**Where is the tunnel port?** On the computer where **WRM runs**:

- **WRM on your own PC** (the usual case): open `http://127.0.0.1:<port>` in your browser, or point any program (DBeaver, RDP client, `psql`, FileZilla…) at `127.0.0.1:<port>`.
- **WRM on a server you reach over the network:** a `127.0.0.1` tunnel is usable only on that server. An administrator can use listen address `0.0.0.0` (or the server's LAN address) and allow the port in the firewall — then everyone who can reach that port can use the tunnel, so prefer a restricted network. Or forward the port to your PC with plain SSH (`ssh -L 2443:127.0.0.1:2443 wrm-server`).

**Example — the web interface of a server behind a site gateway** (PC → WRM → `Site1` → `Site1-Webserver:443`):

1. Connection `Site1` (SSH to the site gateway).
2. *Edit `Site1` → 🔀 SSH tunnels → + Local (-L)*: listen `127.0.0.1:2443`, target `10.1.0.20:443` (the web server as seen from `Site1`), open as `https`.
3. *🔀 Tunnels → ▶ Start → 🌐 Open* → `https://127.0.0.1:2443`.

Even simpler: save it as a [web interface connection](#web-interfaces-http--https-connections) with jump host `Site1`.

### Web interfaces (HTTP / HTTPS connections)

Save the web interfaces of your devices — iDRAC / iLO / IPMI, switches, firewalls, NAS, hypervisors, monitoring — as connections with protocol **🌐 Web interface (HTTPS)** or **(HTTP)**: host:port, an optional **path** (e.g. `/login.html`) and an optional **jump host**.

- **Double-click** (or right-click → *Open web interface*) opens it in a new browser tab.
- **Without a jump host** the browser opens the address directly (`https://10.0.0.200/login.html`).
- **With a jump host** WRM opens a **temporary tunnel** on the WRM machine (`127.0.0.1:<free port>` → host:port as seen from the last jump host) and opens that. The tunnel is reused while you work and stops after 30 minutes without traffic; the Tunnels panel shows it as *temporary*.
- *Test connection* checks that the port answers (through the jump hosts).
- TLS goes end to end from your browser to the device; for self-signed device certificates the browser shows its usual warning.
- As with all tunnels, the temporary tunnel is on the WRM machine: when WRM runs on another computer, WRM shows the address instead of opening it (see *Where is the tunnel port?* above).

### Import from mRemoteNG and OpenSSH

*Settings → Data*:

- **Import from mRemoteNG** — the `confCons.xml` file (`%APPDATA%\mRemoteNG\confCons.xml`, or *File → Export* in mRemoteNG):
  - containers become folders (`DC1 / Rack A`); SSH1/SSH2 connections become SSH connections, HTTP/HTTPS connections become web interface connections (with the path of the URL);
  - **passwords are decrypted** (current AES-GCM format, older AES-CBC format, *FullFileEncryption*) and stored encrypted in WRM. If the file is protected with a **master password**, WRM asks for it;
  - inherited settings (*Inherit* user name, password, port, domain, SSH tunnel) are resolved from the parent folders;
  - the **SSH tunnel** setting (`SSHTunnelConnectionName`) becomes the **jump host**;
  - RDP, VNC, Telnet, ICA and other protocols WRM does not open are listed as *skipped*; connections that already exist (same name, host and protocol) are skipped too, so importing again is safe.
- **Import OpenSSH config** — `~/.ssh/config` (folder *SSH config*):
  - `Host` entries with `HostName`, `User`, `Port`; defaults from `Host *` and from wildcard patterns;
  - **`ProxyJump`** (also `user@host:port`; a jump host that is not its own `Host` entry is created) and `ProxyCommand ssh -W …` become the **jump host**;
  - **`LocalForward`, `RemoteForward`, `DynamicForward`** become tunnels that start with the terminal;
  - `IdentityFile` becomes *Key file* authentication (only if key files on the WRM server are allowed for you; otherwise a note asks you to add the key or a password);
  - `Match` blocks and wildcard-only hosts are skipped.

After the import a summary lists the imported connections, folders, jump host links and tunnels, what was skipped and what to check (e.g. connections without a user name or password).

### Windows, tabs & snapping

Every opened connection is a **window** inside the workspace. There is also a **tab** for it in the tab bar.

- **+ New Window** (or **Alt+N**) opens an empty window with quick buttons for your connections and **saved sessions**.
- **Title bar:** drag to move, **double-click** to maximize or restore. Buttons: *snap layouts*, **↻ Reconnect** (SSH windows), 📡 *share terminal* (in a collaboration room), minimize, maximize, close.
- **Snapping:** left, right, **top and bottom halves**, the four **quarters**, and **maximize**. Use the buttons in the title bar, the compact *snap layouts* menu in narrow windows, or the tab's right-click menu.
- **Drag-to-edge snapping:** drag a window with the pointer to an edge (half) or a corner (quarter) of the workspace. A blue preview shows where it will go.
- **Windows can never leave the workspace:** moving and resizing stop at all four edges. When the browser is resized or zoomed, or a panel changes size, windows are moved and shrunk to stay visible. Snapped and maximized windows keep their slot.
- **Resize** from any edge or corner. Dragging a snapped window restores its previous size.
- **Tabs:** click to focus or restore, **drag to reorder**, **middle-click** to close, **right-click** for *Reconnect*, *Duplicate window*, snap layouts and *Close*. SSH tabs show the connection state (🟡 connecting, 🟢 connected, 🔴 disconnected).
- **Close confirmation** can be *Always*, *Never* or *Only when connected* (Settings).
- The browser warns you before leaving the page while terminals are connected or a transfer is running.

### Snippets & run on connect

Snippets are saved commands. Open them in the right panel (**⚡ Snippets** tab), or press **Ctrl+Shift+Space** in a terminal (or ⚡ in its title bar) for a quick picker.

- **Click** a snippet: it is typed into the focused terminal and run. **Shift+click** (or ⎘) types it without Enter, so you can edit it first. **Right-click:** run, insert, copy, edit, duplicate, delete.
- **Quick picker:** type to filter, ↑ ↓ to choose, **Enter** runs, **Shift+Enter** inserts, **Esc** closes, **＋ New snippet** saves the selected terminal text as a snippet.
- **Groups** (Diagnostics, Services, …), a description, and the scope: **all connections**, the connections of **one folder**, or **one connection**. The picker only offers snippets that apply to the terminal's connection.
- **Multi-line snippets:** every line is one command.
- **Variables** are filled in for each terminal: `{{host}}`, `{{port}}`, `{{user}}`, `{{name}}` (connection), `{{folder}}`, `{{wrm_user}}`, `{{date}}`, `{{time}}`. **Prompts** are asked before the snippet runs, with a live preview: `{{?Service}}` or with a default `{{?Service=nginx}}`. Example: `journalctl -u {{?Service=nginx}} --since "{{?Since=1 hour ago}}" --no-pager`.
- **Run on connect:** tick *Run on connect* and WRM types the snippet into the shell as soon as a terminal of a matching connection is connected (also after a reconnect) — `sudo -i`, `cd /srv/app`, `tmux attach || tmux new`, `export KUBECONFIG=…`. Right-click a connection or a folder → **⚡ Run on connect…** creates one. Order: all-connection snippets, then folder, then connection snippets. A message lists what ran; the audit log records it (`terminal.auto_run`).
- **Shared snippets:** administrators can share a snippet with all users (read-only for them, all connections, never auto-run) — a team runbook in the terminal.
- **✨ Add example snippets** creates a starter set (disk, memory, top processes, failed services, journal errors, listening ports, big files, service status/restart…).
- Snippets are included in *Export* / *Import*, copied with *Duplicate*, and removed with their connection or folder.
- Dangerous-looking commands (see below) ask for confirmation before they run.

### Broadcast input

Type the same command into many servers at once — like *MultiSSH*, *cssh* or mRemoteNG's multi-SSH.

1. Open the terminals (e.g. all web servers of a cluster).
2. Top bar → **📣 Broadcast** → select the terminals → **Start broadcast** (asks once).
3. Everything you type in one of them now goes to all of them. An orange bar shows how many terminals receive it; included windows get an orange frame. **📣** in a window's title bar adds or removes that terminal; **■ Stop** ends it.

Safety:

- **Dangerous commands ask first:** before Enter is sent for a line that looks dangerous — `rm -r/-f`, `shutdown`/`reboot`/`halt`, `init 0/6`, `systemctl stop/restart/disable…`, `mkfs`, `dd … of=`, `fdisk`/`parted`, `wipefs`, writing to `/dev/sd*` or `/etc/…`, `kill -9`, `pkill`, `chmod 777`, `chown -R`, `iptables -F`, `ufw disable`, `DROP/TRUNCATE TABLE`, `DELETE FROM`, `crontab -r`, package removal, `docker rm`, `kubectl delete`, `terraform destroy`, a fork bomb … — WRM shows the command and the number of terminals; **Cancel is the default button**.
- **Pasting several lines** into a broadcast asks for confirmation too.
- **Snippets** run on every terminal of the broadcast, each with its own `{{host}}`/`{{user}}`, after one confirmation.
- The check works on what you typed in this broadcast (not on commands recalled from the shell history with ↑).
- Each terminal records when it joined and left a broadcast (`terminal.broadcast`, with the number of terminals) and, as always, the session recording shows what happened on each server. Administrators can turn the feature off (`broadcast_enabled`).

### Live up/down status

WRM checks every saved connection in the background and shows the result in the sidebar:

- 🟢 **up** (hover: latency, since when, SSH server version), 🟡 **slow** (more than 300 ms), 🔴 **down** (pulsing; hover: since when and why — *connection refused*, *timeout*, *host name not found* …), a **ring** for connections behind a jump host.
- A folder shows how many of its connections are down (**2↓**). When something is down, a red bar above the tree shows **only the connections that are down** with one click.
- **Notifications:** a message when a server goes down or comes back, and a desktop notification while WRM is in the background (*Settings → General → Notify when a server goes down*).
- **🔄 Check status now** (right-click a connection or a folder) checks immediately — also **through the jump hosts** — and shows the result per connection.
- **How it checks:** a TCP connection to the connection's port (SSH 22, FTP 21, web 80/443 or the port you set). SSH checks read the server banner and answer with `SSH-2.0-WRM_status_check` before closing, FTP checks read the 220 greeting and send `QUIT` — the same as Nagios `check_ssh`/`check_ftp`. Nothing logs in. A failed check is repeated once before a server counts as down. Each distinct host:port is checked once per round, however many connections and users point to it.
- **Behind a jump host** a connection shows the state of its jump host (a red ring when the jump host is down). With the policy *Check through jump hosts*, WRM logs in to the jump host once per round and checks the ports behind it through that SSH connection.
- Untick **Monitor up/down status** in a connection to leave it out (e.g. servers that are often off on purpose).
- Policies: `status_enabled` (on), `status_interval_seconds` (60, 15–3600), `status_jump_checks` (off).

> **fail2ban / IDS:** the check is an ordinary connection without a login. Default fail2ban `sshd` filters do not count it, but the *aggressive* / *ddos* modes may. Add the WRM server to `ignoreip`, raise the interval, or untick monitoring for such hosts.

### SSH terminal

- Full **xterm-256color** terminal (xterm.js 5, served by WRM itself — no CDN) with the JetBrains Mono font. Works with `nano`, `vim`, `less`, `htop`, `mc`, `tmux`, colors, mouse reporting and UTF-8 (č, ć, ž, š, đ, emoji…).
- **Correct size everywhere:** the PTY size follows the window on connect, maximize, snap, resize, sidebar toggle, font change and browser zoom.
- **Large output is safe:** binary WebSocket frames plus flow control, so there are no disconnects on big `cat`/`tail`/`journalctl` output.
- **Reconnect:**
  - **↻** button in the title bar (also while connected, e.g. when a session hangs), **Ctrl+Shift+R**, or the tab menu.
  - When a terminal is closed, a banner offers **Reconnect**, and **Enter** reconnects too.
  - **Automatic reconnect** when the SSH connection or the connection to the web server drops (backoff 1 s → 30 s, up to 10 attempts), and immediately when the network comes back. It can be turned off in Settings. If the shell simply exits (`exit`), it does not reconnect automatically.
- **Keepalive:** WRM pings the browser every 25 s and the SSH server every 30 s and detects dead connections.
- **Copy on select:** selecting text copies it to the clipboard and adds it to the [Clipboard panel](#clipboard-panel). **Ctrl+Shift+C** copies too, and **Ctrl+Shift+V** / normal paste pastes.
- **Search in the scrollback:** **Ctrl+Shift+F** (Enter = next, Shift+Enter = previous).
- **Clickable URLs** in the output.
- **Browser shortcuts are blocked while a terminal is focused:** **Ctrl+W** is sent to the terminal (nano "Where Is", bash "delete word") instead of closing the browser tab. Ctrl+T, Ctrl+N and Ctrl+Q are blocked too.
- Configurable **font size** (11–18 px) and **scrollback** (5k–50k lines).
- **"Open terminal here"** from the file manager opens a new terminal and `cd`s into that folder.

### File manager (SFTP / FTP / FTPS)

- **Breadcrumb path bar:** click any part to jump there, or click empty space to type a path (Enter). **↑ Up** and **⟳ Refresh** (F5).
- **Columns:** name (with icons by file type, symlink marker), **Modified**, size (B … TB), type. Click a column header to sort (folders always first, natural sort).
- **Selection:**
  - click to select, **Ctrl/Cmd + click** to toggle, **Shift + click** for a range, checkbox, *select all*, **Ctrl+A**
  - keyboard: arrows, Home/End, **Enter** (open), **Backspace** (up), **Delete**, **F2** (rename), **F5**, **Ctrl+F** (search), Esc
- **Double-click:** folder opens; text file opens in the [editor](#text-editor); other files download.
- **Upload:**
  - **⬆ Upload** (many files at once), **🗂 Folder** (a whole folder with subfolders), or **drag & drop** files or folders from your computer onto the list
  - progress bar with speed, cancel button
  - existing files are detected, and you are asked whether to **overwrite** them
  - optional per-file size limit (Settings; administrators can enforce a limit on the server)
- **Download:** single files, or **folders as ZIP** (SFTP). **⬇ Download** downloads all selected items.
- **📁 New Folder**, **✏️ Rename** (F2), **🗑 Delete** (selected files and folders; folders recursively, with confirmation).
- **⚡ Transfer** sends the selection to another server (see [Server-to-server transfer](#server-to-server-transfer)).
- **⌨ Terminal** opens a terminal in the current folder (SSH/SFTP only).
- **Right-click menu:** Open, Edit, Download / Download as ZIP, Transfer to server, Open terminal here, Copy path, Rename, Delete. It works on multiple selected items.
- **Status bar:** number of folders and files, selection count and size, connection name.
- In narrow windows the toolbar switches to icons only.

### Search in files

The search box at the top of every file manager:

- **Typing filters the current folder instantly.** Wildcards work: `*.log`, `nginx*.conf`, `?`.
- **Enter** (or ⌕) searches **all subfolders** of the current folder on the server. Results stream in live with a progress indicator (folders scanned, current folder) and can be **stopped** at any time.
- Search options (▾):
  - **File names** (default) walks the tree over SFTP/FTP. Works on every server, including Windows SFTP and FTP.
  - **File contents** runs `grep` on the server over SSH: case-insensitive, fixed string, text files only. It needs a POSIX shell with `grep` (Linux/BSD/macOS servers). Not available for FTP.
  - include/skip **hidden** files; result **limit** (200 / 500 / 2000 / 5000).
- When searching from `/`, the virtual filesystems `/proc`, `/sys`, `/dev`, `/run` and `/snap` are skipped. A search stops after 3 minutes or 1,000,000 scanned entries, and tells you if it was truncated.
- Results show the file, its folder, size and date. **Double-click** a result to open its folder with the file highlighted. Results support selection, the context menu, download and delete.

### Text editor

- Opens text and config files up to **5 MB** directly in the browser (double-click, or right-click → *Edit*). Binary files show a warning first.
- **Ctrl+S** saves (uploads back to the server), **Tab** inserts a tab, **Esc** closes and asks if there are unsaved changes.
- Keeps the file's **line endings** (LF or CRLF) and UTF-8 encoding. Shows line count and a *modified* marker.

### Server-to-server transfer

Copy files and folders **directly from one SSH/SFTP server to another**. The data flows through the WRM server, not through your browser.

- Open with **⚡ Transfer** in a file manager or from the context menu. The dual-panel dialog shows the source (with your selection) on the left and the destination connection and folder on the right.
- **If a file exists:** *Ask each time* (checks for conflicts first and asks once for all), *Overwrite all*, *Skip existing*, *Overwrite if newer*.
- Folders are copied recursively with their structure. File modification times are preserved. Up to **4 files are copied in parallel**.
- Live progress: per-file progress, total bytes, speed (sliding window), **ETA**, results list with errors.
- **Minimize** the dialog to a small floating progress bar (draggable, always kept on screen) and keep working. **Cancel** stops the transfer on the server too.
- A sound plays when the transfer finishes (success or error).
- FTP connections are not supported as source or destination.

### Workspace sessions

Right panel → **Sessions**:

- **+ Save Session** saves **all open windows**: connection, terminal or file manager, current folder, snap slot or position (stored relative to the workspace size, so it fits other screens), minimized and focused state. The suggested name is built from the connections.
- **Click a saved session to open it.** If windows are already open, choose **Replace** (close them) or **Add** (open next to them). Saved sessions are also offered in every empty *New Window*.
- **Right-click** a session: *Open (replace)*, *Open (add)*, **Save current windows here** (update), *Rename*, *Lock/Unlock* (a locked session cannot be overwritten, and deleting it needs confirmation), *Delete*.
- Sessions belong to your account and are stored on the server, so they are available in every browser.

### Clipboard panel

Right panel → **Clipboard**: the last 20 texts you copied from terminals (and copied paths).

- Click an entry to copy it again. **Paste** sends it to the focused terminal.
- The history is kept per workspace session in your browser (local storage), not on the server. Use **Clear** to empty it.

### Sharing: members, roles & links

Share connections (or whole folders) with other WRM users or with people **without an account** — the recipients use them **with your stored credentials without ever seeing them**.

Right-click a connection → **🤝 Share**, a folder → **Share Folder**, select several connections and use **🤝 Share**, or *Shares → + New share*. The share dialog has:

- **Connections in this share:** tick connections and/or whole folders (a folder also includes connections added to it later). A share without connections is still a room for chat, voice and terminal sharing.
- **Who can open the link:**

  | Mode | Who |
  |---|---|
  | 👥 **Members only** (default) | Only the users you add as members. They sign in with their account. |
  | 🏢 **Everyone signed in** | Any user of this WRM server who has the link. |
  | 🌐 **Anyone with the link** | Also people without an account (guests). Administrators can disable this mode. |

- **Members** with an individual role each; **role for everyone else** (non-members in the other two modes); optional **password** (members never need it); optional listing under “Shared with me” for all users; **expiry** (1 h … 30 days or a custom date).
- **Roles:**

  | | Observer | Viewer | Operator | Moderator | Owner |
  |---|:-:|:-:|:-:|:-:|:-:|
  | Chat & voice call | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Watch shared terminals | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Browse & download files | — | ✔ | ✔ | ✔ | ✔ |
  | Open terminals | — | — | ✔ | ✔ | ✔ |
  | Upload, edit, rename, delete, transfer | — | — | ✔ | ✔ | ✔ |
  | Share own terminal, receive keyboard control | — | — | ✔ | ✔ | ✔ |
  | Change roles, mute, remove, ban people | — | — | — | ✔ | ✔ |
  | Share settings & members | — | — | — | — | ✔ |

  Roles are **enforced by the server** on every request; the UI only hides what would be refused anyway.
- **🤝 Shares** (top bar) lists your shares with mode, members, expiry and **how many people are online**. Actions: **Open room**, **Copy link**, **Edit** (everything above), **People** (everyone who ever opened the share, with last seen and IP — change a person's role or **ban** them), **New link** (the old link stops working immediately), **Pause/Resume**, **Delete**.
- **Revocation is immediate:** deleting, pausing or rotating a share, removing a member, lowering a role, banning, or disabling a user closes the affected terminals and room connections at once.

### Real-time collaboration

Everyone who opens the **same share link** (including you as the owner) joins its live room. The **collaboration bar** shows the share, your role, avatars of the people online (a green ring shows who is speaking) and the controls.

- **Identities are decided by the server:** signed-in users appear with their account name; guests choose a name on first visit and are marked **guest** (they cannot pose as a registered user).
- **👥 People** panel: who is in the voice call and who is online, with roles, mute/deafen state, shared terminals, raised hands, connection quality (P2P/relay, latency) and a per-person volume slider. Moderators get a **⋯** menu: *change role*, *mute for everyone*, *stop their terminal sharing*, *lower hand*, *remove from room*, *ban from share*.
- **💬 Chat** with **history** (last messages are kept, also for people who join later), timestamps, clickable links, unread badge and notifications. **📎 File exchange** up to the policy limit (images are previewed; other files are always downloaded, never opened in the page).
- **✋ Raise hand.**
- **📡 Terminal sharing:** press 📡 in the title bar of an SSH window. Several terminals can be shared at once. Others click **👁 Watch**: the viewer gets a **snapshot of the current screen** (with colors) and then the live output, in the sharer's exact size. Terminal data is only sent to people who watch it.
- **🎮 Remote keyboard control:** a viewer clicks **Request control**; the sharer sees *Allow / Deny*. The sharer can also give control directly from the People panel and revoke it any time. Keystrokes go only to the terminal that was granted and only while it is shared. Only roles with *Operator* or higher can receive control.
- The room reconnects automatically after network interruptions and restores watching, sharing and the voice call.

### Voice calls

Join the call with **🎙 Join voice** in the collaboration bar. It works like Discord or Jitsi:

- **Mute** (🎤, **Ctrl+Shift+M**), **deafen** (🎧, **Ctrl+Shift+D**), **leave** (📞). Moderators can mute someone for everyone.
- **Push-to-talk** with a key of your choice (e.g. F8), or voice activity.
- **Microphone and speaker selection**, live **input level meter**, test sound, **noise suppression**, **echo cancellation** and **automatic gain control** (*Settings → Voice & audio*, or ⚙ in the call).
- **Speaking indicators**, **per-person volume**, join/leave sounds, connection quality per person (good/fair/poor, P2P or relay, round-trip time).
- **Listen-only:** without a microphone (or permission) you can still join and listen.
- **Robust networking:** peer-to-peer WebRTC with automatic ICE restarts when the network changes, automatic rejoin after reconnects, and the **built-in TURN relay** (UDP and TCP) for networks where direct connections are blocked. External STUN/TURN servers (e.g. coturn on port 443/TLS for very strict networks) can be added in *Admin → Voice & network*.
- **Private:** audio is end-to-end encrypted (DTLS-SRTP); the relay cannot decrypt it.
- Browsers allow the microphone only on **HTTPS** pages (or `localhost`). Up to about 12 people per call work well (mesh); the limit is a policy.

### Audit log, session recording & file transfers

WRM keeps a complete, tamper-evident record of who did what, where and when:

- **Audit log** (*Admin panel → Audit log*): sign-ins (also failed), account and policy changes, connections, shares, room joins and moderation, terminal sessions, file transfers and changes, host keys, exports, viewing of recordings. Entries are linked to the **connection** and the **terminal session** (click **▶ #id** to open the session). Filter by text, event type, user and date range; **CSV export**; **Verify integrity**.
- **Append-only and tamper-evident:** the database refuses to change audit entries, file transfer records and recordings, and to delete anything younger than 7 days. Old entries are removed only by the retention job (`audit_retention_days`, `recording_retention_days`). Every entry contains the SHA-256 of the previous one (**hash chain**), so a changed or removed entry is detected by *Verify integrity*. Every entry is also written as an `AUDIT …` line to the server log (journald/syslog → SIEM).
- **Terminal sessions** (*Admin panel → Sessions & recordings*; every user sees their own in *Settings → Session history*): who, from which IP, which server and remote user, when, how long, how it ended (closed, connection lost, ended by an administrator, interrupted by a server restart).
- **Session recording:** the output of every SSH terminal is recorded in the open **asciinema** format (asciicast v2), compressed, with its SHA-256 stored in the database. **▶ Replay** in the browser — play/pause, seek, speed 0.5–16×, *skip idle time* — or **⬇ download** the `.cast` file and play it with `asciinema play`. The user sees *“This session is recorded”* and a **● REC** badge.
  - **No passwords in recordings:** servers do not echo passwords, so they never appear in the output. Keystrokes are recorded only when an administrator turns on *Record keystrokes*, and typing at password, passphrase and PIN prompts is masked.
  - Recording never slows the terminal (written in the background); a size limit per session stops recording (not the terminal).
  - **Who can see recordings:** administrators all; users their own sessions and sessions of other people on **their** connections (e.g. guests of a share). Viewing or downloading a recording is itself recorded in the audit log.
- **File transfers** (*Admin panel → File transfers*): every upload, download (including files opened in the editor and every file of a ZIP download) and server-to-server copy with source, destination, size, **SHA-256** and status; CSV export.

### Settings

**⚙ Settings** in the top bar has tabs:

| Tab | What it contains |
|---|---|
| General | Language (English / Hrvatski), accent color, close-window confirmation, **notifications when a server goes down**, browser-side upload limit |
| Terminal | Font size, scrollback, auto-reconnect |
| Security | Change password, two-factor authentication (enable/disable, new recovery codes), signed-in devices |
| Voice & audio | Microphone, speaker, level meter, noise suppression, echo cancellation, gain control, input mode / push-to-talk key, call sounds |
| Session history | Your terminal sessions and sessions on your connections, with replay and `.cast` download |
| Data | Export connections (without secrets), export **with** passwords & keys (asks for your password; policy), import (also jump hosts and tunnels), **import from mRemoteNG** and **OpenSSH config** |

Personal preferences are stored per browser; everything security-related is stored on the server.

### Admin panel

*Settings → 🛡 Admin panel* (administrators only):

- **Overview:** version, uptime, users (admins, 2FA), HTTPS, voice relay status, encryption key source, **security warnings** (no HTTPS, open registration, admins without 2FA, host keys off…), **open terminals** (who, which server, from which IP — with *End*), **active SSH tunnels** (owner, connection and route, listen → target, traffic — with *Stop*), **live rooms** (participants, voice, shared terminals).
- **Users:** create users (temporary password generated if you leave it empty), display name, make/remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete (with everything the user owns). The last active administrator cannot be removed.
- **Shares:** every share on the server — pause/resume or delete.
- **Security policies:** registration, required 2FA, password length, lockout, session idle/maximum time, host-key policy, server key files, server-side upload limit, secret export, guest links, chat file size, audit log on/off and retention, session recording (on/off, keystrokes, retention, size limit), **broadcast input**, **live status** (on/off, interval, checks through jump hosts), **SSH tunnels** (on/off, who may use them, listening on network addresses, remote forwarding, idle stop).
- **Voice & network:** voice on/off, participants per call, built-in TURN relay (port, public IP, host name, relay ports, private networks), additional STUN/TURN servers.
- **Host keys:** remembered SSH host keys and FTPS certificates; forget an entry after a server was reinstalled.
- **Audit log:** searchable and filterable (event type, user, date range), linked to sessions, **CSV export**, **Verify integrity** (hash chain).
- **Sessions & recordings:** every terminal session, with replay and `.cast` download.
- **File transfers:** every transferred file with size and SHA-256, CSV export.

Any policy can also be **forced by an environment variable** (`WRM_<KEY>`), e.g. for configuration management; it is then shown locked in the admin panel.

### Mobile & responsive UI

- The top bar and the collaboration bar adapt to the available width (labels collapse into icons).
- On phones (< 760 px): the connection list and the Sessions/Clipboard panel become **slide-in drawers** (☰ and 🗂 buttons), windows open **full screen**, and a single tap opens a connection. Chat and People open as full-width panels.

---

## Keyboard shortcuts

| Where | Shortcut | Action |
|---|---|---|
| Anywhere | **Alt+N** | New window |
| Anywhere | **Esc** | Close menus / dialogs / panels |
| Voice call | **Ctrl+Shift+M** / **Ctrl+Shift+D** | Mute / deafen (also while a terminal has focus) |
| Voice call | your push-to-talk key (e.g. **F8**) | Hold to talk |
| Terminal | **Ctrl+Shift+F** | Search in terminal output |
| Terminal | **Ctrl+Shift+Space** | Snippet picker (Enter run, Shift+Enter insert) |
| Snippets panel | **Click** / **Shift+click** | Run in the focused terminal / insert without Enter |
| Terminal | **Ctrl+Shift+C** / select text | Copy |
| Terminal | **Ctrl+Shift+R** | Reconnect SSH |
| Terminal | **Enter** (when closed) | Reconnect |
| Terminal | **Ctrl+W** | Sent to the terminal (browser tab is not closed) |
| File manager | **Enter / Backspace** | Open / go up |
| File manager | **↑ ↓ Home End**, **Shift** | Move / extend selection |
| File manager | **Ctrl+A**, **Delete**, **F2**, **F5** | Select all, delete, rename, refresh |
| File manager | **Ctrl+F**, **Enter** in search | Filter, search subfolders |
| Editor | **Ctrl+S**, **Tab** | Save, insert tab |
| Window title | **Double-click** | Maximize / restore |
| Tab | **Middle-click**, **right-click** | Close, tab menu |

---

## Configuration

**Environment variables**

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Listening port |
| `LISTEN_ADDR` | `:PORT` | Full listen address, e.g. `127.0.0.1:8080` (overrides `PORT`) |
| `DB_PATH` | `./remote_manager.db` | SQLite database file |
| `ENCRYPTION_KEY` | – | Your own key for stored secrets (first 32 bytes are used; same derivation as v9). |
| `ENCRYPTION_KEY_FILE` | – | Read the key from a file (64 hex characters = raw 256-bit key, anything else is hashed). |
| *(neither set)* | `<DB_PATH>.key` | A random key is generated there on first start (mode `0600`). |
| `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE` | – | Serve HTTPS with this certificate and key (HSTS is sent). |
| `HTTPS_SELF_SIGNED` | – | `1` = create and use a self-signed certificate (`wrm-selfsigned.crt/.key` next to the database). Browsers warn once; good for LANs and voice calls without a CA. |
| `HTTPS_SELF_SIGNED_HOSTS` | – | Extra host names/IPs for the self-signed certificate, comma-separated. |
| `WRM_RECORDINGS_DIR` | `<DB dir>/recordings` | Where session recordings are stored |
| `WRM_TRUST_PROXY` | – | `1` = trust `X-Real-IP` / `X-Forwarded-For` / `X-Forwarded-Proto` / `X-Forwarded-Host` from your reverse proxy (client IPs in rate limits and audit, secure cookies). Only when WRM is reachable **only** through that proxy. |
| `WRM_ALLOWED_ORIGINS` | – | Extra hostnames allowed as `Origin` (comma-separated), e.g. when a proxy rewrites the `Host` header |
| `WRM_ALLOW_ANY_ORIGIN` | – | `1` disables the origin check (not recommended) |
| `WRM_<POLICY>` | – | Forces a policy from the admin panel, e.g. `WRM_REGISTRATION=open`, `WRM_REQUIRE_2FA=all`, `WRM_TURN_PUBLIC_IP=203.0.113.10`, `WRM_TURN_ENABLED=0`, `WRM_HOST_KEY_POLICY=strict`. |

**Policies** (*Admin panel*, or `WRM_<KEY>` in upper case)

| Key | Default | Meaning |
|---|---|---|
| `registration` | `closed` | `closed` / `open` self-registration |
| `require_2fa` | `off` | `off` / `admins` / `all` |
| `password_min_length` | `8` | 6–128 |
| `login_max_failures` | `10` | account lock (15 min) after this many wrong passwords |
| `session_idle_hours` | `168` | sign out after inactivity |
| `session_max_days` | `30` | maximum session age |
| `host_key_policy` | `tofu` | `tofu` / `strict` / `off` |
| `allow_server_keys` | `0` | allow *Key file* / *Auto* for non-admins |
| `max_upload_mb` | `0` | server-side upload limit per file (0 = unlimited) |
| `allow_secret_export` | `1` | allow users to export their secrets (after re-entering the password) |
| `allow_link_shares` | `1` | allow “anyone with the link” shares for guests |
| `chat_file_max_mb` | `5` | file size in collaboration chat (0 = off) |
| `voice_enabled` | `1` | voice calls on/off |
| `voice_max_participants` | `12` | people per call |
| `turn_enabled` | `1` | built-in TURN relay |
| `turn_port` | `3478` | relay port (UDP + TCP) |
| `turn_public_ip` | auto | public IPv4 of the server when behind NAT |
| `turn_host` | WRM host | host name browsers use to reach the relay |
| `turn_relay_ports` | `49152-65535` | UDP port range for relayed audio |
| `turn_allow_private` | `0` | allow relaying to private/loopback networks |
| `ice_servers` | Google + Cloudflare STUN | additional STUN/TURN servers (JSON) |
| `audit_enabled` | `1` | record security events (administrator actions are always recorded). Also `AUDIT_ENABLED` |
| `audit_retention_days` | `365` | how long audit log, file transfers and chat history are kept (minimum 7) |
| `session_recording` | `1` | record terminal sessions. Also `SESSION_RECORDING_ENABLED` |
| `session_recording_input` | `0` | also record keystrokes (masked at password prompts) |
| `recording_retention_days` | `90` | how long recordings are kept (minimum 7) |
| `recording_max_mb` | `100` | maximum recording size per session (0 = unlimited) |
| `tunnels_enabled` | `1` | SSH tunnels and web interfaces behind jump hosts; off stops running tunnels. Also `TUNNELS_ENABLED` |
| `tunnel_users` | `all` | `all` / `admins`: who may use tunnels |
| `tunnel_bind_any` | `0` | allow non-administrators to listen on network addresses (`0.0.0.0`, LAN IP) instead of `127.0.0.1` |
| `tunnel_remote_forward` | `admins` | `off` / `admins` / `all`: remote port forwarding (`-R`) |
| `tunnel_idle_minutes` | `0` | stop tunnels started by hand after this many minutes without traffic (0 = never) |
| `broadcast_enabled` | `1` | broadcast input (typing into several terminals at once). Also `BROADCAST_ENABLED` |
| `status_enabled` | `1` | live up/down status of connections. Also `STATUS_ENABLED` |
| `status_interval_seconds` | `60` | seconds between status rounds (15–3600) |
| `status_jump_checks` | `0` | check connections behind jump hosts by logging in to the jump host |

**Command line**

```
wrm -version                       print the version
wrm -reset-password USER           set a new temporary password (must be changed at sign-in)
wrm -reset-password USER -reset-2fa  … and turn off two-factor authentication
wrm -healthcheck                   check /healthz of the local server (exit code 0 = healthy; for containers)
```

`GET /healthz` (no sign-in) returns `{"status":"ok","version":…}` while the server and its database work, otherwise HTTP 503. Use it for load balancers and monitoring.

---

## Docker, service, reverse proxy & firewall

### Docker

The repository contains a `Dockerfile` (static binary on Alpine, runs as an unprivileged user, data in the volume `/data`, built-in health check) and a `docker-compose.yml`:

```bash
docker compose up -d                       # build and start
WRM_TURN_PUBLIC_IP=203.0.113.10 docker compose up -d   # with voice relay: the IPv4 browsers reach this host at
```

Images of releases are also published to the GitHub Container Registry: `docker run -d -p 8080:8080 -v wrm-data:/data ghcr.io/vedranius/wrm-pro:latest`.

- `/data` holds the database, the encryption key (`remote_manager.db.key`), recordings and the self-signed certificate. Back up the volume and keep a **separate** copy of the key — or provide your own key as a Docker secret (`ENCRYPTION_KEY_FILE=/run/secrets/…`, see the comments in `docker-compose.yml`).
- The compose file enables `HTTPS_SELF_SIGNED=1`. Behind a TLS reverse proxy remove it and set `WRM_TRUST_PROXY=1`.
- Voice relay from a container needs `WRM_TURN_PUBLIC_IP` and the published relay ports (`49160-49200/udp`, matching `WRM_TURN_RELAY_PORTS`).

### Service

**systemd (Linux)**, e.g. `/etc/systemd/system/wrm.service`:

```ini
[Unit]
Description=Web Remote Manager PRO
After=network-online.target

[Service]
User=wrm
WorkingDirectory=/opt/wrm
Environment=LISTEN_ADDR=127.0.0.1:8080
Environment=WRM_TRUST_PROXY=1
Environment=ENCRYPTION_KEY_FILE=/etc/wrm/encryption.key
ExecStart=/opt/wrm/wrm-pro-v10.3.0-linux-amd64
Restart=on-failure
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=/opt/wrm
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload && sudo systemctl enable --now wrm
```

**nginx** in front of WRM (TLS, WebSockets, streaming and big uploads):

```nginx
server {
    listen 443 ssl http2;
    server_name wrm.example.com;
    # ssl_certificate ...; ssl_certificate_key ...;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;                # required: origin check
        proxy_set_header X-Real-IP $remote_addr;    # with WRM_TRUST_PROXY=1
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 3600s;                   # long-lived terminals
        proxy_buffering off;                        # live transfer/search progress
        client_max_body_size 0;                     # large uploads
    }
}
```

**Firewall for voice calls:** allow **UDP and TCP 3478** (TURN) and the **UDP relay range** (`49152-65535`, or narrow it with `turn_relay_ports`) to the WRM server. Behind NAT (cloud VM, home router) forward those ports and set `turn_public_ip`. The TURN port is not proxied through nginx. For networks that only allow HTTPS, run a TURN server with TLS on port 443 (e.g. coturn) and add it in *Admin → Voice & network*.

---

## Security model

See also **[SECURITY.md](SECURITY.md)** (how to report vulnerabilities, hardening checklist).

- **Accounts:** bcrypt hashes, optional/required TOTP 2FA with single-use recovery codes, closed self-registration by default, account lockout and per-IP rate limiting, constant-time responses for unknown users, forced password change after an admin reset, disable/delete accounts with immediate sign-out everywhere.
- **Sessions:** 256-bit random tokens in HttpOnly/SameSite cookies (`Secure` over HTTPS), stored hashed; idle and absolute expiry; device list and remote sign-out.
- **Secrets at rest:** AES-256-GCM with a random per-installation key (or your own). Secrets are **never sent to the browser** — not to the owner, not to share recipients. The database and key files are created with mode `0600`.
- **Web security:** CSRF protection (custom request header + Origin check on every state-changing API call), WebSocket origin check, strict **Content-Security-Policy** (no external scripts, `frame-ancestors 'none'`), `X-Frame-Options`, `nosniff`, `Referrer-Policy: no-referrer` (share tokens never leak via Referer), `Permissions-Policy`, HSTS on HTTPS, request-size limits, no directory listings, all assets served locally (no CDN / supply-chain dependency at runtime).
- **SSH & FTPS:** host keys verified (TOFU/strict); FTPS certificates validated against public CAs or pinned. Server-side key files are admin-only by default.
- **Sharing:** four roles enforced server-side on every file operation, terminal and room action; members-only / signed-in / guest modes; signed password cookies (invalidated when the password changes); expiry, pause and link rotation; bans; immediate revocation of open terminals and room connections.
- **Collaboration:** server-assigned identities, per-client send queues (a slow client cannot stall a room), message rate limits and size limits, terminal data only to watchers, keystrokes only to the granting sharer, file downloads never rendered inline.
- **Voice:** WebRTC with DTLS-SRTP; the built-in TURN relay accepts only short-lived HMAC credentials issued to room participants, limits allocations, and refuses private/loopback/link-local peers by default (no pivoting into internal networks).
- **Audit:** sign-ins (also failed), account and policy changes, connections, shares, room joins and moderation, terminal sessions, file transfers (with SHA-256) and changes, host keys, exports, viewing of recordings — stored **append-only** and **hash-chained** in the database (verifiable) and written as `AUDIT …` lines to the server log (journald/syslog → SIEM). Secrets in audit details are redacted.
- **Session recording:** terminal output in asciicast v2, gzip, mode `0600`, SHA-256 in the database; no keystrokes unless enabled (then masked at password prompts); visible to administrators, the session's user and the connection's owner; every view is audited.
- **Snippets & broadcast:** snippets are per user (shared snippets only by administrators, read-only for others, never auto-run); run-on-connect snippets come from the connection's owner and are recorded in the audit log; broadcast input is client-side typing into the user's own terminals, with confirmation of dangerous commands and an audit entry for every start and stop per terminal.
- **Live status:** checks are plain TCP connections from the WRM server (no credentials, except optional checks through jump hosts with the jump host's own saved login); users only see the states of their own connections.
- **Tunnels & jump hosts:** every hop is authenticated and host-key-verified; tunnels belong to the owner of the connection (only the owner starts them; administrators can stop any); listen on `127.0.0.1` by default; ports below 1024 refused; network addresses and remote forwards only for administrators (policies); a tunnel can be turned off globally or limited to administrators; every start/stop/error is audited with the traffic. A tunnel port on a network address is **not** protected by WRM sign-in — treat it like an open port of that machine.
- **Transport:** use HTTPS (certificate, reverse proxy, or `HTTPS_SELF_SIGNED=1`). Without it passwords and terminal traffic between browser and WRM are not encrypted and browsers block the microphone.

---

## Limitations

- No **RDP/VNC** (graphical desktops): WRM is for SSH, SFTP, FTP and FTPS.
- No **SSO/LDAP/SAML** yet (local accounts with 2FA).
- Voice calls are a **mesh**: fine up to about 12 people; larger meetings would need an SFU.
- No SSH **agent forwarding**. Jump hosts must be SSH servers (no HTTP/SOCKS proxies as jump hosts).
- Tunnel ports open on the **WRM machine**: when WRM runs on a server, a tunnel listening on `127.0.0.1` is not reachable from your PC (see [Where is the tunnel port?](#ssh-tunnels-port-forwarding)).
- Server-to-server transfer and ZIP download work with SSH/SFTP servers only (not FTP). *File contents* search needs `grep` on the server.
- The built-in TURN relay is IPv4 and not available in the Android build (use an external TURN server there).

Ideas and pull requests for any of these are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## API reference

All endpoints (except sign-in, `/api/auth/config`, version and share pages) need the session cookie. JSON in, JSON out, unless noted. **Every non-GET `/api` request must send `X-WRM-Request: 1`** (CSRF protection) and, if present, an `Origin` of this server.

**Accounts**

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/auth/config` | `{first_run, registration_open, version}` |
| POST | `/api/auth/register` | Create account (first user = admin; afterwards only if registration is open) |
| POST | `/api/auth/login` | `{username, password}` → user, or `{mfa_required, mfa_token}` |
| POST | `/api/auth/mfa` | `{mfa_token, code}` (TOTP or recovery code) |
| POST | `/api/auth/logout` | Sign out |
| GET | `/api/auth/me` | Current user, restrictions and policies |
| POST | `/api/auth/password` | `{current_password, new_password}` |
| POST | `/api/auth/2fa/setup` · `enable` · `disable` · `recovery` | Two-factor management |
| GET / DELETE | `/api/auth/sessions[/{id}]` | Signed-in devices / sign out one or all others |
| GET | `/api/users?q=` | User directory (for adding members) |

**Connections, folders, sessions**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/connections` | List (without secrets) / create |
| GET / PUT / DELETE | `/api/connections/{id}` | Read / update (empty secret = keep; `clear_password`, `clear_private_key`) / delete |
| POST | `/api/connections/{id}/duplicate` | Copy on the server |
| POST | `/api/connections/bulk` | `{action: "move"\|"delete", ids, folder_id}` |
| POST | `/api/connections/test` | Test a connection (reports host key problems) |
| POST | `/api/hostkeys/accept` | `{host, fingerprint}` accept a new/changed host key |
| GET / POST | `/api/folders` · DELETE `/api/folders/{id}` | Folders |
| GET / POST | `/api/sessions` · PUT / DELETE `/api/sessions/{id}` | Workspace sessions |
| GET | `/api/config/export` | Export without secrets (with jump hosts and tunnels) |
| POST | `/api/config/export` | `{password}` export with secrets |
| POST | `/api/config/import` | Import → `{imported, skipped, tunnels}` |
| POST | `/api/config/import/mremoteng` | `{xml, password}` import an mRemoteNG `confCons.xml` → `{imported, folders, jump_hosts, tunnels, skipped, notes}`; HTTP 400 with `need_password: true` when a master password is needed |
| POST | `/api/config/import/sshconfig` | `{text, folder}` import an OpenSSH config |

Connections have `jump_id` (jump host connection id or `null`) and `web_path`; lists also return `route` (e.g. `"vpn-gw → bastion-dc1"`) and `tunnels` (number of configured tunnels).

**Snippets & status**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/snippets` | My snippets and shared snippets / create `{name, command, description, group, scope: "all"\|"folder"\|"connection", scope_id, auto_run, shared}` |
| PUT / DELETE | `/api/snippets/{id}` | Update / delete (owner; administrators also shared snippets) |
| GET | `/api/status` | `{enabled, interval, jump_checks, connections: {id: {state: "up"\|"down"\|"unknown", latency_ms, since, checked_at, error, banner, via, via_state}}}` |
| POST | `/api/status/check` | `{ids: [...]}` check now, also through jump hosts (once every 3 s per user) |

Connections also have `monitor` (live status on/off). The terminal WebSocket accepts `{"type":"broadcast","on":true\|false,"peers":n}` (audit) and sends `{"type":"autorun","snippets":[…]}` after run-on-connect snippets ran.

**Tunnels & web interfaces**

| Method | Endpoint | Description |
|---|---|---|
| GET / PUT | `/api/connections/{id}/tunnels` | Tunnels of a connection / replace them: `[{id?, name, kind: "local"\|"remote"\|"dynamic", bind_host, bind_port, target_host, target_port, open_scheme, open_path, start_mode: "manual"\|"connect"\|"always"}]` |
| GET | `/api/tunnels[?all=1]` | `{tunnels: [...], policy}` my tunnels with state (`running`, `state`, `listen`, `active`, `total`, `bytes_up`, `bytes_down`, `started_by`, …); administrators: `all=1` → every running tunnel |
| POST | `/api/tunnels/{key}/start` · `…/stop` | Start (owner) / stop (owner or administrator) |
| POST | `/api/connections/{id}/open-web` | Web interface connection → `{url, direct, route}` (opens a temporary tunnel when it has a jump host) |

**Remote files** (`id` = connection id; for shared connections add `share_token`; the role must allow the action)

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/remote/list?id&path` | List a folder |
| GET | `/api/remote/list-recursive?id&path` | All files below a path |
| GET | `/api/remote/search?id&path&q&mode=name\|content&max&hidden` | Search (NDJSON stream) |
| GET | `/api/remote/download?id&path` | Download a file |
| GET | `/api/remote/download-dir?id&path` | Download a folder as ZIP (SFTP) |
| POST | `/api/remote/upload?id&path&overwrite` | Multipart upload, streamed |
| POST | `/api/remote/mkdir?id&path` | Create folder |
| POST | `/api/remote/rename?id&old&new` | Rename / move |
| DELETE | `/api/remote/delete?id&path&dir=0\|1` | Delete file or folder |
| POST | `/api/remote/transfer` | Server-to-server transfer (NDJSON progress) |

**Sharing**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/shares` | My shares / create `{name, connection_ids, folder_ids, access_mode, guest_role, members:[{username, role}], password, public, expires_at}` |
| GET / PUT / DELETE | `/api/shares/{id}` | Details / update (also `active`, `clear_password`) / delete |
| POST | `/api/shares/{id}/rotate` | New link |
| GET / POST | `/api/shares/{id}/members` · PUT / DELETE `…/members/{userID}` | Members and roles |
| GET | `/api/shares/{id}/participants` · PUT `…/participants/{pid}` | People history; `{role, banned}` |
| GET | `/api/share/{token}` | Share info, my role and permissions, connections |
| POST | `/api/share/{token}/unlock` · `…/profile` | Password unlock · guest name |
| GET | `/api/shared` | “Shared with me” |

**Administration** (administrators)

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/admin/users` | Users / create |
| PUT / DELETE | `/api/admin/users/{id}` | `{is_admin, disabled, display_name}` / delete |
| POST | `/api/admin/users/{id}/reset-password` · `reset-2fa` · `revoke-sessions` · `unlock` | Account actions |
| GET / PUT | `/api/admin/settings` | Policies |
| GET | `/api/admin/status` | Status, open terminals, rooms, relay |
| POST | `/api/admin/terminals/{id}/kill` | End a terminal |
| GET | `/api/admin/shares` | All shares |
| GET / DELETE | `/api/admin/known-hosts[/{id}]` | Host keys |
| GET | `/api/admin/audit?q&action&user&conn&session&from&to&before_id&limit[&format=csv]` | Audit log |
| GET | `/api/admin/audit/verify` | Verify the hash chain `{ok, checked, first_id, last_id, unsigned, broken_at, reason}` |
| GET | `/api/admin/transfers?q&user&conn&direction&from&to&before_id&limit[&format=csv]` | File transfers |
| GET | `/api/version` | Server version |
| GET | `/healthz` | Health check (no sign-in) |

**Terminal sessions & recordings** (administrators: all; users: own sessions and sessions on their connections)

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/recordings?mine&q&user&conn&host&status&from&to&before_id&limit` | Terminal sessions with recording info |
| GET | `/api/recordings/{id}` | One session and its audit entries |
| GET | `/api/recordings/{id}/cast[?download=1]` | The recording as asciicast v2 (`application/x-asciicast`) |

**WebSockets**

| Endpoint | Description |
|---|---|
| `/ws/ssh?id&cols&rows[&share_token]` | Terminal. Binary frames = terminal data; text frames = JSON control (`resize`, `pause`, `resume`, `ping` → server; `status` (with `recording`, `session_id`), `error`, `exit`, `hostkey` ← server). Close codes: 1000 shell exited, 4001 connect failed / not allowed, 4002 SSH connection lost, 4003 access revoked |
| `/ws/events` | Live notifications for the signed-in user (`sessions_changed`, `tunnels_changed`, `snippets_changed`, `status_changed`) |
| `/ws/share/{token}` | Collaboration room. Client → server: `chat`, `file`, `set-name`, `screen` (`start/stop/data/resize/snapshot`), `watch`, `unwatch`, `control-request/grant/revoke/release`, `remote-input`, `voice-join/leave/signal/state`, `hand`, `mod` (`set-role/mute/unmute/stop-share/lower-hand/kick/ban`). Server → client: `welcome`, `participants`, `chat`, `file`, `system`, `screen`, `watch-request`, `control`, `control-request`, `remote-input`, `voice-peers/joined/left/signal`, `force-mute`, `role`, `kicked`, `closed`, `error` |

---

## Building from source & releases

Requirements: **Go** (version in `remote-manager/go.mod`). No C compiler is needed (pure-Go SQLite).

```bash
git clone https://github.com/vedranius/web-browser-RDM-public.git
cd web-browser-RDM-public/remote-manager
go test ./...
go build -o wrm-server .
./wrm-server
```

Cross-compile any platform, e.g. `GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o wrm-linux-arm64 .`, or build all 13 targets with `bash remote-manager/build-all.sh`.

**Project structure**

```
remote-manager/
  main.go         HTTP server, routes, database schema & migrations, connections, folders, sessions, transfer
  security.go     server secret, encryption key management, security headers, CSRF, self-signed TLS, validation
  users.go        sign-in, 2FA step, sessions, password change, account administration
  totp.go         TOTP (RFC 6238), QR code, recovery codes
  settings.go     admin policies (database + WRM_<KEY> environment overrides)
  audit.go        audit log (hash chain, append-only triggers), CSV export, integrity check, retention
  recording.go    terminal sessions, asciicast recorder, file transfer log, recordings & transfers API
  hostkeys.go     SSH host key verification (TOFU/strict), FTPS certificate pinning
  shares.go       sharing model: access modes, roles, members, participants, authorization
  collab.go       collaboration rooms: presence, chat, terminal sharing, control, voice signalling, moderation
  turn.go         built-in TURN relay for voice (turn_android.go: stub for Android)
  admin.go        open-terminal registry, admin status
  ssh_ws.go       SSH terminal WebSocket (binary frames, flow control, keepalive, reconnect codes)
  files.go        file API: list, download, ZIP, streamed upload, mkdir, rename, delete
  search.go       recursive name / content search (NDJSON stream)
  sftp_pool.go    pooled SFTP connections with liveness checks
  conntest.go     "Test connection" endpoint
  jump.go         jump hosts: chains, dialing through hops, validation
  tunnels.go      SSH tunnels (local, remote, SOCKS5), manager, reconnects, tunnels API, web interfaces
  importers.go    import from mRemoteNG confCons.xml and OpenSSH config
  snippets.go     snippets (saved commands), variables, run on connect
  status.go       live up/down status monitor and API
  helpers.go      origin check, login rate limiting, shared helpers
  security_test.go  unit tests (TOTP, roles, share access, CSRF, encryption migration, settings)
  upgrade_test.go   unit tests (upgrading databases from earlier builds, old URLs)
  integration_test.go  in-process SSH/SFTP server: connect → audit → recording, transfers, append-only audit
  tunnels_test.go      jump host chains, local/remote/SOCKS tunnels, policies, web interfaces
  snippets_test.go     snippets API, sharing, variables, run on connect, broadcast audit
  status_test.go       live status: up/down, banners, jump hosts, check now
  importers_test.go    mRemoteNG (encryption formats, inheritance, master password) and OpenSSH config import
  static/index.html          the entire web UI (embedded into the binary)
  static/brand/              logo, favicon and app icons
  static/vendor/             xterm.js + addons and fonts (served locally, see THIRD-PARTY-LICENSES.txt)
.github/workflows/build.yml  CI: vet, gofmt, tests (race), JS syntax check, builds all platforms, Docker image, releases on tags
Dockerfile, docker-compose.yml  container image and one-command deployment
ARCHITECTURE.md              how WRM is built; extension plan status and touch points
docs/brand/                  logo kit (SVG + PNG: mark, lockup, app icons, favicon, social preview)
```

**CI/CD:** every push builds all 13 targets and the Docker image. Pushing a tag `v*` creates a GitHub Release with all binaries, a macOS universal binary and `SHA256SUMS.txt`, and publishes the image `ghcr.io/vedranius/wrm-pro`.

---

## Troubleshooting

| Problem | Fix |
|---|---|
| *Request blocked by CSRF protection* / terminals do not connect behind a proxy | Pass the `Host` header (`proxy_set_header Host $host;`) and WebSocket upgrade headers, or set `WRM_ALLOWED_ORIGINS` |
| *HOST KEY VERIFICATION FAILED* | The server's SSH key changed. If you know why (reinstall), click *Trust this key* in the terminal (connection owner/admin) or forget the host in *Admin → Host keys* |
| *Microphone requires HTTPS* / no microphone | Use HTTPS (certificate, reverse proxy or `HTTPS_SELF_SIGNED=1`) and allow the microphone in the browser; you can also join as a listener |
| Voice: *Connecting…* forever or *Connection problem* | Open UDP+TCP 3478 and the relay port range on the WRM server; behind NAT set `turn_public_ip` and forward the ports; in very strict networks add a TURN server on 443/TLS (*Admin → Voice & network*) |
| Registration tab missing | Self-registration is closed (default). An administrator creates accounts, or opens registration in *Security policies* |
| Locked out as administrator | `wrm -reset-password <user>` (add `-reset-2fa` if needed) |
| *Connection failed: unable to authenticate* | Check user/password/key with **Test connection**; for *Key file*/*Auto* the key must exist on the **WRM server** and the account must be an administrator (or the policy allows it) |
| No *● REC* badge / no recordings | *Admin panel → Security policies → Session recording* is off, or forced off by `WRM_SESSION_RECORDING` / `SESSION_RECORDING_ENABLED`. Check the server log for `recording:` errors (e.g. the recordings folder is not writable) |
| Recordings take too much disk space | Lower `recording_retention_days` or `recording_max_mb`, or move them with `WRM_RECORDINGS_DIR` |
| *Verify integrity* reports a broken chain | An audit entry was changed or removed outside WRM (directly in the database). Restore the database from a backup and investigate who had access to the server |
| Stored passwords stopped working after a restart | The encryption key changed (`ENCRYPTION_KEY`, `ENCRYPTION_KEY_FILE` or the `.key` file next to the database). Restore it, or re-enter the passwords |
| *bastion cannot reach 10.0.0.21:22* (or *jump host "…": …*) | The jump host cannot open a connection to the next hop: check the host:port as seen **from the jump host**, its firewall, and `AllowTcpForwarding yes` in its `sshd_config` |
| Tunnel shows 🔴 *Error* / *administratively prohibited* | The SSH server refuses forwarding (`AllowTcpForwarding`, `PermitOpen`, `GatewayPorts` for remote listeners on other addresses than loopback), or the target is not reachable from the server |
| *address already in use* when starting a tunnel | Another program (or tunnel) uses that port on the WRM machine. Choose another port or leave it empty (automatic) |
| A snippet does not appear in the ⚡ picker | Its scope is a folder or connection other than the terminal's. The panel shows all snippets (others dimmed) |
| Run on connect does not run | It runs only for the owner's snippets that apply to that connection and have no `{{?…}}` prompts; check the audit log for `terminal.auto_run` |
| Broadcast did not ask before a dangerous command | The check sees what you typed during the broadcast, not a command recalled from the shell history (↑) |
| All connections stay grey (no status dot) | Live status is off (`status_enabled`), or the connection has *Monitor up/down status* unticked. The first round starts a few seconds after WRM starts |
| A server shows 🔴 but SSH works | The port in the connection differs from the real one, or a firewall allows SSH only from some addresses (not from the WRM server). Behind a jump host use a jump host instead of a direct connection |
| fail2ban bans the WRM server | Add WRM's IP to `ignoreip`, raise `status_interval_seconds`, or untick monitoring for those hosts |
| Tunnel runs but the page does not open from my PC | The tunnel listens on `127.0.0.1` of the **WRM machine**. See [Where is the tunnel port?](#ssh-tunnels-port-forwarding) |
| No *🔀 Tunnels* button / *not allowed* | Tunnels are off or limited to administrators (*Admin panel → Security policies → SSH tunnels*, or `WRM_TUNNELS_ENABLED` / `WRM_TUNNEL_USERS`) |
| mRemoteNG import: *wrong master password* | Enter the master password the file was protected with in mRemoteNG. Files without a master password are opened automatically |
| *File contents* search fails | The server needs `grep` and a POSIX shell; use *File names* search instead |
| Uploads fail at a certain size | The administrator's upload limit (`max_upload_mb`) or nginx `client_max_body_size` |
| Transfer/search progress appears only at the end | Disable proxy buffering (`proxy_buffering off;`) |

---

## Contributing

Bug reports, ideas, translations and pull requests are welcome! Please read **[CONTRIBUTING.md](CONTRIBUTING.md)**. It explains how to build, test and submit changes, and the contributor terms.

If WRM is useful to you, consider supporting it on **[Ko-fi](https://ko-fi.com/vedranius)** ☕.

---

<sub>Web Remote Manager PRO · © 2026 vedranius · [PolyForm Noncommercial 1.0.0 or PolyForm Internal Use 1.0.0](LICENSE) · hosting, reselling & bundling: [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)</sub>
