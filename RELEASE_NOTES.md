<p align="center"><img src="https://raw.githubusercontent.com/vedranius/web-browser-RDM-public/v11.4.0/docs/brand/png/lockup/wrm-lockup-on-dark-664w.png" alt="WRM PRO — Web Remote Manager" width="332"></p>

## Web Remote Manager PRO v11.4.0 — Git: every file of an installation, partial and incremental checks

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/vedranius)

Browser-based remote server management: SSH terminal with snippets, broadcast input and live status, SFTP/FTP/FTPS file manager, jump hosts, proxies and SSH tunnels, remote desktops (RDP/VNC/Telnet), out-of-band management and serial consoles, SSH key management and a credentials vault, tags and inventory import, quick connect, runbooks and network tools, notifications by e-mail, chat and push, a Git workspace that compares deployed services with GitLab / GitHub or offline bundles, updates them and installs them on new servers, sharing, real-time collaboration with voice calls, audit trail and session recording. One self-contained binary (or container) with an embedded web UI, installable as an app.

### 📂 Every file of an installation

- The details of an installation now list **every file in the server folder** (the catalog's ignored directories are honoured) and **every file of the repository** at the target ref — not just the files the catalog tracks.
- Each row shows the size, the modification time, the state and **why**: tracked, **excluded by the catalog (with the pattern)**, not in the include list, protected (with the pattern), only on the server, only in Git, a **symlink** (its target is shown, never followed out of the installation) or **unreadable**.
- Filter the list (all / differing / tracked only / not tracked) and search it.
- **Open any file** (view only, size-limited; binary files are detected and compared by hash).
- **Diff server ↔ Git for any file in the repository**, including files the catalog excludes such as `setup.py` — clearly labelled *informational, not updated*. Unified or side-by-side, line endings ignored like the hashing, usable on a phone.
- Files the catalog does not track **never change the state** of an installation.
- **Commits behind:** the commits between the server's version (the commit in `VERSION.md`) and the target, from GitLab or GitHub.

### ⚡ Partial and incremental checks

- **Filter the overview** by service, server / host, environment and state (review, needs update, up to date, error, not checked).
- **Pick what to check:** tick installations, whole rows (servers) or columns (services), then *Check selected*; or *Check visible*, *Check this cell*, or *Check only stale* (never checked, older than 1 hour / 24 hours / 7 days). *Check now* still checks everything.
- **Results stream in** per installation, with a spinner on each cell and a **Cancel** button. Two new settings bound the work: `git_check_parallel` (8) and `git_check_per_server` (2).
- Everything you did not check **keeps its result** with "checked X ago"; a failed check keeps the **last good result** next to the error.
- The same on the *Installations* tab.

### 🛠 Fixes

- Upgrading a database from v11.2.0 or earlier to v11.3.0 set the Git sources aside (`git_sources_old_<time>`) instead of adding the new columns. v11.4.0 adds the columns in place and **restores the set-aside sources** once, when the new table is still empty.

### ⬆️ Upgrading from v11.3.0

Replace the binary. The database gets two new columns on `git_installs` (`last_ok_at`, `last_ok_state`); the previous binary still starts on it. Git sources lost by the v11.3.0 upgrade come back on the first start. Refresh the targets once so the repository files outside the catalog are listed too.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v11.4.0/CHANGELOG.md) for details.

---

## Included from v11.3.0 — Git: .gitignore helper and optional CI integration

#### 🧹 .gitignore helper

- A new **.gitignore** tab in the Git workspace collects candidates per service: files that exist **only on servers** (from the latest checks, with size and the servers they were seen on), **standard patterns** for Python, Node, Go, Java, Docker and IDE files (pre-selected when the repository uses them) and the service's **protected** per-host globs.
- **Warnings** about committed secrets or per-host configuration (`.env`, `*.pem`, `*.key`, `*credentials*`, files matching the protected globs) with the advice to commit a `.example` template instead.
- Edit the service's catalog **exclude / protected** lists on the same screen, and **preview** the new `.gitignore` as a diff against the current one.
- The output is a **download**. A **merge request** (GitLab) or **pull request** (GitHub) on a new branch is created **only when you choose it** and type the repository name; it needs an optional, encrypted **write token** on the Git source and the new policy **`git_gitignore_mr`** (administrators by default). Audited as `git.gitignore_mr`.

#### 🔗 Optional CI integration

- **Webhook per Git source:** GitLab push / tag events, GitHub push / release events (HMAC-signed) or a POST from Jenkins or any CI start an **immediate check** of the affected services. Rate-limited, audited, replays refused.
- **Bundles from CI artifacts:** fetch the newest bundle from a GitLab job artifact, a Jenkins artifact or any HTTPS URL (with an auth header), by hand or on a schedule, and use it as the offline source.
- **Deploy method "CI pipeline"** for services built elsewhere: *Update* triggers a **Jenkins** job (`buildWithParameters`) or a **GitLab pipeline** (trigger token) with `server`, `install_path` and `version`; WRM follows the status and shows the result and a link in the run history. Same confirmation, policies and audit as a normal update.
- Everything is opt-in: WRM keeps working without any CI.

#### ⬆️ Upgrading from v11.2.0

Replace the binary. The database gets two new tables (`git_feeds`, `git_pipelines`) and new columns on `git_sources`; the previous binary still starts on it. Only administrators may open .gitignore merge requests until you change `git_gitignore_mr` in *Admin → Policies*. Behind a reverse proxy that filters paths, forward `/api/hooks/git/` for webhooks.

---

## Included from v11.2.0 — Git: new servers (install and transfer)

#### 🆕 Install a service on a new server

- **Install…** on the *Installations* tab of the Git workspace opens a wizard: pick **any SSH connection**, the **services** (each at its catalog version, another branch or tag, or from an imported **offline bundle**) and the **target directory** (suggested from the server roots).
- **Pre-checks first:** free disk space, a writable directory, `tar`, the `python3` / `node` versions when the service has such files, and an **existing installation** at the target — refused unless you confirm the overwrite, and then the replaced files are backed up to `.deploy-bak/`.
- **Per-host files done right:** protected files are never taken from the repository. Fill them **from a template** in the repository (`config.ini.example`, `.env.sample`, `settings.example.py`, …: its placeholders and `key = value` lines become form fields), **copy** them from another installation of the same service, or **type** them. Secret values can come **from the vault**; they are used on the server only and never stored with the run.
- WRM writes the files in **one verified transfer** (staged, checked by hash, then renamed into place), runs the usual checks, rolls back automatically on failure, and writes **`VERSION.md`** and the **`updates.jsonl`** line like an update does.
- **Optional unit:** a systemd unit or supervisor program from a small template (user, working directory, command) — opt-in, with a warning, never over an existing one. **Starting** the service afterwards is opt-in too, with a health check.

