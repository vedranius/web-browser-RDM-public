## Web Remote Manager PRO v14.4.0 — S3 activity and costs, AI in the Chains builder

### S3 requests you can see

- **Live S3 activity panel** beside the storage browser. It lists every request WRM sends to the bucket, as it happens: the operation, key, answer and duration, and what caused it (you, the background check, a chain, an AI connection). Each operation is counted for today and in total. The panel is shown by default; ✕ or 📡 hides it.
- **Provider picker with request costs** in the storage form: Hetzner, DigitalOcean Spaces, Amazon S3, Backblaze B2, Cloudflare R2, MinIO and Other. Picking one shows whether, in which class and for how much the provider bills each S3 operation WRM uses, quoting and linking the provider's own price page:
  - Hetzner and Backblaze B2 bill none of them;
  - DigitalOcean's price list has no request fee;
  - Amazon S3 bills every request except DELETE;
  - Cloudflare R2 bills two classes after a monthly free allowance.

  The full table is in [docs/STORAGES.md](https://github.com/vedranius/web-browser-RDM-public/blob/v14.4.0/docs/STORAGES.md#what-each-provider-bills).

### AI in the Chains builder

- **🤖 AI on a chain or draft.** Say what to add, change or check. The AI reads the chain and its last results and writes a draft. A change to a chain waits as *Changes for …* until you press **Apply changes**; the chain keeps its id and history.
- **Test once** runs a draft's checks a single time and shows the results on the diagram, without saving.
- **New tools.** WRM's AI panel gets `chain_get` and `chain_test`, and `chain_draft` can change an existing chain. MCP clients get the same four tools behind the new opt-in **`chains`** scope: an administrator adds it to the allowed scopes, then a user enables it on an AI connection.

### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.4.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.4.0) asset. Preserve the database and encryption key. Two columns are added (`storages.provider`, `chain_drafts.base_chain_id`); an older binary still starts on the database.

[Full changelog](CHANGELOG.md)

## Included from v14.3.2 — Edit and remove storages, fewer S3 requests

A bug-fix release for the storage section: a saved storage can be edited and removed again, and WRM sends fewer background requests to the bucket.

#### Fixed

- **Edit and remove a storage.** Each storage in the list has a ✎ button and the open storage an *Edit* button. Both open the storage form, which also has *Remove storage*. Before this release the form could only be opened for a new storage. After saving, the open storage uses its new permissions straight away, and on a phone the list sits above the files.

#### Changed

- **Fewer background S3 requests.** The background bucket check (one `HEAD` per storage) runs every 15 minutes instead of 5; `storage_check_minutes` changes it and `0` turns it off. A chain's storage check asks the bucket at most every 5 minutes (`chain_storage_check_s`) instead of on every 30-second poll.
- **Renaming a storage or changing its permissions sends no S3 request**, so it works while the endpoint is down. Changes to the endpoint, region, bucket, prefix or keys are still checked against the bucket.
- [docs/STORAGES.md](docs/STORAGES.md) lists every S3 request WRM sends per action. Hetzner Object Storage and DigitalOcean Spaces Standard do not bill API requests.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.3.2](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.3.2) asset. Preserve the database and encryption key. No database change; two new settings, `storage_check_minutes` and `chain_storage_check_s`.

[Full changelog](CHANGELOG.md)

## Included from v14.3.1 — Log reads and MCP responsiveness

A small bug-fix release for two problems seen in daily use: the AI refusing an ordinary log file, and WRM becoming unreachable while Claude's MCP connection works in parallel with several browser tabs.

#### Fixed

- **A log glob is not a private key.** `follower.log*` (and other ordinary log globs) was refused as a PuTTY private key because it could theoretically expand to `follower.log.ppk`. Extension-only rules now apply to explicit paths; `server.ppk`, `.key`, `.p12` and `.pfx` files are still refused.
- **WRM stays reachable with parallel MCP calls and several tabs.** The database pool grew from 8 to 32 connections with idle expiry, and connection/folder lists load their labels in one query instead of one per row. Eight MCP tool calls plus browser tabs could previously hold every connection in a nested lookup, leaving the UI waiting until a service restart; closing tabs is no longer needed.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.3.1](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.3.1) asset. Preserve the database and encryption key. No database or setting change.

[Full changelog](CHANGELOG.md)

## Included from v14.3.0 — Complete storage folders, no share links, and the Chains fix

A storage folder is now loaded and counted completely, the browser gets search, sorting and bulk actions, the public S3 share links are removed for security, and the Chains wizard no longer fails with four or more servers.

#### Storage browser

- **The whole folder, not the first page.** The window follows every S3 continuation token, shows the running count while loading, then the total object and folder count and size with "all objects in this folder are shown"; above a 100-page safety limit it offers *Load more*.
- **Search, sort, select.** Search filters the complete folder ("showing X of Y"), sort by name, size or date ascending/descending, select-all and bulk delete work on the visible selection, and *Copy path* copies an object key.
- **Test connection before saving** checks the endpoint, bucket and keys from the add form (`POST /api/storages/test`); it no longer needs a saved storage.

#### Removed: S3 share links

- The `Share` button, `GET/POST /api/storages/{id}/shares`, `DELETE .../shares/{id}` and the public `/s3/{token}` page are gone. Links created by an older version stop working; storage access always needs a signed-in account with the storage's permissions. The unused `storage_shares` table may remain and can be dropped.

