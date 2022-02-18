package ui

import (
	"github.com/gonutz/wui/v2"
	"os"
)

func Show() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -10,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerSize(294, 221)
	window.SetTitle("Home Music")

	exitBtn := wui.NewButton()
	exitBtn.SetBounds(201, 184, 85, 25)
	exitBtn.SetText("Exit")
	exitBtn.SetOnClick(func() {
		os.Exit(0)
	})
	window.Add(exitBtn)

	playBtn := wui.NewButton()
	playBtn.SetBounds(10, 182, 85, 26)
	playBtn.SetText("Play")
	window.Add(playBtn)

	cmbInput := wui.NewComboBox()
	cmbInput.SetBounds(9, 38, 150, 21)
	cmbInput.SetItems([]string{
		"Input1",
		"Input2",
		"Input3",
	})
	cmbInput.SetSelectedIndex(0)
	window.Add(cmbInput)

	lblInput := wui.NewLabel()
	lblInput.SetBounds(8, 16, 150, 13)
	lblInput.SetText("Input:")
	window.Add(lblInput)

	lblOutput := wui.NewLabel()
	lblOutput.SetBounds(10, 83, 150, 13)
	lblOutput.SetText("Output:")
	window.Add(lblOutput)

	cmbOutput := wui.NewComboBox()
	cmbOutput.SetBounds(10, 105, 150, 21)
	cmbOutput.SetItems([]string{
		"Output1",
		"Output2",
		"Output3",
	})
	cmbOutput.SetSelectedIndex(0)
	window.Add(cmbOutput)

	window.Show()
}
