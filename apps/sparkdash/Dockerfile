# syntax=docker/dockerfile:1

FROM docker.io/library/node:22-bookworm-slim AS builder
ARG REVISION

ADD --keep-git-dir=false https://github.com/MiaAI-Lab/sparkDash.git#${REVISION} /app

WORKDIR /app

ENV VITE_HISTORY_HOURS=8

RUN npm ci --no-audit --no-fund && \
    npm run build && \
    npm prune --omit=dev --no-audit --no-fund

FROM docker.io/library/node:22-bookworm-slim

ENV \
    BIND_HOST=0.0.0.0 \
    HOME=/home/node \
    LLM_PORT=8888 \
    NODE_ENV=production \
    PORT=5555

USER root
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        catatonit \
        iproute2 \
        openssh-client \
        procps \
        sshpass \
        tzdata && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/dist ./dist-baked
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json /app/package-lock.json ./
COPY --from=builder /app/server ./server
COPY --from=builder /app/src/shared ./src/shared
COPY --from=builder /app/src/components/ShowcasePage/showcasePrompts.ts ./src/components/ShowcasePage/showcasePrompts.ts
COPY --from=builder /app/config ./defaults
COPY . /

RUN mkdir -p /app/config /home/node/.ssh && \
    chmod 0700 /home/node/.ssh && \
    chown -R node:node /app/config /home/node

USER node:node

VOLUME ["/app/config"]

EXPOSE 5555

ENTRYPOINT ["/usr/bin/catatonit", "--", "/entrypoint.sh"]
