# Test Suite

This directory contains comprehensive tests for the Gradle Parser library, including unit tests, integration tests, and performance benchmarks.

## 📁 Test Structure

```
test/
├── README.md                   # This file
├── unit/                       # Unit tests
│   ├── parser_test.go          # Parser unit tests
│   ├── api_test.go             # API unit tests
│   ├── models_test.go          # Model validation tests
│   └── editor_test.go          # Editor unit tests
├── integration/                # Integration tests
│   ├── gradle_files/           # Test Gradle files
│   ├── parsing_test.go         # End-to-end parsing tests
│   ├── editing_test.go         # End-to-end editing tests
│   └── performance_test.go     # Performance benchmarks
├── fixtures/                   # Test data and fixtures
│   ├── simple/                 # Simple Gradle files
│   ├── complex/                # Complex multi-module projects
│   ├── android/                # Android project samples
│   ├── spring-boot/            # Spring Boot project samples
│   └── kotlin/                 # Kotlin DSL samples
└── scripts/                    # Test automation scripts
    ├── run-tests.sh            # Run all tests
    ├── benchmark.sh            # Performance benchmarks
    └── coverage.sh             # Generate coverage reports
```

## 🚀 Running Tests

### Prerequisites

- Go 1.19 or later
- Git (for test data)

### Quick Start

```bash
# Run all tests
./scripts/run-tests.sh

# Run unit tests only
go test ./unit/...

# Run integration tests only
go test ./integration/...

# Run with coverage
./scripts/coverage.sh

# Run performance benchmarks
./scripts/benchmark.sh
```

### Individual Test Categories

#### Unit Tests

```bash
cd test/unit
go test -v
```

Unit tests cover:
- Parser functionality
- API methods
- Data model validation
- Editor operations
- Error handling

#### Integration Tests

```bash
cd test/integration
go test -v
```

Integration tests cover:
- End-to-end parsing workflows
- Real Gradle file processing
- Multi-module project handling
- Performance characteristics

#### Performance Benchmarks

```bash
cd test/integration
go test -bench=. -benchmem
```

Benchmarks measure:
- Parsing speed for different file sizes
- Memory usage optimization
- Configuration impact on performance
- Batch processing efficiency

## 📊 Test Coverage

The test suite aims for comprehensive coverage:

- **Parser Core**: 95%+ coverage
- **API Layer**: 90%+ coverage
- **Data Models**: 100% coverage
- **Editor Functions**: 90%+ coverage

### Generating Coverage Reports

```bash
# Generate HTML coverage report
./scripts/coverage.sh

# View coverage in browser
open coverage.html
```

## 🧪 Test Categories

### 1. Parser Tests

Test the core parsing functionality:

```go
func TestBasicParsing(t *testing.T) {
    content := `
    plugins {
        id 'java'
    }
    group = 'com.example'
    version = '1.0.0'
    `
    
    result, err := api.ParseString(content)
    assert.NoError(t, err)
    assert.Equal(t, "com.example", result.Project.Group)
    assert.Equal(t, "1.0.0", result.Project.Version)
}
```

### 2. API Tests

Test high-level API functions:

```go
func TestGetDependencies(t *testing.T) {
    deps, err := api.GetDependencies("../fixtures/simple/build.gradle")
    assert.NoError(t, err)
    assert.Greater(t, len(deps), 0)
}
```

### 3. Model Tests

Test data structure validation:

```go
func TestProjectModel(t *testing.T) {
    project := &model.Project{
        Group:   "com.example",
        Name:    "test-project",
        Version: "1.0.0",
    }
    
    assert.True(t, project.IsValid())
}
```

### 4. Editor Tests

Test structured editing capabilities:

```go
func TestUpdateDependencyVersion(t *testing.T) {
    original := `implementation 'mysql:mysql-connector-java:8.0.28'`
    expected := `implementation 'mysql:mysql-connector-java:8.0.31'`
    
    result, err := api.UpdateDependencyVersion(
        "test.gradle", "mysql", "mysql-connector-java", "8.0.31")
    assert.NoError(t, err)
    assert.Contains(t, result, expected)
}
```

### 5. Integration Tests

Test complete workflows:

