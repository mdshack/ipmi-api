#!/bin/bash
set -e

VERSION=${1:-latest}
REGISTRY=${DOCKER_REGISTRY:-"mdshack"}

echo "Building IPMI API Docker image..."

# Build production image
docker build -t ipmi-api:${VERSION} .

# Tag with registry if provided
if [ ! -z "$REGISTRY" ]; then
    docker tag ipmi-api:${VERSION} ${REGISTRY}/ipmi-api:${VERSION}
    echo "Tagged as ${REGISTRY}/ipmi-api:${VERSION}"
fi

echo "Build complete: ipmi-api:${VERSION}"

# Optional: Run security scan
if command -v trivy &> /dev/null; then
    echo "Running security scan..."
    trivy image ipmi-api:${VERSION}
fi