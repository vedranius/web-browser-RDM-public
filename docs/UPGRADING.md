# Upgrading WRM PRO

In general: replace the executable, AppImage or container image and keep the database and the encryption key. Database changes are additive, so the previous version still starts on the database after an upgrade. The notes below list what changes per version, newest first; see also the [changelog](../CHANGELOG.md).

## Upgrading from v14.11.7

Replace the executable/AppImage/container and keep the database and encryption key. No database migration. Two new settings: `appearance_lock` (`0`) and `appearance_shared_themes` (`[]`, managed from *Settings → Appearance*); `/api/auth/me` returns them in `appearance.lock` and `appearance.shared`. Nothing changes until an administrator uses them.

## Upgrading from v14.11.6

Replace the executable/AppImage/container and keep the database and encryption key. No database, API or setting change.

## Upgrading from v14.11.5

Replace the executable/AppImage/container and keep the database and encryption key. No database or API change. The clipboard history in the browser is converted to the new format on first load; pins and the recently used snippets are kept in the browser only.

## Upgrading from v14.11.4

Replace the executable/AppImage/container and keep the database and encryption key. The database gets one new column, `share_links.allow_tunnel_edit` (off); an older binary still starts on it. Existing shares keep exactly what they allowed: *start and stop tunnels* stays as it was, and adding or editing tunnels stays off until the owner turns it on.

## Upgrading from v14.11.3

Replace the executable/AppImage/container and keep the database and encryption key. No database change. New MCP tool `wait_for_approval`; while a tool call waits for an approval WRM sends a `notifications/message` (logger `wrm.approval`) on the call's event stream. Update the WRM PRO Claude plugin to 1.2.4 (`/plugin` in Claude Code, or upload the new `claude-plugin-wrm-pro.zip`) and download the Claude Desktop extension again to get the desktop notifications; set `WRM_MCP_NOTIFY=0` for the bridge to turn them off.

## Upgrading from v14.11.2

Replace the executable/AppImage/container and keep the database and encryption key. No database change and no API change. The WRM PRO Claude plugin is now 1.2.3: update it in Claude Code (`/plugin`) or upload the new `claude-plugin-wrm-pro.zip`; download the Claude Desktop extension again from *Settings → AI connections* to get the same bridge. `~/.wrm/mcp.json` stays as it is, so no new sign-in is needed.

## Upgrading from v14.11.1

