---
id: updates-backup
title: Updates, backups and installation
section: Administration
order: 170
summary: One-click signed updates, what to back up, running WRM as a service, in Docker and behind a reverse proxy.
keywords: update, upgrade, version, release, self-update, signature, backup, restore, database, key, service, systemd, windows service, docker, reverse proxy, nginx, https, certificate, port, install
ui: #modal-update, .app-ver
open: updates
admin: true
---
## Update
Click the version in the top bar. WRM shows what is new; administrators press **Update**: the release is downloaded, its **signature** checked, and WRM restarts with it. A Docker install is updated by pulling the new image.

## Back up
Keep copies of both files next to each other:
- `remote_manager.db` (everything: connections, users, settings, history),
- `remote_manager.db.key` (the key that decrypts stored passwords and keys). **Without it stored secrets are lost.**
Copy them while WRM is stopped, or use the SQLite backup of your choice. A previous version still starts on a newer database.

## Run it
- **As a service:** `wrm -install-service` (systemd, launchd, rc.d or a Windows service).
- **Docker:** `docker compose up -d` with the file from the repository (includes guacd).
- **Reverse proxy:** pass WebSockets through; set `ai_mcp_public_url` when the outside address differs.
The README on GitHub has every option (`PORT`, `DB_PATH`, TLS, environment variables).
