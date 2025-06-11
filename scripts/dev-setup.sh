#!/bin/bash

# Development setup script for terraform-provider-hashicorp-ovh
# This script sets up a complete development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_header() {
    echo -e "${BLUE}[SETUP]${NC} $1"
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        print_error "Unsupported OS: $OSTYPE"
        exit 1
    fi
    print_status "Detected OS: $OS"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Homebrew on macOS
install_homebrew() {
    if [[ "$OS" == "macos" ]] && ! command_exists brew; then
        print_header "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
}

# Install Go
install_go() {
    if ! command_exists go; then
        print_header "Installing Go..."
        if [[ "$OS" == "macos" ]]; then
            brew install go
        elif [[ "$OS" == "linux" ]]; then
            # Install Go on Linux
            GO_VERSION="1.21.0"
            wget -c "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O - | sudo tar -xz -C /usr/local
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            export PATH=$PATH:/usr/local/go/bin
        fi
    else
        print_status "Go is already installed: $(go version)"
    fi
}

# Install GoReleaser
install_goreleaser() {
    if ! command_exists goreleaser; then
        print_header "Installing GoReleaser..."
        if [[ "$OS" == "macos" ]]; then
            brew install goreleaser
        elif [[ "$OS" == "linux" ]]; then
            # Install GoReleaser on Linux
            curl -sfL https://goreleaser.com/static/run | bash -s -- --version
            sudo mv ./bin/goreleaser /usr/local/bin/
        fi
    else
        print_status "GoReleaser is already installed: $(goreleaser --version | head -n1)"
    fi
}

# Install tfplugindocs
install_tfplugindocs() {
    if ! command_exists tfplugindocs; then
        print_header "Installing tfplugindocs..."
        go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
    else
        print_status "tfplugindocs is already installed"
    fi
}

# Install golangci-lint
install_golangci_lint() {
    if ! command_exists golangci-lint; then
        print_header "Installing golangci-lint..."
        if [[ "$OS" == "macos" ]]; then
            brew install golangci-lint
        elif [[ "$OS" == "linux" ]]; then
            curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
        fi
    else
        print_status "golangci-lint is already installed: $(golangci-lint --version)"
    fi
}

# Install GPG
install_gpg() {
    if ! command_exists gpg; then
        print_header "Installing GPG..."
        if [[ "$OS" == "macos" ]]; then
            brew install gnupg
        elif [[ "$OS" == "linux" ]]; then
            sudo apt-get update && sudo apt-get install -y gnupg
        fi
    else
        print_status "GPG is already installed: $(gpg --version | head -n1)"
    fi
}

# Install pinentry-mac on macOS
install_pinentry() {
    if [[ "$OS" == "macos" ]] && ! command_exists pinentry-mac; then
        print_header "Installing pinentry-mac..."
        brew install pinentry-mac
        
        # Configure GPG to use pinentry-mac
        mkdir -p ~/.gnupg
        echo "pinentry-program /opt/homebrew/bin/pinentry-mac" >> ~/.gnupg/gpg-agent.conf
        gpg-connect-agent reloadagent /bye || true
    fi
}

# Setup Go environment
setup_go_env() {
    print_header "Setting up Go environment..."
    
    # Ensure GOPATH is set
    if [ -z "$GOPATH" ]; then
        export GOPATH=$HOME/go
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    fi
    
    # Create Go workspace directories
    mkdir -p $GOPATH/{bin,src,pkg}
    
    print_status "Go environment configured"
}

# Install project dependencies
install_dependencies() {
    print_header "Installing project dependencies..."
    
    go mod download
    go mod tidy
    
    print_status "Dependencies installed"
}

