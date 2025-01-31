#!/bin/bash

# This script automates the installation of the mohua tool.
# It checks for the specified version (or fetches the latest one),
# downloads the appropriate binary, and installs it on the system.

# Check for required tools: curl, tar, unzip
# These tools are necessary for downloading and extracting the mohua binary.
if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed. Please install curl and try again."
    exit 1
fi

if ! command -v tar &> /dev/null; then
    echo "Error: tar is not installed. Please install tar and try again."
    exit 1
fi

if ! command -v unzip &> /dev/null; then
    echo "Error: unzip is not installed. Please install unzip and try again."
    exit 1
fi

# Determine the version of mohua to install.
# If no version is specified as a command line argument, fetch the latest version.
if [ -z "$1" ]; then
    VERSION=$(curl -s https://api.github.com/repos/watany-dev/mohua/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Failed to fetch the latest version."
        exit 1
    fi
else
    VERSION=$1
fi

# Remove any leading 'v' from the version string.
VERSION=${VERSION#v}

# Detect the architecture of the current system.
ARCH=$(uname -m)
case $ARCH in
    x86_64|amd64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    i386|i686)     ARCH="i386" ;;
    *) echo "Error: Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Detect the operating system of the current system.
OS=$(uname -s)
case $OS in
    Linux) OS="Linux" ;;
    Darwin) OS="Darwin" ;;
    MINGW*|MSYS*|CYGWIN*) OS="Windows" ;;
    *) echo "Error: Unsupported OS: $OS"; exit 1 ;;
esac

# Determine the file extension based on the operating system.
if [ "$OS" == "Windows" ]; then
    EXT="zip"
else
    EXT="tar.gz"
fi

# Construct the download URL for the mohua binary based on the version, OS, and architecture.
FILE_NAME="mohua_${OS}_${ARCH}.${EXT}"
URL="https://github.com/watany-dev/mohua/releases/download/v${VERSION}/${FILE_NAME}"

# Download the mohua binary.
echo "Downloading mohua from: $URL"
if ! curl -L -o "$FILE_NAME" "$URL"; then
    echo "Error: Failed to download mohua. URL: $URL"
    exit 1
fi

# Extract and install mohua.
echo "Installing mohua..."
if [ "$EXT" == "tar.gz" ]; then
    if ! tar -xzf "$FILE_NAME"; then
        echo "Error: Failed to extract mohua."
        exit 1
    fi
    if [ "$OS" != "Windows" ]; then
        if ! sudo mv mohua /usr/local/bin/mohua; then
            echo "Error: Failed to install mohua to /usr/local/bin."
            exit 1
        fi
    fi
elif [ "$EXT" == "zip" ]; then
    if ! unzip "$FILE_NAME"; then
        echo "Error: Failed to extract mohua."
        exit 1
    fi
    if [ "$OS" == "Windows" ]; then
        if ! mv mohua.exe /usr/local/bin/mohua.exe; then
            echo "Error: Failed to install mohua.exe to /usr/local/bin."
            exit 1
        fi
    fi
else
    echo "Error: Unknown file extension: $EXT"
    exit 1
fi

# Clean up by removing the downloaded file.
rm "$FILE_NAME"

echo "mohua installation complete."
echo "Run 'mohua --help' to see how to use mohua."
