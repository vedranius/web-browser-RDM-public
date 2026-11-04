## Web Remote Manager PRO v14.0.0 — Storages, and everything since v13.14.10

One combined release of the parallel work: a new product area, its release review, and the security and reliability work of v13.14.11.

### New: Storages (S3)

- S3-compatible object storage (Hetzner Object Storage, DigitalOcean Spaces, Amazon S3, MinIO, Backblaze B2) as a section of its own, next to the servers. Each storage is confined to a prefix; its secret is encrypted at rest and never returned; a self-signed endpoint is trusted only by its pinned SHA-256 fingerprint.
- Browse, download, preview images and play video with a subtitle track; edit small text objects (JSON, configuration, subtitles up to 2 MB) with BOM and line endings kept and a conflict check.
- Write, upload and delete are allowed per storage; a delete needs an explicit confirmation and step-up where the security level asks for it. Every change is audited.
- Read-only, expiring, revocable share links to one prefix; `storage.down` / `storage.up` notifications; export and import.
- AI connections reach storages only through the opt-in scopes `storage_read`, `storage_write` and `storage_delete`; changes need the automatic mode.

### Fixed in the release review

- Test connection on a saved storage, saving a storage with a prefix-scoped key, renaming objects and server-side encryption on AWS, and the 2 MB editor limit.
- The storage share page, download names, egress checks and the kill switch were hardened, and an AI connection in the ask mode can no longer change a storage without an approval.

### Security and reliability (from v13.14.11)

- Share links keep bearer tokens out of URLs (HttpOnly cookie per share) and are validated before collaboration mode opens.
- Saved connection credentials stay bound to their destination; AI approvals are rechecked before they run; MFA checks are atomic; MCP OAuth consent needs an explicit start, and grant edits need step-up.
- Outbound requests and addon channels are tighter; Windows recursive-delete bypasses are blocked for the AI; network tool audit entries no longer hold credentials.
- Claude MCP reconnects after a WRM restart, AI access lists have *Select all*, pending AI approvals raise desktop notifications, and file, tunnel, session and folder actions report failures.

### Maintenance

- `pkg/sftp` 1.13.12, `modernc.org/sqlite` 1.61.0, updated CI actions, container base `node:26-alpine3.22`; the nightly ZAP scan signs in again.

### Next

- **Chains**: live monitoring of a whole service chain (servers, storage, configuration and the active path) with an interactive diagram and notifications, built from a workflow described to WRM's AI. Accepted in the roadmap; not part of this release.

### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.0.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.0.0) asset. Preserve the existing database and encryption key. Two tables are added (an older binary still starts on the database). To let an AI connection use a storage, add the storage scopes to `ai_mcp_scopes` and then to that connection.

[Full changelog](CHANGELOG.md)

## Included from v13.14.11 — secure shares and reliable AI connections

- Claude MCP reconnects after a WRM service restart, preserving its saved authorization and server selections.
- AI connection consent lists include Select all, with individual server deselection.
- Share links keep bearer tokens out of request URLs, use per-share HttpOnly cookies, and validate access before entering collaboration mode.
- Security hardening covers saved credentials in connection tests, approval rechecks, MFA, MCP OAuth consent, outbound requests, addon channels, Windows destructive commands, and network audit logs.
- Pending AI approvals raise desktop notifications that link to the matching approval.
- File, tunnel, session, and folder actions report request failures clearly.
