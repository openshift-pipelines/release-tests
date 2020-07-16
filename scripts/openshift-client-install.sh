#!/usr/bin/env bash
set -e

targetdir=/usr/local/bin
version=latest

URL=https://mirror.openshift.com/pub/openshift-v4/clients/ocp/${version}

versionnumber=$(curl f -s ${URL}/release.txt |sed -n '/Version:/ { s/.*:[ ]*//; p ;}')

[[ -z ${versionnumber} ]] && {
    echo "Could not detect version"
    exit 1
}

mkdir -p ${targetdir}

case $(uname -o) in
    *Linux)
        platform=linux
        ;;
    Darwin)
        platform=mac
        ;;
    *)
        echo "Could not detect platform: $(uname -o)"
        exit 1
esac
platform=linux

echo -n "Downloading openshift-clients-${version}: "
curl -s -L ${URL}/openshift-client-${platform}-${versionnumber}.tar.gz|tar -xzf- -C ${targetdir} oc kubectl
echo "Done."