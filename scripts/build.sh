#!/bin/bash

go build -ldflags "-X main.currentVersion=${RELEASE_TAG}" -o dpid ./cmd/cli