#### 🚚 Transfer an installation to another server

- **Transfer…** copies an installation from server A to server B **with its per-host configuration**, but **without** logs, `.deploy-bak`, caches (`__pycache__`, `*.pyc`), the catalog's ignored directories and backup-looking paths (`*_BKP`, `run.py.bak`, …). The wizard lists what stays behind.
- The files stream directly from A to B through WRM (no temporary copies on WRM's disk), are verified by hash, and existing files at the target are backed up first. `VERSION.md` gets the new host and user.

#### 🔐 Policies, audit and notifications

- New policies **`git_install`** and **`git_transfer`** (administrators by default, changeable in *Admin → Policies*). Every install and transfer is audited (`git.install`, `git.transfer`) and sends the usual run notifications (started, done, failed).

#### ⬆️ Upgrading from v11.1.0

Replace the binary. The database is unchanged (installs and transfers are stored as runs); the previous binary still starts on it. Only administrators may install and transfer until you change the new policies in *Admin → Policies*. Production targets and overwrites need the server name typed.

---

## Included from v11.1.0 — Git: update, upgrade, rollback, restarts

#### 🚀 Update, upgrade and rollback from the Git workspace

- **Update** an installation to the newest version of its ref, **upgrade** it to another branch or tag (with a summary of what changes), or **roll back** to a backup — from the *Installations* tab or the details of an installation, over WRM's SSH connections (jump hosts, proxies, the vault). Servers need nothing but POSIX tools.
- **You choose the files:** old and missing files are ticked; files changed by hand only when you tick them (with a warning); protected per-host files never.
- **Safe by default:** a fresh comparison, a **dangling-import check** for Python, a **backup** in `.deploy-bak/`, **atomic writes** that keep mode, owner and line endings (CRLF stays CRLF), checks (`py_compile`, `node --check`, `sh -n` or your own command) and an **automatic rollback** if anything fails.
- **Production needs two confirmations:** a summary, then you type the server name.
- **Rolling** over many servers, stopping at the first failure, with **live progress** per step.
- **Works with file-based deploy tools:** WRM writes `VERSION.md`, `.deploy-bak/<id>-<time>/` and `updates.jsonl` in their format. *Stamp* writes `VERSION.md` where an installation is already current.

#### 🔄 Service restarts (opt-in) and maintenance windows

- Restarts are **off by default**. Per run: only show what needs a restart, restart right after a successful update, at a chosen time, or after a delay — a warning names the units and servers. A **health check** (`systemctl is-active`, `supervisorctl status`) follows every restart.
- **Schedule** updates and restarts into a maintenance window. Scheduled runs survive a WRM restart, can be cancelled, are skipped if WRM was down too long, and **notify before and after**.
- The new **Runs** tab shows the history and what is scheduled.
- Administrators: policies `git_update`, `git_upgrade`, `git_rollback`, `git_restart` (administrators by default) and `git_schedule_grace_minutes`. Every update, upgrade, rollback and restart is audited.

#### ⬆️ Upgrading from v11.0.0

Replace the binary. The database gets two new tables (`git_runs`, `git_pending_restarts`); the previous binary still starts on it. Only administrators may update, upgrade, roll back and restart until you change the new policies in *Admin → Policies*. Subscribe to the new *Git* events in *Settings → Notifications* if you want them.

---

## Included from v11.0.0 — Git workspace: compare services

#### ⎇ Git workspace (phase 1: compare, read-only)

- **Which version runs where?** The new **⎇ Git** button opens a full-screen workspace that compares the services installed on your servers with **GitLab** or **GitHub** (cloud or self-hosted) — **file by file**, over WRM's existing SSH connections, so jump hosts, proxies and the vault apply. Servers need no git, Python or internet access.
- **Offline bundles:** when WRM cannot reach the Git server, import a `.tar.gz` bundle built elsewhere; WRM shows its age and source. WRM also **exports** the same format for servers that only a file-based deploy tool can reach.
- **Service catalog:** project, ref (`tag:latest`, `tag:v2`, a branch), subdirectory, include / exclude, **protected** per-host files (`config*`, `*.ini`, `.env`) and fingerprint files. WRM **suggests an entry from the repository tree**; catalogs import and export as JSON.
- **Discovery:** one SSH call per server finds the installations below `/opt`, `/srv`, … by their fingerprint files, skips copies and backups (`*_BKP`, `backup`, `.deploy-bak`, …) and labels environments from the path (`test/App` → *test*, `App_prod` → *prod*).
- **Comparison:** every file is *ok*, *old* (*2 behind* — it matches an earlier tag), *modified* (changed by hand), *missing*, *extra* or *protected*. Installations are *needs update*, *review* or *up to date*. Details show the version in `VERSION.md`, the systemd / supervisor units and a coloured **diff**.
- **Refs that make sense:** only tags on the chosen branch count, versions sort numerically (`v1.10` after `v1.9`), and you can check against the catalog, one branch for all, or each project's default branch.
- **Few API calls:** one tree listing per version, and each file content is downloaded once, ever.
- **Monitoring:** turn on periodic checks and get **one notification per round** about **new versions**, **drift** and **unreachable servers** through your notification channels. Checks never change anything.
- Administrators: policies `git_enabled` (hides it completely), `git_checks`, `git_check_interval_minutes`, `git_backup_words`.

Updating, upgrading, rollback and restarts follow in the next phases.

#### ⬆️ Upgrading from v10.10.0

Replace the binary. The database gets new tables (`git_sources`, `git_trees`, `git_blobs`, `git_workspace`, `git_targets`, `git_installs`); the previous binary still starts on it. The Git workspace is on for everyone; set `git_enabled` to off (or `WRM_GIT_ENABLED=0`) to hide it.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v11.0.0/CHANGELOG.md) for details.

---

## Included from v10.10.0 — folder bookmarks

#### ★ Folder bookmarks

