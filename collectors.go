package main

import (
	"bytes"
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
		cmd = exec.Command("bash", "-c", strings.Join(append([]string{command}, args...), " "))
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

func IP() string {
	var ip []byte

	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		log.Printf("failed to get external ip %s\n", err.Error())
	}

	defer resp.Body.Close()

	ip, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read external ip %s\n", err.Error())
	}

	return string(bytes.Trim(ip, "\n\r\t "))
}

func KeyboardLayout() (layout string) {
	layout = "US"
	args := []string{"q"}
	out := ExecCommand("xset", args, false)

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

func DPMS() (dpms string) {
	args := []string{"q"}
	out := ExecCommand("xset", args, false)

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "DPMS is") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if fields[2] == "Enabled" {
					dpms = "DPMS ON"
				} else {
					dpms = "DPMS OFF"
				}
			}
		}
	}

	return dpms
}
