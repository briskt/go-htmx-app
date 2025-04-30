#!/usr/bin/env bash

# exit if any line fails
set -e

# print script lines as they are executed
set -x

sqlc generate
templ generate
npx --yes @tailwindcss/cli -i ./public/assets/css/tailwind.css -o ./public/assets/css/tailwind.output.css --minify
go build -o ./bin/app ./cmd/server
