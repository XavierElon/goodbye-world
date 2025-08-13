#!/bin/bash
export AWS_ACCOUNT_ID="275136276893"
export AWS_REGION="us-west-1"
export REPOSITORY_NAME="goodbye-world"
export ECR_IMAGE_URI="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$REPOSITORY_NAME:latest"

echo "Environment variables set:"
echo "AWS_ACCOUNT_ID: $AWS_ACCOUNT_ID"
echo "AWS_REGION: $AWS_REGION"
echo "REPOSITORY_NAME: $REPOSITORY_NAME"
echo "ECR_IMAGE_URI: $ECR_IMAGE_URI"
