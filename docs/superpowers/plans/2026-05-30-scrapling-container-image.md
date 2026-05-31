# Scrapling Container Image Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `apps/scrapling/` to the `kalenarndt/containers` repo, producing a semver-tagged, rootless, multi-arch Scrapling MCP image at `ghcr.io/kalenarndt/scrapling`.

**Architecture:** Single-stage Ubuntu image. Scrapling is installed via `uv` into a dedicated venv (avoids Ubuntu's PEP 668 externally-managed-environment block), then Playwright Chromium is baked into a world-readable `/browsers` at build time so the rootless runtime needs no initContainer or PVC. Follows the existing `plex` app's layout and conventions.

**Tech Stack:** Docker Buildx Bake (HCL), Ubuntu `noble`, Python 3 + `uv`, Playwright Chromium, Go `testcontainers` smoke test, mise task runner.

**Spec:** `docs/superpowers/specs/2026-05-30-scrapling-container-image-design.md`

**Reference files (read before starting):**
- `apps/plex/Dockerfile`, `apps/plex/docker-bake.hcl`, `apps/plex/container_test.go` — the model app (Ubuntu, root build → nobody runtime, catatonit entrypoint).
- `apps/actions-runner/docker-bake.hcl` — renovate annotation placement on the `VERSION` variable.
- `testhelpers/` — `GetTestImage` and `TestCommandSucceeds` signatures.
- `.mise.toml` `[tasks.local-build]` — the local build+test command.

---

### Task 1: Create `docker-bake.hcl`

**Files:**
- Create: `apps/scrapling/docker-bake.hcl`

- [ ] **Step 1: Write the bake file**

```hcl
target "docker-metadata-action" {}

variable "APP" {
  default = "scrapling"
}

variable "VERSION" {
  // renovate: datasource=github-releases depName=D4Vinci/Scrapling
  default = "0.4.8"
}

variable "SOURCE" {
  default = "https://github.com/D4Vinci/Scrapling"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    VERSION = "${VERSION}"
  }
  labels = {
    "org.opencontainers.image.source" = "${SOURCE}"
  }
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}
```

- [ ] **Step 2: Verify bake config parses**

Run: `cd apps/scrapling && docker buildx bake image-all --print`
Expected: prints resolved JSON with `platforms: ["linux/amd64","linux/arm64"]` and `args.VERSION = "0.4.8"`. No HCL parse error.

---

### Task 2: Create the smoke test

**Files:**
- Create: `apps/scrapling/container_test.go`

- [ ] **Step 1: Write the test**

```go
package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/scrapling:rolling")
	testhelpers.TestCommandSucceeds(t, image, nil, "scrapling", "--version")
}
```

Notes:
- Import path is `github.com/kalenarndt/containers/testhelpers` (the module was renamed from `home-operations`).
- `TestCommandSucceeds(t, image, config, entrypoint, args...)` runs the container with the given entrypoint, waits for exit, and asserts exit code 0. A bare HTTP probe against the MCP endpoint is flaky (it returns non-200 without an MCP session), so we validate the install + entrypoint resolve instead.

- [ ] **Step 2: Verify it compiles (and fails for the right reason — no image yet)**

Run: `go vet ./apps/scrapling/...`
Expected: exit 0, no "too many arguments"/signature errors. (The test cannot run until the image is built in Task 4 — that is expected.)

---

### Task 3: Create the `Dockerfile`

**Files:**
- Create: `apps/scrapling/Dockerfile`

- [ ] **Step 1: Write the Dockerfile**

```dockerfile
# syntax=docker/dockerfile:1

FROM docker.io/library/ubuntu:noble
ARG VERSION
ARG VENDOR

ENV DEBIAN_FRONTEND="noninteractive" \
    PYTHONUNBUFFERED="1" \
    PYTHONDONTWRITEBYTECODE="1" \
    PLAYWRIGHT_BROWSERS_PATH="/browsers" \
    VIRTUAL_ENV="/opt/venv" \
    PATH="/opt/venv/bin:${PATH}"

USER root
WORKDIR /app

RUN \
    apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests \
        ca-certificates \
        catatonit \
        python3 \
        python3-venv \
        tzdata \
    && python3 -m venv "${VIRTUAL_ENV}" \
    && pip install --no-cache-dir uv \
    && uv pip install "scrapling[ai]==${VERSION}" \
    && playwright install --with-deps chromium \
    && chmod -R 755 "${PLAYWRIGHT_BROWSERS_PATH}" "${VIRTUAL_ENV}" \
    && apt-get clean \
    && rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/*

COPY . /

ENV HOME="/tmp" \
    UV_CACHE_DIR="/tmp/uv-cache"

USER nobody:nogroup
WORKDIR /tmp

ENTRYPOINT ["/usr/bin/catatonit", "--", "scrapling"]
CMD ["mcp", "--http", "--host", "0.0.0.0", "--port", "8000"]
```

Design notes (why each non-obvious choice):
- **Dedicated venv at `/opt/venv`** instead of `uv pip install --system`: Ubuntu `noble`'s system Python is PEP 668 externally-managed, which blocks system installs. A venv sidesteps that cleanly and keeps the install isolated. `VIRTUAL_ENV` + `PATH` make `uv`/`scrapling`/`playwright` resolve from it.
- **`pip install uv`**: `python3-venv` provides `pip` inside the venv; we bootstrap `uv` there, then use `uv pip install` for speed.
- **`scrapling[ai]`**: the `[ai]` extra pulls in the fetchers stack + the built-in MCP server (`scrapling mcp`).
- **`playwright install --with-deps chromium`**: runs as root, installs Chromium's apt OS libraries *and* downloads the Chromium binary into `PLAYWRIGHT_BROWSERS_PATH=/browsers`. `chmod -R 755 /browsers` makes it readable by the runtime `nobody` user — this is what removes the downstream initContainer/PVC hack.
- **Chromium only** (no camoufox) to keep the image lean, matching the current deployment.
- **`nobody:nogroup` + `HOME=/tmp`**: rootless runtime with a writable home/uv-cache.
- **`COPY . /`** pulls in the shared `.dockerignore` that the build rsyncs into the app dir.

- [ ] **Step 2: Verify Dockerfile syntax with a no-op print**

Run: `cd apps/scrapling && docker buildx bake image --print`
Expected: resolved JSON, no Dockerfile parse error.

---

### Task 4: Build and test the image locally

**Files:** none (uses the mise task)

- [ ] **Step 1: Build the amd64 image and run the smoke test**

Run (from repo root): `mise run local-build scrapling`
What it does: rsyncs `include/` into `apps/scrapling/`, runs `docker buildx bake ... --load`, then `go test -v ./apps/scrapling/...` with `TEST_IMAGE` set to the freshly built `scrapling:0.4.8`.
Expected: build succeeds; test output shows `--- PASS: Test` (the container runs `scrapling --version`, prints the version, exits 0).

- [ ] **Step 2: If the build fails on `playwright install`**

Symptom: `playwright: command not found` or missing chromium deps.
Fix path: confirm `scrapling[ai]` actually pulls Playwright (`uv pip show playwright` inside a `docker run --rm -it scrapling:0.4.8-build sh`). If the `[ai]` extra does not include Playwright on this version, add `playwright` explicitly: change the install line to `uv pip install "scrapling[ai]==${VERSION}" playwright`. Re-run Step 1.

- [ ] **Step 3: Sanity-check the MCP server actually serves (manual, optional)**

Run: `docker run --rm -p 8000:8000 scrapling:0.4.8` then in another shell `curl -fsS -o /dev/null -w '%{http_code}\n' localhost:8000` (any HTTP response, even 4xx, proves it is listening). Ctrl-C to stop.
Expected: an HTTP status code is returned (not connection-refused).

---

### Task 5: Commit

**Files:** `apps/scrapling/`

- [ ] **Step 1: Stage and commit only the new app**

```bash
git add apps/scrapling/Dockerfile apps/scrapling/docker-bake.hcl apps/scrapling/container_test.go
git commit -m "feat(scrapling): add rootless multi-arch MCP image

Builds scrapling[ai] on Ubuntu with Chromium baked into a
world-readable /browsers, giving semver tags and removing the need
for the downstream browser-install initContainer.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

Note: do not `git add -A` — the repo may have other unrelated working-tree changes.

- [ ] **Step 2: (Optional) trigger a remote multi-arch build to publish**

Run: `mise run remote-build scrapling true`
What it does: `gh workflow run release.yaml -f app=scrapling -f release=true` — builds amd64+arm64 and pushes `ghcr.io/kalenarndt/scrapling` with semver + `rolling` tags.
Expected: workflow dispatched; verify in `gh run list --workflow=release.yaml`.

---

## Follow-up (separate change, in the `expanse-ops` repo — NOT part of this plan)

Once `ghcr.io/kalenarndt/scrapling` is published, simplify
`kubernetes/apps/ai/toolhive/mcp-servers/scrapling-mcp/`:
- Repoint `image:` (and the initContainer image) to the new image at a pinned semver/digest.
- Delete the `install-browsers` initContainer.
- Delete `pvc.yaml` (the `scrapling-browsers` PVC) and its volume + mounts.
- Drop the `PLAYWRIGHT_BROWSERS_PATH` env and browsers volume; keep only the `/tmp` emptyDir.

This is tracked here for context only — implement it as its own task in the other repo.

---

## Self-Review

- **Spec coverage:** layout (Task 1–3), Ubuntu base (Task 3), `scrapling[ai]` + chromium-in-`/browsers` (Task 3), versioning + renovate annotation (Task 1), `TestCommandSucceeds` test (Task 2), local build/test (Task 4), publish (Task 5 Step 2), follow-up deployment cleanup (documented, out of scope). All spec sections map to a task.
- **Placeholder scan:** no TBD/TODO; every file has complete content; the one conditional (Task 4 Step 2) is a concrete fallback with an exact command, not a placeholder.
- **Type/name consistency:** `APP`/`VERSION`/`SOURCE` variable names consistent between `docker-bake.hcl` and the Dockerfile `ARG VERSION`; image ref `ghcr.io/kalenarndt/scrapling:rolling` consistent between test and publish; `PLAYWRIGHT_BROWSERS_PATH=/browsers` consistent across Dockerfile steps; testhelpers signature `TestCommandSucceeds(t, image, nil, entrypoint, args...)` matches `testhelpers/`.
