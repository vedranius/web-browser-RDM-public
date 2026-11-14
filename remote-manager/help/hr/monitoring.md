---
id: monitoring
title: Nadzor i stanje servera
section: Rad
order: 70
summary: Stanje radi/ne radi uživo, prozor Nadzor (servisi, funkcionalnost, metrike) i upozorenja kad se nešto pokvari.
keywords: nadzor, monitoring, stanje, radi, ne radi, ping, latencija, upozorenje, alarm, metrike, cpu, memorija, disk, servisi, funkcionalnost, zdravlje, zabbix, provjere, status, alert, metrics, health, checks
ui: #btn-monitoring, .mon-wrap, #mon-modal
open: monitoring
---
## Stanje uživo u bočnoj traci
WRM u pozadini provjerava svaku spremljenu vezu: 🟢 radi (zadržite miš za latenciju i SSH verziju), 🟡 sporo, 🔴 ne radi (zadržite miš za razlog: odbijeno, istek vremena, nepoznat host). Mapa pokazuje koliko ih ne radi; crvena traka pokazuje samo servere koji ne rade. Obavijesti na radnoj površini uključite u **Postavke → Općenito**.

## Prozor Nadzor
**📊 Nadzor** u gornjoj traci:
- **Serveri:** stanje i povijest vaših veza, s metrikama iz Zabbix agenta: čitano kroz SSH, kroz mrežni put ili iz API-ja postojećeg Zabbix servera (*Postavke nadzora* po serveru; pomaže *Prepoznaj OS i instaliranog agenta*). *Provjeri sada* čita odmah.
- **Servisi i funkcionalnost:** Supervisor workeri, neuspjele systemd jedinice, disk i memorija po serveru te health URL-ovi ili TCP provjere servisa iz kataloga.
- **Logovi:** preglednik logova (vidi *Logovi*).

## Upozorenja izvan WRM-a
**Postavke → Obavijesti** bira što primate (server ne radi/radi, dospjela rotacija, nove verzije, stanje lanca…) i gdje: e-mail, Telegram, Slack/Mattermost, Teams, Discord, ntfy, Gotify, Pushover ili potpisani webhook. Administratori dodaju kanale u **Admin → Obavijesti**.

## Pitajte AI
U prozoru servera **AI** odgovara na pitanja kao *zašto je opterećenje visoko?* čitajući procese, servise i logove (samo čita, osim ako ne odobrite promjenu).
