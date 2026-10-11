<div align="center">

## Containers

_Container images built for my homelab_

</div>

<div align="center">

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kalenarndt/containers/release.yaml?style=for-the-badge&label=Release)

</div>

This repository builds rootless, semantically versioned images for apps that publish no usable image of their own. Images are pushed to [GitHub Packages](https://github.com/kalenarndt?tab=packages&repo_name=containers).

## Images

| Image                          | Upstream                                                      | Platforms     |
| ------------------------------ | ------------------------------------------------------------- | ------------- |
| `ghcr.io/kalenarndt/sparkdash` | [MiaAI-Lab/sparkDash](https://github.com/MiaAI-Lab/sparkDash) | `linux/amd64` |

## Tags

Each release is tagged with the full upstream version (`2.32.0`), `major.minor`, `major` and `rolling`. Tags move; pin to the digest for immutability:

```
ghcr.io/kalenarndt/sparkdash:2.32.0@sha256:<digest>
```

[Renovate](https://github.com/renovatebot/renovate) keeps both the upstream versions here and digest pins downstream up to date.

## Verify an image

Release builds carry a [build-provenance attestation](https://github.com/actions/attest). Check that an image came from this repo's CI:

```sh
gh attestation verify --repo kalenarndt/containers oci://ghcr.io/kalenarndt/sparkdash:2.32.0
```

## sparkdash

A dashboard for NVIDIA DGX Spark machines. It polls each unit over SSH, so it runs anywhere that can reach them; it does not need to run on a Spark.

Differences from the upstream `docker compose` setup:

- Runs as `node` (`1000:1000`) and works with a read-only root filesystem.
- Binds `0.0.0.0` by default. Set `SPARKDASH_TOKEN`: without it, anyone who can reach the port can change settings and power units off.
- No host mounts or `privileged`. The "This host" unit is therefore unavailable; add every Spark as a remote unit.
- On an empty `/app/config` volume, the entrypoint seeds upstream's default files. It never overwrites existing files.
- Set `SSH_KEY_SOURCE` to a mounted private key. The entrypoint copies it to `~/.ssh/id_ed25519` with mode `0600`, so a group-readable Kubernetes secret mount works.

Writable paths with `readOnlyRootFilesystem: true`:

| Path              | Purpose                                                                      |
| ----------------- | ---------------------------------------------------------------------------- |
| `/app/config`     | Unit registry, encrypted secrets and `.secrets-key`. Persist and back it up. |
| `/home/node/.ssh` | SSH key and `known_hosts`                                                    |
| `/tmp`            | SSH control sockets                                                          |

## Credits

Forked from [home-operations/containers](https://github.com/home-operations/containers), which provides the build and release tooling.
