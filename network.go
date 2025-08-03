package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func executePing(host string) (PingResult, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", host)
	case "darwin", "linux":
		cmd = exec.Command("ping", "-c", "1", host)
	default:
		return PingResult{}, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	start := time.Now()
	output, err := cmd.Output()
	duration := time.Since(start)

	if err != nil {
		return PingResult{Successful: false, Time: duration}, nil
	}

	outputStr := strings.ToLower(string(output))
	success := strings.Contains(outputStr, "time=") || strings.Contains(outputStr, "ttl=") || strings.Contains(outputStr, "bytes from")

	return PingResult{Successful: success, Time: duration}, nil
}
