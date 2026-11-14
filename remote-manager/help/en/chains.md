---
id: chains
title: Chains
section: Operations
order: 80
summary: Watch a whole working system (servers, storages, endpoints and the links between them) on one live diagram.
keywords: chain, chains, diagram, flow, pipeline, ingest, cdn, primary, backup, scenario, state, check, dependency, topology, build with ai
ui: #btn-chains, .chn-wrap, .chains
open: chains
---
A chain watches a system instead of single servers, for example ingest → transcoding → storage → CDN with a primary and a backup path.

## Build a chain
**⛓ Chains → New chain**, then:
- **Build with AI:** describe the system in words; the AI looks at your connections and storages and drafts the chain. Nothing is monitored until you review the draft on the diagram and create it.
- **Step by step:** pick nodes (your SSH connections, storages, outside endpoints), add checks from the catalog (process, port, HTTP, file age, disk, log pattern…), connect them and describe **states** (for example *running on backup*) with conditions.
- **Templates** and an advanced JSON editor are there too.

## Use it
The diagram shows each node and link in colour with the time of the last check. Click a node for its checks and history; *Test now* runs all checks once. Notifications can tell you when the chain changes state.

## AI on a chain
From the chain's window, ask the AI to explain a state or to propose a change; it returns a draft you apply yourself.
