---
id: directory-sign-in
title: Directory sign-in (Entra ID, AD, LDAP, SCIM)
section: Administration
order: 160
summary: Let people sign in with Microsoft Entra ID, Okta, Google or any OpenID Connect provider, or LDAP/Active Directory, and provision with SCIM.
keywords: sso, single sign-on, oidc, openid, entra, azure ad, active directory, ad, ldap, adfs, okta, keycloak, google workspace, scim, groups, domain, provisioning
ui: #modal-sso, #modal-ldap, #modal-scim, #modal-sso-test
open: admin:sso
admin: true
---
**Admin panel → Directory sign-in**.

## OpenID Connect (Entra ID, ADFS, Okta, Keycloak, Google)
1. **Add a provider** and choose its kind; the editor says where to create the client there.
2. Register the **redirect address** shown (`https://<WRM>/api/auth/oidc/callback`).
3. Enter the issuer, client ID and secret; **Test connection** checks every step.
4. **Rules** map a claim (`roles`, `groups`, `hd`…) to a WRM role and teams. Without rules only linked accounts may sign in.

## LDAP / Active Directory
**Add directory:** kind (Active Directory or other LDAP), `ldaps://` server (or `ldap://` with StartTLS), a service account that may read users and groups, the user and group base, and rules from groups to roles. *Test* signs in a user of your choice.

## SCIM
Entra ID or Okta can create, change and disable WRM accounts automatically: **SCIM → New token**, then give the address and token to the directory.

Directory users cannot sign in with a local password unless you allow it; the break-glass administrator always can.
