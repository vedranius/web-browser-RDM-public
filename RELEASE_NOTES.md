## Web Remote Manager PRO v14.1.0 — Chains: live monitoring, a live diagram and an AI builder

Chains is a new core section: WRM watches a whole working system instead of single servers, draws it as a live diagram, and the built-in assistant can draft the chain from your description. Everything is additive; upgrading from v14.0.0 keeps the database and the encryption key.

### New: Chains

- **A chain** is the set of servers, storages and outside endpoints that form one system (for example ingest → transcoding → content → CDN, with a primary and a backup path). It holds nodes (saved SSH connections, storages, outside endpoints), checks per node from a fixed catalog, links (data, read, write, failover, dependency) and scenarios: ordered rules over the results, the first match being the active path.
- **Read-only checks, nothing installed:** systemd unit and timer, Docker container, process count, listening port, file and its age, file count in a directory, one configuration value (INI, JSON, YAML, env, `key = value` or a regular expression), mount (rw/ro), crontab line and WRM's monitoring state — all in one SSH round trip per server. From WRM: DNS (optional resolver), TCP, HTTP(S) status or a JSON value, TLS expiry, and a storage's bucket or its newest object under a prefix.
- **Background checks** every `interval_s` (default 30 s, minimum 15 s) using the same slots as monitoring. A change counts only after `confirm` sightings, so flapping does not notify; up to 2000 changes are kept per chain.
- **Live diagram:** automatic layout with positions saved when you drag a node; servers as boxes, storages as buckets, outside endpoints as pills; animated flows coloured by a check or node; a link's `scenarios` list shows the backup path idle; a scenario band over time; node and link details with check history; zoom, fit and pan; *Check now*.
- **Starter templates:** one service, and a primary + backup path with a bucket and a CDN.
- **Notifications:** `chain.scenario_changed`, `chain.check_changed`, `chain.node_down` and `chain.node_up` through the existing channels and Workboard; `chain_changed` refreshes an open window.
- **AI builder:** describe the system in WRM's assistant. `chain_status` shows your chains, connections, storages and the check catalog; `chain_draft` writes a validated draft using catalog checks and only your own connections and storages. A draft is not monitored and starts nothing: it appears under *Drafts from the AI* in Chains, opens on the same diagram, and you approve it with *Create chain* (or discard it). One draft per assistant session, audited as `chain.drafted`.

### Security

- Nothing is installed on a server and a chain can never run a command: the check script is fixed text and parameters reach it as data.
- A chain can use only its owner's connections and storages, configuration keys that look like secrets are refused, and only extracted values are kept.
- A draft is validated again when it is approved, and a draft cannot add a check type, a foreign connection or a command. The builder tools are available only in WRM's built-in assistant.

### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.1.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.1.0) asset. Preserve the database and encryption key. Three tables are added (`chains`, `chain_samples`, `chain_drafts`); an older binary still starts on the database. The **Chains** button appears in the top bar (kill switch `chains_enabled`, default interval `chain_interval_s`). AI connection scopes do not change.

[Full changelog](CHANGELOG.md)

## Included from v14.0.0 — Storages, and everything since v13.14.10

One combined release of the parallel work: a new product area, its release review, and the security and reliability work of v13.14.11.

#### New: Storages (S3)

- S3-compatible object storage (Hetzner Object Storage, DigitalOcean Spaces, Amazon S3, MinIO, Backblaze B2) as a section of its own, next to the servers. Each storage is confined to a prefix; its secret is encrypted at rest and never returned; a self-signed endpoint is trusted only by its pinned SHA-256 fingerprint.
- Browse, download, preview images and play video with a subtitle track; edit small text objects (JSON, configuration, subtitles up to 2 MB) with BOM and line endings kept and a conflict check.
- Write, upload and delete are allowed per storage; a delete needs an explicit confirmation and step-up where the security level asks for it. Every change is audited.
- Read-only, expiring, revocable share links to one prefix; `storage.down` / `storage.up` notifications; export and import.
- AI connections reach storages only through the opt-in scopes `storage_read`, `storage_write` and `storage_delete`; changes need the automatic mode.

#### Fixed in the release review

- Test connection on a saved storage, saving a storage with a prefix-scoped key, renaming objects and server-side encryption on AWS, and the 2 MB editor limit.
- The storage share page, download names, egress checks and the kill switch were hardened, and an AI connection in the ask mode can no longer change a storage without an approval.

#### Security and reliability (from v13.14.11)

- Share links keep bearer tokens out of URLs (HttpOnly cookie per share) and are validated before collaboration mode opens.
- Saved connection credentials stay bound to their destination; AI approvals are rechecked before they run; MFA checks are atomic; MCP OAuth consent needs an explicit start, and grant edits need step-up.
- Outbound requests and addon channels are tighter; Windows recursive-delete bypasses are blocked for the AI; network tool audit entries no longer hold credentials.
- Claude MCP reconnects after a WRM restart, AI access lists have *Select all*, pending AI approvals raise desktop notifications, and file, tunnel, session and folder actions report failures.

#### Maintenance

- `pkg/sftp` 1.13.12, `modernc.org/sqlite` 1.61.0, updated CI actions, container base `node:26-alpine3.22`; the nightly ZAP scan signs in again.

#### Next

- **Chains**: live monitoring of a whole service chain (servers, storage, configuration and the active path) with an interactive diagram and notifications, built from a workflow described to WRM's AI. Accepted in the roadmap; not part of this release.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.0.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.0.0) asset. Preserve the existing database and encryption key. Two tables are added (an older binary still starts on the database). To let an AI connection use a storage, add the storage scopes to `ai_mcp_scopes` and then to that connection.

[Full changelog](CHANGELOG.md)

## Included from v13.14.11 — secure shares and reliable AI connections

- Claude MCP reconnects after a WRM service restart, preserving its saved authorization and server selections.
- AI connection consent lists include Select all, with individual server deselection.
- Share links keep bearer tokens out of request URLs, use per-share HttpOnly cookies, and validate access before entering collaboration mode.
- Security hardening covers saved credentials in connection tests, approval rechecks, MFA, MCP OAuth consent, outbound requests, addon channels, Windows destructive commands, and network audit logs.
- Pending AI approvals raise desktop notifications that link to the matching approval.
- File, tunnel, session, and folder actions report request failures clearly.
