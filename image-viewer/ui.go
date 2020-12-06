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

func (a *App) loadEditorTab() *container.TabItem {
	// editor tab
	a.sliderBrightness = newEditingSlider(-100, 100)
	a.sliderBrightness.dragEndFunc = func(f float64) { a.changeParameter(&a.img.brightness, gift.Brightness(float32(f))) }
	a.editBrightness = newEditingOption(
		"Brightness: ",
		a.sliderBrightness,
		0,
	)

	a.sliderContrast = newEditingSlider(-100, 100)
	a.sliderContrast.dragEndFunc = func(f float64) { a.changeParameter(&a.img.contrast, gift.Contrast(float32(f))) }
	a.editContrast = newEditingOption(
		"Contrast: ",
		a.sliderContrast,
		0,
	)

	a.sliderHue = newEditingSlider(-180, 180)
	a.sliderHue.dragEndFunc = func(f float64) { a.changeParameter(&a.img.hue, gift.Hue(float32(f))) }
	a.editHue = newEditingOption(
		"Hue: ",
		a.sliderHue,
		0,
	)

	a.sliderColorBalanceR = newEditingSlider(-100, 500)
	a.sliderColorBalanceR.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbRed, gift.ColorBalance(
			float32(f), float32(a.sliderColorBalanceG.Value), float32(a.sliderColorBalanceB.Value)))
	}
	a.editColorBalanceR = newEditingOption(
		"Red: ",
		a.sliderColorBalanceR,
		0,
	)

	a.sliderColorBalanceG = newEditingSlider(-100, 500)
	a.sliderColorBalanceG.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbGreen, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(f), float32(a.sliderColorBalanceB.Value)))
	}
	a.editColorBalanceG = newEditingOption(
		"Green: ",
		a.sliderColorBalanceG,
		0,
	)

	a.sliderColorBalanceB = newEditingSlider(-100, 500)
	a.sliderColorBalanceB.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbBlue, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(a.sliderColorBalanceG.Value), float32(f)))
	}
	a.editColorBalanceB = newEditingOption(
		"Blue: ",
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
	a.resetBtn.Disable()
	return container.NewTabItem("Editor", container.NewScroll(
		widget.NewVBox(
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
}

func (a *App) loadInformationTab() *container.TabItem {
	a.widthLabel = widget.NewLabel("Width: ")
	a.heightLabel = widget.NewLabel("Height: ")
	a.imgSize = widget.NewLabel("Size: ")
	a.imgLastMod = widget.NewLabel("Last modified: ")
	// a.informationWidgets.SetMinSize(fyne.NewSize(150, a.mainWin.Canvas().Size().Height))
	return container.NewTabItem("Information", container.NewScroll(
		widget.NewVBox(
			a.widthLabel,
			a.heightLabel,
			a.imgSize,
			a.imgLastMod,
		),
	))
}

func (a *App) loadMainUI() fyne.CanvasObject {
	a.mainWin.SetMaster()
	// main menu
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open", a.openFileDialog),
			fyne.NewMenuItem("Save", a.saveFileDialog),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Undo", a.undo),
			fyne.NewMenuItem("Redo", a.redo),
			fyne.NewMenuItem("Preferences", a.loadSettingsUI),
		),
		fyne.NewMenu("View",
			fyne.NewMenuItem("Focus Mode (Ctrl+F)", a.focusMode),
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
	// ctrl+f for focus mode
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.focusMode() })

	// ctrl+o to open file
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.openFileDialog() })

	// ctrl+s to save file
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.saveFileDialog() })

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

	// image canvas
	a.image = &canvas.Image{}
	a.image.FillMode = canvas.ImageFillContain

	a.split = container.NewHSplit(
		a.image,
		container.NewAppTabs(
			a.loadInformationTab(),
			a.loadEditorTab(),
		),
	)
	a.split.SetOffset(0.85)
	layout := container.NewBorder(nil, a.loadStatusBar(), nil, nil, a.split)
	a.loadPreferences()
	return layout
}

func (a *App) focusMode() {
	if !a.focus {
		a.statusBar.Hide()
		a.split.Hide()
		a.mainWin.SetContent(fyne.NewContainer(a.image))
	} else {
		a.mainWin.SetContent(container.NewBorder(nil, a.statusBar, nil, nil, a.split))
		a.statusBar.Show()
		a.split.Show()

	}
}
