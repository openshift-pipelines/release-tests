#!/usr/bin/env bash
# Get an asset from a GitHub project release.
# First arg is the repo name
# Second arg is the binary substring, i.e: darwin or linux or x86_64-linux
# etc.. a subtring in the binary we want so we grab the one we want
set -eu
REPO=${1:-}
BINARY_SUBSTRING=${2:-}

TMP=$(mktemp /tmp/.mm.XXXXXX)
clean() { rm -f ${TMP}; }
trap clean EXIT

[[ -z ${REPO} || -z ${BINARY_SUBSTRING} ]] && {
    echo "Need a ${REPO} ${BINARY_SUBSTRING}"
    exit 1
}

api_url=https://api.github.com/repos/${REPO}/releases/latest
curl -f -L -s ${api_url} >${TMP}
latest_version=$(cat ${TMP} | python -c "import sys, json;x=json.load(sys.stdin);print(x['tag_name'])")
[[ -z ${latest_version} ]] && {
    echo "Could not find the latest version in ${api_url}"
    exit 1
}

asset=$(cat ${TMP} | \
            python -c "import sys, json;x=json.load(sys.stdin);print([ r['browser_download_url'] for r in x['assets'] if '${BINARY_SUBSTRING}' in r['name']][0])")

[[ -z ${asset} ]] && {
    echo "Could not find an asset named ${asset} in ${api_url}"
    exit 1
}

echo ${asset}