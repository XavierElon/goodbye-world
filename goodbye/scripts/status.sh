#!/bin/bash

# Quick status check script for ECS service
# This script shows the current status without making any changes

# Source environment variables
if [[ -f "scripts/environment.sh" ]]; then
    source scripts/environment.sh
else
    echo "❌ Environment file not found. Please run from the project root."
    exit 1
fi

echo "🔍 Checking ECS service status..."
echo "=================================="

# Get service information
aws ecs describe-services \
    --cluster "$CLUSTER_NAME" \
    --services "$SERVICE_NAME" \
    --region "$AWS_REGION" \
    --query 'services[0]' \
    --output json | jq '{
        serviceName: .serviceName,
        status: .status,
        desiredCount: .desiredCount,
        runningCount: .runningCount,
        pendingCount: .pendingCount,
        taskDefinition: .taskDefinition,
        deployments: [.deployments[] | {
            id: .id,
            status: .status,
            rolloutState: .rolloutState,
            rolloutStateReason: .rolloutStateReason,
            taskDefinition: .taskDefinition,
            runningCount: .runningCount,
            desiredCount: .desiredCount
        }]
    }'
