#!/bin/bash

# This script runs some validation, and then publishes client-js to NPM.
# It's only intended to be run from a Travis build.

set -e
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
trap 'CMD=${last_command} RET=$?; if [[ $RET -ne 0 ]]; then echo "\"${CMD}\" command failed with exit code $RET."; fi' EXIT
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH/../client-js

# This script only publishes tags of the form "client-js-{major}.{minor}.{patch}".
# If the tag doesn't start with "client-js-", it's probably trying to publish another package.
PUBLISH_PACKAGE=$(echo ${TRAVIS_TAG} | cut -c 1-10)
if [[ $PUBLISH_PACKAGE != "client-js-" ]]; then
  exit 0
fi

# If the tag begins with "client-js-" and the rest of it isn't
# a semantic version, the tag is malformed and it should be an error.
TRAVIS_VERSION="$(echo ${TRAVIS_TAG} | cut -c 11-)"
if [[ ! ${TRAVIS_VERSION} =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "\"${TRAVIS_TAG}\" is an invalid tag. It must end with a semantic version. client-js will not be published."
  exit 1
fi

# Make sure the git tag matches the version in client-js/package.json.
# If not, quit early.
PACKAGE_VERSION="$(cat ./package.json | grep version | sed 's/[version": ,]*//' | sed 's/["\, ]*$//')"
if [[ ! ${PACKAGE_VERSION} = ${TRAVIS_VERSION} ]]; then
  echo "Git tag \"${TRAVIS_TAG}\ uses version \"${TRAVIS_VERSION}\", which does not match package.json version \"${PACKAGE_VERSION}\". client-js will not be published"
  exit 1
fi

PACKAGE_NAME="$(cat ./package.json | grep name | sed 's/[name": ,]*//' | sed 's/["\, ]*$//')"
echo "Publishing version ${PACKAGE_VERSION} of \"${PACKAGE_NAME}\" to https://registry.npmjs.org/"
rm -f ~/.npmrc
echo "@wikisophia:registry=https://registry.npmjs.org/" >> ~/.npmrc
echo "//registry.npmjs.org/:_authToken=${NPM_TOKEN}" >> ~/.npmrc
npm run build
npm publish
