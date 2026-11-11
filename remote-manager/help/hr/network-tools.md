---
id: network-tools
title: Network tools
section: Operations
order: 95
summary: Port check, ping, traceroute, DNS, HTTP/TLS checks and more, from WRM or from any of your servers.
keywords: tools, network, port, ping, traceroute, dns, http, tls, certificate, check, reach, firewall, scan, nslookup, dig, curl
ui: #modal-tools, #btn-tools
open: tools
---
**🧰 Tools** in the top bar (or right-click a connection → *Network tools from this server…*).

## Tools
- **Port check:** up to 100 ports (`22,80,443`, `8000-8010`): open, closed or no answer.
- **Ping** and **Traceroute.**
- **DNS:** A, AAAA, MX, TXT, CNAME and more.
- **HTTP / TLS:** status, headers, redirects and the certificate with its expiry.

## From where
Run a tool **from the WRM server** or **from one of your SSH servers** (*can app-01 reach the database on 5432?*). From a server, WRM runs read-only commands over SSH; jump hosts apply.

Administrators decide who may use the tools (policy `network_tools`). **📣 Broadcast** is in the same menu (see *SSH terminal*).
