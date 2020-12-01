package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/disintegration/gift"
)

func (a *App) loadStatusBar() *widget.Box {
	a.imagePathLabel = widget.NewLabel("Path: ")
	a.statusBar = widget.NewVBox(
		widget.NewSeparator(),
		widget.NewHBox(
			a.imagePathLabel,
			layout.NewSpacer(),
		))
	return a.statusBar
}

func (a *App) loadInformationWidgets() *container.Scroll {
	a.widthLabel = widget.NewLabel("Width: ")
	a.heightLabel = widget.NewLabel("Height: ")
	a.informationWidgets = container.NewScroll(
		widget.NewHBox(
			widget.NewVBox(
				widget.NewLabelWithStyle("Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				a.widthLabel,
				a.heightLabel,
			),
			widget.NewSeparator(),
		),
	)
	a.informationWidgets.SetMinSize(fyne.NewSize(150, a.mainWin.Canvas().Size().Height))
	return a.informationWidgets
}

func (a *App) loadEditControls() *container.Scroll {
	a.sliderBrightness = widget.NewSlider(-100, 100)
	a.editBrightness = NewEditingOption(
		"Brightness: ",
		a.sliderBrightness,
		func(f float64) { a.changeParameter(&a.img.brightness, gift.Brightness(float32(f)), a.autochange) },
		0,
	)

	a.sliderContrast = widget.NewSlider(-100, 100)
	a.editContrast = NewEditingOption(
		"Contrast: ",
		a.sliderContrast,
		func(f float64) { a.changeParameter(&a.img.contrast, gift.Contrast(float32(f)), a.autochange) },
		0,
	)

	a.sliderHue = widget.NewSlider(-180, 180)
	a.editHue = NewEditingOption(
		"Hue: ",
		a.sliderHue,
		func(f float64) { a.changeParameter(&a.img.hue, gift.Hue(float32(f)), a.autochange) },
		0,
	)

	a.sliderColorBalanceR = widget.NewSlider(-100, 500)
	a.editColorBalanceR = NewEditingOption(
		"Red: ",
		a.sliderColorBalanceR,
		func(f float64) {
			a.changeParameter(&a.img.cbRed, gift.ColorBalance(
				float32(f), float32(a.sliderColorBalanceG.Value), float32(a.sliderColorBalanceB.Value)), a.autochange)
		},
		0,
	)

	a.sliderColorBalanceG = widget.NewSlider(-100, 500)
	a.editColorBalanceG = NewEditingOption(
		"Green: ",
		a.sliderColorBalanceG,
		func(f float64) {
			a.changeParameter(&a.img.cbGreen, gift.ColorBalance(
				float32(a.sliderColorBalanceR.Value), float32(f), float32(a.sliderColorBalanceB.Value)), a.autochange)
		},
		0,
	)

	a.sliderColorBalanceB = widget.NewSlider(-100, 500)
	a.editColorBalanceB = NewEditingOption(
		"Blue: ",
		a.sliderColorBalanceB,
		func(f float64) {
			a.changeParameter(&a.img.cbBlue, gift.ColorBalance(
				float32(a.sliderColorBalanceR.Value), float32(a.sliderColorBalanceG.Value), float32(f)), a.autochange)
		},
		0,
	)

	a.applyBtn = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), a.apply)
	a.resetBtn = widget.NewButtonWithIcon("Reset All", theme.ContentClearIcon(), a.reset)

	// hide all widgets until a file was opened
	// a.scrollEditingWidgets.Content.Hide()
	// a.informationWidgets.Content.Hide()

	// group widgets in a scroll container
	a.scrollEditingWidgets = container.NewScroll(
		widget.NewHBox(
			widget.NewSeparator(),
			widget.NewVBox(
				widget.NewLabelWithStyle("Editor", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("General: "),
				a.editBrightness,
				a.editContrast,
				a.editHue,
				widget.NewSeparator(),
				widget.NewLabel("Color Balance: "),
				a.editColorBalanceR,
				a.editColorBalanceG,
				a.editColorBalanceB,
				widget.NewSeparator(),
				a.resetBtn,
				layout.NewSpacer(),
				a.applyBtn,
			),
		),
	)
	a.scrollEditingWidgets.SetMinSize(fyne.NewSize(130, a.mainWin.Canvas().Size().Height))
	return a.scrollEditingWidgets
}

func (a *App) loadMainUI() fyne.CanvasObject {
	a.mainWin.SetMaster()
	// main menu
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", a.openFile),
			fyne.NewMenuItem("Save", a.saveFile),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Preferences", a.loadSettingsUI),
		),
		fyne.NewMenu("View",
			fyne.NewMenuItem("Information", func() {
				if a.informationWidgets.Visible() {
					a.informationWidgets.Hide()
					a.app.Preferences().SetBool("informationPanelVisible", false)
				} else {
					a.informationWidgets.Show()
					a.app.Preferences().SetBool("informationPanelVisible", true)
				}
			}),
			fyne.NewMenuItem("Editor", func() {
				if a.scrollEditingWidgets.Visible() {
					a.scrollEditingWidgets.Hide()
					a.app.Preferences().SetBool("editorVisible", false)
				} else {
					a.scrollEditingWidgets.Show()
					a.app.Preferences().SetBool("editorVisible", true)
				}
			}),
			fyne.NewMenuItem("Statusbar", func() {
				if a.statusBar.Visible() {
					a.statusBar.Hide()
					a.app.Preferences().SetBool("statusBarVisible", false)
				} else {
					a.statusBar.Show()
					a.app.Preferences().SetBool("statusBarVisible", true)
				}
			}),
			fyne.NewMenuItem("Fullscreen (F11)", func() {

			}),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("Help", func() {

			}),
			fyne.NewMenuItem("About", func() {
				dialog.ShowInformation("About", "A simple image viewer with editing functionality", a.mainWin)
			}),
		),
	)
	a.mainWin.SetMainMenu(mainMenu)

	a.image = &canvas.Image{}
	a.image.FillMode = canvas.ImageFillContain

	layout := container.NewBorder(nil, a.loadStatusBar(), a.loadInformationWidgets(), a.loadEditControls(), a.image)
	a.loadPreferences()
	return layout
}
