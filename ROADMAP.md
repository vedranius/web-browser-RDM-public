# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

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