#### Fixed

- **Chains wizard with four or more servers** no longer fails with "check id must be 1-40 characters" or "too many rules in one list": ids are shortened with a collision suffix, and a large all-checks-OK state is nested in groups of 20.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.3.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.3.0) asset. Preserve the database and encryption key. No database or setting change; existing storages and chains are unchanged, but old share links no longer work.

[Full changelog](CHANGELOG.md)

## Included from v14.2.0 — Chains you can build, with an AI that reads the servers

Chains leaves JSON behind: a guided builder, a workflow-like diagram and a wizard that inspects the chosen servers before the assistant drafts anything. A storage is no longer required.

#### New: the guided builder

- *New chain* offers **Build with AI**, **Build step by step**, templates and advanced JSON. The guided editor picks servers, storages and outside services from the saved connections, adds checks with a form for each catalog type, connects nodes and edits states (scenarios) with a condition editor; *Normal / Backup / Degraded* can be generated.
- Drag from a node's connector dot to another node to create a link, or use *Connect to…* in the node panel (also on touch). The toolbar adds a node or a link; the node panel adds/removes checks; nodes and links can be deleted in place. The same works on AI drafts before approval.
- Scenario rules are shown as sentences, not JSON, and custom rules are preserved when the editor opens.

#### The AI wizard discovers the system first

- Pick servers (storages optional) and, if you want, describe the system. WRM runs one fixed, read-only discovery per server: OS, running non-system services, Docker containers, listening ports and network/object-storage mounts (rclone, s3fs, NFS, CIFS, SSHFS). The facts are shown in the wizard and given to the assistant, which uses real service names, ports and mounts as check targets.
- No AI provider, or an answer that fails validation, is not a dead end: the same facts build a deterministic starter.
- `POST /api/chains/ai-draft` and `POST /api/chains/discover` are the new endpoints; drafts can be edited with `PUT /api/chains/drafts/{id}`.

#### No storage required, nothing attached

- A chain of servers is complete on its own. Templates leave an unpicked storage node out, and the wizard accepts zero storages.
- A discovered network/object-storage mount becomes a `mount` check on its server; when the source names a reachable host, an outside endpoint with a DNS check is added. WRM does not mount it, does not create a storage entry, never asks for its credentials and never uses a `storage_id` for it. Storage nodes remain only for storages you already saved.

#### Security

- Discovery is fixed text over the selected connections, read-only and bounded: names, ports and mount points only — no file contents, environment variables or credentials.
- Every node must use a connection or storage the caller owns; the AI answer passes the same `validateChain` and facts are treated as untrusted data in the prompt.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.2.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.2.0) asset. Preserve the database and encryption key. No database or setting change; existing chains and the v14.1.0 tables are unchanged.

[Full changelog](CHANGELOG.md)

## Included from v14.1.0 — Chains: live monitoring, a live diagram and an AI builder

Chains is a new core section: WRM watches a whole working system instead of single servers, draws it as a live diagram, and the built-in assistant can draft the chain from your description. Everything is additive; upgrading from v14.0.0 keeps the database and the encryption key.

#### New: Chains

- **A chain** is the set of servers, storages and outside endpoints that form one system (for example ingest → transcoding → content → CDN, with a primary and a backup path). It holds nodes (saved SSH connections, storages, outside endpoints), checks per node from a fixed catalog, links (data, read, write, failover, dependency) and scenarios: ordered rules over the results, the first match being the active path.
- **Read-only checks, nothing installed:** systemd unit and timer, Docker container, process count, listening port, file and its age, file count in a directory, one configuration value (INI, JSON, YAML, env, `key = value` or a regular expression), mount (rw/ro), crontab line and WRM's monitoring state — all in one SSH round trip per server. From WRM: DNS (optional resolver), TCP, HTTP(S) status or a JSON value, TLS expiry, and a storage's bucket or its newest object under a prefix.
- **Background checks** every `interval_s` (default 30 s, minimum 15 s) using the same slots as monitoring. A change counts only after `confirm` sightings, so flapping does not notify; up to 2000 changes are kept per chain.
- **Live diagram:** automatic layout with positions saved when you drag a node; servers as boxes, storages as buckets, outside endpoints as pills; animated flows coloured by a check or node; a link's `scenarios` list shows the backup path idle; a scenario band over time; node and link details with check history; zoom, fit and pan; *Check now*.
- **Starter templates:** one service, and a primary + backup path with a bucket and a CDN.
- **Notifications:** `chain.scenario_changed`, `chain.check_changed`, `chain.node_down` and `chain.node_up` through the existing channels and Workboard; `chain_changed` refreshes an open window.
- **AI builder:** describe the system in WRM's assistant. `chain_status` shows your chains, connections, storages and the check catalog; `chain_draft` writes a validated draft using catalog checks and only your own connections and storages. A draft is not monitored and starts nothing: it appears under *Drafts from the AI* in Chains, opens on the same diagram, and you approve it with *Create chain* (or discard it). One draft per assistant session, audited as `chain.drafted`.

#### Security

- Nothing is installed on a server and a chain can never run a command: the check script is fixed text and parameters reach it as data.
- A chain can use only its owner's connections and storages, configuration keys that look like secrets are refused, and only extracted values are kept.
- A draft is validated again when it is approved, and a draft cannot add a check type, a foreign connection or a command. The builder tools are available only in WRM's built-in assistant.

#### Upgrade

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
