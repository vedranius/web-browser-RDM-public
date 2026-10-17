## Web Remote Manager PRO v12.6.19 — OpenCode with its own tools, in a sandbox

### Added

- **Both ways to run OpenCode:** next to the contained OpenCode agent (its own tools off, everything through WRM), an administrator can now allow **🐳 OpenCode's own tools, in a Docker sandbox** (`ai_opencode_sandbox`: off by default, administrators only, or everyone). The agent then gets its own shell, file and web tools, but only inside a container on the WRM host. Your servers are still reached only through WRM.
- **Every use asks you:** each shell command, file change and web fetch in the sandbox appears as an approval card with the exact command, diff or address. You approve it once or refuse it, with a note if you like; it cannot be edited and is never approved automatically, not even with full access.
- **A locked-down container:** no capabilities, no privilege escalation, read-only system, memory, CPU and process limits, only the session's own folder mounted, the port open on `127.0.0.1` only. It is removed when the session ends or WRM stops. The image is your choice (`ai_opencode_sandbox_image`, default `debian:stable-slim`).
- **Set up with the same one click:** *Download and install* also prepares the sandbox (the Linux build of OpenCode for Docker, the image, a test run). The admin page shows whether the sandbox is ready and what is missing.

### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Nothing changes until an administrator sets `ai_opencode_sandbox` and presses *Download and install* again under *Admin panel → AI assistant → OpenCode agents*. The sandbox needs Docker on the WRM host (Medium's step 2 can install it) and WRM itself not running in a container.

