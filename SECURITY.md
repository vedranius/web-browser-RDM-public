# Security policy

## Supported versions

Security fixes go into the **latest release** only. Please upgrade before you report a problem, and include the version you use (`wrm -version`, `/api/version`, or the top bar).

| Version | Supported |
|---|---|
| v12.x | ✔ |
| v11.x and older | ✘ (upgrade: see *Upgrading* in the [README](README.md#table-of-contents)) |

## Reporting a vulnerability

**Do not open a public issue for security problems.**

- Use GitHub's **private vulnerability reporting**: *Security → Report a vulnerability* in this repository.
- If that is not available, open an issue that says only that you have a security report and asks for a private contact. Do not include details.

Please include:

- the affected version and platform,
- what an attacker can do, and what they need first (no account, guest link, signed-in user, share member, admin…),
- steps to reproduce or a proof of concept,
- if you have them, suggestions for a fix.

You will get an answer as soon as possible, usually within a week. Once a fix is released, the release notes credit you if you want. Please give a reasonable time to release a fix before you publish anything.

**In scope:** the WRM server, web UI and container image in this repository. For example: authentication and 2FA bypass, privilege escalation between users or share roles, access to stored secrets or to other people's session recordings, changing or deleting audit records without detection, cross-site scripting or request forgery, reaching servers or networks you should not reach (the TURN relay, SSH tunnels and jump hosts included — e.g. using another user's jump host or tunnel, or bypassing the tunnel policies), getting the **AI assistant** to run something its mode does not allow (escaping the read-only classifier, the approvals, the automatic-mode limits or the destructive-pattern list, or changing the mode from server output), reading another user's AI transcripts, provider keys reaching the browser, and denial of service with small effort.

**Out of scope:** problems that need an administrator account or access to the server's files, missing HTTPS when the operator did not configure it, social engineering, reports from automated scanners that come without a working attack, and vulnerabilities in the servers you connect to.

## Hardening checklist

WRM is secure by default in most respects. For production and company use, also check the items below. The admin panel (*Overview*) warns about several of them.

**Transport and exposure**

- [ ] Serve WRM over **HTTPS** only: a reverse proxy with a real certificate, `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE`, or at least `HTTPS_SELF_SIGNED=1`. Without HTTPS, passwords and terminal traffic between browser and WRM are not encrypted, and browsers block the microphone.
- [ ] Behind a reverse proxy, bind WRM to localhost (`LISTEN_ADDR=127.0.0.1:8080`), pass the `Host` header through, and set `WRM_TRUST_PROXY=1`. Set `WRM_TRUST_PROXY` **only** when clients cannot reach WRM directly.
- [ ] Do not set `WRM_ALLOW_ANY_ORIGIN=1`.
- [ ] Put WRM behind a VPN or an IP allow-list when only internal users need it.
- [ ] Firewall: expose only the HTTPS port. Open the TURN ports (3478 UDP/TCP and the relay range) only if you use voice calls across networks. Otherwise set `turn_enabled=0`, or set `voice_enabled=0` to turn off voice entirely.

**Keys and data**

- [ ] **Back up the encryption key** (`<DB_PATH>.key`, or your `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`) **separately** from database backups. Someone who has both can decrypt the stored secrets.
- [ ] Keep the database, the key file and the self-signed certificate readable only by the WRM service user. WRM creates them with mode `0600`.
- [ ] Run WRM as an unprivileged user. `-install-service` writes a systemd unit with `User=`, `NoNewPrivileges`, `ProtectSystem=strict` (the data and binary folders writable), `ProtectHome=read-only` and `PrivateTmp`; the service's settings file `<data dir>/wrm.env` is mode `0600` (it may hold `ENCRYPTION_KEY`). The Docker image runs as an unprivileged user; provide the encryption key as a Docker secret (`ENCRYPTION_KEY_FILE`) instead of keeping it only in the data volume.

**Updates**

- [ ] The one-click update installs only a release file whose **SHA-256 matches the release's `SHA256SUMS.txt`** and which answers `-version` with the release version; it is refused otherwise. Both come from the same GitHub release, so the check protects against broken or swapped downloads, not against a compromised release repository: keep `update_repo` at the official repository (or your own mirror) and `self_update` at `admins` (or `off` with `WRM_SELF_UPDATE=off`). Every step is audited (`system.update_*`).
- [ ] The update check sends no user data (only the repository path and a fixed User-Agent). Turn it off in air-gapped networks with `WRM_UPDATE_CHECK=off`.
- [ ] Only administrators may write to the binary's folder: a newer `wrm-pro-v…` file copied there is started by the service on its next start (after a `-version` check).

**Accounts**

- [ ] Keep **self-registration closed** (the default) and create accounts in *Admin panel → Users*.
- [ ] Set **`require_2fa`** to `admins` or `all`.
- [ ] Use a sensible `password_min_length` (12 or more), and the lockout (`login_max_failures`) and session lifetime policies (`session_idle_hours`, `session_max_days`) your company requires.
- [ ] Disable accounts of people who leave (*Disable* signs them out everywhere at once) and review *Admin panel → Users* regularly.

**Connections and sharing**

- [ ] Keep `host_key_policy` at `tofu` or `strict`. **Never** use `off` in production.
- [ ] Keep `allow_server_keys=0`, so only administrators can use key files stored on the WRM server.
- [ ] Set `allow_secret_export=0` if users must not export their stored passwords and keys.
- [ ] Set `allow_link_shares=0` if guests without an account must not join shares.
- [ ] Prefer **Members only** shares with the lowest role that is enough (*Viewer* to look at files, *Observer* for calls and watching). Give shares an expiry date.

**Terminal features and live status**

- [ ] Decide about **broadcast input** (`broadcast_enabled`). It only types into the user's own open terminals, but one command reaches many servers; dangerous commands ask first. Every start and stop is in the audit log (`terminal.broadcast`).
- [ ] Review **shared snippets** (administrators publish them for everybody) like any runbook. Run-on-connect snippets come from the connection's owner and are audited (`terminal.auto_run`).
- [ ] **Folder bookmarks** marked as shared are visible (read-only) to everybody who uses your connections through a share: do not put secrets in their names, paths or notes.
- [ ] **Live status** opens a TCP connection to every monitored host:port once per interval (no login). Tell the teams who watch those servers, add the WRM server to fail2ban `ignoreip` where *aggressive* filters are used, and keep `status_jump_checks` off unless logins to the jump hosts every round are acceptable.

**SSH keys & credentials vault**

- [ ] Give shared credentials a **host list** (`Only for hosts`): a password is sent to the server at login, so without one a colleague could send it to a server of their own. With a host list, the people it is shared with can only use matching hosts and jump hosts.
- [ ] Prefer keys (deploy them with WRM) over shared passwords, and **revoke** keys of people who leave — *Who has access* shows what is left on a server.
- [ ] Keep `allow_secret_export` off if users should not be able to export private keys or show vault passwords (administrators always can; every export is audited).
- [ ] Use **read-only** GitLab / GitHub / Gitea tokens for the Git workspace (e.g. GitLab `read_api` + `read_repository`, a GitHub fine-grained token with *Contents: read* (and *Metadata: read* for the repository activity), a Gitea token with `read:repository`). They are stored encrypted and never returned, but anyone with the WRM account can use them to read the repositories. Git checks only read files on servers; limit them with `git_checks` or hide the workspace with `git_enabled`.
- [ ] Anyone allowed Git checks can **open any file** of an installation folder in the Git workspace (the SSH login's read rights apply; protected files are shown from the server, never from Git; symlinks are not followed and paths cannot leave the installation). Each opened file is audited (`git.file_viewed`); limit who may do this with `git_checks`.
- [ ] Git **runs** (update, upgrade, rollback, restart) write on servers with the SSH login of the connection: keep the policies `git_update`, `git_upgrade`, `git_rollback` and `git_restart` at *admins* (the default) unless others need them. Production installations (tag `env:prod`) need the typed server name; every run is audited per installation (`git.update`, `git.upgrade`, `git.rollback`, `git.restart`, `git.stamp`). Restarts use `sudo -n` only when the SSH user is allowed it without a password — grant just the `systemctl restart` / `supervisorctl` commands it needs.
- [ ] Git **installs and transfers** create directories, write per-host configuration and (opt-in) unit files on the target: keep `git_install` and `git_transfer` at *admins* (the default). Typed per-host values are never stored with a run; vault values are read on the WRM side and only for servers the credential's host list allows. Overwriting an existing installation and production targets need the typed server name; replaced files are backed up first. Audit `git.install` / `git.transfer`.
- [ ] Git **environment deploys** take a **lock** on each destination: `mkdir <state>/.lock` (atomic) with an owner file holding the deploy id, the WRM user, the WRM host name and the server's time — no secrets. A deploy removes only a lock whose owner file names its own deploy id. A **stale** lock (older than max(30 minutes, the transfer timeout) by the server's clock) is removed only after an explicit confirmation, by one remote command that checks the owner id and the age again (audited as `git.lock_removed`); remove it only when you know no deploy is running there. The state directory (`<path>/.deploy-bak` by default) holds backups of replaced files — keep it out of public web roots (the access doctor warns about paths that look like one).
- [ ] The remote commands of environment deploys are fixed POSIX scripts; file names travel on **stdin**, never in the command line, and every path is shell-quoted. A rollback never acts through a symlink on the server and moves files aside instead of deleting them. **`post_deploy`** is a command from the catalog that runs in the app directory with the SSH login's rights — only when it is ticked for a run (and shown before it runs); treat editing a service's environments like editing a deploy script. Dry runs take no lock and write nothing.
- [ ] The Git **activity, deploy history and "what is where"** views only read: the provider API with the source's read token (the token stays on the server; the browser gets commit metadata and web links), and `history.jsonl` / `state.json` / `VERSION.md` over SSH under the `git_checks` policy. Their **CSV exports** contain author names, servers and paths — treat them like the audit log; values a spreadsheet would run as a formula (`=`, `+`, `-`, `@`, tab, CR at the start) are prefixed with `'`.
- [ ] Git **write tokens** (for .gitignore merge requests only) are optional: give them only the scope needed to push a branch and open a merge / pull request, and keep `git_gitignore_mr` at *admins*. WRM writes to a Git server only on an explicit, confirmed request (audit `git.gitignore_mr`).
- [ ] Git **webhooks** (`/api/hooks/git/<id>`) take no session: the hook ID and its secret authenticate them (GitHub HMAC, GitLab token, or the generic token / timestamped HMAC). Prefer HMAC signatures, keep the secret in the CI's secret store, renew it when it leaks (*New secret*), and watch `git.webhook_rejected` in the audit log. Webhooks only start checks (read-only on servers).
- [ ] **CI pipeline** tokens (Jenkins API token, GitLab trigger token) and artifact feed headers are stored encrypted and sent only to their configured host (no redirects to another host). A CI-deployed update goes through the same policies and confirmations as a normal one.
- [ ] Use a **read-only** NetBox token for inventory imports; a remembered token is stored encrypted, but anyone with the WRM account can use it to read NetBox.
- [ ] Set a rotation reminder on shared passwords and rotate them when a colleague leaves; check the audit log for `credential.rotation_incomplete`.

**Out-of-band management & serial ports**

- [ ] Use BMC accounts with the least role that does the job (*Operator* for power actions) and a vault credential with a host list for them.
- [ ] Keep `serial_ports` at `admins` (or `off`): serial ports belong to the WRM machine.
- [ ] Watch `bmc.power` events; power actions and consoles are limited to the connection's owner.

**Quick connect, notes & network tools**

- [ ] Set `network_tools` to `admins` (or `off`) if users should not probe networks from the WRM server or from their servers; every run is audited as `nettool.run` (tool, target, source).
- [ ] Notes are stored in clear text in the database (like connection names): keep passwords in the vault, not in notes.
- [ ] Quick connections are audited as `connection.quick` and removed 24 hours after their last use; they use the same host key checks as saved connections.

**Remote desktop (RDP / VNC / Telnet)**

- [ ] Run guacd only on the WRM machine (or in WRM's network namespace) and keep it bound to `127.0.0.1`: guacd has no authentication of its own, so anyone who reaches port 4822 can open connections through it.
- [ ] Prefer **NLA** or **TLS** security for RDP and turn off *Accept the server certificate* where servers have proper certificates.
- [ ] Turn off the remote clipboard (connection option) on servers where copying data out must not be possible; sessions are recorded like terminals.
- [ ] Turn remote desktops off (`desktop_enabled=0`) if nobody needs them.

**SSH tunnels and jump hosts**

- [ ] Decide who may use tunnels: `tunnel_users=admins`, or `tunnels_enabled=0` if nobody needs them (`WRM_TUNNELS_ENABLED=0`). Tunnels reach whatever the SSH servers reach; they are as powerful as the SSH accounts behind them.
- [ ] Keep `tunnel_bind_any=0` and `tunnel_remote_forward=admins` (the defaults). A tunnel on `0.0.0.0` or a network address is **not protected by WRM sign-in**: anyone who reaches that port uses the tunnel. If an administrator opens one, restrict it with the host firewall.
- [ ] When WRM runs on a shared server, remember that `127.0.0.1` tunnel ports are reachable by **every local user and process** of that server. Run WRM on a dedicated host or container, or limit tunnels to administrators.
- [ ] Set `tunnel_idle_minutes` so forgotten manual tunnels close themselves, and review running tunnels in *Admin panel → Overview*.
- [ ] On the bastions, allow only the forwarding you need (`AllowTcpForwarding`, `PermitOpen`, `GatewayPorts no`), as you would for OpenSSH clients.
- [ ] Watch for `tunnel.start`, `tunnel.error`, `tunnel.configured` and `web.open` events in the audit log.

**Proxies**

- [ ] Decide who may define proxies (`proxies=admins` or `off`); everybody may use proxies shared with them. A proxy reaches whatever its network reaches, like a jump host.
- [ ] SOCKS and HTTP proxy logins travel in plain text to the proxy. Prefer proxies on trusted networks, reach them through a jump host (the hop is encrypted), and give shared proxies their own accounts. WRM never sends a shared proxy's password through a grantee's own jump hosts.
- [ ] A "WRM SOCKS tunnel" proxy uses its owner's SSH login and cannot be shared. Review running tunnels (*Admin panel → Overview*) and the `tunnel.*` and `proxy.*` events.
- [ ] In Docker, a proxy on the host (`host.docker.internal`) should listen only on the Docker bridge, not on every interface.
- [ ] WRM refuses a proxy address on its own machine that is a running SSH tunnel of another user; tunnels on `0.0.0.0` (administrators, `tunnel_bind_any`) are still reachable from the network, see *SSH tunnels* below.

**Notifications**

- [ ] Channel secrets (bot tokens, webhook URLs, SMTP passwords) are encrypted at rest and never returned; treat a leaked webhook URL like a password and replace it.
- [ ] Use STARTTLS or TLS for e-mail; WRM does not send an SMTP password over a plain connection to another machine.
- [ ] Notification texts contain connection names, hosts and error messages: send them only to channels the recipients may read. With *Users may enter their own recipient* users can send their own notifications to any address of that kind.
- [ ] Sign webhooks (`X-WRM-Signature`, HMAC-SHA256) and check the signature on the receiving side. Watch `admin.notify_channel_*` events.

**Audit trail and session recording**

- [ ] Keep the audit log on (`audit_enabled`, the default). Run *Admin panel → Audit log → 🔏 Verify integrity* regularly, and after any incident. A broken chain means that someone changed the database directly.
- [ ] Decide about **session recording** (on by default). Tell your users that terminal sessions are recorded (in many countries this is required), and set `recording_retention_days` to match your retention and privacy rules.
- [ ] Turn on keystroke recording (`session_recording_input`) only if you need it. Typing at password prompts is masked, but other secrets typed on the command line would be recorded.
- [ ] Recordings may contain sensitive output (configuration files, logs). Keep the recordings folder (`WRM_RECORDINGS_DIR`, mode `0700`, files `0600`) on an encrypted disk, include it in protected backups, and give administrator rights only to people who may see them. Every view and download of a recording is audited.

**AI assistant**

- [ ] Keep the assistant off (`ai_assistant=off`, the default) until you have decided who may use it, with which provider, and where its data goes. Start with `admins`.
- [ ] Keep `ai_modes` at `read_only,ask` (the default). Allow `auto` only where you want it, and narrow it with **mode rules**: e.g. `{"match":"tag","value":"prod","modes":["read_only"]}` for production. The most restrictive matching rule wins.
- [ ] Keep `ai_redact_output` on. Redaction is pattern-based: it catches common secret formats, not every secret. Treat everything the assistant reads as sent to the provider.
- [ ] Choose providers deliberately (`ai_provider_kinds`, `ai_models`): an enterprise agreement (zero data retention, a regional endpoint, Bedrock or Vertex in your own cloud account, a gateway, or a local model) where server data must not leave your control. Keep `ai_personal_keys=off` if users must not send server data to their own accounts.
- [ ] Use dedicated API keys for WRM with spending limits at the provider; rotate them like any credential. They are encrypted at rest and never sent to the browser.
- [ ] Give the SSH logins the assistant uses the least privilege: it acts with the connection's login (and `sudo -n` only where that login may). Add organisation-specific dangerous commands to `ai_blocked_commands` and secret files to `ai_read_deny_paths`.
- [ ] Keep `ai_allow_power=off` unless reboots through the assistant are needed (then each one needs an approval).
- [ ] Watch `ai.*` audit events (`ai.tool_denied`, `ai.approved` with `edited`, `ai.auto_allowed`, `ai.session_killed`), set `ai_transcript_retention_days`, and know where the kill switch is (*Admin panel → AI assistant*, or `WRM_AI_KILL_SWITCH=1`).

**Monitoring**

- [ ] Forward the server log (`AUDIT …` lines) to your SIEM, or export the audit log and file transfers regularly (*Admin panel → Audit log / File transfers → CSV*). Set `audit_retention_days` to match your retention rules.
- [ ] Monitor `GET /healthz` (HTTP 200 while the server and database work).
- [ ] Watch for `auth.login_failed`, `auth.account_locked`, `hostkey.mismatch`, `admin.*`, `share.*`, `tunnel.*`, `desktop.*`, `terminal.broadcast`, `snippet.*`, `bookmark.*`, `ssh_key.*` (deploy, revoke, export) `credential.*` (grants, reveal, rotation), `proxy.*` (created, shared, tested), `folder.updated` (default jump host / proxy), `admin.notify_channel_*` and `bmc.*` (power actions) events.
- [ ] Keep WRM up to date. Releases are published on the [Releases page](https://github.com/vedranius/web-browser-RDM-public/releases); verify downloads with `SHA256SUMS.txt`.

## AI assistant: threat model

The built-in assistant (v12.0.0) lets a language model act on a server. WRM treats the model as an **untrusted, possibly manipulated actor** and puts the decisions on the server:

- **No shell for the model.** The model can only call WRM's fixed tools. Each call becomes an action (read a file, run a command, edit a file) that the **permission engine** evaluates on the server before anything runs: the session's mode, the read-only classifier, the automatic-mode allow / deny lists and limits, the always-blocked destructive patterns and the unreadable sensitive files. The browser only shows approvals and passes the user's decision back; the engine does not depend on the panel (the MCP integration of v12.1.0 uses the same engine).
- **Prompt injection from server output.** Logs, files and command output can contain text written by anyone (an attacker who can write a log line, a web request, a file in a shared directory). Tool results are given to the model as **untrusted data**, framed and labelled as such, and the system prompt tells the model never to follow instructions in them. This reduces but cannot eliminate manipulation, so it is not what WRM relies on: **only the user can change the mode** (through the authenticated API, never through a tool), there is no tool that changes permissions, unknown tools are refused, and every action is evaluated by the engine regardless of why the model asked for it. In *ask* mode a manipulated model can only *propose* a change, which the user sees exactly; in *auto* mode it can only do what the user's allow list permits, within the limits, and never anything on the destructive list.
- **Read-only means read-only.** The classifier is an allowlist: commands it does not know, and anything with redirections, command or process substitution, here-documents, subshells, background jobs or variables, count as changes. Wrappers (`sudo -n`, `timeout`, `xargs`) are unwrapped and the inner command is checked. Sensitive files (password hashes, private keys, cloud credentials, process environments) are refused also through wildcards and symlinks (`read_file` resolves the real path first). A file the login can read through other means (e.g. a custom script) is outside what any classifier can know: give the SSH login only the rights it needs.
- **Destructive patterns** are checked on the parsed command and again on the raw text with quotes removed, so a pattern hidden in `sh -c "…"` or `eval` is still caught; they are refused in every mode, also after an approval and after the user edits a command. The list is maintained in `ai_destructive.go` and every rule has a test.
- **Approvals.** A pending approval shows the exact command or a unified diff; only the session's user decides; an edited command or content is evaluated again; an unanswered approval is denied after `ai_approval_timeout_seconds`; a file is written only if it is unchanged since the diff was made. Every request, decision (approved, denied, edited, timed out, cancelled), who and when, is audited.
- **Data leaving to providers.** Everything the assistant reads — command output, file contents, the OS description, and the connection notes if the user opts in — is sent to the chosen provider. With `ai_redact_output` (default on) secrets are replaced before sending (the patterns of the audit redaction — password, secret, token, credential, private key, passphrase … — plus private key blocks, `Authorization` headers, passwords in URLs, AWS / GitHub / GitLab / Slack / OpenAI / Anthropic / Google keys, JWTs, password hashes). Redaction is best effort; administrators choose which providers and models are allowed, and personal keys are off by default. Stored WRM secrets (connection passwords, keys, vault entries) are never given to the model.
- **Keys.** Provider keys, AWS secrets and service account keys are encrypted at rest (`encryptValue`), never returned by the API (only *is set* flags), never written to the audit log, and scrubbed from provider error messages.
- **Kill switch.** A session can be stopped (the running request and command are cancelled) or ended by its user; administrators end one session, all sessions of a user, or all sessions, turn the assistant off per user, or turn on `ai_kill_switch`, which ends everything and blocks new sessions until it is turned off. Pending approvals of an ended session are cancelled.
- **Recording.** Transcripts are redacted, kept for `ai_transcript_retention_days` and visible to the session's user and administrators; every session is also recorded like a terminal session.

## How WRM protects your data

The [Security model](README.md#security-model) section of the README describes the protections in detail: accounts, sessions, encryption of stored secrets, web security headers, host key verification, jump hosts and tunnels, share roles, collaboration, voice and TURN relay, and the audit log.
