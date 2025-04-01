#!/bin/sh

set -e  # Stop execution on error

# Check if Go is already installed
if command -v go >/dev/null 2>&1; then
    echo "Go is already installed: $(go version)"
    sleep 0.5
else
    echo "Go is not installed. Proceeding with installation..."

    # Detect package manager
    if command -v apt >/dev/null 2>&1; then
        PKG_MANAGER="apt"
    elif command -v pacman >/dev/null 2>&1; then
        PKG_MANAGER="pacman"
    elif command -v dnf >/dev/null 2>&1; then
        PKG_MANAGER="dnf"
    elif command -v yum >/dev/null 2>&1; then
        PKG_MANAGER="yum"
    elif command -v zypper >/dev/null 2>&1; then
        PKG_MANAGER="zypper"
    elif command -v brew >/dev/null 2>&1; then
        PKG_MANAGER="brew"
    else
        echo "Could not detect a package manager. Please install Go manually."
        PKG_MANAGER=""
    fi

    # Install Go if a package manager is found
    if [ -n "$PKG_MANAGER" ]; then
        case "$PKG_MANAGER" in
            apt)
                sudo apt update && sudo apt install -y golang
                ;;
            pacman)
                sudo pacman -Sy --noconfirm go
                ;;
            dnf)
                sudo dnf install -y golang
                ;;
            yum)
                sudo yum install -y golang
                ;;
            zypper)
                sudo zypper install -y go
                ;;
            brew)
                brew install go
                ;;
        esac
        echo "Go has been successfully installed!"
    fi
fi

# Print Go version if installed
if command -v go >/dev/null 2>&1; then
    # echo "Current Go version: $(go version)"
    echo "Checking golang"
else
    echo "Go is still not installed. Please check for issues."
fi

sleep 0.5

echo "Building fetchify from sources"

go build -o fetchify

sleep 0.5

echo "Installing fetchify in /usr/bin/"
sleep 0.5
sudo cp fetchify /usr/bin
echo "Fetchify successfully installed. Use: fetchify --help"