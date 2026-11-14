---
id: files
title: Datoteke, uređivač i prijenosi
section: Osnove
order: 35
summary: SFTP/FTP upravitelj datoteka, pretraga u datotekama, uređivač teksta s AI-jem, pregledi i kopiranje sa servera na server.
keywords: datoteke, sftp, ftp, ftps, prijenos, učitaj, preuzmi, uređivač, uredi, pretraga, traži, kopiraj, server na server, pregled, dozvole, chmod, zabilješka, files, upload, download, editor, transfer
ui: .fm-wrap, .fm-toolbar, #modal-editor, #modal-transfer, .ed-wrap
open: help:files
---
SFTP, FTP i FTPS veze otvaraju upravitelj datoteka; na SSH vezi koristite **📂 Datoteke** u naslovnoj traci terminala.

## Pregledavanje i upravljanje
- Kliknite dio **trake putanje** da skočite tamo ili kliknite prazno mjesto i upišite putanju. **↑** ide gore, **⟳** (F5) osvježava.
- Povucite datoteke s računala u prozor da ih **učitate**; desni klik → **Preuzmi** (mape kao ZIP).
- Desni klik za preimenovanje, brisanje, novu mapu ili datoteku, dozvole, kopiranje putanje, **Pregled** (slike, video, PDF, tekst) i **Uredi**.
- **★** sprema duboke mape po serveru, mapi, oznaci ili za sve servere.

## Pretraga u datotekama
Okvir za pretragu filtrira trenutnu mapu dok tipkate (`*.log`, `nginx*.conf`). **Enter** pretražuje sve podmape na serveru; rezultati stižu postupno i pretraga se može zaustaviti. ▾ dodaje pretragu sadržaja te filtre veličine i datuma.

## Uređivač teksta
Otvara tekstualne datoteke do 5 MB: isticanje sintakse, **Ctrl+F / Ctrl+H** traži i zamijeni, **Ctrl+G** idi na redak, *Provjeri JSON*, **Ctrl+S** spremi. Pri spremanju prepoznaje je li netko drugi u međuvremenu promijenio datoteku. **✦ AI** u uređivaču otvara asistenta za datoteke: pregledava ili objašnjava datoteku i predlaže promjene kao diff koji odobravate.

## Prijenos sa servera na server
**⚡ Prijenos** kopira datoteke i mape s jednog SSH/SFTP servera izravno na drugi, kroz WRM, a ne kroz vaš preglednik. Odaberite što se događa kad datoteka postoji (pitaj, prepiši, preskoči, prepiši ako je novija). Svaki prijenos je u revizijskom zapisu sa svojim SHA-256.
