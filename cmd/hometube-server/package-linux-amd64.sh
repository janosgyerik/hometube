#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

GOOS=linux GOARCH=amd64 ./package.sh
