package core

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"
	"os"
	"os/exec"
	"strings"
)

type ADB struct {
	AdbCmd    string
	PortRange *PortRange
}

type AdbDevice struct {
	Name    string
	Address string
}

type AdbPairingData struct {
	Name     string
	Password string
}

func NewAdb(adbCmd string, portRange *PortRange) *ADB {
	adb := &ADB{
		AdbCmd:    adbCmd,
		PortRange: portRange,
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

func (adb *ADB) PairQR(connectAfterPair bool) error {
	pairingData, err := generatePairingData()
	if err != nil {
		return err
	}
	err = showQr(pairingData)
	if err != nil {
		return err
	}
	device, err := DiscoverAdbPairing()
	if err != nil {
		return err
	}
	if pairingData.Name != device.Instance {
		return fmt.Errorf("pairing name mismatch")
	}

	pairAddress := fmt.Sprintf("%s:%d", device.IPv4, device.Port)
	err = adb.Pair(pairAddress, pairingData.Password)
	if err != nil {
		return err
	}

	if connectAfterPair {
		port, err := DiscoverAdbPort(device.IPv4, *adb.PortRange)
		if err != nil {
			return err
		}
		connectAddress := fmt.Sprintf("%s:%d", device.IPv4, port)
		err = adb.Connect(connectAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (adb *ADB) ListDevices() ([]AdbDevice, error) {
	var devices []AdbDevice
	err := adb.ExecCallback(func(outBuf, errBuf strings.Builder) error {
		_, _ = fmt.Fprintln(os.Stderr, errBuf.String())

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

			ipPort := parts[0] // gui IP:port

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

func generatePairingData() (*AdbPairingData, error) {
	// WIFI:T:ADB;S:studio-y&1rAp#B9u;P:{GN1!WXb!ut{;;
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	pw, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	return &AdbPairingData{
		Name:     "ADB_WIFI_" + id,
		Password: pw,
	}, nil
}

func showQr(data *AdbPairingData) error {
	qrText := fmt.Sprintf("WIFI:T:ADB;S:%s;P:%s;;", data.Name, data.Password)
	qr, err := qrcode.New(qrText, qrcode.Medium)
	if err != nil {
		return err
	}

	qrcodeString := qr.ToSmallString(false)
	fmt.Println(qrcodeString)
	return nil
}
