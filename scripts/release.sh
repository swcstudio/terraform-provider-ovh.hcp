#!/bin/bash

# Release script for terraform-provider-hashicorp-ovh
# This script automates the release process using GoReleaser

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_requirements() {
    print_status "Checking requirements..."
    
    if ! command -v goreleaser &> /dev/null; then
        print_error "GoReleaser is not installed. Install it with: brew install goreleaser"
        exit 1
    fi
    
    if ! command -v git &> /dev/null; then
        print_error "Git is not installed"
        exit 1
    fi
    
    if ! command -v gpg &> /dev/null; then
        print_error "GPG is not installed"
        exit 1
    fi
    
    print_status "All requirements satisfied"
}

# Set up environment variables
setup_environment() {
    print_status "Setting up environment variables..."
    
    export GPG_TTY=$(tty)
    export GPG_FINGERPRINT=4B2412305828FBD7
    
    if [ -z "$GITHUB_TOKEN" ]; then
        print_error "GITHUB_TOKEN environment variable is not set"
        print_error "Please set it with: export GITHUB_TOKEN=your_github_token"
        exit 1
    fi
    
    print_status "Environment variables configured"
}

# Verify Git state
check_git_state() {
    print_status "Checking Git state..."
    
    if [ -n "$(git status --porcelain)" ]; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        exit 1
    fi
    
    # Check if we're on main branch
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        print_warning "You're not on the main branch (current: $current_branch). Continue? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_error "Aborted by user"
            exit 1
        fi
    fi
    
    print_status "Git state is clean"
}

# Verify GPG key
check_gpg_key() {
    print_status "Verifying GPG key..."
    
    if ! gpg --list-secret-keys $GPG_FINGERPRINT &> /dev/null; then
        print_error "GPG key $GPG_FINGERPRINT not found in keyring"
        exit 1
    fi
    
    print_status "GPG key verified"
}

# Get version from user
get_version() {
    if [ -z "$1" ]; then
        print_error "Usage: $0 <version>"
        print_error "Example: $0 v0.1.0"
        exit 1
    fi
    
    VERSION=$1
    
    # Ensure version starts with 'v'
    if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
        print_error "Version must be in format v0.0.0 (e.g., v0.1.0, v1.2.3-beta1)"
        exit 1
    fi
    
    print_status "Version: $VERSION"
}

# Create and push tag
create_tag() {
    print_status "Creating and pushing tag $VERSION..."
    
    if git tag -l | grep -q "^$VERSION$"; then
        print_error "Tag $VERSION already exists"
        exit 1
    fi
    
    git tag $VERSION
    git push origin $VERSION
    
    print_status "Tag $VERSION created and pushed"
}

# Run goreleaser
run_release() {
    print_status "Running GoReleaser..."
    
    if ! goreleaser release --clean; then
        print_error "GoReleaser failed"
        exit 1
    fi
    
    print_status "Release completed successfully!"
}

# Main function
main() {
    print_status "Starting release process for terraform-provider-hashicorp-ovh"
    
    get_version $1
    check_requirements
    setup_environment
    check_git_state
    check_gpg_key
    
    print_status "All checks passed. Creating release for $VERSION..."
    print_warning "This will create a tag and publish a release. Continue? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_error "Aborted by user"
        exit 1
    fi
    
    create_tag
    run_release
    
    print_status "🎉 Release $VERSION completed successfully!"
    print_status "Check the release at: https://github.com/swcstudio/terraform-provider-hashicorp-ovh/releases/tag/$VERSION"
}

# Run main function with all arguments
main "$@"