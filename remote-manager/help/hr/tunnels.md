---
id: tunnels
title: SSH tuneli i web sučelja
section: Pristup
order: 60
summary: Lokalno, udaljeno i SOCKS prosljeđivanje portova, panel Tuneli i veze na web sučelja iza servera.
keywords: tunel, prosljeđivanje portova, port forwarding, ssh -L, ssh -R, ssh -D, socks, lokalni, udaljeni, dinamički, web sučelje, https, idrac web, switch web, preglednik, tunnel, local, remote, dynamic, web interface, browser
ui: #modal-tunnels, #btn-tunnels, #tun-help
open: tunnels
---
## Vrste tunela
- **Lokalni (-L):** port na WRM računalu vodi do `host:port` kako ga vidi server. Primjer: `127.0.0.1:2443 → 10.0.0.5:443` otvara web stranicu switcha iza servera.
- **SOCKS (-D):** proxy za vaš preglednik ili alate koji dolazi do svega do čega dolazi server.
- **Udaljeni (-R):** port na serveru vodi natrag do nečega do čega dolazi WRM (treba pravilo `tunnel_remote_forward`).

## Postavljanje i korištenje
1. **Uredi vezu → 🔀 SSH tuneli → Dodaj** (ili desni klik → *Dodaj tunel…*).
2. Pokrećite ih i zaustavljajte u **🔀 Tuneli** u gornjoj traci. Pokrenuti tuneli se sami ponovno spajaju.
3. **Otvori** otvara web tunel u novoj kartici preglednika; gumb za kopiranje kopira adresu.

Adresa slušanja `127.0.0.1` (zadano) drži tunel na WRM računalu. Mrežna adresa (`0.0.0.0`, LAN IP) dopušta drugim računalima da ga koriste; postaviti je smiju samo administratori, osim ako pravilo `tunnel_bind_any` to dopušta. Ostavite port prazan da dobijete slobodan.

## Veze na web sučelja
Veza s protokolom **Web sučelje (HTTPS/HTTP)** dvoklikom otvara tu stranicu. S jump hostom WRM sam otvara privremeni tunel i zaustavlja ga nakon 30 minuta bez prometa.
