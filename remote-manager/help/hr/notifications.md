---
id: notifications
title: Obavijesti
section: Prilagodba
order: 135
summary: Poruke u aplikaciji i upozorenja e-mailom, chatom, pushom ili webhookom te kako ih držati urednima.
keywords: obavijest, obavijesti, upozorenje, poruka, email, e-mail, telegram, slack, teams, discord, ntfy, gotify, pushover, webhook, zvono, sandučić, tiho, notification, alert, message, bell, inbox, quiet
ui: #st-notify, #notify-pane, #modal-notify-ch
open: settings:st-notify
---
## U WRM-u
Kratke poruke pojavljuju se dolje desno; važne ostaju i u sandučiću Workboarda. Dok AI aplikacija radi, njezini koraci idu u AI pregled uživo umjesto zasebnih poruka.

## Izvan WRM-a
**Postavke → Obavijesti** bira što primate (server ne radi i ponovno radi, dospjela rotacija lozinke, nove verzije servisa, promjene stanja lanca, sigurnosna upozorenja za administratore…) i na kojem kanalu.

Administratori dodaju kanale u **Admin → Obavijesti**: e-mail (SMTP, postavlja se u *Admin → Pošta*), Telegram, Slack / Mattermost, Microsoft Teams, Discord, ntfy, Gotify, Pushover ili potpisani webhook. *Test* šalje jednu poruku.

Obavijesti na radnoj površini dok je WRM u pozadini: **Postavke → Općenito → Obavijesti kad server prestane raditi** (preglednik jednom traži dopuštenje).
