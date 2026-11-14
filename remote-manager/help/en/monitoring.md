---
id: monitoring
title: Monitoring and server status
section: Operations
order: 70
summary: Live up/down status, the Monitoring window (services, functionality, metrics) and alerts when something breaks.
keywords: monitoring, status, up, down, ping, latency, alert, metrics, cpu, memory, disk, services, functionality, health, zabbix, checks
ui: #btn-monitoring, .mon-wrap, #mon-modal
open: monitoring
---
## Live status in the sidebar
WRM checks every saved connection in the background: 🟢 up (hover for latency and SSH version), 🟡 slow, 🔴 down (hover for why: refused, timeout, unknown host). A folder shows how many are down; the red bar shows only the servers that are down. Turn on desktop notifications in **Settings → General**.

## The Monitoring window
**📊 Monitoring** in the top bar:
- **Servers:** status and history of your connections, with metrics from a Zabbix agent: read through SSH, through the network route, or from an existing Zabbix server API (*Monitoring settings* per server; *Detect OS and installed agent* helps). *Check now* reads at once.
- **Services and functionality:** Supervisor workers, failed systemd units, disk and memory per server, and health URLs or TCP checks of catalog services.
- **Logs:** the log viewer (see *Logs*).

## Alerts outside WRM
**Settings → Notifications** chooses what you receive (server down/up, rotation due, new versions, chain state…) and where: e-mail, Telegram, Slack/Mattermost, Teams, Discord, ntfy, Gotify, Pushover or a signed webhook. Administrators add the channels in **Admin → Notifications**.

## Ask the AI
On a server window, **AI** answers questions such as *why is the load high?* by reading processes, services and logs (it only reads unless you approve a change).
