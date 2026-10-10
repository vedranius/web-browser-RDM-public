# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

## v11.7.0: run as a service, self-update, file manager "here" or "/"

**A. Run as a service (all OS)**
- On an interactive first start (stdin is a TTY, not already a service, not in Docker, not `-no-service-prompt`) ask in the console whether to install WRM as a service that starts at boot: Yes, No, Don't ask again. The answer is remembered next to the database.
- Flags `-install-service`, `-uninstall-service` and `-service-status` for scripts and administrators.
- Windows: a real Windows service (SCM, restart on failure), with elevation or a clear "administrator rights needed" message.
- Linux: a hardened systemd unit (`NoNewPrivileges`, `ProtectSystem`, working directory, an env file for `PORT` / `DB_PATH`); a clear message when systemd is missing.
- macOS: launchd, as a LaunchDaemon (system) or LaunchAgent (user).
- FreeBSD / OpenBSD: an rc.d script.
- The service keeps the data directory and settings of the interactive start.
- Newest binary in the folder: on a service start (and with self-update) a newer valid `wrm-pro-v<semver>-<os>-<arch>[.exe]` in the binary's directory (it must answer `-version` with that version) is used; the UI offers "Restart to v…".

**B. Self-update from GitHub releases**
- Daily and on-demand check of the configured repository's releases (`update_check`, `WRM_UPDATE_REPO`, HTTP proxy, no user data sent).
- An "update available" badge on the version in the top bar and on the sign-in page, with release notes and a link.
- One-click update for administrators (policy `self_update`): download the asset for this OS / architecture, verify its SHA-256 against the release checksums, run `-version`, replace the binary atomically (Windows: rename to `.old`), restart (service manager or re-exec), reconnect the page, keep the previous binary for a one-click rollback, audit every step (`system.update_*`).
- Docker: never self-replace; show the `docker pull` / compose instructions.
- The database is only changed by the new version's normal additive migrations.

**C. File manager from the terminal: "here" or "/"**
- The 📂 button of an SSH terminal becomes a split button: Open here, Open /, Open home.
- "Here" works without OSC 7: WRM asks the server for the shell's working directory over a separate SSH exec channel; OSC 7 stays the fast path; when the directory cannot be found, WRM says so and offers home or `/`.

## v12.0.0: AI assistant core and the built-in assistant (enterprise grade)

- A per-session **AI permission model**:
  - **Read-only:** inspect logs, packages, services and config; no changes.
  - **Ask before every change:** every write command or file edit needs an explicit approval in the UI, showing the exact command or diff.
  - **Automatic within limits:** only after an explicit opt-in per session, with an allow / deny list of commands and paths and a time limit; destructive patterns are always blocked.
- Admin policies limit which modes users may choose, per connection, folder or tag.
- Full audit trail: every AI request, tool call, approval and denial. The session is recorded like a terminal recording.
- A kill switch.
- A **built-in assistant panel** next to the terminal, using the user's own provider:
  - the Anthropic API (a personal key, or the organization's key managed by administrators; the latest Claude models);
  - the OpenAI API;
  - Azure OpenAI, AWS Bedrock and Google Vertex for enterprises;
  - any OpenAI-compatible endpoint (local models).
  - Keys are encrypted at rest and never sent to the browser.
- The assistant sees only what the permission mode allows, runs commands through WRM's SSH layer with the guards above, and explains every step.
- Note for the docs: consumer Claude and ChatGPT chat subscriptions cannot be used by third-party apps. Use API keys, an enterprise gateway, or v12.1.0.

## v12.1.0: AI desktop apps (MCP)

- WRM is an **MCP server**: streamable HTTP with OAuth / personal tokens, plus a small stdio bridge. Claude Desktop, ChatGPT (connectors / MCP) and other MCP clients can open a WRM terminal session to a chosen server.
- Session-scoped tokens with the same permission modes and approvals: approval prompts appear in the WRM UI and as MCP elicitation where the client supports it.
- Per-tool scopes:
  - `read_logs`;
  - `run_readonly`;
  - `run_with_approval`;
  - `edit_file_with_approval`;
  - `transfer`.
- Time-limited sessions with a visible "AI connected" indicator, one-click revoke, and the full audit described above.
- Setup guides for the common desktop apps.
