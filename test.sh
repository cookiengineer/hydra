#!/usr/bin/env bash

# Arch Linux:
#   sudo pacman -S podman
#
# Debian / Ubuntu:
#   sudo apt install podman

set -euo pipefail;

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)";

IMAGE_TAG="hydra-xorg-test";

echo "==> Building test container (${IMAGE_TAG})...";
podman build --pull -t "${IMAGE_TAG}" -f "${PROJECT_DIR}/Containerfile" "${PROJECT_DIR}";

echo "";
echo "==> Running integration tests in isolated Xvfb container...";
podman run --rm "${IMAGE_TAG}";

