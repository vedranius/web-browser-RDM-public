# WRM PRO — Claude Code plugin

Connects Claude Code to a WRM PRO instance over MCP: logs, files, services and commands, under
WRM's permissions, approvals and audit. Claude never receives the server credentials.

## What is in it

- **MCP server `wrm`**: your own WRM over HTTP (`https://<your WRM>/mcp`) with a token from
  WRM → Settings → AI connections. Every tool call runs on WRM: it checks the connection's mode
  and scopes, asks you to approve changes, redacts secrets and writes the audit log.
- **Skill `wrm-safe-operator`**: how to work through WRM (sessions, reading before changing,
  saying what will run before every change, what to do after a denial).
- **Skill `linux-incident-triage`**: a read-first checklist for a server or service that is down,
  slow, erroring, out of disk or memory. It works with WRM's tools (`search_logs` greps every
  discovered log at once) or with a plain shell.
- **Command `/wrm-pro:connect`**: set or change the WRM address and token.

### Where the data goes

The plugin itself has no code that runs on its own and sends nothing anywhere. Claude talks only
to the WRM address you enter; WRM talks to your servers over SSH. Server passwords and keys stay
in WRM and are never given to Claude. What Claude reads (command output, log lines, files) goes to
the model provider of your Claude account, like anything else in the chat, after WRM's redaction.

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

## Skills without the plugin

The two skills are plain Agent Skills folders (`skills/<name>/SKILL.md`). Agents that read skills
from a folder can use them directly; the skills CLI installs them from this repository:

```sh
npx skills add vedranius/web-browser-RDM-public
```

## MCP Registry

`server.json` in the repository root describes the WRM MCP endpoint for the
[MCP Registry](https://registry.modelcontextprotocol.io) as `io.github.vedranius/wrm-pro`: a
remote server whose host is yours (`https://{wrm_host}/mcp`) and an `Authorization` header with
your token. Every release publishes it automatically: the release workflow starts
`.github/workflows/mcp-registry.yml`, which signs in with the repository's GitHub Actions OIDC
token (no stored secret) and runs `mcp-publisher publish`. For an existing tag, run the
**MCP Registry** workflow by hand (Actions → MCP Registry → Run workflow, tag `vX.Y.Z`).

## Claude Desktop

WRM → Settings → AI connections → **Claude Desktop** downloads `wrm-pro.mcpb`; open it to install
the local bridge. The HTTP endpoint above works for any other MCP client too.