# Setup local development configuration
setup_local_config() {
    print_header "Setting up local development configuration..."
    
    # Create .env.example file for environment variables
    cat > .env.example << 'EOF'
# OVH API Configuration
OVH_ENDPOINT=ovh-eu
OVH_APPLICATION_KEY=your-app-key
OVH_APPLICATION_SECRET=your-app-secret
OVH_CONSUMER_KEY=your-consumer-key
OVH_PROJECT_ID=your-project-id

# HashiCorp Cloud Platform
HCP_CLIENT_ID=your-hcp-client-id
HCP_CLIENT_SECRET=your-hcp-client-secret

# Development Configuration
TF_LOG=DEBUG
TF_LOG_PATH=./terraform.log

# Release Configuration (for maintainers)
GITHUB_TOKEN=your-github-token
GPG_FINGERPRINT=your-gpg-fingerprint
EOF

    # Create local Terraform configuration for testing
    mkdir -p examples/local-dev
    cat > examples/local-dev/main.tf << 'EOF'
terraform {
  required_providers {
    hashicorp-ovh = {
      source = "swcstudio/hashicorp-ovh"
      version = "~> 0.1.0"
    }
  }
}

provider "hashicorp-ovh" {
  # Configuration will be loaded from environment variables
}

# Example resource - uncomment and modify as needed
# resource "hashicorp_ovh_nomad_cluster" "example" {
#   name         = "dev-cluster"
#   region       = "eu-west-1"
#   server_count = 1
#   client_count = 2
# }
EOF

    cat > examples/local-dev/terraform.tf << 'EOF'
# Terraform configuration for local development
terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}
EOF

    print_status "Local configuration created"
}

# Setup Git hooks
setup_git_hooks() {
    print_header "Setting up Git hooks..."
    
    mkdir -p .git/hooks
    
    # Pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for terraform-provider-hashicorp-ovh

set -e

echo "Running pre-commit checks..."

# Format code
echo "Formatting Go code..."
gofmt -s -w .

# Run linter
echo "Running linter..."
golangci-lint run

# Run tests
echo "Running tests..."
go test -v ./...

echo "Pre-commit checks passed!"
EOF

    chmod +x .git/hooks/pre-commit
    
    print_status "Git hooks configured"
}

# Run initial tests
run_tests() {
    print_header "Running initial tests..."
    
    # Format code
    make fmt
    
    # Run linter
    if command_exists golangci-lint; then
        make lint
    else
        print_warning "Skipping linter (golangci-lint not installed)"
    fi
    
    # Run unit tests
    make test
    
    print_status "Tests completed successfully"
}

# Create development shortcuts
create_shortcuts() {
    print_header "Creating development shortcuts..."
    
    # Create a simple development script
    cat > dev.sh << 'EOF'
#!/bin/bash
# Quick development commands for terraform-provider-hashicorp-ovh

case "$1" in
    "build")
        echo "Building provider..."
        make build
        ;;
    "test")
        echo "Running tests..."
        make test
        ;;
    "install")
        echo "Installing provider locally..."
        make install
        ;;
    "docs")
        echo "Generating documentation..."
        make docs
        ;;
    "release")
        if [ -z "$2" ]; then
            echo "Usage: $0 release <version>"
            echo "Example: $0 release v0.1.0"
            exit 1
        fi
        echo "Creating release $2..."
        ./scripts/release.sh $2
        ;;
    *)
        echo "Available commands:"
        echo "  build    - Build the provider"
        echo "  test     - Run tests"
        echo "  install  - Install provider locally"
        echo "  docs     - Generate documentation"
        echo "  release  - Create a new release"
        ;;
esac
EOF

    chmod +x dev.sh
    
    print_status "Development shortcuts created"
}

# Print setup completion message
print_completion() {
    print_header "🎉 Development environment setup completed!"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. Copy .env.example to .env and configure your API keys:"
    echo "   cp .env.example .env"
    echo ""
    echo "2. Build the provider:"
    echo "   make build"
    echo ""
    echo "3. Run tests:"
    echo "   make test"
    echo ""
    echo "4. Install locally for testing:"
    echo "   make install"
    echo ""
    echo "5. Use the development shortcuts:"
    echo "   ./dev.sh build"
    echo "   ./dev.sh test"
    echo "   ./dev.sh install"
    echo ""
    echo -e "${BLUE}Available make targets:${NC}"
    make help
    echo ""
    echo -e "${BLUE}Project structure:${NC}"
    echo "├── internal/           - Provider implementation"
    echo "├── examples/           - Example configurations"
    echo "├── docs/              - Documentation"
    echo "├── scripts/           - Development scripts"
    echo "└── charts/            - Helm charts"
    echo ""
    print_status "Happy coding! 🚀"
}

# Main setup function
main() {
    print_header "🛠️  Setting up development environment for terraform-provider-hashicorp-ovh"
    echo ""
    
    detect_os
    install_homebrew
    install_go
    install_goreleaser
    install_tfplugindocs
    install_golangci_lint
    install_gpg
    install_pinentry
    setup_go_env
    install_dependencies
    setup_local_config
    setup_git_hooks
    run_tests
    create_shortcuts
    print_completion
}

# Run main function
main "$@"