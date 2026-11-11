---
id: users-roles-teams
title: Users, roles and teams
section: Administration
order: 150
summary: Create accounts, give roles, build teams, require re-authentication for dangerous actions and keep a break-glass admin.
keywords: user, users, account, create user, role, administrator, admin, operator, auditor, team, teams, permission, policy, step-up, break-glass, disable, delete, reset password, reset 2fa, registration
ui: #modal-admin, #modal-user-new, #modal-team, #modal-user-delete
open: admin:users
admin: true
---
**Settings → 🛡 Admin panel → Users** (administrators).

## Accounts
- **New user:** name, role, optional e-mail (an activation link when mail is set up) or a temporary password that must be changed at first sign-in.
- **⋯** on a user: change role, reset password or 2FA, sign out everywhere, unlock, disable, delete (with everything the user owns, or handed over to another account).
- Registration is closed after the first account unless you open it (*Policies*).

## Roles
- **Administrator:** everything.
- **Operator:** features limited to operators by policies (for example deploys).
- **User:** own connections and what is shared with them.
- **Security auditor:** reads the audit log, sessions and file edits; uses no features.

## Teams
**Teams** groups people; policies can name `teams:<name>` and shares can be given to a team.

## Policies
**Admin → Policies** decides who may use each feature (off, administrators, operators, teams or everybody) and sets limits. A policy forced by a `WRM_<KEY>` environment variable is shown as locked.

## Step-up and break-glass
Dangerous actions (policies, roles, secrets) ask for the password and 2FA again (*step-up*). A **break-glass** administrator can still sign in during a lockdown or when directory sign-in fails; keep one with a strong password and 2FA.
