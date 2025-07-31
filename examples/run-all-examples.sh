#!/bin/bash

# Run All Examples Script
# 运行所有示例脚本

echo "🚀 Gradle Parser Examples Test Suite"
echo "====================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to run an example
run_example() {
    local dir=$1
    local name=$2
    
    echo -e "${BLUE}📁 Running $name...${NC}"
    echo "Directory: $dir"
    
    if [ ! -d "$dir" ]; then
        echo -e "${RED}❌ Directory $dir not found${NC}"
        return 1
    fi
    
    if [ ! -f "$dir/main.go" ]; then
        echo -e "${RED}❌ main.go not found in $dir${NC}"
        return 1
    fi
    
    cd "$dir"
    
    # Run the example with timeout
    timeout 30s go run main.go
    local exit_code=$?
    
    cd ..
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ $name completed successfully${NC}"
        return 0
    elif [ $exit_code -eq 124 ]; then
        echo -e "${YELLOW}⏰ $name timed out (30s limit)${NC}"
        return 1
    else
        echo -e "${RED}❌ $name failed with exit code $exit_code${NC}"
        return 1
    fi
}

# Check if we're in the examples directory
if [ ! -f "README.md" ] || [ ! -d "sample_files" ]; then
    echo -e "${RED}❌ Please run this script from the examples directory${NC}"
    exit 1
fi

# Initialize counters
total=0
passed=0
failed=0

echo "🔍 Checking prerequisites..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go is available: $(go version)${NC}"

# Check if sample files exist
if [ ! -f "sample_files/build.gradle" ]; then
    echo -e "${RED}❌ Sample files not found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Sample files are available${NC}"
echo ""

# List of examples to run
examples=(
    "01_basic:Basic Parsing Example"
    "02_dependencies:Dependency Analysis Example"
    "03_plugins:Plugin Detection Example"
    "04_repositories:Repository Analysis Example"
    "05_complete:Complete Features Example"
    "06_editor:Structured Editing Example"
    "07_advanced:Advanced Features Example"
)

echo "📋 Found ${#examples[@]} examples to run"
echo ""

# Run each example
for example in "${examples[@]}"; do
    IFS=':' read -r dir name <<< "$example"
    
    echo "----------------------------------------"
    total=$((total + 1))
    
    if run_example "$dir" "$name"; then
        passed=$((passed + 1))
    else
        failed=$((failed + 1))
    fi
    
    echo ""
done

# Summary
echo "========================================"
echo "📊 Test Results Summary"
echo "========================================"
echo -e "Total examples: ${BLUE}$total${NC}"
echo -e "Passed: ${GREEN}$passed${NC}"
echo -e "Failed: ${RED}$failed${NC}"

if [ $failed -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All examples completed successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "• Explore the documentation: ../docs/"
    echo "• Try modifying the examples for your use case"
    echo "• Check out the API reference: ../docs/api/"
    exit 0
else
    echo ""
    echo -e "${YELLOW}⚠️  Some examples failed. Please check the output above.${NC}"
    echo ""
    echo "Troubleshooting:"
    echo "• Make sure all dependencies are installed: go mod tidy"
    echo "• Check that sample files are present in sample_files/"
    echo "• Verify Go version compatibility (1.19+)"
    exit 1
fi
