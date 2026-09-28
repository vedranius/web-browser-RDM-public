# Contributing to Web Remote Manager PRO

Thanks for helping to make WRM better! Bug reports, ideas, documentation, translations and code
are all welcome.

## Ways to contribute

- **Report a bug** — open an issue using the *Bug report* template. Include your WRM version
  (`/api/version` or the top bar), OS, browser and steps to reproduce.
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
go vet ./...
go build -o wrm-server .
PORT=8080 ./wrm-server          # open http://localhost:8080
```

The whole frontend is `remote-manager/static/index.html` (embedded into the binary). Before sending
a PR please make sure that:

- `go vet ./...` passes and new Go files are `gofmt`-ed,
- the page has no JavaScript errors (the CI extracts the inline script and runs `node --check`),
- new UI text has both `en` and `hr` entries in `LANGS`,
- you tested the change in a browser (terminal, file manager and sessions if you touched them).

CI (GitHub Actions) builds every platform for each pull request.

## Contributor terms

By submitting a contribution (code, documentation, translations or other material) to this
repository you agree that:

1. You wrote it yourself, or otherwise have the right to submit it under these terms.
2. Your contribution is licensed under the project's license
   ([PolyForm Noncommercial 1.0.0](LICENSE)), so everyone can use it on the same terms as the rest
   of the project.
3. You additionally grant the project owner (**vedranius**) a perpetual, worldwide, non-exclusive,
   royalty-free, irrevocable license to use, modify, distribute, sublicense and relicense your
   contribution as part of the project, **including under commercial licenses**. This keeps the
   project able to offer commercial licenses (see [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md)),
   which helps fund further development.
4. You keep the copyright of your contribution and may use your own contribution however you like.

If you do not agree with these terms, please open an issue to discuss instead of sending code.

## Code of conduct

Be respectful and constructive. Harassment or abusive behaviour is not tolerated.
