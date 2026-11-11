---
id: tunnels
title: SSH tunnels and web interfaces
section: Access
order: 60
summary: Local, remote and SOCKS port forwarding, the Tunnels panel and web interface connections behind servers.
keywords: tunnel, port forwarding, ssh -L, ssh -R, ssh -D, socks, local, remote, dynamic, web interface, https, idrac web, switch web ui, browser
ui: #modal-tunnels, #btn-tunnels, #tun-help
open: tunnels
---
## Kinds of tunnels
- **Local (-L):** a port on the WRM machine leads to `host:port` as seen from the server. Example: `127.0.0.1:2443 → 10.0.0.5:443` opens a switch's web page behind the server.
- **SOCKS (-D):** a proxy for your browser or tools that reaches everything the server reaches.
- **Remote (-R):** a port on the server leads back to something WRM reaches (needs the policy `tunnel_remote_forward`).

## Set up and use
1. **Edit connection → 🔀 SSH tunnels → Add** (or right-click → *Add tunnel…*).
2. Start and stop them in **🔀 Tunnels** in the top bar. Running tunnels reconnect by themselves.
3. **Open** opens a web tunnel in a new browser tab; the copy button copies the address.

The listen address `127.0.0.1` (default) keeps the tunnel to the WRM machine. A network address (`0.0.0.0`, a LAN IP) lets other computers use it; only administrators may set one unless the policy `tunnel_bind_any` allows it. Leave the port empty to get a free one.

## Web interface connections
A connection with protocol **Web interface (HTTPS/HTTP)** opens that page with a double-click. With a jump host WRM opens a temporary tunnel itself and stops it after 30 minutes without traffic.
