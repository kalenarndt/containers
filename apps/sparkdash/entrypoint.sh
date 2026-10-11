#!/usr/bin/env bash
set -euo pipefail

# Seed upstream defaults into an empty config volume; never overwrite user state.
for f in /app/defaults/*; do
    [[ -e "/app/config/${f##*/}" ]] || cp "${f}" /app/config/
done

# Secret mounts with fsGroup are group-readable; ssh and sparkDash want 0600.
if [[ -n "${SSH_KEY_SOURCE:-}" && -f "${SSH_KEY_SOURCE}" ]]; then
    install -m 0600 "${SSH_KEY_SOURCE}" "${HOME}/.ssh/id_ed25519"
fi

exec node /app/server/index.js "$@"
