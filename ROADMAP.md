# Roadmap

Planned work that is agreed but not built yet. Released work is in [CHANGELOG.md](CHANGELOG.md).

## v10.8.1 — fixes

- **mRemoteNG and inventory import of larger files.** The security middleware limits every API body to 1 MB (`bodyLimitFor` in `security.go`); only `/api/config/import` has 32 MB. `/api/config/import/mremoteng`, `/api/config/import/sshconfig` and `/api/inventory/*` therefore fail with "Bad JSON" from about 1 MB on.
  - Give `/api/config/import/*` and `/api/inventory/*` a proper limit (32 MB).
  - A clear error when a file is too large (its size and the limit), and a size check in the browser before sending.
  - Tests with a multi-MB `confCons.xml` and a large CSV.
- **Install as an app:** explain in the UI what to do in the current browser.
  - Exact steps per browser (Chrome/Edge, Safari, iOS, Android).
  - The real reason when installing is not possible (plain `http://` to an IP address, Firefox on the desktop).
  - Hide the *Install* button when it cannot work.
  - README: a short HTTPS recipe (Caddy).

## v10.9.0 — proxy as a connection setting

Requested by users who want to replace PuTTY: some servers are reached through jump hosts, others through a SOCKS or HTTP proxy.

### Saved proxies
A proxy is defined once and picked per connection, like a vault credential. It has:
- a name;
- a type: **SOCKS5**, **SOCKS4/4a** or **HTTP CONNECT**;
- host and port;
- an optional user name and password (encrypted, or a vault credential);
- *DNS at the proxy* (PuTTY: "Do DNS name lookup at proxy end").

Proxies can be shared with other users without revealing the password. Export and import keep them by name.

### In the connection dialog
A **Proxy** field next to *Connect via (jump host)*: none / a saved proxy / *+ New proxy…*.

### What goes through the proxy
- Everything a connection does over TCP:
  - SSH terminal, SFTP, search, transfers, tunnels;
  - key deploy and password rotation;
  - FTP/FTPS;
  - RDP/VNC/Telnet (through a local relay, like jump hosts today);
  - web interfaces;
  - *Test connection*;
  - live status;
  - network tools from a server;
  - Redfish BMCs.
- Not IPMI and Serial-over-LAN: they use UDP. The UI says so.
- **Proxy and jump host together:** the proxy is reached through the jump hosts, then the target. The route is shown and audited, e.g. `bastion → socks5://10.1.1.1:1080 → app-01`.
- Proxy type **"WRM SOCKS tunnel of connection X"**: a dynamic tunnel that WRM starts when needed.

### Docker
Inside a container `127.0.0.1` is the container itself.
- `docker-compose.yml` gets `extra_hosts: host.docker.internal:host-gateway`.
- The README explains the setup.
- WRM warns when a proxy at `127.0.0.1` is set while it runs in a container.

### Folder defaults
A folder can carry a default jump host and proxy. Its connections, including new ones, inherit them unless they set their own.

### Security and errors
- Proxy passwords are encrypted and never sent to the browser.
- Policy `proxies`: who may define proxies.
- Use is audited with the route.
- Clear messages that separate proxy errors from target errors: proxy unreachable, authentication required or failed, SOCKS reply codes, HTTP status.

### Tests
In-process SOCKS5 and HTTP CONNECT proxies, with and without authentication, covering:
- terminal, SFTP, the desktop relay, web interfaces and status;
- jump host + proxy together.

### Later (v10.9.x)
- **Import PuTTY sessions** from a `.reg` export (`HKCU\Software\SimonTatham\PuTTY\Sessions`): host, port, user, key and proxy settings. Proxies are created and linked to the connections.
- With mRemoteNG, connections that point to a PuTTY session (`PuttySession`) get that session's proxy.

## Idea, awaiting decisions — Git: compare and deploy services

A separate **Git** workspace in WRM: a button in the top bar opens it full screen, and its code loads only then. It stays out of the connection tree and the terminal UI, and policy `git_enabled=0` hides it completely.

It is generic: it works with any GitLab or GitHub (cloud or self-hosted) that a user connects. WRM ships no services of its own and no organisation-specific configuration.

