//go:build ui

package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/floholz/agos/internal/core"
)

var fyneWindow fyne.Window
var agosApp *core.AgosApp

func LaunchUI(agos *core.AgosApp) {
	agosApp = agos
	a := app.NewWithID("dev.floholz.agos")
	fyneWindow = a.NewWindow("AGOS")
	fyneWindow.SetIcon(GetAgosIcon())
	fyneWindow.Resize(fyne.NewSize(600, 800))

	gotoHomePage()

	fyneWindow.SetMainMenu(mainMenu())
	fyneWindow.ShowAndRun()
}

func mainMenu() *fyne.MainMenu {
	agosMenu := fyne.NewMenu("Menu",
		fyne.NewMenuItem("About", gotoAboutPage),
	)
	pairMenu := fyne.NewMenu("Pair",
		fyne.NewMenuItem("New ...", gotoPairingPage),
		fyne.NewMenuItem("Pair with QR", gotoPairQrPage),
	)
	return fyne.NewMainMenu(agosMenu, pairMenu)
}

func gotoHomePage() {
	gotoDevicesPage()
}

func gotoDevicesPage() {
	devices, err := agosApp.Adb.ListDevices()
	if err != nil {
		devices = []core.AdbDevice{}
	}

	header := container.NewBorder(
		nil, widget.NewSeparator(), nil,
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), gotoDevicesPage),
			widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
				// TODO: implement connect logic + ui
			}),
		),
		widget.NewRichTextFromMarkdown("## Connected Devices"),
	)

	list := widget.NewList(
		func() int {
			return len(devices)
		},
		func() fyne.CanvasObject {
			return NewDeviceListItem()
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			textObjs := o.(*fyne.Container).Objects[0].(*fyne.Container).Objects
			textObjs[0].(*widget.Label).SetText(devices[i].Name)
			textObjs[1].(*canvas.Text).Text = devices[i].Address
			btnObjs := o.(*fyne.Container).Objects[1].(*fyne.Container).Objects
			btnObjs[0].(*widget.Button).OnTapped = func() {
				scErr := agosApp.ScrCpy.Run(devices[i].Address)
				if scErr != nil {
					fyne.LogError("Failed to start scrcpy", err)
				}
			}
			btnObjs[1].(*widget.Button).OnTapped = func() {
				err = agosApp.Adb.Disconnect(devices[i].Address)
				if err != nil {
					fyne.LogError("Failed to disconnect device", err)
				}
				gotoDevicesPage()
			}
		},
	)

	content := NewPaddedPage(container.NewBorder(
		header,
		nil, nil, nil,
		list,
	))
	fyneWindow.SetContent(content)
}

func gotoPairingPage() {
	var ip, port, code string
	ipBinding := binding.BindString(&ip)
	portBinding := binding.BindString(&port)
	codeBinding := binding.BindString(&code)
	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("IP", widget.NewEntryWithData(ipBinding)),
			widget.NewFormItem("Port", widget.NewEntryWithData(portBinding)),
			widget.NewFormItem("Code", widget.NewEntryWithData(codeBinding)),
		),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButtonWithIcon("", theme.GridIcon(), gotoPairQrPage),
			widget.NewButton("Pair Device", func() {
				err := agosApp.Adb.Pair(fmt.Sprintf("%s:%s", ip, port), code)
				if err != nil {
					fyne.LogError("Failed to pair device", err)
				}
			}),
		),
	)
	fyneWindow.SetContent(NewSubPageWithHeading("Pair a new device", content))
}

func gotoPairQrPage() {
	qrImg, pairData, err := agosApp.Adb.GeneratePairQR()
	if err != nil {
		fyne.LogError("Failed to generate pairing QR code", err)
		return
	}
	go func() {
		err = agosApp.Adb.PairQR(pairData, true)
		if err != nil {
			fyne.LogError("Failed to start pairing discovery", err)
		}
	}()

	image := canvas.NewImageFromImage(qrImg)
	image.FillMode = canvas.ImageFillOriginal

	content := container.NewVBox(
		container.New(
			layout.NewCustomPaddedLayout(50, 50, 50, 10),
			container.NewStack(image),
		),
		container.New(
			layout.NewCustomPaddedLayout(20, 0, 100, 100),
			widget.NewButton("Cancel Pairing", func() {
				// TODO: cancel pairing discovery
				gotoPairingPage()
			}),
		),
	)
	fyneWindow.SetContent(NewSubPageWithHeading("Pair a new device", content))
}

func gotoAboutPage() {
	fyneWindow.SetFixedSize(true)
	aboutText := widget.NewRichTextFromMarkdown(AboutMarkdown)
	aboutText.Wrapping = fyne.TextWrapWord

	logoImg := canvas.NewImageFromFile("assets/logo-color.png")
	logoImg.FillMode = canvas.ImageFillContain
	logoImg.SetMinSize(fyne.Size{Width: 128, Height: 128})

	content := container.NewBorder(
		nil, nil,
		container.NewVBox(logoImg),
		nil,
		aboutText,
	)
	paddedLayout := layout.NewCustomPaddedLayout(50, 50, 40, 40)
	fyneWindow.SetContent(NewSubPage(container.New(paddedLayout, content)))
	fyneWindow.Resize(fyne.NewSize(600, 800))
	fyneWindow.SetFixedSize(false)
}
