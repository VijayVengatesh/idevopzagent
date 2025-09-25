#!/bin/bash

# ==============================
# Build script for iDevopzAgent
# ==============================

# Project root
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd $ROOT_DIR || exit

# Output directory
BIN_DIR="$ROOT_DIR/bin"
mkdir -p "$BIN_DIR"

# Application name
APP_NAME="metrics-agent"

# Architectures to build
LINUX_ARCHS=("386" "amd64" "arm" "arm64")
WINDOWS_ARCHS=("386" "amd64")
MACOS_ARCHS=("amd64" "arm64")

# Function to build for Linux
build_linux() {
  for ARCH in "${LINUX_ARCHS[@]}"; do
    OUTPUT="$BIN_DIR/$APP_NAME-linux-$ARCH"
    echo "Building Linux $ARCH..."
    GOOS=linux GOARCH=$ARCH go build -o "$OUTPUT" ./cmd/myapp
    if [[ $? -ne 0 ]]; then
      echo "Failed to build Linux $ARCH"
    fi
  done
}

# Function to build for Windows
build_windows() {
  for ARCH in "${WINDOWS_ARCHS[@]}"; do
    OUTPUT="$BIN_DIR/$APP_NAME-windows-$ARCH.exe"
    echo "Building Windows $ARCH..."
    GOOS=windows GOARCH=$ARCH go build -o "$OUTPUT" ./cmd/myapp
    if [[ $? -ne 0 ]]; then
      echo "Failed to build Windows $ARCH"
    fi
  done
}

# Function to build for macOS
build_macos() {
  for ARCH in "${MACOS_ARCHS[@]}"; do
    OUTPUT="$BIN_DIR/$APP_NAME-darwin-$ARCH"
    echo "Building macOS $ARCH..."
    GOOS=darwin GOARCH=$ARCH go build -o "$OUTPUT" ./cmd/myapp
    if [[ $? -ne 0 ]]; then
      echo "Failed to build macOS $ARCH"
    fi
  done
}

# Clean old binaries
echo "Cleaning old binaries..."
rm -rf "$BIN_DIR"/*
mkdir -p "$BIN_DIR"

# Build all
build_linux
build_windows
build_macos

echo "Build completed!"
echo "Binaries are in $BIN_DIR"