#!/usr/bin/env bash

# TODO: The below steps don't really work
# LATEST_VERSION=$(git tag -l | tail -1)
# LATEST_MINOR=$(echo $LATEST_VERSION | sed 's/v1.0.//')
# NEW_MINOR="$(($LATEST_MINOR + 1))"
# TODO: See above - NEW_MINOR should be calculated but currently you should set it manually
NEW_MINOR=15
NEW_VERSION=v1.0.$NEW_MINOR
echo $NEW_VERSION

git tag $NEW_VERSION
git push origin tags/$NEW_VERSION

GOOS=darwin go build -o build/jira-ticket-mac_$NEW_VERSION
GOOS=linux go build -o build/jira-ticket-linux_$NEW_VERSION
