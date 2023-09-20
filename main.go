package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type DataUpdater struct {
	Command  func() string
	Channel  chan string
	Interval time.Duration
}

func updateData(updater DataUpdater) {
	for {
		data := updater.Command()
		updater.Channel <- data
		time.Sleep(updater.Interval)
	}
}

func main() {
	var dpms, layout, ip, clock string

	dataUpdaters := []DataUpdater{
		{Command: DPMS, Channel: make(chan string), Interval: 10 * time.Second},
		{Command: KeyboardLayout, Channel: make(chan string), Interval: 1 * time.Second},
		{Command: IP, Channel: make(chan string), Interval: 600 * time.Second},
		{Command: Clock, Channel: make(chan string), Interval: 60 * time.Second},
	}

	for _, updater := range dataUpdaters {
		go updateData(updater)
	}

	go func() {
		for {
			select {
			case data := <-dataUpdaters[0].Channel:
				dpms = data
			case data := <-dataUpdaters[1].Channel:
				layout = data
			case data := <-dataUpdaters[2].Channel:
				ip = data
			case data := <-dataUpdaters[3].Channel:
				clock = data
			}
			status := fmt.Sprintf("󰌵 %s | 󰌌 %s | 󱇱 %s |  %s", dpms, layout, ip, clock)
			ExecCommand("xsetroot", []string{"-name", status}, false)
			// fmt.Println(status)
		}
	}()

	select {}
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

func ExecCommand(command string, args []string, shell bool) string {
	cmd := exec.Command(command, args...)
	if shell {
		cmd = exec.Command("bash", "-c", strings.Join(append([]string{command}, args...), " "))
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(out))
}
