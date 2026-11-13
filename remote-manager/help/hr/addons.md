---
id: addons
title: Dodaci
section: Prilagodba
order: 120
summary: Proširite WRM dodacima (paneli, teme, kod na serveru), instalirajte s popisa zajednice ili izradite svoj s AI-jem.
keywords: dodatak, dodaci, plugin, proširenje, instaliraj, zip, zajednica, tema, panel, wasm, izradi dodatak, razvoj, dozvole, instanca, addon, extension, install, community, theme, develop, permissions, instance
ui: #st-addons, #addons-pane
open: addons
---
Dodatak je ZIP s `addon.json`: panel (prozor), teme boja i, za dodatke organizacije, WebAssembly kod na serveru. Pri instalaciji se ništa ne pokreće, a novi dodatak se instalira **isključen**.

## Instalacija
**Postavke → Dodaci**:
- **Zajednica** navodi objavljene dodatke s razinom povjerenja; preuzimanje se provjerava prema njegovom SHA-256.
- **Učitaj ZIP** instalira vaš. Administratori instaliraju za sve ili samo za sebe; ostali instaliraju **osobne** dodatke (pravilo `addon_personal`).
- **Uključi** pokazuje **list dozvola**: što dodatak smije (mrežna ishodišta, tajne, obavijesti, kod na serveru). Ažuriranje koje traži više ponovno čeka vaše odobrenje.

## Korištenje
*Otvori panel* ga otvara kao prozor na platnu (ili ga usidrite kao panel). **Instance** čuvaju zasebne postavke, npr. jedan *Server metrics* po klasteru.

## Izradite svoj
- **S AI-jem:** zamolite AI asistenta (ili AI aplikaciju preko MCP-a) da izradi dodatak; piše **nacrt** koji pregledate kao diff u Postavke → Dodaci i sami instalirate.
- **Ručno:** krenite od predloška (*Novi dodatak*), koristite razvojni server i `addon lint`; format je u `docs/ADDONS.md` na GitHubu.
