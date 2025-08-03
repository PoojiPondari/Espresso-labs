#!/bin/bash

# macOS Installation Script for Command Executor
# This script installs the application and sets up start-on-boot functionality

set -e

INSTALL_PATH="/usr/local/bin"
SERVICE_NAME="com.commandexecutor.service"
PLIST_PATH="/Library/LaunchDaemons/$SERVICE_NAME.plist"

echo "Installing Command Executor for macOS..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)" 
   exit 1
fi

# Create installation directory if it doesn't exist
if [ ! -d "$INSTALL_PATH" ]; then
    mkdir -p "$INSTALL_PATH"
    echo "Created installation directory: $INSTALL_PATH"
fi

# Copy executable to installation directory
EXECUTABLE_NAME="command-executor-darwin-amd64"
SOURCE_PATH="./build/$EXECUTABLE_NAME"
DEST_PATH="$INSTALL_PATH/$EXECUTABLE_NAME"

if [ -f "$SOURCE_PATH" ]; then
    cp "$SOURCE_PATH" "$DEST_PATH"
    chmod +x "$DEST_PATH"
    echo "Copied executable to: $DEST_PATH"
else
    echo "Executable not found at: $SOURCE_PATH"
    echo "Please run the build script first."
    exit 1
fi

# Create launchd plist file for start-on-boot
cat > "$PLIST_PATH" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$SERVICE_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$DEST_PATH</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/commandexecutor.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/commandexecutor.error.log</string>
    <key>UserName</key>
    <string>root</string>
    <key>GroupName</key>
    <string>wheel</string>
</dict>
</plist>
EOF

echo "Created launchd plist file: $PLIST_PATH"

# Set proper permissions
chmod 644 "$PLIST_PATH"
chown root:wheel "$PLIST_PATH"

# Load the service
launchctl load "$PLIST_PATH"
echo "Loaded launchd service: $SERVICE_NAME"

# Start the service
launchctl start "$SERVICE_NAME"
echo "Started service successfully"

# Create log files
touch /var/log/commandexecutor.log
touch /var/log/commandexecutor.error.log
chmod 644 /var/log/commandexecutor.log
chmod 644 /var/log/commandexecutor.error.log

echo "Installation completed successfully!"
echo "The application is now running as a launchd service and will start automatically on boot."
echo "Service name: $SERVICE_NAME"
echo "HTTP server running on: http://localhost:8080"
echo "Log files:"
echo "  - /var/log/commandexecutor.log"
echo "  - /var/log/commandexecutor.error.log" 