#!/bin/bash
# This script is used to release a new version of the project.
# It will generate the documentation, commit the changes, push the changes to the remote repository and tag the release.
# The first argument is the version number and the second argument is the commit message.
if [[ -z "$1" ]]
then
  exit
fi
if [[ -z "$2" ]]
then
  exit
fi
git commit -am "$2"
git push origin
git tag -a -f "$1" -m "$2"
git push origin -f "$1"
