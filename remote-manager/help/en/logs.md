---
id: logs
title: Logs
section: Operations
order: 75
summary: Find every log on a server, view and follow them, search many logs at once and ask the AI log assistant.
keywords: log, logs, journal, journalctl, syslog, tail, follow, search, grep, error, discovery, docker logs, supervisor, log assistant
ui: .lv-wrap, .logviewer
open: monitoring:logs
---
**📊 Monitoring → Logs**, then choose a server.

## Discovery
WRM finds the logs by itself over the connection's SSH account (no sudo): systemd units and their journal, Supervisor programs, Docker containers and log files under the scan paths. Code files and Python virtual environments are skipped. The list can be sorted and filtered by service, operator, folder and name.

## View and follow
Click a file (or a journal unit) to open it at the end; **Follow** shows new lines live. Filter by level (errors, warnings…), *Find in this log* (`/`, then `n`/`N`), *Jump to time*, *Load older lines*, copy or download the shown lines, or *Open in terminal*.

## Search many logs at once
Tick files (or none for all) and search: WRM searches them on the server (text or regular expression, case, whole words, context lines, a time range) and shows the hits per file with a timeline; click a bar to zoom in, click a hit to open that place.

## AI log assistant
**✦ AI** in the log viewer opens a chat that only reads: *why did demo-api fail last night?* It searches all discovered logs, quotes the lines with file and time, and suggests what to check next. It never changes anything.

## Log settings
**Settings** in the log viewer: scan paths, files to exclude, how long a discovery stays valid, and for administrators the limits for everybody.