Replace the executable/AppImage/container and keep the database and encryption key. No database change and no API change. The ▦ *Panels* button is no longer in the top bar while the Workboard is available: its list is now the *Panels* section of the **▤ Workboard** menu (a panel's ⋯ → *Panels…* opens it too).

## Upgrading from v14.11.0

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: a new table `connection_order` (your position of each connection in its folder); a previous version still starts on the database and lists connections by name again. New endpoint `POST /api/connections/order`; `GET /api/connections` returns your order and a `sort_order` per placed connection.

## Upgrading from v14.7.2

Replace the executable/AppImage/container and keep the database and encryption key; nothing is migrated. *Shares*, *Lock now* and *Sign out* are now in the account menu on your avatar (top bar, left of the AI button). Between 761 and 1024 px the right panel is a drawer opened with 🗂.

## Upgrading from v14.7.1

Replace the executable/container and keep the database and encryption key. The database gets two new columns on `share_links` (`allow_tunnels`, `allow_edit`, both off); an older binary still starts on it. Existing shares behave as before until you turn the switches on.

## Upgrading from v14.7.0

Replace the executable/container and keep the database and encryption key. No database change. New setting `hsts` (default on, as before). If you start a **new** installation and create the administrator from another computer, enter the setup code from WRM's console/log or from `first-admin-setup-code.txt` next to the database; existing installations are not affected.

## Upgrading from v14.6.3

Replace the executable/AppImage/container and keep the database and encryption key; nothing is migrated. Addons may declare the new capability `chains.read` (apiVersion 2, with a panel); the permission sheet shows it, and an update of an installed addon that adds it waits for a new approval. New API `GET /api/addons/{key}/instances/{id}/chains`, `GET …/chains/{chain}` and `POST …/chains/{chain}/check`; the panel context's `features` adds `chains`. See [Addons](../README.md#addons).

## Upgrading from v14.6.2

Replace the executable/AppImage/container and keep the database and encryption key; nothing is migrated. `GET /api/proxies` adds `tunnel_offers` (your dynamic tunnels not yet used by a proxy). Log discovery always skips folders named `venv*`, `.venv*`, `virtualenv*`, `site-packages`, `dist-packages`, `node_modules`, `__pycache__`, `.git`, `.tox`, `.mypy_cache` and `.pytest_cache`, in addition to the excluded names of *Log settings*.

## Upgrading from v14.6.1

Replace the executable/AppImage/container and keep the database and encryption key. The table `auth_links` and the `users` column `activation_pending` are added; v14.6.1 still starts on the database. New policies `password_reset_mail` (`all`; the Enterprise level requires `users`, so such an installation shows it as a check to fix in *Admin panel → Security*) and `registration_verify_email` (on). New public API `/api/auth/forgot`, `/api/auth/reset`, `/api/auth/activate`; `GET /api/auth/config` adds `password_reset` and `registration_email`, `POST /api/auth/register` takes `email`, `POST /api/admin/users` takes `email` and `invite`, `GET /api/admin/users` adds `activation_pending`. See [Forgot password, invitations and registration](../README.md#forgot-password-invitations-and-registration).

## Upgrading from v14.6.0

Replace the executable/AppImage/container and keep the database and encryption key. Two tables (`mail_settings`, `mail_outbox`) and the `users` columns `email`, `email_verified`, `email_pending`, `email_token`, `email_token_at` are added; v14.6.0 still starts on the database. Nothing is sent until an administrator sets up and turns on *Admin → E-mail*. New policies `mail_security_notices` and `user_email_self` (both on), new API `/api/admin/mail`, `/api/me/email`, `/api/auth/verify-email`; `GET /api/admin/users` adds `email`, `email_verified`, `email_pending` and `PUT /api/admin/users/{id}` takes `email`. See [E-mail](../README.md#e-mail).

## Upgrading from v14.5.1

Replace the executable/AppImage/container and keep the database and encryption key. Two tables are added (`server_heartbeat`, `cli_reset_requests`); v14.5.1 still starts on the database. New settings `cli_password_reset` (`alert`; the Business level requires `stopped`, Enterprise `confirm`, so an installation at those levels shows it as a check to fix in *Admin panel → Security*) and `alert_cli_reset` (on), new flag `-reset-code`, new API `/api/admin/cli-resets`. AI sessions can no longer run `wrm -reset-password`, `-reset-2fa` or `-lockdown`. See [Password reset on the server](../README.md#password-reset-on-the-server).

## Upgrading from v14.5.0

Replace the executable/AppImage/container and keep the database and encryption key. Nothing is migrated, no setting or API changes. Window sizes are kept in the browser's local storage (`wrm_win_sizes`), storage column and panel widths in `wrm_stg_cols`; ↺ on a window, or clearing the site data, returns to the defaults.

## Upgrading from v14.4.0

Replace the executable/AppImage/container and keep the database and encryption key. Nothing is migrated and no setting changes. Recognition uses the same fixed read-only discovery as the AI wizard (`POST /api/chains/discover`) and runs only when you add servers in the builder or press *Recognise services*; turn off *Recognise … right after adding* in the builder to add servers without it. *Test now* uses `POST /api/chains/test` with a one-check definition and saves nothing. The browser follows AI and MCP sessions over one stream per tab, `GET /api/ai/events`; the per-session stream stays for other clients.

## Upgrading from v14.3.2

Replace the executable/AppImage/container and keep the database and encryption key. Two columns are added (`storages.provider`, `chain_drafts.base_chain_id`); an older binary still starts on the database. Existing storages get their provider from the endpoint (Hetzner, DigitalOcean, AWS, Backblaze, Cloudflare R2, otherwise *Other*); change it in the storage form. AI connections do not get the new **`chains`** scope automatically. An administrator adds it to `ai_mcp_scopes` (*Admin panel → AI connections → allowed scopes*), and then a user enables it on an AI connection.

## Upgrading from v14.3.1

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Two new settings: `storage_check_minutes` (default 15) replaces the fixed 5-minute background bucket check (set it to `0` to send no background requests at all), and `chain_storage_check_s` (default 300) limits how often a chain's storage check asks the bucket.

## Upgrading from v14.3.0

Replace the executable/AppImage/container and keep the database and encryption key. No database or setting change: this is a bug-fix release for AI/MCP responsiveness and sensitive-path matching; existing connections, storages and chains are unchanged.

## Upgrading from v14.2.0

Replace the executable/AppImage/container and keep the database and encryption key. No database or setting change: the storage browser uses the same tables, and existing storages and chains are unchanged. **Read-only S3 share links no longer work**: the Share button, its API and `/s3/` were removed for security. Delete any link you shared from an older version; storage access now needs a signed-in account. The unused `storage_shares` table may remain and can be dropped.

## Upgrading from v14.1.0

Replace the executable/AppImage/container and keep the database and encryption key. No database or setting change: the guided builder and the discovery run in the existing Chains section and use the connections you already saved. Existing chains are unchanged. Discovery is read-only and happens only when you open the wizard or press *Look at the servers*; a discovered mount is shown as a check, not attached and not turned into a storage.

## Upgrading from v14.0.0

Replace the executable/AppImage/container and keep the database and encryption key. Three tables are added (`chains`, `chain_samples`, `chain_drafts`); an older binary still starts on the database. The **Chains** button appears in the top bar: turn the section and its background checks off with `chains_enabled=0`, and set the default interval with `chain_interval_s` (default 30 s, minimum 15 s). Nothing on your servers changes: chains use the SSH connections and storages you already saved, install nothing and run no commands. AI connection scopes do not change; the chain builder is in WRM's built-in assistant, not in MCP clients.

## Upgrading from v13.14.11

Replace the executable/AppImage/container and keep the database and encryption key. Two tables are added (`storages`, `storage_shares`); an older binary still starts on the database. The **Storages** button appears in the top bar (turn the section off with `storages_enabled=0`). AI connections do not get the new storage scopes: an installation that never changed `ai_mcp_scopes` now defaults to the five server scopes plus `addons.manage`; add `storage_read`, `storage_write` or `storage_delete` there, and then to an AI connection, to let it reach a storage; writes and deletes by an AI need the automatic mode. Upgrading from v13.14.10 or older also brings the v13.14.11 changes (share links with a per-share cookie, credentials bound to their destination; see the [changelog](../CHANGELOG.md)).

## Upgrading from v13.14.9

Replace the executable/AppImage/container and keep the database and encryption key. No schema change. On Windows, a self-update that failed with *put the new binary in place: … being used by another process* has to be done once by hand: stop WRM, replace the `.exe` (keep its file name if a shortcut or service points to it), delete leftover `.wrm-update-*.exe` files and start WRM again. From v13.14.10 on, the self-update waits for the lock and copies the file if the lock stays.

## Upgrading from v13.14.8

Replace the executable/AppImage/container and keep the database and encryption key. No database change. In the managed Docker runtime of the original deploy tool, an empty `DEPLOY_HOSTNAME` now becomes the WRM computer's name at the next start of the tool; a value set in *Original runtime settings* is kept. To use the original tool's Git section, set `DEPLOY_CONFIG_REPO` there (see [Git workspace: Medium original deploy tool](../README.md#git-workspace-medium-original-deploy-tool)).

## Upgrading from v13.14.0

Replace the executable/AppImage/container and keep the database and encryption key. No new tables or columns. On the first start, when `addon_instances` is empty, WRM copies the rows of the newest non-empty `addon_instances_old_<time>` table (left there by the start-up bug of v13.9.0–v13.14.0) back with their ids and logs *Addon instances: restored …*; the old copies are kept and can be dropped by hand once the addons look right. WRM now creates the folder `<DB_PATH>.diag/` (goroutine dumps, newest 20 kept) on every stop and restart.

## Upgrading from v13.13.1

From v13.13.0, also read the next section. Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: a new table `file_edits` (append-only like the audit log, kept for `audit_retention_days`) and new columns `folders.auto_sort` and `ssh_keys.certificate`; a previous version still starts on the database. New settings `file_edit_history` (`1`) and `file_edit_max_kb` (`1024`). New endpoints `GET / POST /api/folders/sort` and `GET /api/admin/file-edits[/{id}]`; `PUT /api/keys/{id}` also takes `certificate`. Git checks no longer compare files the catalog does not track: their state is *outside the check* instead of *same* / *differs*, and log files (`*.log`, `*.log.*`, files below a `log` or `logs` folder) are never compared nor deployed unless the catalog's include list names them. A key file on the WRM server with a `<key>-cert.pub` next to it now logs in with that certificate first.

## Upgrading from v13.13.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change; a previous version still starts on the database. The WRM PRO Claude plugin is now 1.2.1 (icon path only). If adding the marketplace failed in the Claude app, add it again with the lowercase repository name `vedranius/web-browser-rdm-public`.

## Upgrading from v13.12.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change; a previous version still starts on the database. New endpoint `POST /api/mcp/setup-all` (audited as `mcp.setup_script`). The WRM PRO Claude plugin 1.2.0 no longer has the settings `wrm_url` and `wrm_token`: after updating the plugin, run `/wrm-pro:wrm` (or `/wrm-pro:connect`) once and approve the sign-in in the browser; drop `--config` from scripts that install it. A Claude Desktop extension downloaded from this version keeps running when its token is refused and offers `wrm_connect`.

## Upgrading from v13.11.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change; a previous version still starts on the database. New command-line options `-addon-lint`, `-addon-lint-json`, `-addon-new`, `-addon-id`, `-addon-dev`, `-addon-dev-port` and `-addon-dev-private`, and new endpoints under `/api/addons/dev/`. The addon ids `dev`, `install`, `drafts`, `community` and `users` are now reserved; an installed addon with one of those ids keeps working but cannot be installed again under that id.

## Upgrading from v13.10.0

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: new tables `addon_http_grants`, `addon_secrets` and `addon_http_log`; a previous version still starts on the database (it ignores them, and addons with an `http` block were not installable there). New policy `addon_http` (`all`): who may use addons that reach web sites (`off`, `admins`, `all`). Existing addons are unchanged; an addon gets network access only after an update that declares `http` and a new approval.

## Upgrading from v13.9.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change; a previous version still starts on the database. Git workspace installations whose connection was deleted are no longer listed, and deleting a connection now removes its installations. *Show saved login* and *Copy password* follow the existing *Allow export of passwords & keys* policy (`allow_secret_export`, on by default; administrators always can) and the step-up policy.

## Upgrading from v13.8.0

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: new tables `addons_personal`, `addon_instances`, `addon_store` and `addon_windows` and new columns `files`, `approved` and `pending` in `addons`; a previous version still starts on the database (it ignores personal addons, instances and stored values). Installed addons keep working; their files are pinned by SHA-256 the next time they are installed. New policy `addon_personal` (`all`): set it to `admins` or `off` if users must not install their own addons.

## Upgrading from v13.7.0

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: a new column `saved` in `ai_sessions` and an index; a previous version still starts on the database. New settings `ai_mcp_later_approval` (`run`: an approval given in WRM after the AI app stopped waiting runs the call at once; `repeat` keeps the earlier behaviour, where the AI must repeat the call) and `ai_keep_saved_sessions` (`1`). *Deny* in an AI app's approval question is now a final denial; closing the question without an answer still keeps the request in WRM. New endpoints: `GET /api/ai/sessions?view=all`, `POST /api/ai/sessions/{id}/save`, `GET /api/ai/sessions/{id}/export`, and the key `ai` of `/api/me/ui/{key}`. A decision on a request that is no longer waiting answers `409` with `decided`.

## Upgrading from v13.6.0

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: new columns `locked_at` and `lock_failures` in `auth_sessions` and `lock_pin_hash` in `users`; a previous version still starts on the database. New settings `screen_lock_mode` (`server`), `screen_lock_max_minutes` (`0`), `screen_lock_pin` (`1`), `screen_lock_max_failures` (`0`), `appearance_default_theme` (`graphite`) and `appearance_custom_css` (`1`); the Business and Enterprise security levels require `screen_lock_mode=server` and an idle lock of at most 30 / 15 minutes, and Enterprise turns the PIN off. **Two behaviour changes:** a locked browser session can no longer reach the API until it is unlocked (set `screen_lock_mode` to `remember` or `browser` for the looser behaviour), and a colour theme from an addon is no longer applied to everyone: each user picks it in *Settings → Appearance*, or the administrator makes it the default (`appearance_default_theme=addon:<addon>:theme`). The accent colour chosen in *Settings → General* before v13.7 is carried over into the theme once.

## Upgrading from v13.5.0

Replace the executable/AppImage/container and keep the database and encryption key. No migration: the log settings gain `main_folders` (empty) and `guess_operators` (on) inside the existing per-server JSON, and a previous version ignores them. Labels of logs under a configured main folder are taken from it before the default patterns.

## Upgrading from v13.4.4

Replace the executable/AppImage/container and keep the database and encryption key. Additive migration: a new table `log_settings` (per connection, encrypted) and new columns `title`, `purpose`, `resumed_from` and `file_path` in `ai_sessions`; a previous version still starts on the database. Log folders and journal units configured in the Zabbix dialog are read into the new *Log settings* until you save them there. New settings `ai_show_thinking` and `ai_auto_titles` (both `1`). Dates are shown in the EU format by default; pick another in *Settings → Date and time format* (per browser).

## Upgrading from v12.9.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change. New policies: `medium_default_engine` (`native`) and `medium_build_user` (empty). If the original deploy tool was already imported, WRM sets `medium_default_engine=original` once at the first start and logs it, so projects keep their engine; set it to `native` in *Admin → Policies* to deploy with WRM by default. A project with an engine chosen explicitly keeps it. Pulling deploy-config from Git needs `DEPLOY_CONFIG_REPO` in the Medium settings and a WRM Git source for that repository (with a write token for commit and push).

## Upgrading from v12.8.1

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Editing deploy-config in WRM uses the existing `medium_native` policy (default `admins`); users without it see the files read-only.

## Upgrading from v12.7.1

Replace the executable/AppImage/container and keep the database and encryption key. No database change. The Git check reads the version from the code by default; `VERSION.md` is only compared when the new checkbox is on. Deployments and updates still write `VERSION.md`. The first check after the upgrade may need a target refresh to load the older version trees; until then the per-file history is used.

## Upgrading from v12.7.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Existing `deploy-config` files keep working; an environment may now add `servers: [name, ...]` to limit its projects to those servers. Importing a config ZIP still replaces the config directory and removes its Git history, so use Git for changes once the directory is a clone.

## Upgrading from v12.6.26

Replace the executable/AppImage/container and keep the database and encryption key. The database gets the table `medium_project_engine`; a previous version still starts on it (projects switched to the native engine are then deployed by the original tool again, which reads the same server state). Every project stays on the Docker runtime until someone allowed by the new policy `medium_native` (default `admins`) switches it to *WRM native*. New settings `medium_command_timeout` (300) and `medium_transfer_timeout` (3600). With step-up on, native deployments, rollbacks and lock releases ask for a confirmation (`medium.deploy`, `medium.rollback`, `medium.lock_release`).

## Upgrading from v12.6.25

Replace the executable/AppImage/container and keep the database and encryption key. The database gets the table `siem_destinations`; a previous version still starts on it (it does not enforce a lockdown, so end a lockdown before going back). Nothing is exported until an administrator adds a SIEM destination. New settings `alert_failed_signins` (10), `alert_failed_signins_minutes` (10), `alert_break_glass`, `alert_new_admin`, `alert_secret_reveal`, `alert_ai_full_access` and `alert_audit_off` (all on); new notification event `security.alert`; new command line flags `-lockdown` / `-lockdown-reason` and environment variable `WRM_LOCKDOWN`. Security events (`security.*`) are now recorded even with `audit_enabled` off.

## Upgrading from v12.6.24

Replace the executable/AppImage/container and keep the database and encryption key. The database gets the tables `ldap_directories`, `ldap_identities`, `scim_clients`, `scim_users`, `scim_groups` and `scim_group_members`; a previous version still starts on it (accounts that sign in with a directory password then cannot sign in, since WRM keeps no password for them). Nothing changes until an administrator adds an LDAP directory or a SCIM client. New setting `ldap_sync_minutes` (15). Changing an OIDC provider's issuer now asks for its client secret again.

## Upgrading from v12.6.23

Replace the executable/AppImage/container and keep the database and encryption key. The database gets the tables `oidc_providers` and `user_identities` and the columns `users.sso_only`, `auth_sessions.identity_id` and `team_members.source`; a previous version still starts on it (directory accounts then cannot sign in, since they have no password). Nothing changes until an administrator adds an identity provider. New settings `local_login` (`all`), `oidc_revalidate_minutes` (60) and `oidc_unverified_hours` (12). The Enterprise level now requires `local_login` = `break_glass`; it takes effect only while a directory sign-in is enabled, so nobody is locked out of an installation without one. Releases from this version on are signed and published automatically when a version is merged.

## Upgrading from v12.6.22

Replace the executable/AppImage/container and keep the database and encryption key. The database gets the tables `teams`, `team_members` and `share_teams` and the columns `users.role`, `users.break_glass` and `auth_sessions.step_up_at`; a previous version still starts on it (it sees administrators as before, and auditors and operators as ordinary users). Existing accounts keep their rights: administrators stay administrators, everyone else is a *User*. New settings `require_step_up` (off; the Business and Enterprise levels require it, so an installation at one of those levels shows it in the security posture until it is turned on), `step_up_minutes` and `break_glass_ips`. Scripts that call `PUT /api/admin/settings`, change roles or delete users with step-up on get `403 {"step_up": true}` unless the session confirmed with `POST /api/auth/step-up` (signing in counts as a confirmation for `step_up_minutes`). `DELETE /api/admin/users/{id}` still deletes everything the user owns unless `?transfer_to=<user id>` is given.

## Upgrading from v12.6.21

Replace the executable/AppImage/container and keep the database and encryption key. The database gets one new table, `security_deviations`; a previous version still starts on it. Existing installations stay at the **Home** level and no setting changes until an administrator picks another level in *Admin panel → Security* (or sets `WRM_INITIAL_SECURITY_LEVEL` / `WRM_SECURITY_LEVEL`). New settings `security_level`, `require_https` and `link_share_max_days`; new endpoints `GET /api/admin/security` and `POST /api/admin/security/level`. Loosening a policy below the chosen level now needs a reason (`_deviation_notes` in `PUT /api/admin/settings`); scripts that set settings at the Home level are not affected.

## Upgrading from v12.6.20

Replace the executable/AppImage/container and keep the database and encryption key. No database change. New setting `update_signing_keys` (empty: only the built-in key). This version has no built-in release key yet, so its self-update still verifies the SHA-256 only and the update dialog says *Signature: not checked*; the version that adds the key turns the signature check on. New command-line options `-release-keygen`, `-release-sign`, `-release-verify`.

## Upgrading from v12.6.19

Replace the executable/AppImage/container and keep the database and encryption key. No database or setting change. Building from source now needs **Go 1.26 or newer** (`go.mod` asks for the Go 1.27.2 toolchain, which `go` downloads by itself when `GOTOOLCHAIN` allows it).

## Upgrading from v12.6.11

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Nothing changes until an administrator sets `ai_opencode_sandbox` (`admins` or `all`) and presses **Download and install** again in *Admin panel → AI assistant → OpenCode agents*, which prepares the Docker sandbox. New settings `ai_opencode_sandbox` and `ai_opencode_sandbox_image`; `POST /api/ai/sessions` takes `sandbox: true`.

## Upgrading from v12.6.10

Replace the executable/AppImage/container and keep the database and encryption key. The database gets one new column, `users.ai_delegate_training` (a user's own choice to allow helper models that may train on prompts); a previous version still starts on it. Nothing changes until an administrator adds an **OpenCode Zen** provider, turns on delegation (`ai_delegate`) or installs OpenCode in *Admin panel → AI assistant → OpenCode agents*. `ai_provider_kinds` now includes `opencode`; if you set it yourself, add `opencode` to allow the new provider. New settings `ai_delegate*` and `ai_opencode*`; new endpoints `/api/ai/delegate` and `/api/admin/opencode`. OpenCode itself is stored next to the database in `<DB>.opencode/`.

## Upgrading from v12.6.9

Replace the executable/AppImage/container and keep the database and encryption key. No database change. The Workboard is now the canvas itself: your v12.6.9 tiles open as addon windows beside your terminals the first time (from the top right), and tiles you pulled out into panels stay where they are. *Broadcast* moved from the top bar into **🧰 Tools**. `/api/me/ui/{key}` now accepts values up to 256 KiB.

## Upgrading from v12.6.8

Replace the executable/AppImage/container and keep the database and encryption key. The database gets one new table, `user_events` (each user's in-app notifications, newest 200); a previous version still starts on it and ignores it. Nothing changes on the page until you open **▤ Workboard** and add an addon. New policy `workboard_ssh_metrics` (on): turn it off if the *Server load* addon must not log in to servers over SSH; servers with Monitoring still show their values. New endpoints `/api/me/events` and `/api/workboard/metrics`.

## Upgrading from v12.6.7

Replace the executable/AppImage/container and keep the database and encryption key. The database gets one new table, `user_ui` (your panel layout); a previous version still starts on it and ignores it. Nothing moves: your first layout is built from what the browser remembered in v12.6.7 (the right panel's tab and width, where the AI live view was docked, floating or hidden). New endpoint `/api/me/ui/{key}`.

## Upgrading from v12.6.6

Replace the executable/AppImage/container and keep the database and encryption key. No database or API change. The AI live view now opens as a **🤖 AI** tab in the right panel instead of a box at the bottom left; pick *Floating* in its title bar (or the 🤖 menu) for the old behaviour. Medium's Docker runtime now trusts the WRM host's root certificates (policy `medium_host_ca`, default on); if runtime preparation failed with *certificate verify failed*, press *Prepare runtime* again.

## Upgrading from v12.6.5

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Medium now uses the managed Docker runtime by default on Linux and macOS too. An installation whose host runtime (Python virtual environment) is already prepared keeps it; set `WRM_MEDIUM_RUNTIME=docker` to move it to Docker, or `native` to keep the host runtime for good. Nothing is installed until an administrator presses **Install Docker** in Medium step 2; the new policy `medium_docker_install` (default on) switches that button and the automatic installation off.

## Upgrading from v12.6.4

Replace the executable/AppImage/container and keep the database and encryption key. One additive table (`medium_activity`); a previous binary still starts on the database. `POST /api/git/medium/start` now answers `202` and starts in the background (poll `GET /api/git/medium`), and a deploy-config or SSH ZIP uploaded while Medium is busy or not prepared yet is kept and applied at the next start (`202`, `pending`). Command output is shown in the Medium activity log to the user who ran it and, for the shared runtime, to every Medium user; set `medium_log_output` to `admins` or `none` to limit that.

## Upgrading from v12.6.3

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Blocked commands keep asking under full access until an administrator sets `ai_full_access_blocked` to `deny` (*Admin panel → AI connections → Blocked commands with full access*).

## Upgrading from v12.6.2

Replace the executable/AppImage/container and keep the database and encryption key. One additive column (`users.ai_full_access`); a previous binary still starts on the database. If `ai_full_access` was already on, it now applies to every AI session (`ai_full_access_mode` = `all`, the default); set the mode to *AI connections with the Full access scope* (`connection`) in *Admin panel → AI connections* to keep the v12.6.2 behaviour, or to *Users who turn it on themselves* (`user`).

## Upgrading from v12.6.1

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Full access stays off until an administrator turns on `ai_full_access` (*Admin panel → AI connections*); then edit an AI connection and tick **⚠ Full access**. Open sessions of that connection end and the app reopens them with the new rules.

## Upgrading from v12.6.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change. Existing AI connections keep their limits. To let an app work without limits, edit its AI connection in *Settings → AI connections*, press **Allow everything** (it also ticks **No time or action limit**), confirm and save; open sessions end and the app opens them again with the new rules. Administrators can forbid connections without limits with `ai_auto_until_grant` = off.

## Upgrading from v12.5.1

Replace the executable/AppImage/container and keep the database, encryption key and Medium directories. One additive table (`medium_server_links`) is created at start; a previous binary still starts on the database. The native preview uses your existing WRM Git sources and SSH connections and writes nothing to servers. Medium deployments keep running through the original tool.

## Upgrading from v12.5.0

Replace the executable/AppImage/container and keep the database and encryption key. No database change. AI connections (MCP) stay off by default; once an administrator turns them on (*Admin panel → AI connections*), **Claude Desktop** users pick *Claude Desktop* in *Settings → AI connections*, download the extension and install it. An old `claude_desktop_config.json` entry with the stdio bridge keeps working. Token lifetimes can now be chosen in **hours or days**. The defaults change from 8 hours (at most 24) to **7 days (at most 30 days)**: `ai_mcp_token_hours` 168, `ai_mcp_max_token_hours` 720. A value an administrator saved explicitly stays in force; lower them again in *Admin panel → AI connections* if your policy needs short lifetimes.

## Upgrading from v12.4.5

Replace the executable/AppImage/container and keep the database, encryption key and private Medium directories/volumes. No database migration. Medium settings saved empty now mean "original default" and are no longer passed to the tool; earlier `.previous` copies of replaced SSH profiles are deleted at the next Medium start. Stamp no longer refuses silently: choose **Also stamp installations that are not up to date** to stamp an installation as it is.

## Upgrading from v12.4.4

Replace the executable/AppImage/container, preserving your database, encryption key and private runtime directories/volumes. No database migration is needed. In Medium use the combined ZIP fields and **Prepare and run**; empty fields reuse existing imports. Setup failure retains approved archives for retry. Existing original-tool deployment history and behavior are retained.

In Monitoring select **Logs** and a server. Discovery uses the same SSH account without sudo; it finds known service `log`/`logs` directories, bounded configured service roots and Supervisor log directories. Manual roots remain editable from the log toolbar. Date/timezone filters and keyword search read historical segments across matching logs, rather than only their recent tail. Progress and bounds remain visible. Zabbix is needed for host metrics, not SSH logs.

Windows first setup can require Docker/WSL license, OS permissions, virtualization setup or a restart on the WRM host. A Windows service cannot display a desktop installer: install Docker once for its account, then resume. Docker in Windows-container mode is not switched automatically because other containers may be using it. See the [guide](V12.4.5_GUIDE.md).

## Upgrading from v12.4.3b

Replace the executable/container with v12.4.4 and retain the DB, encryption key and private runtime directories/volumes. Migrations are additive. Existing explicit policy values remain in force; MCP is still off by default. Functionality moves to **Monitoring → Functionality**. Configure each server's agent transport/log roots; automatic metrics polling is opt-in.

Client GitLab cache requires a client token and browser CORS permission or the companion on the computer that can access GitLab. Native Medium's Git/build subprocesses still need their direct route. [v12.4.4 setup, limits and troubleshooting](V12.4.4_GUIDE.md).

## Upgrading from v12.4.3

Download the v12.4.3b Windows executable, stop WRM and replace the previous executable while retaining its database and encryption key. **The v12.4.3 updater cannot parse the `b` suffix, so this first hotfix needs a manual download.** v12.4.3b adds letter-hotfix support for subsequent checks/updates.

For Windows 10/11, install/start [Docker Desktop with Linux containers](https://docs.docker.com/desktop/setup/install/windows-install/) under the account running WRM. In **Git → Medium**, import the original tool ZIP or, if it is already imported, click **Prepare runtime / retry**. WRM pulls the pinned Linux image and installs the original lockfile in a Linux volume. Import your deploy-config and SSH profile and start Medium. Standard `~/.ssh/key` references work inside its private Linux home. No Windows Python/rsync installation is needed. Browser clients need no deploy tools.

Docker Desktop is not supported on Windows Server: run WRM/Medium on a supported Linux host or Linux VM/container and access it from the Windows browser. A remote Docker daemon is not supported by this bind-mount/localhost bridge. Docker engine access must belong to the account running WRM; installing Docker for a different desktop account does not grant a Windows service access automatically.

Back up the database/key and private `.medium` source/config/preferences **plus** Docker volumes named `wrm-medium-<instance>-runtime` and `wrm-medium-<instance>-u<user>`. Stop Medium before backing up user volumes. Runtime source/venv are read-only; user volumes retain native edits/history across restarts. Re-importing config/SSH replaces that profile explicitly without deleting deploy history. Existing native Linux runtime data is copied into a new user volume on its first Docker start. Source and SSH keys remain private.

Linux keeps its native runtime when dependencies exist; use `WRM_MEDIUM_RUNTIME=docker` to select the managed runtime or `native` to select native explicitly. Windows always selects Docker. Custom Docker paths/images use `WRM_MEDIUM_DOCKER` and `WRM_MEDIUM_CONTAINER_IMAGE`. The Windows CLI is also discovered at standard Docker Desktop installation paths without requiring a PATH refresh. Additional build toolchains/users may require a trusted custom Linux image; the pinned image uses its unprivileged `wrm` account for builds. Unix host-agent forwarding remains a native-runtime option; managed Docker uses the private SSH profile.

## Upgrading from v12.4.0

Replace the binary/AppImage or update the container. Additive migrations create `git_config_profiles` and add `folders.sort_order`. Existing folders retain alphabetical order until reordered. Back up the database, matching encryption key and adjacent `.medium` / `.bitwarden` directories; they contain private runtime/configuration data. Older binaries ignore the new fields, but cannot use the new features.

Medium is enabled for administrators by default. Import the original deploy-tool ZIP once as an administrator, then each operator imports their deploy-config and optional SSH profile. The Docker image includes Python 3.12+, git, OpenSSH, rsync and the official Bitwarden CLI. Standalone binaries keep these dependencies optional: install them on the WRM host before enabling the corresponding feature. Medium uses the original engine's own plans, locks, backups and history; it does not import those into Simple/Advanced runs.

The earlier client GitLab relay/replay and AI/navigation/color/icon requests are delivered in v12.4.4; see the [operation guide](V12.4.4_GUIDE.md).

## Upgrading from v12.3.6

Replace the binary or use the in-app updater; update the AppImage/container as appropriate. Startup adds a nullable folder parent and the `git_server_health` cache table. Existing folders, catalogs and route defaults keep their settings; previous binaries ignore the new fields.

The new **Functionality** view is available from Git Simple and Advanced. Run a Git check to populate it, and enable periodic monitoring/notification subscriptions if desired. Runtime commands read service state with the SSH login's permissions; no services are started or restarted. For AI availability, open **Settings → AI assistant**; provider connectivity and effective account permission are shown separately.

## Upgrading from v12.2.8

Replace the binary. One additive table (`git_hostnames`) is created automatically; v12.2.8 still starts on the same database. No configuration changes.

- **Git Simple mode is automatic.** Services and environments the checker found on your servers appear as ready cards without defining environments per service; **Save detected environments** stores them in the catalog (audited). Advanced is unchanged and your own definitions keep priority.
- **Servers are shown by hostname** (read once over the existing SSH login and cached, falling back to the connection host) instead of the connection name.

## Upgrading from v12.2.7

Replace the binary. No database or configuration changes; v12.2.7 still starts on the same database.

- The **version label in the top bar** ("WRM v12.2.8") now opens the **Updates** dialog, so an administrator can always press **Check now**, even when no update is known and the green badge is hidden. Non-administrators see the versions only.

## Upgrading from v12.2.6

Replace the binary. One additive table (`ai_ollama_runs`) is created automatically; v12.2.6 still starts on the same database. No configuration changes.

- The **AI provider wizard** has presets for xAI, DeepSeek, Together AI, Fireworks AI, Cerebras, Cohere and Perplexity, refreshed model suggestions and a **Load the model list** button (the provider's `/models`). Models can always be typed freely.
- **Ollama, step by step:** detection, status, installed models, per-OS install / background-service / download commands and, optionally, **running them through WRM**. The last is off by default: an administrator enables the policy `ai_ollama_manage` (`WRM_AI_OLLAMA_MANAGE=1`). Each run needs a confirmation, is permission-checked on the server and written to the audit log (`admin.ai_ollama_*`). See the configuration table and the API reference below.

## Upgrading from v12.2.5

Replace the binary. No database or configuration changes; v12.2.5 still starts on the same database.

- **Settings → AI connections** now has a step-by-step **wizard** for connecting Codex CLI, Claude Code, Gemini CLI, Cursor / VS Code, Claude and ChatGPT: it creates the token, shows (and lets you download) the ready command or config file and tests the connection. The AI assistant has a matching **provider wizard**. See *Connect an AI app* below.
- The MCP server explains what "wrm" is to its clients, and permission denials now say which scope is missing and how to fix it.

## Upgrading from v12.2.4

Replace the binary. No database or configuration changes; v12.2.4 still starts on the same database.

- The **AI assistant window** inside a terminal window can now be docked left / right / top / bottom, floated, moved, resized and expanded over the whole terminal window. See *AI assistant* below.

## Upgrading from v12.2.3

Replace the binary. No database or configuration changes; v12.2.3 still starts on the same database.

- The **session recording replay** window can now be resized (drag its bottom-right corner) and expanded to **fullscreen** (⛶ button, **Esc** leaves). The terminal is scaled to the window width, so there is no sideways scrolling. See *Session recording* below.

## Upgrading from v12.2.2

Replace the binary. No database or configuration changes; v12.2.2 still starts on the same database.

- Every tab of the Advanced Git workspace, of Settings and of the Admin panel now has a short explanation of what it is for. See [Git workspace: tab explanations](../README.md#git-workspace-tab-explanations).

## Upgrading from v12.2.1

Replace the binary. No database or configuration changes; v12.2.1 still starts on the same database.

- The Git workspace has a **Simple / Advanced** switch in its header. **Simple** shows one guided card per service environment (what runs on the server vs. the latest version on the branch, a clear status) with **Install / Update / Deploy again** and **Roll back to the previous version**. **Advanced** is the unchanged Git workspace and stays the default. See [Git workspace: Simple mode](../README.md#git-workspace-simple-mode).

## Upgrading from v12.2.0

Replace the binary. No database or configuration changes.

- New terminal windows open **maximized** with the cursor in the terminal; windows of a saved session keep their layout.
- Selecting more text than fits on the screen now copies the **whole selection** (it is copied when you release the mouse).
- A click or double-click on the title bar of a maximized window no longer shrinks it; only dragging restores it.

## Upgrading from v12.1.0
Replace the binary. There are no database changes; v12.1.0 still starts on the same database.

- **Gemini and other OpenAI-compatible providers** work with tools: the provider's `extra_content` (Gemini's `thought_signature`) of each tool call is stored in the session and sent back unchanged, which fixes *HTTP 400: Function call is missing a thought_signature*.
- **Temporary provider errors** no longer end the turn at once: HTTP 408 / 429 / 5xx (for example 503 *high demand*) are retried up to three times with a back-off, and a request that is too large for the provider (HTTP 413, or a token limit such as Groq's) shortens the long command output of the conversation and is retried. The assistant shows a short notice for both.

## Upgrading from v12.0.0
Replace the binary. The database gets new tables (`mcp_tokens`, `mcp_clients`) and the column `ai_sessions.mcp_token_id`; v12.0.0 still starts on it and ignores them.

- **AI desktop apps (MCP)** are **off by default**. An administrator turns them on in *Admin panel → AI connections* (`ai_mcp_enabled`); by default only administrators may create AI connections (`ai_mcp` = `admins`), with the modes `read_only,ask,auto` (`ai_mcp_modes`), a default lifetime of 8 hours (`ai_mcp_token_hours`) and at most 24 hours (`ai_mcp_max_token_hours`; since v12.5.1 the defaults are 7 and 30 days). Users then find *Settings → AI connections*. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).
- MCP clients reach WRM at `https://<your WRM>/mcp`. Hosted connectors (claude.ai, ChatGPT) connect **from the provider's cloud**: WRM must be reachable from the internet over HTTPS for them. Behind a reverse proxy set `ai_mcp_public_url` (the address the apps use) and pass `/mcp`, `/oauth/` and `/.well-known/` through. Desktop apps with the stdio bridge only need to reach WRM from your computer.
- The kill switch of the AI assistant (`ai_kill_switch`) also stops every AI connection. Turning off the built-in assistant (`ai_assistant`) does not turn off MCP, and the other way round.
- The top bar shows **🤖 AI connected** while an AI app is connected; at narrower widths the top-bar labels hide a little earlier than before so nothing overlaps.

## Upgrading from v11.7.0

Replace the binary. The database gets new tables (`ai_providers`, `ai_sessions`, `ai_messages`) and the column `users.ai_blocked`; v11.7.0 still starts on it and ignores them.

- **AI assistant** (off by default). An administrator turns it on in *Admin panel → AI assistant* (`ai_assistant` = `admins` or `all`) and adds the organisation's provider (Anthropic, OpenAI, Azure OpenAI, AWS Bedrock, Google Vertex AI or an OpenAI-compatible endpoint). Users then get an **AI** button in the title bar of SSH terminals and file managers. See [AI assistant](../README.md#ai-assistant).
- New modes are restrictive by default: `ai_modes` = `read_only,ask` (the automatic mode must be allowed explicitly), personal API keys are off (`ai_personal_keys`), output is redacted before it is sent to the provider (`ai_redact_output`).
- **Consumer chat subscriptions (Claude.ai, ChatGPT) cannot be used** by WRM or any other third-party application: use an API key, an enterprise gateway, or the MCP integration of v12.1.0 ([AI desktop apps](../README.md#ai-desktop-apps-mcp)).
- AI sessions appear in *Sessions & recordings* with an **✦ AI** mark (the transcript is recorded like a terminal) and in the audit log as `ai.*` events.
- The top bar no longer overlaps the *New window* button at about 1100 px.

## Upgrading from v11.6.1

Replace the binary (or use the new one-click update from the next release on); the database does not change.

- **Run as a service:** on an interactive start WRM asks once whether to install itself as a service that starts at boot (systemd, launchd, rc.d or a Windows service; *Don't ask again* is remembered in `<DB_PATH>.service.json`). Scripts use `-install-service`, `-uninstall-service` and `-service-status`; `-no-service-prompt` skips the question. See [Service](../README.md#service). `sudo wrm -setup-docker [-service-user ACCOUNT]` (Linux) installs Docker Engine for Medium, lets the WRM account use it and restarts the WRM service.
- **Updates:** WRM checks the project's GitHub releases once a day (no user data is sent). The version in the top bar shows **⬆ Update vX.Y.Z** with the release notes; administrators update with one click (SHA-256 checked against `SHA256SUMS.txt`, previous binary kept for a rollback). Air-gapped installations: `WRM_UPDATE_CHECK=off`. Docker never replaces itself. See [Updates & self-update](../README.md#updates--self-update).
- **📂 in an SSH terminal** is now a menu: *Open here*, *Open /*, *Open home*. "Here" works without configuring OSC 7. See [Windows, tabs & snapping](../README.md#windows-tabs--snapping).
- A systemd unit that runs WRM with `ProtectSystem=strict` needs the binary's folder in `ReadWritePaths` for the one-click update (units written by `-install-service` have it).

## Upgrading from v11.6.0

Replace the binary; the database does not change.

- **📂 in the title bar of SSH terminals** opens the file manager of the same server, in the terminal's current directory when the shell reports it (OSC 7) ([#24](https://github.com/vedranius/web-browser-RDM-public/issues/24)). See [Windows, tabs & snapping](../README.md#windows-tabs--snapping).

## Upgrading from v11.5.0

Replace the binary; the database does not change.

- **Activity tab** in the Git workspace with three views — *Git activity* (pushes and commits per branch and author from the Git server), *Deploy history* (the `history.jsonl` of every environment destination merged with the WRM runs of installations without environments) and *What is where* (branch, commit, version, who and when per environment and server, commits behind the branch head, drift). Each view has a **CSV export** of what its filters show. See [Git workspace: activity, deploy history & what is where](../README.md#git-workspace-activity-deploy-history--what-is-where).
- **Gitea sources:** a Git source can now be a Gitea server (API v1, `Authorization: token …`). Targets, comparisons, deploys and the activity work as with GitHub; a .gitignore pull request is not available for Gitea yet (download the file).
- The activity reads the provider API with the source's **read-only token**; for GitHub repository activity (pushes, force pushes) the token needs read access to the repository's metadata, otherwise only commits are shown. The deploy history and the matrix read the destinations over SSH like the *Environments* tab and need the `git_checks` policy; without it only WRM's own records are shown.
- No new settings or policies.

## Upgrading from v11.4.1

Replace the binary. The database gets one new table (`git_env_exclusions`); the previous binary still starts on it and ignores it. Catalogs and installations without environments work exactly as before.

- **Environments (optional):** a service of the catalog can get named environments — ordered destinations `server:/path` (the server is the name of an SSH connection) with rules: typed confirmation, `allowed_refs`, extra ignore patterns, an optional `post_deploy` command and the number of backups to keep. Add them on the new *Environments* tab (✎) or in the catalog JSON (`apps[].environments`). See [Git workspace: environments, deploy plans & safe deploys](../README.md#git-workspace-environments-deploy-plans--safe-deploys).
- **Server-side state:** an environment deploy keeps its state next to the installation, by default in `<path>/.deploy-bak` (the existing backup directory): `state.json`, `history.jsonl`, `deploys/<id>.json`, `backups/rollback-<id>/` and the lock `.lock/`. Backups keep the existing `<id>-<YYYYMMDD-HHMMSS>` layout, and `VERSION.md` / `updates.jsonl` are written as before, so file-based tools keep working.
- **Anchored patterns:** an exclude or ignore pattern that starts with `/` now matches the whole relative path only (`/lib/a.py` does not match `x/lib/a.py`). *Exclude permanently* writes such patterns; existing patterns are unchanged.
- The same policies apply: `git_update` (deploy), `git_upgrade` (deploy of another ref), `git_rollback`, `git_restart`. Environment deploys run right away (no scheduling); use a dry run to preview.

## Upgrading from v11.4.0

Replace the binary; the database does not change.

- **Clipboard history per user:** the *Clipboard* panel now keeps its history per signed-in user, and signing out removes it from the browser ([#30](https://github.com/vedranius/web-browser-RDM-public/issues/30)). The browser-wide history of earlier versions is deleted on the first sign-in, so the panel starts empty once.

## Upgrading from v11.3.0

Replace the binary. The database gets two new columns of `git_installs` (`last_ok_at`, `last_ok_state`); the previous binary still starts on it and ignores them.

- **Git sources lost by the v11.3.0 upgrade come back:** upgrading a database from v11.2.0 or earlier to v11.3.0 set the `git_sources` table aside as `git_sources_old_<time>`. On the first start, v11.4.0 restores those sources (with their encrypted tokens) when no source exists yet; the old table is renamed to `git_sources_restored_<time>`.
- **Every file of an installation:** the details list every file of the server folder and of the repository, with the reason for each state, and diff any repository file (also the ones the catalog excludes). **Refresh the targets once** so the targets keep their full commit and the repository files outside the catalog are listed too. See [Git workspace: every file and partial checks](../README.md#git-workspace-every-file-and-partial-checks).
- **Partial checks:** select installations, servers or services and check only those (or only stale ones); results stream in and can be cancelled. Two new settings: `git_check_parallel` (8) and `git_check_per_server` (2).

## Upgrading from v11.2.0

Replace the binary. The database gets two new tables (`git_feeds`, `git_pipelines`) and four new columns of `git_sources` (`write_token`, `hook_id`, `hook_secret`, `hook_at`); the previous binary still starts on it and ignores them.

- **.gitignore helper:** a new *.gitignore* tab in the Git workspace. The output is a download; a merge / pull request needs an optional **write token** on the Git source and the new policy `git_gitignore_mr` (**administrators only** by default). See [Git workspace: .gitignore helper](../README.md#git-workspace-gitignore-helper).
- **Optional CI integration:** a webhook per Git source, bundles fetched from CI artifacts, and the deploy method *CI pipeline* per service. All of it is off until you set it up; WRM works without any CI. See [Git workspace: CI integration (optional)](../README.md#git-workspace-ci-integration-optional).
- Webhooks arrive at `/api/hooks/git/<id>`: if WRM is behind a reverse proxy that only forwards some paths, forward this one too.
- Checks now also record the size of files found only on servers; installations show it after their next check.

## Upgrading from v11.1.0

Replace the binary. The database is unchanged: installs and transfers are stored as runs in `git_runs`; the previous binary still starts on it.

- **New servers from the Git workspace:** *Install…* on the *Installations* tab (a wizard: server, services with their version, target directory, pre-checks, per-host files from templates / copies / typed, optional systemd or supervisor unit, optional start) and *Transfer…* for an installation (copy it to another server with its per-host configuration). See [Git workspace: new servers (install & transfer)](../README.md#git-workspace-new-servers-install--transfer).
- New policies `git_install` and `git_transfer` — **administrators only** by default; change them in *Admin panel → Policies*.
- Installs and transfers send the existing run notifications (*Git: … started / succeeded / failed*) and are audited as `git.install` / `git.transfer`.

## Upgrading from v11.0.0

Replace the binary. The database gets two new tables (`git_runs`, `git_pending_restarts`); the previous binary still starts on it and ignores them.

- **Deploying from the Git workspace:** *Update*, *Upgrade…*, *Rollback…*, *Restart…* and *Stamp* on the *Installations* tab and in the details of an installation, a new *Runs* tab, and scheduled runs. See [Git workspace: update, upgrade, rollback & restarts](../README.md#git-workspace-update-upgrade-rollback--restarts).
- New policies `git_update`, `git_upgrade`, `git_rollback`, `git_restart` — **administrators only** by default; change them in *Admin panel → Policies* — and `git_schedule_grace_minutes` (30).
- New notification events for runs and restarts (*Git: an update, upgrade or rollback started / succeeded / failed*, *a service restart is scheduled / done*, *a scheduled run starts soon*) in *Settings → Notifications*.
- WRM writes on servers only when you start or schedule a run; checks stay read-only.

## Upgrading from v10.10.0

Replace the binary. The database gets new tables (`git_sources`, `git_trees`, `git_blobs`, `git_workspace`, `git_targets`, `git_installs`); the previous binary still starts on it and ignores them.

- **Git workspace:** the new **⎇ Git** button in the top bar compares the services installed on your servers with GitLab / GitHub or an offline bundle, file by file and read-only. See [Git workspace](../README.md#git-workspace-compare-services).
- It is on for everyone. Administrators can hide it (`git_enabled` off, or `WRM_GIT_ENABLED=0`) or limit checks to administrators (`git_checks`).
- New notification events *Git: a new version is available*, *Git: files changed by hand on a server* and *Git: a server is unreachable for checks* in *Settings → Notifications*.

## Upgrading from v10.9.1

Replace the binary. The database gets a new table `bookmarks`; the previous binary still starts on it and ignores it.

- **Folder bookmarks:** ★ next to the path bar of the file manager and in a terminal's title bar. *Bookmark this directory…* saves the current folder for one connection, a folder, a tag or all connections; paths may use `~`, `$USER`, `{host}` and `{name}`. A bookmark can be the **start directory** of a connection. See [Folder bookmarks](../README.md#folder-bookmarks).
- **WinSCP bookmarks:** *★ → Manage bookmarks… → Import WinSCP.ini* adds the remote directory bookmarks of a WinSCP configuration as global bookmarks.
- Bookmarks are part of *Export* / *Import*.

## Upgrading from v10.9.0

Replace the binary. The database gets a new table `putty_sessions` and a column `connections.putty_session`; the previous binary still starts on it and ignores them.

- **Import from PuTTY:** *Settings → Data → Import from PuTTY* reads a `.reg` export of your saved PuTTY sessions, with their proxies. See [Import from mRemoteNG, PuTTY and OpenSSH](../README.md#import-from-mremoteng-putty-and-openssh).
- mRemoteNG connections imported **before** v10.9.1 do not remember their *PuttySession*, and importing the file again skips them as existing: give them the proxy by hand (or through their folder's default proxy).

## Upgrading from v10.8.1

Replace the binary. The database gets new tables (`proxies`, `proxy_grants`, `notify_channels`, `notify_subscriptions`, `notify_prefs`, `notify_pending`, `notify_state`) and columns (`connections.proxy_id`, `folders.jump_conn_id`, `folders.proxy_id`); the previous binary still starts on it and ignores them (connections then go without their proxy and folder defaults).

- **Proxies:** define SOCKS5, SOCKS4/4a or HTTP CONNECT proxies (or a SOCKS tunnel of a connection) under *🔑 SSH keys & credentials → 🌐 Proxies* and choose one in a connection (*Proxy*, next to *Jump host*). See [Proxies](../README.md#proxies-socks--http).
- **Folder defaults:** right-click a folder → *Folder settings…* to give it a default jump host and proxy. Connections in the folder that have no jump host of their own then use the folder's (existing connections without a jump host included). Choose *no jump host (not the folder's)* where a connection must stay direct.
- **Notifications:** an administrator adds channels in *Admin panel → Notifications* (e-mail, Telegram, Slack / Mattermost / Rocket.Chat, Teams, Discord, ntfy, Gotify, Pushover, webhook); users subscribe in *Settings → Notifications*. Nothing is sent until both happened. See [Notifications](../README.md#notifications-e-mail-chat-push-webhook).
- **Docker:** `docker-compose.yml` now maps `host.docker.internal` to the Docker host (`extra_hosts`), for proxies that run on the host.
- Deleting a jump host now makes the connections that used it fall back to their folder's default jump host (or go direct).

## Upgrading from v10.8.0

Replace the binary. Nothing changes in the database.

- **Larger imports work:** mRemoteNG, OpenSSH config and inventory (CSV / Excel) files up to **20 MB** import again (from about 1 MB on they failed with *Bad JSON*). A larger file gets a clear message with its size and the limit, already in the browser.
- **Install as an app** (*Settings → General*) now shows the steps for your browser, or the reason why installing is not possible here (plain `http://` to an address other than `localhost`, an untrusted certificate, Firefox on the desktop). The *Install* button appears only when the browser can install right away. To get HTTPS quickly, see [HTTPS with Caddy](../README.md#https-with-caddy).
- A reverse proxy in front of WRM must allow request bodies of at least 32 MB for imports (`client_max_body_size 0;` in the nginx example does).

## Upgrading from v10.7

Replace the binary. The database gets two new columns (`connections.notes`, `connections.temp_until`); the previous binary still starts on it (it shows quick connections as normal ones).

- **⚡ Quick connect** field at the top of the sidebar: `user@host:port` + Enter.
- **📝 Notes** in the connection dialog and in the right-click menu; a 📝 next to a connection opens them.
- **🧰 Tools** in the top bar, and *Network tools from this server…* / *Check ports of this host* in the right-click menu (policy `network_tools`, default: all users).
- **Install as an app**: *Settings → General*, or the install icon of the browser's address bar (needs HTTPS with a trusted certificate, or `localhost`).

## Upgrading from v10.6

Replace the binary. The database gets one new column (`connections.bmc`); the previous binary still starts on it.

- Connections can have a **BMC** (*Edit → Out-of-band management*); right-click → *Power & status (BMC)…* and *Serial console*. For IPMI install `ipmitool` on the WRM server or on the jump host (`apt install ipmitool`).
- New protocol **Serial port** for console cables on the WRM machine (administrators by default, policy `serial_ports`).
- Saved workspace sessions now really keep remote desktop windows (and console windows).

## Upgrading from v10.5

Replace the binary. The database gets three new columns (`connections.tags`, `ext_id`, `ext_tags`) and a table `inventory_sources`; the previous binary still starts on it (it ignores the tags).

- Connections have **tags**; the sidebar shows a tag filter and an environment badge (`env:prod` …). See [Tags & inventory import](../README.md#tags--inventory-import-csv-excel-netbox).
- **Settings → Data → Import inventory**: CSV / Excel files and NetBox.

## Upgrading from v10.4

Replace the binary. The database gets new tables (`ssh_keys`, `ssh_key_deployments`, `credentials`, `credential_grants`) and two new columns (`connections.key_id`, `connections.credential_id`); the previous binary still starts on it (connections that use the new login types then cannot log in).

- New **🔑 Keys** button: your SSH key store and the credentials vault. See [SSH keys & credentials vault](../README.md#ssh-keys--credentials-vault).
- Connections can log in with an **SSH key from the key store** or a **credential from the vault**; remote desktops and FTP can use vault credentials too.
- Deleting an account now also deletes its snippets, tunnels, keys and credentials.

## Upgrading from v10.3

Replace the binary. The database gets one new column (`connections.options`); the previous binary still starts on it.

- **Remote desktops need guacd** next to WRM: `apt install guacd` (Debian/Ubuntu), `docker run -d --network host guacamole/guacd`, or the `guacd` service that `docker-compose.yml` now contains. Without guacd everything else keeps working; opening an RDP/VNC/Telnet connection then explains what is missing. See [Remote desktop](../README.md#remote-desktop-rdp-vnc-telnet).
- Importing from **mRemoteNG** now also brings RDP, VNC and Telnet connections (with domain, console, colors and RD Gateway).
- SQLite now waits up to 10 s for a busy database on every connection (fixes rare *database is locked* errors under load).

## Upgrading from v10.2

Replace the binary. The database gets a new table (`snippets`) and one new column (`connections.monitor`); the previous binary still starts on it.

- **File names no longer carry a suffix:** releases are `v10.3.0`, binaries `wrm-pro-v10.3.0-linux-amd64` etc. Update scripts and systemd units that use the file name.
- **Live status is on by default:** WRM opens a TCP connection to every saved host:port once a minute (SSH and FTP: it reads the greeting and closes politely, like Nagios `check_ssh`). Servers see it in their logs. Change the interval or turn it off in *Admin panel → Security policies → Live status*, or untick *Monitor up/down status* in a connection. See [Live up/down status](../README.md#live-updown-status) about fail2ban.
- **Broadcast input** is available to everybody; turn it off with `broadcast_enabled` if you do not want it. Every start and stop is in the audit log.

## Upgrading from v10.1

Replace the binary. The database gets a new table (`connection_tunnels`) and two new columns on `connections` (`jump_conn_id`, `web_path`); nothing else changes. To go back, run the v10.1 binary on the same database; it does not know these features (connections behind a jump host then try to connect directly, web interface connections cannot be opened, tunnels do not run).

- New: [jump hosts](../README.md#jump-hosts-bastions), [SSH tunnels](../README.md#ssh-tunnels-port-forwarding), [web interface connections](../README.md#web-interfaces-http--https-connections) and [import from mRemoteNG / OpenSSH](../README.md#import-from-mremoteng-putty-and-openssh).
- **Tunnels are on by default for all users**, listening on `127.0.0.1` of the WRM machine only. Remote forwards (`-R`) and listening on network addresses are limited to administrators. Review *Admin panel → Security policies → SSH tunnels* — e.g. turn tunnels off (`WRM_TUNNELS_ENABLED=0`) or limit them to administrators if WRM runs on a shared server.
- Jump hosts need `AllowTcpForwarding yes` on the bastion (the OpenSSH default), exactly like `ssh -J`.

## Upgrading from v10.0

Replace the binary. The database is extended automatically (new tables, columns and triggers only — nothing is changed or removed; the previous binary still runs on the upgraded database).

- **Terminal sessions are now recorded by default** (`session_recording`). Users see *“This session is recorded”* and a **● REC** badge on the terminal. Inform your users (in many countries this is required), or turn it off in *Admin panel → Security policies → Session recording*. Keystrokes are **not** recorded unless you turn on *Record keystrokes*.
- Recordings use disk space next to the database (`recordings/`, typically a few KB to a few MB per session; limit per session `recording_max_mb`, kept `recording_retention_days` = 90 days). Back up or exclude that folder as you need.
- The audit log is now **append-only and hash-chained**. *Admin panel → Audit log → Verify integrity* checks it. Entries written before the upgrade have no hash and are skipped by the check.
- Every upload, download (also files opened in the editor and files in ZIP downloads) and server-to-server copy is logged in *Admin panel → File transfers* with size and SHA-256.
- New: `/healthz` for load balancers and monitoring, `wrm -healthcheck` for container health checks, a `Dockerfile` and `docker-compose.yml`.

## Upgrading from v9

Replace the binary and start it with the same database and the same `ENCRYPTION_KEY` (if you set one). Everything is migrated automatically. Things to know:

- **Stored secrets are re-encrypted.** Without `ENCRYPTION_KEY`, v9 used a public default key. v10 generates a random key file (`remote_manager.db.key`) and re-encrypts all secrets with it on the first start. Back that file up. (If you did set `ENCRYPTION_KEY`, nothing changes.)
- **Self-registration is closed by default.** Existing accounts keep working. Open it again in *Admin panel → Security policies* if you want it.
- **Passwords and keys are no longer sent to the browser.** The *Edit connection* dialog shows “saved and encrypted — leave empty to keep it”. *Export* is without secrets; *Export with passwords & keys* asks for your password.
- **SSH host keys are verified** (trust on first use). The first connection to each server remembers its key; a changed key is refused until you accept it.
- **Key file / Auto (~/.ssh)** authentication (keys stored on the WRM server) is limited to administrators unless the policy *Allow server key files for all users* is on.
- **Existing shares keep their behaviour** (anyone with the link, role *Operator*). Edit them to use the new access modes and roles.
- **Behind a reverse proxy:** state-changing API calls now require the `Origin` to match the `Host` (like WebSockets already did). Pass the `Host` header, or set `WRM_ALLOWED_ORIGINS`. Set `WRM_TRUST_PROXY=1` so rate limits and the audit log see real client IPs.
- The per-IP sign-in limit is now 20 failures / 10 minutes; accounts additionally lock for 15 minutes after 10 wrong passwords (policy).
- If the database comes from an earlier build that already had tables named like the new v10 tables (`audit_log`, `collab_messages`, …) but with a different layout, WRM keeps them as `<table>_old_<time>` and creates new ones (logged as *Upgraded table …*). Nothing is deleted.
- The app is served at `/`. Old bookmarks to `/static/` are forwarded there.
