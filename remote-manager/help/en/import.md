---
id: import
title: Import and export
section: Basics
order: 25
summary: Bring connections in from mRemoteNG, PuTTY, OpenSSH, CSV/Excel, NetBox or another WRM.
keywords: import, export, mremoteng, confcons, putty, ssh config, csv, excel, xlsx, netbox, backup, migrate, move
ui: #st-data, #modal-inventory
open: settings:st-data
---
## From other programs (Settings → Data)
- **mRemoteNG:** choose `confCons.xml` (`%APPDATA%\mRemoteNG\confCons.xml`, or *File → Export* in mRemoteNG). Containers become folders, SSH/RDP/VNC/Telnet/HTTP connections keep their options, and passwords are decrypted and stored encrypted in WRM. If the file has a master password, WRM asks for it. Check the preview, then import.
- **PuTTY:** export the registry key `HKCU\Software\SimonTatham\PuTTY\Sessions` to a `.reg` file and import it; proxies come along.
- **OpenSSH:** import `~/.ssh/config`; `Host` entries become connections, `ProxyJump` becomes a jump host.

If something is not imported as you expect, the preview names each skipped entry and why. Ask the help AI with the message; it can explain the cause.

## Inventory (CSV, Excel, NetBox)
**Settings → Data → Inventory** (or the sidebar ⋯ menu):
- **CSV / Excel (.xlsx):** the first row holds column names; WRM guesses which column is the name, host, port, user, folder and tags. Check the mapping in the preview.
- **NetBox:** the API address and a read-only token. Sites, racks, roles and tenants become tags; importing again updates what was imported before.

## Export and backup
- **Export connections** writes your connections, folders, jump hosts and tunnels as a file without secrets. **Export with passwords & keys** includes them; it asks for your password and the policy must allow it.
- **Import** of a WRM export restores it on another WRM or after a reinstall.
- Administrators back up the whole WRM: the database file and its encryption key (see *Updates and backups*).