- The same deep directories on many servers, one click away — like WinSCP's bookmarks. **★ next to the path bar** of the file manager: *Bookmark this directory…*, your bookmarks for this connection, *Manage bookmarks…*.
- **Scope:** one connection, the connections of a folder, the connections with a tag, or **all SSH / SFTP / FTP connections**.
- **Variables:** `~`, `$USER`, `{host}`, `{name}` — `/opt/app/servers/{name}/downloads` opens the right folder on every server.
- **In the terminal:** ★ in the title bar and the Ctrl+Shift+Space picker type `cd -- '<path>'`, safely shell-quoted. With broadcast input, 📣 sends it to every terminal of the group.
- **Start directory:** mark a bookmark and the file manager opens there and a terminal `cd`s there on connect.
- A missing directory gives a clear message instead of an error. Bookmarks can be renamed, reordered, moved between scopes, shared with share members (global ones, read-only), exported and imported with your configuration.
- **Import from WinSCP:** *Manage bookmarks → Import WinSCP.ini*.

#### ⬆️ Upgrading from v10.9.1

Replace the binary. The database gets a new table (`bookmarks`); the previous binary still starts on it.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.10.0/CHANGELOG.md) for details.

---

## Included from v10.9.1 — PuTTY session import

#### 📥 Import from PuTTY

- *Settings → Data → Import from PuTTY* reads a `.reg` export of your saved sessions (`regedit` → `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions` → *Export*, UTF-16 or UTF-8).
- SSH and Telnet sessions become connections with host, port and user; *Default Settings* supplies the defaults. Key files get a hint to convert the `.ppk`.
- **Proxies come along:** SOCKS4, SOCKS5 and HTTP proxies (with login and DNS setting) become saved proxies, one per distinct proxy, reused when you already have it, and linked to the connections. PuTTY's *SSH proxy* becomes the jump host.
- **With mRemoteNG:** connections that use a PuTTY session (*PuttySession*) get that session's proxy, whichever file you import first.

#### ⬆️ Upgrading from v10.9.0

Replace the binary. The database gets a new table (`putty_sessions`) and column (`connections.putty_session`); the previous binary still starts on it.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.9.1/CHANGELOG.md) for details.

---

## Included from v10.9.0 — proxies, folder defaults, notifications

#### 🌐 Proxies as a connection setting

- Define a proxy once and pick it per connection, like a vault credential: **SOCKS5** (user / password, DNS at the proxy), **SOCKS4 / 4a**, **HTTP CONNECT** (Basic auth), or the **SOCKS tunnel of one of your connections**, which WRM starts when needed.
- **Everything goes through it:** terminal, SFTP, FTP, tunnels, key deploy and rotation, RDP / VNC / Telnet, web interfaces, *Test connection*, live status, network tools from the server, Redfish.
- **With a jump host** the proxy is reached through it. The route is shown and audited: `bastion → socks5://10.1.1.1:1080 → app-01`.
- Clear errors: is it the **proxy** (unreachable, login) or the **target** (*cannot reach 10.0.0.5:22: connection refused (SOCKS5 reply 5)*)?
- Share proxies without revealing their login; policy `proxies` decides who may define them. IPMI / Serial-over-LAN (UDP) are refused through a proxy.
- **Docker:** `host.docker.internal` points to the Docker host, and WRM warns about `127.0.0.1` proxies inside a container.

#### 📁 Folder defaults

- Right-click a folder → *Folder settings…*: a **default jump host** and **default proxy** for all its connections, also new ones.
- A connection shows what it inherits (*— as the folder: bastion —*) and can override it (*no jump host (not the folder's)*). Export and import keep it.

#### 🔔 Notifications

- Administrators add channels: **e-mail (SMTP)**, **Telegram**, **Slack / Mattermost / Rocket.Chat**, **Microsoft Teams**, **Discord**, **ntfy**, **Gotify**, **Pushover**, or a signed **webhook** — with a *Send test* button.
- Users choose the events and channels in *Settings → Notifications*, with their own address where allowed and optional **quiet hours**.
- First events: **a server went down / is up again**, **credential rotation due / incomplete**. One check round sends **one message**, not a flood.

#### ⬆️ Upgrading from v10.8.1

Replace the binary. The database gets new tables and columns (proxies, folder defaults, notifications); the previous binary still starts on it. Connections without a jump host of their own now use their folder's default once you set one. Docker users: pull the new `docker-compose.yml` for `host.docker.internal`.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.9.0/CHANGELOG.md) for details.

---

## Included from v10.8.1 — fixes: larger imports, install hints

v10.8.1 is a fix release.

#### 📥 Larger imports work again

- **mRemoteNG** (`confCons.xml`), **OpenSSH config** and **inventory** files (CSV / Excel) of more than about 1 MB failed with *Bad JSON*. All imports now take files up to **20 MB**.
- A larger file gets a clear message with its size and the limit, already in the browser before anything is sent (HTTP 413 from the API).

#### 📲 Install as an app: what to do in your browser

- *Settings → General → Install as an app* shows the steps for the browser you use: Chrome / Edge / Brave, Safari on macOS (*File → Add to Dock*), iPhone / iPad (*Share → Add to Home Screen*), Android.
- When installing cannot work, it says why: WRM opened over plain `http://` on an IP address or name (browsers need HTTPS with a trusted certificate), an untrusted certificate, or Firefox on the desktop.
- The *Install* button appears only when it works.
- The README has a short **HTTPS with Caddy** recipe (a public name with Let's Encrypt, or Caddy's own CA for internal names and IP addresses).

#### ⬆️ Upgrading from v10.8.0

Replace the binary. Nothing changes in the database. Behind a reverse proxy, allow request bodies of at least 32 MB for imports.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/CHANGELOG.md) for details.

---

## Included from v10.8.0 — quick connect, notes, network tools & installable app

#### ⚡ Quick connect

- Type `root@10.0.0.5:22` into the field at the top of the sidebar and press Enter — or `rdp://admin@10.0.0.20`, `vnc://…`, `telnet://sw1`, `https://10.0.0.11/ui`.
- Log in with a password, a **vault credential** or a **stored SSH key**, optionally through a **jump host**.
- The connection waits under *Quick connections* for 24 hours after its last use; 💾 keeps it as a normal connection.

#### 📝 Notes & runbooks

- Every connection has notes: what runs there, how to restart it, who is on call, links to the wiki.
- A small markdown (headings, lists, code, links), shown next to the connection, in its terminal windows and in the context menu.
- Only the owner sees them; exported and duplicated with the connection.

#### 🧰 Network tools

- **Port check** (up to 100 ports and ranges), **ping**, **traceroute**, **DNS** and an **HTTP/TLS check** (status, certificate, issuer, days left, trust).
- From the WRM server **or from any of your servers, through its jump hosts** — "can app-01 reach the database on 5432?" Port and HTTP checks need nothing installed on the server.
- Audited; policy `network_tools` (all users by default).

#### 📲 Installable app

- *Settings → General → Install as an app*: WRM in its own window, from the dock, start menu or home screen (HTTPS with a trusted certificate, or `localhost`).

#### 🔧 Fixed

- The top bar no longer overlaps its buttons at medium window widths.

#### ⬆️ Upgrading from v10.7

Replace the binary. The database gets two new columns (`connections.notes`, `connections.temp_until`).

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.0/CHANGELOG.md) for details.

