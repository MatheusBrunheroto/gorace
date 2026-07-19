#!/bin/bash
set -e

BINARY_NAME="gorace"
INSTALL_DIR="/usr/local/bin"

echo "[*] Checking for Go installation..."
if ! command -v go &> /dev/null; then
    echo "[x] Go is not installed. Please install Go first: https://go.dev/dl/"
    exit 1
fi

if [ -d ".git" ]; then
    echo "[*] Existing git repo detected, pulling latest changes..."
    git pull
fi

echo "[*] Building $BINARY_NAME from source..."
go build -o "$BINARY_NAME" .

echo "[*] Installing to $INSTALL_DIR (requires sudo)..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "[+] $BINARY_NAME installed/updated successfully!"
echo "[+] Run '$BINARY_NAME -h' or '$BINARY_NAME --help' to get started."