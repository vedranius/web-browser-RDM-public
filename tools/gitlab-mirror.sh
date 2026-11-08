#!/usr/bin/env bash
# Mirrors this repository to a GitLab project and, for a release tag, creates the matching
# GitLab release whose download links point to the signed GitHub release assets.
#
#   tools/gitlab-mirror.sh push            push every branch and v* tag (branches deleted on
#                                          GitHub are deleted on GitLab too)
#   tools/gitlab-mirror.sh release vX.Y.Z  create or update the GitLab release for a tag
#
# Environment:
#   GITLAB_MIRROR_URL    https URL of the GitLab project, e.g. https://gitlab.com/group/wrm-pro.git
#   GITLAB_MIRROR_TOKEN  GitLab project or personal access token (scopes write_repository + api)
#   GITHUB_REPOSITORY    owner/name of the GitHub repository (release mode, for asset links)
#   GITLAB_MIRROR_DEFAULT_BRANCH  branch never force-pushed (default: main)
#
# The token is sent in an HTTP header, never in a URL, so it does not appear in git output.
set -euo pipefail

mode="${1:-}"
url="${GITLAB_MIRROR_URL:-}"
token="${GITLAB_MIRROR_TOKEN:-}"
default_branch="${GITLAB_MIRROR_DEFAULT_BRANCH:-main}"

if [ -z "$url" ] || [ -z "$token" ]; then
  echo "::notice::GitLab mirror is not configured (repository variable GITLAB_MIRROR_URL and secret GITLAB_MIRROR_TOKEN). Nothing to do."
  exit 0
fi

case "$url" in
  https://*) ;;
  file://*) [ "${GITLAB_MIRROR_ALLOW_FILE:-}" = 1 ] || { echo "::error::GITLAB_MIRROR_URL must be an https URL"; exit 1; } ;;
  *) echo "::error::GITLAB_MIRROR_URL must be an https URL"; exit 1 ;;
esac

auth_header="Authorization: Basic $(printf 'oauth2:%s' "$token" | base64 | tr -d '\n')"
if [ -n "${GITHUB_ACTIONS:-}" ]; then
  echo "::add-mask::${auth_header#Authorization: Basic }"
fi

push_all() {
  # Local refs that should exist on GitLab: every origin branch and every v* tag.
  git fetch --quiet --prune --tags origin '+refs/heads/*:refs/remotes/origin/*'
  local remote_heads local_heads refspecs=() b
  local_heads=$(git for-each-ref --format='%(refname:lstrip=3)' refs/remotes/origin | grep -vx 'HEAD' || true)
  for b in $local_heads; do
    if [ "$b" = "$default_branch" ]; then
      refspecs+=("refs/remotes/origin/$b:refs/heads/$b")      # never forced: diverging history fails loudly
    else
      refspecs+=("+refs/remotes/origin/$b:refs/heads/$b")
    fi
  done
  git -c "http.extraHeader=$auth_header" push --porcelain "$url" "${refspecs[@]}" 'refs/tags/v*:refs/tags/v*'

  # Branches that are gone from GitHub are removed from the mirror.
  remote_heads=$(git -c "http.extraHeader=$auth_header" ls-remote --heads "$url" | awk '{sub("refs/heads/","",$2); print $2}')
  local stale=()
  for b in $remote_heads; do
    if ! grep -qxF "$b" <<<"$local_heads"; then
      [ "$b" = "$default_branch" ] && continue
      stale+=(":refs/heads/$b")
    fi
  done
  if [ "${#stale[@]}" -gt 0 ]; then
    git -c "http.extraHeader=$auth_header" push --porcelain "$url" "${stale[@]}"
  fi
  echo "Mirrored ${#refspecs[@]} branch(es) and the v* tags; removed ${#stale[@]} stale branch(es)."
}

release() {
  local tag="$1"
  if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "::error::Not a release tag: $tag"; exit 1
  fi
  local rest host path api project gh_json payload code
  rest="${url#https://}"
  host="${rest%%/*}"
  path="${rest#*/}"
  path="${path%.git}"
  path="${path%/}"
  project=$(jq -rn --arg p "$path" '$p|@uri')
  api="https://${host}/api/v4/projects/${project}"

  gh_json=$(gh release view "$tag" --repo "$GITHUB_REPOSITORY" --json name,body,url,assets)
  payload=$(jq -c --arg tag "$tag" '{
      tag_name: $tag,
      name: (.name // $tag),
      description: ((.body // "") + "\n\n---\nSigned binaries, SHA256SUMS.txt and its signature: " + .url),
      assets: {links: [.assets[] | {name: .name, url: .url, link_type: (if (.name|test("\\.(txt|sig|json|xml)$")) then "other" else "package" end)}]}
    }' <<<"$gh_json")

  code=$(curl -sS -o /dev/null -w '%{http_code}' -H "PRIVATE-TOKEN: $token" "$api/releases/$tag")
  if [ "$code" = 200 ]; then
    # Links cannot be replaced in one call; update the text, links added on creation stay.
    curl -fsS -o /dev/null -X PUT -H "PRIVATE-TOKEN: $token" -H 'Content-Type: application/json' \
      --data "$(jq -c '{name, description}' <<<"$payload")" "$api/releases/$tag"
    echo "Updated GitLab release $tag."
  else
    curl -fsS -o /dev/null -X POST -H "PRIVATE-TOKEN: $token" -H 'Content-Type: application/json' \
      --data "$payload" "$api/releases"
    echo "Created GitLab release $tag with $(jq '.assets.links|length' <<<"$payload") download link(s)."
  fi
}

case "$mode" in
  push) push_all ;;
  release) push_all; release "${2:-}" ;;
  *) echo "usage: $0 push | release vX.Y.Z" >&2; exit 2 ;;
esac
