# Goodbye World API

A simple Rust web API built with Axum that returns a "Goodbye, World!" message.

## Features

- Built with Rust and Axum framework
- Dockerized for easy deployment
- Exposes a `/goodbye` endpoint
- Returns JSON response

## Prerequisites

- Docker installed and running
- Rust toolchain (for local development)

## Quick Start with Docker

### 1. Build the Docker Image

```bash
docker build -t goodbye-app .
```

### 2. Run the Container

```bash
docker run -d -p 3000:3000 --name goodbye-container goodbye-app
```

### 3. Test the API

```bash
curl http://localhost:3000/goodbye
```

Expected response:

```json
{ "message": "Goodbye, World!", "status": "success" }
```

### 4. Stop and Remove the Container

```bash
docker stop goodbye-container
docker rm goodbye-container
```

## Development

### Local Development

1. Install Rust dependencies:

```bash
cargo build
```

2. Run locally:

```bash
cargo run
```

3. Test the endpoint:

```bash
curl http://localhost:3000/goodbye
```

### Project Structure

```
goodbye/
├── Cargo.toml          # Rust dependencies and package config
├── Dockerfile          # Multi-stage Docker build
├── src/
│   └── main.rs         # Main application code
└── README.md           # This file
```

## API Endpoints

- `GET /goodbye` - Returns a goodbye message in JSON format

## Docker Details

The Dockerfile uses a multi-stage build approach:

1. **Builder Stage**: Uses `rust:slim` to compile the Rust application
2. **Runtime Stage**: Uses the same `rust:slim` image for consistency and minimal size

## AWS ECS Deployment

This application can be deployed to AWS ECS Fargate for production use.

### Prerequisites

- AWS CLI configured with appropriate permissions
- ECR repository for the Docker image
- ECS cluster and service configured
- VPC with public subnets and internet gateway

### Deployment Steps

1. **Push to ECR**:

```bash
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com
docker tag goodbye-app:latest $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/goodbye-world:latest
docker push $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/goodbye-world:latest
```

2. **Update ECS Service** (if needed):

```bash
aws ecs update-service \
  --cluster "goodbye-cluster" \
  --service "goodbye-service" \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=ENABLED}" \
  --region $AWS_REGION
```

3. **Check Service Status**:

```bash
aws ecs describe-services \
  --cluster "goodbye-cluster" \
  --services "goodbye-service" \
  --region $AWS_REGION
```

### Troubleshooting: Public IP Not Assigned

**Problem**: ECS Fargate task running but no public IP assigned.

**Common Causes**:

- Missing security group in ECS service network configuration
- Security group doesn't allow inbound traffic on port 3000
- VPC subnets don't have auto-assign public IP enabled

**Solution**:

1. **Verify Security Group Configuration**:

```bash
# Check if security group exists and allows port 3000
aws ec2 describe-security-groups \
  --filters "Name=vpc-id,Values=vpc-xxx" \
  --region $AWS_REGION
```

2. **Update ECS Service with Security Group**:

```bash
aws ecs update-service \
  --cluster "goodbye-cluster" \
  --service "goodbye-service" \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=ENABLED}" \
  --region $AWS_REGION
```

3. **Check Task Public IP**:

```bash
# Get current task ARN
TASK_ARN=$(aws ecs list-tasks \
  --cluster "goodbye-cluster" \
  --service-name "goodbye-service" \
  --region $AWS_REGION \
  --query "taskArns[0]" \
  --output text)

# Check public IP via network interface
aws ecs describe-tasks \
  --cluster "goodbye-cluster" \
  --tasks $TASK_ARN \
  --region $AWS_REGION \
  --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
  --output text

# Get public IP from network interface
aws ec2 describe-network-interfaces \
  --network-interface-ids eni-xxx \
  --region $AWS_REGION \
  --query "NetworkInterfaces[0].Association.PublicIp" \
  --output text
```

4. **Test Application**:

```bash
curl "http://PUBLIC_IP:3000/"
curl "http://PUBLIC_IP:3000/goodbye"
```

**Required ECS Service Network Configuration**:

```json
{
  "awsvpcConfiguration": {
    "subnets": ["subnet-xxx", "subnet-yyy"],
    "securityGroups": ["sg-xxx"],
    "assignPublicIp": "ENABLED"
  }
}
```

**Required Security Group Rules**:

- **Inbound**: Port 3000, Protocol TCP, Source 0.0.0.0/0
- **Outbound**: All traffic allowed (default)

## Troubleshooting

### Container Not Starting

- Check if port 3000 is available: `lsof -i :3000`
- Verify the image was built successfully: `docker images goodbye-app`

### Build Failures

- Ensure Docker has sufficient disk space
- Clear Docker cache if needed: `docker system prune -a`

### Port Conflicts

- Use a different port: `docker run -d -p 8080:3000 --name goodbye-container goodbye-app`
- Then access via: `curl http://localhost:8080/goodbye`

## License

This project is open source and available under the MIT License.
