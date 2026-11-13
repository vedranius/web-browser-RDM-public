---
id: medium
title: Medium deployevi
section: Rad
order: 88
summary: Deployajte projekte opisane u deploy-configu (serveri, okruženja, projekti) WRM-ovim ugrađenim mehanizmom.
keywords: medium, deploy, deploy-config, yaml, okruženje, projekt, probno pokretanje, zaključavanje, sigurnosna kopija, vraćanje, povijest, build, mehanizam, environment, project, dry run, lock, backup, rollback, history, engine
ui: .medium-wrap, .md-engine
open: git
---
**⎇ Git → Medium** deploya projekte opisane u `deploy-config` (`servers.yaml`, `environments.yaml`, `projects/*.yaml`).

## Postavljanje
1. Uvezite svoj **deploy-config ZIP** (ili povežite njegov Git repozitorij s `DEPLOY_CONFIG_REPO`). WRM ga odmah provjerava i navodi datoteku i ključ greške.
2. **Serveri** se povezuju s vašim spremljenim SSH vezama iste adrese i korisnika ili ručno. Ključevi ostaju u WRM-u.
3. Po projektu odaberite gdje se radi build: na WRM hostu, na Linux build serveru preko SSH-a ili bez builda.

## Deploy
1. Odaberite projekt, okruženje i granu ili tag; **Usporedi** pokazuje što se razlikuje na svakom serveru.
2. **Probno pokretanje** pokazuje plan bez ikakve promjene.
3. **Deploy** uzima zaključavanje na serveru, radi sigurnosnu kopiju, prenosi, provjerava i bilježi povijest. **Vrati** se vraća na raniji deploy.

Neki koraci traže ponovnu provjeru identiteta (ponovno vašu lozinku). Medium pokreće build i post-deploy naredbe iz konfiguracije, pa ga trebaju koristiti samo operateri od povjerenja (pravilo `medium_native`).
