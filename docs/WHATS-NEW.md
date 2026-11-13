# What's new in earlier versions

The [README](../README.md) shows the five latest versions. Earlier highlights are kept here, newest first; every change is in the [changelog](../CHANGELOG.md), and upgrade notes per version are in [UPGRADING.md](UPGRADING.md).

## What's new in v14.11.4

**AI apps wait for your approval, and you see it on the desktop.** Claude Code, Codex and other AI apps stop waiting for a tool call after their own timeout (Codex after 60 seconds), and then used to carry on without the step. Now WRM tells the AI that the request still waits in WRM and gives it a new tool, `wait_for_approval`, which it calls again and again until you approve or deny; only then does it go on. The WRM PRO plugin (and every AI app set up with *All AI apps at once*) also shows each approval request as a **notification of the operating system** (a Windows notification, macOS Notification Center, `notify-send` on Linux), also when no WRM tab is open or in front; a click opens WRM. Update the plugin to 1.2.4. See [AI desktop apps (MCP)](#ai-desktop-apps-mcp).

## What's new in v14.11.3

**AI apps stop asking for a new approval too early.** The WRM PRO Claude plugin (and the Claude Desktop extension, which share one bridge) asked you to approve the AI connection again whenever WRM did not answer for a moment, for example while WRM restarted, updated or hung, or while your computer was off the network, although the access you chose (say 7 days) was still valid. Now only a token that WRM refuses (expired or revoked) leads to a new sign-in: otherwise the AI app is told that WRM is unreachable, the bridge reconnects by itself with the saved access, and `wrm_connect` does not start a new sign-in while that access is valid. Update the plugin to 1.2.3 (or download the Claude Desktop extension again) to get it. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v14.11.2

**Workboard and panels in one place.** The **▤ Workboard** menu now lists every panel (*Connections*, *Sessions*, *AI sessions*, *Snippets*, *Clipboard*, the AI view and docked addons) with where it is, so you show, hide, lock and reset the layout from the same menu where you add addons and arrange the canvas. The separate ▦ *Panels* button is gone from the top bar, which leaves room for the other buttons; **＋ Add to the Workboard** lists the panels next to the addons, so a hidden panel comes back with *Show*. A locked layout shows 🔒 on the Workboard button. See [Workboard: the shared canvas and addons](../README.md#workboard-the-shared-canvas-and-addons).

## What's new in v14.11.1

**Your own order in the connection tree.** Drag a connection above or below another one to put it there; drop it on a connection in another folder to move it into that folder at that place. The order is yours alone: other users keep theirs. Right-click a connection for **Move connection up / down** and **Sort connections here A–Z**, or a folder for **Sort connections in this folder A–Z**. See [Sorting folders and connections](../README.md#connections--folders).

## What's new in v14.11.0

**An interface for every screen.** The top bar keeps its labels: less used buttons move into the labelled **⋯** menu first, so a 1920 px screen shows every button with its name again and a laptop keeps the ones you use most. Your avatar opens an **account menu** with *Settings*, *Appearance*, *Shares*, *Lock now* and *Sign out*. On tablets and narrow laptops (761–1024 px, also a zoomed-in browser) the right panel with *Sessions*, *Clipboard* and *Snippets* slides in from the 🗂 button, so the work area keeps its width. Big dialogs can be **minimised** to a chip at the bottom of the screen (▁), and **notifications are grouped**: repeats get a ×N counter and older ones fold into one *more notifications* toast. See [Screen sizes, zoom and touch](../README.md#mobile-responsive-ui--installable-app).

## What's new in v14.7.2

**Shares a team can manage together** ([#192](https://github.com/vedranius/web-browser-RDM-public/issues/192)). Two new switches in the share dialog, both off by default: *Operators may add and edit connections and folders* and *Operators may start and stop tunnels*. Signed-in members with the operator role or higher then get **Edit connection…** on a shared connection, **New connection here… / New folder… / Rename folder…** on a shared folder and **Tunnels…** on the share. Everything stays in the owner's account, members never see or pick the owner's saved secrets, and changing a connection's host or protocol removes its saved password. See [Sharing: members, roles & links](../README.md#sharing-members-roles--links).

## What's new in v14.7.1

**Security fixes and the tunnel/HSTS bug.** On a fresh install the administrator account can only be created from the network with a one-time **setup code** that WRM prints to its console/log and writes to `first-admin-setup-code.txt` next to its database (not needed on the WRM host itself), so nobody else on the network can claim a new server first. A tunnel on 0.0.0.0 can be opened through an IP address of the WRM host, and *Admin → Policies* has an HSTS switch, so a self-signed web interface behind a tunnel opens in the browser ([#190](https://github.com/vedranius/web-browser-RDM-public/issues/190)). The security overview warns about tunnels open to the network, the origin check compares the port, and one AI client closing a shared AI session no longer stops the others.

## What's new in v14.7.0

**Addons can show your chains.** A new addon capability, `chains.read`, lets an addon window read the [chains](../README.md#chains-live-service-chain-monitoring) of whoever opens it: the checked values, the active scenario and the timeline of changes, and ask for a check now. WRM keeps checking in the background as before, so an addon can show the state of a whole system the way your team needs it (a table per environment, a report, a sync to a wiki) without SSH, commands or secrets of its own. The permission sheet names the new right, and `wrm -addon-dev` can feed sample chains from a JSON file. See [Addons](../README.md#addons).

## What's new in v14.6.3

**Fixes for proxies, the screen lock and log discovery.** A dynamic SSH tunnel (`-D`) of your connection now appears in every proxy list (connection, folder, *Proxies*) and becomes a proxy when picked, also with an automatic port ([#188](https://github.com/vedranius/web-browser-RDM-public/issues/188)). The lock screen keeps the keyboard on its password field even when a terminal reconnects behind it ([#137](https://github.com/vedranius/web-browser-RDM-public/issues/137)). Log discovery skips Python virtual environments (also renamed copies such as `venv.bak.2026-07-16`) and folders that do not answer instead of stopping, and an open log viewer or search finds the logs again by itself when the discovery expired instead of failing with *Find logs first*.

## What's new in v14.6.2

**Forgot password and account activation by e-mail.** The sign-in page offers *Forgot password?*: WRM mails a one-time link (valid one hour) to the account's confirmed address, and an account with 2FA also needs a code. Administrators can create a user with an e-mail address and *Send an invitation*, so the user chooses the password through a link instead of being given one. With open registration, a new account starts once its e-mail address is confirmed. At the Enterprise level administrators cannot use *Forgot password*. See [Forgot password, invitations and registration](../README.md#forgot-password-invitations-and-registration).

## What's new in v14.6.1

**WRM sends mail.** *Admin → E-mail* sets up one account with presets for Gmail, Outlook.com and Microsoft 365 (or any SMTP server) and a test button. Users confirm their e-mail address with a link and get a mail when their password is changed or reset, or two-factor authentication is turned on or off. See [E-mail](../README.md#e-mail).

## What's new in v14.6.0

**A password reset on the server is no longer a quiet way in.** `wrm -reset-password` now tells every administrator, and how strict it is follows the security level: *Alert* (Home) resets at once, *Only while WRM is stopped* (Business) refuses while WRM runs, and *Confirmation* (Enterprise) files a request that another administrator approves in Admin → Security, that the user confirms with their own 2FA code, or that can be completed after one hour if nobody rejected it. An AI working through WRM can never run it. See [Password reset on the server](../README.md#password-reset-on-the-server).

## What's new in v14.5.1

**Resizable windows and storage columns.** The big windows opened from the menus can be shown over the whole screen (⤢ above their top-right corner) or resized by their bottom-right corner, and keep their size. The storage browser sorts by clicking the column headers, widens columns by dragging their borders (double-click fits the longest name), and its side panels can be resized. See [Storages](../README.md#storages-s3).

## What's new in v14.5.0

**Chains recognise your servers.** The builder lists servers by folder with search and adds several at once, without quick connections and without adding a server twice. Right after adding them, WRM looks at the servers read-only and turns their services, containers, named ports (PostgreSQL (5432), HTTPS (443) …) and network mounts into one-click check suggestions. Mounts between servers suggest the links. A URL, `host:port` or host name becomes an outside endpoint with the fitting checks, and **Test now** tries one check without saving. See [Chains](../README.md#chains-live-service-chain-monitoring).

**No more freezing while an AI works through MCP.** The MCP live view used one connection per session, and over plain HTTP a browser has only six per server for all its tabs, so WRM stopped responding once an AI app had opened six sessions. All AI and MCP views of a tab now share one stream.

## What's new in v14.4.0

**See every S3 request, and what it costs.** The storage window has a live **📡 S3 activity** panel beside the file browser, shown by default and hidden with one button. It lists every request WRM sends to the bucket as it happens (operation, key, answer, time, and whether you, the background check, a chain or an AI connection caused it) and counts each operation for today and in total. The storage form has a **Provider** picker (Hetzner, DigitalOcean Spaces, Amazon S3, Backblaze B2, Cloudflare R2, MinIO, other). It says whether that provider bills the requests WRM uses, and how much, quoting the provider's price page. The activity panel prices the counted requests the same way.

**The AI works inside the Chains builder.** The 🤖 AI button on a chain or draft takes one instruction ("add a TLS check for example.com", "why is the backup path red?"). The AI reads the chain and its last results and writes a draft. A draft that changes an existing chain says so, and the chain changes only when you press **Apply changes**; it keeps its id and history. **Test once** runs a draft's checks a single time and shows the results on the diagram, without saving anything. WRM's AI panel gets `chain_get` and `chain_test`, and `chain_draft` can now change a chain. The same four tools are available to MCP clients (Claude Desktop, Claude Code and others) behind the new opt-in **`chains`** scope.

## What's new in v14.3.2

**A storage can be edited and removed again.** Each storage in the list has a ✎ (edit) and a 🗑 (remove) button, and the open storage has *Edit* in the toolbar; the form also has *Remove storage*. Renaming or changing permissions no longer contacts the bucket, so it also works while the endpoint is down. The background bucket check runs every 15 minutes instead of 5 (`storage_check_minutes`, `0` turns it off), a chain's storage check asks the bucket at most every 5 minutes instead of on every 30-second poll (`chain_storage_check_s`), and [docs/STORAGES.md](STORAGES.md) lists every S3 request WRM sends.

## What's new in v14.3.1

**Log reads and WRM responsiveness with Claude/MCP.** A log glob such as `follower.log*` is no longer mistaken for a PuTTY private key; explicit `.ppk` files are still refused. The database pool now has headroom for parallel MCP calls and several browser tabs, and connection/folder lists load their labels in one query instead of one per row, so WRM no longer becomes unreachable until a restart.

## What's new in v14.3.0

**Storage folders are shown completely.** The browser follows every S3 continuation token instead of stopping at the first page: it reports the running count while loading, then the total object and folder count and size with a clear "all objects in this folder are shown" line (and *Load more* above a 100-page safety limit). Search filters the whole folder, sort goes by name/size/date, and select-all, bulk delete and copy-path are available. **Test connection** now checks the values in the add form before saving.

**The S3 share links are gone for security.** The Share button, its API and the public `/s3/` page were removed; links created by an older version stop working, and storage access always needs a signed-in account with the storage's permissions.

**The Chains wizard works with four or more servers.** Check ids are kept inside the 40-character limit and a large "all checks OK" state is nested instead of exceeding the rule limit.

## What's new in v14.2.0

**Chains are now built, not typed.** A new guided builder picks servers from your saved connections, adds checks with a form, connects nodes (or drag from a node's connector dot, like a workflow canvas) and edits states (scenarios) without JSON. Templates and the advanced JSON editor are still there, but nothing has to start as JSON.

**The AI wizard looks at the servers first.** Pick servers, optionally describe the system, and WRM reads the OS, running services, Docker containers, listening ports and network/object-storage mounts (rclone, s3fs, NFS, CIFS, SSHFS) with fixed read-only checks; the assistant then drafts a chain from real service names and ports instead of invented ones. Without an AI provider the same facts build a deterministic starter.

**No storage required.** A chain of servers is complete on its own; discovered mounts become `mount` checks and, when the source names a host, an outside endpoint with a DNS check. WRM never attaches to a discovered storage, never asks for its credentials and never creates a storage entry for it.

## What's new in v14.1.0

**Chains: a whole system watched as one.** WRM now links the servers, storages and outside endpoints that form one working system (for example ingest → transcoding → content → CDN with a primary and a backup path), checks them in the background with read-only checks over the SSH connections it already has, and shows the chain as a live diagram whose links animate with traffic and turn red when something fails. Scenarios decide the active path (primary, backup, degraded), and every change is kept in a timeline and notified. **Describe the system to WRM's AI assistant**: `chain_status` and `chain_draft` let it write a draft with catalog checks only; you review that draft on the same diagram and approve it before any monitoring starts. Nothing is installed on a server, and a chain can never run a command. See [Chains](../README.md#chains-live-service-chain-monitoring).

## What's new in v14.0.0

**One combined release: Storages, its release review, and the security work of v13.14.11.** Everything below arrives together when you upgrade from v13.14.10 or older; from v13.14.11 the new part is Storages. Next on the [roadmap](../ROADMAP.md) is **Chains**, live monitoring of a whole service chain.

**Storages: S3-compatible object storage as a section of its own.** Add a Hetzner, DigitalOcean Spaces, Amazon S3, MinIO or Backblaze B2 bucket (confined to a prefix), browse it, preview images, play video with a sibling subtitle track, download, and edit small text objects (JSON, configuration, subtitles) with BOM/CRLF kept and a conflict check. Uploads, writes and deletes are allowed per storage; a delete needs an explicit confirmation. A background check notifies `storage.down` / `storage.up`, and AI connections get `storage_list/read/write/delete` behind the opt-in scopes `storage_read`, `storage_write` and `storage_delete` (changes only in the automatic mode). See [Storages (S3)](../README.md#storages-s3) and [docs/STORAGES.md](STORAGES.md).

## What's new in v13.14.11

**AI connections survive WRM restarts.** Claude's bridge recovers its in-memory MCP session with the saved grant, and the server list stays selected. Grant creation and OAuth consent also have **Select all** controls.

**Share links and security checks are stronger.** New share links keep bearer tokens out of request URLs; saved credentials stay bound to their connection destination; AI approvals are rechecked before execution; outbound requests and addon messages are constrained.

**Pending AI approvals are harder to miss.** WRM can show desktop notifications that focus the related approval.

## What's new in v13.14.10

**Self-update works on Windows when an antivirus scanner holds the new file.** WRM waits for the lock for about 15 seconds and copies the file into place if it stays locked, instead of stopping with *being used by another process*.

## What's new in v13.14.9

**The original deploy tool's Git section is explained, and locks name the right computer.** The original tool shows its *Git* section (pull, commit and push of deploy-config) only when deploy-config is a Git clone; a deploy-config imported as a ZIP has no Git history. Medium step 3 now says so and points to `DEPLOY_CONFIG_REPO` in *Original runtime settings*, whose fields show in grey what an empty field means. In the managed Docker runtime `DEPLOY_HOSTNAME` now defaults to the WRM computer's name instead of the container ID, so locks and deploy history name a real computer.

## What's new in v13.14.1

**Addon settings survive restarts, and a frozen WRM leaves a record.** Since v13.9.0 every start of WRM set the addon instances aside (for example a Grafana panel's configurations), so addons came back empty after a restart or an update; v13.14.1 keeps them and restores the most recent ones with their settings and window positions. WRM also writes what it was doing to `remote_manager.db.diag/` next to the database whenever it is stopped or restarted and when its own watchdog sees that it no longer answers (see [Troubleshooting](../README.md#troubleshooting)).

## What's new in v13.14.0

**The Git workspace checks only what Git tracks, folders sort themselves, and every file edit is on record.** The Git overview names servers by their hostname or connection name instead of an IP. A check compares only the files the catalog tracks: logs and other files outside the check are listed as *outside the check* and never compared, and the .gitignore helper has a **Logs** group. Every Git filter can narrow to a **folder on the servers** (a tree with counts, like the log viewer's folder filter), and the connection-folder filter includes subfolders. Right-click a folder to **sort its subfolders A–Z** or **keep them sorted** ([#134](https://github.com/vedranius/web-browser-RDM-public/issues/134)). *Admin → File edits* keeps **every file edit with its diff**, from the editor and from the AI, filterable by user, server, session and AI ([#129](https://github.com/vedranius/web-browser-RDM-public/issues/129)). SSH keys can log in with an **OpenSSH user certificate**, also the short-lived ones from **opkssh** ([#119](https://github.com/vedranius/web-browser-RDM-public/issues/119)). See [Connections & folders](../README.md#connections--folders), [SSH user certificates and opkssh](../README.md#ssh-user-certificates-and-opkssh), [File edit history](../README.md#file-edit-history) and [Git workspace: every file and partial checks](../README.md#git-workspace-every-file-and-partial-checks).

## What's new in v13.13.1

**Adding the WRM PRO marketplace in the Claude app.** Enter the repository in lowercase, `vedranius/web-browser-rdm-public`: the Claude app names the marketplace after the repository, and GitHub accepts the name in any case. If the app still says only *Failed to add marketplace*, [Troubleshooting](../README.md#troubleshooting) says what to check and where its log keeps the reason; *Upload plugin* with `claude-plugin-wrm-pro.zip` from the release installs the same plugin.

## What's new in v13.13.0

**The WRM PRO plugin installs from the Claude app with one pick, and one script connects every AI app.** The plugin needs no settings any more: it carries the WRM bridge, and the first `/wrm-pro:wrm` asks for the WRM address and opens WRM in the browser, where you sign in and approve; the token goes straight to your computer, never into the chat. `/wrm-pro:wrm` then lists your servers and everything WRM can do. Before, the plugin depended on settings the Claude app cannot ask for, and its HTTP connector would have connected from the cloud, which cannot reach a WRM on your own network. *Settings → AI connections → 🧰 All AI apps at once* downloads one script (Windows or Linux/macOS) that signs in once and adds WRM to Claude Code, the Claude app, Codex, Gemini CLI, VS Code, Cursor, Windsurf and OpenCode. Releases carry `claude-plugin-wrm-pro.zip` for *Customize → Plugins → Upload plugin*. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v13.12.0

**Tools for writing addons, also with an AI.** `wrm -addon-lint ./my-addon` checks an addon before you install it: the manifest against its schema, every file, and the code for what the sandbox would refuse (remote scripts, inline handlers, `eval`, `fetch`, `localStorage`, tokens in the ZIP, sites missing from `http.origins`, a panel that never says *ready*). Every finding has a code, the file, the line and a hint, so an AI can fix its own addon; `-addon-lint-json` gives the same as JSON, and a failed install in *Settings → Addons* opens the report by itself. `wrm -addon-new grafana-panel ./my-addon` starts from one of four templates (blank, fetch-and-render, grafana-panel, wasm-hook) that come with `wrm-bridge.js`, a small helper for the bridge. `wrm -addon-dev ./my-addon` opens the panel on your computer with the real sandbox and CSP, a stand-in for WRM (settings, web requests, secrets from environment variables) and a reload whenever a file changes. WRM also serves the manifest schema and an `llms.txt` guide for AI assistants. See [Addons](../README.md#addons) and the [addon developer guide](ADDONS.md#tools-for-writing-addons).

## What's new in v13.11.0

**Addons can reach web sites, through WRM and only where an administrator said yes.** An addon now lists the sites it needs in `addon.json` (for example `https://grafana.example.com`, with the methods and paths it uses). WRM makes the request on the panel's behalf, so the panel itself still has no network: anything not listed is refused, every name is resolved by WRM and private or local addresses are blocked unless the administrator ticked *private network* for that site, and redirects are checked again. API keys and passwords for those sites are **secrets** that you type into WRM (🔑 *Secrets* on the instance): per instance, with your own value in a shared instance if you need one. They are stored encrypted, added to the request by WRM and never shown to the panel. A personal addon's sites wait for an administrator's approval. Recent requests are listed in the window's diagnostics. See [Addons](../README.md#addons) and the [addon developer guide](ADDONS.md).

## What's new in v13.10.0

**Saved logins on demand, hostname names and a calmer top bar.** *Edit connection* has **👁 Show saved login**: after you confirm with your WRM password (or your domain password), WRM shows what the connection logs in with. Right-click a connection for **Copy username**, **Copy password** (same confirmation) and **Copy IP address**. A new connection can take its **name from the server's hostname**. Tabs show just the connection's name, the version next to the logo stays visible, the user name no longer jumps every 5 seconds, and the Git workspace server filter has no blank entries any more. See [Connections & folders](../README.md#connections--folders).

## What's new in v13.9.0

**Addon platform v2: your own addons, several windows each, and settings that stay.** Every user can now install **personal addons** that only they see (an administrator decides who may with `addon_personal` and can turn any of them off). An addon with a panel opens as **instances**, each with its own settings and window, for example two dashboards side by side; administrators can set up shared instances for everybody. Settings are kept **on the WRM server** with export and import, and windows come back where you left them after a reload or on another device. Addon **themes** appear in *Settings → Appearance* for whoever picks them. Before an addon is enabled WRM shows a **permission sheet** generated from `addon.json`; an update that asks for more waits for approval again. apiVersion 2 panels run without any network of their own, every file is pinned by SHA-256, and a watchdog offers *Reload* when a panel hangs. See [Addons](../README.md#addons) and the [addon developer guide](ADDONS.md).

## What's new in v13.8.0

**AI sessions, your own AI settings and approvals that do what they say.** The new **🧠 AI sessions** panel lists every AI session you had, from the assistant, the log and file assistants, OpenCode agents and AI apps over MCP, with search (also in the conversation), filters, a transcript viewer, **Continue** (also an MCP session, whose steps become the earlier context), **Save** (kept beyond the transcript retention) and export to Markdown or JSON. Everybody sees only their own sessions. *⚙ My AI settings* in the assistant keep your own default provider (your own key, e.g. your OpenAI or Anthropic account), model, mode, engine and whether you see the model's thinking. Over MCP an approval given in WRM after the AI app stopped waiting now **runs the call at once** (`ai_mcp_later_approval`), *Deny* in the AI app is final, an app's open question is closed when you decide in WRM, and deciding a request that is no longer waiting says who decided it and where. The work log keeps the whole thinking text instead of its last lines. See [AI sessions](../README.md#ai-sessions) and [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v13.7.0

**Themes, accessibility and a screen lock you can rely on.** *Settings → Appearance* offers eleven themes (dark, OLED black, Fjord, Solar, Forest, Ember, light Daylight, Paper, Solar light and two high-contrast themes), follows the system's light/dark setting if you want, and lets you change every colour, the fonts, the interface size, corners and spacing, save your own themes, export and import them, and add your own CSS. Themes from addons are picked there per user instead of being forced on everyone. For colour vision deficiency (deuteranopia, protanopia, tritanopia, achromatopsia) WRM switches status colours, diffs and optionally the terminal palette to colours you can tell apart, can add ✓ ✕ ! signs and shapes so nothing relies on colour alone, and offers more contrast, a terminal minimum contrast, reduced motion and transparency and a strong focus outline. *Settings → Security → Screen lock* sets when WRM locks (inactivity with a warning, a hidden window, the computer locking or sleeping, on every start, a keyboard shortcut) and how the lock screen looks; **Lock now** locks at once while terminals, transfers, tunnels, AI and MCP keep running. The lock is now kept by the server: a reload or a new tab stays locked, directory (LDAP) accounts can unlock, and a personal PIN can unlock instead of the password. See [Appearance, themes and accessibility](../README.md#appearance-themes-and-accessibility) and [Screen lock](../README.md#screen-lock).

## What's new in v13.6.0

**Log folders, a folder tree and smarter operators.** *Monitoring → Logs* can show the logs as a **folder tree** under their main folders (`/var/log`, `/home/scripts` …) and has a **Folder** filter for one folder and everything below it; *Log settings* take your own **main folders**. Operators are recognised from words that repeat in folder and service names (`nova`, `NOVA2`, `novatv` → one operator group), and every log says whether it sits in a log folder or an app folder. A log folder and a script folder with the same name (different case) no longer hide each other's logs. See [Logs: viewer, search and log settings](../README.md#logs-viewer-search-and-log-settings).

## What's new in v13.5.0

**A professional log viewer, a smart editor and an AI that shows its work.** *Monitoring → Logs* lists every log of a server in about a second (one read-only SSH command), sortable and filterable by operator and service, and never lists scripts such as `.py`. Logs open in tabs with level colours, smart highlighting, find, level filters and **Follow**; the **Search** tab greps all selected logs on the server for a From/To range in the EU date format and sorts the hits on a timeline. *Log settings* are now separate from Zabbix. The file manager's editor highlights syntax, finds and replaces, shows your changes before saving, warns when the file changed on the server and has an **AI assistant** that reviews and edits the open file with your approval. The log viewer has an AI log assistant (`search_logs`). The AI work log shows the assistant's thinking and each step, chats get a short title with date and time, and an ended chat can be **continued**. See [Logs: viewer, search and log settings](../README.md#logs-viewer-search-and-log-settings), [Text editor](../README.md#text-editor) and [AI assistant](../README.md#ai-assistant).

v13.0.0 to v13.4.4 (log engine, addons, screen lock, OpenCode Go, Claude quick start, Claude Desktop and Claude Code setup) are described in the [changelog](../CHANGELOG.md).

## What's new in v12.9.1

**Medium runs without the original tool.** The WRM engine is now the default for projects nobody chose an engine for, its card comes first on the Medium page, and the original deploy tool with its Docker runtime is folded away as optional. Upgraded installations that already use the original tool keep it as their default until an administrator switches **Medium: default engine** to `native`. The deploy-config editor can now pull the configuration from its Git repository, show your local changes and **Commit and push** them through the WRM Git source, after checking that the whole configuration loads. WRM-host builds can run as an unprivileged account (`medium_build_user`) with a clean environment.

## What's new in v12.9.0

**Medium shows what changes in each file and lets you edit deploy-config.** In the WRM engine's plan, **Details** opens the server's version of a changed file against the branch, line by line. Files that differ only in CRLF/LF are marked *line endings only* and counted separately, as in the original tool. *Edit deploy-config* edits `servers.yaml`, `environments.yaml` and the project files in WRM: a file is saved only when the whole configuration still loads, and comments, layout and line endings stay as you wrote them.

## What's new in v12.8.1

**Security update without a database change.** The build toolchain is now Go 1.27.3, which fixes **GO-2026-6809 / CVE-2026-111943**: `crypto/x509` now rejects known compromised intermediate certificates that an older system trust store may still contain. The gosec G702 finding in the OpenCode Docker sandbox was reviewed and accepted as a false positive: the Docker path is operator-configured and the container name is generated by WRM, both passed as separate arguments without a shell. See the [release notes](../RELEASE_NOTES.md) for the provenance of both findings.

## What's new in v12.8.0

**Folders now behave like a proper remote manager.** Deleting a folder removes its subfolders and every connection inside in one step, without moving anything to the root first. Folders and connections can be moved from the context menu or dragged, including several selected connections at once, and the new position is saved. Shared and collaboration trees keep their nested folder structure; connection names sort case-insensitively, while folders keep their manual order.

**Every reported detail is covered.** The Settings window opens at the size of the Admin window and wide tab rows have scroll arrows. Terminal windows have **Edit** and **Tunnels** buttons in the title bar. The file manager can change permissions (including recursively and with an optional sudo password), and the transfer destination can create a folder. Additional uploads can be queued while one is running. Tunnels can be filtered by port, description or host, with an active-only switch, and proxy settings are reachable from **Tools** and the Tunnels panel. Save buttons show progress and reject a second click, and **+ Connection** on a folder is pre-filled with that folder and its defaults.

## What's new in v12.7.4

**RDP works on Windows without preparing guacd by hand.** When guacd does not answer and its address is local (`127.0.0.1:4822`), WRM starts the official `guacamole/guacd:1.5.5` container in Docker Desktop, publishes it on loopback only, and keeps it running with Docker (`wrm-guacd`, `unless-stopped`). The first RDP/VNC/Telnet connection can start it; later starts happen with WRM. *Admin → Overview* shows the state and has **Start guacd**, `-setup-guacd` does the same from the command line, and the desktop window explains when Docker Desktop is missing or stopped. A custom guacd on another machine is never touched. Linux keeps using its system guacd or the `docker-compose.yml` service.

**Every text file previews with its formatting.** Preview now detects text from the file content, so shell/Python scripts, XML/SMIL, JSON, YAML, configuration files and extensionless files open as text instead of being refused. A bundled offline PrismJS highlighter colors the formats it knows; binary files still get the clear answer to download them. The 256 KiB preview limit and escaped-text safety policy are unchanged.

## What's new in v12.7.3

**Medium builds where you choose.** Projects with a build step can now deploy with the **WRM native** engine. For each project, choose where the build runs: on the WRM host, on a chosen Linux build server over SSH, or not at all. The plan shows the place and the build's platform and architecture before anything runs; a destination with a different platform or architecture is refused with the reason. A build server receives only the prepared source in a temporary directory (no Git metadata and no SSH keys), and the directory is removed afterwards.

**Monitoring → Logs finds every service.** Service logs are now discovered under `/home/<operator>/<operator>/<service>/log` and `/logs`, with the correct Service and Operator labels. Non-compressed service logs, including FTP `.txt` files, are listed; virtual environments, build/archive directories and old backups are skipped. Discovery and listing have larger time budgets, so large hosts no longer show incomplete results.

## What's new in v12.7.2

**The version is read from the code.** The Git workspace no longer trusts `VERSION.md` to say what is deployed. It compares every tracked file on the server with the repository tree at the target and at the earlier tags: all files matching one version show that version, a consistent older version is named exactly from its commit tree, and a mix or a hand-edited file shows **Mixed / modified** with the per-file details. A new checkbox next to *Check now* (off by default) also checks `VERSION.md` and warns when it disagrees. See [Git workspace: every file and partial checks](../README.md#git-workspace-every-file-and-partial-checks).

## What's new in v12.7.1

**Medium: per-environment server lists.** `environments.yaml` now accepts the original deploy tool's `servers: [name, ...]` key, so an environment can be limited to chosen servers. A project target in that environment cannot deploy to any other server; an unknown, duplicate or empty list is refused with the file and key named. Existing configs without the key keep every server. See [Medium native engine](../README.md#medium-native-engine).

## What's new in v12.7.0

**Medium deploys natively.** Switch a Medium project to the **WRM native** engine and WRM deploys it itself, from your WRM Git source to every server of the environment through your saved WRM connections, without Docker and without the original tool. Each server gets a lock, a backup of everything that changes, a SHA-256 check of every uploaded file, `post_deploy` and a history entry, all in the original tool's format, so both engines can take turns on the same servers and roll back each other's deployments. A **dry run** compares again under the lock and changes nothing; changes can be **left out** of a deployment; **History and rollback** shows what runs where, rolls back the last deployment and releases a lock someone left behind. Projects with a build step gained native builds in v12.7.3. See [Medium native engine](../README.md#medium-native-engine).

## What's new in v12.6.26

**SIEM export, security alerts and lockdown.** WRM sends every audit entry to your **SIEM**: syslog over TLS (RFC 5424), a JSON lines file, Splunk HEC, Elastic, Microsoft Sentinel or any HTTPS collector, with user, role, directory identity, source address and the audit hash chain, and catches up after an outage. The admin sees per destination whether it delivers and what is waiting. **Security alerts** tell administrators about repeated failed sign-ins, break-glass sign-ins, new administrators, revealed secrets, AI full-access approvals and the audit being turned off. **Stop user** ends everything an account has open in one click, and **lockdown** keeps everyone except break-glass administrators out, switches AI off and pauses shares, from the UI, `wrm -lockdown` or `WRM_LOCKDOWN=1`. See [SIEM export, alerts and lockdown](../README.md#siem-export-alerts-and-lockdown).

## What's new in v12.6.25

**Active Directory, LDAP and SCIM.** People sign in with their **Active Directory or LDAP** name and password on the normal sign-in page; WRM checks the password with the directory over LDAPS or StartTLS and never stores it. Rules turn directory groups (nested AD groups included) into WRM roles and teams, and WRM **checks every linked account regularly**: someone deleted, disabled, locked or removed from the WRM groups loses access at once. **SCIM 2.0** lets Microsoft Entra ID or Okta create, change and disable WRM accounts and groups immediately. Whichever way an account is stopped, WRM ends its sessions, terminals, tunnels, rooms, AI sessions and AI connections. See [LDAP / Active Directory and SCIM](../README.md#ldap--active-directory-and-scim).

## What's new in v12.6.24

**Sign in with your organisation's account.** WRM now signs people in through **Microsoft Entra ID, ADFS, Okta, Keycloak, Google Workspace or any OpenID Connect provider** (authorization code with PKCE). In *Admin panel → Directory sign-in* you add the provider, test the connection step by step, and write rules: an Entra **app role** or a **group ID** gives a WRM role and teams. Accounts are created on the first sign-in only for people a rule matches; everyone else is refused with a message that says why. People are recognised by their immutable directory ID, never by their name, so a renamed or re-used name cannot take over an account. A provider can require **MFA proven by the token**. WRM **checks regularly** with the provider that each signed-in user may still sign in: someone disabled or removed from the WRM group loses their sessions and terminals. Password sign-in can be limited to the **break-glass** administrators (the Enterprise level requires it once a directory sign-in is on). **Test sign-in** shows what the provider sends and what WRM would do, without creating anything. See [Directory sign-in (OIDC)](../README.md#directory-sign-in-oidc).

## What's new in v12.6.23

**Roles, teams and a second confirmation for risky actions.** Every account now has a role: *Administrator*, *Security auditor* (reads the audit log, terminal sessions, recordings, file transfers and the security posture, and changes nothing), *Operator* or *User*. Administrators create **teams** in *Admin panel → Teams*; a share can include a whole team with one role, and policies that say who may use a feature accept `operators` and `teams:<names>` besides `off / admins / all`. With **step-up** on (required from the Business level), changing policies, roles, the security level, deleting or disabling users, resetting passwords or 2FA and revealing stored secrets ask for your password (and 2FA code) again. A **break-glass** administrator account must use 2FA, can be limited to chosen networks, and every sign-in is audited and sent to the other administrators. When you delete a user, their connections, folders, credentials, keys, shares and proxies can **pass to their team lead or another account** instead of being deleted. See [Roles, teams, step-up and break-glass](../README.md#roles-teams-step-up-and-break-glass).

## What's new in v12.6.22

**Security levels: Home, Business, Enterprise.** Pick how strict WRM is on the first-run screen, at install (`-security-level`, `WRM_INITIAL_SECURITY_LEVEL`, forced with `WRM_SECURITY_LEVEL`) or later in **Admin panel → Security**. A level sets minimums for HTTPS, two-factor authentication, passwords, sessions, registration, host keys, secret export, guest links, AI full access, audit and session recording. Before you switch, a preview shows every change, what blocks it and whom it affects; a setting you keep below the level is a **deviation** with a written reason, on record and audited. The **security posture** page grades the installation item by item with the exact fix, and exports a report for auditors (JSON or PDF). See [Security levels and posture](../README.md#security-levels-and-posture).

## What's new in v12.6.21

**Signed releases.** Every release from now on carries an ed25519 signature of its `SHA256SUMS.txt`, a software bill of materials (CycloneDX) and GitHub build provenance. Once the project's public key is built in, the one-click update refuses a release whose signature is missing or wrong, so a changed file on the release page cannot be installed. `wrm -release-verify SHA256SUMS.txt` checks a download by hand. Every night, fuzz tests feed random input to the import parsers (mRemoteNG, SSH config, CSV/XLSX, PuTTY, WinSCP, configuration files, deploy-config YAML) and OWASP ZAP scans a running WRM, signed out and as an administrator. The fuzz tests already found and fixed a crash in the PuTTY import.

## What's new in v12.6.20

**Security checked on every change.** WRM is now built with a supported Go release (1.27) and up-to-date libraries, and every push and pull request runs automatic security checks: known vulnerabilities in Go and the libraries WRM uses (govulncheck), GitHub CodeQL for the server and the web UI, gosec and staticcheck, secret scanning (gitleaks), checks of the build workflows themselves, a dependency review, a vulnerability scan of the Docker image and the OpenSSF Scorecard. Dependabot proposes library updates every week, and Claude can review pull requests for security problems. The Medium proxy uses Go's safer `Rewrite` hook. No change for users; see [Continuous security checks](../SECURITY.md#continuous-security-checks) and the [security roadmap](../ROADMAP.md) (signed releases, security levels for home, business and enterprise, roles and teams, sign-in with Microsoft Entra ID / AD, LDAP and SCIM, SIEM export).

## What's new in v12.6.11

**OpenCode inside WRM.** Add the **OpenCode Zen** provider and use OpenCode's free models in the AI assistant without an account, or paste an OpenCode account's Zen key for all its models. A **Go subscription** has its own **OpenCode Go** provider, so a Go key is sent to Go's endpoints instead of Zen. Choose **OpenCode agent** as the assistant's engine (agents *build* and *plan*): WRM downloads OpenCode itself after one click by an administrator (sha512 checked) and runs it on the WRM host with its own shell, file and web tools switched off, so it works on your servers only through WRM and every command still waits for your approval. The assistant can also hand text jobs to **helper models**, picked automatically by the kind of job (free first), each labelled with how it treats your data; models that may train on prompts stay off until you switch them on. See [OpenCode Zen, Go, agents and helper models](../README.md#opencode-zen-go-agents-and-helper-models).

## What's new in v12.6.10

**One canvas, arranged your way.** The window area is one canvas: terminals, file managers, desktops, Workboard addons and panels share it. Addons are windows you move, resize, snap, pin or dock; panels can be placed anywhere on the canvas, not only at the edges. The **▤ Workboard** menu adds addons and holds the canvas tools: grid, tile and cascade, lock, saved layouts (optionally with your terminals) and export / import. New addons: links, saved sessions and a timer. The top bar is tidier: *Broadcast* is in **🧰 Tools**, *Shares* is an icon, and **Settings · Panels · Workboard** sit together. A closed terminal's banner has a **✕ Close window** button. See [Workboard: the shared canvas and addons](../README.md#workboard-the-shared-canvas-and-addons).

## What's new in v12.6.9

**Workboard: what matters, at a glance.** **▤ Workboard** in the top bar opens a page of tiles you put together yourself from an addon catalog, so nothing appears until you add it: server status, CPU / RAM / disk bars of chosen servers (from Monitoring, or read over SSH), WRM's notifications inside WRM, clocks for your time zones, the last commits of a Git service, running tunnels, quick connect buttons and notes. Any tile can be pulled out into a panel and docked, stacked or floated like the other panels. See [Workboard and addons](../README.md#workboard-the-shared-canvas-and-addons).

## What's new in v12.6.8

**Your panels, where you want them.** *Connections*, *Sessions*, *Snippets*, *Clipboard* and the AI live view are panels with a tab you can drag: onto another panel's tab bar to stack them as tabs, onto the top or bottom part of a panel to split the column, onto an edge strip to dock it left, right or under the windows, or anywhere else to let it float. **Panels** (▦) in the top bar shows and hides them, and **Lock the layout** keeps everything in place. The layout is stored for your account on the server, so it follows you to every browser. See [Panel layout](../README.md#panel-layout-dock-stack-float-and-lock).

## What's new in v12.6.7

**A top bar that fits and an AI panel that docks.** The top bar now fits at any width and browser zoom: labels go first, then less used buttons move into a **⋯** menu, so nothing overlaps. The AI live view (what an AI app runs on your servers, with approvals and output) no longer covers your windows: it is a **🤖 AI** tab next to *Sessions* and *Snippets*, a dock under the windows, a floating window, or hidden, chosen with its title bar buttons or by dragging it to an edge, and remembered per user. Approval cards still always pop up. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp). Every text on the page can be selected and copied again, and Medium's Docker runtime trusts your computer's certificates, so it works behind proxies and antivirus that inspect HTTPS (see [Docker for Medium](../README.md#docker-for-medium-installed-by-wrm)).

## What's new in v12.6.6

**Docker everywhere, installed by WRM.** The original deploy tool now runs in WRM's managed Docker runtime on Linux, macOS and Windows alike, and WRM installs Docker itself when it is missing: an administrator presses **Install Docker** in Medium step 2 once (accepting Docker Desktop's licence on Windows and macOS), and WRM installs Docker Engine from the Linux distribution's packages or Docker Desktop from docker.com (signature checked), starts it, waits for it and prepares the runtime, every step with its output in the activity log. Afterwards a tool ZIP import does all of it by itself. Where WRM has no administrator rights it shows the one command to run, `sudo wrm -setup-docker`. See [Docker for Medium](../README.md#docker-for-medium-installed-by-wrm).

## What's new in v12.6.5

**Medium tells you what it is doing.** The Medium view now has a numbered **step-by-step guide** for the original deploy tool (1 import the tool ZIP, 2 prepare its runtime, 3 your deploy-config, 4 your SSH profile, 5 start/stop), each with a light, what it waits for and what to do when it fails. While something runs you see **what runs, who started it, since when and every step with its time**, with the output of Docker, pip and the original tool (secrets redacted). A deploy-config ZIP is checked by WRM at once and kept for the next start while the runtime is still being prepared, so you no longer wait to upload it. A new **activity log** lists every Medium action with its result, duration, steps and output, and WRM's console logs each step. The read-only WRM preview is a separate, clearly labelled panel. See [Medium step by step](../README.md#medium-step-by-step-status-and-activity-log).

## What's new in v12.6.4

**Your blocked commands, your choice.** With full access on, the administrator now decides whether their own list of blocked commands (`ai_blocked_commands`) asks for approval like everything else (default) or stays blocked even with full access (`ai_full_access_blocked`, *Admin panel → AI connections*). See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v12.6.3

**Full access with one switch.** The administrator now chooses how full access applies (`ai_full_access_mode`): by default the **Full access for the AI** switch alone opens it for every AI session; it can instead apply only to AI connections with the **⚠ Full access** scope, or to users who turn it on for themselves in *Settings → AI connections*. Claude still asks before anything WRM would otherwise refuse. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v12.6.2

**Full access, but Claude still asks.** When an administrator allows it (`ai_full_access`), an AI connection can get the **⚠ Full access** scope. WRM then refuses nothing outright: destructive and blocked commands, deny lists, shutdown and secret files come to you as an approval marked *FULL ACCESS*, every time and never automatically. Read-only connections stay read-only. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v12.6.1

**"Allow everything" now means everything.** In an AI connection's automatic mode, **Allow everything** fills in `*` for commands and `/**` for files, and **No time or action limit** keeps the connection automatic for as long as it is valid. Chains with `&&`, variables and redirections no longer stop at an approval when `*` is allowed; blocked commands, your deny lists and sensitive files still do. When WRM does ask and nobody answers in time (or the app declines the question), the request stays in WRM to **approve later**, and the AI repeats it once you have. Files that may hold secrets are read after your approval, with secrets redacted. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v12.6.0

**Medium's native engine, first step.** The Medium view has a read-only *Native engine* panel: WRM reads your `deploy-config` as it is, matches its servers to your saved WRM connections, prepares the chosen branch or tag from your WRM Git source and shows what a deployment would change on each server, without Docker and without writing anything. **Compare with the original tool** runs the original comparison next to it and shows whether both give the same result. Deploying natively and rollback followed in v12.7.0. See [Medium native engine](../README.md#medium-native-engine).

## What's new in v12.5.1

**Claude Desktop and Claude Code in one step.** *Settings → AI connections → Claude Desktop* now downloads a **WRM PRO extension** (`wrm-pro.mcpb`) with your WRM address and token already filled in: open it, press **Install**, done. It works with a WRM on your own network and over plain HTTP, because Claude Desktop talks to WRM from your computer (a claude.ai custom connector connects from Anthropic's cloud and needs a public HTTPS address). For Claude Code there is a **WRM PRO plugin**: `claude plugin install wrm-pro --marketplace vedranius/web-browser-RDM-public --config wrm_url=… --config wrm_token=…`. Token lifetimes are entered in hours or days and default to 7 days (at most 30), so you are not asked to approve again every day. While an AI app works on a server, a **live panel** in WRM shows each command, the approval buttons and the output; after a denial the AI asks you how to go on. See [AI desktop apps (MCP)](../README.md#ai-desktop-apps-mcp).

## What's new in v12.5.0

**Stamp as it is:** the Stamp dialog shows which files differ and why an installation cannot be stamped. Tick **Also stamp installations that are not up to date** to write VERSION.md with each file's actual version and a note that the installation is partial.

**Back in dialogs:** dialogs opened from another dialog (for example Details → Update → confirmation) have **← Back** (Alt+←) and keep what you ticked or typed. A code review of the v12.4.x features fixed production confirmations, Medium start/stop/settings and imports, log search with non-UTF-8 bytes and special characters, and the client GitLab cache. The next versions run Medium without Docker on WRM's own engine ([roadmap](../ROADMAP.md)).

## What's new in v12.4.5

**Medium:** select the trusted tool/config ZIPs and click **Prepare and run**. WRM retains the imports, starts Docker, prepares the original dependencies, validates config, synchronizes the private profile and opens the unchanged original tool. First Windows setup can invoke the signed official Docker installer; OS/vendor prompts and required restarts remain visible. Existing custom configuration Git settings are replaced by the explicitly uploaded local config ZIP.

**Monitoring → Logs:** discovery starts automatically and streams results by server. Service log folders and Supervisor paths supplement manual roots. Filter by server/operator/service and event time, then search a keyword across all matching files from their beginning. Results show file modification time, matching event range, highlighted excerpts and open/follow actions. Progress, permission errors and partial searches remain visible; continue a bounded scan when necessary. Host metrics, service checks and SSH logs have separate setup requirements.

Overview/Installations keep filters, check controls and actions together while scrolling. Environments adds combined search/service/environment/server filters.

[Operation and limits](V12.4.5_GUIDE.md) · [Validation](RELEASE_SUMMARY_v12.4.5.md).

## What's new in v12.4.4

Monitoring with Zabbix metrics/history and service Functionality, multi-window logs, client GitLab token access with persistent snapshots, easier MCP setup and editable grants, structured configuration inventory, terminal selection-to-AI/Open terminal here, visual folder/connection labels, longer provider timeouts, sticky filters and replay/import fixes.

[Operation and limits](V12.4.4_GUIDE.md) · [Release validation](RELEASE_SUMMARY_v12.4.4.md) · [Original Medium tool parity](MEDIUM_PARITY_v12.4.3.md). General configuration import now recognizes mRemoteNG XML and PuTTY registry files ([#5](https://github.com/vedranius/web-browser-RDM-public/issues/5)).
