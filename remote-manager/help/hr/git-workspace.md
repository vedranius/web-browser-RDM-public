---
id: git-workspace
title: Git radni prostor
section: Rad
order: 85
summary: Vidite koja verzija kojeg servisa radi gdje, usporedite servere s GitLabom/GitHubom/Giteom, ažurirajte, deployajte i sigurno vratite.
keywords: git, gitlab, github, gitea, verzija, deploy, ažuriranje, nadogradnja, vraćanje, rollback, usporedba, diff, servis, katalog, instalacija, okruženje, izdanje, bundle, gitignore, ci, webhook, jednostavni način, version, update, upgrade, compare, environment, release, simple mode
ui: #btn-git, .git-ws, .gw-wrap
open: git
---
**⎇ Git** u gornjoj traci otvara radni prostor. Serverima ne treba git: WRM čita API Git pružatelja i radi preko njihovih SSH veza.

## Početak (kartica Postavke radnog prostora)
1. **Dodajte Git izvor:** GitLab, GitHub ili Gitea s **tokenom samo za čitanje**. *Testiraj* navodi projekte koje vidi.
2. **Katalog:** koji su projekti servisi, u kojoj su mapi na serverima, koje se datoteke prate ili štite (datoteke po hostu).
3. WRM pronalazi instalacije na vašim SSH serverima i uspoređuje ih **datoteku po datoteku**.

## Jednostavno ili napredno
- **Jednostavno:** jedna kartica po servisu i okruženju: što radi na kojem serveru, što je najnovije na grani i jedan gumb **Ažuriraj** s pregledom.
- **Napredno:** Pregled, Servisi, Instalacije, Okruženja, Pokretanja, Aktivnost, .gitignore i Postavke. Svaka kartica ima kratko objašnjenje ispod zaglavlja.

## Sigurne promjene na serverima
- **Ažuriraj / Nadogradi… / Vrati…:** pokretanja rade instalaciju za instalacijom i staju na prvoj grešci. Stare i nedostajuće datoteke su označene; datoteke *promijenjene ručno* zapisuju se samo ako ih označite; zaštićene nikad. Sigurnosne kopije idu u `.deploy-bak/`.
- **Produkcija** (`env:prod`) traži da upišete ime servera.
- **Okruženja:** poredani popis servera s pravilima; deploy je pregledani plan, probno pokretanje, zaključavanje na serveru, povijest i vraćanje.
- **Instaliraj… / Prenesi…** postavljaju servis na novi server.

## Bez mreže i CI
Offline paketi (`.tar.gz`) rade bez pristupa Gitu. Neobavezni webhookovi iz GitLaba/GitHuba pokreću trenutnu provjeru.

**Medium** (alat za deploy temeljen na deploy-configu) ima svoj članak.
