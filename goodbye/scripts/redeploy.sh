#!/bin/bash

# Redeploy script for Goodbye Rust application
# This script builds the Docker image, pushes to ECR, and updates ECS

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required environment variables are set
check_env_vars() {
    local required_vars=("AWS_REGION" "ECR_IMAGE_URI" "ECR_REPOSITORY_NAME" "CLUSTER_NAME" "SERVICE_NAME" "TASK_DEFINITION_FAMILY")
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            print_error "Environment variable $var is not set"
            exit 1
        fi
    done
    
    print_success "All required environment variables are set"
}

# Check if AWS CLI is installed and configured
check_aws_cli() {
    if ! command -v aws &> /dev/null; then
        print_error "AWS CLI is not installed. Please install it first."
        exit 1
    fi
    
    if ! aws sts get-caller-identity &> /dev/null; then
        print_error "AWS CLI is not configured. Please run 'aws configure' first."
        exit 1
    fi
    
    print_success "AWS CLI is configured"
}

# Check if Docker is running
check_docker() {
    if ! docker info &> /dev/null; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    
    print_success "Docker is running"
}

# Build Docker image
build_image() {
    print_status "Building Docker image..."
    
    # Build with platform specification for ECS compatibility
    docker build --platform linux/amd64 -t "$ECR_REPOSITORY_NAME:latest" .
    
    if [[ $? -eq 0 ]]; then
        print_success "Docker image built successfully"
    else
        print_error "Failed to build Docker image"
        exit 1
    fi
}

# Tag image for ECR
tag_image() {
    print_status "Tagging image for ECR..."
    
    docker tag "$ECR_REPOSITORY_NAME:latest" "$ECR_IMAGE_URI:latest"
    
    if [[ $? -eq 0 ]]; then
        print_success "Image tagged successfully"
    else
        print_error "Failed to tag image"
        exit 1
    fi
}

# Push image to ECR
push_to_ecr() {
    print_status "Pushing image to ECR..."
    
    # Get ECR login token
    print_status "Getting ECR login token..."
    aws ecr get-login-password --region "$AWS_REGION" | docker login --username AWS --password-stdin "$ECR_IMAGE_URI"
    
    # Push the image
    docker push "$ECR_IMAGE_URI:latest"
    
    if [[ $? -eq 0 ]]; then
        print_success "Image pushed to ECR successfully"
    else
        print_error "Failed to push image to ECR"
        exit 1
    fi
}

# Update ECS task definition
update_task_definition() {
    print_status "Updating ECS task definition..."
    
    # Get the current task definition
    local current_task_def=$(aws ecs describe-task-definition \
        --task-definition "$TASK_DEFINITION_FAMILY" \
        --region "$AWS_REGION" \
        --query 'taskDefinition' \
        --output json)
    
    # Filter out read-only fields that can't be used when registering a new task definition
    local updated_task_def=$(echo "$current_task_def" | \
        jq 'del(.taskDefinitionArn, .revision, .status, .requiresAttributes, .compatibilities, .registeredAt, .registeredBy)' | \
        jq --arg IMAGE "$ECR_IMAGE_URI:latest" \
        '.containerDefinitions[0].image = $IMAGE')
    
    # Register new task definition
    local new_task_def_arn=$(aws ecs register-task-definition \
        --cli-input-json "$updated_task_def" \
        --region "$AWS_REGION" \
        --query 'taskDefinition.taskDefinitionArn' \
        --output text)
    
    if [[ $? -eq 0 ]]; then
        print_success "Task definition updated: $new_task_def_arn"
        NEW_TASK_DEF_ARN="$new_task_def_arn"
    else
        print_error "Failed to update task definition"
        exit 1
    fi
}

# Update ECS service
update_ecs_service() {
    print_status "Updating ECS service..."
    
    # Update the service with the new task definition
    aws ecs update-service \
        --cluster "$CLUSTER_NAME" \
        --service "$SERVICE_NAME" \
        --task-definition "$NEW_TASK_DEF_ARN" \
        --region "$AWS_REGION"
    
    if [[ $? -eq 0 ]]; then
        print_success "ECS service updated successfully"
    else
        print_error "Failed to update ECS service"
        exit 1
    fi
}

# Wait for service to stabilize
wait_for_service_stability() {
    print_status "Waiting for service to stabilize..."
    
    # Set a timeout for the wait command (5 minutes)
    timeout 300 aws ecs wait services-stable \
        --cluster "$CLUSTER_NAME" \
        --services "$SERVICE_NAME" \
        --region "$AWS_REGION" || {
        print_warning "Service stability check timed out after 5 minutes"
        print_status "Checking current service status..."
    }
    
    # Always check the final status regardless of timeout
    get_service_status
}

# Get service status
get_service_status() {
    print_status "Getting service status..."
    
    local service_info=$(aws ecs describe-services \
        --cluster "$CLUSTER_NAME" \
        --services "$SERVICE_NAME" \
        --region "$AWS_REGION" \
        --query 'services[0]' \
        --output json)
    
    local service_status=$(echo "$service_info" | jq -r '.status')
    local desired_count=$(echo "$service_info" | jq -r '.desiredCount')
    local running_count=$(echo "$service_info" | jq -r '.runningCount')
    local pending_count=$(echo "$service_info" | jq -r '.pendingCount')
    
    print_status "Service Status: $service_status"
    print_status "Desired Count: $desired_count"
    print_status "Running Count: $running_count"
    print_status "Pending Count: $pending_count"
    
    # Check if deployment is successful
    local deployments=$(echo "$service_info" | jq -r '.deployments[] | select(.status == "PRIMARY") | .rolloutState')
    if [[ "$deployments" == "COMPLETED" ]]; then
        print_success "Deployment completed successfully!"
    elif [[ "$deployments" == "IN_PROGRESS" ]]; then
        print_warning "Deployment still in progress..."
    else
        print_status "Deployment state: $deployments"
    fi
}

# Main deployment function
main() {
    echo "🚀 Starting Goodbye Rust application redeployment..."
    echo "=================================================="
    
    # Pre-flight checks
    check_env_vars
    check_aws_cli
    check_docker
    
    # Build and push
    build_image
    tag_image
    push_to_ecr
    
    # Deploy to ECS
    update_task_definition
    update_ecs_service
    wait_for_service_stability
    get_service_status
    
    echo "=================================================="
    print_success "Redeployment completed successfully! 🎉"
    print_status "Your application should now be running with the latest code."
}

# Run main function
main "$@"
