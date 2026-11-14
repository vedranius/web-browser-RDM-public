---
id: security
title: Razine sigurnosti i stanje
section: Administracija
order: 155
summary: Odaberite sigurnost Home, Business ili Enterprise, pogledajte odstupanja i jednim klikom popravite slabe postavke.
keywords: sigurnost, razina, home, business, enterprise, stanje, https, učvršćivanje, odstupanje, upozorenja, obavezni 2fa, pravila lozinki, ključevi hosta, csp, security, level, posture, hardening, deviation, warnings, 2fa required, password policy, host keys
ui: #modal-admin
open: admin:security
admin: true
---
**Admin panel → Sigurnost**.

## Razina sigurnosti
- **Home:** razumne zadane postavke za jednu osobu ili laboratorij.
- **Business:** HTTPS za prijavu, 2FA za administratore, strože lozinke i sesije, snimke.
- **Enterprise:** 2FA za sve, kratke sesije, step-up, revizijski zapis u SIEM.
Odaberite je pri prvom pokretanju, s `wrm -security-level …` ili ovdje (WRM pokaže plan prije nego išta promijeni). `WRM_SECURITY_LEVEL` je nameće.

## Stanje
Stranica navodi svaku postavku slabiju od onoga što razina traži (**odstupanje**) s gumbom *Popravi*. Odstupanje koje namjerno zadržite traži bilješku, koja se bilježi u revizijski zapis. Kartica Pregled također pokazuje sigurnosna upozorenja (nema HTTPS-a, otvorena registracija, administratori bez 2FA…).

## Dobra praksa
- Pokrenite WRM iza HTTPS-a (ugrađeni certifikat ili reverse proxy).
- Radite sigurnosnu kopiju baze **i** `remote_manager.db.key` odvojeno.
- Ostavite provjeru ključeva hosta uključenom i koristite vjerodajnice s dopuštenim hostovima.
- Koristite uloge, timove i dijeljenja umjesto dijeljenja lozinki.
