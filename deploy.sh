#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
release_tag="${1:-}"
if [[ $# -ne 1 || ! "$release_tag" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
    echo "Usage: $0 vMAJOR.MINOR.PATCH (for example v1.2.3)" >&2
    exit 1
fi
if [[ -n "$(git status --porcelain)" ]]; then
    echo 'Commit or stash local changes before releasing.' >&2
    exit 1
fi
if git show-ref --verify --quiet "refs/tags/$release_tag"; then
    echo "Tag $release_tag already exists; choose a new version." >&2
    exit 1
fi
# Check the remote before creating a local tag; fail on network/auth errors.
remote_tag=$(git ls-remote --tags origin "refs/tags/$release_tag")
if [[ -n "$remote_tag" ]]; then
    echo "Tag $release_tag already exists on origin; choose a new version." >&2
    exit 1
fi
git tag -a "$release_tag" -m "Release $release_tag"
git push origin "refs/tags/$release_tag"
echo "Pushed $release_tag. GitHub Actions will check, build, and publish the image."
