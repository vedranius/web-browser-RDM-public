---
id: ai-connections
title: AI apps over MCP
section: AI
order: 110
summary: Let Claude, ChatGPT, Claude Code, Codex, Cursor and other MCP apps work on your servers through WRM, with your approval.
keywords: mcp, model context protocol, claude desktop, claude.ai, chatgpt, claude code, cursor, gemini cli, opencode, codex, connector, token, oauth, ai connection, plugin, mcpb, bridge, approve
ui: #st-mcp, #mcp-pane, #btn-mcp
open: mcp_setup
---
WRM is an MCP server. An AI app on your computer can open a WRM session to a server you allowed and work on it under the same rules as the built-in assistant. The model runs in the app; WRM never needs a provider key for this.

## Connect an app
1. An administrator turns on **AI connections** (*✦ AI → Set up AI*, or *Admin → AI connections*).
2. **Settings → AI connections → Connect an AI app** and pick your app:
   - **Claude Desktop:** download the WRM extension (`.mcpb`) and open it.
   - **claude.ai / ChatGPT:** add WRM's MCP address as a connector (WRM must be reachable from the internet for these).
   - **Claude Code, Codex, Gemini CLI, Cursor, OpenCode:** copy the one command or the config file shown.
   - **All AI apps at once:** one script for Windows or Linux/macOS sets up every app found.
3. Choose the **servers**, the **mode** (read-only, ask, automatic) and how long the access lasts. The app signs in with OAuth or a personal token.

## While it works
**🤖 AI connected** in the top bar shows live sessions; a live view shows each step. Changes wait for your approval in WRM (or in the app when it can ask). Revoke a connection with one click; administrators see all of them in *Admin → AI connections*.

## Help in your AI app
Every AI app connected to WRM also gets the **wrm_help** tool: ask it *how do I add a jump host in WRM?* and it reads this help.
