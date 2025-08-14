#!/bin/bash

# Simple wrapper script for redeployment
# This script sources environment variables and calls the main redeployment script

# Source the environment variables
if [[ -f "scripts/environment.sh" ]]; then
    echo "📁 Loading environment variables from scripts/environment.sh..."
    source scripts/environment.sh
else
    echo "⚠️  Environment file scripts/environment.sh not found!"
    echo "   Please create it based on env-template.txt"
    exit 1
fi

# Call the main redeployment script
echo "🚀 Starting redeployment process..."
./scripts/redeploy.sh
