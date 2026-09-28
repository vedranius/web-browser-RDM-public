# Web Remote Manager PRO (WRM)

**A remote server manager that runs in any web browser.** SSH terminal, SFTP / FTP / FTPS file manager, server-to-server transfers, saved workspaces, sharing and real-time collaboration: one self-hosted binary for your PC, server or team.

**Current version: v9.10.2-mimo** · [Download](https://github.com/vedranius/web-browser-RDM-public/releases/latest) · [Release notes](RELEASE_NOTES.md)

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
3. [How it works](#how-it-works)
4. [Features in detail](#features-in-detail)
   - [Accounts & login](#accounts--login)
   - [Connections & folders](#connections--folders)
   - [Windows, tabs & snapping](#windows-tabs--snapping)
   - [SSH terminal](#ssh-terminal)
   - [File manager (SFTP / FTP / FTPS)](#file-manager-sftp--ftp--ftps)
   - [Search in files](#search-in-files)
   - [Text editor](#text-editor)
   - [Server-to-server transfer](#server-to-server-transfer)
   - [Workspace sessions](#workspace-sessions)
   - [Clipboard panel](#clipboard-panel)
   - [Sharing connections](#sharing-connections)
   - [Real-time collaboration](#real-time-collaboration)
   - [Settings](#settings)
   - [Admin panel](#admin-panel)
   - [Mobile & responsive UI](#mobile--responsive-ui)
5. [Keyboard shortcuts](#keyboard-shortcuts)
6. [Configuration (environment variables)](#configuration-environment-variables)
7. [Running as a service & behind a reverse proxy](#running-as-a-service--behind-a-reverse-proxy)
8. [Security model](#security-model)
9. [Limitations](#limitations)
10. [API reference](#api-reference)
11. [Building from source & releases](#building-from-source--releases)
12. [Troubleshooting](#troubleshooting)
13. [Contributing](#contributing)

---

## What WRM is

WRM is a **single executable** with a built-in web server and a built-in web app. You start it on any machine: your PC, a Raspberry Pi, a VPS or a company server. Then you open it in a browser and manage your remote servers from there:

- **SSH terminals** in the browser (xterm.js, 256 colors, full-screen apps like `nano`, `vim`, `htop`, `mc`).
- **File manager** for **SFTP** (over SSH), **FTP** and **FTPS**: browse, upload, download, rename, delete, edit, search.
- **Server-to-server copy** between two SSH servers, without downloading to your computer first.
- **Workspaces**: many terminal and file windows side by side, tabs, snapping, saved sessions.
- **Sharing & collaboration**: share connections with colleagues or guests via a link, with chat, file exchange, live terminal sharing and remote keyboard control.
- **Multi-user**: every user has their own connections, folders and sessions. The first user becomes admin.

Everything is stored in one local **SQLite** file. There is no external database, no Docker requirement and no cloud account.

**Supported platforms (prebuilt binaries):** Linux (x86-64, arm64, ARMv7/Raspberry Pi, ARMv6, 32-bit), Windows (x86-64, arm64), macOS (Intel, Apple Silicon, Universal), Android arm64 (Termux), FreeBSD (x86-64, arm64), OpenBSD (x86-64).

---

## Quick start

1. Download the binary for your system from **[Releases](https://github.com/vedranius/web-browser-RDM-public/releases/latest)** and optionally verify it with `SHA256SUMS.txt`.
2. Run it:

   **Linux / macOS / FreeBSD / OpenBSD**
   ```bash
   chmod +x wrm-pro-v9.10.2-mimo-linux-amd64
   ENCRYPTION_KEY='choose-your-own-32-character-key!' ./wrm-pro-v9.10.2-mimo-linux-amd64
   ```
   On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

   **Windows** (PowerShell), or just double-click the `.exe`:
   ```powershell
   $env:ENCRYPTION_KEY='choose-your-own-32-character-key!'; .\wrm-pro-v9.10.2-mimo-windows-amd64.exe
   ```

   **Android (Termux)**
   ```bash
   pkg install wget
   wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v9.10.2-mimo/wrm-pro-v9.10.2-mimo-android-arm64
   chmod +x wrm-pro-v9.10.2-mimo-android-arm64 && ./wrm-pro-v9.10.2-mimo-android-arm64
   ```
3. Open **http://localhost:8080** (or `http://<server-ip>:8080`).
4. Click **Register** and create the first account. **The first account becomes the administrator.**
5. Click **+ Connection**, enter host, user and password or key, then double-click the connection to open it.

The database file `remote_manager.db` is created in the current working directory (change with `DB_PATH`). **Set your own `ENCRYPTION_KEY` before adding connections.** See [Security model](#security-model).

---

## How it works

```
 Browser (any device)                      WRM server (one binary)                     Your servers
┌──────────────────────┐   HTTPS/HTTP   ┌──────────────────────────────┐   SSH / SFTP   ┌───────────┐
│ Web UI (embedded     │ ─────────────▶ │ REST API  (/api/...)         │ ─────────────▶ │ Linux/BSD │
│ single-page app,     │                │ WebSockets:                  │   FTP / FTPS   │ servers,  │
│ xterm.js terminal)   │ ◀───────────── │  /ws/ssh     terminal        │ ─────────────▶ │ NAS, ...  │
│                      │   WebSockets   │  /ws/share/  collaboration   │                └───────────┘
└──────────────────────┘                │  /ws/events  live updates    │
                                        │ SQLite: users, connections,  │
                                        │ folders, sessions, shares    │
                                        └──────────────────────────────┘
```

- **The WRM server makes the connections.** Your browser talks only to WRM, and WRM talks to your servers with SSH (`golang.org/x/crypto/ssh`), SFTP (`pkg/sftp`) and FTP (`jlaffaye/ftp`). Your servers only need to be reachable **from the WRM machine**, not from every browser.
- **Terminal:** every terminal window is a WebSocket (`/ws/ssh`) bound to one SSH session with a PTY on the server. Keystrokes and terminal output go as **binary frames**, so large outputs and multi-byte UTF-8 characters are never corrupted. The browser tells the server the terminal size on connect and after every resize, so full-screen programs always fit the window.
- **Flow control:** when a command prints faster than the browser can draw (for example `cat` of a huge log), the browser asks the server to pause and resume. The tab stays responsive and nothing is lost.
- **File operations** go through a small **pool of SFTP connections** per saved connection (reused for 5 minutes and health-checked), so browsing is fast. Uploads are **streamed** straight to the remote server. Downloads and ZIP archives are streamed to your browser.
- **Secrets at rest:** passwords and private keys of connections are stored **AES-256-GCM encrypted** in SQLite. The key comes from `ENCRYPTION_KEY`.
- **Live updates:** `/ws/events` notifies open browser tabs when sessions change, so the Sessions panel stays in sync across tabs.
- **Collaboration:** `/ws/share/{token}` is a "room" per share link that relays chat, files, terminal-sharing data and remote-control keystrokes between participants.

---

## Features in detail

### Accounts & login

- Register with a username and password (min. 6 characters). Passwords are stored as **bcrypt** hashes.
- The **first user becomes administrator** (see [Admin panel](#admin-panel)).
- Login sessions use an HTTP-only cookie (valid 30 days). The browser pings `/api/auth/keepalive` every 60 s while you work. Sessions inactive for 7 days are removed automatically.
- **Brute-force protection:** 8 failed logins from one IP within 10 minutes lock that IP out for 5 minutes.
- If your login expires while terminals or transfers are running, WRM shows a message instead of reloading the page and killing your work.

### Connections & folders

Left sidebar ("PRO MANAGER"):

- **+ Connection** creates a connection:
  - **Protocol:** `SSH` (terminal, plus files over SFTP), `SFTP` (opens the file manager by default), `FTP`, `FTPS` (FTP with explicit TLS; file manager only).
  - **Host : Port** (default port 22 for SSH/SFTP, 21 for FTP/FTPS), **Username**, **Folder**.
  - **Authentication:** *Password*, *Private key (paste)*, *Key file (path on the WRM server)*, or *Auto* (tries `~/.ssh/id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa` of the WRM server user). If you point *Key file* to a `.pub` file, WRM looks for the matching private key.
  - **🔌 Test connection** logs in once and shows the result and latency, without saving.
- **Double-click** a connection to open it: SSH opens a terminal; SFTP/FTP/FTPS open the file manager. On touch devices a single tap opens it.
- **Right-click** a connection: *Open as Terminal / File Manager* (in a new window or an existing window of the same type), *Edit*, *Duplicate*, *Share*, *Delete*.
- **📁 Folder** creates folders. **Drag** connections onto a folder to move them, or onto *"↓ Drop here"* to move them back to the root. Opened folders stay open while you work. Right-click a folder to *Share* or *Delete* it (its connections move to the root).
- **Search box** filters by name, host, username or protocol.
- **Multi-select** with **Ctrl/Cmd + click**, then use the bulk buttons **🤝 Share** or **Delete**.
- A **green dot** next to a connection means one of its terminals is currently connected. Hover to see `user@host`.
- **Shared with me** shows connections other users shared with you (see [Sharing](#sharing-connections)).
- The sidebar can be **collapsed** (◀) and **resized** by dragging its edge. The state is remembered.

### Windows, tabs & snapping

Every opened connection is a **window** inside the workspace. There is also a **tab** for it in the tab bar.

- **+ New Window** (or **Alt+N**) opens an empty window with quick buttons for your connections and **saved sessions**.
- **Title bar:** drag to move, **double-click** to maximize or restore. Buttons: *snap layouts*, **↻ Reconnect** (SSH windows), 📡 *share terminal* (collaboration), minimize, maximize, close.
- **Snapping:** left, right, **top and bottom halves**, the four **quarters**, and **maximize**. Use the buttons in the title bar, the compact *snap layouts* menu in narrow windows, or the tab's right-click menu.
- **Drag-to-edge snapping:** drag a window with the pointer to an edge (half) or a corner (quarter) of the workspace. A blue preview shows where it will go.
- **Windows can never leave the workspace:** moving and resizing stop at all four edges. When the browser is resized or zoomed, or a panel changes size, windows are moved and shrunk to stay visible. Snapped and maximized windows keep their slot.
- **Resize** from any edge or corner. Dragging a snapped window restores its previous size.
- **Tabs:** click to focus or restore, **drag to reorder**, **middle-click** to close, **right-click** for *Reconnect*, *Duplicate window*, snap layouts and *Close*. SSH tabs show the connection state (🟡 connecting, 🟢 connected, 🔴 disconnected).
- **Close confirmation** can be *Always*, *Never* or *Only when connected* (Settings).
- The browser warns you before leaving the page while terminals are connected or a transfer is running.

### SSH terminal

- Full **xterm-256color** terminal (xterm.js 5) with the JetBrains Mono font. Works with `nano`, `vim`, `less`, `htop`, `mc`, `tmux`, colors, mouse reporting and UTF-8 (č, ć, ž, š, đ, emoji…).
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
  - optional per-file size limit (Settings)
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

### Sharing connections

Share connections with other WRM users or with people **without an account**:

- Right-click a connection → **🤝 Share**, a folder → **Share Folder**, or select several connections and use the bulk **🤝 Share** button.
- Options:
  - **Share name**
  - **Password lock** (optional; the recipient must enter it once and it is remembered for 7 days in their browser)
  - **Share with users** (comma-separated usernames; the share appears under *Shared with me* in their sidebar)
  - **Public link**: anyone with the URL can open it. Public shares are also listed under *Shared with me* for **every** user of this WRM server.
- You get a link like `https://your-wrm/share/<token>`. **🤝 Shares** in the top bar lists your shares, lets you copy links and **delete** shares, which revokes access immediately.
- People who open the link see the shared connections and can open **terminals and file managers** for them. WRM uses **your stored credentials**, so the recipients never see the password. They get the same access that account has on the server. Share only with people you trust, and prefer password-locked, non-public shares.

### Real-time collaboration

Everyone who opens the **same share link** (including you as the owner) joins a live room. A blue **Live Collab** bar shows who is online (you are marked *(you)*; guests get names like `Guest-4F2A1B`).

- **💬 Chat** with timestamps and join/leave messages.
- **📎 File exchange** in the chat (up to 2 MB per file).
- **📡 Terminal sharing:** press 📡 in the title bar of your SSH window. Others see a **👁 Watch** button next to your name, which opens a **read-only live copy of your terminal**. People who join later can watch too.
- **🎮 Remote control:** while you share a terminal, a *🎮 Control* button appears next to every other participant. Granting it lets that person type into **your** terminal through their viewer window (a green banner tells them they have control). *🔒 Revoke* takes it back. Control also ends automatically when you stop sharing or leave. Only the person sharing a terminal can grant control over it, and keystrokes go only to that person.

### Settings

**⚙ Settings** in the top bar:

| Setting | What it does |
|---|---|
| Language / Jezik | English or Croatian (Hrvatski) interface |
| Terminal Font Size | 11–18 px for all terminals |
| Close Window Confirmation | Always ask / Never ask / Only when connected |
| Terminal scrollback | 5,000 – 50,000 lines of history |
| Auto-reconnect SSH | On/Off: automatic reconnect after connection drops |
| Accent color | Blue, Violet, Emerald, Amber, Pink or Cyan |
| Max upload size | Per-file limit in MB (0 = unlimited) |
| Export Config | Download your folders and connections as JSON (**contains passwords/keys in plain text**; keep it safe) |
| Import Config | Import such a JSON file. Folders are matched by name, connections are added |
| Admin Panel | (admins only) user management |

Settings are stored per browser.

### Admin panel

Settings → **Admin Panel** (administrators only):

- List of all users with admin flag and join date.
- **Make Admin / Remove Admin**. You cannot remove your own admin rights.
- **Delete user**: removes the account, its login sessions and its shares. Its connections, folders and sessions are no longer accessible. You cannot delete yourself.

The first registered account is always an administrator.

### Mobile & responsive UI

- The top bar adapts to the available width (labels collapse into icons).
- On phones (< 760 px): the connection list and the Sessions/Clipboard panel become **slide-in drawers** (☰ and 🗂 buttons), windows open **full screen**, and a single tap opens a connection.
- Toolbars inside windows adapt to the window width.

---

## Keyboard shortcuts

| Where | Shortcut | Action |
|---|---|---|
| Anywhere | **Alt+N** | New window |
| Anywhere | **Esc** | Close menus / dialogs / editor |
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

## Configuration (environment variables)

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Listening port |
| `LISTEN_ADDR` | `:PORT` | Full listen address, e.g. `127.0.0.1:8080` (overrides `PORT`) |
| `DB_PATH` | `./remote_manager.db` | SQLite database file |
| `ENCRYPTION_KEY` | built-in default | Key for encrypting stored connection passwords/keys (first 32 bytes are used). **Set your own before adding connections.** Changing it later makes existing stored secrets unreadable. |
| `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE` | – | Serve HTTPS directly with this certificate and key |
| `WRM_ALLOWED_ORIGINS` | – | Extra hostnames allowed to open WebSockets (comma-separated), e.g. when a proxy rewrites the `Host` header |
| `WRM_ALLOW_ANY_ORIGIN` | – | `1` disables the WebSocket origin check (not recommended) |

---

## Running as a service & behind a reverse proxy

**systemd (Linux)**, e.g. `/etc/systemd/system/wrm.service`:

```ini
[Unit]
Description=Web Remote Manager PRO
After=network-online.target

[Service]
User=wrm
WorkingDirectory=/opt/wrm
Environment=PORT=8080
Environment=LISTEN_ADDR=127.0.0.1:8080
Environment=ENCRYPTION_KEY=choose-your-own-32-character-key!
ExecStart=/opt/wrm/wrm-pro-v9.10.2-mimo-linux-amd64
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload && sudo systemctl enable --now wrm
```

**nginx** in front of WRM (TLS, WebSockets, streaming and big uploads):

```nginx
server {
    listen 443 ssl;
    server_name wrm.example.com;
    # ssl_certificate ...; ssl_certificate_key ...;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;                # required: WebSocket origin check
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 3600s;                   # long-lived terminals
        proxy_buffering off;                        # live transfer/search progress
        client_max_body_size 0;                     # large uploads
    }
}
```

---

## Security model

Please read this before exposing WRM to a network:

- **Who can register:** anyone who can reach the WRM URL can create an account. Run WRM on a private network or VPN, or put it behind a reverse proxy with its own authentication, if the URL is reachable from the internet.
- **Accounts:** bcrypt password hashes, HTTP-only session cookies, login rate limiting.
- **Stored secrets:** connection passwords and private keys are AES-256-GCM encrypted in the SQLite file. The key is `ENCRYPTION_KEY`: **set your own**. With the built-in default key, anyone who gets your database file can decrypt the secrets. Protect the database file and backups. The owner of a connection can see its password in the *Edit* dialog, and *Export Config* writes secrets in plain text.
- **SSH host keys are not verified** (like `StrictHostKeyChecking=no`). Use WRM on networks you trust, or through a VPN, to avoid man-in-the-middle attacks. FTPS certificates are not verified either.
- **Transport:** use HTTPS (`HTTPS_CERT_FILE`/`HTTPS_KEY_FILE` or a reverse proxy). Otherwise passwords and terminal traffic between the browser and WRM travel unencrypted.
- **WebSockets** only accept connections whose `Origin` matches the host, which protects against cross-site WebSocket hijacking.
- **Sharing:** a share link gives terminal and file access **with your stored credentials** to everyone who can open it. Prefer password-locked, user-specific shares, and delete shares you no longer need.
- **Collaboration control:** only the participant who shares a terminal can grant keyboard control, and can revoke it at any time.
- **Browser dependencies:** the UI loads xterm.js and fonts from public CDNs (jsdelivr, Google Fonts), so the browser needs internet access for the terminal to load.

---

## Limitations

To be clear about what WRM does **not** do (yet):

- No **RDP/VNC** (graphical desktops): WRM is for SSH, SFTP, FTP and FTPS.
- No **voice/audio calls** in collaboration. Chat, file exchange and terminal sharing are text-based.
- No **2-factor authentication** and no option to disable self-registration yet.
- No SSH **agent forwarding**, **port forwarding/tunnels** or **jump hosts** yet.
- Server-to-server transfer and ZIP download work with SSH/SFTP servers only (not FTP).
- *File contents* search needs `grep` on the server.

Ideas and pull requests for any of these are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## API reference

All endpoints (except login/register, version and share pages) need the session cookie. JSON in, JSON out, unless noted.

**Auth**

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/auth/register` | Create account `{username, password}` (first user = admin) |
| POST | `/api/auth/login` | Log in `{username, password}` |
| POST | `/api/auth/logout` | Log out |
| GET | `/api/auth/me` | Current user |
| GET | `/api/auth/keepalive` | Keep the login session alive |

**Connections, folders, sessions**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/connections` | List / create connections |
| PUT / DELETE | `/api/connections/{id}` | Update / delete |
| POST | `/api/connections/bulk` | `{action: "move"\|"delete", ids, folder_id}` |
| POST | `/api/connections/test` | Test a connection (SSH/SFTP/FTP login) |
| GET / POST | `/api/folders` | List / create folders |
| DELETE | `/api/folders/{id}` | Delete folder (connections move to root) |
| GET / POST | `/api/sessions` | List / create workspace sessions (`layout` = JSON of windows) |
| PUT / DELETE | `/api/sessions/{id}` | Update (name, lock, layout) / delete (`?force=true` for locked) |
| GET | `/api/config/export` | Export folders + connections (JSON, secrets in plain text) |
| POST | `/api/config/import` | Import that JSON |

**Remote files** (`id` = connection id; for shared connections add `share_token`)

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/remote/list?id&path` | List a folder |
| GET | `/api/remote/list-recursive?id&path` | All files below a path |
| GET | `/api/remote/search?id&path&q&mode=name\|content&max&hidden` | Search (NDJSON stream) |
| GET | `/api/remote/download?id&path` | Download a file (`inline=1` for the editor) |
| GET | `/api/remote/download-dir?id&path` | Download a folder as ZIP (SFTP) |
| POST | `/api/remote/upload?id&path&overwrite` | Multipart upload, streamed; `relpath` field before a file part = subfolder |
| POST | `/api/remote/mkdir?id&path` | Create folder |
| POST | `/api/remote/rename?id&old&new` | Rename / move |
| DELETE | `/api/remote/delete?id&path&dir=0\|1` | Delete file or folder (recursive) |
| POST | `/api/remote/transfer` | Server-to-server transfer (NDJSON progress stream) |

**Sharing & admin**

| Method | Endpoint | Description |
|---|---|---|
| GET / POST | `/api/shares` | Your shares / create a share |
| DELETE | `/api/shares/{id}` | Delete a share |
| GET / POST | `/api/share/{token}` | Share info / unlock with password |
| GET | `/api/shared` | Shares visible to you ("Shared with me") |
| GET | `/api/users` | Usernames (for sharing) |
| GET | `/api/admin/users` | All users (admin) |
| PUT / DELETE | `/api/admin/users/{id}` | Toggle admin / delete user (admin) |
| GET | `/api/version` | Server version |

**WebSockets**

| Endpoint | Description |
|---|---|
| `/ws/ssh?id&cols&rows[&share_token]` | Terminal. Binary frames = terminal data; text frames = JSON control (`resize`, `pause`, `resume`, `ping` → server; `status`, `error`, `exit` ← server). Close codes: 1000 shell exited, 4001 connect failed, 4002 SSH connection lost |
| `/ws/events` | Live notifications (`sessions_changed`) |
| `/ws/share/{token}?name` | Collaboration room: `chat`, `file`, `screen`, `grant-control`, `revoke-control`, `remote-input`, `presence`, `welcome` |

---

## Building from source & releases

Requirements: **Go** (version in `remote-manager/go.mod`). No C compiler is needed (pure-Go SQLite).

```bash
git clone https://github.com/vedranius/web-browser-RDM-public.git
cd web-browser-RDM-public/remote-manager
go build -o wrm-server .
PORT=8080 ./wrm-server
```

Cross-compile any platform, e.g. `GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o wrm-linux-arm64 .`, or build all at once with `bash remote-manager/build-all.sh`.

**Project structure**

```
remote-manager/
  main.go         HTTP server, auth, connections, folders, sessions, shares, collaboration, transfer
  ssh_ws.go       SSH terminal WebSocket (binary frames, flow control, keepalive, reconnect codes)
  files.go        File API: list, download, ZIP, streamed upload, mkdir, rename, delete
  search.go       Recursive name / content search (NDJSON stream)
  sftp_pool.go    Pooled SFTP connections with liveness checks
  conntest.go     "Test connection" endpoint
  helpers.go      Origin check, login rate limiting, shared helpers
  static/index.html   The entire web UI (embedded into the binary)
.github/workflows/build.yml   CI: vet + JS syntax check, builds all platforms, releases on tags
```

**CI/CD:** every push builds all 13 targets. Pushing a tag `v*` creates a GitHub Release with all binaries, a macOS universal binary and `SHA256SUMS.txt`.

---

## Troubleshooting

| Problem | Fix |
|---|---|
| Terminal says *Access denied* or does not connect behind a proxy | Pass the `Host` header (`proxy_set_header Host $host;`) and WebSocket upgrade headers, or set `WRM_ALLOWED_ORIGINS` |
| Terminal area stays empty | The browser must be able to load xterm.js from `cdn.jsdelivr.net` |
| *Connection failed: unable to authenticate* | Check user/password/key with **Test connection**; for *Key file*/*Auto* the key must exist on the **WRM server** |
| Stored passwords stopped working after a restart | `ENCRYPTION_KEY` changed. Set the old key again, or re-enter the passwords |
| *File contents* search fails | The server needs `grep` and a POSIX shell; use *File names* search instead |
| Uploads fail at a certain size behind nginx | Set `client_max_body_size 0;` (or a larger value) |
| Transfer/search progress appears only at the end | Disable proxy buffering (`proxy_buffering off;`) |

---

## Contributing

Bug reports, ideas, translations and pull requests are welcome! Please read **[CONTRIBUTING.md](CONTRIBUTING.md)**. It explains how to build, test and submit changes, and the contributor terms.

If WRM is useful to you, consider supporting it on **[Ko-fi](https://ko-fi.com/vedranius)** ☕.

---

<sub>Web Remote Manager PRO · © 2026 vedranius · [PolyForm Noncommercial 1.0.0](LICENSE) · commercial use: [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)</sub>
