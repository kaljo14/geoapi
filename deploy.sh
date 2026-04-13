#!/usr/bin/env bash
set -euo pipefail

IMAGE="kaljo14/places-scraper"
PLATFORMS="linux/amd64,linux/arm64"
BUILDER_NAME="geopulse-multiarch"

# --- Resolve version tag ---------------------------------------------------

if [[ -n "${1:-}" ]]; then
    TAG="$1"
else
    LATEST=$(git tag --sort=-v:refname --list 'v*' | head -1)
    if [[ -z "$LATEST" ]]; then
        TAG="v0.1.0"
    else
        # Bump patch: v0.1.2 → v0.1.3
        IFS='.' read -r MAJOR MINOR PATCH <<< "${LATEST#v}"
        TAG="v${MAJOR}.${MINOR}.$((PATCH + 1))"
    fi
    echo "Auto-resolved tag: ${TAG}"
fi

echo "==> Building ${IMAGE}:${TAG} for ${PLATFORMS}"

# --- Ensure buildx builder exists -----------------------------------------

if ! docker buildx inspect "${BUILDER_NAME}" &>/dev/null; then
    echo "==> Creating buildx builder: ${BUILDER_NAME}"
    docker buildx create --name "${BUILDER_NAME}" --use
else
    docker buildx use "${BUILDER_NAME}"
fi

# --- Build & push ----------------------------------------------------------

docker buildx build \
    --platform "${PLATFORMS}" \
    --tag "${IMAGE}:${TAG}" \
    --tag "${IMAGE}:latest" \
    --push \
    .

echo "==> Pushed ${IMAGE}:${TAG} and ${IMAGE}:latest"

# --- Git tag ---------------------------------------------------------------

git tag "${TAG}"
git push origin "${TAG}"

echo "==> Tagged git with ${TAG}"
echo "Done."
