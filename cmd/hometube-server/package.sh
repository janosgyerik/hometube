#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

npm --prefix app run build

go build

cmdname=$(basename "$PWD")

if [[ ${GOOS+x} && ${GOARCH+x} ]]; then
    zipfile=$cmdname-$GOOS-$GOARCH.zip
else
    zipfile=$cmdname.zip
fi

rm -f "$zipfile"
zip -r "$zipfile" hometube-server app/public

rm -f "$cmdname"
