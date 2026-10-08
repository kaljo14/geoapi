# Semantic image releases

Push a Git tag `vMAJOR.MINOR.PATCH` to publish `kaljo14/places-scraper:MAJOR.MINOR.PATCH`
for Linux AMD64 and ARM64. For example, `v1.2.3` publishes
`kaljo14/places-scraper:1.2.3`. Each repository has its own version sequence.
Use a new patch version for fixes, minor version for compatible features, and
major version for breaking changes. Versions are explicit; CI does not infer bumps
from commit messages or `package.json`.

The workflows check/build pushes to `main`. Only release tags publish images.
No new `latest`, SHA, major-only, or minor-only aliases are pushed. Stable tags
must have three numeric components without leading zeroes; prereleases and build
metadata are not supported by this stable production release path.

## GitHub secrets

In this repository's Settings → Secrets and variables → Actions, add:

- `DOCKERHUB_USERNAME`: the Docker Hub account with push access to `kaljo14/places-scraper`.
- `DOCKERHUB_TOKEN`: a Docker Hub personal access token with Read & Write access.
- Optional `MAP_INFRA_DISPATCH_TOKEN`: GitHub token with Contents write access to `kaljo14/map-infra`, to trigger Renovate after publishing.

Create the Docker Hub repository before the first release.

## Publish

Commit and push the workflow changes to `main` first. From a clean checkout of
the commit you intend to release, choose an unused version and run:

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin refs/tags/v1.2.3
```

The version above is an example. Inspect existing Git/Docker Hub tags before
choosing the next version. Never move a released Git tag or rebuild it from different
source; publish a new version for fixes. Configure Docker Hub tag immutability for
release tags if you need the registry to enforce this policy.

Alternatively, `make deploy TAG=v1.2.3` checks the version, working tree, and existing tags,
then pushes the release tag. It no longer builds or pushes images locally.

Image publishing waits for every existing CI test, lint, secret, dependency,
and infrastructure scan job to pass on the tagged commit.

## First semantic deployment

After the image is published, run this from `map-infra` (using the real version):

```bash
python3 scripts/adopt-release.py places-scraper 1.2.3
```

This verifies the registry image includes AMD64 and ARM64, then changes the manifest
to the version plus its image-index digest. Review and merge the change to `main`.
Flux deploys it; Renovate proposes subsequent stable semantic version updates.
Existing `latest` deployments remain on the legacy path until this one-time adoption.
A notification alone cannot convert `latest` into a versioned reference.
