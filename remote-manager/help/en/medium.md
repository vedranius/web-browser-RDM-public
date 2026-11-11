---
id: medium
title: Medium deploys
section: Operations
order: 88
summary: Deploy projects described in a deploy-config (servers, environments, projects) with WRM's native engine.
keywords: medium, deploy, deploy-config, yaml, environment, project, dry run, lock, backup, rollback, history, build, engine
ui: .medium-wrap, .md-engine
open: git
---
**⎇ Git → Medium** deploys projects that are described in a `deploy-config` (`servers.yaml`, `environments.yaml`, `projects/*.yaml`).

## Set up
1. Import your **deploy-config ZIP** (or link its Git repository with `DEPLOY_CONFIG_REPO`). WRM checks it at once and names the file and key of an error.
2. **Servers** are matched to your saved SSH connections with the same address and user, or linked by hand. Keys stay in WRM.
3. Per project choose where builds run: on the WRM host, on a Linux build server over SSH, or no build.

## Deploy
1. Choose project, environment and branch or tag; **Compare** shows what differs on every server.
2. **Dry run** shows the plan without changing anything.
3. **Deploy** takes a server-side lock, makes a backup, transfers, checks and records history. **Rollback** returns to an earlier deploy.

Some steps need step-up re-authentication (your password again). Medium runs build and post-deploy commands from the config, so only trusted operators should use it (policy `medium_native`).
