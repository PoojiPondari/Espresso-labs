package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Commander interface defines the contract for command execution
type Commander interface {
	Ping(host string) (PingResult, error)
	GetSystemInfo() (SystemInfo, error)
}

// PingResult represents the result of a ping operation
type PingResult struct {
	Successful bool          `json:"successful"`
	Time       time.Duration `json:"time"`
}

// SystemInfo represents system information
type SystemInfo struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
}

// CommandRequest represents incoming command requests
type CommandRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// CommandResponse represents the response to command requests
type CommandResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// commander implements the Commander interface
type commander struct{}

// NewCommander creates a new commander instance
func NewCommander() Commander {
	return &commander{}
}

// GetSystemInfo retrieves system information
func (c *commander) GetSystemInfo() (SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return SystemInfo{}, err
	}

	ipAddress, err := getLocalIP()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		Hostname:  hostname,
		IPAddress: ipAddress,
	}, nil
}

// Ping executes a ping command to the specified host
func (c *commander) Ping(host string) (PingResult, error) {
	return executePing(host)
}

// handleRequests creates the HTTP request handler
func handleRequests(cmdr Commander) http.Handler {
	mux := http.NewServeMux()

	// Root handler (GET /)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h2>✅ Command Executor API is running!</h2><p>Use <code>/health</code> or <code>/execute</code> endpoints.</p>`)
	})

	// Health check
	mux.HandleFunc("/health", handleHealth)

	// Command executor
	mux.HandleFunc("/execute", handleCommand(cmdr))

	return mux
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// handleCommand handles command execution requests
func handleCommand(cmdr Commander) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var req CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(CommandResponse{Success: false, Error: "Invalid request format"})
			return
		}

		var response CommandResponse

		switch req.Type {
		case "ping":
			if req.Payload == "" {
				response = CommandResponse{Success: false, Error: "Host is required for ping command"}
			} else {
				result, err := cmdr.Ping(req.Payload)
				if err != nil {
					response = CommandResponse{Success: false, Error: err.Error()}
				} else {
					response = CommandResponse{Success: true, Data: result}
				}
			}

		case "sysinfo":
			result, err := cmdr.GetSystemInfo()
			if err != nil {
				response = CommandResponse{Success: false, Error: err.Error()}
			} else {
				response = CommandResponse{Success: true, Data: result}
			}

		default:
			response = CommandResponse{Success: false, Error: "Unknown command type. Supported types: ping, sysinfo"}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// main starts the HTTP server
func main() {
	commander := NewCommander()
	server := &http.Server{
		Addr:    ":8080",
		Handler: handleRequests(commander),
	}

	log.Println("Starting server on http://localhost:8080/")
	log.Fatal(server.ListenAndServe())
}
