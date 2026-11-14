---
id: connections
title: Veze i mape
section: Osnove
order: 20
summary: Dodajte, uredite, organizirajte, premjestite i pretražite servere; protokoli, oznake i bilješke.
keywords: veza, konekcija, server, dodaj, nova, uredi, mapa, premjesti, povuci, redoslijed, sortiraj, oznaka, tag, pretraga, ssh, rdp, vnc, telnet, web, brzo spajanje, bilješke, status, connection, folder, quick connect
ui: #sidebar, #modal-conn, #modal-folder, #modal-move, #modal-quick, #modal-notes
open: new_connection
---
Veza je jedan spremljeni server: protokol, adresa, prijava i opcije. Veze su u mapama u lijevoj bočnoj traci.

## Dodavanje veze
1. **+ Veza** u bočnoj traci.
2. Odaberite **protokol**: SSH (terminal i datoteke), SFTP/FTP/FTPS (samo datoteke), RDP, VNC, Telnet (udaljena radna površina) ili **Web sučelje** (HTTP/HTTPS stranica switcha, iDRAC-a, NAS-a…).
3. Upišite **host**, **port** i **korisnika**. Za prijavu odaberite lozinku, ključ iz **🔑 Ključeva**, dijeljenu **vjerodajnicu** ili Bitwarden SSH ključ.
4. Neobavezno: **mapa**, **oznake** (`env:prod`, `site:dc1`), **jump host** ili **proxy**, SSH **tuneli**, **bilješke**, upravljanje izvan pojasa (BMC).
5. **Spremi**, pa dvaput kliknite vezu da je otvorite.

## Brzo spajanje
Upišite `korisnik@host:port` (ili `rdp://…`, `vnc://…`, `https://…`) u **⚡ Brzo spajanje** i pritisnite Enter. Čuva se 24 sata pod *Brze veze*; 💾 je trajno sprema.

## Mape i redoslijed
- **📁 Mapa** stvara mapu; desni klik na mapu za podmapu, preimenovanje, dijeljenje, sortiranje A–Ž ili *uvijek sortirano*.
- **Povucite** veze i mape da ih premjestite; spustite između dva retka da je stavite točno tamo. Ctrl+klik označava više njih, pa ih povucite zajedno ili desni klik → **Premjesti u…**.
- Desni klik na vezu: otvori, uredi, dupliciraj, podijeli, mrežni alati, tuneli, bilješke, obriši.

## Pronalaženje servera
- Okvir za **pretragu** filtrira po imenu, hostu, korisniku i oznakama (`#env:prod`).
- **🏷 Traka oznaka** filtrira po oznakama; veze s `env:prod` dobivaju crvenu oznaku **PROD**.
- Svjetlo pokraj svake veze: 🟢 radi, 🟡 sporo, 🔴 ne radi (zadržite miš za razlog). Crvena traka navodi samo servere koji ne rade.

## Bilješke
📝 na vezi čuva procedure, kontakte i poveznice (jednostavni Markdown). Prikazuju se i u prozoru terminala.
