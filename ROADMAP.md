# Roadmap

The agreed plan. Released work is in [CHANGELOG.md](CHANGELOG.md). A version's section is removed from this file in the PR that finishes it.

## v11.6.0 — Git workspace: activity and deploy history per user and branch

- **Git activity per service**: pushes and commits per branch and per author from the provider API (GitLab project events with `action=pushed`; GitHub / Gitea commits per branch), filterable by author, branch and date, with links to the commits.
- **Deploy history view** per service / environment / server: the history of all destinations merged into one view, filterable by who, branch, version and action (deploy / rollback).
- **"What is where" matrix**: which branch / commit / version each environment and server runs, who deployed it and when, and how many commits it is behind the head of its branch.
- **CSV export** of both views.
