# Windows Installation Script for Command Executor
# This script installs the application and sets up start-on-boot functionality

param(
    [string]$InstallPath = "C:\Program Files\CommandExecutor"
)

Write-Host "Installing Command Executor for Windows..." -ForegroundColor Green

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "This script requires administrator privileges. Please run as administrator." -ForegroundColor Red
    exit 1
}

# Create installation directory
if (!(Test-Path $InstallPath)) {
    New-Item -ItemType Directory -Path $InstallPath -Force
    Write-Host "Created installation directory: $InstallPath" -ForegroundColor Yellow
}

# Copy executable to installation directory
$ExecutableName = "command-executor-windows-amd64.exe"
$SourcePath = ".\build\$ExecutableName"
$DestPath = "$InstallPath\$ExecutableName"

if (Test-Path $SourcePath) {
    Copy-Item $SourcePath $DestPath -Force
    Write-Host "Copied executable to: $DestPath" -ForegroundColor Yellow
} else {
    Write-Host "Executable not found at: $SourcePath" -ForegroundColor Red
    Write-Host "Please run the build script first." -ForegroundColor Red
    exit 1
}

# Create Windows Service for start-on-boot
$ServiceName = "CommandExecutor"
$ServiceDisplayName = "Command Executor Service"
$ServiceDescription = "Cross-platform command execution service"

# Check if service already exists
$ExistingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue

if ($ExistingService) {
    Write-Host "Service already exists. Stopping and removing..." -ForegroundColor Yellow
    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    Remove-Service -Name $ServiceName -Force
}

# Create the service
try {
    New-Service -Name $ServiceName `
                -DisplayName $ServiceDisplayName `
                -Description $ServiceDescription `
                -BinaryPathName "`"$DestPath`"" `
                -StartupType Automatic

    Write-Host "Created Windows service: $ServiceName" -ForegroundColor Green
    
    # Start the service
    Start-Service -Name $ServiceName
    Write-Host "Started service successfully" -ForegroundColor Green
    
} catch {
    Write-Host "Failed to create service: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "You may need to run this script as administrator." -ForegroundColor Red
    exit 1
}

# Create firewall rule
try {
    New-NetFirewallRule -DisplayName "Command Executor" `
                        -Direction Inbound `
                        -Protocol TCP `
                        -LocalPort 8080 `
                        -Action Allow `
                        -ErrorAction SilentlyContinue
    
    Write-Host "Created firewall rule for port 8080" -ForegroundColor Green
} catch {
    Write-Host "Failed to create firewall rule: $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "Installation completed successfully!" -ForegroundColor Green
Write-Host "The application is now running as a Windows service and will start automatically on boot." -ForegroundColor Green
Write-Host "Service name: $ServiceName" -ForegroundColor Cyan
Write-Host "HTTP server running on: http://localhost:8080" -ForegroundColor Cyan 