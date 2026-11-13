---
id: network-tools
title: Mrežni alati
section: Rad
order: 95
summary: Provjera porta, ping, traceroute, DNS, HTTP/TLS provjere i više, s WRM-a ili s bilo kojeg vašeg servera.
keywords: alati, mreža, port, ping, traceroute, dns, http, tls, certifikat, provjera, dostupnost, vatrozid, firewall, skeniranje, nslookup, dig, curl, tools, network, certificate, check, reach, scan
ui: #modal-tools, #btn-tools
open: tools
---
**🧰 Alati** u gornjoj traci (ili desni klik na vezu → *Mrežni alati s ovog servera…*).

## Alati
- **Provjera porta:** do 100 portova (`22,80,443`, `8000-8010`): otvoren, zatvoren ili bez odgovora.
- **Ping** i **Traceroute.**
- **DNS:** A, AAAA, MX, TXT, CNAME i drugi.
- **HTTP / TLS:** status, zaglavlja, preusmjeravanja i certifikat s istekom.

## Odakle
Pokrenite alat **s WRM servera** ili **s jednog od vaših SSH servera** (*može li app-01 doći do baze na 5432?*). Sa servera WRM izvodi naredbe samo za čitanje preko SSH-a; jump hostovi vrijede.

Administratori odlučuju tko smije koristiti alate (pravilo `network_tools`). **📣 Broadcast** je u istom izborniku (vidi *SSH terminal*).
