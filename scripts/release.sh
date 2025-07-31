#!/bin/bash

# Gradle Parser Release Script
# 自动化发布脚本

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/scagogogo/gradle-parser"
MAIN_BRANCH="main"

# Functions
print_header() {
    echo -e "${PURPLE}================================${NC}"
    echo -e "${PURPLE}  Gradle Parser Release Tool${NC}"
    echo -e "${PURPLE}================================${NC}"
    echo ""
}

print_step() {
    echo -e "${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if we're in the right directory
check_directory() {
    if [ ! -f "go.mod" ] || [ ! -f "README.md" ]; then
        print_error "Please run this script from the root of the gradle-parser repository"
        exit 1
    fi
}

# Check if git is clean
check_git_status() {
    print_step "Checking git status..."
    
    if [ -n "$(git status --porcelain)" ]; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        git status --short
        exit 1
    fi
    
    print_success "Working directory is clean"
}

# Check if we're on the main branch
check_branch() {
    print_step "Checking current branch..."
    
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "$MAIN_BRANCH" ]; then
        print_error "Please switch to the $MAIN_BRANCH branch before releasing"
        print_error "Current branch: $current_branch"
        exit 1
    fi
    
    print_success "On $MAIN_BRANCH branch"
}

# Pull latest changes
pull_latest() {
    print_step "Pulling latest changes..."
    
    git pull origin "$MAIN_BRANCH"
    print_success "Latest changes pulled"
}

# Run tests
run_tests() {
    print_step "Running tests..."
    
    # Run Go tests
    if ! go test ./...; then
        print_error "Tests failed"
        exit 1
    fi
    
    # Run example tests
    if [ -f "examples/run-all-examples.sh" ]; then
        cd examples
        if ! ./run-all-examples.sh; then
            print_error "Example tests failed"
            exit 1
        fi
        cd ..
    fi
    
    # Run comprehensive test suite
    if [ -f "test/scripts/run-tests.sh" ]; then
        cd test
        if ! ./scripts/run-tests.sh; then
            print_error "Comprehensive test suite failed"
            exit 1
        fi
        cd ..
    fi
    
    print_success "All tests passed"
}

# Get current version
get_current_version() {
    if git tag --list | grep -q "v"; then
        git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
    else
        echo "v0.0.0"
    fi
}

# Validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
        print_error "Invalid version format. Use semantic versioning (e.g., v1.0.0, v1.0.0-beta.1)"
        return 1
    fi
    return 0
}

# Check if version already exists
check_version_exists() {
    local version=$1
    if git tag --list | grep -q "^$version$"; then
        print_error "Version $version already exists"
        return 1
    fi
    return 0
}

# Update changelog
update_changelog() {
    local version=$1
    local date=$(date +%Y-%m-%d)
    
    print_step "Updating CHANGELOG.md..."
    
    # Create backup
    cp CHANGELOG.md CHANGELOG.md.bak
    
    # Update changelog
    sed -i.tmp "s/## \[Unreleased\]/## [Unreleased]\n\n## [$version] - $date/" CHANGELOG.md
    rm -f CHANGELOG.md.tmp
    
    print_success "CHANGELOG.md updated"
}

# Create git tag
create_tag() {
    local version=$1
    local message=$2
    
    print_step "Creating git tag $version..."
    
    git add CHANGELOG.md
    git commit -m "chore: prepare release $version"
    git tag -a "$version" -m "$message"
    
    print_success "Tag $version created"
}

# Push changes
push_changes() {
    local version=$1
    
    print_step "Pushing changes and tags..."
    
    git push origin "$MAIN_BRANCH"
    git push origin "$version"
    
    print_success "Changes and tags pushed"
}

# Generate release notes
generate_release_notes() {
    local version=$1
    local prev_version=$2
    
    print_step "Generating release notes..."
    
    echo "# Release Notes for $version"
    echo ""
    echo "## Changes"
    echo ""
    
    if [ "$prev_version" != "v0.0.0" ]; then
        git log --pretty=format:"- %s (%h)" "$prev_version..$version" | head -20
    else
        echo "- Initial release"
    fi
    
    echo ""
    echo "## Installation"
    echo ""
    echo "\`\`\`bash"
    echo "go get github.com/scagogogo/gradle-parser@$version"
    echo "\`\`\`"
    echo ""
    echo "## Documentation"
    echo ""
    echo "- [API Documentation](https://scagogogo.github.io/gradle-parser/api/)"
    echo "- [User Guide](https://scagogogo.github.io/gradle-parser/guide/)"
    echo "- [Examples](https://github.com/scagogogo/gradle-parser/tree/$version/examples)"
}

# Main release function
release() {
    local version=$1
    local message=${2:-"Release $version"}
    
    print_header
    
    # Validations
    check_directory
    check_git_status
    check_branch
    
    # Get current version
    local current_version=$(get_current_version)
    print_step "Current version: $current_version"
    print_step "New version: $version"
    
    # Validate new version
    if ! validate_version "$version"; then
        exit 1
    fi
    
    if ! check_version_exists "$version"; then
        exit 1
    fi
    
    # Confirm release
    echo ""
    echo -e "${YELLOW}Are you sure you want to release version $version? (y/N)${NC}"
    read -r confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        print_warning "Release cancelled"
        exit 0
    fi
    
    # Pull latest changes
    pull_latest
    
    # Run tests
    run_tests
    
    # Update changelog
    update_changelog "$version"
    
    # Create tag
    create_tag "$version" "$message"
    
    # Push changes
    push_changes "$version"
    
    # Generate release notes
    echo ""
    print_step "Release notes:"
    echo ""
    generate_release_notes "$version" "$current_version"
    
    echo ""
    print_success "Release $version completed successfully!"
    echo ""
    echo -e "${GREEN}🎉 Release $version is now available!${NC}"
    echo -e "${BLUE}📦 GitHub Release: $REPO_URL/releases/tag/$version${NC}"
    echo -e "${BLUE}📚 Documentation: https://scagogogo.github.io/gradle-parser/${NC}"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Check the GitHub Actions workflow"
    echo "2. Verify the release assets"
    echo "3. Update any dependent projects"
    echo "4. Announce the release"
}

# Show usage
usage() {
    echo "Usage: $0 <version> [message]"
    echo ""
    echo "Examples:"
    echo "  $0 v1.0.0"
    echo "  $0 v1.0.1 'Bug fix release'"
    echo "  $0 v1.1.0-beta.1 'Beta release'"
    echo ""
    echo "Version format: vMAJOR.MINOR.PATCH[-PRERELEASE]"
}

# Main script
main() {
    if [ $# -eq 0 ]; then
        usage
        exit 1
    fi
    
    local version=$1
    local message=${2:-"Release $version"}
    
    release "$version" "$message"
}

# Run main function with all arguments
main "$@"
