#!/bin/bash

set -e

cd "$(dirname "$0")"

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen -package api -generate types xmproject.yml > ../internal/api/xmproject.gen.go