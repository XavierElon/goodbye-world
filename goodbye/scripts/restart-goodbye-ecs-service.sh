#!/bin/bash

echo "🔄 Resetting ECS service..."

# 1. Stop and delete service
echo "📤 Stopping service..."
aws ecs update-service \
    --cluster "goodbye-cluster" \
    --service "goodbye-service" \
    --desired-count 0 \
    --region $AWS_REGION

echo "⏳ Waiting for tasks to stop..."
sleep 30

echo "��️  Deleting service..."
aws ecs delete-service \
    --cluster "goodbye-cluster" \
    --service "goodbye-service" \
    --force \
    --region $AWS_REGION

# 2. Build and push new image
echo "🏗️  Building Docker image..."
docker build --platform linux/amd64 -t goodbye-world .

echo "🏷️  Tagging image..."
docker tag goodbye-world:latest $ECR_IMAGE_URI

echo "⬆️  Pushing to ECR..."
docker push $ECR_IMAGE_URI

# 3. Recreate service
echo "🚀 Creating new service..."
aws ecs create-service \
    --cluster "goodbye-cluster" \
    --service-name "goodbye-service" \
    --task-definition "goodbye-task" \
    --desired-count 1 \
    --launch-type "FARGATE" \
    --network-configuration "awsvpcConfiguration={subnets=[$SUBNET_1,$SUBNET_2],securityGroups=[$SG_ID],assignPublicIp=ENABLED}" \
    --region $AWS_REGION

echo "✅ Service reset complete!"
echo "🌐 Check AWS Console for status"