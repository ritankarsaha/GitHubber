#!/bin/bash

# GitHubber Uninstallation Script
# Author: Ritankar Saha <ritankar.saha786@gmail.com>
# Description: Removes GitHubber from the system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="githubber"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.githubber"

# Print header
print_header() {
    echo -e "${BLUE}"
    echo "🗑️  GitHubber Uninstallation Script"
    echo "===================================="
    echo "Author: Ritankar Saha <ritankar.saha786@gmail.com>"
    echo -e "${NC}"
}

# Check if GitHubber is installed
check_installation() {
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        echo -e "${GREEN}✅ GitHubber found at ${INSTALL_DIR}/${BINARY_NAME}${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  GitHubber not found in ${INSTALL_DIR}${NC}"
        return 1
    fi
}

# Remove binary
remove_binary() {
    echo -e "${BLUE}🗑️  Removing GitHubber binary...${NC}"
    if sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}"; then
        echo -e "${GREEN}✅ Binary removed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to remove binary${NC}"
        return 1
    fi
}

# Ask about configuration removal
ask_config_removal() {
    if [ -d "${CONFIG_DIR}" ]; then
        echo -e "${YELLOW}📁 Configuration directory found: ${CONFIG_DIR}${NC}"
        echo -e "${YELLOW}This contains your GitHub tokens and preferences.${NC}"
        echo ""
        read -p "Do you want to remove the configuration directory? (y/N): " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            remove_config
        else
            echo -e "${BLUE}ℹ️  Configuration directory preserved${NC}"
        fi
    else
        echo -e "${BLUE}ℹ️  No configuration directory found${NC}"
    fi
}

# Remove configuration directory
remove_config() {
    echo -e "${BLUE}🗑️  Removing configuration directory...${NC}"
    if rm -rf "${CONFIG_DIR}"; then
        echo -e "${GREEN}✅ Configuration directory removed${NC}"
    else
        echo -e "${RED}❌ Failed to remove configuration directory${NC}"
    fi
}

# Print completion message
print_completion() {
    echo -e "${GREEN}"
    echo "🎉 GitHubber uninstallation completed!"
    echo ""
    echo "Thank you for using GitHubber!"
    echo "If you have any feedback, please visit:"
    echo "https://github.com/ritankarsaha/GitHubber/issues"
    echo ""
    echo "To reinstall, run:"
    echo "curl -sSL https://raw.githubusercontent.com/ritankarsaha/GitHubber/main/scripts/install.sh | bash"
    echo -e "${NC}"
}

# Handle forced uninstallation
force_uninstall() {
    echo -e "${YELLOW}🔥 Force uninstalling GitHubber...${NC}"
    
    # Remove binary (ignore errors)
    sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
    
    # Remove config (ignore errors)
    rm -rf "${CONFIG_DIR}" 2>/dev/null || true
    
    echo -e "${GREEN}✅ Force uninstallation complete${NC}"
}

# Main uninstallation function
main() {
    local force=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force|-f)
                force=true
                shift
                ;;
            --help|-h)
                echo "GitHubber Uninstall Script"
                echo ""
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --force, -f    Force uninstall (remove everything without prompts)"
                echo "  --help, -h     Show this help message"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ Unknown option: $1${NC}"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    print_header
    
    if [ "$force" = true ]; then
        force_uninstall
        print_completion
        exit 0
    fi
    
    # Check if GitHubber is installed
    if ! check_installation; then
        echo -e "${YELLOW}ℹ️  GitHubber doesn't appear to be installed${NC}"
        echo -e "${BLUE}Still want to clean up any remaining files?${NC}"
        read -p "Continue with cleanup? (y/N): " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}👋 Exiting...${NC}"
            exit 0
        fi
    fi
    
    # Confirm uninstallation
    echo -e "${YELLOW}⚠️  This will remove GitHubber from your system${NC}"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}👋 Uninstallation cancelled${NC}"
        exit 0
    fi
    
    # Remove binary
    remove_binary
    
    # Ask about configuration
    ask_config_removal
    
    print_completion
}

# Handle interruption
trap 'echo -e "${RED}\n❌ Uninstallation interrupted${NC}"; exit 1' INT TERM

# Run main function
main "$@"