---

## Included from v10.7.0 — out-of-band management, consoles & serial ports

#### ⚡ Power & status through the BMC

- Set the **BMC** of a server in its connection — **Redfish** (iDRAC, iLO, XClarity, Supermicro, OpenBMC) or **IPMI** (ipmitool) — with its own login or a **vault credential**, directly or **through the jump host**.
- **Power & status:** power state, health and faults, model, serial number, BIOS, CPUs and memory, BMC firmware, power draw, inlet temperature.
- **Power on, shut down, restart, force off, force restart, power cycle, NMI** — only what the BMC offers, each confirmed and written to the audit log.
- Self-signed BMC certificates are pinned on first use, like SSH host keys.

#### 📟 Serial console

- The console of the machine in a terminal window — BIOS, boot loader, kernel messages, a login when the network is down.
- **Redfish BMCs:** over the BMC's SSH with the console command (iDRAC `console com2`, iLO `vsp` — suggested automatically).
- **IPMI:** Serial-over-LAN, on the jump host or the WRM server.
- Recorded and audited like every terminal; only the connection's owner opens it.

#### 🔌 Serial ports

- Console cables on the WRM machine (`ttyUSB0`, `ttyS0`) as connections of their own — speed, data bits, parity, stop bits, flow control — for switches, routers and firewalls.
- Administrators only by default (policy `serial_ports`); Linux.

#### 🔧 Fixed

- Saved workspace sessions now keep remote desktop windows (and the new console windows).

---

## Included from v10.6.0 — tags & inventory import

#### 🏷 Tags

- **Tags on every connection** — free (`web`, `db`) or with a key: `env:prod`, `site:zg1`, `rack:a12`, `role:`, `tenant:`, `platform:`.
- **Filter by tags** with one click in the sidebar's tag bar, or search `#rack:a12`.
- **Production stands out:** a red **PROD** badge and a red edge (staging yellow, test blue, dev green).
- **Bulk tagging** of selected connections (`env:test, -old`).

#### 📥 Import from CSV / Excel

- **CSV** — commas, semicolons or tabs, UTF-8 or Windows-1250 as Excel saves it here (č ć š ž đ stay right) — and **Excel .xlsx** with sheet selection.
- WRM **guesses the columns** (English and Croatian names: *Naziv, IP adresa, Okruženje, Lokacija, Rack…*) and shows a **preview** where every column can be changed.
- Hosts with ports, users, CIDR or URLs; Windows platforms become RDP; environment, site, rack, role and platform become tags; jump hosts by name.

#### 🗄 Import and sync from NetBox

- **Devices and virtual machines** through the NetBox API, with filters (sites, roles, tenants, tags, status), primary IP or DNS name, and the environment from a custom field.
- Sites, racks, roles, tenants, platforms and NetBox tags become **tags**; folders per site, role or tenant.
- **Sync:** load again and update the address, NetBox tags and folder of the connections imported before — your own tags, names and logins stay. Objects that disappeared from NetBox are listed, never deleted.
- Address and token can be remembered (the token encrypted, never shown again).

**Log in right away:** imported connections can get a **vault credential** (`root@dc1`) or a **stored SSH key** as their login.

#### ⬆️ Upgrading from v10.5

Replace the binary. The database gets three new columns for tags and inventory ids and a table for saved inventory sources; the previous binary still starts on it.

---

## Included from v10.5.0 — SSH keys & credentials vault

#### 🔑 SSH key management

- **Key store** (🔑 *Keys* in the top bar): **generate** ED25519, RSA or ECDSA keys, or **import** existing ones (also with a passphrase, also public keys of colleagues). Private keys are stored encrypted and never leave the server unless you export them (password confirmation, audited).
- **Deploy like `ssh-copy-id`** to any number of SSH connections at once — through jump hosts too — and optionally **switch the connections to the key** right away (WRM tests the key login first). Right-click a connection → *Deploy SSH key…*, or select many connections → **🔑 Key**.
- **Who has access:** right-click a connection → *Who has access (authorized_keys)…* lists every key that can log in as that user, with comment, fingerprint and options, names the keys from your key store, marks WRM's own login key and **removes** any other key with one click.
- **Revoke** a key from many servers at once. WRM never removes the key it logs in with.

#### 🗝 Credentials vault

- **One login for many connections:** a credential (`root@dc1`, `Administrator`, the iDRAC admin…) holds a user name and a password and/or an SSH key. Connections choose *Credential from the vault* — SSH, SFTP, RDP, VNC, Telnet and FTP. Change it once, every connection uses the new one.
- **Share without revealing:** grant a credential to colleagues (administrators: to everybody). They use it in their own connections but never see or change the secret. **Host lists** (`*.dc1.example.com`, `10.1.0.0/16`) limit where a credential can be sent — for the people it is shared with, through matching jump hosts only.
- **Password rotation** on all its SSH servers: a pre-flight login check everywhere, then `passwd` on each server, a login with the new password, and only then the new password is stored. **All or nothing**: when a server fails, the others get the old password back; if even that fails, the new password is kept as *pending* so nothing is lost. Live progress per server.
- **Rotation reminders** per credential (days), a badge on 🔑 and a warning in the admin overview; **Show** reveals the password to its owner after a password confirmation.
- Every action is in the audit log — never the secrets.

#### ⬆️ Upgrading from v10.4

Replace the binary. The database gets new tables for keys and credentials and two new columns (`connections.key_id`, `connections.credential_id`); the previous binary still starts on it. Deleting an account now also deletes its snippets, tunnels, keys and credentials.

---

## Included from v10.4.0 — RDP, VNC & Telnet in the browser

#### 🖥 Remote desktop: RDP, VNC, Telnet

