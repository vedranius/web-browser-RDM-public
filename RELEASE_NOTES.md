## Web Remote Manager PRO v9.10.2-mimo

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/vedranius)

Browser-based remote server management: SSH terminal, SFTP/FTP/FTPS file manager, server-to-server transfer, sharing and real-time collaboration — one self-contained binary with an embedded web UI.

### What's new in v9.10.2-mimo — collaboration security fix

- **Security:** in a shared (collaboration) session any participant, including a guest without an account, could give *themselves* keyboard control and type into another participant's shared SSH terminal. Now only the participant who is sharing a terminal can grant control. Keystrokes go only to that sharer. Control is revoked automatically when the sharer stops sharing or leaves. The browser also accepts input only from people it granted control to.
- **Remote control works:** the *🎮 Control* button now appears next to the other participants while you share a terminal. Before, it only appeared next to people who were sharing themselves, so control could never be given to a viewer.
- Participants who join after someone started sharing now see the *👁 Watch* button right away.
- Duplicate names in a room (the same user in two tabs or two guests) get a unique suffix, e.g. `tester (2)`; you are marked with *(you)*.

### Included from v9.10.1-mimo

**Saved sessions really restore your workspace**
- *+ Save Session* now stores every open terminal / file-manager window: connection, mode, folder, position or snap slot, minimized and focused state. Before, only the session **name** was saved, so opening a saved session did nothing.
- Click a saved session (in the Sessions panel or in an empty *New Window*) to open it. If windows are already open you choose **Replace** or **Add**.
- Right-click a session: *Open (replace)*, *Open (add)*, **Save current windows here**, rename, lock, delete. Locked sessions cannot be overwritten. Positions are stored relative to the window area, so they also fit other screen sizes.

**Reconnect button on every SSH window**
- New **↻** button in the title bar of SSH windows, also when connected. It is useful when a session hangs. Shortcut: **Ctrl+Shift+R** inside the terminal. The same action is in the tab right-click menu.

**Windows stay inside the screen**
- Windows can no longer be dragged or resized past **any** edge (top, bottom, left, right). If the browser is zoomed, the resolution or panel sizes change, or a session from a bigger screen is opened, windows are moved and shrunk to fit automatically. The floating transfer panel stays on screen too.

**Snap**
- New snap targets **top half** and **bottom half**, next to left/right halves, the four quarters and maximize. All snap buttons have clear layout icons.
- **Drag-to-edge snapping**: drag a window with the pointer to an edge or corner of the window area. A preview shows the slot, and releasing snaps the window there. Edges give halves, corners give quarters.
- Narrow windows show a compact **snap layouts** menu instead of the full button row.
- Tab right-click menu: *Reconnect*, *Duplicate window*, snap layouts, close.

---

## Included from v9.10.0-mimo

### 🔧 Fixed (v9.10.0)

