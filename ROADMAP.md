# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

## v11.4.0 — Git workspace: every file of an installation, partial and incremental checks

### A. Always show all files of an installation folder
- The installation detail lists **every file** in the server folder (recursive, honouring `ignore_dirs`, with a cap and a "truncated" notice), not just the tracked ones, plus the files that exist only in Git at the target ref.
- Each row: path, size, modification time, state, and **why** it is in that state: tracked / excluded by the catalog (which pattern) / not in the include list / protected (which pattern) / extra (not in Git) / missing on the server / symlink (target shown, never followed outside the root) / unreadable (permission denied).
- Filters on the file list: all / differing / tracked only / not tracked; text search.
- Any file can be **opened** (view only, size-limited; binary files are detected and compared by hash only).
- **Diff server ↔ Git** for any file whose path exists in the repository at the target ref, including files the catalog excludes (labelled "excluded by catalog — informational, not updated"). Side-by-side and unified view, CRLF-insensitive like the hashing, in dark and light theme and at 390 px.
- Differences of excluded / not tracked files never change the installation state; tracked files keep the current semantics.
- **Commits between** the server's version (commit from `VERSION.md`) and the target ref, through the Git provider API, with "N commits behind".
- Server content is read over the existing SSH layer: read only, safely quoted paths, nothing new cached on disk.

### B. Partial and incremental checks
- Overview matrix: filters by service, server / host, environment and state (review / needs update / up to date / error / not checked).
- Selection with checkboxes (rows, columns, cells): "Check selected", "Check visible", "Check this cell", "Check only stale (older than X / never checked)"; "Check all" stays.
- Results **stream in** per installation (a check job polled by the browser), with a spinner per cell and a cancel button; bounded concurrency per server and overall (settings with sensible defaults).
- Cells that are not checked keep their previous results with "checked X ago"; a partial check never wipes other results; the last good result is kept (additive migration).
- The same selection options on the Installations tab.

### C. Quality
- Go tests with the fake Git provider and the fake SSH host (generic fixtures): an excluded file is still diffable, reason labels, symlink / unreadable handling, a partial check keeps other results, the concurrency limit, cancel.
- Browser test (Croatian texts, 390 px, top bar at 1100–1920 px, no console errors); en + hr translations for every new string.

## v11.5.0 — Git workspace: environments, deploy plan and safe deploys

- **Environments**: named groups of ordered destinations (`server:/path`) per service, with per-environment rules: require typing the environment name to confirm, `allowed_refs` (e.g. `refs/heads/main`, `refs/tags/v*`), extra per-environment ignore patterns, an optional `post_deploy` command (explicit, shown before it runs), keep the last N backups.
- **Deploy plan** per server before any update: files to change / add / delete (deleting files removed from the repository is opt-in), excluded files; a plan fingerprint — only the reviewed plan is applied; the plan is computed again under the lock and the deploy stops if it differs.
- **Per-run file exclusion**: a checkbox per changed file to leave it out of this update; remembered for the next compare with who and when; easy to clear.
- **Server-side lock** against concurrent updates from several WRM instances or colleagues: an atomic `mkdir` lock in the state directory with an owner file (who, when, from where); a stale lock is removed only after explicit confirmation.
- **Server-side state** next to the installation (compatible with the existing `VERSION.md` and `.deploy-bak` formats): current version / commit / branch / deployed by / at, an append-only history, backups with a list of added files so a rollback can remove them; rollback in reverse destination order.
- **Overview by environment**: version, commit, who / when, commits behind, and a warning when the servers of one environment differ from each other.
- **Access doctor**: per destination a check of SSH, write permission, the state directory, unreadable subdirectories, a clean non-interactive shell (no banner / rc output that breaks commands), and a warning when the path looks like a public web root.
