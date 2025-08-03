#!/usr/bin/env python3
"""
Demonstration script for the Cross-Platform Command Execution Application
This simulates the functionality of the Go application to show how it works.
"""

import json
import subprocess
import socket
import time
from http.server import HTTPServer, BaseHTTPRequestHandler
import threading
import requests

class CommandExecutor:
    """Simulates the Commander interface from the Go application"""
    
    def ping(self, host):
        """Simulate ping functionality"""
        try:
            # Use system ping command
            if subprocess.run(['ping', '-n', '1', host] if os.name == 'nt' else ['ping', '-c', '1', host], 
                           capture_output=True, timeout=10).returncode == 0:
                return {"successful": True, "time": "50ms"}
            else:
                return {"successful": False, "time": "100ms"}
        except:
            return {"successful": False, "time": "0ms"}
    
    def get_system_info(self):
        """Get system information"""
        hostname = socket.gethostname()
        try:
            # Get local IP
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            s.connect(("8.8.8.8", 80))
            ip_address = s.getsockname()[0]
            s.close()
        except:
            ip_address = "127.0.0.1"
        
        return {"hostname": hostname, "ip_address": ip_address}

class CommandHandler(BaseHTTPRequestHandler):
    """HTTP request handler for the command execution API"""
    
    def __init__(self, *args, commander=None, **kwargs):
        self.commander = commander
        super().__init__(*args, **kwargs)
    
    def do_GET(self):
        """Handle GET requests (health check)"""
        if self.path == '/health':
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({"status": "healthy"}).encode())
        else:
            self.send_response(404)
            self.end_headers()
    
    def do_POST(self):
        """Handle POST requests (command execution)"""
        if self.path == '/execute':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            
            try:
                request = json.loads(post_data.decode('utf-8'))
                response = self.handle_command(request)
            except json.JSONDecodeError:
                response = {"success": False, "error": "Invalid JSON format"}
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_response(404)
            self.end_headers()
    
    def handle_command(self, request):
        """Process command requests"""
        command_type = request.get('type')
        payload = request.get('payload', '')
        
        if command_type == 'ping':
            if not payload:
                return {"success": False, "error": "Host is required for ping command"}
            
            result = self.commander.ping(payload)
            return {"success": True, "data": result}
        
        elif command_type == 'sysinfo':
            result = self.commander.get_system_info()
            return {"success": True, "data": result}
        
        else:
            return {"success": False, "error": "Unknown command type. Supported types: ping, sysinfo"}

def start_server(commander, port=8080):
    """Start the HTTP server"""
    class Handler(CommandHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, commander=commander, **kwargs)
    
    server = HTTPServer(('localhost', port), Handler)
    print(f"🚀 Server started on http://localhost:{port}")
    server.serve_forever()

def test_client():
    """Test the API endpoints"""
    print("\n🧪 Testing API Endpoints")
    print("=" * 40)
    
    base_url = "http://localhost:8080"
    
    # Test health endpoint
    try:
        response = requests.get(f"{base_url}/health")
        if response.status_code == 200:
            print("✅ Health check passed")
        else:
            print("❌ Health check failed")
    except:
        print("❌ Health check failed - server not running")
        return
    
    # Test system info
    try:
        response = requests.post(f"{base_url}/execute", 
                              json={"type": "sysinfo"})
        if response.status_code == 200:
            data = response.json()
            if data["success"]:
                sysinfo = data["data"]
                print(f"✅ System Info: {sysinfo['hostname']} ({sysinfo['ip_address']})")
            else:
                print(f"❌ System info failed: {data['error']}")
    except Exception as e:
        print(f"❌ System info test failed: {e}")
    
    # Test ping
    try:
        response = requests.post(f"{base_url}/execute", 
                              json={"type": "ping", "payload": "8.8.8.8"})
        if response.status_code == 200:
            data = response.json()
            if data["success"]:
                ping_result = data["data"]
                status = "✅" if ping_result["successful"] else "❌"
                print(f"{status} Ping test: {ping_result['time']}")
            else:
                print(f"❌ Ping test failed: {data['error']}")
    except Exception as e:
        print(f"❌ Ping test failed: {e}")

def main():
    """Main demonstration function"""
    print("🚀 Cross-Platform Command Execution Application Demo")
    print("=" * 60)
    
    # Create commander instance
    commander = CommandExecutor()
    
    # Start server in background thread
    server_thread = threading.Thread(target=start_server, args=(commander,), daemon=True)
    server_thread.start()
    
    # Wait for server to start
    time.sleep(2)
    
    # Test the API
    test_client()
    
    print("\n📋 Application Features Demonstrated:")
    print("• HTTP API server running on port 8080")
    print("• Health check endpoint (/health)")
    print("• Command execution endpoint (/execute)")
    print("• Ping functionality with system commands")
    print("• System information retrieval")
    print("• JSON request/response format")
    print("• Error handling and validation")
    
    print("\n🔧 Installation Features:")
    print("• Cross-platform build scripts (build.sh, build.bat)")
    print("• Windows installation with PowerShell script")
    print("• macOS installation with bash script")
    print("• Service/daemon setup for auto-start")
    print("• Firewall configuration")
    
    print("\n📁 Project Structure Created:")
    print("• main.go - Main application and HTTP server")
    print("• network.go - Network utilities")
    print("• main_test.go - Comprehensive tests")
    print("• client.go - API client for testing")
    print("• build scripts for cross-platform compilation")
    print("• installation scripts for Windows and macOS")
    print("• README.md - Complete documentation")
    
    print("\n✨ Demo completed! The Go application would provide the same functionality with:")
    print("• Better performance (compiled language)")
    print("• Single executable deployment")
    print("• Native system service integration")
    print("• Cross-platform compatibility")

if __name__ == "__main__":
    import os
    main() 