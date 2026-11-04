## Web Remote Manager PRO v13.15.0 — Storages (S3)

- **Storages:** S3-compatible object storage (Hetzner, DigitalOcean Spaces, Amazon S3, MinIO, Backblaze B2) as its own section, confined to a prefix, with secrets encrypted at rest and certificate pinning for self-signed endpoints.
- Browse, download, preview images and play video with subtitles; edit small text objects (JSON, configuration, subtitles up to 2 MB) with BOM/CRLF kept and a conflict check.
- Write, upload and delete are allowed per storage; a delete needs an explicit confirmation and step-up where required. Every change is audited.
- Read-only, expiring, revocable share links to one prefix; `storage.down` / `storage.up` notifications; export/import.
- AI connections reach storages only through the opt-in scopes `storage_read`, `storage_write` and `storage_delete`.
- Dependency and CI action updates.

### Upgrade

Use WRM's update action or replace the executable/container with the matching [v13.15.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v13.15.0) asset. Preserve the existing database and encryption key. Two tables are added; to let an AI connection use a storage, add the storage scopes to `ai_mcp_scopes` and then to that connection.

[Full changelog](CHANGELOG.md)

## Included from v13.14.11 — secure shares and reliable AI connections

- Claude MCP reconnects after a WRM service restart, preserving its saved authorization and server selections.
- AI connection consent lists include Select all, with individual server deselection.
- Share links keep bearer tokens out of request URLs, use per-share HttpOnly cookies, and validate access before entering collaboration mode.
- Security hardening covers saved credentials in connection tests, approval rechecks, MFA, MCP OAuth consent, outbound requests, addon channels, Windows destructive commands, and network audit logs.
- Pending AI approvals raise desktop notifications that link to the matching approval.
- File, tunnel, session, and folder actions report request failures clearly.
