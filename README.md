# Web Remote Manager PRO (WRM)

**A remote server manager that runs in any web browser.** SSH terminal, SFTP / FTP / FTPS file manager, server-to-server transfers, saved workspaces, sharing with roles, real-time collaboration with **voice calls**, and enterprise security (2FA, audit log, policies): one self-hosted binary for your PC, server or company.

**Current version: v10.0.0-mimo** · [Download](https://github.com/vedranius/web-browser-RDM-public/releases/latest) · [Release notes](RELEASE_NOTES.md) · [Security](SECURITY.md)

---

## 📜 License: free for noncommercial use

WRM is **source-available** under the **[PolyForm Noncommercial License 1.0.0](LICENSE)**.

| | |
|---|---|
| ✅ **Free, no permission needed** | Personal use, home labs, learning, hobby projects, non-profits, schools and universities, public research, government. You may **use, fork, modify and share** it, including modified versions. Keep the [LICENSE](LICENSE) and its `Required Notice:` lines. |
| 💼 **Needs a commercial license** | Use in a company or for paid work, offering WRM as a hosted service / SaaS, bundling or selling it, or using its code in a commercial product. See **[COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)**. |
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
3. [Upgrading from v9](#upgrading-from-v9)
4. [How it works](#how-it-works)
5. [Features in detail](#features-in-detail)
   - [Accounts, sign-in & two-factor authentication](#accounts-sign-in--two-factor-authentication)
   - [Connections & folders](#connections--folders)
   - [Windows, tabs & snapping](#windows-tabs--snapping)
   - [SSH terminal](#ssh-terminal)
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
   - [Admin panel](#admin-panel)
   - [Mobile & responsive UI](#mobile--responsive-ui)
6. [Keyboard shortcuts](#keyboard-shortcuts)
7. [Configuration](#configuration)
8. [Running as a service, reverse proxy & firewall](#running-as-a-service-reverse-proxy--firewall)
9. [Security model](#security-model)
10. [Limitations](#limitations)
11. [API reference](#api-reference)
12. [Building from source & releases](#building-from-source--releases)
13. [Troubleshooting](#troubleshooting)
14. [Contributing](#contributing)

---

## What WRM is

WRM is a **single executable** with a built-in web server and a built-in web app. You start it on any machine: your PC, a Raspberry Pi, a VPS or a company server. Then you open it in a browser and manage your remote servers from there:

- **SSH terminals** in the browser (xterm.js, 256 colors, full-screen apps like `nano`, `vim`, `htop`, `mc`).
- **File manager** for **SFTP** (over SSH), **FTP** and **FTPS**: browse, upload, download, rename, delete, edit, search.
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
   chmod +x wrm-pro-v10.0.0-mimo-linux-amd64
   ./wrm-pro-v10.0.0-mimo-linux-amd64
   ```
   On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

   **Windows** (PowerShell), or just double-click the `.exe`:
   ```powershell
   .\wrm-pro-v10.0.0-mimo-windows-amd64.exe
   ```

   **Android (Termux)**
   ```bash
   pkg install wget
   wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v10.0.0-mimo/wrm-pro-v10.0.0-mimo-android-arm64
   chmod +x wrm-pro-v10.0.0-mimo-android-arm64 && ./wrm-pro-v10.0.0-mimo-android-arm64
   ```
3. Open **http://localhost:8080** (or `http://<server-ip>:8080`). For voice calls from other computers use HTTPS — the quickest way is `HTTPS_SELF_SIGNED=1` (see [Configuration](#configuration)).
4. **Create the administrator account** (the first account). After that, self-registration is **closed**: the administrator creates accounts in *Admin panel → Users* (or opens registration in *Security policies*).
5. Click **+ Connection**, enter host, user and password or key, then double-click the connection to open it.

On the first start WRM creates, next to the database:

| File | What it is |
|---|---|
| `remote_manager.db` | SQLite database (users, connections, shares, audit log). Mode `0600`. |
| `remote_manager.db.key` | Random 256-bit key that encrypts stored passwords, private keys and 2FA secrets. Mode `0600`. **Back it up separately** — without it stored secrets cannot be decrypted. Or set your own key with `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`. |

Locked out? `./wrm-pro-… -reset-password admin` prints a new temporary password (add `-reset-2fa` to also turn off two-factor authentication).

---

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
  - **Protocol:** `SSH` (terminal, plus files over SFTP), `SFTP` (opens the file manager by default), `FTP`, `FTPS` (FTP with explicit TLS; file manager only).
  - **Host : Port** (default port 22 for SSH/SFTP, 21 for FTP/FTPS), **Username**, **Folder**.
  - **Authentication:** *Password*, *Private key (paste)*, *Key file (path on the WRM server)*, or *Auto* (tries `~/.ssh/id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa` of the WRM server user). *Key file* and *Auto* use keys of the **WRM server itself** and are therefore limited to administrators (policy). Password authentication also answers keyboard-interactive prompts.
  - **🔌 Test connection** logs in once and shows the result, latency and — for a new server — its host key fingerprint.
- **Secrets never leave the server.** Passwords and private keys are write-only: the edit dialog shows “saved and encrypted — leave empty to keep it”, with an option to remove the saved secret. *Duplicate* copies the connection on the server.
- **Host key verification:** the first connection to a server remembers its SSH host key (and, for FTPS, the certificate if it is not signed by a public CA). A **changed key is refused** with a clear warning and both fingerprints; the owner of the connection (or an administrator) can accept the new key after checking it. Policy: *trust on first use* (default), *strict* (new hosts must be approved) or *off*.
- **Double-click** a connection to open it: SSH opens a terminal; SFTP/FTP/FTPS open the file manager. On touch devices a single tap opens it.
- **Right-click** a connection: *Open as Terminal / File Manager* (in a new or an existing window), *Edit*, *Duplicate*, *Share*, *Delete*.
- **📁 Folder** creates folders. **Drag** connections onto a folder to move them, or onto *"↓ Drop here"* to move them back to the root. Right-click a folder to *Share* or *Delete* it (its connections move to the root).
- **Search box** filters by name, host, username or protocol. **Multi-select** with **Ctrl/Cmd + click**, then use **🤝 Share** or **Delete**.
- A **green dot** next to a connection means one of its terminals is connected.
- **Shared with me** lists shares you are a member of (and shares published for all users), with your **role**; ↗ opens the collaboration room.

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

### Settings

**⚙ Settings** in the top bar has tabs:

| Tab | What it contains |
|---|---|
| General | Language (English / Hrvatski), accent color, close-window confirmation, browser-side upload limit |
| Terminal | Font size, scrollback, auto-reconnect |
| Security | Change password, two-factor authentication (enable/disable, new recovery codes), signed-in devices |
| Voice & audio | Microphone, speaker, level meter, noise suppression, echo cancellation, gain control, input mode / push-to-talk key, call sounds |
| Data | Export connections (without secrets), export **with** passwords & keys (asks for your password; policy), import |

Personal preferences are stored per browser; everything security-related is stored on the server.

### Admin panel

*Settings → 🛡 Admin panel* (administrators only):

- **Overview:** version, uptime, users (admins, 2FA), HTTPS, voice relay status, encryption key source, **security warnings** (no HTTPS, open registration, admins without 2FA, host keys off…), **open terminals** (who, which server, from which IP — with *End*), **live rooms** (participants, voice, shared terminals).
- **Users:** create users (temporary password generated if you leave it empty), display name, make/remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete (with everything the user owns). The last active administrator cannot be removed.
- **Shares:** every share on the server — pause/resume or delete.
- **Security policies:** registration, required 2FA, password length, lockout, session idle/maximum time, host-key policy, server key files, server-side upload limit, secret export, guest links, chat file size, audit retention.
- **Voice & network:** voice on/off, participants per call, built-in TURN relay (port, public IP, host name, relay ports, private networks), additional STUN/TURN servers.
- **Host keys:** remembered SSH host keys and FTPS certificates; forget an entry after a server was reinstalled.
- **Audit log:** searchable and filterable (sign-ins, admin actions, connections, shares, collaboration, terminals, files, host keys, export/import), **CSV export**.

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
| `audit_retention_days` | `365` | how long audit log and chat history are kept |

**Command line**

```
wrm -version                       print the version
wrm -reset-password USER           set a new temporary password (must be changed at sign-in)
wrm -reset-password USER -reset-2fa  … and turn off two-factor authentication
```

---

## Running as a service, reverse proxy & firewall

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
ExecStart=/opt/wrm/wrm-pro-v10.0.0-mimo-linux-amd64
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
- **Audit:** sign-ins (also failed), account and policy changes, connections, shares, room joins and moderation, terminals, file changes, host keys, exports — stored in the database and written as `AUDIT …` lines to the server log (journald/syslog → SIEM).
- **Transport:** use HTTPS (certificate, reverse proxy, or `HTTPS_SELF_SIGNED=1`). Without it passwords and terminal traffic between browser and WRM are not encrypted and browsers block the microphone.

---

## Limitations

- No **RDP/VNC** (graphical desktops): WRM is for SSH, SFTP, FTP and FTPS.
- No **SSO/LDAP/SAML** yet (local accounts with 2FA).
- Voice calls are a **mesh**: fine up to about 12 people; larger meetings would need an SFU.
- No SSH **agent forwarding**, **port forwarding/tunnels** or **jump hosts** yet.
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
| GET | `/api/config/export` | Export without secrets |
| POST | `/api/config/export` | `{password}` export with secrets |
| POST | `/api/config/import` | Import |

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
| GET | `/api/admin/audit?q&action&user&from&to&before_id&limit[&format=csv]` | Audit log |
| GET | `/api/version` | Server version |

**WebSockets**

| Endpoint | Description |
|---|---|
| `/ws/ssh?id&cols&rows[&share_token]` | Terminal. Binary frames = terminal data; text frames = JSON control (`resize`, `pause`, `resume`, `ping` → server; `status`, `error`, `exit`, `hostkey` ← server). Close codes: 1000 shell exited, 4001 connect failed / not allowed, 4002 SSH connection lost, 4003 access revoked |
| `/ws/events` | Live notifications for the signed-in user (`sessions_changed`) |
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
  audit.go        audit log, CSV export, retention
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
  helpers.go      origin check, login rate limiting, shared helpers
  security_test.go  unit tests (TOTP, roles, share access, CSRF, encryption migration, settings)
  static/index.html          the entire web UI (embedded into the binary)
  static/vendor/             xterm.js + addons and fonts (served locally, see THIRD-PARTY-LICENSES.txt)
.github/workflows/build.yml  CI: vet, gofmt, tests (race), JS syntax check, builds all platforms, releases on tags
```

**CI/CD:** every push builds all 13 targets. Pushing a tag `v*` creates a GitHub Release with all binaries, a macOS universal binary and `SHA256SUMS.txt`.

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
| Stored passwords stopped working after a restart | The encryption key changed (`ENCRYPTION_KEY`, `ENCRYPTION_KEY_FILE` or the `.key` file next to the database). Restore it, or re-enter the passwords |
| *File contents* search fails | The server needs `grep` and a POSIX shell; use *File names* search instead |
| Uploads fail at a certain size | The administrator's upload limit (`max_upload_mb`) or nginx `client_max_body_size` |
| Transfer/search progress appears only at the end | Disable proxy buffering (`proxy_buffering off;`) |

---

## Contributing

Bug reports, ideas, translations and pull requests are welcome! Please read **[CONTRIBUTING.md](CONTRIBUTING.md)**. It explains how to build, test and submit changes, and the contributor terms.

If WRM is useful to you, consider supporting it on **[Ko-fi](https://ko-fi.com/vedranius)** ☕.

---

<sub>Web Remote Manager PRO · © 2026 vedranius · [PolyForm Noncommercial 1.0.0](LICENSE) · commercial use: [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)</sub>
