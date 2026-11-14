---
id: jump-hosts-proxies
title: Jump hostovi i proxyji
section: Pristup
order: 55
summary: Dođite do servera iza bastiona (kao ssh -J) ili kroz SOCKS/HTTP proxy, i kao zadano za mapu.
keywords: jump host, bastion, skok, proxyjump, ssh -J, preko, proxy, socks, socks5, http connect, zadano za mapu, iza, interna mreža, via, behind, internal network, folder default
ui: #m-jump-help, #m-proxy-help, #modal-proxy
open: vault:proxies
---
## Jump host (bastion)
1. Spremite bastion kao običnu SSH vezu.
2. U vezi internog servera odaberite **Jump host (spoji preko) → bastion**.
3. Dvaput je kliknite: terminal piše *Connecting to 10.0.0.21 via bastion…*. Bastioni mogu imati svoj jump host (lanci); datoteke, prijenosi i tuneli koriste isti put.

## Proxy
Za mreže do kojih se dolazi kroz SOCKS5/4 ili HTTP CONNECT proxy:
1. **🔑 Ključevi → 🌐 Proxyji → Novi proxy** (ili *＋ Novi proxy…* u polju Proxy veze).
2. Adresa, port, neobavezna prijava (ili vjerodajnica), *DNS na proxyju*.
3. Odaberite ga u vezi. Proxy može biti i SOCKS tunel (-D) jedne od vaših veza; WRM ga pokreće kad zatreba.

## Zadano za mapu
Uredite mapu i postavite njezin zadani jump host i proxy: svaka veza u njoj (i nove) ih koristi, osim ako ne odabere svoje.

RDP/VNC radne površine i web sučelja također idu kroz jump host; IPMI koristi UDP pa se izvodi na zadnjem jump hostu.
