# Roadmap

Planned work that is agreed but not built yet. Released work is in [CHANGELOG.md](CHANGELOG.md).

## v10.10.0 — folder bookmarks (server-specific and global)

User request: the same deep directories (e.g. `/opt/app/servers/myServer/downloads/`) are needed on many servers; WinSCP-style bookmarks.

- **Bookmark** = a name and a path, optionally with a folder icon/colour and a note. Scope:
  - **server**: belongs to one connection;
  - **global**: offered on every SSH/SFTP/FTP connection;
  - optionally **folder / tag**: offered on all connections in a WRM folder or with a tag.
- Paths may use variables: `~`, `$USER`, `{host}`, `{name}` (connection name), e.g. `/opt/app/servers/{name}/downloads`.
- **File manager:** a ★ button next to the path bar: *Bookmark this directory* (choose scope) and a list of bookmarks; picking one opens that directory. A missing directory shows a clear message, not an error page.
- **Terminal:** the same list in the window menu and the command palette; picking one types `cd -- '<path>'` (shell-quoted) into the terminal; an option sends it to every terminal of a broadcast group.
- **Start directory:** a bookmark can be marked as the start directory of a connection (file manager opens there; terminal runs `cd` on connect, like run on connect).
- Management: rename, reorder, move between scopes, sharing of global bookmarks with share members (read-only), export/import with the configuration.
- Import of WinSCP bookmarks (`WinSCP.ini` `[Configuration\Bookmarks]`) if simple.
- Tests: API (scopes, ownership, variables), file manager and terminal in the browser.

## v11 — Git: compare and deploy services (decisions taken, not started)

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
- **Offline bundles (required):** for setups where the WRM server cannot reach the Git server.
  - **Import:** upload a `.tar.gz` bundle (a JSON manifest with files, hashes and per-file history, plus the payload) built elsewhere. WRM shows its age and source. Protected files in the payload are ignored, matched as globs.
  - **Export:** WRM builds the same bundle format, for servers that only a file-based tool can reach.

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

### Update and upgrade
- **Update:** bring an installation to the newest version of the ref it follows (same branch, a newer tag or commit).
- **Upgrade / switch:** move it to another branch or major version (e.g. a feature branch → the default branch). This gets an extra warning with a summary of the changes.
- **Rollback:** go back to an earlier backup or version.
- Any of these can be scheduled into a maintenance window. Nothing runs without an explicit choice.

For each selected installation and its selected files:
1. A double confirmation (typing the host name for production).
2. A check for dangling imports (Python).
3. A backup to `.deploy-bak/<time>/`.
4. Atomic writes that keep owner, mode and line endings.
5. Language checks (`py_compile`, `node --check`, `sh -n`, or a custom command).
6. An **automatic rollback** on failure.
7. Writing `VERSION.md` (overall version and a per-file table).
8. Optionally, a service restart (below).

Also:
- Rolling over several servers, stopping at the first failure.
- History with rollback.
- *Stamp* writes `VERSION.md` where an installation is already current.

### Service restart (opt-in only)
- WRM lists the systemd / supervisor units that belong to the installation. Restarting them is a choice per run, **off by default**, with a warning that names the units and servers.
- Choices:
  - do not restart (only show what needs a restart);
  - restart right after a successful update;
  - restart at a chosen time;
  - restart after a delay.
- Scheduled restarts:
  - are stored (they survive a WRM restart) and can be cancelled;
  - send a notification before and after;
  - are skipped if the update failed, or if WRM was down past a grace period;
  - are followed by a health check (`systemctl is-active`, `supervisorctl status`).

### Compatibility with file-based deploy tools
WRM writes `VERSION.md` (overall version and a per-file table), `.deploy-bak/<time>/` backups and an `updates.jsonl` log in a documented format, so a file-based tool on the server can keep working alongside WRM.

### New server
- **Install** from the repository into a chosen directory. Protected files are filled from templates in the repository (`*.example`) or entered in a form; secrets come from the vault.
- **Transfer**: copy an existing installation from another server, with its per-host configuration, but without logs and `.deploy-bak`.
- Optional systemd / supervisor unit from a template.

### .gitignore helper
- Pick files that exist only on servers (*extra*) and add standard patterns per language.
- Warn about per-host configuration committed to the repository.
- Edit the catalog's exclude / protected lists in the same place.
- Output: **a download by default**. A merge request only when the user explicitly chooses it; WRM never writes to a Git server on its own.

### Notifications (a general WRM module — built in v10.9.0)
**Already built in v10.9.0** (`notify.go`: channels, subscriptions, digests, quiet hours; first users live status and credential rotation). v11 only adds the Git events below.
- **Channels:** in WRM, browser notifications, e-mail (SMTP), Telegram, Slack / Mattermost / Rocket.Chat, Microsoft Teams (Workflows webhook), Discord, ntfy, Gotify, Pushover, and a generic webhook (JSON, HMAC-signed).
- **Configuration:** administrators set up the channels; users choose what they want to receive.
- **Behaviour:** one digest per check run instead of a flood, and optional quiet hours.
- **Git events:**
  - a new version is available;
  - drift on a server;
  - a server is unreachable;
  - update / upgrade started, succeeded, failed or rolled back;
  - a restart is scheduled or done.

### Permissions
Each action is a policy that administrators can change in the admin panel:
- check: all users by default;
- update, upgrade, install, transfer, rollback, restart and the .gitignore merge request: administrators by default.

### Security
- Read-only tokens (encrypted; write access only for merge requests).
- Self-signed Git servers are pinned on first use.
- Every update, upgrade, rollback, install, transfer and restart is audited.

### Optional CI integration (later)
WRM does not depend on Jenkins or GitLab CI, but can work with them:
- an incoming webhook (a GitLab / GitHub push or tag, or a Jenkins job) triggers an immediate check instead of waiting for the interval;
- a bundle can be taken from a CI artifact;
- a service whose deployment is a CI pipeline can be updated by triggering that job (Jenkins or GitLab CI) with the server and version as parameters, with its status shown in WRM.

### Phases
1. Sources, offline bundles, catalog, discovery, comparison, matrix, monitoring, notifications (read-only on servers).
2. Update / upgrade / rollback, opt-in restarts, scheduling.
3. New server (install / transfer).
4. .gitignore helper.
5. Optional CI integration.
