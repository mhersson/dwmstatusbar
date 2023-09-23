package dwmstatusbar

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var ExternalIPURL = "https://icanhazip.com"

type CmdInterface interface {
	CombinedOutput() ([]byte, error)
}

type Cmd struct {
	Cmd *exec.Cmd
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
	return c.Cmd.CombinedOutput() //nolint:wrapcheck
}

type Exec struct{}

func (e *Exec) NewCommand(command string, args ...string) CmdInterface { //nolint:ireturn
	cmd := exec.Command(command, args...)

	return &Cmd{Cmd: cmd}
}

type ExecInterface interface {
	NewCommand(command string, args ...string) CmdInterface
}

var ExecCmd ExecInterface = &Exec{}

func ExecCommand(command string, args []string, shell bool) string {
	cmd := ExecCmd.NewCommand(command, args...)
	if shell {
		cmd = ExecCmd.NewCommand("bash", "-c", command)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		Log.Printf("failed to execute command %s\n", err.Error())
	}

	if len(out) == 0 {
		return "No Data"
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
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ExternalIPURL, nil)
	if err != nil {
		Log.Printf("failed to create request %s\n", err.Error())

		return ""
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		Log.Printf("failed to get external ip %s\n", err.Error())

		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Printf("failed to read external ip %s\n", err.Error())

		return ""
	}

	return string(bytes.Trim(body, "\n\r\t "))
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