package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetSystemInfo(t *testing.T) {
	cmdr := NewCommander()
	info, err := cmdr.GetSystemInfo()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info.Hostname == "" {
		t.Error("Expected hostname to be non-empty")
	}

	if info.IPAddress == "" {
		t.Error("Expected IP address to be non-empty")
	}
}

func TestPing(t *testing.T) {
	cmdr := NewCommander()
	result, err := cmdr.Ping("8.8.8.8")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ping time should be reasonable (less than 10 seconds)
	if result.Time > 10*time.Second {
		t.Errorf("Expected ping time to be reasonable, got %v", result.Time)
	}
}

func TestHTTPEndpoints(t *testing.T) {
	commander := NewCommander()
	handler := handleRequests(commander)

	// Test health endpoint
	t.Run("Health", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]string
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["status"] != "healthy" {
			t.Errorf("Expected status 'healthy', got %s", response["status"])
		}
	})

	// Test sysinfo endpoint
	t.Run("SystemInfo", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(`{"type":"sysinfo"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response CommandResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success true, got false. Error: %s", response.Error)
		}

		// Check that data contains system info
		data, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatal("Expected data to be a map")
		}

		if data["hostname"] == "" {
			t.Error("Expected hostname to be non-empty")
		}

		if data["ip_address"] == "" {
			t.Error("Expected ip_address to be non-empty")
		}
	})

	// Test ping endpoint
	t.Run("Ping", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(`{"type":"ping","payload":"8.8.8.8"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response CommandResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success true, got false. Error: %s", response.Error)
		}

		// Check that data contains ping result
		data, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatal("Expected data to be a map")
		}

		if data["successful"] == nil {
			t.Error("Expected successful field to be present")
		}

		if data["time"] == nil {
			t.Error("Expected time field to be present")
		}
	})

	// Test invalid command type
	t.Run("InvalidCommand", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(`{"type":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response CommandResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success false for invalid command")
		}

		if response.Error == "" {
			t.Error("Expected error message for invalid command")
		}
	})

	// Test ping without payload
	t.Run("PingWithoutPayload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/execute", strings.NewReader(`{"type":"ping"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response CommandResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success false for ping without payload")
		}

		if response.Error == "" {
			t.Error("Expected error message for ping without payload")
		}
	})
} 