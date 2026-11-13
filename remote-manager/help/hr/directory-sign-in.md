---
id: directory-sign-in
title: Prijava preko imenika (Entra ID, AD, LDAP, SCIM)
section: Administracija
order: 160
summary: Neka se ljudi prijavljuju s Microsoft Entra ID, Okta, Google ili bilo kojim OpenID Connect pružateljem, ili LDAP/Active Directoryjem, i dodjeljujte račune preko SCIM-a.
keywords: sso, jedinstvena prijava, oidc, openid, entra, azure ad, active directory, ad, ldap, adfs, okta, keycloak, google workspace, scim, grupe, domena, dodjela računa, single sign-on, groups, domain, provisioning
ui: #modal-sso, #modal-ldap, #modal-scim, #modal-sso-test
open: admin:sso
admin: true
---
**Admin panel → Prijava preko imenika**.

## OpenID Connect (Entra ID, ADFS, Okta, Keycloak, Google)
1. **Dodajte pružatelja** i odaberite vrstu; uređivač kaže gdje tamo napraviti klijenta.
2. Registrirajte prikazanu **adresu za preusmjeravanje** (`https://<WRM>/api/auth/oidc/callback`).
3. Upišite issuer, client ID i tajnu; **Testiraj vezu** provjerava svaki korak.
4. **Pravila** povezuju claim (`roles`, `groups`, `hd`…) s WRM ulogom i timovima. Bez pravila se mogu prijaviti samo povezani računi.

## LDAP / Active Directory
**Dodaj imenik:** vrsta (Active Directory ili drugi LDAP), `ldaps://` server (ili `ldap://` sa StartTLS-om), servisni račun koji smije čitati korisnike i grupe, baza korisnika i grupa te pravila od grupa do uloga. *Testiraj* prijavljuje korisnika po vašem izboru.

## SCIM
Entra ID ili Okta mogu automatski stvarati, mijenjati i onemogućavati WRM račune: **SCIM → Novi token**, pa adresu i token dajte imeniku.

Korisnici iz imenika se ne mogu prijaviti lokalnom lozinkom osim ako to dopustite; break-glass administrator uvijek može.
