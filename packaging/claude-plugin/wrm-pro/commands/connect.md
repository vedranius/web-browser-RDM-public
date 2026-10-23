---
description: Connect this Claude Code session to WRM PRO over MCP
argument-hint: "[wrm-url] [token]"
allowed-tools: Bash(claude mcp:*)
---
Connect this session to WRM PRO.

1. If `$1` and `$2` were given, use them as the MCP address and token. Otherwise ask for:
   - **MCP address**: WRM → Settings → AI connections shows it, for example `https://wrm.example.com/mcp`.
   - **Token**: the token of that AI connection (`wrm_pat_…`). It is sent only to the `claude mcp add` command below.
2. Run exactly one command, with the two values quoted:

```sh
claude mcp add --transport http --scope user wrm "<address>" --header "Authorization: Bearer <token>"
```

3. Check it with `claude mcp list` and tell the user which WRM connections the token can reach.
   Every call still goes through WRM's permissions, approvals and audit; the assistant never gets
   the server credentials.

Do not print the token back, do not write it to a file, and do not add it to `.mcp.json` in a
repository.

A local WRM (`http://localhost:8080/mcp`) works the same way, but only with the token header:
Claude Code opens a sign-in page only over HTTPS and otherwise reports
`Refused to open sign-in URL for wrm: must be https`. Adding the header avoids OAuth entirely.
