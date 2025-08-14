#!/bin/bash

echo "Generating Swagger documentation..."

# Remove old docs
rm -rf docs/

# Generate new docs
swag init -g cmd/server/main.go --parseDependency --parseInternal

echo "Swagger docs generated successfully!"
echo "Files created:"
ls -la docs/ 