#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output="${1:-${repo_root}/dist/whatsMyIp.zip}"
workdir="$(mktemp -d)"

cleanup() {
  rm -rf "${workdir}"
}
trap cleanup EXIT

mkdir -p "$(dirname "${output}")"

(
  cd "${repo_root}"
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "${workdir}/bootstrap" ./cmd/whatsmyip
)

(
  cd "${workdir}"
  zip -q -X "${output}" bootstrap
)
