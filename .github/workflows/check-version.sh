#!/bin/bash

set -ex

TAG="${GITHUB_REF_NAME:-${GITHUB_REF#refs/tags/}}"
if [ -z "$TAG" ] || [ "$TAG" = "$GITHUB_REF" ]; then
  echo "Unable to determine tag from GitHub Actions environment"
  exit 1
fi
FILE_VERSION=$(sed -nE 's/var Version = "(.*)"/\1/p' cmd/version.go)

if [ "$TAG" != "v$FILE_VERSION" ]; then
  echo "Version mismatch: tag is $TAG but version.go has $FILE_VERSION"
  exit 1
fi