- **RDP**, **VNC** and **Telnet** connections open in a WRM window, next to terminals and file managers — no plug-in, no client software on the PC.
- **Comfortable:** scale to fit or 1:1, full screen, the RDP desktop adapts to the window size, **Ctrl+Alt+Del**, **clipboard both ways**, **type text as keystrokes** (for login screens and consoles), sound from RDP, credentials prompts, reconnect.
- **RDP options:** domain, NLA/TLS/RDP/Hyper-V security, keyboard layout, colors, start program, console session, certificate, sound, clipboard, wallpaper and **RD Gateway**. **VNC:** colors, pointer, view only. **Telnet:** font, colors, automatic login prompts.
- **Behind jump hosts:** WRM tunnels the desktop through your bastions automatically for the length of the session.
- **Audited and recorded:** desktop sessions are in *Sessions & recordings* and the audit log, and the screen is recorded — **▶ Replay** in the browser or download as `.guac`.
- **Secure by design:** the browser never chooses the target; WRM does the handshake with the stored credentials and forwards only display data and allowed input.
- **mRemoteNG import** now also brings your RDP, VNC and Telnet connections.

#### ⚙️ What you need: guacd

Remote desktops use **guacd**, the proxy daemon of Apache Guacamole, on the same machine as WRM:

- Debian/Ubuntu: `sudo apt install guacd`
- Docker: `docker run -d --name guacd --network host --restart unless-stopped guacamole/guacd`
- `docker compose up -d` from the repository now starts guacd too.

*Admin panel → Overview* shows whether guacd answers. Everything else in WRM works without it.

#### 🔧 Fixed

- SQLite waits for a busy database on every connection now (rare *database is locked* errors under load).

#### ⬆️ Upgrading from v10.3

Replace the binary (and install guacd if you want remote desktops). The database gets one new column (`connections.options`); the previous binary still starts on it.

---

## Included from v10.3.0 — snippets, broadcast input & live status

#### ⚡ Snippets & run on connect

- **Saved commands** in the right panel (⚡ tab, groups, search) and in a **quick picker in every terminal: Ctrl+Shift+Space**. Click runs the snippet in the focused terminal, Shift+click types it without Enter.
- **Variables** per terminal: `{{host}}`, `{{user}}`, `{{name}}`, `{{folder}}`, `{{date}}` … and **prompts** with defaults: `systemctl status {{?Service=nginx}}` asks for the service first and shows a preview.
- **Scope:** all connections, one folder or one connection — the picker shows what fits the terminal.
- **Run on connect:** `sudo -i`, `cd /srv/app`, `tmux attach || tmux new`, … are typed automatically when a terminal of the connection (or folder) connects. Right-click a connection or folder → *⚡ Run on connect…*.
- **Shared snippets** from administrators for the whole team; **✨ example set** with the usual diagnostics (disk, memory, top processes, failed services, journal errors, listening ports, big files …).

#### 📣 Broadcast input

- Select terminals with **📣 Broadcast** and type into all of them at once. An orange bar and window frames show where your keystrokes go; 📣 in a title bar adds or removes a terminal.
- **Safety first:** dangerous-looking commands (`rm -rf`, `shutdown`, `mkfs`, `dd of=`, `iptables -F`, `DROP TABLE`, `kubectl delete`, …) and multi-line pastes ask before they are sent, with *Cancel* as the default button.
- **Snippets on all terminals** at once, each with its own host and user.
- Audited per terminal session (`terminal.broadcast`); administrators can turn it off (`broadcast_enabled`).

#### 🟢 Live up/down status

- A **status dot** for every connection: 🟢 up (latency, since when, SSH version), 🟡 slow, 🔴 down (since when and why), a ring for connections behind a jump host.
- **Down counters on folders**, a one-click **"show only down"** filter, and **notifications** (also desktop notifications in the background) when a server goes down or comes back.
- **🔄 Check status now** for a connection or a whole folder — also through the jump hosts.
- Checks like Nagios `check_ssh`: a TCP connection, the SSH/FTP greeting, a polite close; no login. Each host:port once per round (default every 60 s), with a retry before "down". Optional checks through jump hosts.
- Policies `status_enabled`, `status_interval_seconds`, `status_jump_checks`; per-connection opt-out.

---

## Included from v10.2.0 — jump hosts, SSH tunnels & mRemoteNG import

#### 🪜 Jump hosts (bastions)

- Every connection can go **through another saved SSH connection**: choose *Jump host (connect via)* in the connection. Like `ssh -J` / `ProxyJump` and the *SSH tunnel* setting of mRemoteNG.
- **Chains** up to 5 hops (`vpn-gw → bastion-dc1 → app-01`); loops are refused.
- **Everything works through it:** terminal, file manager, FTP/FTPS, search, editor, server-to-server transfer, *Test connection*, tunnels and web interfaces.
- Every hop uses **its own credentials** and its **host key is verified**. No port is opened on the jump host (`direct-tcpip`, the bastion only needs `AllowTcpForwarding yes`).
- The route is shown everywhere: **⤳** in the sidebar, *Test connection* (`… via bastion-dc1`), the terminal, session history and the audit log.

#### 🔀 SSH tunnels (port forwarding)

- **Local (-L):** a port on the WRM machine → a host:port behind the server, e.g. `127.0.0.1:2443 → 10.0.0.5:443` for the web UI of a switch, iDRAC or iLO.
- **SOCKS proxy (-D):** one proxy for your browser or tools that reaches the whole network behind a server.
- **Remote (-R):** a port on the server → a host:port on the WRM side (administrators, policy).
- Configure them in the connection (**🔀 SSH tunnels**, with **templates**: web HTTPS/HTTP, SSH, RDP, VNC, PostgreSQL, MySQL, iDRAC/iLO/IPMI, SOCKS) and use the new **🔀 Tunnels** panel: state, listen → target, connections, traffic, **🌐 Open**, **📋 Copy address**, **▶ Start / ■ Stop**.
- **Start modes:** manual, **with the terminal** (starts with the first terminal and stops after the last one), **always** (starts with WRM).
- **Reliable:** each tunnel keeps its own SSH connection through the jump hosts, checks it every 30 s and **reconnects automatically** while its port stays open.
- **Safe defaults:** tunnels listen on `127.0.0.1`; network addresses and remote forwards are for administrators; a tunnel can be turned off or limited to administrators; every start, stop and error is **audited** with the traffic. Administrators see and stop all tunnels in *Admin panel → Overview*.

#### 🌐 Web interface connections

- New protocols **Web interface (HTTPS / HTTP)** with host:port and path: save the web UIs of your iDRACs, iLOs, switches, firewalls and NAS next to your servers.
- **Double-click** opens it in a new tab — directly, or, **behind a jump host, through an automatic temporary tunnel** (closed after 30 minutes without traffic).

