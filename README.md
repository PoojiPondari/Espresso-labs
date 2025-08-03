# Cross-Platform Command Execution Application

A cross-platform command execution application written in Go that can be installed on Windows and macOS, capable of executing simple system commands and returning their results via HTTP API.

## Features

- **Cross-Platform Support**: Runs on Windows, macOS, and Linux
- **HTTP API**: RESTful API for command execution
- **System Commands**: Ping network hosts and get system information
- **Auto-Start**: Configurable start-on-boot functionality
- **Service Management**: Runs as a system service/daemon

## Supported Commands

### 1. Network Ping
Ping a specified host and return response time.

**Request:**
```json
{
  "type": "ping",
  "payload": "8.8.8.8"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "successful": true,
    "time": "50ms"
  }
}
```

### 2. System Information
Get hostname and IP address of the system.

**Request:**
```json
{
  "type": "sysinfo"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "hostname": "my-computer",
    "ip_address": "192.168.1.100"
  }
}
```

## API Endpoints

### POST /execute
Execute a command and return results.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "type": "ping|sysinfo",
  "payload": "string (required for ping)"
}
```

**Response:**
```json
{
  "success": true|false,
  "data": "command result object",
  "error": "error message (if success is false)"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Build Instructions

### Prerequisites
- Go 1.21 or later
- Git

### Building the Application

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd job-trends
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build for all platforms:**
   
   **On Linux/macOS:**
   ```bash
   chmod +x build.sh
   ./build.sh
   ```
   
   **On Windows:**
   ```cmd
   build.bat
   ```

4. **Build for specific platform:**
   ```bash
   # Windows
   GOOS=windows GOARCH=amd64 go build -o command-executor.exe .
   
   # macOS (Intel)
   GOOS=darwin GOARCH=amd64 go build -o command-executor .
   
   # macOS (Apple Silicon)
   GOOS=darwin GOARCH=arm64 go build -o command-executor .
   
   # Linux
   GOOS=linux GOARCH=amd64 go build -o command-executor .
   ```

## Installation Guide

### Windows Installation

1. **Build the application** (see Build Instructions above)

2. **Run the installation script as Administrator:**
   ```powershell
   # Open PowerShell as Administrator
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   .\install-windows.ps1
   ```

3. **Verify installation:**
   ```powershell
   Get-Service -Name CommandExecutor
   ```

4. **Test the API:**
   ```powershell
   Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
   ```

### macOS Installation

1. **Build the application** (see Build Instructions above)

2. **Run the installation script as root:**
   ```bash
   sudo chmod +x install-macos.sh
   sudo ./install-macos.sh
   ```

3. **Verify installation:**
   ```bash
   launchctl list | grep commandexecutor
   ```

4. **Test the API:**
   ```bash
   curl http://localhost:8080/health
   ```

## Manual Installation

### Windows (Manual)

1. Copy the executable to a directory (e.g., `C:\Program Files\CommandExecutor\`)
2. Create a Windows service:
   ```powershell
   New-Service -Name "CommandExecutor" -BinaryPathName "C:\Program Files\CommandExecutor\command-executor-windows-amd64.exe" -StartupType Automatic
   Start-Service CommandExecutor
   ```

### macOS (Manual)

1. Copy the executable to `/usr/local/bin/`
2. Create a launchd plist file in `/Library/LaunchDaemons/`
3. Load and start the service:
   ```bash
   sudo launchctl load /Library/LaunchDaemons/com.commandexecutor.service.plist
   sudo launchctl start com.commandexecutor.service
   ```

## Testing

### Run Tests
```bash
go test -v
```

### Manual Testing

1. **Start the application:**
   ```bash
   go run .
   ```

2. **Test health endpoint:**
   ```bash
   curl http://localhost:8080/health
   ```

3. **Test system info:**
   ```bash
   curl -X POST http://localhost:8080/execute \
     -H "Content-Type: application/json" \
     -d '{"type":"sysinfo"}'
   ```

4. **Test ping:**
   ```bash
   curl -X POST http://localhost:8080/execute \
     -H "Content-Type: application/json" \
     -d '{"type":"ping","payload":"8.8.8.8"}'
   ```

## Project Structure

```
job-trends/
├── main.go              # Main application and HTTP server
├── network.go           # Network utilities and ping implementation
├── main_test.go         # Unit tests
├── go.mod               # Go module file
├── build.sh             # Cross-platform build script (Linux/macOS)
├── build.bat            # Cross-platform build script (Windows)
├── install-windows.ps1  # Windows installation script
├── install-macos.sh     # macOS installation script
└── README.md           # This file
```

## Architecture

### Core Components

1. **Commander Interface**: Defines the contract for command execution
2. **HTTP Server**: RESTful API for receiving commands
3. **Platform-Specific Implementation**: Handles OS-specific command execution
4. **Service Management**: Start-on-boot functionality

### Design Patterns

- **Interface Segregation**: Commander interface for command execution
- **Dependency Injection**: Commander injected into HTTP handlers
- **Factory Pattern**: NewCommander() creates platform-specific implementations

## Error Handling

The application includes comprehensive error handling:

- Invalid command types
- Missing required parameters
- Network connectivity issues
- System command failures
- HTTP request validation

## Security Considerations

- Runs as a system service with appropriate permissions
- Firewall rules configured for port 8080
- Input validation for all HTTP requests
- Error messages don't expose sensitive information

## Troubleshooting

### Common Issues

1. **Port already in use:**
   - Check if another service is using port 8080
   - Modify the port in main.go if needed

2. **Permission denied:**
   - Ensure installation scripts are run with administrator/root privileges

3. **Service not starting:**
   - Check system logs for error messages
   - Verify executable permissions

4. **Ping not working:**
   - Ensure network connectivity
   - Check firewall settings

### Logs

- **Windows**: Event Viewer → Windows Logs → Application
- **macOS**: `/var/log/commandexecutor.log` and `/var/log/commandexecutor.error.log`

## Development

### Adding New Commands

1. Add new method to Commander interface
2. Implement the method in commander struct
3. Add case in handleCommand function
4. Update tests

### Building for Distribution

```bash
# Create release builds
./build.sh

# Create zip archives
cd build
zip command-executor-windows.zip command-executor-windows-amd64.exe
zip command-executor-macos.zip command-executor-darwin-*
```

## License

This project is provided as-is for educational and demonstration purposes.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Support

For issues and questions, please create an issue in the repository. 


https://github.com/user-attachments/assets/41b4c594-15b1-43de-8401-fd328219faa7

