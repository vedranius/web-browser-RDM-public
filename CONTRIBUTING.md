# Contributing to Web Remote Manager PRO

Thanks for helping to make WRM better! Bug reports, ideas, documentation, translations and code
are all welcome.

## Ways to contribute

- **Report a bug** — open an issue using the *Bug report* template. Include your WRM version
  (`/api/version` or the top bar), OS, browser and steps to reproduce.
- **Report a security problem** — privately, as described in [SECURITY.md](SECURITY.md). Please
  do not open a public issue for it.
- **Suggest an idea / feature** — open an issue using the *Feature request* template and describe
  the problem it solves.
- **Improve translations** — the UI strings live in the `LANGS` object in
  `remote-manager/static/index.html` (English and Croatian today; new languages are welcome).
- **Send a pull request** — for anything bigger than a small fix, please open an issue first so we
  can agree on the approach.

## Development

Requirements: Go (version from `remote-manager/go.mod`), Node.js is only needed for the JS syntax
check.

```bash
cd remote-manager
gofmt -l .                      # must print nothing
go vet ./...
go test -race ./...             # unit + integration tests (security, roles, 2FA, encryption, recording, audit…)
go build -o wrm-server .
HTTPS_SELF_SIGNED=1 PORT=8080 ./wrm-server   # open https://localhost:8080
```

Use a separate `DB_PATH` for development so you do not touch your real database. Voice calls need
HTTPS (or `localhost`) and two browser profiles or devices to test.

Code layout:

| File | What it does |
|---|---|
| `main.go` | Startup, routes, database schema & migrations, connections API, server-to-server transfer |
| `security.go` | Encryption of secrets, signed values, security headers, CSRF, HTTPS |
| `users.go`, `totp.go` | Accounts, sessions, sign-in, 2FA, admin user management |
| `settings.go`, `audit.go`, `admin.go` | Policies (feature flags), audit log (hash chain, append-only), admin panel API, terminal registry |
| `recording.go` | Terminal sessions, asciicast session recorder, file transfer log, recordings & transfers API |
| `shares.go`, `collab.go` | Shares, roles & permissions, collaboration rooms (chat, terminal sharing, voice signalling) |
| `turn.go` | Built-in TURN relay for voice calls |
| `ssh_ws.go`, `files.go`, `sftp_pool.go`, `search.go`, `conntest.go`, `hostkeys.go` | Terminal, file manager, SFTP connection pool, search, connection test, host key verification |
| `static/index.html` | The whole web UI (embedded into the binary); `static/vendor/` holds xterm.js and fonts, `static/brand/` the logo and icons |
| `*_test.go` | `security_test.go` and `upgrade_test.go` (unit), `integration_test.go` (in-process SSH/SFTP server: connect → audit → recording, transfers) |

See [ARCHITECTURE.md](ARCHITECTURE.md) for how the parts fit together and where planned features belong.

Before sending a PR please make sure that:

- `gofmt -l .` prints nothing, and `go vet ./...` and `go test -race ./...` pass,
- the page has no JavaScript errors (the CI extracts the inline script and runs `node --check`),
- new UI text has both `en` and `hr` entries in `LANGS`,
- every new API endpoint or WebSocket message checks **who** may use it (signed-in user, share role
  / permission, admin) on the server. The UI hiding a button is not access control,
- secrets (passwords, keys, tokens) are never returned to the browser or written to the log,
- security-relevant actions write an audit entry (`auditLog…`), and nothing updates or deletes audit data outside the retention job,
- new features that change behaviour sit behind a setting (feature flag), and the existing flow works with the flag off,
- you tested the change in a browser (terminal, file manager, sessions, collaboration and voice if
  you touched them).

CI (GitHub Actions) builds every platform for each pull request.

## Contributor terms

By submitting a contribution (code, documentation, translations or other material) to this
repository you agree that:

1. You wrote it yourself, or otherwise have the right to submit it under these terms.
2. Your contribution is licensed under the project's licenses
   ([PolyForm Noncommercial 1.0.0 and PolyForm Internal Use 1.0.0](LICENSE)), so everyone can use
   it on the same terms as the rest of the project.
3. You additionally grant the project owner (**vedranius**) a perpetual, worldwide, non-exclusive,
   royalty-free, irrevocable license to use, modify, distribute, sublicense and relicense your
   contribution as part of the project, **including under commercial licenses**. This keeps the
   project able to offer commercial licenses (see [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)),
   which helps fund further development.
4. You keep the copyright of your contribution and may use your own contribution however you like.

If you do not agree with these terms, please open an issue to discuss instead of sending code.

## Code of conduct

Be respectful and constructive. Harassment or abusive behaviour is not tolerated.
