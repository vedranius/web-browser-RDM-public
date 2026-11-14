---
id: audit-siem
title: Revizijski zapis, snimke, SIEM i lockdown
section: Administracija
order: 165
summary: Tko je što i kada napravio, snimke sesija, izmjene datoteka, izvoz u SIEM, sigurnosna upozorenja i hitni lockdown.
keywords: revizija, revizijski zapis, audit, log, povijest, tko, snimka, reprodukcija, sesija, izmjena datoteke, prijenos, siem, syslog, splunk, elastic, upozorenje, lockdown, zaustavi korisnika, prekidač, integritet, lanac hasheva, csv, recording, replay, file edit, alert, stop user, kill switch, integrity, hash chain
ui: #modal-admin, #modal-replay, #modal-siem
open: admin:audit
admin: true
---
## Revizijski zapis
**Admin → Revizijski zapis:** prijave, promjene računa i pravila, veze, dijeljenja, terminalske sesije, prijenosi i izmjene datoteka, AI radnje, pitanja AI pomoći. Filtrirajte po tekstu, događaju, korisniku i datumu; **CSV izvoz**; **Provjeri integritet** provjerava lanac hasheva (unosi se ne mogu mijenjati ni brisati prije isteka vremena čuvanja).

## Sesije i snimke
**Admin → Sesije i snimke:** svaka terminalska sesija i sesija radne površine, s reprodukcijom. Korisnici vide svoje u *Postavke → Povijest sesija*. **Admin → Izmjene datoteka** pokazuje svaku promijenjenu datoteku s diffom.

## Izvoz u SIEM i upozorenja
**Admin → Sigurnost → Izvoz u SIEM:** syslog, HTTP (Splunk HEC, Elastic i drugi) s *Testom*. **Sigurnosna upozorenja** obavještavaju administratore o valovima neuspjelih prijava, novim administratorima, lockdownu i slično.

## Kad nešto nije u redu
- **Zaustavi korisnika:** odjavljuje osobu svugdje, završava njezine terminale, tunele i AI sesije i onemogućava račun.
- **Lockdown:** prijaviti se mogu samo break-glass administratori, AI staje, a dijeljenja se pauziraju (`wrm -lockdown on` radi i sa servera).
- **AI prekidač za isključivanje:** *Admin → AI asistent*.
