# Integration test container for Xorg bindings
#
# Runs hydra's Xorg CGo tests inside an isolated Xvfb display,
# so keyboard/mouse simulation and window manipulation never
# leak into the developer's desktop session.
#
# Build & Run:
#   podman build -t hydra-xorg-test -f Containerfile .
#   podman run --rm hydra-xorg-test
#
# Or use the convenience script:
#   bash test.sh

FROM docker.io/golang:bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    xvfb \
    libx11-dev \
    libxi-dev \
    libxtst-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

ENV DISPLAY=:99

CMD Xvfb :99 -screen 0 1920x1080x24 & \
    sleep 1 && \
    go test -tags=integration -v -count=1 ./adapters/xorg/...
