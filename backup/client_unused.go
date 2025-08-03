package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client for the command executor API
type Client struct {
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HealthCheck checks if the service is healthy
func (c *Client) HealthCheck() error {
	resp, err := c.Client.Get(c.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode health response: %v", err)
	}

	if result["status"] != "healthy" {
		return fmt.Errorf("service is not healthy: %s", result["status"])
	}

	fmt.Println("✅ Service is healthy")
	return nil
}

// GetSystemInfo retrieves system information
func (c *Client) GetSystemInfo() error {
	req := CommandRequest{
		Type: "sysinfo",
	}

	resp, err := c.executeCommand(req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("failed to get system info: %s", resp.Error)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response data format")
	}

	fmt.Printf("✅ System Information:\n")
	fmt.Printf("   Hostname: %s\n", data["hostname"])
	fmt.Printf("   IP Address: %s\n", data["ip_address"])

	return nil
}

// PingHost pings a specified host
func (c *Client) PingHost(host string) error {
	req := CommandRequest{
		Type:    "ping",
		Payload: host,
	}

	resp, err := c.executeCommand(req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("failed to ping host: %s", resp.Error)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response data format")
	}

	successful := data["successful"].(bool)
	timeStr := data["time"].(string)

	if successful {
		fmt.Printf("✅ Ping to %s successful (time: %s)\n", host, timeStr)
	} else {
		fmt.Printf("❌ Ping to %s failed (time: %s)\n", host, timeStr)
	}

	return nil
}

// executeCommand sends a command request to the API
func (c *Client) executeCommand(req CommandRequest) (*CommandResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.Client.Post(
		c.BaseURL+"/execute",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var commandResp CommandResponse
	if err := json.Unmarshal(body, &commandResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &commandResp, nil
}

func main() {
	client := NewClient("http://localhost:8080")

	fmt.Println("🚀 Command Executor API Client")
	fmt.Println("================================")

	// Health check
	if err := client.HealthCheck(); err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		return
	}

	fmt.Println()

	// Get system info
	if err := client.GetSystemInfo(); err != nil {
		fmt.Printf("❌ Failed to get system info: %v\n", err)
	} else {
		fmt.Println()
	}

	// Ping Google DNS
	if err := client.PingHost("8.8.8.8"); err != nil {
		fmt.Printf("❌ Failed to ping host: %v\n", err)
	} else {
		fmt.Println()
	}

	// Ping localhost
	if err := client.PingHost("127.0.0.1"); err != nil {
		fmt.Printf("❌ Failed to ping localhost: %v\n", err)
	}

	fmt.Println("\n✨ All tests completed!")
} 