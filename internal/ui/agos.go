package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var AgosVersion = "0.3.0"

var AboutMarkdown = "## AGOS v" + AgosVersion +
	`
Agos is a wrapper application for adb & scrcpy. It was created to make it easy and simple to screenshare from any 
android device to your pc. To simplify the adb setup process even more, Agos provides the functionality to automatically 
scan for adb ports and also to pair devices using a QR code.    

   
created by [floholz](https://github.com/floholz)   

   
[agos.floholz.dev](https://agos.floholz.dev)   
`

func GetAgosIcon() fyne.Resource {
	res, err := fyne.LoadResourceFromPath("assets/logo.png")
	if err != nil {
		fyne.LogError("Error loading logo.png", err)
		return theme.ComputerIcon()
	}
	return res
}
