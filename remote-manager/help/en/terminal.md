---
id: terminal
title: SSH terminal, snippets and broadcast
section: Basics
order: 30
summary: Work in terminals; search, reconnect, snippets, run on connect, bookmarks and typing into many servers.
keywords: terminal, ssh, console, shell, snippet, command, broadcast, multi, cssh, run on connect, bookmark, search, copy, paste, font, reconnect, recording
ui: .term-wrap, .xterm, #modal-snippet, #modal-broadcast, #modal-bookmark, #snip-list
open: snippets
---
Double-click an SSH connection to open a terminal. It is a full xterm (256 colours, mouse, UTF-8) and works with `vim`, `htop`, `mc` and `tmux`.

## Everyday keys
- **Ctrl+Shift+F** search in the terminal output; **Ctrl+Shift+C / V** copy and paste (selecting text copies it to the Clipboard panel).
- **Enter** in a closed terminal reconnects; **↻** in the title bar does it any time.
- **Ctrl+Shift+Space** (or ⚡ in the title bar) opens the snippet picker.

## Snippets
Saved commands with variables such as `{{host}}`, `{{user}}` or `{{?Service=nginx}}` (asked when used).
1. **Snippets** panel → **＋**, or select text in a terminal and press **＋ New snippet** in the picker.
2. Click a snippet to run it in the focused terminal; **Shift+click** inserts it without Enter.
3. Scope: all connections, one folder or one connection. **Run on connect** snippets run after login (for example `sudo -i`).

## Folder bookmarks
★ in the title bar `cd`s to a saved directory such as `/opt/app/{name}/logs`; in the file manager it opens it.

## Broadcast input
Type the same thing into many terminals: **🧰 Tools → 📣 Broadcast…**, tick the terminals, **Start**. Dangerous commands (`rm -rf`, `shutdown`, `mkfs`…) and pasted multi-line text ask first. **■ Stop** ends it.

## AI in the terminal
**AI** in the title bar opens the assistant next to this terminal (see *AI assistant*). To ask about something on screen without a server session, use **?** → **⌖ Pick** and click the terminal.

## Recording
Sessions are recorded when the policy says so; replay them in **Settings → Session history**.
