//go:build ui

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

func NewSubPage(content ...fyne.CanvasObject) *fyne.Container {
	return NewPaddedPage(
		container.NewBorder(
			container.NewHBox(widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), gotoHomePage)),
			nil, nil, nil,
			content...,
		),
	)
}
func NewSubPageWithHeading(heading string, content ...fyne.CanvasObject) *fyne.Container {
	return NewPaddedPage(
		container.NewBorder(
			container.NewBorder(
				nil, nil,
				container.NewHBox(widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), gotoHomePage)),
				nil,
				container.NewCenter(widget.NewRichTextFromMarkdown("## "+heading)),
			),
			nil, nil, nil,
			content...,
		),
	)
}

func NewPaddedPage(content ...fyne.CanvasObject) *fyne.Container {
	padded := layout.NewCustomPaddedLayout(10, 10, 10, 10)
	return container.New(
		padded,
		content...,
	)
}

// NewDeviceListItem: BorderContainer > (r:HBox > Button, Button), widget.Lable, canvas.Text
func NewDeviceListItem() fyne.CanvasObject {
	nameLabel := widget.NewLabel("Name")
	addressLabel := canvas.NewText("Address", colornames.Grey)
	return container.NewBorder(
		nil, nil, nil,
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.MediaVideoIcon(), nil),
			widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
		),
		container.NewHBox(nameLabel, addressLabel),
	)
}
