package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/gift"
)

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for i := range elements {
		if !encountered[elements[i]] {
			// Record this element as an encountered element.
			encountered[elements[i]] = true
			// Append to result slice.
			result = append(result, elements[i])
		}
	}
	return result
}

func (a *App) nextImage(forward, folder bool) {
	if a.img.OriginalImage == nil || len(a.img.ImagesInFolder) < 2 {
		return
	}

	if forward {
		if a.img.index == len(a.img.ImagesInFolder)-1 {
			return
		}
		a.img.index++
	} else {
		if a.img.index == 0 {
			return
		}
		a.img.index--
	}

	file, err := os.Open(a.img.Directory + "/" + a.img.ImagesInFolder[a.img.index])
	if err != nil {
		dialog.ShowError(err, a.mainWin)
		return
	}
	if folder {
		a.open(file, true)
	} else {
		a.open(file, false)
	}
}

func (a *App) fullscreenMode() {
	if !a.focus {
		a.fullscreenWin = a.app.NewWindow("Image Viewer - " + a.img.Path)
		a.fullscreenWin.SetFullScreen(true)
		a.fullscreenWin.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			switch key.Name {
			case fyne.KeyEscape:
				a.fullscreenWin.Close()
			case fyne.KeyF11:
				a.fullscreenWin.Close()
			case fyne.KeyRight:
				a.nextImage(true, false)
			case fyne.KeyLeft:
				a.nextImage(false, false)
			}
		})
		a.fullscreenWin.SetContent(a.image)
		a.fullscreenWin.Show()
		a.focus = true
	} else {
		a.fullscreenWin.Hide()
		a.focus = false
	}
}

