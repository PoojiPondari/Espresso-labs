@echo off

REM Cross-platform build script for Windows
REM This script builds executables for Windows and macOS

echo Building cross-platform command execution application...

REM Create build directory
if not exist build mkdir build

REM Build for Windows (64-bit)
echo Building for Windows (64-bit)...
set GOOS=windows
set GOARCH=amd64
go build -o build\command-executor-windows-amd64.exe .

REM Build for macOS (64-bit)
echo Building for macOS (64-bit)...
set GOOS=darwin
set GOARCH=amd64
go build -o build\command-executor-darwin-amd64 .

REM Build for macOS (Apple Silicon)
echo Building for macOS (Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -o build\command-executor-darwin-arm64 .

REM Build for Linux (64-bit)
echo Building for Linux (64-bit)...
set GOOS=linux
set GOARCH=amd64
go build -o build\command-executor-linux-amd64 .

echo Build complete! Executables created in build\ directory:
dir build\

pause 