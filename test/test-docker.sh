#!/bin/bash
# test-docker.sh - Docker testing script for timezones API
# Follows CLAUDE.md/SPEC guidelines for isolated testing

set -e

PROJECTNAME="timezones"
TESTPORT=$(shuf -i 64000-64999 -n 1)

echo "🧪 Testing ${PROJECTNAME} using Docker"
echo "📡 Test Port: ${TESTPORT}"
echo ""

# Build development image
echo "📦 Building development Docker image..."
make docker-dev

# Run container
echo "🚀 Starting test container..."
docker run -d \
  --name ${PROJECTNAME}-test-${TESTPORT} \
  -p ${TESTPORT}:80 \
  -v /tmp/${PROJECTNAME}-test:/data \
  -e ADMIN_USER=testadmin \
  -e ADMIN_PASSWORD=testpass123 \
  ${PROJECTNAME}:dev

# Wait for startup
echo "⏳ Waiting for service to start..."
sleep 5

# Test health endpoint
echo "🔍 Testing health endpoint..."
if curl -sf http://localhost:${TESTPORT}/healthz > /dev/null; then
  echo "✅ Health check passed"
else
  echo "❌ Health check failed"
  docker logs ${PROJECTNAME}-test-${TESTPORT}
  docker stop ${PROJECTNAME}-test-${TESTPORT}
  docker rm ${PROJECTNAME}-test-${TESTPORT}
  rm -rf /tmp/${PROJECTNAME}-test
  exit 1
fi

# Test API endpoint
echo "🔍 Testing timezones API..."
if curl -sf http://localhost:${TESTPORT}/api/v1/timezones.json > /dev/null; then
  echo "✅ API endpoint passed"
else
  echo "❌ API endpoint failed"
  docker logs ${PROJECTNAME}-test-${TESTPORT}
  docker stop ${PROJECTNAME}-test-${TESTPORT}
  docker rm ${PROJECTNAME}-test-${TESTPORT}
  rm -rf /tmp/${PROJECTNAME}-test
  exit 1
fi

# Test homepage
echo "🔍 Testing homepage..."
if curl -sf http://localhost:${TESTPORT}/ > /dev/null; then
  echo "✅ Homepage passed"
else
  echo "❌ Homepage failed"
  docker logs ${PROJECTNAME}-test-${TESTPORT}
  docker stop ${PROJECTNAME}-test-${TESTPORT}
  docker rm ${PROJECTNAME}-test-${TESTPORT}
  rm -rf /tmp/${PROJECTNAME}-test
  exit 1
fi

# Show container logs
echo ""
echo "📄 Container logs:"
docker logs ${PROJECTNAME}-test-${TESTPORT}

# Cleanup
echo ""
echo "🧹 Cleaning up..."
docker stop ${PROJECTNAME}-test-${TESTPORT}
docker rm ${PROJECTNAME}-test-${TESTPORT}
rm -rf /tmp/${PROJECTNAME}-test

echo ""
echo "✅ All tests passed!"
echo "🎉 Docker test completed successfully"
