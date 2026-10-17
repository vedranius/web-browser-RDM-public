# Roadmap

Accepted work through **v12.6.22** is delivered in [CHANGELOG.md](CHANGELOG.md).

## Security and enterprise identity (accepted, versions before v12.7.0)

Goal: the same WRM fits a home lab, a small company and a large organisation. Security is checked automatically on every change, the operator chooses a security level, and companies sign in with their directory (Microsoft Entra ID / AD, LDAP, other identity providers) so every action is traceable to a domain user and can be stopped at once.

v12.6.20 delivered the continuous security pipeline, v12.6.21 signed releases, SBOM, build provenance and the nightly fuzz and ZAP tests, v12.6.22 the security levels (Home / Business / Enterprise), deviations and the security posture page (see [CHANGELOG.md](CHANGELOG.md)). Next, in this order:

1. **v12.6.23 – roles and teams:** built-in roles Administrator, Security auditor (audit, recordings and posture read-only), Operator, User; teams own shared connections, folders and credentials; policies accept `teams:<list>` besides `off/admins/all`; step-up re-authentication for dangerous actions (also a rule of the Business and Enterprise levels); a local break-glass administrator; on deprovisioning, a user's connections and credentials pass to their team (or a named admin), audited.
2. **v12.6.24 – OIDC sign-in:** generic OIDC (authorization code + PKCE) for Entra ID, ADFS, Okta, Keycloak, Google Workspace and others; users keyed by the immutable directory ID (`iss`+`sub`, Entra `tid`+`oid`), never by name; Entra **app roles** (and group IDs) mapped to WRM roles and teams with a preview; just-in-time accounts only for mapped groups, others refused with a clear message; MFA proven by the token (`amr`/`acr`) when required; periodic revalidation so a user disabled in the directory loses access; local accounts can be limited to break-glass (a security level can require it); a connection test in the UI that explains every step. Tests use a fake identity provider.
3. **v12.6.25 – LDAP / Active Directory and SCIM:** LDAPS or StartTLS only, service-account search then user bind, nested AD groups, periodic sync of disabled/locked state and group membership; a SCIM 2.0 endpoint so Entra ID / Okta create, change and disable users and groups immediately. Disabling by any route ends terminals, tunnels, rooms, sessions, AI sessions and MCP tokens.
4. **v12.6.26 – SIEM export and containment:** audit to syslog (RFC 5424 over TLS), JSON lines and HTTPS (Splunk HEC, Elastic, Microsoft Sentinel, generic), with identity fields and the hash chain; delivery status visible to the admin; alert rules (failed sign-ins, break-glass sign-in, new admins, secret reveals, AI full-access approvals, audit off); *Stop user* (one click, shows what was ended) and *Lockdown* mode (only break-glass admins, AI kill switch, shares paused) from the UI, `wrm -lockdown` and an environment variable.

Later, when needed: passkeys/WebAuthn, SAML 2.0, Kerberos/SPNEGO, trusted proxy headers, client certificates, encryption key in an external key store with key rotation, a CSP without `'unsafe-inline'`, encrypted backups.

## Medium without Docker (accepted direction, versions to follow)

Goal: Medium keeps the original deploy tool's interface, options and `deploy-config` format, but runs on WRM's own engine. No Docker runtime, no Python virtual environment, no dependency installation, and the same behaviour whether WRM runs on Linux or Windows.

Agreed design:
- **Interface:** the original tool's screens (overview, comparison, details/diff, plan and exclusions, dry run, confirmations, progress, lock, history and rollback, server/environment/project settings, YAML editor, config Git actions) stay as they are. WRM serves them and answers their requests from Go.
- **Configuration:** `deploy-config` (servers, environments, `projects/*.yaml`) stays the source of truth and is read without conversion. The editor keeps comments and layout; validation runs before every deploy.
- **Connections and keys:** destinations use the connections already saved in WRM (keys, vault/Bitwarden keys, host-key checks, jump hosts) instead of a separate SSH profile ZIP. Keys never leave WRM's store.
- **Transfer and state:** WRM's native SSH transfer, lock, backup, integrity check and rollback, writing the server-side state (`current.json`, `history.jsonl`, backups) in the format the original tool reads, so existing history stays valid.
- **Where builds run** (per project, shown in the plan before anything runs):
  - on the WRM host (for example WRM on a Linux server with the needed toolchain);
  - on a chosen Linux build server over SSH, in a temporary directory, with only the prepared source sent (no Git or SSH keys); the result is then transferred to the destinations;
  - no build (deploy the selected files as today).
  Platform and architecture of the build server are checked against the destinations.
- **Trust:** build and `post_deploy` commands still run code from the project/config; Medium stays limited to trusted operators.

Delivery plan (v12.6.0 delivered the deploy-config reader/validator, the native plan and the comparison with the original tool; see [CHANGELOG.md](CHANGELOG.md)):
1. **v12.7.0 – native engine:** deploy, dry run, lock, backup, history and rollback through WRM connections; build location choice; per-project switch between the Docker runtime and the native engine. The Docker path keeps working until the native engine matches it.
2. **v12.8.0 – switch-over:** native engine as the default on Linux and Windows; the Docker runtime becomes optional and is then retired.
