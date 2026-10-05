# Working on WRM PRO (web-browser-RDM-public)

Guidance for Claude Code sessions in this repository. The project owner writes in **Croatian**: reply to them in Croatian. Code, comments, commit messages and repository documentation are in **English**.

## What to build
- **[ROADMAP.md](ROADMAP.md)** holds the agreed plan. A session works on **one version** from it, the first one not yet done, unless told otherwise.
- When that version is finished, **remove its section from ROADMAP.md in the same PR**. Its content then lives in CHANGELOG.md.
- Do not start the next version on your own.

## Project layout (short)
- `remote-manager/`: one Go program (single binary) with the web UI embedded from `remote-manager/static/`.
  - The whole UI is `static/index.html`: vanilla JS, no build step.
  - Translations live in `LANGS.en` / `LANGS.hr`, with later additions as `Object.assign(LANGS.en, {...})` / `Object.assign(LANGS.hr, {...})` blocks placed before `function t(k)`. Every new UI text needs **both** languages.
- **Database:** SQLite with **additive migrations only** (`CREATE TABLE IF NOT EXISTS`, `ALTER TABLE … ADD COLUMN`, `ensureTable`) in `initDB` (`main.go`). A previous binary must still start on the new database.
- **Settings and policies:** `settingSpecs` (`settings.go`). They can be forced by `WRM_<KEY>` environment variables. Per-user policy flags are exposed through `mePayload` (`users.go`).
- **Audit:** `auditLog` / `auditLogRef`. Detail keys containing `password`, `token`, `credential`, `private_key` … are redacted, so use neutral keys such as `cred_id`.
- **Secrets:** encrypted at rest (`encryptValue`), never sent to the browser.
- Read [ARCHITECTURE.md](ARCHITECTURE.md) before larger changes.

## Before every push
Run everything from `remote-manager/` unless noted:

```
gofmt -l .                      # must print nothing
go vet ./...
go test -race -count=1 ./...    # full suite, ~4–5 minutes
node ../tools/check-ui.js       # UI: JS syntax + en/hr translation keys
```

- Cross-build at least `GOOS=windows`, `GOOS=darwin GOARCH=arm64`, `GOOS=freebsd` and `GOARCH=arm` for code that touches OS-specific parts.
- **Test UI changes in a real browser.** Playwright with a preinstalled Chromium is available in the cloud environment. Build the binary, start it with `PORT=… DB_PATH=<scratch dir>/wrm.db`, and drive it.
  - Check: no console errors, Croatian texts, phone width (390 px), and that the top bar does not overlap at 1100–1920 px.
  - Look at the screenshots.

## Version bump and docs (every version)
- **Version in:**
  - `remote-manager/main.go` (`AppVersion`);
  - `remote-manager/static/index.html` (`<title>`, `#auth-ver`, `.app-ver`);
  - `.github/workflows/build.yml` (`APP_VERSION`);
  - `docker-compose.yml`;
  - `remote-manager/build-all.sh`;
  - README (current version, download file names).
- **CHANGELOG.md:** a new section at the top (Keep a Changelog: Added / Changed / Fixed).
- **RELEASE_NOTES.md:** a new top section (same structure as the current one, links to the new tag). The previous top section becomes "Included from vX.Y.Z" (its `###` headings become `####`). Download links point to the new version.
- **README.md:** an "Upgrading from vPrev" section (and its TOC link), feature documentation, configuration table, API reference, project layout, tests list, troubleshooting, as far as the change touches them.
- **ARCHITECTURE.md, SECURITY.md, CONTRIBUTING.md:** where the change touches them.
- **Never** use the suffix `-mimo` anywhere (old releases had it; new ones must not).

## Git, PRs and merging
- Each version is developed on its own branch `claude/vX.Y.Z`, branched from the latest `main` (unless the session prompt says otherwise).
- Commit code first, then docs ("vX.Y.Z: version, changelog, release notes and docs").
- **No model names or identifiers** in commits, PRs, code or docs.
- Open the PR against `main` and follow `.github/pull_request_template.md`.
  - Leave the contributor-terms checkbox **unchecked** and add the note "The project owner decides on the contributor terms checkbox."
- Then watch the PR's CI. Fix failures and push again.
- When every check is green (success or skipped) and the PR is mergeable, **merge it** (merge method `merge`, expected head SHA = the full 40-character SHA). Then stop watching it.
- **Tags and releases cannot be pushed from these sessions** (HTTP 403). After the merge, tell the owner to create the release `vX.Y.Z` at https://github.com/vedranius/web-browser-RDM-public/releases/new from `main`.

## Public repository: what must never be committed
This repository is public.
- No organisation-specific names, internal host names or URLs, real service or project names from anyone's private Git server, tokens, customer data or private files shared in chat.
- Features such as the Git workspace are **generic**. Test fixtures use invented names (`demo-api`, `demo-web`, `example.com`) and fake servers (as in `inventory_test.go`'s fake NetBox and `bmc_test.go`'s fake Redfish).

## Finishing a session
End with a short report to the owner **in Croatian**:
- what the version contains;
- how it was tested;
- the PR link and its merge state;
- the reminder to create the release.
