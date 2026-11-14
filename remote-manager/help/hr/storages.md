---
id: storages
title: Spremišta (S3)
section: Rad
order: 90
summary: Pregledavajte S3-kompatibilne buckete pokraj servera, s podacima o cijeni zahtjeva po pružatelju.
keywords: s3, spremište, bucket, minio, hetzner, digitalocean, spaces, aws, backblaze, r2, objekt, učitaj, preuzmi, cijena, trošak, prefiks, storage, object, upload, download, cost, prefix
ui: #btn-storages, .stg-wrap, #stg-list
open: storages
---
**🗄 Spremišta** u gornjoj traci drži S3-kompatibilne buckete pokraj servera.

## Dodavanje spremišta
1. **＋ Dodaj**, odaberite **pružatelja** (Hetzner, DigitalOcean Spaces, Amazon S3, Backblaze B2, Cloudflare R2, MinIO ili drugi). Popunjava krajnju točku i kaže naplaćuje li taj pružatelj zahtjeve koje WRM šalje.
2. Bucket, regija, neobavezni **prefiks** (WRM ostaje unutar njega), pristupni ključ i tajna (šifrirano spremljeni, više se nikad ne prikazuju).
3. Prekidači za pisanje, učitavanje i brisanje te ograničenje veličine. *Testiraj vezu*.

## Rad s datotekama
Mape se učitavaju u cijelosti (s brojem objekata i ukupnom veličinom); pretraga, sortiranje po imenu, veličini ili datumu, označavanje i skupno brisanje, učitavanje, preuzimanje, pregled. **📡** pokazuje S3 zahtjeve koje je WRM poslao, pa vidite koliko je koštao svaki klik.

## Troškovi
Popisivanje i čitanje su zahtjevi. Pružatelji kao Hetzner i DigitalOcean uključuju zahtjeve u cijenu; AWS S3 i drugi naplaćuju po zahtjevu i za promet. Biranje pružatelja pokazuje trenutno pravilo i vodi na stranicu s cijenama.

AI aplikacije spojene preko MCP-a također mogu čitati i (ako je dopušteno) pisati u spremišta, s istim prekidačima.
