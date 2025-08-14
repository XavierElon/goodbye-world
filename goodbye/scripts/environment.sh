#!/bin/bash

# Environment variables for redeployment script
# Source this file before running the deployment script

# AWS Configuration
export AWS_ACCOUNT_ID=275136276893
export AWS_REGION=us-west-1

# ECR Configuration
export ECR_REPOSITORY_NAME=goodbye-world
export ECR_IMAGE_URI=275136276893.dkr.ecr.us-west-1.amazonaws.com/goodbye-world

# ECS Configuration
export CLUSTER_NAME=goodbye-cluster
export SERVICE_NAME=goodbye-service
export TASK_DEFINITION_FAMILY=goodbye-task

# Optional: Network Configuration (if you need to recreate the service)
# export SUBNET_1=subnet-xxxxxxxxx
# export SUBNET_2=subnet-yyyyyyyyy
# export SG_ID=sg-zzzzzzzzz
