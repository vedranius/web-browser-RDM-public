# Roadmap

Planned work that is agreed but not built yet. Released work is in [CHANGELOG.md](CHANGELOG.md).

## v11 — Git: deploy services (phases 2–5; decisions taken)

Phase 1 (sources, offline bundles, catalog, discovery, file-level comparison, overview matrix, periodic checks with notifications; read-only on servers) shipped in **v11.0.0**, see [CHANGELOG.md](CHANGELOG.md). The phases below build on that **Git** workspace (`git*.go`, `static/git.js`) and keep its rules: generic for any GitLab / GitHub, servers without git, pip or internet access, everything over WRM's SSH connections.

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

### Notifications
The notifications module (v10.9.0) already carries the phase 1 Git events (new version, drift, unreachable server). Still to add:
- update / upgrade started, succeeded, failed or rolled back;
- a restart is scheduled or done.

### Permissions
Each action is a policy that administrators can change in the admin panel (`git_enabled` and `git_checks` exist since v11.0.0):
- update, upgrade, install, transfer, rollback, restart and the .gitignore merge request: administrators by default.

### Security
- Tokens stay read-only (encrypted since v11.0.0); write access only for merge requests.
- Every update, upgrade, rollback, install, transfer and restart is audited.

### Optional CI integration (later)
WRM does not depend on Jenkins or GitLab CI, but can work with them:
- an incoming webhook (a GitLab / GitHub push or tag, or a Jenkins job) triggers an immediate check instead of waiting for the interval;
- a bundle can be taken from a CI artifact;
- a service whose deployment is a CI pipeline can be updated by triggering that job (Jenkins or GitLab CI) with the server and version as parameters, with its status shown in WRM.

### Phases
1. ~~Sources, offline bundles, catalog, discovery, comparison, matrix, monitoring, notifications (read-only on servers).~~ Done in v11.0.0.
2. Update / upgrade / rollback, opt-in restarts, scheduling.
3. New server (install / transfer).
4. .gitignore helper.
5. Optional CI integration.
