---
id: help-and-ai-help
title: Korištenje Pomoći i AI pomoći
section: AI
order: 115
summary: Pretražite pomoć, odaberite bilo koji dio zaslona da pitate o njemu i neka vam AI pomoć objasni ili napravi.
keywords: pomoć, f1, upitnik, pretraga, odaberi, označi, okvir, objasni, kako, pitaj ai, napravi umjesto mene, dokumentacija, priručnik, vodič, help, question mark, search, pick, select, frame, explain, how to, ask ai, documentation, manual, guide
ui: #help-panel, #btn-help
open: help:help-and-ai-help
---
**?** u gornjoj traci ili **F1** otvara Pomoć sa strane; WRM radi dalje pokraj nje, a dijalog može ostati otvoren dok čitate.

## Pronalaženje odgovora
- **Pretraga** na engleskom ili hrvatskom, sa ili bez č, ć, š, ž, đ.
- **Teme** navode svaki članak po odjeljcima; svaki ima i poveznicu na istu stranicu na GitHubu.
- Pomoć se otvara na članku o onome što je pred vama (npr. dijalog koji imate otvoren).

## Odabir dijela zaslona
1. Pritisnite **⌖ Odaberi** u Pomoći.
2. Pokažite gumb, dijalog, prozor ili popis: dobiva okvir. Kliknite da ga odaberete ili **povucite okvir** oko područja. **Esc** odustaje.
3. Pomoć pokazuje članke o tome i stavlja ga u AI chat kao kontekst, pa možete pitati *što ovo radi?* ili *zašto je ovo crveno?* Šalje se samo vidljivi tekst, nikad vrijednosti polja.

## Pitajte AI pomoć
Upišite pitanje svojim riječima, npr. *kako doći do servera iza bastiona?* Odgovor dolazi iz ovih članaka, s putanjama izbornika i koracima.

Ako tražite da **napravi umjesto vas**, nudi gumbe kao *Otvori novu vezu (popunjenu)*, *Otvori Postavke → AI* ili *Otvori app-01 s AI asistentom*. Svaki gumb točno kaže što radi i **ništa se ne događa dok ga ne pritisnete**. Promjena postavke pita još jednom prije spremanja, a promjene na serveru i dalje idu kroz odobrenja AI asistenta.

AI pomoć treba AI servis (vidi *Postavi AI*); pretraga i članci rade i bez njega. Administratori biraju tko je smije koristiti pravilom `ai_help` (zadano: isti ljudi kao za AI asistenta). Pitanja se bilježe u revizijski zapis (`help.ask`), bez njihova teksta.

Poruke o greškama također imaju gumb **?**: otvara Pomoć s upisanom greškom.
