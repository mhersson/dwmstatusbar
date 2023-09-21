package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func ExecCommand(command string, args []string, shell bool) string {
	cmd := exec.Command(command, args...)
	if shell {
		cmd = exec.Command("bash", "-c", command)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to execute command %s\n", err.Error())
	}

	return strings.TrimSpace(string(out))
}

func Clock() string {
	t := time.Now()

	return t.Format("Monday 2006-01-02 15:04")
}

func PIA() string {
	args := []string{"get", "vpnip"}
	out := ExecCommand("piactl", args, false)

	return strings.TrimSpace(out)
}

func ExternalIP() string {
	var ipaddress []byte

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://icanhazip.com", nil)
	if err != nil {
		log.Printf("failed to create request %s\n", err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to get external ip %s\n", err.Error())
	}

	defer resp.Body.Close()

	ipaddress, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read external ip %s\n", err.Error())
	}

	return string(bytes.Trim(ipaddress, "\n\r\t "))
}

func KeyboardLayout() string {
	args := []string{"q"}
	out := ExecCommand("xset", args, false)

	layout := "US"

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "LED mask") {
			fields := strings.Fields(line)
			if len(fields) >= 10 {
				tmpLayout := fields[9]
				if tmpLayout == "00001000" {
					layout = "NO"
				}
			}
		}
	}

	return layout
}

func DPMS() string {
	args := []string{"q"}
	out := ExecCommand("xset", args, false)

	dpms := "DPMS ON"

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "DPMS is") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if fields[2] != "Enabled" {
					dpms = "DPMS OFF"
				}
			}
		}
	}

	return dpms
}
