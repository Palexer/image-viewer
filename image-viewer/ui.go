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

func (a *App) loadInformationWidgets() *widget.Box {
	a.widthLabel = widget.NewLabel("Width: ")
	a.heightLabel = widget.NewLabel("Height: ")
	a.informationWidgets = widget.NewHBox(
		widget.NewVBox(
			widget.NewLabelWithStyle("Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			a.widthLabel,
			a.heightLabel,
		),
		widget.NewSeparator(),
	)
	return a.informationWidgets
}

func (a *App) loadEditControls() *widget.Box {
	a.sliderBrightness = widget.NewSlider(-100, 100)
	a.sliderBrightness.SetValue(0)

	a.sliderContrast = widget.NewSlider(-100, 100)
	a.sliderContrast.SetValue(0)

	a.sliderBrightness.OnChanged = func(f float64) {
		a.changeParameter(&a.img.brightness, gift.Brightness(float32(f)), a.autochange)
	}
	a.sliderContrast.OnChanged = func(f float64) {
		a.changeParameter(&a.img.contrast, gift.Contrast(float32(f)), a.autochange)
	}

	a.applyBtn = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), a.apply)
	a.resetBtn = widget.NewButtonWithIcon("Reset All", theme.ContentClearIcon(), a.reset)

	// disable widgets until a file was opened
	a.resetBtn.Disable()
	a.applyBtn.Disable()
	a.sliderBrightness.Hide()
	a.sliderContrast.Hide()

	a.editControls = widget.NewHBox(
		widget.NewSeparator(),
		widget.NewVBox(
			widget.NewLabelWithStyle("Editor", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Brightness: "),
			a.sliderBrightness,
			widget.NewLabel("Contrast: "),
			a.sliderContrast,
			a.resetBtn,
			layout.NewSpacer(),
			a.applyBtn,
		),
	)
	return a.editControls
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
				if a.editControls.Visible() {
					a.editControls.Hide()
					a.app.Preferences().SetBool("editorVisible", false)
				} else {
					a.editControls.Show()
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
