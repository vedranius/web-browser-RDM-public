---
id: terminal
title: SSH terminal, snippeti i broadcast
section: Osnove
order: 30
summary: Rad u terminalima; pretraga, ponovno spajanje, snippeti, pokretanje pri spajanju, zabilješke mapa i tipkanje u više servera.
keywords: terminal, ssh, konzola, ljuska, shell, snippet, naredba, broadcast, više servera, cssh, pokreni pri spajanju, zabilješka, pretraga, kopiraj, zalijepi, font, ponovno spajanje, snimanje, command, bookmark, recording
ui: .term-wrap, .xterm, #modal-snippet, #modal-broadcast, #modal-bookmark, #snip-list
open: snippets
---
Dvaput kliknite SSH vezu da otvorite terminal. To je potpuni xterm (256 boja, miš, UTF-8) i radi s `vim`, `htop`, `mc` i `tmux`.

## Svakodnevne tipke
- **Ctrl+Shift+F** pretraga ispisa terminala; **Ctrl+Shift+C / V** kopiranje i lijepljenje (označeni tekst se kopira u panel Međuspremnik).
- **Enter** u zatvorenom terminalu ponovno se spaja; **↻** u naslovnoj traci to radi bilo kada.
- **Ctrl+Shift+Space** (ili ⚡ u naslovnoj traci) otvara biranje snippeta.

## Snippeti
Spremljene naredbe s varijablama kao `{{host}}`, `{{user}}` ili `{{?Service=nginx}}` (pita se pri korištenju).
1. Panel **Snippeti** → **＋**, ili označite tekst u terminalu i pritisnite **＋ Novi snippet** u biranju.
2. Kliknite snippet da ga pokrenete u terminalu koji ima fokus; **Shift+klik** ga umeće bez Entera.
3. Doseg: sve veze, jedna mapa ili jedna veza. Snippeti **Pokreni pri spajanju** izvode se nakon prijave (npr. `sudo -i`).

## Zabilješke mapa
★ u naslovnoj traci radi `cd` u spremljeni direktorij kao `/opt/app/{name}/logs`; u upravitelju datoteka ga otvara.

## Broadcast unos
Tipkajte isto u više terminala: **🧰 Alati → 📣 Broadcast…**, označite terminale, **Pokreni**. Opasne naredbe (`rm -rf`, `shutdown`, `mkfs`…) i zalijepljeni višeredni tekst najprije pitaju. **■ Zaustavi** završava.

## AI u terminalu
**AI** u naslovnoj traci otvara asistenta pokraj ovog terminala (vidi *AI asistent*). Da pitate o nečemu na zaslonu bez sesije na serveru, koristite **?** → **⌖ Odaberi** i kliknite terminal.

## Snimanje
Sesije se snimaju kad to pravilo kaže; reproducirajte ih u **Postavke → Povijest sesija**.
