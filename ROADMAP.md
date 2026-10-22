# Roadmap

Accepted work through **v12.9.0** is delivered in [CHANGELOG.md](CHANGELOG.md).

## Security and enterprise identity (accepted, versions before v12.7.0)

Goal: the same WRM fits a home lab, a small company and a large organisation. Security is checked automatically on every change, the operator chooses a security level, and companies sign in with their directory (Microsoft Entra ID / AD, LDAP, other identity providers) so every action is traceable to a domain user and can be stopped at once.

v12.6.20 delivered the continuous security pipeline, v12.6.21 signed releases, SBOM, build provenance and the nightly fuzz and ZAP tests, v12.6.22 the security levels (Home / Business / Enterprise), deviations and the security posture page, v12.6.23 roles, teams, step-up re-authentication, the break-glass administrator and handover on delete, v12.6.24 directory sign-in (OIDC), v12.6.25 LDAP / Active Directory and SCIM, v12.6.26 SIEM export, security alerts, *Stop user* and lockdown (see [CHANGELOG.md](CHANGELOG.md)).

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

Delivery plan (v12.6.0 delivered the deploy-config reader/validator, the native plan and the comparison with the original tool; v12.7.0 native deploy, dry run, lock, backup, history and rollback through WRM connections with a per-project engine switch; see [CHANGELOG.md](CHANGELOG.md)):
1. **v12.7.3 – build locations (delivered):** projects with a build step deploy natively; per project the build runs on the WRM host, on a chosen Linux build server over SSH (temporary directory, only the prepared source sent, no Git or SSH keys) or not at all, chosen and shown in the plan before anything runs, with the build server's platform and architecture checked against the destinations.
2. **v12.9.0 – screens (delivered):** file details (server against branch), line-ending marking (`eol_only`) and the deploy-config editor with validation before saving.
3. **v12.9.1 – switch-over:** config Git actions (state, pull, commit and push of a deploy-config clone, validated before commit), `build_user`, native engine as the default on Linux and Windows; the Docker runtime becomes optional and is then retired.
