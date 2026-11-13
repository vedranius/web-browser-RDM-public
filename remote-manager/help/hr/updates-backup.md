---
id: updates-backup
title: Ažuriranja, sigurnosne kopije i instalacija
section: Administracija
order: 170
summary: Potpisana ažuriranja jednim klikom, što spremati u sigurnosnu kopiju, WRM kao servis, u Dockeru i iza reverse proxyja.
keywords: ažuriranje, nadogradnja, verzija, izdanje, samoažuriranje, potpis, sigurnosna kopija, vraćanje, baza, ključ, servis, systemd, windows servis, docker, reverse proxy, nginx, https, certifikat, port, instalacija, update, upgrade, release, signature, backup, restore, database, service, install
ui: #modal-update, .app-ver
open: updates
admin: true
---
## Ažuriranje
Kliknite verziju u gornjoj traci. WRM pokazuje što je novo; administratori pritisnu **Ažuriraj**: izdanje se preuzima, provjerava mu se **potpis** i WRM se ponovno pokreće s njim. Docker instalacija se ažurira povlačenjem nove slike.

## Sigurnosna kopija
Čuvajte kopije obje datoteke jednu uz drugu:
- `remote_manager.db` (sve: veze, korisnici, postavke, povijest),
- `remote_manager.db.key` (ključ koji dešifrira spremljene lozinke i ključeve). **Bez njega su spremljene tajne izgubljene.**
Kopirajte ih dok je WRM zaustavljen ili koristite SQLite sigurnosnu kopiju po izboru. Prethodna verzija i dalje se pokreće na novijoj bazi.

## Pokretanje
- **Kao servis:** `wrm -install-service` (systemd, launchd, rc.d ili Windows servis).
- **Docker:** `docker compose up -d` s datotekom iz repozitorija (uključuje guacd).
- **Reverse proxy:** propustite WebSockete; postavite `ai_mcp_public_url` kad se vanjska adresa razlikuje.
README na GitHubu ima svaku opciju (`PORT`, `DB_PATH`, TLS, varijable okruženja).
