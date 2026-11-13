---
id: troubleshooting
title: Rješavanje problema
section: Pomoć
order: 190
summary: Što napraviti kad veza, radna površina, AI, uvoz ili sam WRM ne radi.
keywords: problem, greška, ne radi, pokvareno, ne mogu se spojiti, istek vremena, odbijeno, ključ hosta, permission denied, sporo, zamrznuto, rušenje, log, debug, pomoć, error, fails, not working, broken, cannot connect, timeout, refused, host key, slow, freeze, crash, help
ui: .toast.error
open: help:troubleshooting
---
**Najprije:** pritisnite **?** na poruci o grešci ili otvorite Pomoć i pitajte AI pomoć s tom porukom; objašnjava uobičajene uzroke.

## Server se ne spaja
- *Connection refused / timeout:* provjerite adresu i port, vatrozid i treba li serveru **jump host**. **🧰 Alati → Provjera porta** s WRM-a (ili s drugog servera) pokazuje gdje staje.
- *Host key changed:* server je ponovno instaliran ili je nešto između. Ako znate da je ponovno instaliran, zaboravite stari ključ u **Admin → Ključevi hostova** (ili prihvatite novi kad pita).
- *Permission denied:* krivi korisnik, lozinka ili ključ; s ključem provjerite **Tko ima pristup**.

## Udaljena radna površina ostaje crna
guacd nedostaje ili ne može doći do servera. Prozor navodi uzrok; vidi *Udaljena radna površina*.

## AI ne odgovara
Otvorite **✦ AI → 🧭 Postavi AI**: svaki korak ima svjetlo i kaže što nedostaje (AI isključen, nema servisa, neuspjeli test, prekidač za isključivanje, lockdown). Greške pružatelja (krivi ključ, nema kredita, nepoznat model) prikazuju se onako kako ih je servis poslao.

## Uvoz propušta unose
Pregled uvoza navodi svaki preskočeni unos s razlogom. Za mRemoteNG provjerite je li datoteka `confCons.xml` i upišite glavnu lozinku ako je ima.

## WRM je spor ili visi
Ponovno učitajte stranicu. Administratori: kartica Pregled pokazuje otvorene terminale, tunele i AI sesije; log servera (`journalctl -u wrm` ili konzola) pokazuje greške. Prijavite grešku na GitHubu s verzijom i retcima loga.
