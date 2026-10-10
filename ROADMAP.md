# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

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