[Download v12.6.19](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.19) · [OpenCode in WRM](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.19/README.md#opencode-zen-agents-and-helper-models)

### Included from v12.6.11 — OpenCode inside WRM

#### Added

- **OpenCode Zen, free models without a key:** add the **OpenCode Zen** provider (a preset in the wizard) and use OpenCode's free models in the AI assistant without an account. With an OpenCode account, paste its Zen API key (org-wide by an administrator, or for yourself under *My API keys*) and every model of the account is available. Free models are marked in the model lists.
- **OpenCode agents in the assistant:** choose **OpenCode agent** as the engine and its **build** or **plan** agent. WRM downloads OpenCode itself (an administrator presses *Download and install* once; the file is checked against the registry's sha512) and runs it on the WRM host, one private process per session. OpenCode's own shell, file and web tools are switched off: it works on your server **only through WRM**, so every command and edit shows up in the chat and waits for your approval exactly like the WRM assistant (unless full access is on).
- **Helper models, picked automatically:** the assistant can hand a self-contained text job (summarise a long log, explain an error, draft a snippet) to a cheaper or free helper model. WRM chooses the model by the kind of job, free ones first, and tries the next one if a model fails. Each delegation is a card in the chat with the model, size, tokens, cost and how that model treats your data.
- **Clear data-use labels:** 🛡 zero retention, ⚠ may use prompts for training, ❔ unknown terms, 💳 paid. Models that may train on prompts are **off by default**; each user can switch them all on for themselves in the assistant (administrators can also allow them for everyone or never).
- **OpenCode on your computer** can connect to WRM too: the MCP wizard has an OpenCode entry with the `opencode.json` snippet or the OAuth login.

#### Fixed

- Stopping WRM also stops every running OpenCode agent; after a crash WRM cleans up what was left on its next start.

#### Upgrade

Replace the executable/container and keep the database and encryption key. The database gets one new column (`users.ai_delegate_training`); a previous version still starts on it. Nothing changes until an administrator adds an OpenCode Zen provider, turns on delegation (`ai_delegate`) or installs OpenCode (*Admin panel → AI assistant → OpenCode agents*). New settings: `ai_delegate*`, `ai_opencode*`. New endpoints: `/api/ai/delegate`, `/api/admin/opencode`.

[Download v12.6.11](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.11) · [OpenCode in WRM](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.11/README.md#opencode-zen-agents-and-helper-models)

### Included from v12.6.10 — One canvas, arranged your way

#### Added

- **One canvas:** the Workboard is no longer a page over your windows. Everything in the window area shares one canvas: terminals, file managers, desktops, Workboard addons and panels you place on it. Put a server-load addon next to the terminal of that server, a clock in a corner, notes under a file manager, and it stays like that.
- **Addons are windows:** move them, resize them, snap them to halves and quarters, minimize, maximize, switch with the tabs, or dock one as a panel at the left, right or bottom. They come back exactly where you left them, in every browser.
- **Panels anywhere:** drag a panel's tab anywhere on the canvas (or pick *On the canvas (anywhere)* in its ⋯ menu) and it stays there between your windows, not only in the docks at the edges.
- **The ▤ Workboard menu** holds the canvas tools: add addons, **snap to a grid** (10, 20 or 40 px, with dots or lines on the canvas), **tile** or **cascade** all windows, minimize and show them, show or hide the addons, **lock** the whole layout, **saved layouts** (for example *Morning check* and *Deploy day*, optionally with the open terminals), and **export / import** a layout as a file to share it or move it to another WRM.
- **Pin and always on top:** right-click any window's title bar or tab to pin it in place or keep it above the others.
- **New addons:** *Links* to your dashboards and wikis, *Saved sessions* to open one with a click, and a *Timer* (countdown with a notification, or stopwatch). Every addon can get a colour and a text size.
- **Close from the reconnect banner:** when a terminal or desktop connection closes or fails, the banner has a red **✕ Close window** button next to **↻ Reconnect** (#70).

#### Changed

- **Top bar:** *Broadcast* is in the **🧰 Tools** menu with the network tools (the button turns orange while a broadcast is on); *Shares* is only its icon; **Settings, Panels and Workboard** sit together right of *Keys*.
- **Lock the layout** now also keeps windows from being moved or resized.
- Folder bookmarks in the file manager (#17) have been in WRM since v10.10: ★ next to the path bar, per connection, folder, tag or for all connections, with WinSCP import.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Your v12.6.9 Workboard tiles appear as addon windows beside your terminals the first time; panels you pulled out stay where they are. The per-user interface store accepts up to 256 KiB.

[Download v12.6.10](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.10) · [The canvas](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.10/README.md#workboard-the-shared-canvas-and-addons)

### Included from v12.6.9 — Workboard: what matters, at a glance

#### Added

- **Workboard** (▤ in the top bar): a page of tiles over your windows, empty until you add something. Pick addons from the catalog and arrange them: drag a tile by its title to move it, make it wide or narrow, give it your own title, remove it again. Your Workboard follows you to every browser.
- **Addons to start with:**
  - **Server status:** which servers are down, slow or up, with latency and since when; problems first, or only problems. Click a server to open it.
  - **Server load:** CPU, RAM, load, uptime and disk usage of the servers you choose, as bars that turn yellow at 75 % and red at 90 %. Values come from Monitoring when the server has it, otherwise WRM reads them over SSH (Linux) once a minute while the tile is visible. Each server says where its values come from and how old they are; errors are shown as they are.
  - **Notifications:** everything WRM told you (servers down and back, expiring keys, new versions, Git jobs) in WRM itself, also when no e-mail or chat channel is set up. Unread ones are marked and counted on the ▤ button.
  - **Clock:** your local time and the time zones you work with.
  - **Recent commits:** the last 10 (or 5, 20) commits and pushes of a service from your Git workspace (GitLab, GitHub or Gitea), with author, branch and a link.
  - **Running tunnels**, **Quick connect** buttons for your most used servers, and **Notes**.
- **Pull a tile out:** **⤢** turns a tile into a panel you can dock left, right or at the bottom, stack with *Sessions*, or let float, like every other panel. It stays there after a reload; *Back to the Workboard* returns it. A locked layout keeps tiles in place too.
- **Administrators** decide whether server load may be read over SSH: *Workboard: read server load over SSH* in *Admin panel → Policies → Live status* (on by default).

#### Changed

- A tab bar that is too narrow for all its tabs scrolls to the open one.

#### Upgrade

Replace the executable/container and keep the database and encryption key. The database gets one new table (`user_events`); a previous version still starts on it. Nothing new appears until you open the Workboard and add an addon. New endpoints: `/api/me/events`, `/api/workboard/metrics`.

[Download v12.6.9](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.9) · [Workboard](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.9/README.md#workboard-and-addons)

### Included from v12.6.8 — Your panels, where you want them

#### Added

- **Arrange the page yourself.** *Connections*, *Sessions*, *Snippets*, *Clipboard* and the AI live view are now panels with a tab. Drag a panel by its tab:
  - onto another panel's tab bar, and they share it as tabs;
  - onto the top or bottom part of a panel, and the column is split between them;
  - onto the *Left*, *Right* or *Bottom* strip that appears at the edges while you drag, and it docks there (the bottom zone sits under your windows);
  - anywhere else, and it becomes a floating window you move by its tab bar and resize by its corner.
  Drag the line between two panels, or a zone's edge, to resize. Every panel's **⋯** menu also moves it, hides it, locks the layout or resets it.
- **Panels** (▦ in the top bar): every panel with where it is; show or hide it, lock the layout, or go back to the start layout.
- **Lock the layout** so nothing moves by accident: no dragging, splitting or resizing while locked (the button shows 🔒). Tabs, collapsing the side panels, showing and hiding still work.
- **Your layout follows you:** it is stored for your account on the WRM server, so every browser and computer you sign in on shows the same arrangement.

#### Changed

- The AI live view is a panel like the others: dock it left, right or at the bottom, stack it with *Sessions*, let it float or hide it. Its title bar buttons and the 🤖 menu still work.
- On a phone the page keeps its fixed layout with the two drawers.

#### Upgrade

Replace the executable/container and keep the database and encryption key. The database gets one new table (`user_ui`); a previous version still starts on it. Your first layout is built from what this browser remembered in v12.6.7 (the right panel's tab and width, where the AI live view was).

[Download v12.6.8](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.8) · [Panel layout](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.8/README.md#panel-layout-dock-stack-float-and-lock)

### Included from v12.6.7 — A top bar that fits, an AI panel that docks

#### Added

- **The AI live view docks where you want it.** What an AI app runs on your servers (commands, files, WRM's decision, approvals and output) is no longer a box over your windows. Choose with the buttons in its title bar, in the 🤖 menu, or by dragging the title bar onto the right or bottom edge:
  - **Right** (default): a **🤖 AI** tab next to *Sessions* and *Snippets*, with its own width; the tab counts the open AI sessions and turns yellow while one waits for you.
  - **Bottom**: a dock under the windows, resizable by its top edge; several AI sessions sit side by side.
  - **Floating**: a window you move by its title bar and resize by its corner.
  - **Hidden**: a small *Show the AI panel* chip at the bottom left, which pulses while an approval waits.
- WRM remembers the layout, the floating position and size and the dock height for each user in this browser.
- Approval cards still always pop up, whatever the layout, so no request is missed. On a phone the live view is a sheet at the bottom of the screen.
- **Medium works behind proxies and antivirus that inspect HTTPS.** The original deploy tool's Docker runtime now trusts the certificates your WRM computer trusts (the Windows certificate store, the macOS keychains or the Linux CA bundle), so pip no longer fails with *self-signed certificate in certificate chain*. New policy *Medium trusts this computer's certificates* (on by default) and *Extra CA file for Medium* for a proxy root that is not in the computer's store, in *Administration → Policies → Git*.

#### Changed

- **The top bar fits at any width and browser zoom.** It measures itself and first drops the labels of less used buttons; buttons that still do not fit move into a **⋯** menu (with their badges), so no button covers another at 1100–1920 px or at 125–150 % zoom.
- The right panel remembers its width, separately for the AI tab.

#### Fixed

- **You can select and copy every text on the page again**, for example Medium's command output, error messages and hints. Only buttons, tabs, title bars and drag handles stay unselectable.
- The *Snippets* tab no longer looks selected while *Sessions* is open.
- When pip fails because of an untrusted certificate, Medium says so and what to do.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database or API change. If Medium's runtime preparation failed with a certificate error, press *Prepare runtime* again after the upgrade.

[Download v12.6.7](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.7) · [AI live view layout](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.7/README.md#ai-desktop-apps-mcp)

### Included from v12.6.6 — Docker everywhere, installed by WRM

#### Added

- **Medium installs Docker itself.** When Docker is missing on the WRM host, step 2 of the Medium guide says which Docker WRM will install and how. An administrator presses **Install Docker** once; on Windows and macOS that also accepts Docker Desktop's licence (shown with its terms). From then on importing the tool ZIP or pressing *Prepare runtime* installs Docker, starts it, waits for it and prepares the runtime, with no manual steps.
  - **Linux:** Docker Engine from your distribution's packages, its service enabled and started.
  - **macOS:** Docker Desktop from docker.com, signature checked, installed for the user signed in on the Mac.
  - **Windows:** Docker Desktop from docker.com, signature checked, installed quietly for the signed-in user. The first time Windows may ask once to enable WSL and to restart.
- **Every step is visible:** installing, downloading, checking the signature, starting and waiting for Docker, each with its time, the command output and what to do when it fails, also in the activity log.
- **One command when WRM has no rights:** if WRM runs as a normal Linux account, step 2 shows `sudo wrm -setup-docker -service-user <account>`. It installs Docker, lets WRM use it and restarts WRM.
- **New policy** *Install Docker for Medium* (`medium_docker_install`, on by default) and *Medium command output* in *Administration → Policies → Git*.

#### Changed

- The original deploy tool now runs in the same Docker runtime on Linux, macOS and Windows. A host runtime you already prepared on Linux or macOS keeps working; `WRM_MEDIUM_RUNTIME=docker` moves it to Docker.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Nothing is installed until an administrator presses **Install Docker**. Docker Desktop is free for personal use, education, non-commercial open source and small businesses; larger organisations need a paid Docker subscription.

[Download v12.6.6](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.6) · [Docker for Medium](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.6/README.md#docker-for-medium-installed-by-wrm)

### Included from v12.6.5 — Medium tells you what it is doing

#### Added

- **Step by step:** the original deploy tool's card in Medium is a numbered guide: 1 import the deploy-tool ZIP, 2 prepare its runtime, 3 your deploy-config, 4 your SSH profile (optional), 5 start/stop. Each step shows whether it is done, running, waiting or failed, why a button is disabled and what to do when something fails.
- **Signals:** a light for the original tool (off, working, running, needs attention) and one for the read-only WRM preview, which is now a separate panel.
- **What runs, who and since when:** while an operation runs, everyone sees it, and its starter sees every step with its time and the live output of Docker, pip and the original tool. A second click is answered with the operation, the user and the start time instead of "busy".
- **Activity log:** every Medium action with time, user, result and duration; click a row for its steps and output. Secrets in the output are redacted. Administrators decide who sees output with `medium_log_output` (everyone, administrators, nobody).
- WRM's console logs each Medium operation and step.

#### Changed

- Your deploy-config ZIP is checked by WRM immediately and, while the runtime is still being prepared, kept and applied at your next start, so you no longer wait to upload it.
- Start and the tool import run in the background and the page follows them; a failure stays shown with the failed step and the last lines of the output.

#### Fixed

- **Approve later after the AI app closed its session:** approving such a request in WRM failed with "AI session not found". It now works for the request's hour, and the AI runs the same command once when it repeats it in a new session.
- **Search patterns are not files:** a `grep` regular expression, a `sed` or `awk` script and `find -name` patterns no longer count as file names, so a pattern such as `'^[[:space:]]*command=.*'` is not taken for a PuTTY key and does not ask for approval.

#### Upgrade

Replace the executable/container and keep the database and encryption key. One additive table (`medium_activity`); a previous binary still starts on the database. API clients: `POST /api/git/medium/start` now answers `202` (poll `GET /api/git/medium`).

[Download v12.6.5](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.5) · [Medium step by step](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.5/README.md#medium-step-by-step-status-and-activity-log)

### Included from v12.6.4 — Your blocked commands, your choice

#### Added

- **Blocked commands with full access** (`ai_full_access_blocked`, *Admin panel → AI connections*): the administrator decides what happens to their own list of blocked commands (`ai_blocked_commands`) when full access is on.
  - **Ask for approval** (default): as in v12.6.3, such a command comes to you as an approval marked *FULL ACCESS*.
  - **Stay blocked**: such a command is refused even with full access, and the AI is told that an approval cannot help and what you can change instead.

The built-in list of destructive commands, deny lists, shutdown and secret files keep asking under full access either way.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Nothing changes until an administrator sets *Blocked commands with full access* to *Stay blocked*.

[Download v12.6.4](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.4) · [MCP setup](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.4/README.md#ai-desktop-apps-mcp)

### Included from v12.6.3 — Full access with one switch

#### Added

- **The administrator chooses how full access applies** (`ai_full_access_mode`, *Admin panel → AI connections*):
  - **Every AI session** (default): the **Full access for the AI** switch alone opens full access for every AI session, the WRM AI panel and every AI connection, without ticking a scope anywhere.
  - **AI connections with the Full access scope**: as in v12.6.2, a user ticks **⚠ Full access** on an AI connection.
  - **Users who turn it on themselves**: each user gets a **Full access for my AI sessions** switch in *Settings → AI connections*; it applies to all of that user's AI sessions.
- *Settings → AI connections* shows whether full access is on and how it applies.

Full access still never runs anything WRM would refuse on its own: such a command or file comes to you as an approval marked *FULL ACCESS*, every time, and read-only sessions stay read-only.

#### Upgrade

Replace the executable/container and keep the database and encryption key. One additive column (`users.ai_full_access`); a previous binary still starts on the database. **If you had already turned on `ai_full_access` in v12.6.2,** it now applies to every AI session (mode *Every AI session*); set `ai_full_access_mode` to *AI connections with the Full access scope* to keep the v12.6.2 behaviour.

[Download v12.6.3](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.3) · [MCP setup](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.3/README.md#ai-desktop-apps-mcp)

### Included from v12.6.2 — Full access, but Claude still asks

#### Added

- **Full access with approval:** an administrator turns on **Full access for AI connections** (`ai_full_access`, *Admin panel → AI connections*). After that, an AI connection can get the **⚠ Full access** scope. With it, WRM refuses nothing outright: destructive commands, blocked commands, your deny lists, shutdown and reboot, protected paths and secret files all come to you as an approval marked *FULL ACCESS*, with the full command and the rule WRM would otherwise apply. They are asked every time and never run automatically, even with "Allow everything". Ordinary commands keep working as the connection is set up.

#### Changed

A refused command now tells the AI that full access exists and that you or your administrator turn it on, so it can explain the choice to you.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Nothing changes until an administrator turns on `ai_full_access` and you give a connection the *Full access* scope (edit the AI connection; open sessions end and are reopened). A read-only connection stays read-only even with full access.

[Download v12.6.2](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.2) · [MCP setup](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.2/README.md#ai-desktop-apps-mcp)

### Included from v12.6.1 — "Allow everything" means everything, and nothing ends in a dead refusal

#### Added

- **Approve later:** when the AI app declines WRM's question, stops waiting, or nobody answers in time, the request is no longer lost. It stays in WRM for an hour as a card with the full command and **Approve** / **Deny** (also in the live panel). Approve it, tell the AI to continue, and the same command runs once without asking again.
- **Files that may hold secrets:** reading a configuration or credentials file WRM treats as sensitive now asks you instead of failing; after you approve, the AI gets the file with passwords, keys and tokens replaced by `[REDACTED]`. Administrators can keep such files blocked with `ai_sensitive_reads` = deny.

- **Allow everything:** when you create or edit an AI connection in automatic mode, **Allow everything** fills in `*` for commands and `/**` for files, and **No time or action limit** keeps the connection automatic for as long as it is valid (up to the lifetime you chose). The AI then runs commands on the chosen servers without asking, also backups such as `cp -a a a.bak && cp -a b b.bak && ls -la`. Administrators can switch this option off with `ai_auto_until_grant`.

#### Changed

- `*` in the allow list now covers chains, variables, groups and redirections. Blocked commands (`ai_blocked_commands` and the built-in destructive list), your deny lists, sensitive files and commands WRM cannot read still ask or are refused, and the AI connection still reaches only the servers you chose, until it expires.
- Commands that stay blocked (the destructive list, your deny lists, read-only sessions) now say so: the AI is told that an approval cannot help and what you can do instead.
- When WRM does ask and the answer is no (or nobody answers), the AI is told which rule asked, for example "the automatic mode time limit is over" or "not in the session's allow list", what you can change in *Settings → AI connections* so it is not asked again, and that it must tell you this and ask before trying again.

#### Fixed

A long automatic session switched back to asking once its time limit (at most 60 minutes by default) was over, and commands with variables asked even with `*` allowed. When the AI app declined that question, the AI only saw "The user denied this action" although you had allowed everything, and there was nowhere to approve it afterwards.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. Reading files that may hold secrets now asks instead of being refused; set `ai_sensitive_reads` to `deny` in *Admin panel → AI assistant* to keep the old behaviour. Existing AI connections keep their limits: to allow everything without limits, edit the connection, press **Allow everything**, confirm and save (open sessions end; the AI opens them again with the new rules).

[Download v12.6.1](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.1) · [MCP setup](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.1/README.md#ai-desktop-apps-mcp)

### Included from v12.6.0 — Medium's native engine, side by side with the original tool

#### Added

- **Native engine preview (read-only):** open **Medium** in the Git workspace and expand *Native engine (preview, read-only)*. WRM reads your `deploy-config` as it is, with the same validation as the original tool, and lists every server of `servers.yaml` with the WRM connection it will use. Connections are matched by address and user; pick another one where needed.
- **Plan:** choose a project, an environment and a branch or tag, then **Preview natively**. WRM prepares the source from your WRM Git source and compares it with each destination over SSH, the way `rsync -rlc --delete` with the project's exclusions would, and shows new, changed and deleted files, the exclusions, what is on the server (version, ref, lock, recent deployments) and the `post_deploy` command. Nothing is written to any server.
- **Compare with the original tool:** with Medium running, **Compare with the original tool** runs the original comparison and the native one on the same source and says per server whether they give the same result, listing any change found by only one of them. Projects with a build step use the original tool's built source until the native engine arrives.

#### Changed

The Medium bridge has one more authenticated endpoint for the comparison; the original tool itself is not modified.

#### Upgrade

Replace the executable/AppImage/container and keep the database, encryption key and Medium directories. One additive table (`medium_server_links`); a previous binary still starts on the database. The preview needs a WRM Git source (GitLab, GitHub or Gitea API) for the repositories named in `projects/*.yaml` and saved SSH connections for the servers. Deploying is still done by the original tool; the native deploy, build location and rollback follow in v12.7.0 ([roadmap](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.0/ROADMAP.md)).

[Download v12.6.0](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.6.0) · [Medium native engine](https://github.com/vedranius/web-browser-RDM-public/blob/v12.6.0/README.md#medium-native-engine-preview)

### Included from v12.5.1 — Claude Desktop and Claude Code in one step

#### Added

- **Claude Desktop extension:** in *Settings → AI connections* pick **Claude Desktop**, choose the servers and how long the token is valid, then **Download Claude Desktop extension**. Open `wrm-pro.mcpb` and press **Install**: the WRM address and token are already filled in. The extension runs on your computer with Claude Desktop's built-in Node.js, so it reaches a WRM on your own network and over plain HTTP; a self-signed HTTPS certificate can be trusted with a CA file in the extension's settings. Approvals still appear in WRM and, where supported, in Claude.
- **Claude Code plugin:** `claude plugin install wrm-pro --marketplace vedranius/web-browser-RDM-public --config wrm_url=… --config wrm_token=…` (Claude Code 2.1.275 or later). The wizard prints the command with your values; the token is kept in Claude Code's secure storage.
- **See what the AI does:** when an AI app works on a server through WRM, a live panel shows each command, WRM's decision, the approval buttons, the exit code and the output. *Watch live* in the 🤖 menu reopens it; the automatic opening can be turned off there.
- **Lifetime in hours or days:** the wizard, the advanced form and the OAuth approval dialog take the token lifetime in hours or days, so apps need not be approved again every day.

#### Changed

When you deny a command, or nobody answers in time, the AI is told where it was refused and to ask you whether to continue, change the command or stop, rather than trying again on its own. Default AI connection lifetime is 7 days (`ai_mcp_token_hours` 168), at most 30 days (`ai_mcp_max_token_hours` 720).

#### Fixed

Claude Desktop could not be connected to a WRM on a private network: the wizard offered only the claude.ai custom connector, which connects from Anthropic's cloud and needs a public HTTPS address, and `claude_desktop_config.json` does not accept HTTP entries. The wizard now names that option *Claude.ai (cloud connector, public HTTPS)* and offers the extension for Claude Desktop. README troubleshooting explains an empty `mcpServers` block.

#### Upgrade

Replace the executable/container and keep the database and encryption key. No database change. MCP stays off by default (`ai_mcp_enabled`). Token lifetimes default to **7 days, at most 30 days** (were 8 hours, at most 24); values an administrator saved explicitly stay in force and can be changed in *Admin panel → AI connections*. The extension contains the token: delete the downloaded file after installing. Windows packages are cross-compiled; the extension was validated with the official MCP bundle tool and its bridge was run against WRM, but not inside a native Windows Claude Desktop in this release's validation.

[Download v12.5.1](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.5.1) · [MCP setup](https://github.com/vedranius/web-browser-RDM-public/blob/v12.5.1/README.md#ai-desktop-apps-mcp)

### Included from v12.5.0 — Stamping as it is, Back in dialogs and reviewed fixes

#### Added

- **Stamp as it is:** the Stamp dialog lists each installation's differing files and says why it cannot be stamped instead of "No files are selected". **Also stamp installations that are not up to date** writes VERSION.md with each file's actual version (older files their matching version, hand-changed files `(local)`) and a `Napomena` note that the installation is partial.
- **Back in dialogs:** a Git dialog opened from another one (Details → Update, Upgrade → file selection → confirmation, Rollback …) has **← Back** (Alt+←), which returns to the previous dialog with its ticks and typed values. Settings → Import inventory, Settings → Admin panel and Tunnels → Edit connection return to where they were opened.
- **Medium forced stop:** when the original tool cannot report its state, Stop can be confirmed a second time and stops it anyway.

#### Fixed

- **Production confirmation** in the Git workspace compares the typed connection name (it compared a hostname, so production updates could not start).
- **Medium:** one user's start or setup no longer blocks other users' status, proxy and Stop; empty settings keep the original defaults and WRM's SSH options; a successful start or import clears an earlier setup error and older pending imports; an imported config replaces a config repository; replaced SSH profiles (and their keys) are deleted; ZIP names that differ only in case are rejected; status polling continues after a network error; stage texts are translated.
- **Logs:** a non-UTF-8 byte no longer stalls a search; search segments end at a complete line; keyword highlighting is literal (`foo(`, `c++`, `[ERROR]`); empty files complete; Previous segment no longer skips data; logs are discovered again after Monitoring settings change; a closed tab during Check now no longer records a failed sample or a problem notification; a late check does not repaint another tab.
- **Client GitLab cache:** a newer snapshot published by another member is not deleted; WRM relay errors are no longer reported as GitLab CORS problems; one failed status poll no longer blocks later refreshes.
- Following a log in a terminal keeps the shell after Ctrl+C; moving or deleting a folder rechecks folder shares, so shared terminals outside the shared subtree end.

#### Upgrade

Replace the executable/container and keep the database, encryption key and private runtime directories/volumes. No schema migration. Earlier Medium versions kept replaced SSH profiles as `.previous` directories; they are removed at the next start.

The next versions replace Medium's Docker runtime with WRM's native engine (same original interface and deploy-config, WRM connections and keys, chosen build location); see [ROADMAP.md](https://github.com/vedranius/web-browser-RDM-public/blob/v12.5.0/ROADMAP.md). [Downloads](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.5.0).

### Included from v12.4.5 — One-click Medium and usable log search

#### Added

- **Medium:** add the original trusted tool/config ZIPs and optional SSH ZIP, then **Prepare and run**. WRM starts Docker, prepares the original dependencies, validates configuration, synchronizes profiles and opens the unchanged tool. Approved imports survive setup failure and can be resumed without re-uploading.
- **Logs:** automatic discovery from manual roots, service log folders and Supervisor metadata; progressive server results, visible errors, operator/service filters and historical cross-file keyword/date search. Tiles show modification/event times, highlighted excerpts and open/follow actions. Search can be canceled or continued when a bound is reached.
- **Environments:** combined text/service/environment/server filters with sticky controls.

#### Fixed

The whole Overview/Installations filter/check/action area stays at the top during scroll. Monitoring tab responses are isolated. Date-only searches and continuation-event timestamps are handled. Discovery avoids repeated per-file metadata round trips; host metric versus log prerequisites are explained.

#### Upgrade

Retain the database/key and private runtime directories/volumes, then replace the executable/container. No schema migration. Windows setup may display the official Docker installer and vendor/OS prompts; required license choices, WSL/virtualization setup and restarts remain user actions. A Windows service needs a preinstalled Docker account; other Windows containers are never stopped or switched automatically. Linux without installed Docker/native prerequisites needs its administrator's initial package setup. Windows builds are cross-compiled, without a native Windows VM run.

Search is bounded per file/round and excerpts; incomplete results show Continue. Compressed/binary logs remain explicit unsupported text inputs. See [guide](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.5/docs/V12.4.5_GUIDE.md), [validation](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.5/docs/RELEASE_SUMMARY_v12.4.5.md) and [downloads](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.5).

### Included from v12.4.4 — Monitoring, logs and AI workflows

#### Added

- **Monitoring** brings Zabbix server metrics, history, agent detection and service Functionality into one workspace. Use passive agent/agent 2 over SSH or your configured network route, verified certificate TLS, PSK through the existing remote `zabbix_get`, or an existing HTTPS Zabbix server API for active-only agents. Unavailable metrics are unknown and diagnostics remain visible.
- **Logs** supports multiple roots/operators/services, historical segmented search and date/timezone/keyword filters, highlighted results, copying/downloads, terminal actions and up to 16 movable/resizable live windows with tiling/cascade/pause/rotation handling.
- **Client GitLab tokens and cache** let a computer with GitLab access supply bounded repository snapshots to WRM. The GitLab token stays in browser memory or the companion's local token file. The latest complete successful server cache survives client disconnect and WRM restart; refresh is atomic, sharing is explicit and reviewed deployment content stays pinned.
- **MCP onboarding and permissions:** CLI setup scripts, OAuth/cloud guidance, Cursor install link, `wrm` prompt and authorized server picker, native assistant `/wrm`, safe structured configuration inventory, visible Automatic within limits and editable existing grants. Grant edits retain bearer/expiry and terminate old sessions/approvals.
- **Terminal and navigation:** previewed/redacted selection-to-AI context, Open terminal here for POSIX/PowerShell, Copy/Paste/native-menu fallback, hostname groups/combined Functionality filters and folder/connection colors/icons.

#### Fixed

General import recognizes mRemoteNG XML/PuTTY `.reg` ([issue #5](https://github.com/vedranius/web-browser-RDM-public/issues/5)); synthetic 4 MiB XML and UTF-16 registry exports were exercised in the browser. Provider header/request/idle timeouts are configurable for longer local AI requests. Wide replay fits phone/desktop widths. Search/filter bars remain at the top of supported scroll areas, and MCP wizard controls align correctly on narrow screens.

#### Upgrade and operation

Replace the executable/container and retain the database, encryption key and private runtime files. Migrations are additive. Functionality is now **Monitoring → Functionality**. Configure each server's installed Zabbix agent and log roots in Monitoring settings; automatic polling is opt-in. MCP stays disabled by default and existing explicit administrator policies remain in force.

Browser GitLab access needs the GitLab CORS/origin/local-network permissions; use the WRM companion when the browser cannot reach that API. The relay is for WRM repository API reads: native Medium Git clone/build/config Git operations and GitLab writes still need a direct route. Native Medium's Windows Linux-container prerequisite from v12.4.3b remains.

[Download v12.4.4](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.4) · [setup/limits](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.4/docs/V12.4.4_GUIDE.md) · [validation](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.4/docs/RELEASE_SUMMARY_v12.4.4.md). Windows packages are cross-compiled; no native Windows VM execution is claimed.

### Included from v12.4.3b — Windows Medium runtime fix

#### Fixed

Windows `.exe` installations run the original Medium deploy tool inside a managed Linux Docker runtime. Python, Git, OpenSSH and rsync are supplied by its pinned image; separate Windows installations are unnecessary. ZIP imports survive setup failure and can be prepared again without re-uploading. Private Linux volumes preserve deploy history, rollback state, permissions and config edits across restarts. Native Linux mode remains available.

The original tool remains unchanged. Runtime source/venv are read-only, SSH/config are private and builds use an unprivileged user. Runtime environment forwarding preserves Docker CLI context and keeps secret values out of command arguments. Only owned containers are stopped.

#### Windows upgrade

**Download the new Windows executable manually:** v12.4.3 cannot parse a `b` suffix; v12.4.3b supports it for later updates. Retain the DB/key. Start Docker Desktop with **Linux containers** on Windows 10/11, then use **Git → Medium → Prepare runtime / retry** for an existing tool import, or import your trusted tool ZIP first. Import deploy-config/SSH profile and start Medium. A browser client needs no deployment dependencies.

Docker Desktop is not supported on Windows Server; use a supported Linux WRM host/VM/container for Medium. The Docker engine must be local and accessible to the WRM account. Back up private `.medium` files and the managed Linux runtime/user Docker volumes.

#### Validation and packages

Real pinned-image setup, original comparison/build, two actual OpenSSH/rsync deployments, rollback and stop/start history persistence were exercised against a local fixture. The build user could not read private SSH/config. Browser layouts at 390/1100/1440/1920 px had no JavaScript errors. Windows is cross-compiled; no native Windows VM run is claimed. Full validation: [hotfix summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3b/docs/RELEASE_SUMMARY_v12.4.3b.md).

[Download v12.4.3b](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.3b): 16 platform packages plus checksums; Docker `ghcr.io/vedranius/wrm-pro:v12.4.3b` and `latest`. The remaining accepted work is planned for **v12.4.4** in ROADMAP, not included in this hotfix.

### Included from v12.4.3 — Original Medium deployment, config profiles and Bitwarden

#### Added

- **Medium** uses the original supplied deploy tool and deploy-config format unchanged: original screens, comparisons, file selection/exclusions, builds, native plans/confirmations, deployment progress, backups, history, rollback and config Git controls. Import the private tool/config ZIPs in Git → Medium. WRM handles authenticated access, private directories, encrypted runtime preferences and setup status. [Original-tool parity record](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/docs/MEDIUM_PARITY_v12.4.3.md).
- **Functionality config profiles** inspect uploaded reference settings across servers/operators and show searchable reference/actual values, differences, missing settings and masked secrets. INI/config/Python parsing is static; incomplete checks remain explicit.
- **Bitwarden** account sign-in/unlock and direct native SSH item references for SSH/SFTP. Private key retrieval happens on the server during authentication. Master passwords and unlock sessions are not persisted; per-user CLI profiles are encrypted between commands.
- **Folder order** persists after drag/drop, Move up/down and export/import, including nested folders.

#### Fixed and changed

- Git dialogs now stay visible above mobile drawers and keep wide tables within their scroll area.
- The bridge supports original same-origin no-referrer forms and rejects cross-site submissions; stop refuses active original-tool operations.
- Docker includes Python, git, OpenSSH, rsync and the official Bitwarden CLI. Standalone installations need these optional dependencies on the WRM host. Medium defaults to administrators; grant it only to trusted operators because its native build/config commands run as the WRM OS account.

#### Validation and packages

Original tool: **1,581 tests passed, 9 skipped**. Real Chromium exercised private archive imports and native comparison, deploy, history and rollback against test Git/SSH destinations. Config profiles, folder persistence and layouts at 390/1100/1440/1920 px were checked without JavaScript errors. Bitwarden mock-account flows and direct real SSH authentication were tested; a real cloud account was not used. Full WRM race/vet/format/UI checks and platform builds are detailed in the [release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/docs/RELEASE_SUMMARY_v12.4.3.md).

Release packages: 13 standalone targets, macOS universal and Linux amd64/arm64 AppImages, plus `SHA256SUMS.txt`. Docker: `ghcr.io/vedranius/wrm-pro:v12.4.3` and `latest`. The pipeline publishes only after the complete package set and Docker succeed.

#### Upgrading from v12.4.0

Replace the binary/AppImage or update Docker. Startup adds the config-profile table and folder order column. Preserve your database, encryption key and private adjacent `.medium`/`.bitwarden` directories. Tool and config archives are private imports, not public release assets. Git/SSH targets must be reachable from the runtime. The earlier client GitLab relay/replay, AI/navigation and icon/color requests remain accepted backlog, and were not included in this prioritized release.

**Downloads:** [v12.4.3 packages and checksums](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.3). [Changelog](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.3/CHANGELOG.md).

### Included from v12.4.0 — Service functionality, nested folders and remote previews

Inspect service functionality and configuration on your servers, distinguish identically named connections by client/site, and preview remote files directly in WRM.

#### Added

- **Git → Functionality** checks every Supervisor worker and failed systemd service on selected servers, as well as each catalog service's runtime units. It shows installation problems, disk/memory warnings, timestamps and bounded error logs. Optional HTTP/TCP health probes run from the selected server. Unavailable checks remain unverified; problem/recovery notifications are available.
- **Configuration** shows section, variable/function, value and description for each service's config/INI/environment files and JSON. Secrets are masked. Paths stay inside the installation; contents are read on demand with explicit limits and parse-gap warnings.
- **Nested folders** ([#57](https://github.com/vedranius/web-browser-RDM-public/issues/57)): recursive tree, parent selection, complete folder paths in connection pickers/windows/tabs, folder tab filtering and hierarchy-preserving imports. Folder shares include descendants; route defaults remain local to the immediate folder.
- **Remote previews** ([#55](https://github.com/vedranius/web-browser-RDM-public/issues/55)): video/audio streaming with seekable byte ranges over SFTP/FTP, image/PDF viewing and bounded text viewing. Select a path in a terminal and right-click to preview/download it; the file manager offers Preview too. Browser codec support determines playable media formats.

#### Fixed

- **AI entry points and diagnostics** ([discussion #56](https://github.com/vedranius/web-browser-RDM-public/discussions/56)): AI is visible in SSH/SFTP windows; unavailable permissions lead to an explanation. Settings → AI assistant and the top AI button provide setup/status access. Policy changes refresh other open browsers. A successful provider test is distinguished from account permission, and the terminal assistant requires no MCP setup.
- Missing installations stay visible with a warning. Supervisor programs in the same configuration file are associated with the correct service stanza.
- Docker Compose uses the current version consistently.

#### Packages and validation

Linux, Windows, macOS (including universal), Android/Termux, FreeBSD and OpenBSD standalone binaries; Linux amd64/arm64 AppImages; Docker `ghcr.io/vedranius/wrm-pro:v12.4.0` and `latest`. The release pipeline verifies all 16 downloadable packages and `SHA256SUMS.txt` before publication.

Validation details: [v12.4.0 release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.0/docs/RELEASE_SUMMARY_v12.4.0.md).

#### Upgrading from v12.3.6

Replace the binary/AppImage or update the container. Startup adds `folders.parent_id` and a service-health cache table without removing existing data. Previous binaries ignore the new fields; run the previous binary to downgrade. Existing catalogs need no new fields. Git runtime checks use the SSH account's read permissions and do not restart services.

**Downloads:** [v12.4.0 packages and checksums](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.4.0). See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v12.4.0/CHANGELOG.md) for the complete history.

### Included from v12.3.6

This release adds deployment review and server selection to Git Simple, gives AI tools the server's actual date, and improves Linux systemd service installation and failure diagnostics.

#### 🤖 AI assistant

- Ollama setup releases its action lock before reporting completion, so the next setup step can start immediately. This fixes the intermittent HTTP 409 exposed by the release CI run.

- `system_info` now reports the remote server's actual date and time instead of leaving the model to infer it from the kernel build string.
- The assistant is explicitly instructed not to invent tool calls, tool output, logs or internal `<tool_response>` wrappers.
- Enabling or disabling the assistant refreshes its buttons on already-open terminals without reconnecting. Shared connection roles remain enforced; legacy lowercase SSH protocols are accepted.
- The processing indicator now has its English and Croatian translation.
- The indicator remains visible until the model finishes. Recognized tool JSON is removed from the finalized chat message, and incomplete JSON is no longer repaired into executable calls or silently removed from answers.

#### ⎇ Git Simple

- Git Simple now has combined service, hostname, environment, status, branch and version/commit filters. Its deploy review shows changed files and their state per destination, lets you select successfully compared servers and exclude files, and shows when server state was compared. Deploys keep using WRM's fingerprinted plan, backup, lock and audit flow.
- Empty server selections are refused; confirmation is recomputed whenever the selection changes. Canceled or superseded comparisons cannot reopen a review. Rollback shows each selected server's own previous version. Simple mode keeps removed files; Advanced mode offers deletion after backup.
- Connection changes after comparison invalidate the plan, so a reviewed deployment cannot silently follow a changed SSH destination or login.

#### 🐧 Linux service

- Installation checks read/write access to SQLite and its WAL/SHM files, plus write/search access to its directory, before changing service settings. It points to `DB_PATH` when permissions prevent startup.
- If systemd cannot start WRM, installation reports the unit status and recent journal entries; read-only database errors include a specific permissions hint.
- Generated units limit rapid restart loops after a persistent failure. The install summary shows the database path the service will use.
- Service definitions handle spaces, literal dollars and percent characters in paths correctly. Reinstalling clears the unit's old failures and restarts the newly configured binary.

#### 📦 Packages and validation

- **Linux**: amd64, arm64, ARMv7, ARMv6, 386; **AppImage** for amd64 and arm64.
- **Windows**: amd64 and arm64 `.exe` files.
- **macOS**: Intel, Apple Silicon and universal binaries.
- **Android / Termux**: arm64; **FreeBSD**: amd64 and arm64; **OpenBSD**: amd64.
- Docker: `ghcr.io/vedranius/wrm-pro:v12.3.6` and `latest` (Linux amd64).
- `SHA256SUMS.txt` covers all 16 downloadable packages. GitHub also provides source archives.

AppImage quick start: `chmod +x wrm-pro-v12.3.6-linux-amd64.AppImage`, then run it and open **http://localhost:8080**. On ARM64 choose the arm64 AppImage. Without FUSE, add `--appimage-extract-and-run`. Data stays in `${XDG_DATA_HOME:-$HOME/.local/share}/wrm-pro/` unless `DB_PATH` is set. Replace the AppImage to upgrade; use the standalone binary for boot services.

The pre-push review covered race tests, UI regressions, all 13 cross-compilation targets, native Linux/Docker health checks, Git Simple browser deployment and rollback, and a real local Ollama tool call. Other OS binaries were compiled; they were not run on their respective OSes. See [the review report](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/docs/PRE_PUSH_REVIEW_v12.3.5.md).

#### ⬆️ Upgrading from v12.3.4

Replace the binary or use the in-app updater. No database migration is required.

For an existing Linux service, re-run `-install-service` with the new binary to regenerate its unit and apply the path and restart-limit fixes. This restarts that service.

**Downloads:** [all v12.3.6 files and `SHA256SUMS.txt`](https://github.com/vedranius/web-browser-RDM-public/releases/tag/v12.3.6)

See [the release summary](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/docs/RELEASE_SUMMARY_v12.3.6.md) for package and validation details.

See [CHANGELOG.md](https://github.com/vedranius/web-browser-RDM-public/blob/v12.3.6/CHANGELOG.md) for the complete change history.
