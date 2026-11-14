---
id: import
title: Uvoz i izvoz
section: Osnove
order: 25
summary: Uvezite veze iz mRemoteNG-a, PuTTY-ja, OpenSSH-a, CSV/Excela, NetBoxa ili drugog WRM-a.
keywords: uvoz, izvoz, prijenos, migracija, sigurnosna kopija, import, export, mremoteng, confcons, putty, ssh config, csv, excel, xlsx, netbox, backup, migrate
ui: #st-data, #modal-inventory
open: settings:st-data
---
## Iz drugih programa (Postavke → Podaci)
- **mRemoteNG:** odaberite `confCons.xml` (`%APPDATA%\mRemoteNG\confCons.xml` ili *File → Export* u mRemoteNG-u). Kontejneri postaju mape, SSH/RDP/VNC/Telnet/HTTP veze zadržavaju opcije, a lozinke se dešifriraju i šifrirano spremaju u WRM. Ako datoteka ima glavnu lozinku, WRM je traži. Provjerite pregled, pa uvezite.
- **PuTTY:** izvezite ključ registra `HKCU\Software\SimonTatham\PuTTY\Sessions` u `.reg` datoteku i uvezite je; proxyji dolaze s njom.
- **OpenSSH:** uvezite `~/.ssh/config`; unosi `Host` postaju veze, `ProxyJump` postaje jump host.

Ako se nešto ne uveze kako očekujete, pregled navodi svaki preskočeni unos i razlog. Pitajte AI pomoć s tom porukom; može objasniti uzrok.

## Inventar (CSV, Excel, NetBox)
**Postavke → Podaci → Inventar** (ili izbornik ⋯ u bočnoj traci):
- **CSV / Excel (.xlsx):** prvi redak sadrži nazive stupaca; WRM pogađa koji je stupac ime, host, port, korisnik, mapa i oznake. Provjerite mapiranje u pregledu.
- **NetBox:** adresa API-ja i token samo za čitanje. Lokacije, rackovi, uloge i tenanti postaju oznake; ponovni uvoz ažurira ono što je prije uvezeno.

## Izvoz i sigurnosna kopija
- **Izvezi veze** zapisuje veze, mape, jump hostove i tunele u datoteku bez tajni. **Izvoz s lozinkama i ključevima** ih uključuje; traži vašu lozinku i pravilo to mora dopuštati.
- **Uvoz** WRM izvoza vraća ga na drugom WRM-u ili nakon ponovne instalacije.
- Administratori rade sigurnosnu kopiju cijelog WRM-a: datoteku baze i njezin ključ za šifriranje (vidi *Ažuriranja i sigurnosne kopije*).
