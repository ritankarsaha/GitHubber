#!/bin/bash

# GitHubber Installation Script
# Author: Ritankar Saha <ritankar.saha786@gmail.com>
# Description: Cross-platform installation script for GitHubber

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="ritankarsaha/GitHubber"
BINARY_NAME="githubber"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.githubber"

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)   os="linux" ;;
        Darwin*)  os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *) 
            echo -e "${RED}❌ Unsupported operating system: $(uname -s)${NC}"
            exit 1
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) 
            echo -e "${RED}❌ Unsupported architecture: $(uname -m)${NC}"
            exit 1
            ;;
    esac
    
    PLATFORM="${os}-${arch}"
    if [ "$os" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    echo -e "${BLUE}🔍 Detected platform: ${PLATFORM}${NC}"
}

# Check if Go is installed for building from source
check_go() {
    if command -v go >/dev/null 2>&1; then
        local go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
        echo -e "${GREEN}✅ Go found: ${go_version}${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Go not found${NC}"
        return 1
    fi
}

# Check if git is installed
check_git() {
    if command -v git >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Git found${NC}"
        return 0
    else
        echo -e "${RED}❌ Git is required but not installed${NC}"
        echo -e "${YELLOW}Please install Git and try again${NC}"
        exit 1
    fi
}

# Install from GitHub releases (if available)
install_from_release() {
    echo -e "${BLUE}📦 Attempting to install from GitHub releases...${NC}"
    
    # Get latest release URL
    local download_url="https://github.com/${GITHUB_REPO}/releases/latest/download/${BINARY_NAME}-${PLATFORM}"
    if [ "$PLATFORM" = "windows-amd64" ]; then
        download_url="${download_url}.exe"
    fi
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    local binary_path="${temp_dir}/${BINARY_NAME}"
    
    # Download binary
    echo -e "${BLUE}⬇️  Downloading from: ${download_url}${NC}"
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${binary_path}" "${download_url}" || {
            echo -e "${YELLOW}⚠️  Failed to download from releases${NC}"
            return 1
        }
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${binary_path}" "${download_url}" || {
            echo -e "${YELLOW}⚠️  Failed to download from releases${NC}"
            return 1
        }
    else
        echo -e "${RED}❌ Neither curl nor wget found${NC}"
        return 1
    fi
    
    # Make executable and install
    chmod +x "${binary_path}"
    sudo mv "${binary_path}" "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Cleanup
    rm -rf "${temp_dir}"
    
    echo -e "${GREEN}✅ GitHubber installed from release${NC}"
    return 0
}

# Build and install from source
install_from_source() {
    echo -e "${BLUE}🔨 Building from source...${NC}"
    
    # Check if Go is available
    if ! check_go; then
        echo -e "${RED}❌ Go is required to build from source${NC}"
        echo -e "${YELLOW}Please install Go 1.21+ or use a pre-built release${NC}"
        exit 1
    fi
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "${temp_dir}"
    
    # Clone repository
    echo -e "${BLUE}📥 Cloning repository...${NC}"
    git clone "https://github.com/${GITHUB_REPO}.git" GitHubber
    cd GitHubber
    
    # Build
    echo -e "${BLUE}🔨 Building binary...${NC}"
    go build -o "${BINARY_NAME}" ./cmd/main.go
    
    # Install
    sudo mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Cleanup
    cd /
    rm -rf "${temp_dir}"
    
    echo -e "${GREEN}✅ GitHubber built and installed from source${NC}"
}

# Create configuration directory
setup_config() {
    if [ ! -d "${CONFIG_DIR}" ]; then
        echo -e "${BLUE}📁 Creating configuration directory: ${CONFIG_DIR}${NC}"
        mkdir -p "${CONFIG_DIR}"
    fi
}

# Print installation completion message
print_completion() {
    echo -e "${GREEN}"
    echo "🎉 GitHubber installation completed successfully!"
    echo ""
    echo "📋 Next steps:"
    echo "  1. Run 'githubber' to start the application"
    echo "  2. Set up GitHub authentication (optional):"
    echo "     - Export GITHUB_TOKEN environment variable, or"
    echo "     - Use the Settings menu in GitHubber"
    echo ""
    echo "📚 Documentation: https://github.com/${GITHUB_REPO}/wiki"
    echo "🐛 Report issues: https://github.com/${GITHUB_REPO}/issues"
    echo ""
    echo "Created by Ritankar Saha <ritankar.saha786@gmail.com>"
    echo -e "${NC}"
}

# Main installation function
main() {
    echo -e "${BLUE}"
    echo "🚀 GitHubber Installation Script"
    echo "=================================="
    echo "Author: Ritankar Saha <ritankar.saha786@gmail.com>"
    echo -e "${NC}"
    
    # Check prerequisites
    check_git
    detect_platform
    
    # Create config directory
    setup_config
    
    # Try installation methods
    echo -e "${BLUE}🔍 Choosing installation method...${NC}"
    
    # Method 1: Try installing from GitHub releases
    if install_from_release; then
        print_completion
        exit 0
    fi
    
    # Method 2: Build from source
    echo -e "${YELLOW}⚠️  Release installation failed, trying source build...${NC}"
    install_from_source
    print_completion
}

# Handle interruption
trap 'echo -e "${RED}\n❌ Installation interrupted${NC}"; exit 1' INT TERM

# Check if running as root (not recommended)
if [ "$EUID" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Warning: Running as root is not recommended${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Run main function
main "$@"