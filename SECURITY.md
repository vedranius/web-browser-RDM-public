# Security policy

## Supported versions

Security fixes go into the **latest release** only. Please upgrade before you report a problem, and include the version you use (`wrm -version`, `/api/version`, or the top bar).

| Version | Supported |
|---|---|
| v10.x | ✔ |
| v9.x and older | ✘ (upgrade: see *Upgrading from v9* in the [README](README.md#upgrading-from-v9)) |

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

**In scope:** the WRM server and web UI in this repository. For example: authentication and 2FA bypass, privilege escalation between users or share roles, access to stored secrets, cross-site scripting or request forgery, reaching servers or networks you should not reach (the TURN relay included), and denial of service with small effort.

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
- [ ] Run WRM as an unprivileged user. The systemd example in the README adds `NoNewPrivileges`, `ProtectSystem=strict` and `PrivateTmp`.

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

**Monitoring**

- [ ] Forward the server log (`AUDIT …` lines) to your SIEM, or export the audit log regularly (*Admin panel → Audit log → CSV*). Set `audit_retention_days` to match your retention rules.
- [ ] Watch for `auth.login_failed`, `auth.account_locked`, `hostkey.mismatch`, `admin.*` and `share.*` events.
- [ ] Keep WRM up to date. Releases are published on the [Releases page](https://github.com/vedranius/web-browser-RDM-public/releases); verify downloads with `SHA256SUMS.txt`.

## How WRM protects your data

The [Security model](README.md#security-model) section of the README describes the protections in detail: accounts, sessions, encryption of stored secrets, web security headers, host key verification, share roles, collaboration, voice and TURN relay, and the audit log.
