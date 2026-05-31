# Scrapling container image — design

**Date:** 2026-05-30
**Repo:** `home-operations/containers` fork (`kalenarndt`)
**Status:** Approved design, pending implementation plan

## Problem

The `scrapling-mcp` deployment in `expanse-ops` consumes the upstream
`ghcr.io/d4vinci/scrapling` image. Two recent commits document the pain:

- `7dd506d` — upstream only publishes `:latest`/`:main` (no semver tags), forcing
  a manual `@sha256` digest pin with no clean upgrade path.
- `7fa360d` + the `install-browsers` initContainer — the upstream image bakes
  Playwright browsers into root-owned `/root`, unreadable by uid 1000 under
  toolhive's rootless/read-only-rootfs sandbox. The deployment works around this
  with an initContainer that re-downloads chromium into a PVC at
  `PLAYWRIGHT_BROWSERS_PATH`.

The underlying project (`D4Vinci/Scrapling`) is a clean PyPI package with proper
semver GitHub releases (currently **v0.4.8**). Building our own image gives us
semver + rootless + multi-arch and lets us bake chromium into a world-readable
path, eliminating both workarounds.

## Goal

Add `apps/scrapling/` to the containers repo producing a semver-tagged, rootless,
multi-arch (amd64 + arm64) Scrapling MCP image published to
`ghcr.io/kalenarndt/scrapling`.

Non-goal for this spec (tracked as a separate follow-up change): simplifying the
`expanse-ops` deployment once the image is published.

## Design

### Layout

Mirrors the existing `plex`/`octoprint` apps:

```
apps/scrapling/
  Dockerfile
  docker-bake.hcl
  container_test.go
```

The shared `include/.dockerignore` is rsynced into the app dir by the
`app-builder` workflow and pulled in via `COPY . /`.

### Base image

`docker.io/library/ubuntu:noble`. Scrapling drives Playwright/Chromium, which
officially requires glibc; upstream uses `python:3.12-slim-trixie`. Alpine/musl
chromium is unsupported and fragile on arm64. The repo README explicitly permits
Alpine **or** Ubuntu, and the existing `plex` app is Ubuntu-based.

### Dockerfile

- `FROM docker.io/library/ubuntu:noble`, `ARG VERSION`, `ARG VENDOR`,
  `ENV DEBIAN_FRONTEND=noninteractive`.
- Build as `root`:
  - apt-install `python3`, `python3-venv`, `ca-certificates`, `catatonit`,
    `tzdata`, `curl`.
  - `pip install uv` then `uv pip install --system "scrapling[ai]==${VERSION}"`
    (the `[ai]` extra pulls in the fetchers + the built-in MCP server).
  - `ENV PLAYWRIGHT_BROWSERS_PATH=/browsers` (a world-readable path), then
    `playwright install --with-deps chromium` to bake chromium **and its OS
    libraries** into the image. `chmod -R 755 /browsers`.
  - **Chromium only** — skip camoufox to keep the image lean and match the
    current deployment's intent.
  - Strip apt lists/caches; `COPY . /`.
- Runtime as `nobody:nogroup`:
  - `ENV HOME=/tmp UV_CACHE_DIR=/tmp/uv-cache` so the rootless runtime has a
    writable home/cache.
  - `ENTRYPOINT ["/usr/bin/catatonit", "--", "scrapling"]`
  - `CMD ["mcp", "--http", "--host", "0.0.0.0", "--port", "8000"]` — toolhive
    overrides args at deploy time, but a sane default keeps the image usable
    standalone.

Baking chromium into a world-readable `/browsers` is the core fix: the rootless
runtime reads browsers directly from the image, so no initContainer or PVC is
needed downstream.

### Versioning (`docker-bake.hcl`)

- `APP = "scrapling"`, `SOURCE = "https://github.com/D4Vinci/Scrapling"`,
  `VERSION` default set to the current release.
- Renovate annotation so `VERSION` auto-bumps from upstream semver releases:
  `# renovate: datasource=github-releases depName=D4Vinci/Scrapling`
  (placed per the repo's customManager convention — confirm exact placement
  against an existing app during implementation).
- `image`/`image-all` targets set `platforms = ["linux/amd64", "linux/arm64"]`,
  matching `plex`.

### Test (`container_test.go`)

```go
image := testhelpers.GetTestImage("ghcr.io/kalenarndt/scrapling:rolling")
testhelpers.TestCommandSucceeds(t, image, nil, "scrapling", "--version")
```

A bare HTTP GET against the streamable-HTTP MCP endpoint returns a non-200
(needs an MCP session/headers), so an `TestHTTPEndpoint` wait is flaky.
`TestCommandSucceeds` validates the install + entrypoint resolve. An HTTP smoke
test can be added later if desired.

## Follow-up (separate change, after image publishes)

In `kubernetes/apps/ai/toolhive/mcp-servers/scrapling-mcp/`:

- Repoint `image:` (and the initContainer image) to
  `ghcr.io/kalenarndt/scrapling` at a pinned semver/digest.
- Delete the `install-browsers` initContainer.
- Delete `pvc.yaml` (the `scrapling-browsers` PVC) and its volume + mounts.
- Drop the `PLAYWRIGHT_BROWSERS_PATH` env and browsers volume; keep only the
  `/tmp` emptyDir for uv/runtime writes.

Net: ~40 lines of workaround removed; upgrades become a renovate-driven semver
bump instead of a manual digest pin.

## Risks / open items

- **Image size.** Chromium + glibc base is large (hundreds of MB). Acceptable —
  this is a browser-driven scraper; matches upstream.
- **`playwright install --with-deps` apt behaviour** must run as root at build
  time (it does here — the failure noted in `7fa360d` was at *runtime* under the
  rootless sandbox, which we no longer hit).
- **Renovate annotation placement** for the `VERSION` var to confirm against an
  existing app (`docker-bake.hcl` vs Dockerfile comment) during implementation.