#### 📥 Import from mRemoteNG and OpenSSH

- **mRemoteNG `confCons.xml`:** folders, SSH and web connections **with their passwords** (both encryption formats, master password, full-file encryption), inherited settings, and **SSH tunnel → jump host**. RDP/VNC/Telnet are listed as skipped; duplicates are skipped.
- **OpenSSH `~/.ssh/config`:** hosts, users, ports, `Host *` defaults, **`ProxyJump`** / `ProxyCommand ssh -W`, **`LocalForward` / `RemoteForward` / `DynamicForward`** and `IdentityFile`.
- *Settings → Data* shows a summary of what was imported, skipped and what to check. Export/import of WRM's own format now includes jump hosts and tunnels.

#### ⚙️ New policies

`tunnels_enabled` (`TUNNELS_ENABLED`, on), `tunnel_users` (`all` / `admins`), `tunnel_bind_any` (off), `tunnel_remote_forward` (`off` / `admins` / `all`, default `admins`), `tunnel_idle_minutes` (0). All in *Admin panel → Security policies → SSH tunnels*, or forced with `WRM_<KEY>`.

---

## Included from v10.1.0 — audit trail & session recording

#### 🎬 Session recording & replay

- **Every SSH terminal session is recorded** in the open **asciinema** format (asciicast v2), compressed, with its SHA-256 checksum stored in the database.
- **▶ Replay in the browser** (play/pause, seek, 0.5–16× speed, *skip idle time*), or **⬇ download the `.cast`** and play it with `asciinema play`.
- **Transparent for users:** the terminal shows *“This session is recorded”* and a **● REC** badge.
- **No passwords in recordings:** servers do not echo passwords, so they never appear in the output. Keystroke recording is **off** by default. When an administrator turns it on, typing at password, passphrase and PIN prompts is masked.
- **Never slows the terminal:** the recording is written in the background. A per-session size limit stops the recording, not the terminal.
- **Who can watch:** administrators all; users their own sessions and sessions on **their** connections (e.g. guests of a share). Watching or downloading a recording is itself audited.
- **Where to find it:** *Admin panel → Sessions & recordings* (everything), *Settings → Session history* (your own).

#### 🧾 Audit trail

