#!/bin/bash
set -e

# Install tools
go install github.com/swaggo/swag/cmd/swag@v1.16.6
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Regenerate API docs
cd ../api
~/go/bin/swag init -g cmd/api/main.go --output docs

# Fix basePath issue for oapi-codegen (removes "basePath": "/")
sed -i 's/"basePath": "\/",/"basePath": "",/g' docs/swagger.json

# Convert to OpenAPI 3.0 (requires npm)
npx swagger2openapi -o docs/openapi.yaml docs/swagger.json

# Generate Client
cd ../agent
~/go/bin/oapi-codegen -package client -generate types,client -o pkg/client/client.go ../api/docs/openapi.yaml
go mod tidy
