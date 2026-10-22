# WRM PRO — Claude Code plugin

Connects Claude Code to a WRM PRO instance over MCP: logs, files, services and commands, under
WRM's permissions, approvals and audit. Claude never receives the server credentials.

## Install in the chat (two lines, once)

```
/plugin marketplace add vedranius/web-browser-RDM-public
/plugin install wrm-pro@wrm-pro
```

Claude then asks for the two values from **WRM → Settings → AI connections**:

- **MCP address**, for example `https://wrm.example.com/mcp`
- **token** (`wrm_pat_…`)

After that, `/wrm-pro:connect` reconnects or switches to another WRM instance at any time.

## Local WRM (`http://localhost:8080`)

Add it with the token header:

```sh
claude mcp add --transport http --scope user wrm http://localhost:8080/mcp \
  --header "Authorization: Bearer wrm_pat_…"
```

Claude Code opens its OAuth sign-in page only over HTTPS; without a token a local address fails
with `Refused to open sign-in URL for wrm: must be https`. The bearer header avoids OAuth, and
the token is the same kind created in WRM → Settings → AI connections.

## Without the plugin manager

```sh
claude mcp add --transport http --scope user wrm https://wrm.example.com/mcp \
  --header "Authorization: Bearer wrm_pat_…"
```

or, from this folder:

```sh
./scripts/connect-claude.sh https://wrm.example.com/mcp wrm_pat_…
.\scripts\connect-claude.ps1 -Url https://wrm.example.com/mcp -Token wrm_pat_…
```

`/wrm-pro:connect <url> <token>` does the same from inside a Claude Code session.

## Claude Desktop

WRM → Settings → AI connections → **Claude Desktop** downloads `wrm-pro.mcpb`; open it to install
the local bridge. The HTTP endpoint above works for any other MCP client too.