func (a *App) loadStatusBar() *fyne.Container {
	a.leftArrow = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		a.nextImage(false, false)
	})

	a.rightArrow = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		a.nextImage(true, false)
	})
	a.leftArrow.Disable()
	a.rightArrow.Disable()

	a.deleteBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Delte file?", "Do you really want to delete this image?\n This action can't be undone.", func(b bool) {
			if b {
				a.deleteFile()
			}
		}, a.mainWin)
	})
	a.deleteBtn.Disable()

	a.renameBtn = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), a.renameDialog)
	a.renameBtn.Disable()

	a.zoomIn = widget.NewButtonWithIcon("", theme.ZoomInIcon(), a.zoomImageIn)
	a.zoomIn.Disable()

	a.zoomOut = widget.NewButtonWithIcon("", theme.ZoomOutIcon(), a.zoomImageOut)
	a.zoomOut.Disable()

	a.zoomLabel = widget.NewLabel("")

	a.resetZoomBtn = widget.NewButtonWithIcon("", theme.ZoomFitIcon(), a.resetZoom)
	a.resetZoomBtn.Disable()

	a.statusBar = container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			a.leftArrow,
			a.rightArrow,
			layout.NewSpacer(),
			a.zoomLabel,
			a.resetZoomBtn,
			a.zoomOut,
			a.zoomIn,
			a.renameBtn,
			a.deleteBtn,
		),
	)
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

	rotate90Btn := widget.NewButton("Rotate 90°", func() { a.addParameter(gift.Rotate90()) })
	flipVerticalBtn := widget.NewButton("Flip Vertically", func() { a.addParameter(gift.FlipVertical()) })
	flipHorizontalBtn := widget.NewButton("Flip Horizontally", func() { a.addParameter(gift.FlipHorizontal()) })
	resizeBtn := widget.NewButton("Resize", func() {
		var keepAspectRatio bool

		widthEntry := widget.NewEntry()
		heightEntry := widget.NewEntry()
		widthEntry.SetPlaceHolder("Width")
		heightEntry.SetPlaceHolder("Height")
		widthEntry.Validator = func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return fmt.Errorf("input is not a valid number")
			}
			return nil
		}
		heightEntry.Validator = func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return fmt.Errorf("input is not a valid number")
			}
			return nil
		}

		dialog.ShowCustom("Resize", "Cancel", container.NewVBox(
			container.NewHBox(
				widthEntry,
				heightEntry,
			),
			container.NewHBox(
				widget.NewCheck("Keep Aspect Ratio", func(b bool) {
					if b {
						keepAspectRatio = true
					} else {
						keepAspectRatio = false
					}
				}),
				widget.NewButton("Apply", func() {
					if err := widthEntry.Validate(); err != nil {
						dialog.ShowError(err, a.mainWin)
						return
					}
					if err := widthEntry.Validate(); err != nil {
						dialog.ShowError(err, a.mainWin)
						return
					}

					width, _ := strconv.Atoi(widthEntry.Text)
					height, _ := strconv.Atoi(heightEntry.Text)
					if keepAspectRatio {
						a.changeParameter(&a.img.resize, gift.ResizeToFit(width, height, gift.LinearResampling))
					} else {
						a.changeParameter(&a.img.resize, gift.ResizeToFill(width, height, gift.LinearResampling, gift.BottomAnchor))
					}
					a.mainWin.Canvas().Overlays().Top().Hide()
				}),
			),
		), a.mainWin)
	})

	grayscaleBtn := widget.NewButton("Grayscale", func() { a.changeParameter(&a.img.grayscale, gift.Grayscale()) })

	a.sliderSepia = newEditingSlider(0, 100)
	a.sliderSepia.dragEndFunc = func(f float64) { a.changeParameter(&a.img.sepia, gift.Sepia(float32(f))) }
	editSepia := newEditingOption("Sepia: ", a.sliderSepia, 0)

	a.sliderBlur = newEditingSlider(0, 100)
	a.sliderBlur.dragEndFunc = func(f float64) { a.changeParameter(&a.img.blur, gift.GaussianBlur(float32(f))) }
	editBlur := newEditingOption("Blur: ", a.sliderBlur, 0)

	a.resetBtn = widget.NewButtonWithIcon("Reset All", theme.ContentClearIcon(), a.reset)
	a.resetBtn.Disable()

	return container.NewTabItem("Editor", container.NewScroll(
		container.NewVBox(
			widget.NewAccordion(
				widget.NewAccordionItem(
					"General",
					container.NewVBox(
						editBrightness,
						editContrast,
						editHue,
						editSaturation,
					),
				),
				widget.NewAccordionItem(
					"Color Balance",
					container.NewVBox(
						editColorBalanceR,
						editColorBalanceG,
						editColorBalanceB,
					),
				),
				widget.NewAccordionItem(
					"Transform",
					container.NewVBox(
						rotate90Btn,
						flipHorizontalBtn,
						flipVerticalBtn,
						resizeBtn,
					),
				),
				widget.NewAccordionItem(
					"Filter",
					container.NewVBox(
						editSepia,
						editBlur,
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
	return container.NewTabItem("Information", container.NewScroll(
		container.NewVBox(
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
			fyne.NewMenuItem("Save As", a.saveFileDialog),
			// recent,
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Undo", a.undo),
			fyne.NewMenuItem("Redo", a.redo),
			fyne.NewMenuItem("Delete Image", a.deleteFile),
			fyne.NewMenuItem("Keyboard Shortucts", a.showShortcuts),
			fyne.NewMenuItem("Preferences", a.loadSettingsUI),
		),
		fyne.NewMenu("View",
			fyne.NewMenuItem("Fullscreen", func() {
				if a.image.Image == nil {
					return
				}
				a.fullscreenMode()
			}),
			fyne.NewMenuItem("Next Image", func() {
				a.nextImage(true, false)
			}),
			fyne.NewMenuItem("Last Image", func() {
				a.nextImage(false, false)
			}),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("About", func() {
				dialog.ShowCustom("About", "Ok", container.NewVBox(
					widget.NewLabel("A simple image viewer with some editing functionality."),
					widget.NewHyperlink("Help and more information on Github", parseURL("https://github.com/Palexer/image-viewer")),
					widget.NewLabel("v1.2 | License: MIT"),
				), a.mainWin)
			}),
		),
	)
	a.mainWin.SetMainMenu(mainMenu)
	a.loadKeyboardShortcuts()

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
	return layout
}
