package main

import (
	"runtime"
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

// loadEditorTab returns the editor tab
func (a *App) loadEditorTab() *container.TabItem {
	a.sliderBrightness = newEditingSlider(-100, 100)
	a.sliderBrightness.dragEndFunc = func(f float64) { a.changeParameter(&a.img.brightness, gift.Brightness(float32(f))) }
	editBrightness := newEditingOption(
		"Brightness: ",
		a.sliderBrightness,
		0,
	)

	a.sliderContrast = newEditingSlider(-100, 100)
	a.sliderContrast.dragEndFunc = func(f float64) { a.changeParameter(&a.img.contrast, gift.Contrast(float32(f))) }
	editContrast := newEditingOption(
		"Contrast: ",
		a.sliderContrast,
		0,
	)

	a.sliderHue = newEditingSlider(-180, 180)
	a.sliderHue.dragEndFunc = func(f float64) { a.changeParameter(&a.img.hue, gift.Hue(float32(f))) }
	editHue := newEditingOption(
		"Hue: ",
		a.sliderHue,
		0,
	)

	a.sliderSaturation = newEditingSlider(-100, 500)
	a.sliderSaturation.dragEndFunc = func(f float64) { a.changeParameter(&a.img.saturation, gift.Saturation(float32(f))) }
	editSaturation := newEditingOption("Saturation: ", a.sliderSaturation, 0)

	a.sliderColorBalanceR = newEditingSlider(-100, 500)
	a.sliderColorBalanceR.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbRed, gift.ColorBalance(
			float32(f), float32(a.sliderColorBalanceG.Value), float32(a.sliderColorBalanceB.Value)))
	}
	editColorBalanceR := newEditingOption(
		"Red: ",
		a.sliderColorBalanceR,
		0,
	)

	a.sliderColorBalanceG = newEditingSlider(-100, 500)
	a.sliderColorBalanceG.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbGreen, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(f), float32(a.sliderColorBalanceB.Value)))
	}
	editColorBalanceG := newEditingOption(
		"Green: ",
		a.sliderColorBalanceG,
		0,
	)

	a.sliderColorBalanceB = newEditingSlider(-100, 500)
	a.sliderColorBalanceB.dragEndFunc = func(f float64) {
		a.changeParameter(&a.img.cbBlue, gift.ColorBalance(
			float32(a.sliderColorBalanceR.Value), float32(a.sliderColorBalanceG.Value), float32(f)))
	}
	editColorBalanceB := newEditingOption(
		"Blue: ",
		a.sliderColorBalanceB,
		0,
	)

	cropWidth := widget.NewEntry()
	cropWidth.SetPlaceHolder("Width: " + strconv.Itoa(a.img.OriginalImageData.Width))

	cropHeight := widget.NewEntry()
	cropHeight.SetPlaceHolder("Height: " + strconv.Itoa(a.img.OriginalImageData.Height))

	grayscaleBtn := widget.NewButton("Grayscale", func() { a.changeParameter(&a.img.grayscale, gift.Grayscale()) })

	a.sliderSepia = newEditingSlider(0, 100)
	a.sliderSepia.dragEndFunc = func(f float64) { a.changeParameter(&a.img.sepia, gift.Sepia(float32(f))) }
	editSepia := newEditingOption("Sepia: ", a.sliderSepia, 0)

	a.sliderPixelate = newEditingSlider(0, 100)
	a.sliderPixelate.dragEndFunc = func(f float64) { a.changeParameter(&a.img.pixelate, gift.Pixelate(int(f))) }
	editPixelate := newEditingOption("Pixelate: ", a.sliderPixelate, 0)

	a.sliderBlur = newEditingSlider(0, 100)
	a.sliderBlur.dragEndFunc = func(f float64) { a.changeParameter(&a.img.blur, gift.GaussianBlur(float32(f))) }
	editBlur := newEditingOption("Blur: ", a.sliderBlur, 0)

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
			widget.NewAccordion(
				widget.NewAccordionItem(
					"General",
					widget.NewVBox(
						editBrightness,
						editContrast,
						editHue,
						editSaturation,
					),
				),
				widget.NewAccordionItem(
					"Color Balance",
					widget.NewVBox(
						editColorBalanceR,
						editColorBalanceG,
						editColorBalanceB,
					),
				),
				widget.NewAccordionItem(
					"Transform",
					widget.NewVBox(
						cropWidth,
						cropHeight,
					),
				),
				widget.NewAccordionItem(
					"Filter",
					widget.NewVBox(
						editSepia,
						editBlur,
						editPixelate,
						grayscaleBtn,
					),
				),
			),
			layout.NewSpacer(),
			a.resetBtn,
		)))
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
	// set main mod key to super on darwin hosts, else set it to ctrl
	if runtime.GOOS == "darwin" {
		a.mainModKey = desktop.SuperModifier
	} else {
		a.mainModKey = desktop.ControlModifier
	}
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
			fyne.NewMenuItem("About", func() {
				dialog.ShowCustom("About", "Ok", widget.NewVBox(
					widget.NewLabel("A simple image viewer with some editing functionality."),
					widget.NewHyperlink("Help and more information on Github", parseURL("https://github.com/Palexer/image-viewer")),
				), a.mainWin)
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

	// ctrl+q to quit application
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.app.Quit() })

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
	a.split.SetOffset(0.90)
	layout := container.NewBorder(nil, a.loadStatusBar(), nil, nil, a.split)
	a.loadPreferences()
	return layout
}

func (a *App) focusMode() {
	if !a.focus {
		a.statusBar.Hide()
		a.split.Hide()
		a.mainWin.SetContent(fyne.NewContainer(a.image))
        a.focus = true
	} else {
        a.statusBar.Show()
		a.split.Show()
		a.mainWin.SetContent(container.NewBorder(nil, a.statusBar, nil, nil, a.split))
        a.focus = false
	}
}
