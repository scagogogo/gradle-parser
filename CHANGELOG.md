# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation with VitePress
- Multi-language support (English and Chinese)
- Advanced example demonstrating performance optimization
- Comprehensive test suite with unit and integration tests
- CI/CD pipeline with GitHub Actions
- Code quality checks and security scanning
- Performance benchmarking
- Structured editing capabilities
- Source mapping support (planned)

### Changed
- Improved API design for better usability
- Enhanced error handling and validation
- Better project type detection
- Optimized parsing performance

### Fixed
- Various parsing edge cases
- Memory usage optimization
- Dependency scope detection

## [0.1.0] - 2024-01-XX (Initial Release)

### Added
- Core Gradle file parsing functionality
- Support for parsing dependencies, plugins, and repositories
- Project type detection (Android, Kotlin, Spring Boot)
- Basic structured editing capabilities
- Command-line examples and tools
- Go module support
- MIT license

### Features
- **Parsing**: Parse Gradle build files (Groovy DSL)
- **Dependencies**: Extract and analyze project dependencies
- **Plugins**: Detect and analyze applied plugins
- **Repositories**: Parse repository configurations
- **Project Types**: Automatic detection of project types
- **Editing**: Basic structured editing with format preservation
- **Performance**: Optimized for speed and memory usage

### API
- `api.ParseFile()` - Parse Gradle file from filesystem
- `api.ParseString()` - Parse Gradle content from string
- `api.ParseReader()` - Parse Gradle content from io.Reader
- `api.GetDependencies()` - Extract dependencies only
- `api.GetPlugins()` - Extract plugins only
- `api.GetRepositories()` - Extract repositories only
- `api.IsAndroidProject()` - Detect Android projects
- `api.IsKotlinProject()` - Detect Kotlin projects
- `api.IsSpringBootProject()` - Detect Spring Boot projects
- `api.DependenciesByScope()` - Group dependencies by scope
- `api.UpdateDependencyVersion()` - Update dependency versions
- `api.CreateGradleEditor()` - Create structured editor

### Examples
- Basic parsing example
- Dependency analysis example
- Plugin detection example
- Repository parsing example
- Complete feature demonstration
- Structured editing example
- Advanced features example

### Documentation
- Comprehensive API documentation
- User guides and tutorials
- Code examples and best practices
- Multi-language documentation (EN/CN)
- Interactive documentation website

### Testing
- Unit tests for core functionality
- Integration tests with real Gradle files
- Performance benchmarks
- Example validation tests
- Continuous integration testing

### Infrastructure
- GitHub Actions CI/CD pipeline
- Automated testing and quality checks
- Documentation deployment
- Release automation
- Code coverage reporting

## Version History

### Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version when making incompatible API changes
- **MINOR** version when adding functionality in a backwards compatible manner
- **PATCH** version when making backwards compatible bug fixes

### Release Schedule

- **Major releases**: Every 6-12 months
- **Minor releases**: Every 1-3 months
- **Patch releases**: As needed for bug fixes

### Supported Go Versions

- Go 1.19+
- Go 1.20+
- Go 1.21+ (recommended)

### Compatibility

- **Gradle**: 4.0+ (Groovy DSL)
- **Kotlin DSL**: Basic support (planned for v0.2.0)
- **Android Gradle Plugin**: 3.0+
- **Spring Boot**: 2.0+

## Migration Guide

### From v0.0.x to v0.1.0

This is the initial stable release. No migration needed.

### Future Migrations

Migration guides will be provided for breaking changes in major versions.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for information on how to contribute to this project.

## Support

- **Documentation**: https://scagogogo.github.io/gradle-parser/
- **Issues**: https://github.com/scagogogo/gradle-parser/issues
- **Discussions**: https://github.com/scagogogo/gradle-parser/discussions

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