The design targets servers **without git, pip or internet access**: WRM fetches code through the Git provider's API and works on servers over its existing SSH connections, so jump hosts, the vault and proxies all apply.

### Service catalog
Per service:
- project;
- ref: `tag:latest`, `tag:vX` or a branch;
- tag filter (regex);
- subdirectory of the repository that matches the installation root;
- include / exclude globs;
- **protected** per-host files that are never overwritten (`config*`, `*.ini`, `.env`);
- **fingerprint** files used to recognise an installation;
- kind (app, library, tool).

How it is filled in:
- WRM suggests a catalog entry from the repository tree.
- Catalogs can be imported and exported as JSON.

Refs:
- Branch-aware tags: only tags reachable from the chosen branch count, with a warning and a fallback to the branch head.
- Ref choice per run: the catalog, one branch for all, per service, or each project's default branch.

### Sources
- The GitLab API v4 or the GitHub REST API: groups or organisations (with subgroups), projects, branches, tags.
- **Fewer API calls.** One recursive tree listing per ref gives the blob ID of every file. File contents are fetched only for blobs not seen before and cached by blob ID forever (the same blob in 20 versions is fetched once). The normalised hash is computed from that content.
- **Offline source:** import a bundle (a JSON manifest with files, hashes and per-file history, plus the payload) built elsewhere, for setups where the WRM server cannot reach the Git server.

### Discovery
- One SSH call per server walks the configured roots (e.g. `/opt`, `/srv`, a scripts directory) to a limited depth.
- It finds installations by their fingerprint files.
- It skips copies and backups (`*_BKP`, `*_OLD`, `backup`, `.deploy-bak`, …; configurable).
- It derives an environment label from the path.

### File-level comparison
- The server returns normalised content hashes of its files (POSIX `sha256sum`, CRLF→LF, no Python needed).
- WRM compares them with the target ref and with the **history of each file** on that branch (cached by blob).
- **File states:**
  - ok;
  - old (a known earlier version, *n versions behind*);
  - modified (local change);
  - missing;
  - extra (only on the server);
  - protected.
- **Installation states:**
  - *needs update*;
  - *review* (local changes);
  - *up to date*.
- Details show:
  - a coloured diff;
  - the systemd / supervisor units that point to the installation directory;
  - the version recorded in `VERSION.md`.

### Overview and monitoring
- A servers × installations matrix with filters (needs update / review / up to date, environment, folder, tag).
- Periodic checks while WRM runs. They **never update on their own**; they notify about new versions and about drift (files changed by hand on a server).

### Update
For each selected installation and its selected files:
1. A double confirmation (typing the host name for production).
2. A check for dangling imports (Python).
3. A backup to `.deploy-bak/<time>/`.
4. Atomic writes that keep owner, mode and line endings.
5. Language checks (`py_compile`, `node --check`, `sh -n`, or a custom command).
6. An **automatic rollback** on failure.
7. Writing `VERSION.md` (overall version and a per-file table).
8. An optional service restart and health check.

Also:
- Rolling over several servers, stopping at the first failure.
- History with rollback.
- *Stamp* writes `VERSION.md` where an installation is already current.

### New server
- **Install** from the repository into a chosen directory. Protected files are filled from templates in the repository (`*.example`) or entered in a form; secrets come from the vault.
- **Transfer**: copy an existing installation from another server.
- Optional systemd / supervisor unit from a template.

### .gitignore helper
- Pick files that exist only on servers (*extra*) and add standard patterns per language.
- Warn about per-host configuration committed to the repository.
- Edit the catalog's exclude / protected lists in the same place.
- Output: a merge request, or a download.

### Security
- Read-only tokens (encrypted; write access only for merge requests).
- Self-signed Git servers are pinned on first use.
- Policies `git_enabled`, `git_checks`, `git_deploy`.
- Every update, rollback and install is audited.

### Phases
1. Sources, catalog, discovery, comparison, matrix, monitoring (read-only).
2. Update and rollback.
3. New server (install / transfer).
4. .gitignore helper.
