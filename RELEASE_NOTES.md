## Web Remote Manager PRO v14.11.0 — An interface for every screen

### Added

- **Account menu** on your avatar: your name and role, *Settings*, *Appearance*, *Shares*, *Lock now* and *Sign out* in one place.
- **Minimise dialogs**: the ▁ button above a big dialog turns it into a chip at the bottom of the screen; click the chip to continue where you left off.
- **Grouped notifications**: repeated messages show a ×N counter and older ones fold into one *more notifications* toast.
- **Tablets get the work area back**: between 761 and 1024 px (a tablet, a narrow laptop or a zoomed-in browser) *Sessions*, *Clipboard*, *Snippets* and the docked AI view slide in from the 🗂 button.

### Changed

- **The top bar keeps its labels.** Instead of turning every button into a bare icon as soon as one label did not fit, WRM now moves the less used buttons into the labelled **⋯** menu first. A 1920 px screen shows every button with its name again, and a laptop keeps the names of the ones you use most.

### Fixed

- Top bar buttons no longer overlap on phones.

### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.11.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.11.0) asset. Preserve the database and encryption key; nothing is migrated.

[Full changelog](CHANGELOG.md)

## Included from v14.7.2 — Shares a team can manage together

#### Added

- **Team management of shares** ([#192](https://github.com/vedranius/web-browser-RDM-public/issues/192)): the share dialog has two new switches, both off by default. *Operators may add and edit connections and folders* lets signed-in members with the operator role or higher edit shared connections and add connections and folders inside the shared folders (right-click a shared connection or folder). *Operators may start and stop tunnels* adds **⇄ Tunnels…** to the share, with *Start* / *Stop* for the tunnels you configured. Guests never get these rights.

#### Security

- Everything members add stays in your account. Members never see or pick your saved secrets, and changing a connection's host or protocol removes its saved password, so it is never sent to a server a member chose. Every change is audited under the member's name.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.7.2](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.7.2) asset. Preserve the database and encryption key; two columns are added to `share_links` and an older binary still starts on the database.

## Included from v14.7.1 — Security fixes, tunnels behind HSTS

#### Added

- **Setup code for a new installation**: until the first account exists, creating the administrator from another computer needs a one-time code that WRM prints to its console/log and writes to `first-admin-setup-code.txt` next to its database. On the WRM host itself (localhost) no code is needed. Before, whoever reached a fresh server first became its administrator.
- **Tunnels behind HSTS** ([#190](https://github.com/vedranius/web-browser-RDM-public/issues/190)): a tunnel on 0.0.0.0 can be opened through an IP address of the WRM host, and *Admin → Policies* can turn HSTS off, so a self-signed web interface behind a tunnel opens in the browser.
- **Security overview** warns about tunnels open to the network; a SOCKS tunnel on 0.0.0.0 is marked as a failure.

#### Fixed

- Two simultaneous first registrations can no longer both become administrators.
- The origin check compares the port, so another service on the WRM host is not trusted as WRM.
- One AI client closing a shared AI session no longer stops the others; `~` in AI paths is expanded reliably under parallel calls.
- Built with Go 1.27.4 (net/http fix GO-2026-6859).

## Included from v14.7.0 — Addons that show your chains

#### Added

- **`chains.read` for addons**: an addon window can now read your [chains](README.md#chains-live-service-chain-monitoring): the checked values, the active scenario and what changed when. WRM keeps checking the chains in the background, so an addon can present a whole system the way your team needs it (a table per environment, a report, an update of a wiki page) without SSH, commands or secrets of its own.
- **Check now** from the addon window, with the same limit as in the Chains window, recorded in the audit log.
- The permission sheet says exactly what the addon may read. Every viewer sees only their own chains, also in a shared instance.
- For addon authors: `wrm.chains()`, `wrm.chain(id)` and `wrm.chainCheck(id)` in `wrm-bridge.js`, and sample chains in `wrm -addon-dev` through `WRM_ADDON_CHAINS`.

## Included from v14.6.3 — Proxies from SSH tunnels, a steady lock screen, log discovery that keeps going

#### Fixed

- **SSH tunnels as proxies** ([#188](https://github.com/vedranius/web-browser-RDM-public/issues/188)): a dynamic tunnel (`-D`) of your connection now shows up in the connection's and folder's **Proxy** list and under *Proxies*. Pick it and WRM makes it a proxy, also when its port is *auto*.
- **Lock screen keeps the keyboard** ([#137](https://github.com/vedranius/web-browser-RDM-public/issues/137)): a terminal reconnecting behind the lock, a click on the page or returning to the tab no longer takes the cursor away from the password field.
- **Log discovery keeps going**: Python virtual environments (also renamed ones such as `venv.bak.2026-07-16`) are skipped, an unreadable or silent folder is skipped with a note, and an open log no longer fails with *Find logs first* after an hour: WRM finds the logs again by itself.

## Included from v14.6.2 — Forgot password and activation by e-mail

#### Forgot password

- **Forgot password?** on the sign-in page: enter your user name or confirmed e-mail address and WRM mails a link that works once, for one hour. The page says the same whether or not an account matched, so it does not reveal who has an account.
- On the link you choose a new password. With two-factor authentication you also enter a code (or a recovery code), so someone who only got into your mailbox cannot change it. Your other sign-ins end and you get a notice by mail.
- Directory and single sign-on accounts never get a link, and nobody does during a lockdown. At the **Enterprise** level administrators cannot use it at all (policy `password_reset_mail`); they reset each other's passwords instead.

#### New accounts

- **Invitations**: when creating a user, enter an e-mail address and tick *Send an invitation*. The user chooses the password through the link (valid three days), which also confirms the address. No more passing temporary passwords around.
- **Open registration with confirmation**: a new account asks for an e-mail address and starts once the link sent to it is opened. Accounts never confirmed are removed after a week.

## Included from v14.6.1 — WRM sends mail

#### Set it up once

- **Admin → E-mail**: pick **Gmail**, **Outlook.com**, **Microsoft 365** or another SMTP server, enter the user name and an app password, the sender and WRM's address, and press *Send a test message*. WRM says what each provider needs (an app password; for Microsoft 365, SMTP AUTH allowed for the mailbox). It works the same for WRM on your own computer and on a server, as long as that machine can reach the mail server.
- The password is stored encrypted and never shown again. Mail waits in a queue and is retried, so a slow mail server never slows WRM down.

#### For every user

- **E-mail address with confirmation**: *Settings → Security → E-mail address* sends a link; the address counts once you open it. Administrators can set addresses in *Admin → Users*, still confirmed by the user.
- **Security notices**: you get a mail when your password is changed or reset (by you, an administrator or on the server), when two-factor authentication is turned on or off, and when a new address is set. Each message is in English and Croatian.

## Included from v14.6.0 — Password reset on the server, protected

Until now anyone who could run WRM on its server could reset any password, the administrator's included, with `wrm -reset-password`, and nobody was told. That is still possible for a person who really needs it, but it can no longer happen quietly.

#### How it is protected

- **By security level** (Admin → Policies → Accounts → *Password reset on the server*):
  - **Alert administrators** (Home): the reset happens at once and every administrator gets a security alert.
  - **Only while WRM is stopped** (Business): the command refuses while WRM is running, so a takeover means a visible outage.
  - **Confirmation** (Enterprise): the command files a request. Another administrator approves it in Admin → Security, or the user confirms it with their own 2FA code (`-reset-code 123456`), or it can be run again after one hour if nobody rejected it.
- **Every reset is reported.** Resets, refused attempts, requests and wrong codes go to every administrator's inbox and notification channels, the audit log and SIEM.
- **AI can't do it.** An AI or MCP session can never run `wrm -reset-password`, `-reset-2fa` or `-lockdown` on a server, even with full access, and a temporary password in its output is masked.

## Included from v14.5.1 — Resizable windows and storage columns

#### Windows

- **Full screen and any size.** The big windows opened from the menus (storages, chains, admin, settings, monitoring and other large dialogs) have a small bar above their top-right corner: ⤢ for full screen, ⤡ back, ↺ for the default size. Drag the bottom-right corner to resize. Each window keeps its size in your browser.

#### Storage browser

- **Columns like a file manager.** Click *Name*, *Size* or *Date* to sort, drag a column border to widen it, and double-click the *Name* border to fit the longest file name. Long names scroll sideways instead of being cut. The storage list and the S3 activity panel can be resized by dragging their edges.
- The search field and the sort list follow the dark theme, and the sort direction arrow shows the real direction.

## Included from v14.5.0 — Chains that recognise your servers

#### A server list that makes sense

- **Servers by folder, with search.** The Chains builder, the AI wizard and the templates list your saved SSH servers by folder with their host and tags. Search by name, host, folder or tag, tick several and add them at once. Quick connections are not listed, a server already in the chain cannot be added twice, and two servers with the same name in different folders are marked.

#### Recognition and suggestions

- **What runs on a server becomes a check with one click.** Right after adding servers, WRM looks at them read-only (services, Docker containers, listening ports and network mounts) and suggests checks. Ports carry their service name, e.g. PostgreSQL (5432). *Add the recommended* takes the usual set.
- **Suggested links.** A network mount from another server of the chain suggests the link between the two. A mount from a host outside the chain suggests adding that host as an outside endpoint.
- **Smart outside endpoints.** Type a URL, `host:port` or a host name, and WRM adds the fitting checks (HTTP and certificate, TCP port, or DNS). Saved HTTP/HTTPS connections can be added from *Web connections*.
- **Test now** in the check dialog runs only that check once, without saving.

#### Fixes

- **WRM no longer freezes while an AI works through MCP.** Every MCP session's live view held its own connection, and over plain HTTP a browser has only six per server for all tabs; at six sessions the page stopped responding. All AI views of a tab now share one connection.
- The builder's fields follow the dark theme, the Croatian confirm count is grammatical, the automatic states are named in your language, and the fields explain themselves.
- The lock button follows a language change. Monitoring cards show the server's name as the title, and the doubled gear icon on *Settings* is gone.

## Included from v14.4.0 — S3 activity and costs, AI in the Chains builder

#### S3 requests you can see

- **Live S3 activity panel** beside the storage browser. It lists every request WRM sends to the bucket, as it happens: the operation, key, answer and duration, and what caused it (you, the background check, a chain, an AI connection). Each operation is counted for today and in total. The panel is shown by default; ✕ or 📡 hides it.
- **Provider picker with request costs** in the storage form: Hetzner, DigitalOcean Spaces, Amazon S3, Backblaze B2, Cloudflare R2, MinIO and Other. Picking one shows whether, in which class and for how much the provider bills each S3 operation WRM uses, quoting and linking the provider's own price page:
  - Hetzner and Backblaze B2 bill none of them;
  - DigitalOcean's price list has no request fee;
  - Amazon S3 bills every request except DELETE;
  - Cloudflare R2 bills two classes after a monthly free allowance.

  The full table is in [docs/STORAGES.md](https://github.com/vedranius/web-browser-RDM-public/blob/v14.4.0/docs/STORAGES.md#what-each-provider-bills).

#### AI in the Chains builder

- **🤖 AI on a chain or draft.** Say what to add, change or check. The AI reads the chain and its last results and writes a draft. A change to a chain waits as *Changes for …* until you press **Apply changes**; the chain keeps its id and history.
- **Test once** runs a draft's checks a single time and shows the results on the diagram, without saving.
- **New tools.** WRM's AI panel gets `chain_get` and `chain_test`, and `chain_draft` can change an existing chain. MCP clients get the same four tools behind the new opt-in **`chains`** scope: an administrator adds it to the allowed scopes, then a user enables it on an AI connection.

#### Upgrade

Use WRM's update action or replace the executable/container with the matching [v14.4.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v14.4.0) asset. Preserve the database and encryption key. Two columns are added (`storages.provider`, `chain_drafts.base_chain_id`); an older binary still starts on the database.

[Full changelog](CHANGELOG.md)

## Included from v14.3.2 — Edit and remove storages, fewer S3 requests

A bug-fix release for the storage section: a saved storage can be edited and removed again, and WRM sends fewer background requests to the bucket.

#### Fixed

- **Edit and remove a storage.** Each storage in the list has a ✎ (edit) and a 🗑 (remove) button, and the open storage has *Edit* in the toolbar. The storage form also has *Remove storage*. Before this release the form could only be opened for a new storage. After saving, the open storage uses its new permissions straight away, and on a phone the list sits above the files.

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