- **"Connection error / disconnected" when a command prints a lot of output** (e.g. `cat` / `tail -f` of a big log with č, ć, ž, š, đ or emoji). Terminal data is now sent as binary WebSocket frames. Before, a 32 KB chunk could split a multi-byte UTF-8 character and the browser closed the connection.
- **nano / vim / less / htop drawn in a small part of the window until the browser was zoomed.** The terminal size is now sent to the server when connecting and on every change: maximize, snap, window resize, sidebar toggle, browser zoom, font size. Before, maximize/snap/resize only resized the terminal in the browser.
- Server-to-server transfer of folders failed when WRM runs on **Windows** (remote paths were built with `\`).
- SSH session stayed "open" after typing `exit`; pasting a lot of text could drop input; terminal modes were 7-bit (now CS8 + IUTF8).
- Collaboration server could crash (concurrent WebSocket write) when granting/revoking control.
- Event WebSocket leaked subscribers and was reachable without login.
- The "Edit Connection" dialog title showed the raw key `edit_connection_title`.
- Download of files with quotes or non-ASCII characters in the name.
- File sizes above 1 GB were shown in MB only.

### ✨ New (v9.10.0)

**Terminal**
- **Reconnect**: banner with a *Reconnect* button, or press **Enter** in a closed terminal.
- **Automatic reconnect** with backoff after SSH or network drops, and when the network comes back (Settings → Auto-reconnect).
- Connection state indicator (connecting / connected / disconnected) on windows, tabs and sidebar connections.
- Flow control for very fast output; SSH keepalive detects dead connections.
- Search in terminal output (**Ctrl+Shift+F**), clickable URLs, configurable scrollback.

**File manager**
- **Search**: type to filter the current folder instantly; press **Enter** to search all subfolders **by name** (wildcards: `*.log`, `nginx*.conf`) or **by content** (grep on the server). Results stream live, can be stopped, and double-click jumps to the file.
- Breadcrumb path, *Modified* column, file icons, status bar with selection size.
- Multi-select (click, Ctrl/Shift+click, Ctrl+A) and keyboard: arrows, Enter, Backspace, Delete, F2 (rename), F5, Ctrl+F.
- **Upload many files or whole folders**, **drag & drop** from the desktop, progress bar with speed, overwrite confirmation. Uploads are streamed to the server (no RAM buffering).
- **Built-in text editor** (Ctrl+S, keeps LF/CRLF).
- Context menu: *Open terminal here*, *Copy path*, bulk download / delete / transfer.
- Pooled SFTP connections → much faster browsing.

**Connections & UI**
- **Test connection** button, **Duplicate** connection, **FTPS** (explicit TLS).
- Modern dark UI, configurable accent color, modern dialogs.
- **Responsive**: the top bar adapts to the space available; on phones the sidebar and sessions become slide-in drawers and windows open full screen; tap to open connections.
- Double-click a window title to maximize; dragging a maximized window restores it; maximized/snapped windows follow layout changes.
- Setting for max upload size per file.

**Security**
- Login rate limiting (8 failures per IP → 5 min lockout).
- WebSocket origin check (`WRM_ALLOWED_ORIGINS`, `WRM_ALLOW_ANY_ORIGIN=1` to disable).
- Optional HTTPS with `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE`.

### 📦 Downloads

| Platform | Architecture | Binary |
|---|---|---|
| Linux | x86-64 | `wrm-pro-v9.10.2-mimo-linux-amd64` |
| Linux | arm64 | `wrm-pro-v9.10.2-mimo-linux-arm64` |
| Linux | ARMv7 (Raspberry Pi) | `wrm-pro-v9.10.2-mimo-linux-armv7` |
| Linux | ARMv6 | `wrm-pro-v9.10.2-mimo-linux-armv6` |
| Linux | 32-bit | `wrm-pro-v9.10.2-mimo-linux-386` |
| Windows | x86-64 | `wrm-pro-v9.10.2-mimo-windows-amd64.exe` |
| Windows | arm64 | `wrm-pro-v9.10.2-mimo-windows-arm64.exe` |
| macOS | Intel | `wrm-pro-v9.10.2-mimo-darwin-amd64` |
| macOS | Apple Silicon | `wrm-pro-v9.10.2-mimo-darwin-arm64` |
| macOS | Universal | `wrm-pro-v9.10.2-mimo-darwin-universal` |
| Android | arm64 (Termux) | `wrm-pro-v9.10.2-mimo-android-arm64` |
| FreeBSD | x86-64 | `wrm-pro-v9.10.2-mimo-freebsd-amd64` |
| FreeBSD | arm64 | `wrm-pro-v9.10.2-mimo-freebsd-arm64` |
| OpenBSD | x86-64 | `wrm-pro-v9.10.2-mimo-openbsd-amd64` |

Verify integrity with `SHA256SUMS.txt`.

### 🚀 Quick start

**Linux / macOS**
```bash
chmod +x wrm-pro-v9.10.2-mimo-linux-amd64
./wrm-pro-v9.10.2-mimo-linux-amd64
# open http://localhost:8080 — the first registered user becomes admin
```

**Windows** — double-click `wrm-pro-v9.10.2-mimo-windows-amd64.exe`, or in PowerShell:
```powershell
$env:PORT=9000; .\wrm-pro-v9.10.2-mimo-windows-amd64.exe
```

**Android (Termux)**
```bash
pkg install wget
wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v9.10.2-mimo/wrm-pro-v9.10.2-mimo-android-arm64
chmod +x wrm-pro-v9.10.2-mimo-android-arm64
PORT=8080 ./wrm-pro-v9.10.2-mimo-android-arm64
```

### ⚙️ Environment variables

```
PORT=8080                  # listening port (default 8080)
LISTEN_ADDR=:8080          # full listen address (overrides PORT)
DB_PATH=./remote_manager.db
ENCRYPTION_KEY=...         # 32-char key for stored connection secrets — set your own!
HTTPS_CERT_FILE=...        # enable HTTPS (together with HTTPS_KEY_FILE)
HTTPS_KEY_FILE=...
WRM_ALLOWED_ORIGINS=a.example.com   # extra hosts allowed to open WebSockets (reverse proxies)
WRM_ALLOW_ANY_ORIGIN=1     # disable the WebSocket origin check
```

> **Upgrading behind a reverse proxy:** WebSockets now require the browser's `Origin` to match the `Host` (or `X-Forwarded-Host`) header. If terminals stop connecting after the upgrade, pass the Host header through (`proxy_set_header Host $host;` in nginx) or set `WRM_ALLOWED_ORIGINS`.

---

### 📜 License

Web Remote Manager PRO is **source-available** under the [PolyForm Noncommercial License 1.0.0](https://github.com/vedranius/web-browser-RDM-public/blob/main/LICENSE): free for personal, educational, non-profit and other noncommercial use, including forks and modifications. **Commercial use requires a separate license**; see [COMMERCIAL-LICENSE.md](https://github.com/vedranius/web-browser-RDM-public/blob/main/COMMERCIAL-LICENSE.md). Contributions are welcome; see [CONTRIBUTING.md](https://github.com/vedranius/web-browser-RDM-public/blob/main/CONTRIBUTING.md).

☕ **Like WRM?** Support its development on **[Ko-fi](https://ko-fi.com/vedranius)** — thank you!
