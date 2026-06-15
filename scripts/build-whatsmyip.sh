#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output="${1:-${repo_root}/dist/whatsMyIp.zip}"
workdir="$(mktemp -d)"
source_date_epoch="$(git -C "${repo_root}" log -1 --format=%ct -- go.mod go.sum cmd/whatsmyip)"

cleanup() {
  rm -rf "${workdir}"
}
trap cleanup EXIT

mkdir -p "$(dirname "${output}")"
(
  cd "${repo_root}"
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o "${workdir}/bootstrap" ./cmd/whatsmyip
  perl -e 'utime $ARGV[0], $ARGV[0], $ARGV[1] or die "utime: $!\\n"' "${source_date_epoch}" "${workdir}/bootstrap"
)

(
  cd "${workdir}"
  zip -q -X "${output}" bootstrap
)
