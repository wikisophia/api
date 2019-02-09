#!/bin/bash

set -e
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
trap 'CMD=${last_command} RET=$?; if [[ $RET -ne 0 ]]; then echo "\"${CMD}\" command failed with exit code $RET."; fi' EXIT
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"

cd $SCRIPTPATH/../client-js
rm -f ~/.npmrc
touch ~/.npmrc
echo "@wikisophia:registry=https://registry.npmjs.org/" >> ~/.npmrc
echo "//registry.npmjs.org/:_authToken=${NPM_TOKEN}" >> ~/.npmrc
npm run build
npm publish
