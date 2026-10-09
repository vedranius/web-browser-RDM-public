# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

## v11.5.0 — Git workspace: environments, deploy plan and safe deploys

- **Environments**: named groups of ordered destinations (`server:/path`) per service, with per-environment rules: require typing the environment name to confirm, `allowed_refs` (e.g. `refs/heads/main`, `refs/tags/v*`), extra per-environment ignore patterns, an optional `post_deploy` command (explicit, shown before it runs), keep the last N backups.
- **Deploy plan** per server before any update: files to change / add / delete (deleting files removed from the repository is opt-in), excluded files; a plan fingerprint — only the reviewed plan is applied; the plan is computed again under the lock and the deploy stops if it differs.
- **Per-run file exclusion**: a checkbox per changed file to leave it out of this update; remembered for the next compare with who and when; easy to clear.
- **Server-side lock** against concurrent updates from several WRM instances or colleagues: an atomic `mkdir` lock in the state directory with an owner file (who, when, from where); a stale lock is removed only after explicit confirmation.
- **Server-side state** next to the installation (compatible with the existing `VERSION.md` and `.deploy-bak` formats): current version / commit / branch / deployed by / at, an append-only history, backups with a list of added files so a rollback can remove them; rollback in reverse destination order.
- **Overview by environment**: version, commit, who / when, commits behind, and a warning when the servers of one environment differ from each other.
- **Access doctor**: per destination a check of SSH, write permission, the state directory, unreadable subdirectories, a clean non-interactive shell (no banner / rc output that breaks commands), and a warning when the path looks like a public web root.