```go
func TestCompleteWorkflow(t *testing.T) {
    // Parse project
    result, err := api.ParseFile("../fixtures/spring-boot/build.gradle")
    assert.NoError(t, err)
    
    // Verify project type
    assert.True(t, api.IsSpringBootProject(result.Project.Plugins))
    
    // Update dependency
    newContent, err := api.UpdateDependencyVersion(
        "../fixtures/spring-boot/build.gradle",
        "org.springframework.boot", "spring-boot-starter-web", "2.7.2")
    assert.NoError(t, err)
    
    // Verify update
    assert.Contains(t, newContent, "2.7.2")
}
```

## 🎯 Test Data

### Fixture Files

The test suite includes various Gradle file fixtures:

#### Simple Projects
- Basic Java project
- Minimal configuration
- Single module

#### Complex Projects
- Multi-module setup
- Custom configurations
- Advanced plugin usage

#### Framework-Specific
- Android applications
- Spring Boot projects
- Kotlin multiplatform

#### Edge Cases
- Empty files
- Malformed syntax
- Large files (performance testing)

### Test Data Management

```bash
# Update test fixtures
git submodule update --init --recursive

# Add new test case
cp your-gradle-file.gradle test/fixtures/custom/
```

## 🔧 Test Configuration

### Environment Variables

```bash
# Enable verbose testing
export GRADLE_PARSER_TEST_VERBOSE=true

# Set test timeout
export GRADLE_PARSER_TEST_TIMEOUT=30s

# Enable performance profiling
export GRADLE_PARSER_TEST_PROFILE=true
```

### Custom Test Settings

Create `test/config.json`:

```json
{
  "timeout": "30s",
  "verbose": true,
  "coverage": true,
  "benchmarks": true,
  "fixtures": {
    "simple": "fixtures/simple",
    "complex": "fixtures/complex"
  }
}
```

## 📈 Performance Testing

### Benchmark Categories

1. **Parsing Speed**
   - Small files (< 1KB)
   - Medium files (1-10KB)
   - Large files (> 10KB)

2. **Memory Usage**
   - Standard parsing
   - Optimized parsing
   - Batch processing

3. **Configuration Impact**
   - Different parser options
   - Feature toggles
   - Memory vs. speed tradeoffs

### Running Benchmarks

```bash
# All benchmarks
go test -bench=. -benchmem ./integration/

# Specific benchmark
go test -bench=BenchmarkParsing -benchmem ./integration/

# With CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./integration/

# With memory profiling
go test -bench=. -memprofile=mem.prof ./integration/
```

## 🐛 Debugging Tests

### Verbose Output

```bash
go test -v ./...
```

### Test-Specific Debugging

```bash
# Run single test
go test -run TestSpecificFunction ./unit/

# With race detection
go test -race ./...

# With coverage
go test -cover ./...
```

### Debugging Failed Tests

```bash
# Show detailed failure information
go test -v -failfast ./...

# Generate test binary for debugging
go test -c ./unit/
./unit.test -test.run TestSpecificFunction
```

## 🤝 Contributing Tests

### Adding New Tests

1. **Choose appropriate category** (unit vs integration)
2. **Follow naming conventions** (`TestFunctionName`)
3. **Include edge cases** and error scenarios
4. **Add documentation** for complex test cases
5. **Update fixtures** if needed

### Test Guidelines

1. **Test one thing at a time**
2. **Use descriptive test names**
3. **Include both positive and negative cases**
4. **Mock external dependencies**
5. **Keep tests fast and reliable**

### Example Test Structure

```go
func TestParseComplexGradleFile(t *testing.T) {
    // Arrange
    testFile := "../fixtures/complex/build.gradle"
    
    // Act
    result, err := api.ParseFile(testFile)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result.Project)
    assert.Greater(t, len(result.Project.Dependencies), 0)
    
    // Verify specific expectations
    assert.Equal(t, "com.example", result.Project.Group)
    assert.True(t, api.IsSpringBootProject(result.Project.Plugins))
}
```

## 📚 Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Assertion Library](https://github.com/stretchr/testify)
- [Go Benchmarking](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [Test Coverage](https://golang.org/doc/tutorial/add-a-test)
