---
id: jump-hosts-proxies
title: Jump hosts and proxies
section: Access
order: 55
summary: Reach servers behind a bastion (like ssh -J) or through a SOCKS/HTTP proxy, also as folder defaults.
keywords: jump host, bastion, proxyjump, ssh -J, via, proxy, socks, socks5, http connect, folder default, behind, internal network
ui: #m-jump-help, #m-proxy-help, #modal-proxy
open: vault:proxies
---
## Jump host (bastion)
1. Save the bastion as a normal SSH connection.
2. In the internal server's connection choose **Jump host (connect via) → the bastion**.
3. Double-click it: the terminal says *Connecting to 10.0.0.21 via bastion…*. Bastions can have their own jump host (chains); files, transfers and tunnels use the same route.

## Proxy
For networks reached through a SOCKS5/4 or HTTP CONNECT proxy:
1. **🔑 Keys → 🌐 Proxies → New proxy** (or *＋ New proxy…* in the connection's Proxy field).
2. Address, port, optional login (or a credential), *DNS at the proxy*.
3. Choose it in a connection. A proxy can also be a SOCKS tunnel (-D) of one of your connections; WRM starts it when needed.

## Folder defaults
Edit a folder and set its default jump host and proxy: every connection in it (also new ones) uses them unless it chooses its own.

RDP/VNC desktops and web interfaces go through the jump host too; IPMI uses UDP and runs on the last jump host instead.
