package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type AgosApp struct {
	AdbCmd    string
	ScrcpyCmd string
}

type AdbDevice struct {
	Name    string
	Address string
}

func NewAgosApp() *AgosApp {
	app := &AgosApp{
		AdbCmd:    "adb",
		ScrcpyCmd: "scrcpy",
	}

	if runtime.GOOS == "windows" {
		app.AdbCmd = "bin\\adb.exe"
		app.ScrcpyCmd = "bin\\scrcpy.exe"
	} else {
		app.AdbCmd = "./bin/adb"
		app.ScrcpyCmd = "./bin/scrcpy"
	}
	return app
}

func (app *AgosApp) StartAdb() error {
	adbStartCmd := exec.Command(app.AdbCmd, "start-server")
	adbStartCmd.Stdout = os.Stdout
	adbStartCmd.Stderr = os.Stderr
	return adbStartCmd.Run()
}

func (app *AgosApp) AdbListDevices() ([]AdbDevice, error) {
	cmd := exec.Command(app.AdbCmd, "devices", "-l")
	cmd.Stderr = os.Stderr
	var outbuf strings.Builder
	cmd.Stdout = &outbuf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	output := outbuf.String()
	var devices []AdbDevice

	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] { // skip header
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		ipPort := parts[0] // full IP:port

		var model string
		for _, part := range parts {
			if strings.HasPrefix(part, "model:") {
				model = strings.TrimPrefix(part, "model:")
				break
			}
		}

		devices = append(devices, AdbDevice{
			Name:    model,
			Address: ipPort,
		})
	}

	return devices, nil
}

func (app *AgosApp) AdbConnect(target string) error {
	connectCmd := exec.Command(app.AdbCmd, "connect", target)
	connectCmd.Stdout = os.Stdout
	connectCmd.Stderr = os.Stderr

	if err := connectCmd.Run(); err != nil {
		fmt.Println("Failed to connect via adb:", err)
		return err
	}

	fmt.Println("Successfully connected to", target)
	return nil
}

func (app *AgosApp) RunScrcpy(ip string) error {
	scrcpyCmdExec := exec.Command(app.ScrcpyCmd, "-s", ip)
	scrcpyCmdExec.Stdout = os.Stdout
	scrcpyCmdExec.Stderr = os.Stderr

	if err := scrcpyCmdExec.Run(); err != nil {
		fmt.Println("Failed to run scrcpy:", err)
		return err
	}

	fmt.Println("scrcpy session ended.")
	return nil
}
