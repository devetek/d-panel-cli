#!/bin/bash

# setup to fetch private packages
echo ${NETRC} > ~/.netrc
chmod og-rw ~/.netrc

go build -ldflags "-X main.currentVersion=${RELEASE_TAG}" -o dpid ./cmd/cli