# Changelog

All notable changes to Web Remote Manager PRO. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Release notes with downloads are on the
[Releases page](https://github.com/vedranius/web-browser-RDM-public/releases).

## [11.6.0] — 2026-10-09 — Git: activity and deploy history per user and branch

### Added
- **Activity tab** in the Git workspace with three views (`static/git.js`): *Git activity*, *Deploy history* and *What is where*; every view has filters and a **CSV export** of what the filters show. English and Croatian; works at phone width (tables scroll inside their cards).
- **Git activity per service** (`git_activity.go`, `GET /api/git/activity`):
  - GitLab: project events with `action=pushed` (branch and tag pushes with the commit count and the range) and the commits of the listed branches; GitHub: the commits of the listed branches and, where the token allows it, the repository activity (`/activity`: pushes, force pushes, branch creation and deletion, merges); Gitea: the commits of the listed branches.
  - Branches: the chosen branch, or the service's, the target's, the default branch, `fallback_branches` and the branches of the environments' `allowed_refs` (at most 6). A commit on several listed branches is one item with its branches; a GitHub push by a login is attributed to the name seen on that login's commits.
  - Filters: branch, author (name, user name or e-mail, case-insensitive), dates (`since` / `until`, default the last 30 days, at most 366). Summaries **per author** (commits, pushes, branches, last activity) and **per branch** (commits, pushes, authors, last activity, link to the branch's commits). Links to the commit page of every commit and the compare page of every push range (GitLab `/-/commit/`, `/-/compare/`; GitHub and Gitea `/commit/`, `/compare/`).
  - **Paging and rate limits:** listings follow the provider's next links (`Link: rel="next"`, GitLab `X-Next-Page`, GitHub's activity cursors) up to 10 pages per listing, stop early when `RateLimit-Remaining` / `X-RateLimit-Remaining` reports 5 or fewer calls left, and stop at the first commit older than the range (servers that ignore `since`). HTTP 429 and 403 with no calls left become a note with the reset time (`Retry-After`, `RateLimit-Reset`, `X-RateLimit-Reset`); the listing is marked `truncated` / `rate_limited`.
  - Cached per user, source, service, branch and dates for 3 minutes (`refresh=1` reads again); the author filter and the summaries are computed from the cached listing.
  - `GET /api/git/activity/export.csv`.
- **Deploy history** (`git_where.go`): `GET /api/git/history/cell` reads `history.jsonl` of every destination of an environment (one request per service and environment, loaded lazily three at a time, policy `git_checks`), `GET /api/git/history/runs` lists the WRM runs (update, upgrade, rollback, install, transfer) of installations without environments (environment runs are left out: their destinations keep the history). Rows carry the action, result, who, branch, version, commit (linked), the version before, counts and the per-file lists. Filters: service, environment, server, who, branch, version or commit, action (`deploy`, `rollback`, `update`, `upgrade`, `install`, `transfer`), dates. Cells are cached for 2 minutes. `GET /api/git/history/export.csv` merges everything (the server parts only with `git_checks`).
- **What is where** (`git_where.go`): `GET /api/git/where/cell` per environment — per destination the branch (`state.json`, else `VERSION.md`'s *Grana/ref*; a tag falls back to the service's branch), commit, version, who and when, the deploy id, the lock, **commits behind the branch head** (`branchHead` + the compare API, one call per branch head and per distinct commit, cached for 2 minutes) with links to the commit, the branch and the compare page, and **drift**: destinations whose commit and version differ from the most common state of the environment are marked. `GET /api/git/where/installs` lists installations without environment from their last check (no SSH). `GET /api/git/where/export.csv`.
- **Gitea as a Git source** (`kind: "gitea"`, API base `<url>/api/v1`, `Authorization: token …`, paging with `limit`): projects, branches, tags, commits, the paginated recursive tree, blobs and compare work as with GitHub, so targets, comparisons, deploys and the activity work too. A .gitignore pull request is not available for Gitea (the helper says so; the download works).
- `GET /api/git/env/cell` also returns `ref_md`, `user_md` and `updated_md` from `VERSION.md`.
- Tests: `git_activity_test.go` (fake GitLab, GitHub and Gitea: paging over `X-Next-Page`, `Link` headers and cursors, tracked branches, deduplicated commits, branch / author / date filters, servers that ignore `since`, author and branch summaries, links, cache and refresh, few calls left, HTTP 429, CSV header, byte order mark, escaping and formula neutralisation) and `git_where_test.go` (history merged from two destinations on two connections with WRM runs, filters, cache and refresh, CSV with file lists, the `git_checks` policy; commits behind the branch head, drift, installations without environment, CSV); `TestGitTargets` also runs against a fake Gitea.

## [11.5.0] — 2026-10-09 — Git: environments, deploy plans and safe deploys

### Added
- **Environments per service** (catalog `apps[].environments`, `git_env.go`; optional — catalogs and installations without environments work as before):
  - Ordered **destinations** `{server, path}` (`server` = the name of an SSH connection; paths checked like install directories) and per-environment rules: **`confirm`** (type the environment name; always for `prod`, `production`, `prd`, `live`), **`allowed_refs`** (`refs/heads/main`, `refs/tags/v*`; a pattern without `refs/` matches the short name), extra **`ignore`** patterns, an optional **`post_deploy`** command (runs in the app directory only when chosen for a run), **`keep_backups`** (1–50, default 5) and an optional **`state_dir`**.
  - Validation: two destinations on one server where one app path lies inside the other are rejected unless the outer service's excludes / environment ignores cover the inner path; a state directory inside any app path is rejected; one path cannot belong to two services.
- **Deploy plan** (`POST /api/git/env/plan`): per destination one scan over SSH (hash with CR removed, number of CRs, symlinks, unreadable files, `VERSION.md`, `state.json`, the lock) compared with the target or another ref: `new`, `changed` (`local` when it matches no version of the repository, else how many versions behind), **`eol`** (*same (CRLF/LF)*: only the line endings differ — counted separately and sent, which normalises them), `same`, `removed` (in the repository at the deployed commit, not any more; deleted only on request), `extra`, `protected`, `ignored` (with the pattern). Type changes (a directory that becomes a file or back) are listed with the server files that must be deleted; a symlink that would be replaced, files changed by hand, a missing path (the first deploy creates it and its parents) and a lock are warnings. Each destination plan has a **fingerprint**; plans live in memory for 2 hours and are **applied once**.
- **Environment deploy runs** (`git_env_run.go`; `POST /api/git/runs` with `env`): kind `update` (or `upgrade` for another ref), the destinations in order, stopping at the first failure. Per destination: a **server-side mkdir lock** with an owner file (deploy id, WRM user, WRM host, server time), the **plan computed again under the lock** (another fingerprint stops the deploy before anything is written), the selection (minus per-run exclusions; deletes only when chosen), one **tar stream** into `<path>/.wrm-incoming-<id>` verified by hash, a **backup** of replaced and deleted files to `<state>/<deploy id>/`, deletes and renames (mode and owner of replaced files kept; directories emptied by deletes removed and recorded with mode and owner), checks, automatic restore on failure, the **manifest** `deploys/<id>.json`, **`state.json`** (`current`, `previous`), the `updates.jsonl` line, `post_deploy` when chosen, a **DEPLOY entry in `history.jsonl`** (the previous state and the new / changed / deleted / excluded files, at most 500 per kind, with the real counts) and pruning of older backups. Only the deploy's **own lock** is released (its id must be in the owner file).
- **Dry run** for deploys and rollbacks: the exact options of the real run, no lock, nothing written, no typed confirmation; the run log lists what would happen. No notifications for dry runs.
- **Failure handling:** the transfer timeout is at least 1 hour; an **interrupted transfer** is recorded as partial (in the run and in `history.jsonl`) with the files that had arrived (removed again) and writes no new state; a **failed `post_deploy`** is recorded with a note (the files are on the server, the new state is current) and the remaining destinations are not started.
- **Rollback per deploy id** (`kind: rollback` with `env.deploy_id`, chosen destinations, default all, **in reverse destination order**): under the lock it checks that the deploy is still the latest there and that its backup and manifest exist; files the deploy added are **moved** to `<state>/backups/rollback-<id>/`, the current versions of changed files are kept there too, backed-up and deleted files come back atomically, removed directories are re-created with their mode and owner. **Nothing is done through a symlink**: such items are skipped and the rollback is incomplete (state unchanged; it can be repeated). Then `post_deploy` when chosen, the stored `previous` becomes current, and a **ROLLBACK** entry with `rolled_back=<id>` is written. The last 5 rollback backups are kept.
- **Stale locks:** older than max(30 minutes, the transfer timeout) by the server's clock; removed only after a confirmation in the UI (`POST /api/git/env/unlock`) by one server command that checks the owner id and the age again. Audit `git.lock_removed`.
- **Per-run exclusions** are remembered per environment with who and when (table `git_env_exclusions`) and pre-unticked on the next compare; `DELETE /api/git/env/exclusions` clears them. **Permanent exclude** (`POST /api/git/env/exclude`, scope `service` = the catalog's excludes, `env` = the environment's ignore list) writes an anchored pattern with glob characters escaped, unless an existing pattern already covers the file. Audit `git.exclude_added`.
- **Overview by environment** (`GET /api/git/env/overview`, `GET /api/git/env/cell` — one request per cell, loaded lazily): per destination version, commit, who / when, commits behind the target, the lock; a drift warning when the servers of one environment differ.
- **History** (`GET /api/git/env/history`): `history.jsonl` of every destination, one row per destination, newest first, with the current deploy of each destination.
- **Access doctor** (`POST /api/git/env/doctor`): SSH, a clean non-interactive shell, the SSH user (a warning for root), the app directory (OK when a missing path's nearest existing parent is writable), the state directory, unreadable subdirectories, signs of another deploy method (`.git`, `.svn`, `.hg`, `releases` / `current`, a symlinked path), a public web root, `sha256sum` / `rsync` / `tar` with versions, the lock.
- UI: a new **Environments** tab in the Git workspace (lazy cells, ✎ editor, ▶ Deploy with the plan, per-file checkboxes, *⊘ service* / *⊘ environment*, options, dry run, typed confirmation, *Back to the plan*; 🕘 History with expandable file lists and *Roll back…*; 🩺 Doctor). Run details show the deploy id, dry-run badge, partial transfers, skipped symlinks, `post_deploy` and notes; English and Croatian.
- Tests: `git_env_test.go` (plans with CRLF-only files, dry run, typed confirmation, one-shot plans, fingerprint mismatch, lock contention, stale lock, releasing only one's own lock, remembered and permanent exclusions with escaping, `allowed_refs`, nested-path and state-directory validation, type changes, partial transfer, `post_deploy` failure, rollback with moved-aside files, skipped symlinks, re-created directories and kept backups, doctor, drift) on the fake SSH host and the fake GitLab.

### Changed
- Exclude / ignore / protected patterns that start with `/` are **anchored**: they match the whole relative path only.
- The backup list of the per-installation *Rollback…* skips the `backups` and `deploys` directories of the environment state.
- `PUT /api/git/catalog/apps/{name}` keeps the environments of a service when the request does not send `environments`.

## [11.4.1] — 2026-10-09 — Clipboard history per user

### Fixed
- **The clipboard history leaked between accounts on the same browser** ([#30](https://github.com/vedranius/web-browser-RDM-public/issues/30)). The *Clipboard* panel stored its history in one browser-wide `localStorage` key, so after signing out and signing in as another user, the previous user's copied texts were still listed. The history is now stored per user (`wrm_clipboard_v9:u<user id>`), signing out removes the signed-in user's history from the browser, and the old browser-wide key is deleted on the next sign-in (it cannot be attributed to a user).

## [11.4.0] — 2026-10-09 — Git: every file of an installation, partial and incremental checks

### Added
- **Every file of an installation** (`git_files.go`, the details of an installation):
  - `GET /api/git/installs/{id}/files` lists every file of the server folder (one `find` over the folder with POSIX tools, `stat -c` / `stat -f` for size and time, hashes as in checks; `.git`, `.deploy-bak`, `node_modules`, `__pycache__`, virtual environments and the catalog's `ignore_dirs` pruned; backup-looking paths skipped; at most 10 000 entries, then *truncated*) and every file of the repository at the target ref below the service's subdirectory (the cached tree of the target commit; targets now keep the full commit, `commit_full`).
  - Per file: kind, size, modification time, on server / in Git / tracked, a **state** (`ok` / `old` / `modified` for tracked files; `same` / `differs` / `unknown` for repository files the catalog does not track; `extra`, `missing`, `protected`, `symlink`, `unreadable`) and a **reason** with the matching pattern (`tracked`, `excluded`, `not_included`, `protected`, `extra`, `tool`, `compiled`, `not_tracked`, `no_target`). Untracked repository files are compared through the blob cache, fetching at most 200 unknown blobs per listing.
  - Symlinks are listed with their target and never followed (an *outside* flag when they point out of the installation); unreadable files and directories are shown as such.
  - `GET /api/git/installs/{id}/file?path=&side=server|git` opens one file (view only, up to ~1.9 MB, binary detected by NUL bytes, otherwise the size and hash). The server side never follows a symlink and refuses paths that lead out of the installation (`pwd -P` check); paths are always shell-quoted. Audit `git.file_viewed`.
  - `GET /api/git/installs/{id}/diff` works for **any file in the repository at the target**, including files the catalog excludes or does not include (`informational: true`, with the reason and pattern), and files missing on the server; binary or large files are compared by hash (`same`). Protected files are still not compared.
  - `GET /api/git/installs/{id}/commits`: the commits between the commit in `VERSION.md` (or its version) and the target, newest first, with the count (GitLab `repository/compare`, GitHub `compare/a...b`).
  - UI: filters *All / Differing / Tracked only / Not tracked*, search, *Open* (server or Git side) and *Diff* per file, unified and side-by-side diff, the informational label, *Commits behind* with the list.
- **Check jobs** (`git_jobs.go`): `POST /api/git/check/jobs` with `install_ids`, `conn_ids`, `all`, `discover`, `refresh`, `stale_minutes` or `only_never`; `GET /api/git/check/jobs/current?since=N` returns the job with every installation's state (queued / running / done / error / cancelled) and the installations finished after `N`; `POST /api/git/check/jobs/current/cancel` skips the queued installations and aborts running scans without saving them. Each installation is scanned and saved on its own; concurrency is bounded per server and overall by the new settings **`git_check_per_server`** (2) and **`git_check_parallel`** (8; also used by periodic checks). One check at a time per user, as before. Audit `git.checked` (with the job) and `git.check_cancelled`.
- UI: overview filters by **service** and **server** (the search also matches host and environment); checkboxes per installation, server row and service column; *Check selected*, *Check visible*, *Check this cell*, *Check only stale* (never checked, older than 1 hour / 24 hours / 7 days) and *Cancel*; a spinner per queued or running cell; "checked X ago" and the last good result on every cell; the same selection and buttons on the *Installations* tab. *Check now* runs as a job too.
- Tests: `git_files_test.go` (the file list with reasons and patterns, an excluded file that is diffable, CRLF, a file only in Git, symlinks inside and outside, an unreadable file, binary files, quoting, commits, other users — against the fake GitLab and GitHub and the fake SSH host; check jobs: streaming, a partial check that keeps other results, stale and never-checked selections, the per-server and overall limits, cancel, an unreachable server that keeps the last good result, policies; the restore of set-aside Git sources).

### Changed
- A failed check of an installation keeps the files and counts of the last good result; the new columns `git_installs.last_ok_at` / `last_ok_state` remember it.
- Installations carry the connection's host (`host`).

### Fixed
- Upgrading a database from v11.2.0 or earlier to v11.3.0 set the `git_sources` table aside (as `git_sources_old_<time>`) instead of adding its new columns, so Git sources were lost. Columns are now added in place before the table check, and the sources of a set-aside table are restored once when `git_sources` is still empty.

## [11.3.0] — 2026-10-09 — Git: .gitignore helper and optional CI integration (phases 4–5)

### Added
- **.gitignore helper** (`git_gitignore.go`, a new *.gitignore* tab in the Git workspace) per service:
  - **Candidates:** files in state *extra* from the latest checks with their **size** and the **servers** they were seen on, and a suggested pattern (known directories such as `logs/`, known extensions such as `*.log`, or the anchored file); **standard patterns** per stack (Python, Node, Go, Java, Docker, IDE), pre-selected when the repository tree shows the stack; the catalog's **protected globs**; more patterns by hand.
  - The repository's **current `.gitignore`** (in the service's subdirectory) read through the API at the head of the service's branch.
  - **Warnings** about committed secret-looking files (`.env`, `*.pem`, `*.key`, `*credentials*`, `*secret*`, `id_rsa*`, …) and files matching the protected globs, with the advice to commit a `.example` template instead (`config.example.ini`, `.env.example`).
  - The service's catalog **exclude / protected lists** edited on the same screen.
  - **Preview** as a diff against the current file (only new patterns, under a dated comment, line endings kept) and a **download** by default.
  - **Merge request (GitLab) / pull request (GitHub)** on a new branch `wrm/gitignore-<service>-<time>` **only when chosen explicitly**, with a confirmation that names the repository (its name is typed and checked by the server). It needs an optional **write token per Git source** (encrypted, used only for this). Policy **`git_gitignore_mr`** (administrators by default, `WRM_GIT_GITIGNORE_MR`, *Admin → Policies*); `mePayload` flag `git_gitignore_mr`. Audit **`git.gitignore_mr`**.
- **Optional CI integration** (`git_ci.go`; all opt-in, WRM works without any CI):
  - **Incoming webhook per Git source** (`POST /api/hooks/git/<id>`, no session, outside the CSRF check): GitLab push / tag / release events (`X-Gitlab-Token`), GitHub push / release events (HMAC `X-Hub-Signature-256`, `ping` answered), and a generic POST (`X-WRM-Token`, or `X-WRM-Signature` over `<timestamp>.<body>` with `X-WRM-Timestamp` within 5 minutes). It queues an **immediate check** of the affected services (targets refreshed, installations compared, notifications as for periodic checks); deliveries during a check are merged. **Rate-limited** (30 per minute per hook and per client address), **replay protection** by delivery ID (`X-GitHub-Delivery`, `X-Gitlab-Event-UUID` / `Idempotency-Key`, `X-WRM-Delivery`, the generic signature), audit `git.webhook`, `git.webhook_rejected`, `git.webhook_enabled`, `git.webhook_disabled`. The secret is shown once and stored encrypted.
  - **Bundles from CI artifacts** (table `git_feeds`): a URL (GitLab job artifacts API, Jenkins artifact URL, any HTTPS URL) with an optional auth header (encrypted; `Basic user:token` is encoded), fetched by hand or every N minutes by the Git scheduler, imported as an offline bundle when it changed (the feed keeps its newest bundle). Audit `git.bundle_fetched`, `git.bundle_fetch_failed`, `git.feed_added`, `git.feed_changed`, `git.feed_deleted`.
  - **Deploy method "CI pipeline" per service** (table `git_pipelines`, the service's *Edit* dialog, policy `git_update`): *Update* / *Upgrade* trigger a **Jenkins** job (`buildWithParameters` with user + API token or the job's trigger token; queue item and build followed) or a **GitLab pipeline** (trigger token; status read with the source's read token) with `server`, `install_path` and `version`. Same plan, confirmation, policies and audit as a normal update (`git.update` / `git.upgrade` with `via: ci`, link, result); the run history shows the CI status and a link. Audit `git.pipeline_saved`, `git.pipeline_deleted`.
- API: `GET /api/git/gitignore/{app}`, `POST …/preview`, `…/download`, `…/mr`; `POST|DELETE /api/git/sources/{id}/hook`; `POST /api/git/feeds`, `PUT|DELETE /api/git/feeds/{id}`, `POST /api/git/feeds/{id}/fetch`; `PUT|DELETE /api/git/pipelines/{app}`; `POST /api/hooks/git/{hook_id}`. Git sources take `write_token`; the workspace state lists `feeds` and `pipelines` (no secrets).
- Tests: `git_gitignore_test.go` (candidates from a check on the fake SSH host, warnings, preview, download, merge request only with a write token, the explicit request and the typed repository — fake GitLab — and a pull request on the fake GitHub, policy), `git_ci_test.go` (webhook tokens, signatures, rejections, replays and rate limit; artifact feed import, unchanged, error and schedule; Jenkins job and GitLab pipeline runs, failure, policy).

### Changed
- Checks record the **size** of files found only on servers (`wc -c` in the hashing script).
- Write calls to the Git API share the error handling of read calls (`gitProvider.do`).
- The roadmap is done: `ROADMAP.md` keeps a short note for new ideas.

## [11.2.0] — 2026-10-09 — Git: new servers, install and transfer (phase 3)

### Added
- **Install wizard** in the Git workspace (*Installations → Install…*, `git_provision.go`, `git_template.go`): install services on **any SSH connection** into a chosen directory.
  - **Services and versions:** each service at its target (as in the catalog), another branch or tag (Git API) or from an imported **offline bundle** (`ref: "bundle:<source id>"`). **Target directory** per service, suggested from the server roots; system directories are refused.
  - **Pre-checks** on the target (`precheckDirs`): free disk space, a writable directory (or nearest existing parent), `tar`, the `python3` / `node` versions when the service has such files, and an **existing installation** — refused unless the overwrite is confirmed (typed server name), then the replaced files are backed up to `.deploy-bak/`.
  - **Per-host files:** protected files are never taken from the repository. Each is filled **from a template** of the repository (`<name>.example`, `.sample`, `.template`, `.tmpl`, `.tpl`, `.dist`, `.default`, `settings.example.py` → `settings.py`, or the protected file itself; placeholders `{{ name }}`, `${NAME}` / `${NAME:-default}`, `__NAME__`, `@NAME@`, `<NAME>` and `key = value` / `key: value` lines become form fields, INI keys as `section.key`), **copied** from another installation of the same service, **typed**, or left out. Field values can come from the **vault** (password or user name, only for servers the credential's host list allows). Filled files get mode `640`. Typed values and contents are **not stored** with the run.
  - **Writing:** the files travel as one tar stream into a staging directory inside the target, are verified by hash and renamed into place; language checks and the custom check run; any failure rolls back. Then `VERSION.md` and the `updates.jsonl` line, as an update writes them, and the installation is registered (as added by hand) and compared.
  - **Optional unit** (opt-in, with a warning): a systemd unit (`User=`, `WorkingDirectory=`, `ExecStart=`, `daemon-reload`, optional `enable`) or a supervisor program (`supervisorctl reread`), written as root or with `sudo -n`, never over an existing file. **Optional start** afterwards (policy `git_restart`) with the health check of v11.1.
- **Transfer wizard** (*Transfer…* for an installation): copies an installation from server A to server B (or another directory) **with its per-host configuration**, but without logs, `.deploy-bak`, caches, the catalog's `ignore_dirs`, backup-looking paths and the usual prunes (`.git`, `node_modules`, virtual environments). The wizard lists what stays behind. The files stream from A to B through WRM (`tar -c` → `tar -x`, nothing on WRM's disk) with the same pre-checks, overwrite confirmation, backup, verification, checks and rollback; `VERSION.md` gets the new host and user.
- **Policies** `git_install` and `git_transfer` (administrators by default, changeable in *Admin → Policies*, `WRM_GIT_INSTALL` / `WRM_GIT_TRANSFER`); `mePayload` flags `git_install`, `git_transfer`; the workspace state's `deploy` map has `install` and `transfer`.
- **Audit** `git.install` and `git.transfer` per installation (result, files, backup, versions, source, the new installation). Installs and transfers send the run notifications (`git.deploy_started`, `git.deploy_done`, `git.deploy_failed`) with the kind *install* / *transfer*.
- API: `POST /api/git/provision/prepare` (targets, per-host slots with template fields, server roots, pre-checks; or the source files and exclusions of a transfer); `POST /api/git/runs` takes `kind: "install" | "transfer"` with `provision`.
- Tests: `git_provision_test.go` (templates, exclusions, directories, units; install from the fake GitLab with a vault value, a typed file, a systemd unit and a start; refused and confirmed overwrite with backup; install from a bundle with a copied per-host file; transfer between two fake SSH hosts with exclusions and a supervisor program; policies).

### Changed
- Runs carry the source of a transfer (`source_conn_name`, `source_path`) and the installation an install created (`install_id`).

## [11.1.0] — 2026-10-09 — Git: update, upgrade, rollback, restarts (phase 2)

### Added
- **Deploying from the Git workspace** (`git_deploy.go`, `git_runs.go`, new tables `git_runs` and `git_pending_restarts`): *Update*, *Upgrade…*, *Rollback…*, *Restart…* and *Stamp* on the *Installations* tab (for the ticked rows) and in the details of an installation. Everything runs over WRM's SSH connections (jump hosts, vault credentials, proxies) with POSIX tools.
  - **Update** brings an installation to the newest version of the ref it follows; **Upgrade** switches it to another branch or tag (`main`, `tag:v2`) with an extra warning and a summary of the changes (files to change, new, changed locally, no longer in the target); an option makes the new ref the services' ref override. **Rollback** restores a chosen `.deploy-bak` backup. **Stamp** writes `VERSION.md` where an installation is up to date and the file is missing (or overwrites it when forced).
  - **Plan first** (`POST /api/git/plan`): the files of each installation with their state. Default selection: *old* + *missing*. Files changed by hand are written only when ticked explicitly (with a warning); protected files never.
  - **Steps per installation:** a fresh scan (a file changed since the check stops the run), a **dangling-import check** for Python (a `from mod import name` in another file of the installation whose name the new `mod.py` no longer defines; top-level `def` / `class` / assignments / imports, star imports and module `__getattr__` disable it; can be skipped per run), a **backup** to `<install>/.deploy-bak/<bundle id or wrm-YYYYMMDD-HHMM>-<YYYYMMDD-HHMMSS>/<rel>`, **atomic writes** (a temporary file in the same directory copied from the old file, so mode and owner stay, then a rename; the existing line endings are kept, new files follow most files of the installation; the written content is verified by hash), **checks** (`python3 -m py_compile` with the `.pyc` files kept out of the installation, `node --check`, `sh -n` / `bash -n`, or a custom command per service that is remembered), an **automatic rollback** on any failure (files restored, new files and directories removed, the backup dropped), then **`VERSION.md`** and a line in **`updates.jsonl`**.
  - **Compatibility with file-based deploy tools:** `VERSION.md` (Croatian labels, padded per-file table sorted by path, `(per-host)` for protected files) and `updates.jsonl` (one line, sorted keys, Python `json.dumps` style) in the first writable of `/var/log/deploytool/updates.jsonl`, `<install>/.deploy-bak/updates.jsonl` and `/tmp/deploytool-updates.jsonl`. Rollback reads which files a backup holds from WRM's run history or the `updates.jsonl` line; a file that did not exist before is deleted.
  - **Rolling** over several installations, stopping at the first failure; **live progress** per installation and step; *Stop after the current server*.
  - **Double confirmation:** a summary, and for production installations (tag `env:prod` / `prod` or a *prod* environment directory) the server name must be typed — checked again by the server.
  - **Scheduling into maintenance windows:** updates, upgrades and rollbacks run now or at a chosen time. Scheduled runs are stored, survive a WRM restart, can be cancelled, and are skipped when WRM was down past `git_schedule_grace_minutes` (30). A run interrupted by a WRM stop is marked *interrupted*.
  - **Service restart (opt-in, off by default):** the systemd / supervisor units found by the check. Per run: do not restart (the installation shows *restart pending*), restart right after a successful update, at a chosen time, or after a delay. A warning names the units and servers. Scheduled restarts are stored runs too (skipped if the update failed) and are followed by a **health check** (`systemctl is-active`, `supervisorctl status`). `sudo -n` is used when the SSH user is not root and may use it.
  - **Runs tab:** history and scheduled runs with their installations, versions and states; details with every step.
  - **Notifications** (new events): `git.deploy_started`, `git.deploy_done`, `git.deploy_failed` (failed or rolled back), `git.restart_scheduled`, `git.restart_done` (with the health check) and `git.job_upcoming` (15 minutes before a scheduled run).
  - **Policies** (administrators by default, changeable in *Admin → Policies*): `git_update` (also *Stamp*), `git_upgrade`, `git_rollback`, `git_restart`; `git_schedule_grace_minutes`. `mePayload` flags `git_update`, `git_upgrade`, `git_rollback`, `git_restart`.
  - **Audit:** `git.update`, `git.upgrade`, `git.rollback`, `git.restart`, `git.stamp` (one entry per installation with the result, files, backup and versions), `git.run_scheduled`, `git.run_cancelled`, `git.run_skipped`.
  - API: `POST /api/git/plan`, `GET /api/git/installs/{id}/backups`, `GET|POST /api/git/runs`, `GET /api/git/runs/{id}`, `POST /api/git/runs/{id}/cancel`; live update event `git_run_changed`.
- Tests: `git_deploy_test.go` (formats of `VERSION.md` and `updates.jsonl`, dangling imports; against the fake SSH host: update with CRLF, mode and protected files kept, production confirmation, stamp, rollback, failed check → automatic rollback, dangling import, ticked modified files, policies; scheduled restart with notices, health check, cancel and grace period).

### Changed
- The comparison of one installation is shared by checks and runs (`evalInstall`); installations carry `restart_pending`.

## [11.0.0] — 2026-10-09 — Git workspace: compare services (phase 1)

### Added
- **Git workspace** (`git.go`, `git_source.go`, `git_bundle.go`, `git_scan.go`, `git_monitor.go`, `/api/git`): a full-screen workspace opened from the **⎇ Git** button in the top bar. Its code (`static/git.js`) loads only when it is opened. Tabs *Overview*, *Services*, *Installations* and *Settings*. **Read-only on servers** in this version: nothing is updated, installed or restarted.
  - **Sources:** GitLab API v4 and GitHub REST (cloud or self-hosted, groups / organisations with subgroups) with a read-only token, encrypted at rest and never sent to the browser. Self-signed Git servers are pinned on first use. One recursive tree listing per commit (cached forever) gives every file's blob ID; file contents are fetched only for blobs not seen before and cached by blob ID (new tables `git_sources`, `git_trees`, `git_blobs`).
  - **Offline bundles:** import a `.tar.gz` (`<name>/bundle.json` + `payload/<app>/<rel>`, the format of file-based deploy tools; `deploytool.py` and `README.txt` are ignored) when WRM cannot reach the Git server, with its age and source; protected files of the payload are never used. **Export** writes the same format (`bundle.json` with indent 1 and sorted keys, payload without protected files) for servers that only a file-based tool can reach.
  - **Service catalog:** project, ref (`tag:latest`, `tag:vX` or a branch), branch, tag filter, subdirectory = installation root, include / exclude globs, **protected** per-host files, fingerprint files, install hints and kind (app, library, tool). **Suggestion from the repository tree** (subdirectory, fingerprint, protected and exclude globs, kind). Import (replace or merge) and export as catalog JSON. Services of imported bundles are used even when the catalog does not list them.
  - **Refs:** branch-aware tags (only tags whose commit is on the branch count; `tag:latest` falls back to the branch head with a warning), numeric version sort (`v1.10` > `v1.9`), fallback branches, and a ref choice per run: as in the catalog, one branch for all, each project's default branch, or per service.
  - **Discovery over SSH:** one `find` per server over the configured roots (depth 3) recognises installations by their fingerprint files (or install hints), skips copies and backups (`*_BKP`, `backup`, `.deploy-bak`, `kopija`, … as path segments; configurable) and derives an environment label from the path (`/srv/scripts/test/App` → `test`, `App_prod` → `prod`). Installations can also be added by hand.
  - **File-level comparison:** the server hashes its files with POSIX tools (`tr -d '\r' < f | sha256sum`, no Python); WRM compares them with the target and with the **history of every file** (the earlier tags of the branch). File states *ok*, *old* (*n behind*), *modified*, *missing*, *extra*, *protected*; installation states *needs update*, *review*, *up to date*. Details show the version from `VERSION.md`, the systemd / supervisor units that mention the directory, the files and a coloured **diff** (server → target).
  - **Overview:** targets per service (version, latest tag, source, warnings) and a servers × services matrix with filters (state, environment, folder, tag, text).
  - **Periodic checks** while WRM runs (per user, opt-in, every `git_check_interval_minutes`): they never change anything and notify through the notifications module — new events **`git.new_version`**, **`git.drift`** (files changed by hand) and **`git.unreachable`**, one digest per round.
  - Policies `git_enabled` (hides the workspace completely), `git_checks` (who may run checks: everyone by default), `git_check_interval_minutes` (60, 0 = off) and `git_backup_words`; `mePayload` flags `git` and `git_checks`.
  - Audit events `git.source_added`, `git.source_changed`, `git.source_deleted`, `git.catalog_saved`, `git.service_saved`, `git.settings_changed`, `git.bundle_imported`, `git.bundle_exported`, `git.checked`, `git.install_added`, `git.diff_viewed`.
- `tools/check-ui.js` also checks `static/git.js` (syntax and translation keys).
- Tests: `git_test.go` (fake GitLab and GitHub servers: targets, branch-aware tags, fallbacks, blob cache, suggestion, tokens; bundle export / import round trip and unsafe archives; catalog import / export; discovery, comparison, `VERSION.md`, units, diff and notifications against a fake SSH host; policies and ownership).

## [10.10.0] — 2026-10-08 — folder bookmarks

### Added
- **Folder bookmarks** (`bookmarks.go`, `/api/bookmarks`, new table `bookmarks`): named remote directories with an optional colour and note, for **one connection**, the connections of **a folder**, the connections with **a tag**, or **all SSH / SFTP / FTP connections** (global).
  - **File manager:** a ★ button next to the path bar with *Bookmark this directory…*, the bookmarks of the connection and *Manage bookmarks…*. A bookmark whose directory is missing on the server gives a clear message and keeps the current folder open.
  - **Terminal:** ★ in the window's title bar and the bookmarks in the Ctrl+Shift+Space picker type `cd -- '<path>'` (shell-quoted on the server; a leading `~/` stays unquoted so the shell expands it). In a broadcast group, 📣 next to a bookmark sends the `cd` to every terminal, each resolved for its own connection (`POST /api/bookmarks/resolve`).
  - **Path variables** filled in by the server per connection: `~`, `$USER` / `${USER}` / `{user}`, `{host}`, `{name}` — e.g. `/opt/app/servers/{name}/downloads`.
  - **Start directory:** a bookmark can be the start directory of its connections. The file manager opens there; a terminal types the `cd` on connect (also after a reconnect), before the run-on-connect snippets. The most specific one wins (connection, tag, folder, global). *★ Start directory…* in the context menu of a connection and a folder.
  - **Management:** rename, reorder, move between scopes, delete. Global bookmarks can be **shared with share members** (read-only, `GET /api/bookmarks?conn=…&share_token=…`).
  - **Import of WinSCP bookmarks:** *Manage bookmarks → Import WinSCP.ini* (`POST /api/bookmarks/import-winscp`) reads the `[Configuration\Bookmarks\Remote\…]` sections as global bookmarks; existing ones are skipped.
  - Bookmarks are part of the configuration export and import (folders and connections are mapped), copied with *Duplicate*, and removed with their connection, folder or user.
  - Audit events `bookmark.created`, `bookmark.updated`, `bookmark.deleted`, `bookmark.imported`; live update event `bookmarks_changed`.
- The file manager opens `~/…` paths below the home directory (SFTP) or the login directory (FTP).
- Tests: `bookmarks_test.go` (scopes, ownership, validation, variables, `cd` quoting, reorder, duplicate / delete, export / import, share access, WinSCP.ini, start directory on connect).

## [10.9.1] — 2026-10-08 — PuTTY session import

### Added
- **Import from PuTTY** (*Settings → Data*, `POST /api/config/import/putty`, `putty.go`): a registry export (`.reg`, UTF-16 LE with a BOM as written by regedit, or UTF-8) of `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions`.
  - URL-encoded session names are decoded (`My%20Server`); `HostName` (also `user@host`), `PortNumber`, `UserName` and `Protocol` are imported into the folder *PuTTY*. *Default Settings* is not imported but supplies the defaults of the other sessions.
  - SSH and Telnet sessions become connections; raw, rlogin and serial sessions are listed as skipped. `PublicKeyFile` gives a note (convert the `.ppk` with PuTTYgen and add the key).
  - **Proxies:** SOCKS4 / SOCKS5 / HTTP settings (`ProxyMethod`, `ProxyHost`, `ProxyPort`, `ProxyUsername`, `ProxyPassword`, `ProxyDNS`) become saved proxies linked to the connections. One proxy per distinct setting; an identical proxy the user already owns is reused, so importing again creates none. PuTTY's *SSH proxy* becomes the jump host (by name); Telnet and local-command proxies get a note. The policy `proxies` applies.
  - **With mRemoteNG:** connections whose `PuttySession` (also inherited) names an imported PuTTY session get its proxy, whichever file is imported first (new table `putty_sessions`, new column `connections.putty_session`).
  - The import result and the audit event `config.imported` also count `proxies` and `proxy_links`.
- Tests: `putty_test.go` (UTF-16 and UTF-8 files, defaults, proxies and their reuse, SSH proxy → jump host, mRemoteNG linking in both orders).

## [10.9.0] — 2026-10-06 — proxies, folder defaults, notifications

### Added
- **Saved proxies** (`proxy.go`, *🔑 SSH keys & credentials → 🌐 Proxies*, `/api/proxies`): SOCKS5 with an optional user name and password (RFC 1929) and *DNS at the proxy*, SOCKS4 / SOCKS4a, HTTP CONNECT with Basic authentication, and the type **WRM SOCKS tunnel of a connection** (a dynamic tunnel that WRM starts when needed and stops after 30 minutes without traffic).
  - The login is encrypted or a vault credential and never sent to the browser; proxies are shared like credentials (grants, everybody) without revealing it.
  - A **Proxy** field in the connection dialog next to *Connect via (jump host)*, with *＋ New proxy…* and the route of the connection.
  - **Every TCP path goes through it:** the SSH hops in `dialSSH` (terminal, SFTP, search, transfers, tunnels, key deploy, rotation, network tools from a server), FTP/FTPS (control and data connections), the local relay of RDP / VNC / Telnet, web interfaces, live status and *Check now*, Redfish, *Test connection*.
  - **Proxy + jump host:** the proxy is reached through the jump hosts; jump hosts can have proxies themselves. The route is shown and audited, e.g. `bastion → socks5://10.1.1.1:1080 → app-01` (sidebar, dialog, terminal, *Sessions & recordings*, audit details).
  - Errors tell proxy and target apart: proxy unreachable, authentication required / failed, SOCKS5 and SOCKS4 reply codes, HTTP status (`407`, `502`, …).
  - IPMI and Serial-over-LAN (UDP) through a proxy are refused with a clear message.
  - Policy `proxies` (`all` / `admins` / `off`): who may define proxies. Audit events `proxy.created`, `proxy.updated`, `proxy.deleted`, `proxy.granted`, `proxy.grant_revoked`, `proxy.tested`.
  - Export lists the user's proxies (passwords only with secrets) and names them in connections (`proxy_ref`); import creates missing ones and links by name. *Duplicate* keeps the proxy.
  - **Docker:** `docker-compose.yml` maps `host.docker.internal` to the Docker host (`extra_hosts`); WRM warns when a proxy at 127.0.0.1 is saved while it runs in a container (`/.dockerenv`).
- **Folder defaults:** a folder can carry a default jump host and a default proxy (*Folder settings…* in the folder's context menu, `PUT /api/folders/{id}`). Its connections, also new ones, inherit them unless they choose their own (*no jump host / no proxy (not the folder's)*). The connection dialog shows the inherited value (*— as the folder: bastion —*), the tree shows ⤳ after the folder name. Export and import keep them. Audit event `folder.updated`.
- **Notifications module** (`notify.go`):
  - channels configured by administrators (*Admin panel → Notifications*): SMTP e-mail (STARTTLS / TLS), Telegram, Slack / Mattermost / Rocket.Chat (incoming webhooks), Microsoft Teams (Workflows webhook, Adaptive Card), Discord, ntfy, Gotify, Pushover and a generic JSON webhook signed with `X-WRM-Signature: sha256=HMAC-SHA256`; secrets encrypted and never returned; **Send test** per channel; last delivery and last error;
  - users subscribe per event and channel (*Settings → Notifications*), optionally with their own recipient (e-mail, Telegram chat, ntfy topic, Pushover key), and set **quiet hours** in their time zone;
  - one **digest** per user and channel for everything of one check round, at most one message per `notify_min_interval_seconds`; waiting messages are stored and retried with backoff; texts in English or Croatian;
  - first events: a server went down / is up again (live status; owner only, monitoring opt-out respected) and credential rotation due / incomplete;
  - policies `notifications_enabled` and `notify_min_interval_seconds`; audit events `admin.notify_channel_created` / `_updated` / `_deleted` / `_tested`, `notify.settings_changed`, `notify.tested`.
- Tests: `proxy_test.go` (in-process SOCKS5 and HTTP CONNECT proxies with and without authentication for terminal, SFTP, status, jump host + proxy, web interface, desktop relay, WRM tunnel proxy; errors, policy, sharing, Docker warning; folder defaults) and `notify_test.go` (fake HTTP endpoints and a fake SMTP server; digests, rate limit, quiet hours, retries, events).

### Fixed
- **FTP / FTPS behind a jump host** (and now through a proxy): EPSV data connections went to the address the library took from the control connection (the jump host channel or the proxy) instead of the FTP server, and FTPS data connections were not encrypted after `PROT P`. Both are fixed (`routedFTPDial`).

### Changed
- E-mail notifications log in with AUTH PLAIN or, when the server offers only that (e.g. Exchange), AUTH LOGIN.
- `connections.jump_conn_id` NULL now means "as the folder" (still direct without a folder default); deleting a jump host makes its connections fall back to the folder default.
- The admin overview shows proxies and notification channels; the audit filter has *Proxies*, *Folders* and *Notifications*.

## [10.8.1] — 2026-10-05 — fixes: larger imports, install hints

### Fixed
- **mRemoteNG, OpenSSH config and inventory imports of larger files.** The security middleware limited every API body to 1 MB except `/api/config/import`, so `/api/config/import/mremoteng`, `/api/config/import/sshconfig` and `/api/inventory/*` failed with *Bad JSON* from about 1 MB on.
  - One limit for all imports: files up to 20 MB (`maxImportFile`), request bodies up to 32 MB (`maxImportBody`, room for JSON escaping and base64) for `/api/config/import`, `/api/config/import/*` and `/api/inventory/*` (`bodyLimitFor`). The separate inventory limits (10 MB file, 30 MB body) are gone.
  - A file or body above the limit gets **HTTP 413** with its size and the limit (*The file is too large (21.0 MB): the import limit is 20.0 MB*) instead of *Bad JSON* (`decodeImportJSON`).
  - The browser checks the size before sending (mRemoteNG / OpenSSH config import, WRM configuration import, inventory file).

### Changed
- **Install as an app** (*Settings → General*) explains what to do in the current browser:
  - the steps for Chrome / Edge / Brave (the install icon in the address bar or the menu), Safari on macOS (*File → Add to Dock*), iPhone / iPad (*Share → Add to Home Screen*) and Android;
  - the real reason when installing is not possible: plain `http://` to an address other than `localhost` (HTTPS with a trusted certificate is needed), an untrusted certificate, Firefox on the desktop;
  - the *Install* button is shown only when the browser can install right away.
- README: *HTTPS with Caddy* in the reverse proxy section (public name with Let's Encrypt, or `tls internal` for internal names and IP addresses), the import size limit, troubleshooting rows for installing and for too large files.

### Added
- Tests (`import_limits_test.go`), through the real middleware chain: a multi-MB `confCons.xml` and a large CSV import; oversized bodies (with and without `Content-Length`) and oversized files get the clear 413; the browser's limit matches the server's.

## [10.8.0] — 2026-10-02 — quick connect, notes, network tools & installable app

### Added
- **Quick connect** (`quick.go`): a field at the top of the sidebar takes `user@host:port` or a URL-like target (`ssh://`, `sftp://`, `rdp://`, `vnc://`, `telnet://`, `https://`) and connects without the connection dialog:
  - a small dialog for the login (password, vault credential or stored key) and an optional jump host; `POST /api/connections/quick`;
  - the connection is a **temporary** connection of its owner (`connections.temp_until`) — same checks, host keys, recording and audit (`connection.quick`) — listed under *Quick connections*, kept 24 hours after its last use (open terminals keep it alive), at most 20 per user, not monitored, not exported;
  - 💾 / *Save as a connection* turns it into a normal connection (`PUT … {temporary: false}`).
- **Notes & runbook** per connection (`connections.notes`, up to 20 000 characters): a small markdown (headings, lists, bold, code, code blocks, links) rendered as escaped text; edited in the connection dialog or a notes window; 📝 in the sidebar, the terminal window bar and the context menu; owner only (not share members); kept by export / import and *Duplicate*; lists only carry `has_notes`.
- **Network tools** (`nettools.go`, 🧰 *Tools*, `POST /api/nettools`): port check (up to 100 ports and ranges; open / closed / no answer), ping, traceroute (`traceroute`, `tracepath` or `tracert`), DNS (addresses, CNAME, MX, TXT, PTR) and an HTTP/TLS check (status, server, redirect target, TLS version, certificate subject, issuer, expiry and days left, names, trust) — **from the WRM server or from one of the user's SSH connections through its jump hosts** (port and HTTP checks over SSH channels, the other tools with the server's commands); targets are validated as host names or addresses, one run per user at a time, audit event `nettool.run`, policy `network_tools` (off / admins / all, default all); context menu entries *Network tools from this server…* and *Check ports of this host*.
- **Installable app (PWA)** (`pwa.go`): `/manifest.webmanifest` and a service worker (`/sw.js`) that caches the app's own static files per version and shows an offline page; *Settings → General → Install as an app*.
- Tests: quick connect (targets, temporary connections, expiry with and without open terminals, limit, export), notes (save, limit, list flag, duplicate), network tools (ports and HTTP/TLS from WRM and through SSH, DNS, ping and traceroute on a server, validation, ownership, policy) and the PWA endpoints.

### Fixed
- The top bar overlapped its buttons at medium widths (and the version label never collapsed); it now gives way step by step.

## [10.7.0] — 2026-10-02 — out-of-band management, consoles & serial ports

### Added
- **Out-of-band management** (`bmc.go`): a BMC per connection (`connections.bmc`, JSON with the password encrypted, or a vault credential):
  - **Redfish** (iDRAC, iLO, XClarity, Supermicro, OpenBMC): power state, health, maker, model, serial number, BIOS, CPUs, memory, BMC firmware, power draw, inlet temperature; certificates pinned on first use like SSH host keys (shared with FTPS: `pinnedTLSConfig`);
  - **IPMI** with `ipmitool` (`chassis status`, `mc info`, `fru print`), the password in the environment (`-E`);
  - power actions on, graceful shutdown/restart, force off/restart, power cycle and NMI — only those the BMC offers — with confirmation and audit events `bmc.power` / `bmc.power_failed`;
  - through the connection's jump hosts: Redfish over an SSH channel of the last hop, `ipmitool` on the last hop (password over stdin).
- **Serial consoles** in terminal windows (`console.go`): over the BMC's SSH with the vendor's console command (iDRAC `console com2`, iLO `vsp`, …, suggested from the BMC maker), or IPMI **Serial-over-LAN** on the jump host or on the WRM server in a pseudo terminal; owner only, audited and recorded (`protocol` console / sol); saved in workspace sessions.
- **Serial ports** of the WRM server: protocol `SERIAL` (device name, speed, data bits, parity, stop bits, flow control), raw terminal, *Test connection*; policy `serial_ports` (off / admins / all, default admins), Linux.
- Policies `bmc_enabled` (`BMC_ENABLED`) and `serial_ports`; ipmitool in the admin overview; audit filters for `bmc.` and `inventory.`.
- Tests: Redfish against a fake BMC (status, actions, pinning, changed certificate, wrong password, jump host, credential), IPMI with a fake ipmitool (locally and on a jump host), Serial-over-LAN (PTY and jump host), the BMC SSH console and a serial port on a pseudo terminal.

### Changed
- The terminal WebSocket runs on terminal backends (SSH shell, BMC console, SOL, serial port) instead of an SSH session only.

### Fixed
- Saved workspace sessions did not include remote desktop windows.

## [10.6.0] — 2026-10-02 — tags & inventory import (CSV, Excel, NetBox)

### Added
- **Tags** on connections (`connections.tags`): free tags and tags with a key (`env:`, `site:`, `rack:`, `role:`, `tenant:`, `platform:`, `cluster:`), normalised to lower case, at most 30:
  - a tags field with suggestions in the connection dialog, bulk tagging (`POST /api/connections/bulk` action `tag`, add / remove);
  - a tag filter bar in the sidebar (several tags = all must match), tag search (`#tag`, several words);
  - environment badges (PROD / staging / test / dev colours) and a red edge for production connections;
  - kept by export / import and *Duplicate*.
- **Inventory import** (`inventory.go`, *Settings → Data → Import inventory*):
  - **CSV** with automatic separator (`,` `;` tab `|`, Excel's `sep=` line) and encoding detection (UTF-8, UTF-16, Windows-1250), and **Excel .xlsx** (read with archive/zip + encoding/xml, sheet selection, shared and inline strings);
  - a column mapping guessed from English and Croatian header names, checked and changed in a preview;
  - host fields with ports, users, CIDR suffixes, IPv6 and URLs; Windows platforms become RDP; environment / site / rack / role / platform columns become tags; jump hosts by name;
  - **NetBox** devices and virtual machines through the REST API (filters, all pages, primary IP or DNS name, environment custom field), with **sync**: `connections.ext_id` remembers the NetBox object, a later import updates address, NetBox tags and folder (`ext_tags`) and lists objects that disappeared; address and token can be remembered per user (token encrypted, table `inventory_sources`);
  - for both: the login of the imported connections (none, a **vault credential** or a **stored SSH key**), default user and protocol, folders by column / site / role / tenant / one folder, extra tags.
- Audit events `inventory.imported`, `inventory.netbox_imported`, `connection.tagged`.
- Tests: tags (normalisation, bulk, export/import), CSV encodings and separators, an XLSX import, NetBox import and sync against a fake NetBox with pagination.

### Fixed
- Texts that still said mRemoteNG RDP/VNC/Telnet connections are skipped (they are imported since 10.4).

## [10.5.0] — 2026-10-02 — SSH keys & credentials vault

### Added
- **SSH key store** (🔑 *Keys* in the top bar, `keys.go`):
  - generate ED25519, RSA 2048/3072/4096 or ECDSA 256/384/521 keys, or import OpenSSH/PEM private keys (with passphrase) and public keys; private keys are encrypted at rest;
  - new auth method **SSH key from the key store** (`KEY_REF`, `connections.key_id`); last use is recorded;
  - **deploy** a key to many SSH connections at once like `ssh-copy-id` (through jump hosts, idempotent, keeps other lines, restores SELinux contexts) and optionally switch them to the key after a test login;
  - **revoke** a key from many servers (never the key WRM logs in with);
  - **who has access**: read a connection's `~/.ssh/authorized_keys`, match the keys of the key store, mark WRM's login key, remove single entries;
  - public key copy/download, rename, delete (refused while in use), export of the private key with re-authentication and an optional new passphrase.
- **Credentials vault** (`credentials.go`):
  - a credential holds a user name with a password and/or a stored key; new auth method **Credential from the vault** (`CREDENTIAL`, `connections.credential_id`) for SSH, SFTP, RDP, VNC, Telnet and FTP;
  - sharing with selected users (administrators: everybody) without revealing the secret; host lists (glob patterns, CIDR) limit where it can be used — for grantees also the jump hosts;
  - **password rotation** on all SSH servers that use it: pre-flight login check, `passwd` over a PTY, verification of the new login, all-or-nothing with rollback, a *pending* password if a rollback fails, live progress per server; a check-only mode;
  - rotation reminders (badge, admin overview warning), owner-only *Show* with re-authentication, usage list.
- Audit events `ssh_key.*` and `credential.*`; audit filter entries for both; admin overview card with key, deployment and credential counts.
- Bulk action **🔑 Key** for selected connections; context menu entries *Who has access…* and *Deploy SSH key…*.
- Export/import reference keys and credentials by name (`key_ref`, `credential_ref`) instead of copying them.
- Tests with a fake SSH server that runs scripts in `sh` and simulates `passwd`: key store, deploy/revoke/who has access, vault sharing, host and jump-host restrictions, rotation with success, rollback and incomplete rollback.

### Changed
- Deleting an account also deletes its snippets, tunnel definitions, keys and credentials (and its grants).
- RDP, VNC, Telnet and FTP connections offer only password and vault logins in the connection dialog.

## [10.4.0] — 2026-10-02 — RDP, VNC & Telnet in the browser

### Added
- **Remote desktop** connections — **RDP**, **VNC** and **Telnet** — in WRM windows, through guacd (Apache Guacamole proxy daemon):
  - WRM performs the guacd handshake with the stored credentials and options and relays the Guacamole protocol over a WebSocket (`/ws/desktop`, subprotocol `guacamole`);
  - browser → guacd input is limited to an allow-list of instructions; internal pings are answered by WRM;
  - guacamole-common-js 1.5.0 (Apache 2.0) is bundled and loaded on demand.
- **Desktop windows:**
  - scale-to-fit or 1:1 with scroll bars, full screen, dynamic resolution for RDP (display update);
  - keyboard and mouse (also touch), **Ctrl+Alt+Del**, **send to the remote clipboard**, **type text as keystrokes**;
  - remote clipboard into WRM's clipboard panel, sound from RDP, prompts for credentials the server requires, reconnect;
  - tabs with connection state; restored by saved workspace sessions.
- **Connection options:**
  - RDP: domain, security (NLA/TLS/RDP/Hyper-V), keyboard layout, resize behaviour, colors, start program, certificate, console session, sound, clipboard, wallpaper, **RD Gateway**;
  - VNC: colors, pointer, view only, clipboard;
  - Telnet: font size, colors, login/password prompt patterns.
- **Behind jump hosts:** a temporary tunnel through the jump hosts is opened for guacd for the duration of the session.
- **Audit and recording:**
  - desktop sessions in *Sessions & recordings* (`protocol` rdp/vnc/telnet) and audit events `desktop.open`, `desktop.close`, `desktop.error`;
  - the screen stream is recorded (gzip, SHA-256) and replayed in the browser with play/pause/seek, or downloaded as `.guac`.
- **Admin overview** shows whether guacd answers and its version; *Test connection* checks the port and guacd.
- **Policies** `desktop_enabled` (`REMOTE_DESKTOP_ENABLED`), `guacd_address` (`GUACD_ADDRESS`), `desktop_tunnel_bind`.
- `docker-compose.yml` contains a `guacd` service sharing WRM's network.
- **mRemoteNG import** brings RDP (domain, console, colors, RD Gateway, sound), VNC (view only) and Telnet connections instead of skipping them.
- **Live status** checks RDP/VNC/Telnet ports too.
- **Tests:**
  - Guacamole codec;
  - desktop relay against a fake guacd (handshake parameters, allow-list, ping, wrong password, policy, recording and download);
  - VNC through a jump host.

### Fixed
- SQLite `busy_timeout` is now set on every pooled database connection (it reached only one before), avoiding rare "database is locked" errors under concurrent load.

## [10.3.0] — 2026-10-02 — snippets, broadcast input & live status

Daily terminal work for data center admins (extension plan phases 4–6, first part). Release names no longer carry a suffix: this release is `v10.3.0`.

### Added
- **Snippets** — saved commands:
  - in the right panel (⚡ tab, grouped, searchable) and in a quick picker in every terminal (**Ctrl+Shift+Space** or ⚡ in the title bar);
  - click runs in the focused terminal, Shift+click inserts without Enter;
  - scope: all connections, one folder, or one connection; multi-line commands;
  - variables `{{host}} {{port}} {{user}} {{name}} {{folder}} {{wrm_user}} {{date}} {{time}}` and prompts `{{?Label}}` / `{{?Label=default}}` with a preview;
  - shared snippets for all users (administrators); example set; export/import; copied with *Duplicate*.
- **Run on connect:** snippets typed into the shell by the server once a terminal of a matching connection is connected (also after reconnects), e.g. `sudo -i`, `cd /srv/app`, `tmux attach`. Created from the connection or folder context menu; audited as `terminal.auto_run`.
- **Broadcast input:**
  - type into several terminals at once (📣 in the top bar, per-window toggle, orange frame and status bar);
  - dangerous commands (rm -rf, shutdown, mkfs, dd, iptables -F, DROP TABLE, kubectl delete, …) and multi-line pastes ask first, with Cancel as the default;
  - snippets run on every terminal of the broadcast with their own variables;
  - every start and stop is audited per terminal (`terminal.broadcast`); policy `broadcast_enabled` (`BROADCAST_ENABLED`).
- **Live up/down status:**
  - a background monitor checks every saved connection (TCP connect; SSH and FTP read the greeting and close politely, like Nagios check_ssh) once per host:port and round, retrying before reporting down;
  - status dots in the sidebar with latency, banner and "since"; down counter per folder; "show only down" filter;
  - notifications (also desktop notifications in the background) when a server goes down or comes back;
  - *Check status now* for a connection or a whole folder, also through jump hosts;
  - connections behind jump hosts show their jump host's state, or are checked through it (`status_jump_checks`);
  - per-connection opt-out (*Monitor up/down status*); policies `status_enabled` (`STATUS_ENABLED`), `status_interval_seconds`, `status_jump_checks`; API `GET /api/status`, `POST /api/status/check`; event `status_changed`.
- **Tests** for snippets (API, sharing, validation, variables, run on connect, export/import), broadcast audit and the status monitor (up/down, banners, jump hosts, check now, FTP greeting).

### Changed
- Release tags, binaries and the Docker image no longer have a suffix after the version number (`v10.3.0`, `wrm-pro-v10.3.0-linux-amd64`). CI names the release files after the tag.
- Confirmation dialogs for dangerous actions focus *Cancel*, so Enter cannot confirm them by accident.
- The "terminal connected" marker in the sidebar is blue (green is now the live status).

### Fixed
- Deleting several connections at once also stops their tunnels, removes their tunnel definitions and detaches connections that used them as jump host (single deletes already did).

## [10.2.0] — 2026-10-02 — jump hosts, SSH tunnels & mRemoteNG import

Phase 5 (fleet operations) of the extension plan, first part: what mRemoteNG, PuTTY and `ssh -J` users need to manage servers behind bastions (see [ARCHITECTURE.md](ARCHITECTURE.md)).

### Added
- **Jump hosts** (ProxyJump):
  - any connection can go through another saved SSH connection, in chains of up to 5 hops; loops, foreign and non-SSH jump hosts are refused;
  - used by terminals, the file manager (SFTP and FTP/FTPS), search, editor, server-to-server transfer, *Test connection*, tunnels and web interfaces;
  - every hop authenticates with its own credentials and its host key is verified;
  - the route is shown in the sidebar (⤳), connection test, terminal, session history (`jump_path`) and audit log.
- **SSH tunnels** (port forwarding) per connection:
  - local (`-L`), remote (`-R`) and dynamic SOCKS5 (`-D`) tunnels, through the jump hosts too;
  - listen address and port (0 = automatic), target, *Open as* http/https + path, start mode *manual*, *with terminal* or *always* (also after a restart);
  - templates: web interface (HTTPS/HTTP), SSH, RDP, VNC, PostgreSQL, MySQL/MariaDB, iDRAC/iLO/IPMI, SOCKS;
  - a hint in plain words for every tunnel;
  - each tunnel keeps its own SSH connection with a 30 s keepalive and automatic reconnect (backoff 2–60 s) while its port stays open;
  - statistics: open and total connections, bytes up/down;
  - **🔀 Tunnels** panel (top bar, with a counter of running tunnels): state, listen → target, traffic, *Open*, *Copy address*, *Start/Stop*; administrators can show and stop the tunnels of all users;
  - connection context menu: *Tunnels*, *Start / Stop all tunnels*, *Add tunnel…*;
  - live updates through `/ws/events` (`tunnels_changed`).
- **Web interface connections** (protocols `HTTPS` / `HTTP`, with a path): a double-click opens the page in a new tab, directly or, behind a jump host, through an automatic temporary tunnel (30 minutes idle limit). *Test connection* checks that the port answers.
- **Import from mRemoteNG** (`confCons.xml`):
  - AES-GCM and legacy AES-CBC passwords, master password (checked with `Protected`, asked for when needed), full-file encryption;
  - inheritance of user name, password, port, domain and SSH tunnel;
  - containers become folders; `SSHTunnelConnectionName` becomes the jump host;
  - unsupported protocols and duplicates are reported as skipped.
- **Import from OpenSSH config:** `Host` blocks, `Host *` and wildcard defaults, `HostName`/`User`/`Port`, `ProxyJump` (implicit jump hosts are created), `ProxyCommand ssh -W`, `LocalForward`/`RemoteForward`/`DynamicForward` (tunnels that start with the terminal), `IdentityFile` (when server key files are allowed).
- **Policies:** `tunnels_enabled` (`TUNNELS_ENABLED`), `tunnel_users`, `tunnel_bind_any`, `tunnel_remote_forward`, `tunnel_idle_minutes`. Running tunnels that a policy change no longer allows are stopped.
- **Audit events:** `tunnel.configured`, `tunnel.start`, `tunnel.stop` (with traffic), `tunnel.error`, `tunnel.reconnected`, `web.open`; connection changes record the jump host.
- **Admin panel → Overview:** active SSH tunnels with *Stop*. **Admin status API** includes the tunnels.
- **API:** `GET/PUT /api/connections/{id}/tunnels`, `GET /api/tunnels[?all=1]`, `POST /api/tunnels/{key}/start|stop`, `POST /api/connections/{id}/open-web`, `POST /api/config/import/mremoteng`, `POST /api/config/import/sshconfig`. Connections have `jump_id`, `web_path`, `route` and `tunnels`.
- **Tests:**
  - the in-process SSH test server also handles `direct-tcpip` and `tcpip-forward`;
  - new tests for two-hop jump chains (terminal, SFTP, connection test, validation), local/SOCKS/remote tunnels (traffic, ownership, policies, audit, start modes) and web interfaces;
  - new tests for mRemoteNG import (both encryption formats, master password, inheritance) and OpenSSH config import.

### Changed
- Export and import of WRM's own format include jump hosts and tunnels (IDs are remapped).
- *Duplicate* copies the jump host, path and tunnels of a connection. Deleting a connection stops its tunnels and makes connections that used it as jump host direct.
- Editing a connection restarts its running tunnels with the new settings.
- The connection dialog is wider and fully translated (labels were partly English in Croatian).

### Fixed
- The connection dialog no longer moves the focus back to *Name* when you already started typing in another field.

## [10.1.0] — 2026-10-01 — audit trail & session recording

First phase of the extension plan: WRM as a lightweight PAM for teams (see [ARCHITECTURE.md](ARCHITECTURE.md)).

### Added
- **Session recording** of every SSH terminal in asciinema format (asciicast v2, gzip, SHA-256 in the database). Recording is written by a background goroutine through a bounded queue, so the terminal is never slowed down. There is a size limit per session (`recording_max_mb`) and a retention period (`recording_retention_days`).
- **Replay in the browser** (play/pause, seek, 0.5–16×, skip idle time) and `.cast` download. Viewing and downloading are audited.
- **Optional keystroke recording** (`session_recording_input`, off by default). Typing at password, passphrase and PIN prompts is masked.
- **● REC** badge and a *“This session is recorded”* notice in the terminal.
- **Terminal sessions log** (`terminal_sessions`): user, IP, server, remote user, start/end, duration, status (closed, failed, connection lost, ended by an administrator, interrupted by a restart), exit code, bytes. Sessions left open by a crash are closed at the next start, and their partial recordings are kept.
- **File transfer log** (`file_transfers`): every upload, download (also editor opens and each file of a ZIP download) and server-to-server copy, with source, destination, size, SHA-256 and status.
- **Admin panel:** *Sessions & recordings* and *File transfers* tabs with filters and CSV export. The *Audit log* tab gained date range filters, recording events, links to sessions and **Verify integrity**.
- **Settings → Session history:** your own sessions and sessions on your connections.
- **Audit integrity:**
  - audit entries now reference the connection and the terminal session;
  - every entry carries a SHA-256 hash chain, checked by `GET /api/admin/audit/verify`;
  - database triggers make `audit_log`, `file_transfers` and `session_recordings` append-only (no UPDATE; no DELETE of entries younger than 7 days);
  - secret-looking keys in audit details are redacted.
- **Feature flags:** `audit_enabled` (`AUDIT_ENABLED`), `session_recording` (`SESSION_RECORDING_ENABLED`), `session_recording_input`, `recording_retention_days`, `recording_max_mb`. Environment aliases are supported for the plan's flag names.
- **API:** `GET /api/recordings`, `/api/recordings/{id}`, `/api/recordings/{id}/cast`, `GET /api/admin/transfers`, `GET /api/admin/audit/verify`, `GET /healthz`.
- **Operations:**
  - `Dockerfile` (unprivileged user, `/data` volume, health check) and `docker-compose.yml`;
  - CI builds and smoke-tests the image and publishes release images to `ghcr.io/vedranius/wrm-pro`;
  - `wrm -healthcheck`;
  - `WRM_RECORDINGS_DIR`.
- **Brand:** new WRM PRO logo for the favicon, app icons, sign-in screen, top bar, sidebar, README, release notes and GitHub social preview. A *WRM Orange* accent color. The logo kit is in `docs/brand/`.
- **Documentation:**
  - `ARCHITECTURE.md` (stack, auth, SSH layer, credentials, collaboration, migrations, extension plan status and touch points);
  - this changelog;
  - README, SECURITY and CONTRIBUTING updated.
- **Tests:**
  - an integration test with an in-process SSH + SFTP server (connect → audit entries → replayable recording without passwords; transfer checksums; access to recordings);
  - unit tests for the recorder (UTF-8 split across chunks, masking, resize, size limit) and for the append-only, hash-chained audit log.

### Changed
- Downloads by connection owners are now audited too (before, only downloads through shares were).
- Server-to-server transfer reports a read error as a failed file instead of a partial "ok".
- The audit log CSV export has new columns: connection, session and hash.

## [10.0.1] — 2026-09-28

### Fixed
- "404 page not found" at `/static/`, the v9 address of the app. It now redirects to `/`.
- Upgrading a database from an earlier build with an incompatible `audit_log` / `collab_messages` table printed `Index warning: no such column: ts / share_id`, and the audit log and chat history could not be written. Such tables are now kept as `<table>_old_<time>` and created again.

### Added
- The startup log shows the address to open in the browser.

## [10.0.0] — 2026-09-28

### Added
- **Voice calls** in collaboration rooms (WebRTC):
  - mute, deafen, push-to-talk, device selection, noise suppression, echo cancellation and gain control;
  - speaking indicators, per-person volume, connection quality;
  - a built-in TURN relay (UDP/TCP 3478) with short-lived credentials.
- **Users and roles:**
  - admin user management (create, disable, reset password/2FA, sign out everywhere, unlock, delete);
  - share access modes (members only / signed-in users / anyone with the link);
  - five roles enforced by the server (observer, viewer, operator, moderator, owner);
  - members, people history, bans, expiry, link rotation, immediate revocation.
- **Collaboration:** collaboration bar, people panel with moderation, chat history, file exchange, raise hand, several shared terminals with screen snapshot, remote keyboard control.
- **Security:**
  - TOTP 2FA with recovery codes and a 2FA policy; closed self-registration; account lockout; hashed session tokens with idle/maximum lifetime; device list;
  - a random AES-256-GCM key per installation (migrated from the v9 default key); secrets never sent to the browser; secret export only after re-entering the password;
  - SSH host key verification (trust on first use, or strict) and FTPS certificate pinning;
  - CSRF protection, CSP and security headers, request limits, local assets (no CDN), self-signed HTTPS option, trusted proxy support;
  - an audit log with CSV export.
- **Admin panel:** overview with warnings, users, shares, policies, voice & network, host keys, audit log. Policies can be forced with `WRM_<KEY>` environment variables.
- `wrm -version`, `wrm -reset-password USER [-reset-2fa]`; unit tests in CI with the race detector; `SECURITY.md`.

## [9.10.2]

### Security
- In a collaboration session any participant could give themselves keyboard control of another participant's shared terminal. Now only the sharer can grant control, keystrokes go only to that sharer, and control ends when sharing stops.

### Fixed
- The *Control* button appears next to participants while you share a terminal; late joiners see *Watch* at once; duplicate names get a suffix.

## [9.10.1]

### Added
- Saved sessions restore all windows (connection, mode, folder, position, snap, minimized, focus).
- Reconnect button on SSH windows (Ctrl+Shift+R).
- Top/bottom half snap targets, drag-to-edge snapping, snap layouts menu.

### Fixed
- Windows stay inside the screen when the zoom, resolution or panel sizes change.

## [9.10.0]

### Added
- Terminal: reconnect and automatic reconnect, connection state, flow control, SSH keepalive, search in output, clickable URLs, configurable scrollback.
- File manager: instant filter and recursive name/content search, breadcrumb path, multi-select, keyboard shortcuts, folder and drag & drop uploads with progress, built-in text editor, pooled SFTP connections.
- Connections: test connection, duplicate, FTPS. Modern dark responsive UI with an accent color; mobile drawers.
- Security: login rate limiting, WebSocket origin check, optional HTTPS.

### Fixed
- Disconnects on large outputs with multi-byte characters (binary frames).
- Full-screen programs drawn at the wrong size (the PTY size follows the window).
- Folder transfer when WRM runs on Windows; several smaller issues.

[10.7.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.7.0
[10.6.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.6.0
[10.5.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.5.0
[10.4.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.4.0
[10.3.0]: https://github.com/vedranius/web-browser-RDM-public/releases/tag/v10.3.0
[10.2.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.2.0&expanded=true
[10.1.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.1.0&expanded=true
[10.0.1]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.0.1&expanded=true
[10.0.0]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v10.0.0&expanded=true
[9.10.2]: https://github.com/vedranius/web-browser-RDM-public/releases?q=v9.10.2&expanded=true
