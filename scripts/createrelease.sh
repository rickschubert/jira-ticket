#!/usr/bin/env bash

TAG_VERSION=$(git describe --tags)
GOOS=darwin go build -o build/jira-ticket-mac_$TAG_VERSION
GOOS=linux go build -o build/jira-ticket-linux_$TAG_VERSION