- **Terminal sessions:** who, from which IP, which server and remote user, start, duration, and how the session ended (closed, connection lost, ended by an administrator, interrupted by a restart).
- **File transfers:** every upload, download (also files opened in the editor and every file of a ZIP download) and server-to-server copy, with source, destination, size and **SHA-256**. See *Admin panel → File transfers*, with CSV export.
- **Append-only, tamper-evident audit log.** Database triggers refuse to change audit entries, transfers and recordings, and refuse to delete anything younger than 7 days. Only the retention job removes old data. Every entry contains the hash of the previous one, and **🔏 Verify integrity** checks the whole chain.
- **Better audit search:** date range, events linked to the connection and the terminal session (click **▶ #id** to replay), recording events, and secrets in event details are always redacted.
- **Feature flags:** `audit_enabled` (`AUDIT_ENABLED`) and `session_recording` (`SESSION_RECORDING_ENABLED`) are on by default. Also new: `session_recording_input`, `recording_retention_days` (90) and `recording_max_mb` (100). Administrator actions are always recorded.

#### 🐳 Operations

- **Docker:** `Dockerfile` (static binary, unprivileged user, data volume `/data`, health check) and **`docker compose up -d`**. Release images are published to `ghcr.io/vedranius/wrm-pro`.
- **`/healthz`** for load balancers and monitoring; **`wrm -healthcheck`** for container health checks.
- Recordings location: `WRM_RECORDINGS_DIR` (default `recordings/` next to the database).

#### 🎨 New logo

The WRM PRO logo is now used for the favicon, the app icons, the sign-in screen, the top bar, the README and the release. A **WRM Orange** accent color is available in *Settings → General*. The complete logo kit (SVG and PNG) is in `docs/brand/`.

#### 📐 Architecture & tests

- [`ARCHITECTURE.md`](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.0/ARCHITECTURE.md) describes how WRM is built. It also maps the extension plan (RBAC, SSO, broadcast input, fleet status, …) to the code.
- A new **integration test** runs a real SSH + SFTP server inside the test. It checks connect → audit entries → recording (replayable, no passwords), transfer checksums and the append-only audit trail, and it runs in CI with the race detector.

#### ⬆️ Upgrading from v10.0

Replace the binary. The database is only **extended**: new tables, columns and triggers are added, and nothing is changed or removed. The previous binary still runs on the upgraded database.

- **Session recording is on by default.** Tell your users about it (in many countries this is required), or turn it off in *Admin panel → Security policies → Session recording*.
- Recordings use disk space in `recordings/` next to the database. A session usually takes a few KB to a few MB. Recordings are kept for 90 days, with at most 100 MB per session.
- Audit entries written before the upgrade have no hash. *Verify integrity* checks the chain from the first new entry on.

---

## Included from v10.0.x

### 🔧 v10.0.1

- `/static/` (the v9 address) forwards to the app instead of returning 404.
- A database from an earlier build with an incompatible `audit_log` / `collab_messages` table is upgraded at startup. The old table is kept as `<table>_old_<time>`.
- The server log shows the address to open in the browser.

### v10.0.0

#### 🎙 Voice calls in collaboration rooms (back, rebuilt)

Everyone in a share room can talk. It works like Discord or Jitsi:

- **🎙 Join voice** in the collaboration bar; **mute** (Ctrl+Shift+M), **deafen** (Ctrl+Shift+D), **leave**.
- **Push-to-talk** with a key of your choice (e.g. F8), or voice activity.
- **Microphone and speaker selection**, live input level meter, test sound, **noise suppression**, **echo cancellation**, **automatic gain control**.
- **Speaking indicators** (green ring), **per-person volume**, join/leave sounds, connection quality per person (P2P or relay, round-trip time).
- **Listen-only**: without a microphone, or without permission to use it, you can still join and listen.
- **Stable on real networks**: WebRTC with "perfect negotiation", automatic ICE restarts when the network changes, and automatic rejoin after reconnects.
- **Built-in TURN relay** (UDP and TCP, port 3478): calls work behind NAT and strict firewalls with no extra software. It accepts only short-lived credentials issued to room participants. By default it refuses to relay to private or loopback networks.
- **End-to-end encrypted** audio (DTLS-SRTP). The relay cannot decrypt it.
- Moderators can **mute someone for everyone**. Admins can turn voice on or off, limit participants per call (default 12), and add external STUN/TURN servers.

> Browsers allow the microphone only on **HTTPS** pages (or `localhost`). The quickest way to get HTTPS is `HTTPS_SELF_SIGNED=1`. For calls across the internet, open **UDP+TCP 3478** and the **UDP relay range** (default `49152-65535`) in the firewall. Behind NAT, set `turn_public_ip`.

#### 👥 Users, members & roles

- **Admin panel → Users**: create users (a temporary password is generated if you leave it empty), set the display name, make or remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete. The last active administrator cannot be removed.
- **Shares have three access modes**: 👥 *Members only* (default), 🏢 *Everyone signed in*, 🌐 *Anyone with the link* (guests; admins can disable this mode).
- **Members** with an individual role each, plus a role for everyone else. There are **five roles**:

  | | Observer | Viewer | Operator | Moderator | Owner |
  |---|:-:|:-:|:-:|:-:|:-:|
  | Chat & voice | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Watch shared terminals | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Browse & download files | — | ✔ | ✔ | ✔ | ✔ |
  | Open terminals, change files, transfer | — | — | ✔ | ✔ | ✔ |
  | Share own terminal, receive keyboard control | — | — | ✔ | ✔ | ✔ |
  | Change roles, mute, remove, ban people | — | — | — | ✔ | ✔ |
  | Share settings & members | — | — | — | — | ✔ |

  The **server enforces roles** on every file operation, terminal and room action.
- **People** view per share: everyone who ever opened it, with last seen and IP. From there you can change a person's role or **ban** them.
- **Share expiry** (1 hour to 30 days, or a custom date), **pause/resume**, **new link** (the old link stops working immediately).
- **Revocation takes effect at once.** Deleting, pausing or rotating a share, removing a member, lowering a role, banning, or disabling a user closes their open terminals and room connections.

#### 💬 Collaboration

- A **new collaboration bar** shows the share, your role, avatars of the people online, and the voice controls.
- **👥 People panel**: who is in the call and who is online, with roles, mute and deafen state, shared terminals, raised hands, connection quality and volume. Moderators have a **⋯** menu for *change role*, *mute*, *stop sharing*, *lower hand*, *remove* and *ban*.
- **Chat history** is kept, so people who join later see it. The chat has timestamps, clickable links, an unread badge and notifications. **File exchange** up to the policy limit: images get a preview, and other files are always downloaded, never opened in the page.
- **✋ Raise hand.**
- **Terminal sharing**: you can share several terminals at once. Viewers get a **snapshot of the current screen** (with colors) and then the live output. Terminal data is sent only to people who watch it.
- **Remote keyboard control**: *Request control* → *Allow / Deny*. Keystrokes go only to the granted terminal, and only while it is shared.
- **Identities come from the server.** Signed-in users appear under their account name. Guests pick a name and are marked *guest*.
- The room **reconnects automatically** and restores watching, sharing and the voice call.
- **Per-participant send queues**: one slow client can no longer stall a room.

#### 🔐 Security hardening

- **Two-factor authentication (TOTP)** with QR code and single-use **recovery codes**. The policy can **require 2FA** for admins or for everyone.
- **Self-registration is closed by default.** Only the first account (the administrator) registers itself.
- **Accounts lock** after repeated wrong passwords (policy). Per-IP rate limiting. Unknown users get constant-time responses. After an admin reset, the user **must change the password**.
- **Sessions**: 256-bit tokens, **stored hashed**. Idle and maximum lifetime (policy). A list of **signed-in devices** with remote sign-out.
- **Stored secrets** (connection passwords, private keys, 2FA secrets) are encrypted with **AES-256-GCM** and a **random key per installation** (`remote_manager.db.key`, or your own `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`).
- **Secrets are never sent to the browser**, not even to their owner. The *Edit connection* dialog shows "saved and encrypted — leave empty to keep it". *Export* leaves secrets out. *Export with passwords & keys* asks for your password again.
- **SSH host key verification** (trust on first use, or strict). A changed key is refused and shown with its fingerprint. **FTPS certificates** are validated, or pinned on first use.
- **Server-side key files** (*Key file* / *Auto (~/.ssh)*) are limited to administrators by default.
- **CSRF protection** on every state-changing API call. **Content-Security-Policy**, `X-Frame-Options`, `nosniff`, `Referrer-Policy: no-referrer`, `Permissions-Policy`, and HSTS on HTTPS. Request-size limits. No directory listings.
- **All assets are served locally** (xterm.js and fonts are embedded). There is no CDN dependency at runtime.
- **Share passwords** are held in signed cookies, which become invalid when the password changes. Guests get their own identity cookie. Share connection lists stay hidden until the share password is entered.
- **HTTPS made easy**: `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE`, or `HTTPS_SELF_SIGNED=1`. Behind a reverse proxy, `WRM_TRUST_PROXY=1` gives correct client IPs and secure cookies.
- **Audit log** of sign-ins (failed ones too), account and policy changes, connections, shares, room joins and moderation, terminals, file changes, host keys, exports and imports. It is searchable in the admin panel, exportable as **CSV**, and written as `AUDIT …` lines to the server log (journald/syslog → SIEM).
- **Recovery from the command line**: `wrm -reset-password USER [-reset-2fa]`.

#### 🛡 Admin panel

*Settings → 🛡 Admin panel* has these tabs:
- **Overview**: version, uptime, users, HTTPS, relay status, key source, **security warnings**, **open terminals** (with *End*), live rooms.
- **Users**
- **Shares**: every share on the server.
- **Security policies**
- **Voice & network**
- **Host keys**
- **Audit log**

Any policy can also be **forced by an environment variable** (`WRM_<KEY>`, e.g. `WRM_REQUIRE_2FA=all`). A forced policy is shown as locked in the panel.

#### ✨ Other improvements

- **Settings** are reorganised into tabs: General, Terminal, Security, Voice & audio, Data.
- **Connections**: SSH keyboard-interactive authentication (servers that ask for the password that way). *Duplicate* copies the stored secrets on the server, so they never pass through the browser. Deleting a connection closes its open terminals.
- **Uploads** respect the server-side size limit (`max_upload_mb`), and partial files are removed.
- Translations (English / Hrvatski) cover the whole new UI.
- A **unit test suite** (TOTP, recovery codes, roles, share access, CSRF, encryption and migration, signed values, settings, input validation) runs in CI with the race detector.

---

#### ⬆️ Upgrading from v9

Replace the binary and start it with the same database, and the same `ENCRYPTION_KEY` if you set one. Everything is migrated automatically.

- **Stored secrets are re-encrypted.** If `ENCRYPTION_KEY` was not set, v9 used a public default key. v10 generates a random key file, `remote_manager.db.key`, next to the database and re-encrypts all secrets with it. **Back that file up separately from the database.** If you did set `ENCRYPTION_KEY`, nothing changes.
- **Self-registration is now closed.** Existing accounts keep working. To reopen it: *Admin panel → Security policies*.
- **SSH host keys are now verified.** The first connection to each server remembers its key.
- **Key file / Auto (~/.ssh)** authentication is now admin-only, unless you turn on the policy *Allow server key files for all users*.
- **Existing shares keep their behaviour**: anyone with the link can open them, with the role *Operator*. Edit a share to switch it to members only or to change roles.
- **Behind a reverse proxy**, state-changing API calls now require the `Origin` to match the `Host`. Pass the `Host` header through (`proxy_set_header Host $host;`) or set `WRM_ALLOWED_ORIGINS`. Also set `WRM_TRUST_PROXY=1`.
- For voice calls, open **UDP+TCP 3478** and the relay port range in the firewall (see above).

---

---

### 📦 Downloads

| Platform | Architecture | Binary |
|---|---|---|
| Linux | x86-64 | `wrm-pro-v11.0.0-linux-amd64` |
| Linux | arm64 | `wrm-pro-v11.0.0-linux-arm64` |
| Linux | ARMv7 (Raspberry Pi) | `wrm-pro-v11.0.0-linux-armv7` |
| Linux | ARMv6 | `wrm-pro-v11.0.0-linux-armv6` |
| Linux | 32-bit | `wrm-pro-v11.0.0-linux-386` |
| Windows | x86-64 | `wrm-pro-v11.0.0-windows-amd64.exe` |
| Windows | arm64 | `wrm-pro-v11.0.0-windows-arm64.exe` |
| macOS | Intel | `wrm-pro-v11.0.0-darwin-amd64` |
| macOS | Apple Silicon | `wrm-pro-v11.0.0-darwin-arm64` |
| macOS | Universal | `wrm-pro-v11.0.0-darwin-universal` |
| Android | arm64 (Termux) | `wrm-pro-v11.0.0-android-arm64` |
| FreeBSD | x86-64 | `wrm-pro-v11.0.0-freebsd-amd64` |
| FreeBSD | arm64 | `wrm-pro-v11.0.0-freebsd-arm64` |
| OpenBSD | x86-64 | `wrm-pro-v11.0.0-openbsd-amd64` |

Verify integrity with `SHA256SUMS.txt`. The Android build has no built-in TURN relay; configure an external TURN server there if you need one.

### 🚀 Quick start

**Linux / macOS**
```bash
chmod +x wrm-pro-v11.0.0-linux-amd64
HTTPS_SELF_SIGNED=1 ./wrm-pro-v11.0.0-linux-amd64
# open https://<server>:8080 — create the administrator account (the first account)
```
On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

**Windows**: double-click `wrm-pro-v11.0.0-windows-amd64.exe`, or in PowerShell:
```powershell
$env:HTTPS_SELF_SIGNED=1; .\wrm-pro-v11.0.0-windows-amd64.exe
```

**Android (Termux)**
```bash
pkg install wget
wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v11.0.0/wrm-pro-v11.0.0-android-arm64
chmod +x wrm-pro-v11.0.0-android-arm64 && ./wrm-pro-v11.0.0-android-arm64
```

**Docker**
```bash
docker run -d --name wrm -p 8080:8080 -v wrm-data:/data -e HTTPS_SELF_SIGNED=1 ghcr.io/vedranius/wrm-pro:v10.9.0
# or, from the source tree:  docker compose up -d
```

On the first start WRM creates `remote_manager.db` (the database) and `remote_manager.db.key` (the encryption key) next to each other, both with mode `0600`. **Back up the key file.**

### ⚙️ Most important settings

```
PORT=8080 / LISTEN_ADDR=:8080     listening address
DB_PATH=./remote_manager.db       database
ENCRYPTION_KEY / ENCRYPTION_KEY_FILE   your own key for stored secrets (else <DB_PATH>.key)
HTTPS_CERT_FILE + HTTPS_KEY_FILE  HTTPS with your certificate
HTTPS_SELF_SIGNED=1               HTTPS with a generated self-signed certificate
WRM_TRUST_PROXY=1                 behind a reverse proxy (real client IPs, secure cookies)
WRM_ALLOWED_ORIGINS=host          extra allowed Origin hosts
WRM_<POLICY>=value                force a policy, e.g. WRM_REQUIRE_2FA=all, WRM_TURN_PUBLIC_IP=203.0.113.10
AUDIT_ENABLED / SESSION_RECORDING_ENABLED   audit log / session recording on (1) or off (0)
TUNNELS_ENABLED=0                 turn SSH tunnels off (WRM_TUNNEL_USERS=admins: administrators only)
STATUS_ENABLED=0                  turn the live up/down status off (WRM_STATUS_INTERVAL_SECONDS=300: check every 5 minutes)
BROADCAST_ENABLED=0               turn broadcast input off
GUACD_ADDRESS=127.0.0.1:4822      where guacd runs (remote desktop); REMOTE_DESKTOP_ENABLED=0 turns it off
WRM_RECORDINGS_DIR=/path          where session recordings are stored (default: recordings/ next to the database)
WRM_NETWORK_TOOLS=admins          network tools for administrators only (off: turn them off)
```

For the full documentation (Docker, reverse proxy, systemd, firewall, API), see the [README](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/README.md). To report a vulnerability, see [SECURITY.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/SECURITY.md).

---

### 📜 License

Web Remote Manager PRO is **source-available** under the [PolyForm Noncommercial License 1.0.0 or the PolyForm Internal Use License 1.0.0](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/LICENSE). It is free for personal, educational, non-profit and other noncommercial use, and **free for companies that use it as a work tool**, including paid work for their customers. **Offering WRM as a hosted service, charging for its use, reselling or bundling it requires a commercial license**; see [COMMERCIAL-LICENSE.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/COMMERCIAL-LICENSE.md). Contributions are welcome; see [CONTRIBUTING.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.8.1/CONTRIBUTING.md).

☕ **Like WRM?** Support its development on **[Ko-fi](https://ko-fi.com/vedranius)**. Thank you!
