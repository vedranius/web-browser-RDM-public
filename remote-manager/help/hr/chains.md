---
id: chains
title: Lanci
section: Rad
order: 80
summary: Pratite cijeli radni sustav (servere, spremišta, krajnje točke i veze među njima) na jednom dijagramu uživo.
keywords: lanac, lanci, dijagram, tok, pipeline, ingest, cdn, primarni, rezervni, scenarij, stanje, provjera, ovisnost, topologija, izradi s ai, chain, flow, backup, state, check, dependency, topology
ui: #btn-chains, .chn-wrap, .chains
open: chains
---
Lanac prati sustav umjesto pojedinačnih servera, npr. ingest → transkodiranje → spremište → CDN s primarnim i rezervnim putem.

## Izrada lanca
**⛓ Lanci → Novi lanac**, zatim:
- **Izradi s AI-jem:** opišite sustav riječima; AI pogleda vaše veze i spremišta i napravi nacrt lanca. Ništa se ne nadzire dok ne pregledate nacrt na dijagramu i ne napravite ga.
- **Korak po korak:** odaberite čvorove (vaše SSH veze, spremišta, vanjske krajnje točke), dodajte provjere iz kataloga (proces, port, HTTP, starost datoteke, disk, uzorak u logu…), povežite ih i opišite **stanja** (npr. *radi na rezervi*) s uvjetima.
- Tu su i **predlošci** i napredni JSON uređivač.

## Korištenje
Dijagram prikazuje svaki čvor i vezu u boji s vremenom zadnje provjere. Kliknite čvor za njegove provjere i povijest; *Testiraj sada* jednom pokreće sve provjere. Obavijesti vam mogu javiti kad lanac promijeni stanje.

## AI na lancu
Iz prozora lanca pitajte AI da objasni stanje ili predloži promjenu; vraća nacrt koji sami primijenite.
