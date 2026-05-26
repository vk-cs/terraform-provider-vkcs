#!/bin/bash

REPO_DIR="$(CDPATH='' cd "$(dirname "$0")"/.. && pwd)"
cd "${REPO_DIR}"

# shellcheck disable=SC2034 # documented constant, not consumed inside this script
DEFAULT_BRANCH="master"

if [[ "$(uname)" == "Darwin" ]]; then
  echo "(Using BSD sed)"
  SED="sed -E"
else
  echo "(Using GNU sed)"
  SED="sed -r"
fi

ALLOW_DIRTY="${ALLOW_DIRTY:-}"
for arg in "$@"; do
  case "${arg}" in
    --allow-dirty) ALLOW_DIRTY=1 ;;
    *) echo "Unknown argument: ${arg}" >&2; exit 64 ;;
  esac
done

if [[ -z "${ALLOW_DIRTY}" && "$(git status --short)" != "" ]]; then
  echo "Error: working tree is dirty (pass --allow-dirty or set ALLOW_DIRTY=1 to bypass)" >&2
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
