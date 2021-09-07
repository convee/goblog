#!/bin/sh

COMMIT=$(git rev-parse --short HEAD)
TIME=$(date +%Y%m%d%H)

case "$1" in
'build')
  GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.Commit=$COMMIT"
  ;;
'tar')
  tar zcvf  blog"$TIME".tar.gz ./blog ./tpl ./static ./conf
  ;;
*)
  echo "usage: $0 {build|tar}"
  exit 1
  ;;
esac
