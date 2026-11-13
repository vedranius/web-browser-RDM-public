---
id: users-roles-teams
title: Korisnici, uloge i timovi
section: Administracija
order: 150
summary: Napravite račune, dodijelite uloge, složite timove, tražite ponovnu provjeru za opasne radnje i čuvajte break-glass administratora.
keywords: korisnik, korisnici, račun, novi korisnik, uloga, administrator, admin, operater, revizor, tim, timovi, dozvola, pravilo, step-up, break-glass, onemogući, obriši, reset lozinke, reset 2fa, registracija, user, account, role, auditor, team, permission, policy, disable, delete, reset password, registration
ui: #modal-admin, #modal-user-new, #modal-team, #modal-user-delete
open: admin:users
admin: true
---
**Postavke → 🛡 Admin panel → Korisnici** (administratori).

## Računi
- **Novi korisnik:** ime, uloga, neobavezni e-mail (poveznica za aktivaciju kad je pošta postavljena) ili privremena lozinka koja se mora promijeniti pri prvoj prijavi.
- **⋯** na korisniku: promjena uloge, reset lozinke ili 2FA, odjava svugdje, otključavanje, onemogućavanje, brisanje (sa svime što korisnik posjeduje ili predano drugom računu).
- Registracija je zatvorena nakon prvog računa, osim ako je ne otvorite (*Pravila*).

## Uloge
- **Administrator:** sve.
- **Operater:** funkcije koje pravila ograničavaju na operatere (npr. deployevi).
- **Korisnik:** vlastite veze i ono što je s njim podijeljeno.
- **Sigurnosni revizor:** čita revizijski zapis, sesije i izmjene datoteka; ne koristi funkcije.

## Timovi
**Timovi** grupiraju ljude; pravila mogu navesti `teams:<ime>`, a dijeljenja se mogu dati timu.

## Pravila
**Admin → Pravila** odlučuje tko smije koristiti koju funkciju (isključeno, administratori, operateri, timovi ili svi) i postavlja ograničenja. Pravilo nametnuto varijablom okruženja `WRM_<KLJUČ>` prikazuje se kao zaključano.

## Step-up i break-glass
Opasne radnje (pravila, uloge, tajne) ponovno traže lozinku i 2FA (*step-up*). **Break-glass** administrator se može prijaviti i tijekom lockdowna ili kad prijava preko imenika ne radi; čuvajte jednog s jakom lozinkom i 2FA.
