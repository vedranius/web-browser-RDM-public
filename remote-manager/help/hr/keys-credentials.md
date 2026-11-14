---
id: keys-credentials
title: SSH ključevi, vjerodajnice i Bitwarden
section: Pristup
order: 50
summary: Generirajte i rasporedite SSH ključeve, dijelite prijave bez otkrivanja, rotirajte lozinke, koristite certifikate i Bitwarden.
keywords: ključ, ssh ključ, ed25519, rsa, rasporedi ključ, ssh-copy-id, vjerodajnica, lozinka, trezor, dijeli prijavu, rotacija, tko ima pristup, certifikat, opkssh, bitwarden, passphrase, ppk, key, credential, password, vault, rotate, certificate
ui: #modal-vault, #modal-keynew, #modal-keydeploy, #modal-cred, #modal-rotate, #modal-authkeys, #btn-vault
open: vault
---
**🔑 Ključevi** u gornjoj traci otvara trezor. Sve u njemu je šifrirano spremljeno na WRM serveru i nikad se ne šalje natrag u preglednik.

## SSH ključevi
1. **Generiraj ključ** (preporuka ED25519) ili **Uvezi ključ** (OpenSSH ili PEM; passphrase se traži jednom). PuTTY `.ppk`: najprije ga u PuTTYgenu izvezite kao OpenSSH ključ.
2. **Rasporedi**: odaberite veze i WRM dodaje javni ključ u njihov `authorized_keys` (kao `ssh-copy-id`), i kroz jump hostove.
3. U vezi odaberite **Prijava: ključ iz spremišta ključeva**.
4. **Tko ima pristup** na SSH vezi navodi ključeve u njezinom `authorized_keys` i koji su od njih vaši.

## Vjerodajnice
Vjerodajnica je jedna prijava (korisnik + lozinka ili ključ) koju koristi više veza, npr. `root@dc1`. Promijenite je jednom i sve veze koriste novu.
- **Podijelite** je s kolegama: mogu je koristiti, ali je nikad ne vide.
- **Dopušteni hostovi** (`*.dc1.example.com`, `10.0.0.0/8`) sprječavaju slanje dijeljene lozinke na tuđi server.
- **Rotacija:** *Provjeri* se prijavljuje svugdje bez ikakve promjene; *Rotiraj* postavlja novu lozinku na svim njezinim SSH serverima i vraća staru ako jedan ne uspije.

## Certifikati i Bitwarden
- OpenSSH korisnički certifikat (i iz `opkssh login`) može se spremiti uz ključ: **📜 Certifikat**.
- **Bitwarden:** prijavite se u svoj trezor u **Ključevi → Bitwarden**, otključajte ga i u vezi odaberite *Bitwarden SSH ključ*; ključ se dohvaća na serveru samo pri prijavi.
