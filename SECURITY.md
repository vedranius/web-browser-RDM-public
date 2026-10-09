# Security policy

## Supported versions

Security fixes go into the **latest release** only. Please upgrade before you report a problem, and include the version you use (`wrm -version`, `/api/version`, or the top bar).

| Version | Supported |
|---|---|
| v11.x | ✔ |
| v10.x and older | ✘ (upgrade: see *Upgrading* in the [README](README.md#table-of-contents)) |

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

**In scope:** the WRM server, web UI and container image in this repository. For example: authentication and 2FA bypass, privilege escalation between users or share roles, access to stored secrets or to other people's session recordings, changing or deleting audit records without detection, cross-site scripting or request forgery, reaching servers or networks you should not reach (the TURN relay, SSH tunnels and jump hosts included — e.g. using another user's jump host or tunnel, or bypassing the tunnel policies), and denial of service with small effort.

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
- [ ] Run WRM as an unprivileged user. The systemd example in the README adds `NoNewPrivileges`, `ProtectSystem=strict` and `PrivateTmp`. The Docker image runs as an unprivileged user; provide the encryption key as a Docker secret (`ENCRYPTION_KEY_FILE`) instead of keeping it only in the data volume.

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
- [ ] Use **read-only** GitLab / GitHub tokens for the Git workspace (e.g. GitLab `read_api` + `read_repository`, a GitHub fine-grained token with *Contents: read*). They are stored encrypted and never returned, but anyone with the WRM account can use them to read the repositories. Git checks only read files on servers; limit them with `git_checks` or hide the workspace with `git_enabled`.
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

**Monitoring**

- [ ] Forward the server log (`AUDIT …` lines) to your SIEM, or export the audit log and file transfers regularly (*Admin panel → Audit log / File transfers → CSV*). Set `audit_retention_days` to match your retention rules.
- [ ] Monitor `GET /healthz` (HTTP 200 while the server and database work).
- [ ] Watch for `auth.login_failed`, `auth.account_locked`, `hostkey.mismatch`, `admin.*`, `share.*`, `tunnel.*`, `desktop.*`, `terminal.broadcast`, `snippet.*`, `bookmark.*`, `ssh_key.*` (deploy, revoke, export) `credential.*` (grants, reveal, rotation), `proxy.*` (created, shared, tested), `folder.updated` (default jump host / proxy), `admin.notify_channel_*` and `bmc.*` (power actions) events.
- [ ] Keep WRM up to date. Releases are published on the [Releases page](https://github.com/vedranius/web-browser-RDM-public/releases); verify downloads with `SHA256SUMS.txt`.

## How WRM protects your data

The [Security model](README.md#security-model) section of the README describes the protections in detail: accounts, sessions, encryption of stored secrets, web security headers, host key verification, jump hosts and tunnels, share roles, collaboration, voice and TURN relay, and the audit log.
