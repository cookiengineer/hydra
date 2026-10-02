#!/usr/bin/env bash

# Hydra test runner. Delegates to the centralized Makefile.
#
#   bash test.sh unit;
#   bash test.sh integration;
#   bash test.sh e2e hydratwo;
#   bash test.sh all;

set -euo pipefail;

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)";

MODE="${1:-all}";

case "${MODE}" in

    unit)
        make -C "${PROJECT_DIR}" unit;
        ;;

    integration)
        make -C "${PROJECT_DIR}" integration;
        ;;

    e2e)
        make -C "${PROJECT_DIR}" e2e HOST="${2:-hydratwo}";
        ;;

    all)
        make -C "${PROJECT_DIR}" unit;
        echo "";
        make -C "${PROJECT_DIR}" integration;
        ;;

    *)
        echo "Usage: bash test.sh [unit|integration|e2e|all] [host]";
        exit 1;
        ;;

esac
