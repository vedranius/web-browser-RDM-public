---
id: notifications
title: Notifications
section: Personalise
order: 135
summary: In-app messages and alerts by e-mail, chat, push or webhook, and how to keep them tidy.
keywords: notification, alert, message, email, e-mail, telegram, slack, teams, discord, ntfy, gotify, pushover, webhook, bell, inbox, quiet
ui: #st-notify, #notify-pane, #modal-notify-ch
open: settings:st-notify
---
## In WRM
Short messages appear at the bottom right; important ones also stay in the inbox of the Workboard. While an AI app works, its steps go to the AI live view instead of separate messages.

## Outside WRM
**Settings → Notifications** chooses what you receive (server down and back up, password rotation due, new versions of services, chain state changes, security alerts for administrators…) and on which channel.

Administrators add channels in **Admin → Notifications**: e-mail (SMTP, set in *Admin → Mail*), Telegram, Slack / Mattermost, Microsoft Teams, Discord, ntfy, Gotify, Pushover or a signed webhook. *Test* sends one message.

Desktop notifications while WRM is in the background: **Settings → General → Notify when a server goes down** (the browser asks for permission once).
