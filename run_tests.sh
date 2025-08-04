#!/bin/bash

echo "🚀 Running Unit Tests for Ecom_MiniGo Project"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run tests for a service
run_service_tests() {
    local service_name=$1
    local service_path=$2
    
    echo -e "\n${YELLOW}Testing $service_name...${NC}"
    
    if [ -d "$service_path" ]; then
        cd "$service_path"
        
        # Check if there are test files
        if find . -name "*_test.go" | grep -q .; then
            echo "Found test files in $service_name"
            
            # Run tests with coverage
            go test -v -cover ./...
            
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}✅ $service_name tests passed${NC}"
            else
                echo -e "${RED}❌ $service_name tests failed${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  No test files found in $service_name${NC}"
        fi
        
        cd - > /dev/null
    else
        echo -e "${RED}❌ Directory $service_path not found${NC}"
    fi
}

# Run tests for each service
run_service_tests "User Service" "user_service/service"
run_service_tests "Product Service" "product_service/service"
run_service_tests "Order Service" "order_service/service"
run_service_tests "API Gateway" "api-gateway"

echo -e "\n${GREEN}🎉 All tests completed!${NC}"

# Optional: Generate coverage report
echo -e "\n${YELLOW}Generating coverage report...${NC}"
echo "To view detailed coverage, run:"
echo "go test -coverprofile=coverage.out ./..."
echo "go tool cover -html=coverage.out -o coverage.html"
echo "open coverage.html" 