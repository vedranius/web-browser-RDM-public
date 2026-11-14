---
id: remote-desktop
title: Udaljena radna površina, konzole i BMC
section: Osnove
order: 40
summary: RDP, VNC i Telnet u pregledniku, upravljanje izvan pojasa (iDRAC, iLO, IPMI) i serijske konzole.
keywords: rdp, vnc, telnet, udaljena radna površina, remote desktop, windows, guacd, guacamole, bmc, idrac, ilo, ipmi, redfish, napajanje, uključi, isključi, serijski, konzola, kvm, power, serial, console
ui: .guac-wrap, #modal-bmc
open: new_connection
---
## RDP, VNC i Telnet
Napravite vezu s protokolom RDP, VNC ili Telnet i dvaput je kliknite. Zaslon se otvara u WRM prozoru pokraj terminala; tipkovnica, miš, međuspremnik i (za RDP) zvuk prolaze. I ovdje radi jump host.

Ovi protokoli trebaju **guacd** (Apache Guacamole proxy) uz WRM:
- Docker/compose: već je uključen.
- Linux: instalirajte paket `guacd` ili pokrenite njegov kontejner; Windows i macOS: `wrm -setup-guacd` ga pokreće u Docker Desktopu.
Ako se radna površina ne otvori, prozor kaže nedostaje li guacd ili je server odbio prijavu.

## Upravljanje izvan pojasa (BMC)
U vezi: **Uredi → ⚡ Upravljanje izvan pojasa**. Odaberite **Redfish** (iDRAC, iLO, XClarity, Supermicro, OpenBMC) ili **IPMI**, adresu BMC-a i prijavu (ili dijeljenu vjerodajnicu). Zatim desni klik na vezu → **Upravljanje**: uključivanje, isključivanje i ponovno pokretanje, stanje, inventar i **serijska konzola**, i kad operacijski sustav ne odgovara. WRM po potrebi dolazi do BMC-a kroz jump host veze.

## Serijski portovi
Administratori mogu dopustiti serijske portove WRM računala (pravilo `serial_ports`) za konzole switcheva i routera.
