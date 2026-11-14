---
id: connections
title: Connections and folders
section: Basics
order: 20
summary: Add, edit, organise, move and search your servers; protocols, tags and notes.
keywords: connection, server, add, new, edit, folder, move, drag, order, sort, tag, search, ssh, rdp, vnc, telnet, web, quick connect, notes, status
ui: #sidebar, #modal-conn, #modal-folder, #modal-move, #modal-quick, #modal-notes
open: new_connection
---
A connection is one saved server: protocol, address, login and options. Connections live in folders in the left sidebar.

## Add a connection
1. **+ Connection** in the sidebar.
2. Choose the **protocol**: SSH (terminal and files), SFTP/FTP/FTPS (files only), RDP, VNC, Telnet (remote desktop), or **Web interface** (HTTP/HTTPS page of a switch, iDRAC, NAS…).
3. Enter **host**, **port** and **user**. For the login choose a password, a key from **🔑 Keys**, a shared **credential**, or a Bitwarden SSH key.
4. Optional: **folder**, **tags** (`env:prod`, `site:dc1`), a **jump host** or **proxy**, SSH **tunnels**, **notes**, out-of-band management (BMC).
5. **Save**, then double-click the connection to open it.

## Quick connect
Type `user@host:port` (or `rdp://…`, `vnc://…`, `https://…`) into **⚡ Quick connect** and press Enter. It is kept 24 hours under *Quick connections*; 💾 saves it for good.

## Folders and order
- **📁 Folder** creates a folder; right-click a folder for a subfolder, rename, share, sort A–Z or *keep sorted*.
- **Drag** connections and folders to move them; drop between two rows to put it exactly there. Ctrl+click selects several, then drag them together or use right-click → **Move to…**.
- Right-click a connection: open, edit, duplicate, share, network tools, tunnels, notes, delete.

## Find a server
- The **search** box filters by name, host, user and tags (`#env:prod`).
- The **🏷 tag bar** filters by tags; `env:prod` connections get a red **PROD** badge.
- The light next to each connection: 🟢 up, 🟡 slow, 🔴 down (hover for why). A red bar lists only the servers that are down.

## Notes
📝 on a connection holds procedures, contacts and links (small Markdown). They are shown in the terminal window too.
