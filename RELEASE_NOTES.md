## Web Remote Manager PRO v13.14.11 — secure shares and reliable AI connections

- Claude MCP reconnects after a WRM service restart, preserving its saved authorization and server selections.
- AI connection consent lists include Select all, with individual server deselection.
- Share links keep bearer tokens out of request URLs, use per-share HttpOnly cookies, and validate access before entering collaboration mode.
- Security hardening covers saved credentials in connection tests, approval rechecks, MFA, MCP OAuth consent, outbound requests, addon channels, Windows destructive commands, and network audit logs.
- Pending AI approvals raise desktop notifications that link to the matching approval.
- File, tunnel, session, and folder actions report request failures clearly.

### Upgrade

Use WRM's update action or replace the executable/container with the matching v13.14.11 asset. Preserve the existing database and encryption key.

[Full changelog](CHANGELOG.md)
