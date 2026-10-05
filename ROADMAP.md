# Roadmap

Planned work that is agreed but not built yet. Released work is in [CHANGELOG.md](CHANGELOG.md).

## v10.8.1 — fixes

- **mRemoteNG and inventory import of larger files.** The security middleware limits every API body to 1 MB (`bodyLimitFor` in `security.go`); only `/api/config/import` has 32 MB. `/api/config/import/mremoteng`, `/api/config/import/sshconfig` and `/api/inventory/*` therefore fail with "Bad JSON" from about 1 MB on.
  - Give `/api/config/import/*` and `/api/inventory/*` a proper limit (32 MB).
  - A clear error when a file is too large (its size and the limit), and a size check in the browser before sending.
  - Tests with a multi-MB `confCons.xml` and a large CSV.
- **Install as an app:** explain in the UI what to do in the current browser.
  - Exact steps per browser (Chrome/Edge, Safari, iOS, Android).
  - The real reason when installing is not possible (plain `http://` to an IP address, Firefox on the desktop).
  - Hide the *Install* button when it cannot work.
  - README: a short HTTPS recipe (Caddy).

## v10.9.0 — proxy as a connection setting

Requested by users who want to replace PuTTY: some servers are reached through jump hosts, others through a SOCKS or HTTP proxy.

### Saved proxies
A proxy is defined once and picked per connection, like a vault credential. It has:
- a name;
- a type: **SOCKS5**, **SOCKS4/4a** or **HTTP CONNECT**;
- host and port;
- an optional user name and password (encrypted, or a vault credential);
- *DNS at the proxy* (PuTTY: "Do DNS name lookup at proxy end").

Proxies can be shared with other users without revealing the password. Export and import keep them by name.

### In the connection dialog
A **Proxy** field next to *Connect via (jump host)*: none / a saved proxy / *+ New proxy…*.

### What goes through the proxy
- Everything a connection does over TCP:
  - SSH terminal, SFTP, search, transfers, tunnels;
  - key deploy and password rotation;
  - FTP/FTPS;
  - RDP/VNC/Telnet (through a local relay, like jump hosts today);
  - web interfaces;
  - *Test connection*;
  - live status;
  - network tools from a server;
  - Redfish BMCs.
- Not IPMI and Serial-over-LAN: they use UDP. The UI says so.
- **Proxy and jump host together:** the proxy is reached through the jump hosts, then the target. The route is shown and audited, e.g. `bastion → socks5://10.1.1.1:1080 → app-01`.
- Proxy type **"WRM SOCKS tunnel of connection X"**: a dynamic tunnel that WRM starts when needed.

### Docker
Inside a container `127.0.0.1` is the container itself.
- `docker-compose.yml` gets `extra_hosts: host.docker.internal:host-gateway`.
- The README explains the setup.
- WRM warns when a proxy at `127.0.0.1` is set while it runs in a container.

### Folder defaults
A folder can carry a default jump host and proxy. Its connections, including new ones, inherit them unless they set their own.

### Security and errors
- Proxy passwords are encrypted and never sent to the browser.
- Policy `proxies`: who may define proxies.
- Use is audited with the route.
- Clear messages that separate proxy errors from target errors: proxy unreachable, authentication required or failed, SOCKS reply codes, HTTP status.

### Tests
In-process SOCKS5 and HTTP CONNECT proxies, with and without authentication, covering:
- terminal, SFTP, the desktop relay, web interfaces and status;
- jump host + proxy together.

### Later (v10.9.x)
- **Import PuTTY sessions** from a `.reg` export (`HKCU\Software\SimonTatham\PuTTY\Sessions`): host, port, user, key and proxy settings. Proxies are created and linked to the connections.
- With mRemoteNG, connections that point to a PuTTY session (`PuttySession`) get that session's proxy.

## Idea, awaiting decisions — Git: compare and deploy services

A separate **Git** workspace in WRM, outside the connection tree and the terminal UI. It is generic: it works with any GitLab or GitHub (cloud or self-hosted) that a user connects. WRM ships no services of its own and no organisation-specific configuration. The details are fixed once the open questions are answered.

1. Connect a GitLab group or GitHub organisation/user (URL + read-only token).
   - WRM lists the repositories.
   - The user picks which ones are services.
2. **Installations:** which service runs on which server (connection), in which folder, on which branch or tag. Version detection: git checkout, a version file, a command, Docker image tags.
3. **Batch check:** a servers × services matrix with versions and status (up to date, behind, ahead, local changes, unknown).
4. **Periodic check** while WRM runs. The user is notified when a server runs an old version.
5. **Upgrade:**
   - choose a branch or tag and see the commits and release notes in between;
   - update through git on the server, or WRM transfers the code;
   - post-install steps, rollback, rolling over several servers.
6. **Install or transfer to a new server:**
   - choose a folder;
   - the inputs the service needs (from a `.wrm/deploy.yml` manifest or `.env.example`), stored encrypted.
7. **.gitignore helper:** pick untracked files from a server installation and stack templates, preview, then a merge request or a download.
