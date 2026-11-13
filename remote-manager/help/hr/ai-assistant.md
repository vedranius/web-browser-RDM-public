---
id: ai-assistant
title: AI asistent na serverima
section: AI
order: 105
summary: Pitajte o serveru običnim jezikom; asistent čita kroz WRM-ove alate i pita prije svake promjene.
keywords: ai, asistent, chat, pitaj, server, terminal, samo čitanje, odobri, odobrenje, odbij, automatski, način, model, pružatelj, logovi, asistent za datoteke, opencode, sesije, transkript, prekidač, assistant, ask, read only, approve, deny, automatic, mode, provider, sessions, transcript, kill switch
ui: .ai-panel, #btn-ai, .ai-head
open: ai_setup
---
Kliknite **AI** u naslovnoj traci SSH terminala ili upravitelja datoteka (ili **✦ AI** u gornjoj traci za zadnje korišteni server). Odaberite servis, model i način i pitajte: *zašto nginx ne radi?*, *oslobodi malo mjesta na disku*.

## Kako radi
Model nikad ne dobiva ljusku. Može samo pozivati WRM-ove alate (podaci o sustavu, čitanje datoteke, popis direktorija, logovi, journal, servisi, paketi, procesi, disk, mreža, pokretanje naredbe, pisanje i uređivanje datoteke). **Mehanizam dozvola na serveru** odlučuje o svakom pozivu:
- **Samo čitanje:** samo čita.
- **Pitaj:** svaka promjena pokazuje točnu naredbu ili diff; vi **Odobravate**, **Uređujete** ili **Odbijate**.
- **Automatski u granicama:** popis dopuštenog s ograničenjem vremena i broja radnji koje sami uključite.
Destruktivne naredbe se odbijaju osim ako ih administrator ne dopusti, a tajne u ispisu se maskiraju.

## Drugi asistenti
- **Asistent za logove** u pregledniku logova (samo čita), **asistent za datoteke** u uređivaču (predlaže diffove), **Lanci** (izrađuje nacrt lanca), **AI pomoć** (objašnjava WRM).
- **OpenCode agenti** mogu raditi na serveru u sandboxu, sa svakom radnjom kao WRM odobrenjem.

## Sesije i kontrola
**🧠 AI sesije** (izbornik Paneli) navodi vaše sesije: pretraga, otvaranje transkripta, nastavak, spremanje, izvoz. *Završi sve* ih zaustavlja. Administratori vide svaku sesiju, potrošnju i trošak u *Admin → AI asistent* i imaju **prekidač za isključivanje**.
