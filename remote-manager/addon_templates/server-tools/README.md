# Server tools

A window that works on your own servers through WRM: it lists your SSH/SFTP connections
(`connections.read`), shows their current metrics (`monitoring.read`) and runs a command
(`servers.exec`). The addon always acts as the person who looks at the window, with that
person's own connections. Every command waits for that person's approval in WRM, outside the
panel, unless an administrator marked the addon as trusted. Output is limited to 256 KB and
secret-looking values are masked. Logins, passwords and keys never reach the panel.

1. Change `id` in `addon.json` (lowercase, for example `yourname.server-tools`) and `name`.
2. Check it: `wrm -addon-lint .` (or **Lint** in Settings → Addons → Addon Studio).
3. Install the draft, enable it and open the window.

The bridge calls are `wrm.connections()`, `wrm.metrics(id)`, `wrm.exec(id, command, {timeout_s, reason})`,
`wrm.readFile(id, path)` and `wrm.writeFile(id, path, content, {reason})`; see `llms.txt`.
