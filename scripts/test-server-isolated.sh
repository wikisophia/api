#!/bin/bash

set -e
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
trap 'CMD=${last_command} RET=$?; if [[ $RET -ne 0 ]]; then echo "\"${CMD}\" command failed with exit code $RET."; fi' EXIT
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd ${SCRIPTPATH}/../server

go test ./... -count=1

# gofmt always returns 0, and doesn't support ./... syntax. So run through the files/directories
# manually and capture any output for errors.
FMT=$(gofmt -l -s $(ls -d */ | grep -v "vendor") *.go)
if ! [ -z ${FMT} ]; then
    echo ''
    echo "Some files have bad style. Run the following commands to fix them:"
    for LINE in ${FMT}
    do
        echo "  gofmt -s -w `pwd`/${LINE}"
    done
    echo ''
    exit 1
fi
