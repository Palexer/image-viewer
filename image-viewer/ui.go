package main

import (
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/disintegration/gift"
)

// StringToInt converts a string to an integer
func StringToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

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
	a.sliderBrightness = newEditingSlider(-100, 100)
	a.sliderBrightness.dragEndFunc = func(f float64) { a.changeParameter(&a.img.brightness, gift.Brightness(float32(f))) }
	a.editBrightness = newEditingOption(
		"Brightness",
		a.sliderBrightness,
		0,
	)

	a.sliderContrast = newEditingSlider(-100, 100)
	a.sliderContrast.dragEndFunc = func(f float64) { a.changeParameter(&a.img.contrast, gift.Contrast(float32(f))) }
	a.editContrast = newEditingOption(
		"Contrast",
		a.sliderContrast,
		0,
	)

	a.sliderHue = newEditingSlider(-180, 180)
	a.sliderHue.dragEndFunc = func(f float64) { a.changeParameter(&a.img.hue, gift.Hue(float32(f))) }
	a.editHue = newEditingOption(
		"Hue",
		a.sliderHue,
		0,
	)

	a.sliderColorBalanceR = newEditingSlider(-100, 500)
	a.sliderColorBalanceR.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbRed, gift.ColorBalance(
			float32(f), float32(a.sliderColorBalanceG.Value), float32(a.sliderColorBalanceB.Value)))
	}
	a.editColorBalanceR = newEditingOption(
		"Red",
		a.sliderColorBalanceR,
		0,
	)

	a.sliderColorBalanceG = newEditingSlider(-100, 500)
	a.sliderColorBalanceG.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbGreen, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(f), float32(a.sliderColorBalanceB.Value)))
	}
	a.editColorBalanceG = newEditingOption(
		"Green",
		a.sliderColorBalanceG,
		0,
	)

	a.sliderColorBalanceB = newEditingSlider(-100, 500)
	a.sliderColorBalanceB.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbBlue, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(a.sliderColorBalanceG.Value), float32(f)))
	}
	a.editColorBalanceB = newEditingOption(
		"Blue",
		a.sliderColorBalanceB,
		0,
	)

	cropWidth := widget.NewEntry()
	cropWidth.SetPlaceHolder("Width: " + strconv.Itoa(a.img.OriginalImageData.Width))

	cropHeight := widget.NewEntry()
	cropHeight.SetPlaceHolder("Height: " + strconv.Itoa(a.img.OriginalImageData.Height))

	cropWidth.OnChanged = func(s string) {
		var width, height int
		if cropHeight.Text != "" {
			width, _ = StringToInt(cropHeight.Text)
		} else {
			width = a.img.OriginalImageData.Height
		}
		height, _ = StringToInt(s)
		a.changeParameter(&a.img.cropWidth, gift.CropToSize(height, width, gift.CenterAnchor))
	}

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
				widget.NewLabel("General"),
				a.editBrightness,
				a.editContrast,
				a.editHue,
				widget.NewSeparator(),
				widget.NewLabel("Color Balance"),
				a.editColorBalanceR,
				a.editColorBalanceG,
				a.editColorBalanceB,
				widget.NewSeparator(),
				widget.NewLabel("Transform"),
				cropWidth,
				cropHeight,
				layout.NewSpacer(),
				a.resetBtn,
			),
		),
	)
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
			fyne.NewMenuItem("Undo", a.undo),
			fyne.NewMenuItem("Redo", a.redo),
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
			fyne.NewMenuItem("Fullscreen (Ctrl+F)", a.showFullscreen),
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

	// keyboard shortcuts
	// ctrl+f for fullscreen
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.showFullscreen() })

	// ctrl+o to open file
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.openFile() })

	// ctrl+z to undo
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyZ,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.undo() })

	// ctrl+y to redo
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyY,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.redo() })

	a.image = &canvas.Image{}
	a.image.FillMode = canvas.ImageFillContain

	layout := container.NewBorder(nil, a.loadStatusBar(), a.loadInformationWidgets(), a.loadEditControls(), a.image)
	// container.NewHSplit()
	a.loadPreferences()
	return layout
}

func (a *App) showFullscreen() {
	fullWin := a.app.NewWindow("Fullscreen Image")
	fullWin.SetContent(a.image)
	fullWin.SetFullScreen(true)
	fullWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF11,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { fullWin.Close(); a.mainWin.Content().Refresh() })
	fullWin.Show()
}
