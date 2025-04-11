package core

type AgosApp struct {
	Adb    *ADB
	ScrCpy *ScrCpy
	Config *Config
}

func NewAgosApp() *AgosApp {
	return &AgosApp{
		Adb:    NewAdb(DefaultConfig.AdbPath, &DefaultConfig.PortRange),
		ScrCpy: NewScrCpy(DefaultConfig.ScrCpyPath),
		Config: DefaultConfig,
	}
}

func NewAgosAppWithConfig(config *Config) *AgosApp {
	return &AgosApp{
		Adb:    NewAdb(config.AdbPath, &config.PortRange),
		ScrCpy: NewScrCpy(config.ScrCpyPath),
		Config: config,
	}
}
