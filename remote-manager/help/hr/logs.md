---
id: logs
title: Logovi
section: Rad
order: 75
summary: Pronađite svaki log na serveru, pregledajte ga i pratite, pretražite više logova odjednom i pitajte AI asistenta za logove.
keywords: log, logovi, zapisi, dnevnik, journal, journalctl, syslog, tail, prati, pretraga, grep, greška, otkrivanje, docker logs, supervisor, asistent za logove, follow, search, error, discovery
ui: .lv-wrap, .logviewer
open: monitoring:logs
---
**📊 Nadzor → Logovi**, pa odaberite server.

## Otkrivanje
WRM sam pronalazi logove preko SSH računa veze (bez sudo): systemd jedinice i njihov journal, Supervisor programe, Docker kontejnere i log datoteke ispod putanja za skeniranje. Datoteke s kodom i Python virtualna okruženja se preskaču. Popis se može sortirati i filtrirati po servisu, operateru, mapi i imenu.

## Pregled i praćenje
Kliknite datoteku (ili journal jedinicu) da je otvorite na kraju; **Prati** prikazuje nove retke uživo. Filtrirajte po razini (greške, upozorenja…), *Traži u ovom logu* (`/`, pa `n`/`N`), *Skoči na vrijeme*, *Učitaj starije retke*, kopirajte ili preuzmite prikazane retke ili *Otvori u terminalu*.

## Pretraga više logova odjednom
Označite datoteke (ili nijednu za sve) i tražite: WRM ih pretražuje na serveru (tekst ili regularni izraz, velika/mala slova, cijele riječi, retci konteksta, vremenski raspon) i prikazuje pogotke po datoteci s vremenskom crtom; kliknite stupac da zumirate, kliknite pogodak da otvorite to mjesto.

## AI asistent za logove
**✦ AI** u pregledniku logova otvara chat koji samo čita: *zašto je demo-api pao sinoć?* Pretražuje sve otkrivene logove, citira retke s datotekom i vremenom i predlaže što dalje provjeriti. Nikad ništa ne mijenja.

## Postavke logova
**Postavke** u pregledniku logova: putanje za skeniranje, datoteke koje se isključuju, koliko dugo otkrivanje vrijedi te za administratore ograničenja za sve.
