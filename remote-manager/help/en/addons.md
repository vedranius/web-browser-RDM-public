---
id: addons
title: Addons
section: Personalise
order: 120
summary: Extend WRM with addons (panels, themes, server code), install from the community list or build your own with AI.
keywords: addon, addons, plugin, extension, install, zip, community, theme, panel, wasm, build addon, develop, permissions, instance
ui: #st-addons, #addons-pane
open: addons
---
An addon is a ZIP with `addon.json`: a panel (window), colour themes and, for organisation addons, WebAssembly server code. Nothing runs at install time and a new addon is installed **off**.

## Install
**Settings → Addons**:
- **Community** lists published addons with their trust level; the download is checked against its SHA-256.
- **Upload ZIP** installs your own. Administrators install for everybody or only for themselves; others install **personal** addons (policy `addon_personal`).
- **Enable** shows the **permission sheet**: what the addon may do (network origins, secrets, notifications, server code). An update that asks for more waits for your approval again.

## Use
*Open panel* opens it as a window on the canvas (or dock it as a panel). **Instances** keep separate settings, for example one *Server metrics* per cluster.

## Build your own
- **With AI:** ask the AI assistant (or an AI app over MCP) to build an addon; it writes a **draft** that you review as a diff in Settings → Addons and install yourself.
- **By hand:** start from a template (*New addon*), use the dev server and `addon lint`; the format is in `docs/ADDONS.md` on GitHub.
