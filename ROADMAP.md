# Roadmap

Planned work that is agreed but not built yet. Released work is in [CHANGELOG.md](CHANGELOG.md).

## v11 — Git: deploy services (phases 3–5; decisions taken)

Phase 1 (sources, offline bundles, catalog, discovery, file-level comparison, overview matrix, periodic checks with notifications; read-only on servers) shipped in **v11.0.0**, phase 2 (update, upgrade, rollback, stamp, opt-in restarts, scheduling, the deploy tool's `VERSION.md` / `.deploy-bak` / `updates.jsonl` formats, deploy policies and notifications) in **v11.1.0**, see [CHANGELOG.md](CHANGELOG.md). The phases below build on that **Git** workspace (`git*.go`, `static/git.js`) and keep its rules: generic for any GitLab / GitHub, servers without git, pip or internet access, everything over WRM's SSH connections.

### New server
- **Install** from the repository into a chosen directory. Protected files are filled from templates in the repository (`*.example`) or entered in a form; secrets come from the vault.
- **Transfer**: copy an existing installation from another server, with its per-host configuration, but without logs and `.deploy-bak`.
- Optional systemd / supervisor unit from a template.

### .gitignore helper
- Pick files that exist only on servers (*extra*) and add standard patterns per language.
- Warn about per-host configuration committed to the repository.
- Edit the catalog's exclude / protected lists in the same place.
- Output: **a download by default**. A merge request only when the user explicitly chooses it; WRM never writes to a Git server on its own.

### Permissions
Each action is a policy that administrators can change in the admin panel (`git_enabled` and `git_checks` exist since v11.0.0; `git_update`, `git_upgrade`, `git_rollback` and `git_restart` since v11.1.0):
- install, transfer and the .gitignore merge request: administrators by default.
- Notifications for install and transfer (started, succeeded, failed), like the run events of v11.1.0.

### Security
- Tokens stay read-only (encrypted since v11.0.0); write access only for merge requests.
- Every install and transfer is audited (updates, upgrades, rollbacks and restarts are since v11.1.0).

### Optional CI integration (later)
WRM does not depend on Jenkins or GitLab CI, but can work with them:
- an incoming webhook (a GitLab / GitHub push or tag, or a Jenkins job) triggers an immediate check instead of waiting for the interval;
- a bundle can be taken from a CI artifact;
- a service whose deployment is a CI pipeline can be updated by triggering that job (Jenkins or GitLab CI) with the server and version as parameters, with its status shown in WRM.

### Phases
1. ~~Sources, offline bundles, catalog, discovery, comparison, matrix, monitoring, notifications (read-only on servers).~~ Done in v11.0.0.
2. ~~Update / upgrade / rollback, opt-in restarts, scheduling.~~ Done in v11.1.0.
3. New server (install / transfer).
4. .gitignore helper.
5. Optional CI integration.
