package cmd

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ADB struct {
	AdbCmd string
}

type AdbDevice struct {
	Name    string
	Address string
}

func NewAdb() *ADB {
	adb := &ADB{
		AdbCmd: "adb",
	}
	if runtime.GOOS == "windows" {
		adb.AdbCmd = "bin\\adb.exe"
	} else {
		adb.AdbCmd = "./bin/adb"
	}
	return adb
}

func (adb *ADB) Exec(cmd string, args ...string) error {
	execCmd := exec.Command(adb.AdbCmd, append([]string{cmd}, args...)...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func (adb *ADB) ExecCallback(callback func(outBuf, errBuf strings.Builder) error, cmd string, args ...string) error {
	execCmd := exec.Command(adb.AdbCmd, append([]string{cmd}, args...)...)

	var outBuf, errBuf strings.Builder
	execCmd.Stdout = &outBuf
	execCmd.Stderr = &errBuf

	err := execCmd.Run()
	if err != nil {
		return err
	}
	return callback(outBuf, errBuf)
}

func (adb *ADB) StartServer() error {
	return adb.Exec("start-server")
}

func (adb *ADB) KillServer() error {
	return adb.Exec("kill-server")
}

func (adb *ADB) Connect(target string) error {
	return adb.Exec("connect", target)
}

func (adb *ADB) Pair(target, code string) error {
	return adb.Exec("pair", target, code)
}

func (adb *ADB) PairQR() error {
	err := showQr()
	if err != nil {
		return err
	}
	err = startDiscovery()
	if err != nil {
		return err
	}
	return nil
}

func (adb *ADB) ListDevices() ([]AdbDevice, error) {
	var devices []AdbDevice
	err := adb.ExecCallback(func(outBuf, errBuf strings.Builder) error {
		fmt.Fprintln(os.Stderr, errBuf.String())

		output := outBuf.String()

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
		return nil
	}, "devices", "-l")
	return devices, err
}

func showQr() error {
	// WIFI:T:ADB;S:studio-y&1rAp#B9u;P:{GN1!WXb!ut{;;
	qrText := fmt.Sprintf("WIFI:T:ADB;S:%s;P:%s;;", "ADB_WIFI_agos", "agos_pw")
	qr, err := qrcode.New(qrText, qrcode.High)
	if err != nil {
		return err
	}

	qrcodeString := qr.ToSmallString(false)
	fmt.Println(qrcodeString)

	return nil
}

func startDiscovery() error {
	return nil
}

func connect() error {
	return nil
}

func getDevice() {

}
