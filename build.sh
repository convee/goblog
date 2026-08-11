#!/usr/bin/env bash
set -euo pipefail

case "${1:-}" in
build)
  exec make build
  ;;
tar)
  exec make tar
  ;;
*)
  echo "usage: $0 {build|tar}" >&2
  exit 1
  ;;
esac
