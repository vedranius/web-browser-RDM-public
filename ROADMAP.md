# Roadmap

Accepted work through **v12.6.8** is delivered in [CHANGELOG.md](CHANGELOG.md).

## v12.6.9 – Workboard and addons (agreed, next)

Goal: a **Workboard** page of small tiles (addons) for a better overview, where each tile can also be pulled out as a panel of the v12.6.8 panel layout (docked, stacked or floating). Everything is off by default and added from an addon catalog, so nobody is flooded with tiles.

First addons:
- **Server status:** up/down and latency of chosen connections (from the live status).
- **Disk, CPU and RAM:** bars and short history for chosen servers (from Monitoring / host metrics).
- **Notifications:** the latest WRM notifications and events.
- **Clock:** local time and chosen time zones.
- **Latest commits:** the last ~10 commits pushed to chosen GitLab / GitHub projects of the Git workspace sources.
- More where useful and quiet (for example tunnels and transfers in progress, certificate expiry of HTTPS connections, notes).

The Workboard and tile settings are stored per user (`/api/me/ui/workboard`).

## Medium without Docker (accepted direction, versions to follow)

Goal: Medium keeps the original deploy tool's interface, options and `deploy-config` format, but runs on WRM's own engine. No Docker runtime, no Python virtual environment, no dependency installation, and the same behaviour whether WRM runs on Linux or Windows.

Agreed design:
- **Interface:** the original tool's screens (overview, comparison, details/diff, plan and exclusions, dry run, confirmations, progress, lock, history and rollback, server/environment/project settings, YAML editor, config Git actions) stay as they are. WRM serves them and answers their requests from Go.
- **Configuration:** `deploy-config` (servers, environments, `projects/*.yaml`) stays the source of truth and is read without conversion. The editor keeps comments and layout; validation runs before every deploy.
- **Connections and keys:** destinations use the connections already saved in WRM (keys, vault/Bitwarden keys, host-key checks, jump hosts) instead of a separate SSH profile ZIP. Keys never leave WRM's store.
- **Transfer and state:** WRM's native SSH transfer, lock, backup, integrity check and rollback, writing the server-side state (`current.json`, `history.jsonl`, backups) in the format the original tool reads, so existing history stays valid.
- **Where builds run** (per project, shown in the plan before anything runs):
  - on the WRM host (for example WRM on a Linux server with the needed toolchain);
  - on a chosen Linux build server over SSH, in a temporary directory, with only the prepared source sent (no Git or SSH keys); the result is then transferred to the destinations;
  - no build (deploy the selected files as today).
  Platform and architecture of the build server are checked against the destinations.
- **Trust:** build and `post_deploy` commands still run code from the project/config; Medium stays limited to trusted operators.

Delivery plan (v12.6.0 delivered the deploy-config reader/validator, the native plan and the comparison with the original tool; see [CHANGELOG.md](CHANGELOG.md)):
1. **v12.7.0 – native engine:** deploy, dry run, lock, backup, history and rollback through WRM connections; build location choice; per-project switch between the Docker runtime and the native engine. The Docker path keeps working until the native engine matches it.
2. **v12.8.0 – switch-over:** native engine as the default on Linux and Windows; the Docker runtime becomes optional and is then retired.
