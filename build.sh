#!/usr/bin/env bash

# exit if any line fails
set -e

# print script lines as they are executed
set -x

sqlc generate
go build -o ./bin/app ./cmd/server
