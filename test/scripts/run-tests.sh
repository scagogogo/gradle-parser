#!/bin/bash

# Gradle Parser Test Suite Runner
# 运行 Gradle Parser 测试套件

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
VERBOSE=${GRADLE_PARSER_TEST_VERBOSE:-false}
COVERAGE=${GRADLE_PARSER_TEST_COVERAGE:-true}
BENCHMARKS=${GRADLE_PARSER_TEST_BENCHMARKS:-false}
TIMEOUT=${GRADLE_PARSER_TEST_TIMEOUT:-30s}

echo -e "${BLUE}🧪 Gradle Parser Test Suite${NC}"
echo "=================================="
echo ""

# Function to print section headers
print_section() {
    echo -e "${PURPLE}$1${NC}"
    echo "$(printf '%.0s-' {1..50})"
}

# Function to run tests with proper error handling
run_test() {
    local test_path=$1
    local test_name=$2
    local extra_flags=$3
    
    echo -e "${BLUE}Running $test_name...${NC}"
    
    if [ "$VERBOSE" = "true" ]; then
        extra_flags="$extra_flags -v"
    fi
    
    if go test -timeout="$TIMEOUT" $extra_flags "$test_path"; then
        echo -e "${GREEN}✅ $test_name passed${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name failed${NC}"
        return 1
    fi
}

# Check if we're in the right directory
if [ ! -f "README.md" ] || [ ! -d "../pkg" ]; then
    echo -e "${RED}❌ Please run this script from the test directory${NC}"
    exit 1
fi

# Check prerequisites
print_section "🔍 Checking Prerequisites"

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go is available: $(go version)${NC}"

# Check if we can access the main package
if ! go list ../pkg/api &> /dev/null; then
    echo -e "${RED}❌ Cannot access gradle-parser package${NC}"
    echo "Make sure you're in the test directory of the gradle-parser project"
    exit 1
fi

echo -e "${GREEN}✅ Gradle Parser package is accessible${NC}"
echo ""

# Initialize counters
total_tests=0
passed_tests=0
failed_tests=0

# Unit Tests
print_section "🔬 Unit Tests"
if [ -d "unit" ]; then
    total_tests=$((total_tests + 1))
    if run_test "./unit/..." "Unit Tests" ""; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
else
    echo -e "${YELLOW}⚠️  Unit test directory not found, skipping${NC}"
fi
echo ""

# Integration Tests
print_section "🔗 Integration Tests"
if [ -d "integration" ]; then
    total_tests=$((total_tests + 1))
    if run_test "./integration/..." "Integration Tests" ""; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
    fi
else
    echo -e "${YELLOW}⚠️  Integration test directory not found, skipping${NC}"
fi
echo ""

# Coverage Report
if [ "$COVERAGE" = "true" ]; then
    print_section "📊 Coverage Report"
    echo -e "${BLUE}Generating coverage report...${NC}"
    
    if go test -coverprofile=coverage.out ./... &> /dev/null; then
        if command -v go &> /dev/null; then
            coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
            echo -e "${GREEN}✅ Total coverage: $coverage${NC}"
            
            # Generate HTML report
            if go tool cover -html=coverage.out -o coverage.html &> /dev/null; then
                echo -e "${GREEN}✅ HTML coverage report generated: coverage.html${NC}"
            fi
        fi
    else
        echo -e "${YELLOW}⚠️  Coverage report generation failed${NC}"
    fi
    echo ""
fi

# Benchmarks
if [ "$BENCHMARKS" = "true" ]; then
    print_section "🏃 Performance Benchmarks"
    echo -e "${BLUE}Running performance benchmarks...${NC}"
    
    if [ -d "integration" ]; then
        if go test -bench=. -benchmem ./integration/ > benchmark.out 2>&1; then
            echo -e "${GREEN}✅ Benchmarks completed${NC}"
            echo "Results saved to benchmark.out"
            
            # Show summary
            if grep -q "Benchmark" benchmark.out; then
                echo ""
                echo "Benchmark Summary:"
                grep "Benchmark" benchmark.out | head -5
            fi
        else
            echo -e "${YELLOW}⚠️  Benchmarks failed or no benchmarks found${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Integration test directory not found, skipping benchmarks${NC}"
    fi
    echo ""
fi

# Race Detection
print_section "🏁 Race Detection"
echo -e "${BLUE}Running tests with race detection...${NC}"

race_failed=false
if [ -d "unit" ]; then
    if ! go test -race ./unit/... &> /dev/null; then
        echo -e "${RED}❌ Race conditions detected in unit tests${NC}"
        race_failed=true
    fi
fi

if [ -d "integration" ]; then
    if ! go test -race ./integration/... &> /dev/null; then
        echo -e "${RED}❌ Race conditions detected in integration tests${NC}"
        race_failed=true
    fi
fi

if [ "$race_failed" = "false" ]; then
    echo -e "${GREEN}✅ No race conditions detected${NC}"
fi
echo ""

# Memory Leak Detection (basic)
print_section "🧠 Memory Leak Detection"
echo -e "${BLUE}Running memory leak detection...${NC}"

if go test -memprofile=mem.prof ./... &> /dev/null; then
    echo -e "${GREEN}✅ Memory profile generated: mem.prof${NC}"
    echo "Use 'go tool pprof mem.prof' to analyze memory usage"
else
    echo -e "${YELLOW}⚠️  Memory profiling failed${NC}"
fi
echo ""

# Test Summary
print_section "📋 Test Summary"
echo -e "Total test suites: ${BLUE}$total_tests${NC}"
echo -e "Passed: ${GREEN}$passed_tests${NC}"
echo -e "Failed: ${RED}$failed_tests${NC}"

if [ $failed_tests -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed successfully!${NC}"
    echo ""
    echo "Generated files:"
    [ -f "coverage.out" ] && echo "  • coverage.out - Coverage data"
    [ -f "coverage.html" ] && echo "  • coverage.html - HTML coverage report"
    [ -f "benchmark.out" ] && echo "  • benchmark.out - Benchmark results"
    [ -f "mem.prof" ] && echo "  • mem.prof - Memory profile"
    echo ""
    echo "Next steps:"
    echo "  • Review coverage report: open coverage.html"
    echo "  • Analyze benchmarks: cat benchmark.out"
    echo "  • Profile memory: go tool pprof mem.prof"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Some tests failed. Please check the output above.${NC}"
    echo ""
    echo "Debugging tips:"
    echo "  • Run with verbose output: GRADLE_PARSER_TEST_VERBOSE=true ./run-tests.sh"
    echo "  • Run specific test: go test -v -run TestName ./unit/"
    echo "  • Check race conditions: go test -race ./..."
    exit 1
fi
