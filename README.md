<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/brand/svg/wrm-lockup-on-dark.svg">
    <img src="docs/brand/svg/wrm-lockup-on-light.svg" alt="WRM PRO — Web Remote Manager" width="420">
  </picture>
</p>

# Web Remote Manager PRO (WRM)

**A remote server manager that runs in any web browser.** SSH terminal with **snippets**, **broadcast input** and **live up/down status**, SFTP / FTP / FTPS file manager, **RDP, VNC and Telnet remote desktops** in the browser, **jump hosts** (bastions, chains), **SSH tunnels** (local, remote, SOCKS), **web interfaces** behind jump hosts, **import from mRemoteNG** and `~/.ssh/config`, server-to-server transfers, saved workspaces, sharing with roles, real-time collaboration with **voice calls**, and enterprise security (2FA, policies, a tamper-evident **audit log**, **session recording** with replay, file transfer log): one self-hosted binary (or container) for your PC, server or company.

**Current version: v11.0.0** · [Download](https://github.com/vedranius/web-browser-RDM-public/releases/latest) · [Release notes](RELEASE_NOTES.md) · [Changelog](CHANGELOG.md) · [Security](SECURITY.md) · [Architecture](ARCHITECTURE.md)

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
3. [Upgrading from v10.10.0](#upgrading-from-v10100) · [from v10.9.1](#upgrading-from-v1091) · [from v10.9.0](#upgrading-from-v1090) · [from v10.8.1](#upgrading-from-v1081) · [from v10.8.0](#upgrading-from-v1080) · [from v10.7](#upgrading-from-v107) · [from v10.6](#upgrading-from-v106) · [from v10.5](#upgrading-from-v105) · [from v10.4](#upgrading-from-v104) · [from v10.3](#upgrading-from-v103) · [from v10.2](#upgrading-from-v102) · [from v10.1](#upgrading-from-v101) · [from v10.0](#upgrading-from-v100) · [from v9](#upgrading-from-v9)
4. [How it works](#how-it-works)
5. [Features in detail](#features-in-detail)
   - [Accounts, sign-in & two-factor authentication](#accounts-sign-in--two-factor-authentication)
   - [Connections & folders](#connections--folders)
   - [Quick connect & connection notes](#quick-connect--connection-notes)
   - [SSH keys & credentials vault](#ssh-keys--credentials-vault)
   - [Jump hosts (bastions)](#jump-hosts-bastions)
   - [Proxies (SOCKS / HTTP)](#proxies-socks--http)
   - [SSH tunnels (port forwarding)](#ssh-tunnels-port-forwarding)
   - [Web interfaces (HTTP / HTTPS connections)](#web-interfaces-http--https-connections)
   - [Tags & inventory import (CSV, Excel, NetBox)](#tags--inventory-import-csv-excel-netbox)
   - [Import from mRemoteNG, PuTTY and OpenSSH](#import-from-mremoteng-putty-and-openssh)
   - [Windows, tabs & snapping](#windows-tabs--snapping)
   - [SSH terminal](#ssh-terminal)
   - [Remote desktop: RDP, VNC, Telnet](#remote-desktop-rdp-vnc-telnet)
   - [Out-of-band management, consoles & serial ports](#out-of-band-management-consoles--serial-ports)
   - [Snippets & run on connect](#snippets--run-on-connect)
   - [Folder bookmarks](#folder-bookmarks)
   - [Broadcast input](#broadcast-input)
   - [Live up/down status](#live-updown-status)
   - [Notifications (e-mail, chat, push, webhook)](#notifications-e-mail-chat-push-webhook)
   - [Network tools](#network-tools)
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
   - [Mobile, responsive UI & installable app](#mobile-responsive-ui--installable-app)
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
- **Quick connect** — type `root@10.0.0.5:22` (or `rdp://…`, `vnc://…`) and connect without saving; **notes & runbooks** per connection (procedures, contacts, links); **network tools** — port check, ping, traceroute, DNS and HTTP/TLS check from the WRM server or **from any of your servers** ("can app-01 reach the database on 5432?").
- **Out-of-band management** — power on/off/restart, health and inventory through the server's **BMC** (iDRAC, iLO, XClarity, Supermicro, OpenBMC via **Redfish**, or **IPMI**), and its **serial console** (BMC SSH or IPMI Serial-over-LAN) in a terminal window; **serial ports** of the WRM machine for switch and router consoles.
- **Remote desktops** — **RDP** (Windows), **VNC** and **Telnet** in a browser window, also behind jump hosts, with clipboard, Ctrl+Alt+Del, recording and replay (through guacd, the Apache Guacamole proxy).
- **File manager** for **SFTP** (over SSH), **FTP** and **FTPS**: browse, upload, download, rename, delete, edit, search.
- **SSH keys & credentials vault**: generate SSH keys, put them on many servers at once (like `ssh-copy-id`), see **who has access** to a server; one shared login such as `root@dc1` for many connections, shared with colleagues without revealing it, with **password rotation** on all its servers.
- **Jump hosts** (like `ssh -J`): reach servers behind a bastion, also in chains — for terminals, files, transfers and tunnels.
- **SSH tunnels** (like `ssh -L / -R / -D`, PuTTY, mRemoteNG): reach the web interface of a switch, an iDRAC/iLO or a database behind a server; a SOCKS proxy into a whole management network. **Web interface connections** open such pages with one double-click.
- **Import** your server inventory from **CSV / Excel** files and **NetBox** (with tags for environment, site and rack, and sync), and your connections from **mRemoteNG** (with passwords, folders and SSH tunnels), **PuTTY** (with proxies) and **OpenSSH** `~/.ssh/config`. **Tags** filter the sidebar and mark production servers.
- **Git workspace**: which version of which service runs where — compare installations on your servers **file by file** with **GitLab / GitHub** or an **offline bundle**, see what is *old*, *missing* or *changed by hand*, with diffs and notifications (read-only).
- **Server-to-server copy** between two SSH servers, without downloading to your computer first.
- **Workspaces**: many terminal and file windows side by side, tabs, snapping, saved sessions. **Installable as an app** (PWA) on desktops, tablets and phones.
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
   chmod +x wrm-pro-v11.0.0-linux-amd64
   ./wrm-pro-v11.0.0-linux-amd64
   ```
   On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

   **Windows** (PowerShell), or just double-click the `.exe`:
   ```powershell
   .\wrm-pro-v11.0.0-windows-amd64.exe
   ```

   **Android (Termux)**
   ```bash
   pkg install wget
   wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v11.0.0/wrm-pro-v11.0.0-android-arm64
   chmod +x wrm-pro-v11.0.0-android-arm64 && ./wrm-pro-v11.0.0-android-arm64
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

## Upgrading from v10.10.0

Replace the binary. The database gets new tables (`git_sources`, `git_trees`, `git_blobs`, `git_workspace`, `git_targets`, `git_installs`); the previous binary still starts on it and ignores them.

- **Git workspace:** the new **⎇ Git** button in the top bar compares the services installed on your servers with GitLab / GitHub or an offline bundle, file by file and read-only. See [Git workspace](#git-workspace-compare-services).
- It is on for everyone. Administrators can hide it (`git_enabled` off, or `WRM_GIT_ENABLED=0`) or limit checks to administrators (`git_checks`).
- New notification events *Git: a new version is available*, *Git: files changed by hand on a server* and *Git: a server is unreachable for checks* in *Settings → Notifications*.

## Upgrading from v10.9.1

Replace the binary. The database gets a new table `bookmarks`; the previous binary still starts on it and ignores it.

- **Folder bookmarks:** ★ next to the path bar of the file manager and in a terminal's title bar. *Bookmark this directory…* saves the current folder for one connection, a folder, a tag or all connections; paths may use `~`, `$USER`, `{host}` and `{name}`. A bookmark can be the **start directory** of a connection. See [Folder bookmarks](#folder-bookmarks).
- **WinSCP bookmarks:** *★ → Manage bookmarks… → Import WinSCP.ini* adds the remote directory bookmarks of a WinSCP configuration as global bookmarks.
- Bookmarks are part of *Export* / *Import*.

## Upgrading from v10.9.0

Replace the binary. The database gets a new table `putty_sessions` and a column `connections.putty_session`; the previous binary still starts on it and ignores them.

- **Import from PuTTY:** *Settings → Data → Import from PuTTY* reads a `.reg` export of your saved PuTTY sessions, with their proxies. See [Import from mRemoteNG, PuTTY and OpenSSH](#import-from-mremoteng-putty-and-openssh).
- mRemoteNG connections imported **before** v10.9.1 do not remember their *PuttySession*, and importing the file again skips them as existing: give them the proxy by hand (or through their folder's default proxy).

## Upgrading from v10.8.1

Replace the binary. The database gets new tables (`proxies`, `proxy_grants`, `notify_channels`, `notify_subscriptions`, `notify_prefs`, `notify_pending`, `notify_state`) and columns (`connections.proxy_id`, `folders.jump_conn_id`, `folders.proxy_id`); the previous binary still starts on it and ignores them (connections then go without their proxy and folder defaults).

- **Proxies:** define SOCKS5, SOCKS4/4a or HTTP CONNECT proxies (or a SOCKS tunnel of a connection) under *🔑 SSH keys & credentials → 🌐 Proxies* and choose one in a connection (*Proxy*, next to *Jump host*). See [Proxies](#proxies-socks--http).
- **Folder defaults:** right-click a folder → *Folder settings…* to give it a default jump host and proxy. Connections in the folder that have no jump host of their own then use the folder's (existing connections without a jump host included). Choose *no jump host (not the folder's)* where a connection must stay direct.
- **Notifications:** an administrator adds channels in *Admin panel → Notifications* (e-mail, Telegram, Slack / Mattermost / Rocket.Chat, Teams, Discord, ntfy, Gotify, Pushover, webhook); users subscribe in *Settings → Notifications*. Nothing is sent until both happened. See [Notifications](#notifications-e-mail-chat-push-webhook).
- **Docker:** `docker-compose.yml` now maps `host.docker.internal` to the Docker host (`extra_hosts`), for proxies that run on the host.
- Deleting a jump host now makes the connections that used it fall back to their folder's default jump host (or go direct).

## Upgrading from v10.8.0

Replace the binary. Nothing changes in the database.

- **Larger imports work:** mRemoteNG, OpenSSH config and inventory (CSV / Excel) files up to **20 MB** import again (from about 1 MB on they failed with *Bad JSON*). A larger file gets a clear message with its size and the limit, already in the browser.
- **Install as an app** (*Settings → General*) now shows the steps for your browser, or the reason why installing is not possible here (plain `http://` to an address other than `localhost`, an untrusted certificate, Firefox on the desktop). The *Install* button appears only when the browser can install right away. To get HTTPS quickly, see [HTTPS with Caddy](#https-with-caddy).
- A reverse proxy in front of WRM must allow request bodies of at least 32 MB for imports (`client_max_body_size 0;` in the nginx example does).

## Upgrading from v10.7

Replace the binary. The database gets two new columns (`connections.notes`, `connections.temp_until`); the previous binary still starts on it (it shows quick connections as normal ones).

- **⚡ Quick connect** field at the top of the sidebar: `user@host:port` + Enter.
- **📝 Notes** in the connection dialog and in the right-click menu; a 📝 next to a connection opens them.
- **🧰 Tools** in the top bar, and *Network tools from this server…* / *Check ports of this host* in the right-click menu (policy `network_tools`, default: all users).
- **Install as an app**: *Settings → General*, or the install icon of the browser's address bar (needs HTTPS with a trusted certificate, or `localhost`).

## Upgrading from v10.6

Replace the binary. The database gets one new column (`connections.bmc`); the previous binary still starts on it.

- Connections can have a **BMC** (*Edit → Out-of-band management*); right-click → *Power & status (BMC)…* and *Serial console*. For IPMI install `ipmitool` on the WRM server or on the jump host (`apt install ipmitool`).
- New protocol **Serial port** for console cables on the WRM machine (administrators by default, policy `serial_ports`).
- Saved workspace sessions now really keep remote desktop windows (and console windows).

## Upgrading from v10.5

Replace the binary. The database gets three new columns (`connections.tags`, `ext_id`, `ext_tags`) and a table `inventory_sources`; the previous binary still starts on it (it ignores the tags).

- Connections have **tags**; the sidebar shows a tag filter and an environment badge (`env:prod` …). See [Tags & inventory import](#tags--inventory-import-csv-excel-netbox).
- **Settings → Data → Import inventory**: CSV / Excel files and NetBox.

## Upgrading from v10.4

Replace the binary. The database gets new tables (`ssh_keys`, `ssh_key_deployments`, `credentials`, `credential_grants`) and two new columns (`connections.key_id`, `connections.credential_id`); the previous binary still starts on it (connections that use the new login types then cannot log in).

- New **🔑 Keys** button: your SSH key store and the credentials vault. See [SSH keys & credentials vault](#ssh-keys--credentials-vault).
- Connections can log in with an **SSH key from the key store** or a **credential from the vault**; remote desktops and FTP can use vault credentials too.
- Deleting an account now also deletes its snippets, tunnels, keys and credentials.

## Upgrading from v10.3

Replace the binary. The database gets one new column (`connections.options`); the previous binary still starts on it.

- **Remote desktops need guacd** next to WRM: `apt install guacd` (Debian/Ubuntu), `docker run -d --network host guacamole/guacd`, or the `guacd` service that `docker-compose.yml` now contains. Without guacd everything else keeps working; opening an RDP/VNC/Telnet connection then explains what is missing. See [Remote desktop](#remote-desktop-rdp-vnc-telnet).
- Importing from **mRemoteNG** now also brings RDP, VNC and Telnet connections (with domain, console, colors and RD Gateway).
- SQLite now waits up to 10 s for a busy database on every connection (fixes rare *database is locked* errors under load).

## Upgrading from v10.2

Replace the binary. The database gets a new table (`snippets`) and one new column (`connections.monitor`); the previous binary still starts on it.

- **File names no longer carry a suffix:** releases are `v10.3.0`, binaries `wrm-pro-v10.3.0-linux-amd64` etc. Update scripts and systemd units that use the file name.
- **Live status is on by default:** WRM opens a TCP connection to every saved host:port once a minute (SSH and FTP: it reads the greeting and closes politely, like Nagios `check_ssh`). Servers see it in their logs. Change the interval or turn it off in *Admin panel → Security policies → Live status*, or untick *Monitor up/down status* in a connection. See [Live up/down status](#live-updown-status) about fail2ban.
- **Broadcast input** is available to everybody; turn it off with `broadcast_enabled` if you do not want it. Every start and stop is in the audit log.

## Upgrading from v10.1

Replace the binary. The database gets a new table (`connection_tunnels`) and two new columns on `connections` (`jump_conn_id`, `web_path`); nothing else changes. To go back, run the v10.1 binary on the same database; it does not know these features (connections behind a jump host then try to connect directly, web interface connections cannot be opened, tunnels do not run).

- New: [jump hosts](#jump-hosts-bastions), [SSH tunnels](#ssh-tunnels-port-forwarding), [web interface connections](#web-interfaces-http--https-connections) and [import from mRemoteNG / OpenSSH](#import-from-mremoteng-putty-and-openssh).
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
  - **Protocol:** `SSH` (terminal, plus files over SFTP), `SFTP` (opens the file manager by default), **🖥 RDP**, **🖥 VNC**, **Telnet** (see [Remote desktop](#remote-desktop-rdp-vnc-telnet)), **🔌 Serial port** of the WRM server (see [Serial ports](#out-of-band-management-consoles--serial-ports)), `FTP`, `FTPS` (FTP with explicit TLS; file manager only), **🌐 Web interface** `HTTPS` / `HTTP` (see [Web interfaces](#web-interfaces-http--https-connections)).
  - **Host : Port** (default port 22 for SSH/SFTP, 21 for FTP/FTPS, 443/80 for web interfaces), **Username**, **Folder**, **Tags** (comma separated, e.g. `env:prod, site:dc1, rack:a12, web`; the most used tags are offered below the field).
  - **Jump host (connect via):** another saved SSH connection to go through (see [Jump hosts](#jump-hosts-bastions)).
  - **Proxy:** a saved SOCKS / HTTP proxy for everything the connection does over TCP (see [Proxies](#proxies-socks--http)). With a jump host the proxy is reached through it; the dialog shows the route, e.g. `⤳ Route: bastion → socks5://10.1.1.1:1080 → app-01`.
  - Both show *— as the folder: … —* when the folder has a default (see below); *no jump host / no proxy (not the folder's)* overrides it.
  - **🔀 SSH tunnels:** port forwards of this connection (see [SSH tunnels](#ssh-tunnels-port-forwarding)).
  - **⚡ Out-of-band management:** the BMC of the server (see [Out-of-band management](#out-of-band-management-consoles--serial-ports)).
  - **📝 Notes & runbook:** procedures, contacts, links (see [Notes](#quick-connect--connection-notes)).
  - **Authentication:** *Password*, *Private key (paste)*, *Key file (path on the WRM server)*, *Auto* (tries `~/.ssh/id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa` of the WRM server user), **🔑 SSH key from the key store** or **🗝 Credential from the vault** (see [SSH keys & credentials vault](#ssh-keys--credentials-vault)). *Key file* and *Auto* use keys of the **WRM server itself** and are therefore limited to administrators (policy). Password authentication also answers keyboard-interactive prompts. RDP, VNC, Telnet and FTP log in with a password or a vault credential.
  - **🔌 Test connection** logs in once and shows the result, latency and — for a new server — its host key fingerprint.
- **Secrets never leave the server.** Passwords and private keys are write-only: the edit dialog shows “saved and encrypted — leave empty to keep it”, with an option to remove the saved secret. *Duplicate* copies the connection on the server.
- **Host key verification:** the first connection to a server remembers its SSH host key (and, for FTPS, the certificate if it is not signed by a public CA). A **changed key is refused** with a clear warning and both fingerprints; the owner of the connection (or an administrator) can accept the new key after checking it. Policy: *trust on first use* (default), *strict* (new hosts must be approved) or *off*.
- **Double-click** a connection to open it: SSH opens a terminal; SFTP/FTP/FTPS open the file manager. On touch devices a single tap opens it.
- **Right-click** a connection: *Open as Terminal / File Manager* (in a new or an existing window), *Tunnels*, *Start / Stop all tunnels*, *Add tunnel…*, *Network tools from this server…*, *Who has access (authorized_keys)…*, *Deploy SSH key…*, *Notes & runbook*, *Check ports of this host*, *Edit*, *Duplicate*, *Share*, *Delete*. Web interface connections: *Open web interface*.
- **📁 Folder** creates folders. **Drag** connections onto a folder to move them, or onto *"↓ Drop here"* to move them back to the root. Right-click a folder to *Share* or *Delete* it (its connections move to the root), or **Folder settings…** to rename it and give it a **default jump host** and a **default proxy**: its connections — also new ones and the ones moved into it — use them unless they choose their own (⤳ after the folder name). Export and import keep the defaults.
- **Search box** filters by name, host, username, protocol or tag; `#tag` searches tags only, several words must all match (`#env:prod web`). The **🏷 tag bar** above the list filters by one or more tags (click again to remove). **Multi-select** with **Ctrl/Cmd + click**, then use **🤝 Share**, **🔑 Key** (deploy an SSH key to all of them), **🏷 Tags** (add tags, or remove them with `-tag`) or **Delete**.
- The **dot in front of a connection** is its [live status](#live-updown-status): 🟢 up, 🟡 up but slow (> 300 ms), 🔴 down, a ring = behind a jump host (colored by the jump host's state). A **blue dot after the name** means one of its terminals is connected. **⤳** means it goes through a jump host (hover shows the route), **🔀** that it has tunnels (colored while one runs; click opens them), **WEB** marks web interfaces.
- **Shared with me** lists shares you are a member of (and shares published for all users), with your **role**; ↗ opens the collaboration room.

### Quick connect & connection notes

**⚡ Quick connect** — the field at the top of the sidebar (*⚡ Quick connect: user@host:port*) — connects without filling in the connection dialog:

- Type an address and press **Enter**: `root@10.0.0.5`, `admin@sw1:2222`, `10.0.0.7`, or with a scheme: `ssh://`, `sftp://`, `rdp://admin@10.0.0.20`, `vnc://10.0.0.30:5901`, `telnet://sw1:23`, `https://10.0.0.11/ui`. Without a scheme it is SSH.
- A small dialog asks for the **login** — a password, a **vault credential** or a **key from the key store** — and optionally a **jump host**. Enter connects: a terminal, file manager, remote desktop or web interface opens like for a saved connection (host key check, recording and audit included; audit action `connection.quick`).
- The connection is listed under **⚡ Quick connections** in the sidebar: 💾 (or right-click → *Save as a connection*) opens the connection dialog to give it a name, folder and tags and keep it; ✕ removes it. Quick connections are kept **24 hours after their last use** (an open terminal keeps them), at most 20 per user, are not monitored and not exported.

**📝 Notes & runbook** — every connection has notes: what runs there, how to restart it, who to call, links to the wiki or the ticket:

- *Edit → 📝 Notes & runbook*, or right-click → *Notes & runbook*. A **📝** after the name in the sidebar and in the window bar of its terminals opens them.
- A small markdown: `# headings`, `- lists`, `**bold**`, `` `code` ``, ```` ``` ```` code blocks and links (`https://…`). Everything else is shown as text (no HTML).
- Notes belong to the owner of the connection: share members don't see them. They are exported and imported with the connection and copied by *Duplicate*. Up to 20 000 characters.

### SSH keys & credentials vault

The **🔑 Keys** button in the top bar opens two lists: **SSH keys** and **Credentials**.

**SSH keys (key store)**

- **Generate key:** ED25519 (recommended), RSA 2048/3072/4096 or ECDSA 256/384/521, with a name and a comment. The key is created on the WRM server and stored **encrypted** like all other secrets.
- **Import key:** paste (or choose the file of) an OpenSSH or PEM private key — keys with a passphrase are asked for it once and then stored encrypted with WRM's key. A **public key alone** (e.g. a colleague's `id_ed25519.pub`) can be imported too: it can be deployed to servers, but WRM cannot log in with it. PuTTY `.ppk` keys: export them with PuTTYgen → *Conversions → Export OpenSSH key* first.
- **Use it:** in a connection choose *Authentication → SSH key from the key store*. The key list shows which connections log in with a key, where it is deployed and when it was last used.
- **🚀 Deploy:** puts the public key on the selected SSH connections, like `ssh-copy-id`: WRM logs in with each connection's current login (password, other key or vault credential — through jump hosts too), creates `~/.ssh` if needed and appends the key to `~/.ssh/authorized_keys` once (a key that is already there is not added again; other lines stay untouched; SELinux contexts are restored). With **"Log in with this key from now on"** WRM then logs in with the key to test it and switches the connections to it. Deploy from the key, from a connection (right-click → *Deploy SSH key…*) or for many selected connections (**🔑 Key** in the sidebar).
- **✂ Revoke:** removes the key from `authorized_keys` on the selected servers. A connection that logs in with that key cannot be selected — WRM would lock itself out.
- **📋 Copy public key**, **⬇ .pub** download, **rename**, **delete** (refused while connections or credentials log in with the key — servers keep a deleted key until you revoke it).
- **🔓 Export private key:** after entering your WRM password again, optionally protected with a new passphrase; recorded in the audit log. Allowed when *allow_secret_export* is on (always for administrators).

**Who has access** (right-click an SSH connection → *Who has access (authorized_keys)…*): WRM reads `~/.ssh/authorized_keys` of the user it logs in as and lists every key with type, size, comment, fingerprint and options (`from=…`, `command=…`). Keys of your key store are named, **the key WRM itself logs in with** is marked (and cannot be removed), and every other key can be **removed** with one click. *Add a key…* deploys one there. Keys allowed elsewhere (`AuthorizedKeysFile`, `AuthorizedKeysCommand`, other accounts) are not shown.

**Credentials (vault)**

A credential is a login that many connections share: a user name with a **password**, an **SSH key** of the key store, or both. Typical: `root@dc1`, `Administrator` of the Windows servers, the iDRAC/iLO admin.

- **New credential:** name, user name (optional — without one each connection uses its own user name), password (🎲 generates a strong one), SSH key, description, a **rotation reminder** in days and **Only for hosts**: host names with `*` and `?` or CIDR ranges (`*.dc1.example.com, 10.1.0.0/16`). With a host list the credential works only for matching hosts.
- **Use it:** in a connection choose *Authentication → Credential from the vault*. SSH, SFTP, RDP, VNC, Telnet and FTP connections can use it. Change the credential once — every connection logs in with the new secret.
- **Share with users:** add colleagues (or, administrators, *all users*). They see the credential in their list and can use it in **their own** connections, but they never see or change the secret. Removing a user stops their connections from logging in with it at once. Because a password is sent to the server when logging in, a shared credential should have a host list: then it only works for those hosts — and, for the people it is shared with, only through jump hosts that match the list too. WRM asks before saving a shared credential without a host list.
- **👁 Show** reveals the stored password to its owner after entering the WRM password again (audited; policy *allow_secret_export*). **🔗 Used by** lists the connections that use it (all users' for the owner).

**Password rotation** (🔄 *Rotate…* on a credential, owner only) changes the password on every SSH server that uses the credential:

1. **Check logins** — WRM logs in to every server with the stored password (through jump hosts, in parallel) and changes nothing. Connections of other users that use the credential are included; the same user@host:port is changed only once.
2. **Rotate now** — the same check first: if one server does not accept the stored password, **nothing is changed**. Then, one server after the other, WRM runs `passwd` over SSH (answering its prompts in a terminal), logs in with the new password to verify it, and finally stores it in the vault. The new password is generated (24 characters: upper and lower case letters, digits and symbols) or entered by you.
3. **All or nothing:** when a server fails (password rules, `passwd` not allowed…), the servers already changed get the old password back and the vault keeps it. If even that fails (e.g. a password history rule refuses the old password), the vault keeps the old password and remembers the new one as **pending**: the credential is marked *rotation incomplete*, *Show* reveals both, and the result lists which servers have which password. Set the right password with *Edit* once fixed.

The progress of every server is shown live. Rotation needs SSH connections to servers where the login user may change its own password with `passwd` (Linux, BSD, macOS). Connections of other types (RDP, VNC, FTP…) block rotation — change those passwords yourself and save the new one with *Edit* (*Edit* with a new password counts as rotated). When the rotation reminder is due, the 🔑 button shows a badge and *Admin panel → Overview* warns.

Everything is in the audit log: `ssh_key.created`, `ssh_key.deployed`, `ssh_key.revoked`, `ssh_key.exported`, `ssh_key.deleted`, `credential.created`, `credential.updated`, `credential.granted`, `credential.grant_revoked`, `credential.revealed`, `credential.checked`, `credential.rotated`, `credential.rotation_failed`, `credential.rotation_incomplete` (never the secrets).

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
- Deleting a jump host makes the connections that used it direct again (or use their folder's default jump host).
- A folder can give its connections a **default jump host** (*Folder settings…*); a jump host that is itself in the folder does not go through itself.

### Proxies (SOCKS / HTTP)

Some servers are reached through a SOCKS or HTTP proxy instead of (or after) a bastion, like PuTTY's *Connection → Proxy*. A proxy is defined once and chosen per connection, like a vault credential:

- *🔑 SSH keys & credentials → 🌐 Proxies → New proxy* (or *＋ New proxy…* in the connection dialog's **Proxy** field). Types:
  - **SOCKS5**, with an optional user name and password (RFC 1929) — encrypted, or a **vault credential**;
  - **SOCKS4 / SOCKS4a** (a user id, no password);
  - **HTTP CONNECT** with Basic authentication;
  - **WRM SOCKS tunnel of a connection:** a dynamic tunnel (`-D`) of one of your SSH connections. WRM starts it when a connection needs it and stops it after 30 minutes without traffic.
- **DNS at the proxy** (PuTTY: *Do DNS name lookup at proxy end*, default on): the proxy resolves host names (SOCKS5 domain names, SOCKS4a). Off: WRM resolves them and sends the address.
- **Everything a connection does over TCP goes through it:** SSH terminal, SFTP, search, transfers, tunnels, key deploy and password rotation, FTP/FTPS (also the data connections), RDP / VNC / Telnet (guacd connects to a local relay on `desktop_tunnel_bind`), web interfaces (opened through a local relay on `127.0.0.1`), *Test connection*, live status, network tools from the server, Redfish (when *Through the jump host* is ticked). **IPMI and Serial-over-LAN use UDP** and cannot go through a proxy: WRM refuses that combination.
- **Proxy and jump host together:** the proxy is reached through the jump hosts, then the target:

```
 WRM ──SSH──► bastion ──direct-tcpip──► proxy 10.1.1.1:1080 ──SOCKS CONNECT──► app-01:22
```

  The route is shown in the sidebar, the dialog, the terminal, *Sessions & recordings* and the audit log: `bastion → socks5://10.1.1.1:1080 → app-01`. Jump hosts can have proxies of their own.
- **Clear errors** that tell the proxy and the target apart:
  - `proxy socks5://10.1.1.1:1080: not reachable: connection refused` — the proxy itself;
  - `proxy socks5://…: authentication required` / `authentication failed: wrong user name or password` (HTTP: `407`);
  - `proxy socks5://… cannot reach 10.0.0.5:22: connection refused (SOCKS5 reply 5)` — the proxy works, the target does not answer (HTTP: `502` / `504`).
- **Live status** checks through a *WRM SOCKS tunnel* proxy only while that tunnel runs (it is never started or kept alive by a check); otherwise the state is *unknown*.
- A proxy address on the WRM machine that is a running SSH tunnel of **another user** is refused (*this address is an SSH tunnel of another user on the WRM server*): tunnel listeners have no login of their own.
- **Sharing:** share a proxy with users (or, administrators, with everybody): they can use it but never see its login. A proxy shared **with a password** is not reached through the grantee's own jump hosts (they could read the plain-text proxy login there). WRM tunnel proxies are not shared.
- **Folder default:** *Folder settings…* → *Default proxy*.
- **Policy `proxies`** (`all` / `admins` / `off`): who may define proxies. Everybody can use proxies shared with them.
- Export and import keep proxies by name (the export lists your proxies, passwords only in an export with secrets; an import creates the missing ones); *Duplicate* keeps the proxy.
- **Docker:** inside a container `127.0.0.1` is the container itself. `docker-compose.yml` maps `host.docker.internal` to the Docker host (`extra_hosts: ["host.docker.internal:host-gateway"]`): for `ssh -D 1080` or a proxy on the host use `host.docker.internal:1080`, and let it listen on an address the container reaches (not only `127.0.0.1` of the host). WRM warns when a proxy at `127.0.0.1` is saved while it runs in a container.

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

### Tags & inventory import (CSV, Excel, NetBox)

**Tags** describe connections: free words (`web`, `db`) and tags with a key — `env:` (environment), `site:` (location), `rack:`, `role:`, `tenant:`, `platform:`, `cluster:`. Tags are lower case, without spaces (`Rack A 12` becomes `rack:a-12`), at most 30 per connection.

- Set them in the connection dialog or for many connections at once (**🏷 Tags** after a multi-select; `-tag` removes one).
- The **🏷 tag bar** above the connections lists the most used tags with their counts; click tags to show only connections that have all of them (with *+N* for the rest). The search box finds tags too (`#rack:a12`).
- **Environments stand out:** `env:prod` (also `production`, `prd`, `live`) gets a red **PROD** badge and a red edge, staging/UAT a yellow, test/QA a blue and dev/lab a green badge — so a production server is never mistaken for a test one.
- Tags are kept by *Export* / *Import* and *Duplicate*.

**Import inventory** (*Settings → Data → Import inventory*):

*CSV / Excel*

- Choose a **CSV** file (comma, semicolon, tab or `|` separated — detected automatically; UTF-8, UTF-16 or **Windows-1250** as Excel saves it in Croatian, so `č ć š ž đ` come out right) or an **Excel workbook (.xlsx)** with a sheet selector. Old `.xls` files: save them as `.xlsx` or CSV first.
- WRM reads the first row as column names and **guesses the mapping** (English and Croatian names: *Name/Naziv/Hostname, IP/IP adresa/Address/URL, Port, Protocol/Protokol, User/Korisnik, Password/Lozinka, Folder/Mapa/Grupa, Tags/Oznake, Environment/Okruženje, Site/Location/Lokacija, Rack, Role/Uloga, Platform/OS, Jump host/Bastion, Path*). The preview shows the first rows under a field selector per column — change any column or set it to *ignore*.
- The host column understands `10.0.0.5`, `10.0.0.5:2222`, `root@10.0.0.5`, `10.0.0.5/24`, IPv6 and URLs (`https://idrac-01/` becomes a web interface connection, `ssh://admin@sw1:2222` an SSH connection). A *Platform* containing *Windows* makes an RDP connection. *Environment*, *Site*, *Rack*, *Role* and *Platform* become tags; a *Jump host* column links connections to a jump host by name.

*NetBox*

- Enter the NetBox address and an **API token** (read-only is enough), optionally filters (sites, roles, tenants, a NetBox tag, status), whether to use the **primary IP** or the **name (DNS)** as address, and the name of a custom field that holds the environment. *Load from NetBox* reads **devices** and/or **virtual machines** (all pages) and shows them; objects without an address are marked.
- Sites, racks, roles, tenants, platforms, clusters, the environment and NetBox tags become tags; Windows platforms become RDP connections.
- **Sync:** imported objects remember their NetBox id. Loading again marks them as *imported*; with **Update the ones imported before** their address, NetBox tags and folder follow NetBox (your own tags, the name and the login stay). Objects that disappeared from NetBox are listed — WRM never deletes connections on its own.
- *Remember address and token* stores the address, filters and the token (encrypted) for your account; the token is never shown again.

*For both:* choose the **login** of the imported connections — none (set it later), a **vault credential** (e.g. `root@dc1`) or a **stored SSH key** — a default user name and protocol, **folders** (from a column, one per site / role / tenant, all in one, or none) and extra tags for all. Rows with their own password keep it. Existing connections (same name, host and protocol) are skipped, so importing again is safe. The result lists what was imported, updated and skipped (with reasons).

### Import from mRemoteNG, PuTTY and OpenSSH

*Settings → Data*:

- **Import from mRemoteNG** — the `confCons.xml` file (`%APPDATA%\mRemoteNG\confCons.xml`, or *File → Export* in mRemoteNG):
  - containers become folders (`DC1 / Rack A`); SSH1/SSH2 connections become SSH connections, RDP, VNC and Telnet connections remote desktops (with their options), HTTP/HTTPS connections web interface connections (with the path of the URL);
  - **passwords are decrypted** (current AES-GCM format, older AES-CBC format, *FullFileEncryption*) and stored encrypted in WRM. If the file is protected with a **master password**, WRM asks for it;
  - inherited settings (*Inherit* user name, password, port, domain, SSH tunnel) are resolved from the parent folders;
  - the **SSH tunnel** setting (`SSHTunnelConnectionName`) becomes the **jump host**;
  - ICA, rlogin and other protocols WRM does not open are listed as *skipped*; connections that already exist (same name, host and protocol) are skipped too, so importing again is safe.
- **Import OpenSSH config** — `~/.ssh/config` (folder *SSH config*):
  - `Host` entries with `HostName`, `User`, `Port`; defaults from `Host *` and from wildcard patterns;
  - **`ProxyJump`** (also `user@host:port`; a jump host that is not its own `Host` entry is created) and `ProxyCommand ssh -W …` become the **jump host**;
  - **`LocalForward`, `RemoteForward`, `DynamicForward`** become tunnels that start with the terminal;
  - `IdentityFile` becomes *Key file* authentication (only if key files on the WRM server are allowed for you; otherwise a note asks you to add the key or a password);
  - `Match` blocks and wildcard-only hosts are skipped.
- **Import from PuTTY** — a registry export of the saved sessions (folder *PuTTY*). In `regedit` right-click `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions` → *Export*, or run `reg export HKCU\Software\SimonTatham\PuTTY\Sessions putty.reg`. Both UTF-16 (regedit's default) and UTF-8 files are read:
  - session names are decoded (`My%20Server` → *My Server*); `HostName` (also `user@host`), `PortNumber`, `UserName` and `Protocol` are imported. *Default Settings* is not imported as a connection, but its values are the defaults of the other sessions;
  - **SSH** and **Telnet** sessions become connections; raw, rlogin and serial sessions (a COM port of the Windows PC) are listed as *skipped*;
  - `PublicKeyFile` gives a note: convert the `.ppk` with PuTTYgen (*Conversions → Export OpenSSH key*) and add it under *🔑 SSH keys*;
  - **proxies:** SOCKS4, SOCKS5 and HTTP proxies (`ProxyMethod`, `ProxyHost`, `ProxyPort`, `ProxyUsername`, `ProxyPassword`, `ProxyDNS`) become [saved proxies](#proxies-socks--http) named e.g. *PuTTY socks5 proxy.example.com:1080*, linked to the connections. Sessions with the same proxy share one; a proxy you already have (same type, address, login and DNS setting) is reused, so importing again creates nothing new. PuTTY's *SSH proxy* becomes the **jump host** (by name); Telnet and local-command proxies get a note. Needs the policy `proxies` to allow you to define proxies;
  - **with mRemoteNG:** connections whose *PuttySession* names an imported PuTTY session get that session's proxy (unless they have a proxy already). The order does not matter: the second import links connections that the first one created.

After the import a summary lists the imported connections, folders, jump host links and tunnels, what was skipped and what to check (e.g. connections without a user name or password).

**Size limit:** every import file (WRM export, mRemoteNG, PuTTY, OpenSSH config, CSV / Excel inventory) may be up to **20 MB**. The browser checks the size before sending; a larger file gets *The file is too large (… MB): the import limit is 20.0 MB* (HTTP 413 from the API).

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

### Remote desktop: RDP, VNC, Telnet

Windows servers (RDP), Linux desktops and KVM/IPMI consoles (VNC) and old network gear (Telnet) open in a WRM window next to your SSH terminals — no RDP client, no plug-in.

```
 browser ──WebSocket── WRM ──guacd protocol── guacd ──RDP / VNC / Telnet── server
```

**Setup (once):** WRM uses **guacd**, the proxy daemon of [Apache Guacamole](https://guacamole.apache.org/), which implements the protocols. Run it on the same machine as WRM:

| Where WRM runs | Install guacd |
|---|---|
| Debian / Ubuntu | `sudo apt install guacd` (listens on 127.0.0.1:4822) |
| Docker | `docker compose up -d` — `docker-compose.yml` contains a `guacd` service that shares WRM's network |
| Any Linux with Docker | `docker run -d --name guacd --network host --restart unless-stopped guacamole/guacd` |
| Windows / macOS | run guacd in Docker (Docker Desktop) or on a Linux machine and set `guacd_address` |

*Admin panel → Overview* shows whether guacd answers (and its version); *Test connection* on an RDP/VNC connection checks the server's port **and** guacd.

**Connections:** choose protocol **RDP**, **VNC** or **Telnet**, host (default ports 3389, 5900, 23), user, password and optionally a jump host. **🖥 Remote desktop options**:

- **RDP:** domain, security (automatic, NLA, TLS, RDP, Hyper-V console), keyboard layout of the server, what happens when the window is resized (adapt the desktop size, reconnect, or keep and scale), colors, start program, accept the server certificate, console (admin) session, sound, clipboard, wallpaper, **RD Gateway** (host, port, user, domain).
- **VNC:** colors, local or remote mouse pointer, view only, clipboard.
- **Telnet:** font size, colors, login and password prompt patterns (for automatic login).

**Working in a desktop window:**

- **Double-click** the connection (or right-click → *Open remote desktop*). The window shows the remote screen, scaled to fit; **⤢ Fit** switches to 1:1 with scroll bars, **⛶** to full screen. With RDP the remote desktop adapts its size to the window.
- **Keyboard and mouse** go to the remote computer while the window has the focus (click it). **Ctrl+Alt+Del** has its own button.
- **Clipboard:** text copied on the remote computer appears in WRM's clipboard panel and in your clipboard; **📋 Clipboard** sends text to the remote clipboard. **⌨ Type** types text as keystrokes — for login screens, BIOS/iDRAC consoles and other places without a clipboard.
- **Sound** from RDP plays in the browser.
- If the server asks for credentials that are not stored (e.g. an empty password), WRM asks for them.
- **↻** reconnects; the tab and the window show the connection state.
- **Behind a jump host**, WRM opens a temporary tunnel through the jump hosts for the session (listed as *Remote desktop (temporary)* in the Tunnels panel) and closes it with the session.
- Saved workspace sessions restore desktop windows too.

**Audit and recording:** every desktop session appears in *Sessions & recordings* and the audit log (`desktop.open`, `desktop.close`, `desktop.error`). With session recording on, WRM records the screen stream (no passwords: they are sent to guacd during the handshake and never appear on screen) and **▶ Replay** plays it in the browser with play/pause and seek; *⬇ .guac* downloads it (playable with Guacamole's tools).

**Security:** the browser never picks the target: WRM performs the handshake with the stored host, user, password and options, and afterwards forwards only display data and input instructions from an allow-list. Access follows the same rules as terminals (owner; share members with the *Operator* role or higher). Policies: `desktop_enabled` (`REMOTE_DESKTOP_ENABLED`), `guacd_address` (`GUACD_ADDRESS`, default `127.0.0.1:4822`), `desktop_tunnel_bind` (address of jump-host tunnels for guacd, default `127.0.0.1`).

### Out-of-band management, consoles & serial ports

**The BMC of a server** (iDRAC, iLO, XClarity, Supermicro, OpenBMC, any IPMI controller) is set in the connection: *Edit → ⚡ Out-of-band management*:

- **Type:** **Redfish** (the HTTPS API of current BMCs) or **IPMI** (IPMI over LAN with `ipmitool`).
- **BMC address** (`10.0.100.11`, or `host:port`), **login**: a user name and password (stored encrypted), or a **vault credential** (e.g. one `idrac-root` for all servers).
- **Through the jump host:** when the management network is only reachable from the bastion, WRM reaches the BMC through the connection's jump hosts — Redfish over an SSH channel of the last hop, `ipmitool` runs **on** the last hop (IPMI uses UDP).
- **Serial console command** for the BMC's SSH (`console com2` for iDRAC, `vsp` for iLO, `console 1` for XClarity; suggested from the BMC's maker) and its SSH port.

**⚡ Power & status** (right-click → *Power & status (BMC)…*):

- Power state, health (with fault messages over IPMI), maker, model, serial number, BIOS, CPUs and memory, BMC model and firmware, power draw and inlet temperature (as far as the BMC tells).
- **Power on**, **Shut down** (graceful), **Restart** (graceful), **Force off**, **Force restart**, **Power cycle**, **NMI** — only what the BMC offers. Everything except *Power on* asks first and explains what happens; every action is in the audit log (`bmc.power`, `bmc.power_failed`).
- Redfish certificates are mostly self-signed: WRM pins the certificate on first use like an SSH host key and refuses a changed one (accept it after checking, like a host key).

**📟 Serial console** (right-click → *Serial console*, or the button in *Power & status*) opens a terminal window with the console of the machine — BIOS/UEFI setup, boot loader, kernel messages, a login prompt when the network is down:

- **Redfish BMCs:** WRM logs in to the BMC's SSH with the BMC login and types the console command (iDRAC `console com2`, iLO `vsp` …).
- **IPMI:** Serial-over-LAN (`ipmitool sol activate`, an old session is deactivated first), on the jump host or on the WRM server (Linux). `~.` ends it.
- Consoles are terminal sessions like any other: audited, recorded (`protocol` *console* / *sol*), visible to administrators — and only the connection's owner opens them (not share members). Saved workspace sessions restore console windows.

**🔌 Serial ports of the WRM machine** (USB console cables, `ttyS0`): protocol *Serial port*, the device name (`ttyUSB0`, `ttyACM0`, `ttyS0`) as host and the line settings — speed (default 9600; 115200 is common too), data bits, parity, stop bits, flow control (none, RTS/CTS, XON/XOFF). *Test connection* opens the port; a double-click opens it in a terminal (raw, recorded). Serial ports belong to the WRM server like its key files: by default only administrators may add them (policy `serial_ports`: off / admins / all); WRM must run on Linux, as a user that may open the device (group `dialout`).

Policies: `bmc_enabled` (`BMC_ENABLED`), `serial_ports`.

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

### Folder bookmarks

Bookmarks are named remote directories, like WinSCP's: the same deep folders (`/opt/app/servers/<server>/downloads`) are one click away on every server.

- **★ in the file manager** (next to the path bar): **☆ Bookmark this directory…**, the bookmarks of this connection, and **⚙ Manage bookmarks…**. Picking a bookmark opens that directory. If it does not exist on this server, a message says so and the current folder stays open.
- **★ in a terminal's title bar** and in the **Ctrl+Shift+Space** picker (after the snippets): picking a bookmark types `cd -- '<path>'` into the terminal (shell-quoted; a leading `~/` stays outside the quotes so the shell expands it). In a [broadcast](#broadcast-input), **📣** next to a bookmark sends the `cd` to every terminal of the group, each with its own connection's variables.
- **Scope:** **one connection**, the connections of **a folder**, the connections with **a tag**, or **all SSH / SFTP / FTP connections** (global).
- **Variables**, filled in by the server for each connection: `~` (home directory), `$USER` / `${USER}` / `{user}` (the connection's user name), `{host}` (without the port), `{name}` (connection name). Example: `/opt/app/servers/{name}/downloads`.
- **Start directory:** tick *Start directory* and the file manager of a matching connection opens there; a terminal types the `cd` when it connects (also after a reconnect), before the [run-on-connect](#snippets--run-on-connect) snippets. The most specific start directory wins: connection, then tag, then folder, then global. Right-click a connection or a folder → **★ Start directory…** creates one.
- **Name, colour and note**; **Manage bookmarks** renames, reorders (↑ ↓), moves bookmarks between scopes and deletes them.
- **Sharing:** a global bookmark marked *Share with share members* is offered (read-only) to the people who use your connections through a [share](#sharing-members-roles--links).
- **Import from WinSCP:** *Manage bookmarks → Import WinSCP.ini* reads the remote directory bookmarks (`[Configuration\Bookmarks\Remote\…]`) of a WinSCP configuration file and adds them as global bookmarks, with the WinSCP site in the note. Bookmarks that already exist (same name and path) are skipped.
- Bookmarks are included in *Export* / *Import*, copied with *Duplicate*, and removed with their connection or folder.

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
- **Notifications:** a message when a server goes down or comes back, and a desktop notification while WRM is in the background (*Settings → General → Notify when a server goes down*). Outside the browser: e-mail, chat or push through [Notifications](#notifications-e-mail-chat-push-webhook).
- **Through a proxy:** connections with a proxy are checked through it (and through their jump hosts when *Check through jump hosts* is on).
- **🔄 Check status now** (right-click a connection or a folder) checks immediately — also **through the jump hosts** — and shows the result per connection.
- **How it checks:** a TCP connection to the connection's port (SSH 22, FTP 21, web 80/443 or the port you set). SSH checks read the server banner and answer with `SSH-2.0-WRM_status_check` before closing, FTP checks read the 220 greeting and send `QUIT` — the same as Nagios `check_ssh`/`check_ftp`. Nothing logs in. A failed check is repeated once before a server counts as down. Each distinct host:port is checked once per round, however many connections and users point to it.
- **Behind a jump host** a connection shows the state of its jump host (a red ring when the jump host is down). With the policy *Check through jump hosts*, WRM logs in to the jump host once per round and checks the ports behind it through that SSH connection.
- Untick **Monitor up/down status** in a connection to leave it out (e.g. servers that are often off on purpose).
- Policies: `status_enabled` (on), `status_interval_seconds` (60, 15–3600), `status_jump_checks` (off).

> **fail2ban / IDS:** the check is an ordinary connection without a login. Default fail2ban `sshd` filters do not count it, but the *aggressive* / *ddos* modes may. Add the WRM server to `ignoreip`, raise the interval, or untick monitoring for such hosts.

### Notifications (e-mail, chat, push, webhook)

WRM tells you about events outside the browser too. **Administrators** set up the channels (*Admin panel → Notifications*), **users** choose what they receive and where (*Settings → Notifications*). The in-app messages and browser notifications stay as they are.

- **Channels:**

| Type | What it needs |
|---|---|
| E-mail (SMTP) | server, port, security (STARTTLS / TLS / none), optional login, from, to. A login is only sent over TLS or to a server on the WRM machine |
| Telegram | bot token (@BotFather), chat ID; optionally the address of a local Bot API server |
| Slack / Mattermost / Rocket.Chat | incoming webhook URL (Slack format) |
| Microsoft Teams | a *Workflows* webhook URL (*Post to a channel when a webhook request is received*); sent as an Adaptive Card |
| Discord | webhook URL |
| ntfy | server (ntfy.sh or your own), topic, optional access token |
| Gotify | server, application token |
| Pushover | application token, user or group key |
| Webhook | URL; WRM POSTs JSON (`title`, `text`, `events`). With a signing secret the header `X-WRM-Signature: sha256=<HMAC-SHA256 of the body>` lets the receiver check it |

- Secrets (passwords, tokens, webhook URLs) are **encrypted** and never shown again; leave a field empty to keep it. **📨 Send test** sends a test message at once and shows the service's answer (*HTTP 400 chat not found*, *SMTP login: …*). The last delivery and the last error are listed per channel.
- *Users may enter their own recipient* (e-mail, Telegram, ntfy, Pushover): each user can enter their own address, chat ID, topic or user key (e-mail: plain addresses only); otherwise the channel's is used. A channel without a default recipient needs the user's own. Slack-format and Discord messages cannot mention `@channel` / `@everyone` through connection names.
- **Events:** a server went down, a server is up again (only the owner's own monitored connections; connections with monitoring off are left out), credential rotation is due (again every 7 days while it stays due), credential rotation is incomplete (again every day).
- **One message per round:** what happens in one status round (or within a few seconds) goes out as **one digest** per user and channel, and a channel sends to the same recipient at most once per `notify_min_interval_seconds` (60); what comes in meanwhile goes into the next message.
- **Quiet hours** (per user, in the browser's time zone): messages wait and are sent together when the quiet hours end.
- Messages are in the user's language (English / Hrvatski). Waiting messages survive a restart; a failing channel is retried with a growing delay and gives up after 6 attempts.
- Policies `notifications_enabled` (on) and `notify_min_interval_seconds`. Channel changes and tests are audited (`admin.notify_channel_*`), subscriptions too (`notify.settings_changed`).

### Git workspace: compare services

**⎇ Git** in the top bar opens a full-screen workspace (its code, `static/git.js`, loads only then). It answers *which version of which service runs where, and what was changed by hand* — for any GitLab or GitHub (cloud or self-hosted) that you connect. WRM ships no services or organisation-specific configuration. Servers need no git, Python or internet access: WRM reads the Git provider's API itself and works on servers over its existing SSH connections, so jump hosts, proxies and the vault apply. **This version only reads**: nothing on a server is updated, installed or restarted.

- **Sources** (*Settings*): GitLab (API v4) or GitHub (REST; `github.com` or Enterprise `<url>/api/v3`) with a **read-only token** (encrypted, never sent back to the browser). *Test* lists the projects of the catalog's groups / organisations (with subgroups). Self-signed Git servers are pinned on first use.
- **Few API calls:** one recursive tree listing per version (cached forever, a commit never changes) gives the blob ID of every file; a file's content is downloaded only for a blob WRM has not seen yet and cached by blob ID — the same file in 20 versions is fetched once.
- **Offline bundles** (*Settings → Offline bundles*): when WRM cannot reach the Git server, import a `.tar.gz` built elsewhere — one top directory with `bundle.json` (per service the target ref and the hash and history of every file) and `payload/<service>/<path>`; `deploytool.py` and `README.txt` are ignored, protected files of the payload are never used. The list shows the bundle's age, author and source. *Services → Export bundle* writes the same format for servers that only a file-based deploy tool can reach. *Targets from*: automatic (the Git API, otherwise the newest bundle), the Git API or a bundle.
- **Service catalog** (*Services*): name, project, ref (`tag:latest`, `tag:v2` = the newest `v2.x`, or a branch), branch for tags, tag filter (regex), **subdirectory** of the repository that is the installation root, include / exclude globs, **protected** per-host files (`config*`, `*.ini`, `.env`; matched like Python's `fnmatch` on the path and the file name) that are never overwritten and only shown, **fingerprint** files that recognise an installation, install hints (directory names) and kind (app, library, tool). *Suggest from repository* fills an entry from the repository tree. *Import catalog* (merge or replace) and *Export catalog* use catalog JSON (`gitlab.url` / `groups`, `fallback_branches`, `server_roots`, `apps`, `ignore_dirs`). Services of imported bundles work without a catalog entry.
- **Refs:** only tags whose commit is on the chosen branch count (*branch-aware*); versions sort numerically (`v1.10` after `v1.9`); `tag:latest` without a matching tag uses the branch head, and a missing branch falls back to `fallback_branches` or the default branch — both with a warning. *Ref* in the overview: as in the catalog, one branch for all, or each project's default branch.
- **Discovery** (*Installations → Discover*, and with every check): one SSH call per server walks the server roots (`/opt`, `/srv`, …) to depth 3 and finds installations by their fingerprint files. Paths that look like copies or backups are skipped (segments with `bkp`, `backup`, `bak`, `old`, `orig`, `prev`, `copy`, `kopija`, `stari`, `recyclebin`, `trash`, `snapshot`, policy `git_backup_words`, plus the catalog's `ignore_dirs`). The environment comes from the path: `/srv/scripts/test/App` → *test*, `/opt/App_prod` → *prod*. Installations can also be added by hand.
- **Comparison:** the server hashes its files with `tr -d '\r' < f | sha256sum` (or `shasum` / `openssl`); WRM compares them with the target (SHA-256 with CRLF → LF) and with the earlier versions of every file (the earlier tags of the branch, *History depth* in *Settings*). Files are **ok**, **old** (*n behind*, with the tag it matches), **modified**, **missing**, **extra** (only on the server) or **protected**; an installation is **review** when a file is modified, **needs update** when a file is old or missing, otherwise **up to date**. `VERSION.md` and `updates.jsonl` of deploy tools, `.deploy-bak`, `.git`, `__pycache__` are left out.
- **Details:** the version recorded in `VERSION.md` (`**Verzija**:` line, otherwise the `# <service> - <version>` title; *missing* when there is none), the systemd / supervisor units whose unit file mentions the directory, the file states and a coloured **diff** of a file (server → target).
- **Overview:** the target of every service (version, latest tag, source, warnings) and a **servers × services matrix** with the state, environment, `VERSION.md` version and *n behind* of each installation, filtered by state, environment, folder, tag and text.
- **Periodic checks** (*Settings → Periodic checks*): while WRM runs, the chosen servers are checked every `git_check_interval_minutes` (60). Checks never change anything. They notify (one digest per round) about **a new version** of a service, **drift** (an installation that turns to *review*) and **unreachable servers** — subscribe to the *Git* events in *Settings → Notifications*.
- Everything is per user (sources, catalog, installations). Policies `git_enabled` (hides the workspace), `git_checks` (who may run checks: everyone), `git_check_interval_minutes`, `git_backup_words`. Checks, bundle imports / exports, source and catalog changes and viewed diffs are audited (`git.*`).

### Network tools

**🧰 Tools** in the top bar (or right-click a connection → *Network tools from this server…* / *Check ports of this host*) — the usual first checks when something does not answer, without opening a terminal:

| Tool | What it does |
|---|---|
| 🔍 **Port check** | TCP connection to up to 100 ports (`22, 80, 443`, ranges like `8000-8010`): **open** (with time), **closed** (refused) or **no answer** (filtered by a firewall) |
| 📶 **Ping** | 4 pings, the output of `ping` |
| 🛤 **Traceroute** | `traceroute` (or `tracepath`; `tracert` on Windows), up to 20 hops |
| 🌍 **DNS** | Addresses, CNAME, MX, TXT, and the name of an IP address (PTR) |
| 🔒 **HTTP / TLS** | One request to the URL (redirects are shown, not followed): status, time, server; for HTTPS the TLS version, certificate name, issuer, validity (**days left**), names and whether it is **trusted** |

**From** chooses where the check runs: the **WRM server**, or **one of your SSH connections** — through its jump hosts. That answers "can *app-01* reach the database on 5432?" or "does the firewall between *dmz-web* and the API let 443 through?". From a server, port and HTTP checks go through SSH channels (nothing to install there); ping, traceroute and DNS run the commands of the server (`ping`, `traceroute`/`tracepath`, `getent`/`dig`/`nslookup`).

Only your own connections can be used as a source; one check per user at a time; targets are host names or addresses only (no options or shell characters); every check is audited (`nettool.run`). Policy `network_tools`: all users (default), administrators only, or off.

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
| General | Language (English / Hrvatski), **install as an app**, accent color, close-window confirmation, **notifications when a server goes down**, browser-side upload limit |
| Notifications | Events × channels to receive outside the browser, your own recipient where the channel allows it, quiet hours, *Send test* |
| Terminal | Font size, scrollback, auto-reconnect |
| Security | Change password, two-factor authentication (enable/disable, new recovery codes), signed-in devices |
| Voice & audio | Microphone, speaker, level meter, noise suppression, echo cancellation, gain control, input mode / push-to-talk key, call sounds |
| Session history | Your terminal and remote desktop sessions and sessions on your connections, with replay and download |
| Data | Export connections (without secrets), export **with** passwords & keys (asks for your password; policy), import (also jump hosts and tunnels), **import from mRemoteNG**, **PuTTY** and **OpenSSH config** |

Personal preferences are stored per browser; everything security-related is stored on the server.

### Admin panel

*Settings → 🛡 Admin panel* (administrators only):

- **Overview:** version, uptime, users (admins, 2FA), HTTPS, voice relay status, encryption key source, **security warnings** (no HTTPS, open registration, admins without 2FA, host keys off…), **open terminals** (who, which server, from which IP — with *End*), **active SSH tunnels** (owner, connection and route, listen → target, traffic — with *Stop*), **live rooms** (participants, voice, shared terminals).
- **Users:** create users (temporary password generated if you leave it empty), display name, make/remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete (with everything the user owns). The last active administrator cannot be removed.
- **Shares:** every share on the server — pause/resume or delete.
- **Security policies:** registration, required 2FA, password length, lockout, session idle/maximum time, host-key policy, server key files, server-side upload limit, secret export, guest links, chat file size, audit log on/off and retention, session recording (on/off, keystrokes, retention, size limit), **broadcast input**, **live status** (on/off, interval, checks through jump hosts), **SSH tunnels** (on/off, who may use them, listening on network addresses, remote forwarding, idle stop), **proxies** (who may define them).
- **Notifications:** notification channels (e-mail, Telegram, Slack / Mattermost / Rocket.Chat, Teams, Discord, ntfy, Gotify, Pushover, webhook) with *Send test*, last delivery and error; notifications on/off and the minimum interval.
- **Voice & network:** voice on/off, participants per call, built-in TURN relay (port, public IP, host name, relay ports, private networks), additional STUN/TURN servers.
- **Host keys:** remembered SSH host keys and FTPS certificates; forget an entry after a server was reinstalled.
- **Audit log:** searchable and filterable (event type, user, date range), linked to sessions, **CSV export**, **Verify integrity** (hash chain).
- **Sessions & recordings:** every terminal session, with replay and `.cast` download.
- **File transfers:** every transferred file with size and SHA-256, CSV export.

Any policy can also be **forced by an environment variable** (`WRM_<KEY>`), e.g. for configuration management; it is then shown locked in the admin panel.

### Mobile, responsive UI & installable app

- The top bar and the collaboration bar adapt to the available width (labels collapse into icons, step by step).
- **Installable app (PWA):** *Settings → General → Install as an app* (or the install icon in the address bar; on iPhone/iPad: *Share → Add to Home Screen*) opens WRM in its own window, from the dock, start menu or home screen. The app caches only its own scripts and fonts (per version); connections and data always come live from the server, and without the server it shows an offline page. Browsers allow this only over **HTTPS with a trusted certificate** (or on `localhost`). The *Install as an app* row shows the steps for the current browser, or why installing is not possible (plain `http://`, an untrusted certificate, Firefox on the desktop); the *Install* button appears only when the browser can install right away.
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
| `allow_secret_export` | `1` | allow users to export their secrets, export private keys of the key store and show vault passwords (after re-entering the password) |
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
| `desktop_enabled` | `1` | RDP / VNC / Telnet in the browser. Also `REMOTE_DESKTOP_ENABLED` |
| `guacd_address` | `127.0.0.1:4822` | host:port of guacd. Also `GUACD_ADDRESS` |
| `bmc_enabled` | `1` | out-of-band management (BMC power, status, consoles). Also `BMC_ENABLED` |
| `serial_ports` | `admins` | who may use serial ports of the WRM server: `off`, `admins`, `all` |
| `network_tools` | `all` | who may use the network tools (port check, ping, traceroute, DNS, HTTP/TLS): `off`, `admins`, `all` |
| `proxies` | `all` | who may define saved proxies: `off`, `admins`, `all` (everybody may use proxies shared with them) |
| `notifications_enabled` | `1` | send notifications through the channels of *Admin panel → Notifications*. Also `NOTIFICATIONS_ENABLED` |
| `notify_min_interval_seconds` | `60` | at most one message per user and channel in this time (0–86400); the rest goes into the next digest |
| `git_enabled` | `1` | the Git workspace (⎇ Git); off hides it completely. Also `GIT_ENABLED` |
| `git_checks` | `all` | who may run Git checks (discovery and comparison over SSH, read-only): `off`, `admins`, `all` |
| `git_check_interval_minutes` | `60` | interval of periodic Git checks for users who turned them on (0 = none, up to 10080) |
| `git_backup_words` | `bkp,backup,bak,old,…` | path segments with one of these words are treated as copies / backups by Git discovery and comparison |
| `desktop_tunnel_bind` | `127.0.0.1` | address of the temporary jump-host tunnels guacd connects to (keep 127.0.0.1 when guacd runs on the WRM machine) |

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
- `127.0.0.1` inside the container is the container. The compose file adds `extra_hosts: ["host.docker.internal:host-gateway"]`, so a [proxy](#proxies-socks--http) on the Docker host is `host.docker.internal:<port>` (it must listen on an address the container reaches, e.g. the Docker bridge). With `docker run` add `--add-host=host.docker.internal:host-gateway`.

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
ExecStart=/opt/wrm/wrm-pro-v11.0.0-linux-amd64
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

### HTTPS with Caddy

[Caddy](https://caddyserver.com/) gets and renews a certificate by itself and proxies WebSockets without extra settings — the shortest way to HTTPS, which browsers need to **install WRM as an app**, for the microphone and for secure cookies. WRM listens on `127.0.0.1:8080` with `WRM_TRUST_PROXY=1` (as in the service above). `/etc/caddy/Caddyfile`:

```caddyfile
# A public DNS name pointing to this server (ports 80 and 443 open): a Let's Encrypt certificate.
wrm.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

Only an internal name or an IP address? Let Caddy be its own certificate authority and trust it once on every computer that opens WRM:

```caddyfile
wrm.lan.example.com, 10.0.0.10 {
    tls internal
    reverse_proxy 127.0.0.1:8080
}
```

```bash
sudo systemctl reload caddy
# root certificate of Caddy's local CA, to import into the browsers' / systems' trusted roots:
sudo cat /var/lib/caddy/.local/share/caddy/pki/authorities/local/root.crt
```

Caddy has no request body limit by default, so large imports and uploads pass. With WRM in Docker, point `reverse_proxy` at the published port and remove `HTTPS_SELF_SIGNED=1` from the compose file.

---

## Security model

See also **[SECURITY.md](SECURITY.md)** (how to report vulnerabilities, hardening checklist).

- **Accounts:** bcrypt hashes, optional/required TOTP 2FA with single-use recovery codes, closed self-registration by default, account lockout and per-IP rate limiting, constant-time responses for unknown users, forced password change after an admin reset, disable/delete accounts with immediate sign-out everywhere.
- **Sessions:** 256-bit random tokens in HttpOnly/SameSite cookies (`Secure` over HTTPS), stored hashed; idle and absolute expiry; device list and remote sign-out.
- **Secrets at rest:** AES-256-GCM with a random per-installation key (or your own) — connection passwords and keys, the SSH key store and the credentials vault. Secrets are **never sent to the browser** — not to share recipients, not to users a credential is shared with — except when their owner explicitly exports or shows them after entering the account password again (policy `allow_secret_export`, audited). The database and key files are created with mode `0600`.
- **Web security:** CSRF protection (custom request header + Origin check on every state-changing API call), WebSocket origin check, strict **Content-Security-Policy** (no external scripts, `frame-ancestors 'none'`), `X-Frame-Options`, `nosniff`, `Referrer-Policy: no-referrer` (share tokens never leak via Referer), `Permissions-Policy`, HSTS on HTTPS, request-size limits, no directory listings, all assets served locally (no CDN / supply-chain dependency at runtime).
- **SSH & FTPS:** host keys verified (TOFU/strict); FTPS certificates validated against public CAs or pinned. Server-side key files are admin-only by default.
- **Sharing:** four roles enforced server-side on every file operation, terminal and room action; members-only / signed-in / guest modes; signed password cookies (invalidated when the password changes); expiry, pause and link rotation; bans; immediate revocation of open terminals and room connections.
- **Collaboration:** server-assigned identities, per-client send queues (a slow client cannot stall a room), message rate limits and size limits, terminal data only to watchers, keystrokes only to the granting sharer, file downloads never rendered inline.
- **Voice:** WebRTC with DTLS-SRTP; the built-in TURN relay accepts only short-lived HMAC credentials issued to room participants, limits allocations, and refuses private/loopback/link-local peers by default (no pivoting into internal networks).
- **Audit:** sign-ins (also failed), account and policy changes, connections, shares, room joins and moderation, terminal sessions, file transfers (with SHA-256) and changes, host keys, exports, viewing of recordings — stored **append-only** and **hash-chained** in the database (verifiable) and written as `AUDIT …` lines to the server log (journald/syslog → SIEM). Secrets in audit details are redacted.
- **Session recording:** terminal output in asciicast v2, gzip, mode `0600`, SHA-256 in the database; no keystrokes unless enabled (then masked at password prompts); visible to administrators, the session's user and the connection's owner; every view is audited.
- **Folder bookmarks:** per user; the server fills in the variables and builds the shell-quoted `cd`, so a path cannot inject commands. Share participants see only the owner's bookmarks marked as shared (global, read-only), and only for connections of the share.
- **Snippets & broadcast:** snippets are per user (shared snippets only by administrators, read-only for others, never auto-run); run-on-connect snippets come from the connection's owner and are recorded in the audit log; broadcast input is client-side typing into the user's own terminals, with confirmation of dangerous commands and an audit entry for every start and stop per terminal.
- **Out-of-band management:** BMC passwords are encrypted (or come from the vault, with its host rules); Redfish certificates are pinned on first use, ipmitool gets the password in its environment (`-E`) — on a jump host through stdin, never on a command line; power actions and consoles are audited and limited to the connection's owner; serial ports are an administrator feature by default.
- **Inventory import:** files are parsed in memory (CSV, or `.xlsx` read with size limits) and never stored. NetBox is called from the WRM server over HTTP(S) with your token; redirects and pagination links to another server are refused, responses are size-limited, and a remembered token is stored encrypted and never returned to the browser.
- **SSH keys & vault:** private keys and vault passwords are encrypted at rest and used only on the server. Keys go to servers on stdin of a POSIX `sh` script (never on a command line); `authorized_keys` is changed in place, keeping other lines, permissions and SELinux contexts; WRM refuses to remove the key it logs in with. Shared credentials are usable but never readable by the people they are shared with; a host list limits where they can be sent (target and, for them, jump hosts). Rotation is all-or-nothing with a pre-flight check, verification of every new login and rollback; passwords are scrubbed from `passwd` output and never logged.
- **Remote desktop:** WRM does the guacd handshake with the stored credentials (never sent to the browser) and forwards only display and allow-listed input instructions; jump-host tunnels for guacd listen on loopback and live only as long as the session; sessions are audited and recorded.
- **Live status:** checks are plain TCP connections from the WRM server (no credentials, except optional checks through jump hosts with the jump host's own saved login); users only see the states of their own connections.
- **Tunnels & jump hosts:** every hop is authenticated and host-key-verified; tunnels belong to the owner of the connection (only the owner starts them; administrators can stop any); listen on `127.0.0.1` by default; ports below 1024 refused; network addresses and remote forwards only for administrators (policies); a tunnel can be turned off globally or limited to administrators; every start/stop/error is audited with the traffic. A tunnel port on a network address is **not** protected by WRM sign-in — treat it like an open port of that machine.
- **Proxies:** proxy passwords are encrypted, never sent to the browser and not shown to grantees; a proxy shared with a password is never reached through the grantee's own jump hosts; WRM tunnel proxies cannot be shared; IPMI / SOL never bypass a proxy silently (refused). Creation, changes, grants and tests are audited (`proxy.*`) and every session records its route.
- **Notifications:** channel secrets are encrypted and write-only; users only see enabled channels by name and type; messages go only to the owner of a connection or credential; configuration changes are audited as `admin.*` (always recorded).
- **Quick connect, notes & network tools:** quick connections are ordinary connections of their owner (same checks, encryption, host keys, audit) that expire; notes are only shown to the owner and rendered as escaped text; network tools only take host names or addresses (never options or shell characters), run from the WRM server or the user's own SSH connections, one at a time per user, and are audited (`nettool.run`) — limit them with `network_tools`.
- **Transport:** use HTTPS (certificate, reverse proxy, or `HTTPS_SELF_SIGNED=1`). Without it passwords and terminal traffic between browser and WRM are not encrypted and browsers block the microphone.

---

## Limitations

- Ping and traceroute use the commands of the WRM server or of the source server (`ping`, `traceroute`/`tracepath`); without them the tool says so. The app (PWA) installs only over HTTPS with a trusted certificate or on `localhost`.
- IPMI needs `ipmitool` (on the WRM server or the jump host); Serial-over-LAN on the WRM server and serial ports need WRM on Linux. Power actions go to one server at a time.
- Inventory import reads `.xlsx` and CSV, not old `.xls`; NetBox sync updates and reports, but never deletes connections.
- Password rotation works for SSH logins that may run `passwd` (Linux, BSD, macOS) — not for Windows/RDP, network devices or IPMI. *Who has access* shows `~/.ssh/authorized_keys` of the login user only.
- RDP, VNC and Telnet need **guacd** next to WRM (see [Remote desktop](#remote-desktop-rdp-vnc-telnet)); file transfer and printer redirection of RDP are not offered.
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
| GET / POST | `/api/folders` · PUT / DELETE `/api/folders/{id}` | Folders: `{name, jump_id, proxy_id}` (defaults of its connections; views add `jump_name`, `proxy_name`) |
| GET / POST | `/api/sessions` · PUT / DELETE `/api/sessions/{id}` | Workspace sessions |
| GET | `/api/config/export` | Export without secrets (with jump hosts and tunnels) |
| POST | `/api/config/export` | `{password}` export with secrets |
| POST | `/api/config/import` | Import → `{imported, skipped, tunnels}` |
| POST | `/api/config/import/mremoteng` | `{xml, password}` import an mRemoteNG `confCons.xml` → `{imported, folders, jump_hosts, tunnels, skipped, notes}`; HTTP 400 with `need_password: true` when a master password is needed |
| POST | `/api/config/import/sshconfig` | `{text, folder}` import an OpenSSH config |
| POST | `/api/config/import/putty` | `{data, folder}` import a PuTTY `.reg` export (`data` = the file, base64; UTF-16 or UTF-8) → also `{proxies, proxy_links}` |

Import files may be up to 20 MB (request bodies of `/api/config/import*` and `/api/inventory/*` up to 32 MB); larger ones get HTTP 413 with the size and the limit.

Connections have `jump_id` (jump host connection id or `null`) and `web_path`; lists also return `route` (e.g. `"vpn-gw → bastion-dc1"`) and `tunnels` (number of configured tunnels).

**Remote desktop**

| Endpoint | Description |
|---|---|
| `/ws/desktop?id&width&height&dpi&tz[&share_token]` | WebSocket (subprotocol `guacamole`) for an RDP/VNC/Telnet connection. WRM does the guacd handshake; the browser then exchanges Guacamole instructions (display ← / input →). First message: the tunnel id; then `wrm-session` (session id, recorded, route) |
| `GET /api/recordings/{id}/guac[?download=1]` | Recording of a remote desktop session (Guacamole protocol stream) |

Connections of type RDP/VNC/TELNET have `options` (see [Remote desktop](#remote-desktop-rdp-vnc-telnet)); sessions in `/api/recordings` have `protocol` and `recording.format` (`asciicast-v2+gzip` or `guacamole+gzip`).

**Quick connect, notes & network tools**

| Endpoint | Description |
|---|---|
| `POST /api/connections/quick` | `{target: "user@host:port" or "rdp://…", protocol?, password?, credential_id?, key_id?, jump_id?}` → the new temporary connection (`temporary: true`) |
| `PUT /api/connections/{id}` | `{…, notes, temporary: false}` — `notes` sets the notes, `temporary: false` keeps a quick connection |
| `GET /api/connections/{id}` | Includes `notes` (lists only have `has_notes`) |
| `POST /api/nettools` | `{tool: ports\|ping\|traceroute\|dns\|http, target, ports?, source_id?}` → `{source, target, ms, ports: [{port, state: open\|closed\|filtered\|error, ms, info}] \| output \| dns: {addresses, cname, mx, txt, ptr} \| http: {status, tls, subject, issuer, not_after, days_left, names, trusted, trust_error}}` |
| `GET /manifest.webmanifest`, `GET /sw.js` | Web app manifest and service worker (installable app) |

**Out-of-band management**

| Endpoint | Description |
|---|---|
| `GET /api/connections/{id}/bmc` | `{type, power: on\|off\|unknown, health, manufacturer, model, serial, bios, cpus, cpu_model, memory_gib, bmc_model, bmc_firmware, power_watts, inlet_c, faults, actions, console}`. A changed certificate answers `502` with `hostkey` |
| `POST /api/connections/{id}/bmc` | `{action: on\|shutdown\|restart\|off\|reset\|cycle\|nmi}` |
| `/ws/ssh?id&…&console=bmc\|sol` | The serial console of the connection's server (same protocol as terminals) |

Connections have `bmc: {type: redfish|ipmi, host, username, password (write-only), credential_id, via_jump, console, ssh_port}` (`{type: ""}` removes it; views add `has_password`, `credential_name`). Protocol `SERIAL` connections have the device name as `host` and `options: {baud, data_bits, parity, stop_bits, flow}`.

**Tags & inventory**

| Endpoint | Description |
|---|---|
| `POST /api/connections/bulk` | `{action: "tag", ids, add: [tags], remove: [tags]}` → `{changed}` |
| `POST /api/inventory/parse` | `{file_name, data (base64), sheet, header}` → `{sheets, sheet, encoding, delimiter, headers, rows (preview), total, mapping}` |
| `POST /api/inventory/import` | `{file_name, data, sheet, header, mapping: {field: column}, options}` → import result. Fields: `name host port protocol username password folder tags env site rack role platform jump web_path` |
| `GET /api/inventory/netbox` | Saved NetBox address and filters (`has_token`, never the token) |
| `POST /api/inventory/netbox/preview` | `{query: {url, token, devices, vms, site, role, tenant, tag, status, host_from, env_field, insecure, remember}}` → `{rows, missing}` |
| `POST /api/inventory/netbox/import` | `{query, rows, options}` → import result (`imported`, `updated`, `skipped`, `notes`) |

`options`: `{default_protocol, username, auth: ""|credential|key, credential_id, key_id, folder_mode: column|site|role|tenant|fixed|none, folder, tags, update}`. Connections have `tags` (array; on PUT, leaving it out keeps them) and, when imported from NetBox, `source: "netbox"`.

**SSH keys & credentials vault**

| Endpoint | Description |
|---|---|
| `GET /api/keys` | Your key store (public keys, fingerprints, `used_by`, `deployed`; never private keys) |
| `POST /api/keys` | `{name, comment, generate: {type: ed25519\|rsa\|ecdsa, bits}}` or `{name, key, passphrase}` (private or public key). `400 {need_passphrase: true}` when the key is encrypted |
| `PUT / DELETE /api/keys/{id}` | Rename `{name}` / delete (`409` while in use) |
| `GET /api/keys/{id}/public` | Public key file (`.pub`) |
| `POST /api/keys/{id}/export` | `{password, passphrase}` → `{private_key}` (re-authentication, policy `allow_secret_export`) |
| `POST /api/keys/{id}/deploy` | `{conn_ids, use_key}` → `{results: [{conn_id, name, ok, already, switched, error}], ok, failed}` |
| `POST /api/keys/{id}/revoke` | `{conn_ids}` → `{results: [{conn_id, name, ok, removed, error}], ok, failed}` |
| `GET /api/connections/{id}/authorized-keys` | `{user, path, entries: [{line, key_type, bits, fingerprint, comment, options, key_id, key_name, current}], login_key}` |
| `DELETE /api/connections/{id}/authorized-keys?fingerprint=SHA256:…` | Remove that key from the server (`409` for WRM's own login key) |
| `GET /api/credentials` | Credentials you own, are granted or are shared with everybody (no secrets) |
| `POST /api/credentials` · `PUT / DELETE /api/credentials/{id}` | `{name, username, password, key_id, description, hosts, rotate_days, grants: [user ids], shared_all}`; empty `password` on PUT keeps it (owner only) |
| `GET /api/credentials/{id}/usage` | Connections that use it |
| `POST /api/credentials/{id}/reveal` | `{password}` → `{password, pending_password}` (owner, re-authentication) |
| `POST /api/credentials/{id}/rotate` | `{mode: check\|rotate, new_password}` → starts a background job |
| `GET /api/credentials/{id}/rotation` | `{running, stage: preflight\|change\|rollback\|done, ok, message, hosts: [{key, connections, state, error}]}` |

**Proxies**

| Endpoint | Description |
|---|---|
| `GET /api/proxies` | `{proxies: [...], can_define, container}`: proxies you own, are granted or are shared with everybody (no secrets; `url`, `warning: docker_loopback`) |
| `POST /api/proxies` · `PUT / DELETE /api/proxies/{id}` | `{name, kind: socks5\|socks4\|http\|wrm_tunnel, host, port, username, password, credential_id, remote_dns, tunnel_id, description, grants, shared_all}`; empty `password` on PUT keeps it (owner only; policy `proxies`); DELETE is refused (409) while connections or folders use it |
| `POST /api/proxies/{id}/test` | `{target?: "host:port"}` → `{ok, message, latency_ms}`: connects to the proxy (and through it to the target) |

Connections have `proxy_id` (and `jump_id`): `null` = as the folder, `-1` = none, an id otherwise. Views return the effective `jump_id` / `proxy_id`, the stored `jump_choice` / `proxy_choice`, `jump_inherited` / `proxy_inherited`, `proxy_name` and the `route`. Exports name proxies (`proxy_ref`, and a `proxies` list).

**Notifications**

| Endpoint | Description |
|---|---|
| `GET /api/notify` | `{enabled, events, channels: [{id, name, kind, user_address, address_field}], subscriptions: [{event, channel_id, address}], prefs: {quiet_enabled, quiet_start, quiet_end, tz, lang}, pending}` |
| `PUT /api/notify` | `{subscriptions: [...], prefs: {...}}` replaces your subscriptions and quiet hours |
| `POST /api/notify/test` | `{channel_id, address}` a test message to you (once every 10 s) |
| `GET / POST /api/admin/notify/channels` | channels (secret fields empty, `secrets_set`), channel types and their fields · create `{name, kind, enabled, user_address, config: {…}}` |
| `PUT / DELETE /api/admin/notify/channels/{id}` | change (empty secret fields keep the stored value; `clear: [field]` removes one) · delete |
| `POST /api/admin/notify/channels/{id}/test` | `{address?}` send a test message now (HTTP 502 with the service's error) |

Events: `status.down`, `status.up`, `credential.rotation_due`, `credential.rotation_incomplete`.

Connections with `auth_method` `KEY_REF` have `key_id`, with `CREDENTIAL` a `credential_id` (views add `key_name`, `credential_name`, `credential_user`). Exports name keys and credentials (`key_ref`, `credential_ref`) instead of containing them; an import uses your key or credential with that name.

**Snippets & status**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/snippets` | My snippets and shared snippets / create `{name, command, description, group, scope: "all"\|"folder"\|"connection", scope_id, auto_run, shared}` |
| PUT / DELETE | `/api/snippets/{id}` | Update / delete (owner; administrators also shared snippets) |
| GET / POST | `/api/bookmarks` | My folder bookmarks / create `{name, path, scope: "global"\|"folder"\|"tag"\|"connection", scope_id, tag, color, note, start_dir, shared}` |
| GET | `/api/bookmarks?conn={id}` | Bookmarks for one connection with `resolved` path and `cd` command, plus `start` (also with `share_token`: the owner's shared ones) |
| POST | `/api/bookmarks/resolve` | `{id, conns: [...]}` → `{conn_id: {path, cd}}` (a broadcast group) |
| POST | `/api/bookmarks/reorder` | `{ids: [...]}` |
| POST | `/api/bookmarks/import-winscp` | Body: a `WinSCP.ini` → `{found, imported, skipped}` |
| PUT / DELETE | `/api/bookmarks/{id}` | Update (rename, move between scopes) / delete |
| GET | `/api/status` | `{enabled, interval, jump_checks, connections: {id: {state: "up"\|"down"\|"unknown", latency_ms, since, checked_at, error, banner, via, via_state}}}` |
| POST | `/api/status/check` | `{ids: [...]}` check now, also through jump hosts (once every 3 s per user) |

Connections also have `monitor` (live status on/off). The terminal WebSocket accepts `{"type":"broadcast","on":true\|false,"peers":n}` (audit) and sends `{"type":"autorun","snippets":[…]}` after run-on-connect snippets ran.

**Git workspace**

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/git` | Workspace state: `catalog`, `settings`, `sources` (no tokens), `targets`, `installs`, `bundle_apps`, `last_check_at`, `can_check` |
| PUT | `/api/git/settings` | `{mode: "auto"\|"live"\|"bundle", source_id, bundle_id, ref_mode: "catalog"\|"default"\|"branch", ref_branch, ref_overrides: {app: ref}, history_depth, conn_ids, monitor}` |
| PUT | `/api/git/catalog` | The whole catalog (catalog JSON) |
| POST | `/api/git/catalog/import?mode=replace\|merge` | Import catalog JSON |
| GET | `/api/git/catalog/export` | Catalog JSON (download) |
| PUT / DELETE | `/api/git/catalog/apps/{name}` | Add or change / delete one service |
| POST | `/api/git/sources` | `{kind: "gitlab"\|"github", name, url, token}`; `PUT` / `DELETE /api/git/sources/{id}` (an omitted `token` keeps the stored one) |
| POST | `/api/git/sources/{id}/test` | `{projects}` of the catalog's groups |
| POST | `/api/git/suggest` | `{project, ref}` → `{app, files}` (catalog entry from the repository tree) |
| POST | `/api/git/bundles` | Body: a bundle `.tar.gz` → `{result: {source, files, payload, warnings}, state}` |
| POST | `/api/git/bundles/export` | `{apps: [...]}` → bundle `.tar.gz` |
| POST | `/api/git/targets/refresh` | `{apps}` (empty = all): resolve refs and build targets |
| GET | `/api/git/targets/{app}` | Target with files and per-file history |
| POST | `/api/git/check` | `{conn_ids (empty = the chosen servers), discover, refresh}` → `{result: {servers, installs, targets}, state}` |
| POST | `/api/git/installs` | `{conn_id, path, app, env}` (added by hand) |
| GET / DELETE | `/api/git/installs/{id}` | Details with file states / forget |
| GET | `/api/git/installs/{id}/diff?path=` | `{lines: [{op: " "\|"-"\|"+"\|"@", a, b, s}], binary, truncated}` (server → target) |

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
| `/ws/events` | Live notifications for the signed-in user (`sessions_changed`, `tunnels_changed`, `snippets_changed`, `bookmarks_changed`, `status_changed`) |
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
  jump.go         jump hosts: chains, dialing through hops (and their proxies), validation
  proxy.go        saved proxies: SOCKS5 / SOCKS4(a) / HTTP CONNECT clients, WRM tunnel proxies, routes, API
  notify.go       notifications: channels (SMTP, chat, push, webhook), subscriptions, digests, quiet hours
  tunnels.go      SSH tunnels (local, remote, SOCKS5), manager, reconnects, tunnels API, web interfaces
  importers.go    import from mRemoteNG confCons.xml and OpenSSH config
  putty.go        import of PuTTY sessions (.reg) with proxies, link to mRemoteNG PuttySession
  snippets.go     snippets (saved commands), variables, run on connect
  bookmarks.go    folder bookmarks: scopes, variables, shell-quoted cd, start directory, WinSCP import
  status.go       live up/down status monitor and API
  git.go          Git workspace: storage, settings, targets, API, diff
  git_source.go   GitLab API v4 / GitHub REST client, tree and blob cache, ref resolution, suggestions
  git_bundle.go   offline bundles: import and export (bundle.json + payload)
  git_scan.go     discovery and file hashing over SSH, VERSION.md, units, comparison
  git_monitor.go  periodic Git checks and their notifications
  bmc.go          out-of-band management: Redfish and IPMI status and power, certificate pinning
  quick.go        quick connect (temporary connections and their cleanup), connection notes
  nettools.go     network tools: port check, ping, traceroute, DNS, HTTP/TLS check, from WRM or a server
  pwa.go          web app manifest and service worker (installable app)
  console.go      terminal backends: SSH shell, BMC SSH console, IPMI Serial-over-LAN, serial ports
  tty_linux.go    pseudo terminals and serial port settings (Linux); tty_other.go elsewhere
  inventory.go    tags, CSV/XLSX reader with column mapping, NetBox import and sync
  keys.go         SSH key store: generate/import, deploy/revoke (authorized_keys), who has access
  credentials.go  credentials vault: sharing, host lists, password rotation
  guac.go         remote desktop: guacd handshake, relay, recording (RDP / VNC / Telnet)
  helpers.go      origin check, login rate limiting, shared helpers
  security_test.go  unit tests (TOTP, roles, share access, CSRF, encryption migration, settings)
  upgrade_test.go   unit tests (upgrading databases from earlier builds, old URLs)
  integration_test.go  in-process SSH/SFTP server: connect → audit → recording, transfers, append-only audit
  tunnels_test.go      jump host chains, local/remote/SOCKS tunnels, policies, web interfaces
  snippets_test.go     snippets API, sharing, variables, run on connect, broadcast audit
  bookmarks_test.go    folder bookmarks: scopes, ownership, variables, cd quoting, shares, WinSCP.ini, start directory
  status_test.go       live status: up/down, banners, jump hosts, check now
  git_test.go          Git workspace: fake GitLab / GitHub, branch-aware tags, blob cache, bundles, catalog, discovery and comparison on a fake SSH host, notifications
  quick_test.go        quick connect (targets, expiry, open terminals, limit), notes, network tools from WRM and through SSH, PWA endpoints
  bmc_test.go          Redfish (fake BMC, pinning, jump host, credential), IPMI with a fake ipmitool, SOL (local PTY and jump host), BMC SSH console, serial port on a PTY
  inventory_test.go    tags, CSV encodings/separators, XLSX import, NetBox import and sync (fake NetBox)
  keys_test.go         key store, deploy/revoke/who has access, vault sharing and rotation (fake SSH server with sh and passwd)
  guac_test.go         remote desktop relay against a fake guacd: handshake, filtering, recording, jump hosts
  importers_test.go    mRemoteNG (encryption formats, inheritance, master password) and OpenSSH config import
  putty_test.go        PuTTY .reg import (UTF-16 / UTF-8, defaults, proxies and their reuse, SSH proxy → jump host), link with mRemoteNG in both orders
  proxy_test.go        proxies (in-process SOCKS5 / HTTP CONNECT, with and without auth): terminal, SFTP, status, jump host + proxy, web interface, desktop relay, WRM tunnel proxy, errors, sharing, folder defaults
  review_fixes_test.go  FTP with EPSV through a proxy (fake FTP server), another user's tunnel refused as a proxy, folder defaults (own WRM tunnel, loops refused, PUT keeps choices), per-user status keys, notification fixes (recipient injection, mentions, reminders, AUTH LOGIN)
  notify_test.go       notification channels against fake HTTP endpoints and a fake SMTP server, digests, rate limit, quiet hours, status and rotation events
  import_limits_test.go  import size limits through the real middleware: multi-MB confCons.xml and CSV, HTTP 413 above the limit
  static/index.html          the entire web UI (embedded into the binary)
  static/git.js              the Git workspace (loaded when it is opened)
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
| A bookmark does not appear under ★ | Its scope is another connection, folder or tag. In a share you see only bookmarks the owner shared (global ones) |
| *the directory … does not exist on this server* | The bookmark's path (after `{name}`, `{host}`, `$USER` are filled in) is not there or not readable on that server. Edit the path or use a narrower scope |
| Broadcast did not ask before a dangerous command | The check sees what you typed during the broadcast, not a command recalled from the shell history (↑) |
| All connections stay grey (no status dot) | Live status is off (`status_enabled`), or the connection has *Monitor up/down status* unticked. The first round starts a few seconds after WRM starts |
| A server shows 🔴 but SSH works | The port in the connection differs from the real one, or a firewall allows SSH only from some addresses (not from the WRM server). Behind a jump host use a jump host instead of a direct connection |
| BMC: *refused the login* / *the BMC user may not do this* | Check the BMC user and password (or the vault credential); the user needs the *Operator* or *Administrator* role on the BMC for power actions |
| BMC: *not reachable* | The management network is often only reachable from a bastion: set the jump host on the connection and tick *Through the jump host* |
| IPMI: *ipmitool is not installed* | `apt install ipmitool` on the WRM server, or on the jump host when the BMC is reached through it. IPMI over LAN must be enabled on the BMC |
| Serial console shows nothing | Press Enter. iDRAC: the console needs *Serial communication → Redirection* in the BIOS; IPMI SOL: enable SOL on the BMC and the serial console in the OS (`console=ttyS1,115200`) |
| Serial port: *permission denied* | Run WRM as a user in the `dialout` group (or give it access to the device) |
| No *Install* button / *Install as an app* only shows steps or a reason | The button appears only when the browser offers to install right away; otherwise the row shows the steps for your browser (Chrome/Edge/Brave: the icon in the address bar or the menu; Safari on macOS: *File → Add to Dock*; iPhone/iPad: *Share → Add to Home Screen*; Android: menu → *Install app*). *Not possible here* means: plain `http://` to an address other than `localhost` (use HTTPS with a **trusted** certificate, e.g. [Caddy](#https-with-caddy)), an untrusted certificate (`HTTPS_SELF_SIGNED=1`: trust it first, or use Caddy's `tls internal` CA), or Firefox on the desktop (use Chrome, Edge, Brave or Safari) |
| Import: *The file is too large (… MB): the import limit is 20.0 MB* | Import files may be up to 20 MB. Split the file (e.g. export mRemoteNG folders separately, or split the CSV) and import the parts — existing connections are skipped. Behind a reverse proxy allow bodies of at least 32 MB (nginx `client_max_body_size`) |
| `proxy socks5://…: not reachable` | WRM cannot open a TCP connection to the proxy: check its address and port, and with a jump host that the **jump host** can reach it. In Docker `127.0.0.1` is the container: use `host.docker.internal` |
| `proxy …: authentication required` / `authentication failed` / HTTP `407` | Set or correct the proxy's user name and password (or its vault credential) in *🌐 Proxies* |
| `proxy … cannot reach host:port: … (SOCKS5 reply 5)` / HTTP `502` | The proxy works but the target refused or is unreachable from the proxy. Check the host:port as seen **from the proxy**; turn off *DNS at the proxy* if the proxy cannot resolve internal names |
| *IPMI and Serial-over-LAN use UDP: they cannot go through the proxy* | Use Redfish for that BMC, or untick *Through the jump host* (ipmitool then runs on the WRM server) |
| *proxy … is shared with a password: it cannot be reached through your own jump hosts* | Ask the owner for a proxy without a password for this route, or define your own proxy |
| A connection suddenly goes through a jump host / proxy | Its folder has a default (*Folder settings…*, ⤳ after the folder name). Choose *no jump host / no proxy (not the folder's)* in the connection |
| Git: an installation is not found | Check *Settings → Catalog settings → Server roots* (searched to depth 3), the service's fingerprint files (relative to the installation root, all must exist) and whether the path looks like a backup (`git_backup_words`, `ignore_dirs`). Installations can be added by hand |
| Git: every file is *modified* | The *subdirectory* of the service must be the part of the repository that is the installation root (e.g. `linux`). The server needs `sha256sum`, `shasum` or `openssl` |
| Git: `HTTP 401` / `403` from the Git server | Check the token (read access to the API and repositories) in *Settings → Git sources*; 403 can also be the API rate limit |
| Notifications: nothing arrives | Check *Admin panel → Notifications* (channel enabled, last error, *Send test*), *Settings → Notifications* (events ticked, quiet hours) and the policy `notifications_enabled`. Messages of one round arrive together after a few seconds, at most one per `notify_min_interval_seconds` |
| Notifications: *SMTP login: … unencrypted connection* | Go's SMTP client sends a password only over TLS (or to localhost): choose STARTTLS or TLS |
| Network tools: *ping is not installed* | Install `iputils-ping` / `traceroute` on the WRM server (or on the source server), or use the port check — it needs nothing |
| A quick connection disappeared | Quick connections are deleted 24 hours after their last use. Save the ones you need (💾) |
| CSV import: wrong letters (`Ä�`, `?` instead of `č`) | WRM reads UTF-8, UTF-16 and Windows-1250. For other encodings save the file from Excel as *CSV UTF-8* or as `.xlsx` |
| NetBox: *refused the token* / *API not found* | Use the base address (`https://netbox.example.com`, without `/api`), a token of a user that can view devices and VMs; for a self-signed certificate tick *Accept a self-signed certificate* |
| Deploy SSH key: *cannot write ~/.ssh/authorized_keys* | The login user's home is read-only or its shell is not POSIX (network devices): add the key there by hand. A key that is not accepted afterwards: check `AuthorizedKeysFile` and `PubkeyAcceptedAlgorithms` (old servers may need RSA) in `sshd_config` |
| Rotation: *the stored password does not work on …* | Nothing was changed. Fix the login (or the password with *Edit*) and check again |
| Rotation: *passwd failed* / *rejected by its password rules* | The server's password policy (length, classes, history, minimum age) refused the new password: try your own password, or wait for the minimum age. The other servers got the old password back |
| *rotation incomplete* | Some servers have the new password and could not be changed back: *Show* reveals both; set the servers to one of them, then save it with *Edit* |
| *credential … is shared for hosts matching …: the jump host … does not match* | The credential's owner limited it to some hosts: the jump hosts of your connection must match the list too (ask the owner to add the bastion) |
| Remote desktop: *guacd … is not reachable* | Install and start guacd on the WRM machine (`apt install guacd` or the Docker image) or set `guacd_address`. *Admin panel → Overview* shows its state |
| RDP: *login failed* / black screen and disconnect | Check user, password and **domain**; try *Security: NLA* or *TLS*; older servers need *RDP*. Accept the server certificate (option) for self-signed certificates |
| RDP: wrong characters when typing | Set the **keyboard layout** to the server's layout (or *Unicode*) |
| Remote desktop behind a jump host does not connect | guacd must reach the temporary tunnel on `desktop_tunnel_bind` (127.0.0.1): run guacd on the WRM machine or in the same network namespace (compose file) |
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
