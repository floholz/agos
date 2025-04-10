package cmd

type AgosApp struct {
	Adb    *ADB
	ScrCpy *ScrCpy
}

func NewAgosApp() *AgosApp {
	return &AgosApp{
		Adb:    NewAdb(),
		ScrCpy: NewScrCpy(),
	}
}
