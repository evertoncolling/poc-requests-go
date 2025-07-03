#!/bin/bash
set -e

# Script to run integration tests with proper environment loading

echo "🧪 Running integration tests..."

# Check if .env file exists
if [ -f ".env" ]; then
    echo "📁 Loading environment variables from .env file..."
    set -a
    source .env
    set +a
else
    echo "⚠️  No .env file found, using existing environment variables..."
fi

# Check if required environment variables are set
required_vars=("CLIENT_ID" "CLIENT_SECRET" "TENANT_ID" "CDF_CLUSTER" "CDF_PROJECT")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -ne 0 ]; then
    echo "❌ Missing required environment variables: ${missing_vars[*]}"
    echo ""
    echo "For local development, create a .env file with:"
    echo "CLIENT_ID=your_client_id"
    echo "CLIENT_SECRET=your_client_secret"
    echo "TENANT_ID=your_tenant_id"
    echo "CDF_CLUSTER=your_cluster_name"
    echo "CDF_PROJECT=your_project_name"
    echo ""
    echo "For CI/CD, set these as GitHub repository secrets."
    exit 1
fi

echo "✅ All required environment variables are set"
echo "🏃 Running integration tests..."

# Run the integration tests
go test -v -run TestIntegration ./pkg/api/

echo "✅ Integration tests completed successfully!"