#!/bin/bash

REPO_DIR="$(cd "$(dirname "$0")"/.. && pwd)"
cd "${REPO_DIR}"

DEFAULT_BRANCH="master"

if [[ "$(uname)" == "Darwin" ]]; then
  echo "(Using BSD sed)"
  SED="sed -E"
else
  echo "(Using GNU sed)"
  SED="sed -r"
fi

if [[ "$(git status --short)" != "" ]]; then
  echo "Error: working tree is dirty" >&2
  exit 4
fi

set -e

echo "Preparing changelog for release..."

if [[ ! -f CHANGELOG.md ]]; then
  echo "Error: CHANGELOG.md not found."
  exit 2
fi

# Remove unreleased mark
( set -x; $SED -i.bak "s/ \(unreleased\)//" CHANGELOG.md )

rm CHANGELOG.md.bak

echo "Updating provider schema JSON..."
(
  set -x
  make update_release_schema
)

echo "Running go generate..."
(
    set -x
    make generate
)
