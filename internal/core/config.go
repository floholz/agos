package core

import "runtime"

type Config struct {
	AdbPath    string
	ScrCpyPath string
	PortRange  PortRange
}

type PortRange struct {
	MinPort int
	MaxPort int
}

var DefaultConfig = defaultConfig()

func defaultConfig() *Config {
	config := &Config{
		AdbPath:    "adb",
		ScrCpyPath: "scrcpy",
		PortRange: PortRange{
			MinPort: 32000,
			MaxPort: 48000,
		},
	}
	if runtime.GOOS == "windows" {
		config.AdbPath = "bin\\windows\\adb.exe"
		config.ScrCpyPath = "bin\\windows\\scrcpy.exe"
	}
	return config
}
