#! /bin/bash

set -ex

TAG=$(git describe --tags --abbrev=0)
FILE_VERSION=$(sed -nE 's/var Version = "(.*)"/\1/p' cmd/version.go)

if [ "$TAG" != "v$FILE_VERSION" ]; then
  echo "Version mismatch: tag is $TAG but version.go has $FILE_VERSION"
  exit 1
fi