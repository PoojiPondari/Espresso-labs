#!/bin/bash

# Cross-platform build script for the command execution application
# This script builds executables for Windows and macOS

echo "Building cross-platform command execution application..."

# Create build directory
mkdir -p build

# Build for Windows (64-bit)
echo "Building for Windows (64-bit)..."
GOOS=windows GOARCH=amd64 go build -o build/command-executor-windows-amd64.exe .

# Build for macOS (64-bit)
echo "Building for macOS (64-bit)..."
GOOS=darwin GOARCH=amd64 go build -o build/command-executor-darwin-amd64 .

# Build for macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o build/command-executor-darwin-arm64 .

# Build for Linux (64-bit)
echo "Building for Linux (64-bit)..."
GOOS=linux GOARCH=amd64 go build -o build/command-executor-linux-amd64 .

echo "Build complete! Executables created in build/ directory:"
ls -la build/ 