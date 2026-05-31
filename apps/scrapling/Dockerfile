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
