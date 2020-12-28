package main

import (
	"fmt"
	"net/url"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type editingSlider struct {
	widget.Slider
	dragEndFunc func(float64)
}

func (e *editingSlider) DragEnd() {
	e.dragEndFunc(e.Value)
}

func newEditingSlider(min, max float64) *editingSlider {
	editSlider := &editingSlider{}
	editSlider.Max = max
	editSlider.Min = min
	editSlider.ExtendBaseWidget(editSlider)
	return editSlider
}

// newEditingOption creates a new VBox, that includes an info text and a widget to edit the paramter
func newEditingOption(infoText string, slider *editingSlider, defaultValue float64) *widget.Box {
	label := widget.NewLabel(fmt.Sprintf("%v %2.f", infoText, defaultValue))
	slider.SetValue(defaultValue)
	slider.OnChanged = func(f float64) { label.SetText(fmt.Sprintf("%v%2.f", infoText, slider.Value)) }
	vbox := widget.NewVBox(
		label,
		slider,
	)
	return vbox
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	return link
}

// App represents the whole application with all its windows, widgets and functions
type App struct {
	app     fyne.App
	mainWin fyne.Window

	img        Img
	mainModKey desktop.Modifier
	focus      bool
	lastOpened []string

	image *canvas.Image

	sliderBrightness    *editingSlider
	sliderContrast      *editingSlider
	sliderHue           *editingSlider
	sliderSaturation    *editingSlider
	sliderColorBalanceR *editingSlider
	sliderColorBalanceG *editingSlider
	sliderColorBalanceB *editingSlider
	sliderSepia         *editingSlider
	sliderBlur          *editingSlider

	resetBtn *widget.Button

	split          *widget.SplitContainer
	widthLabel     *widget.Label
	heightLabel    *widget.Label
	imgSize        *widget.Label
	imgLastMod     *widget.Label
	statusBar      *widget.Box
	imagePathLabel *widget.Label
	leftArrow      *widget.Button
	rightArrow     *widget.Button
	deleteBtn      *widget.Button
}

func reverseArray(arr []string) []string {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func (a *App) init() {
	a.img = Img{}
	a.img.init()

	// theme
	switch a.app.Preferences().StringWithFallback("Theme", "Dark") {
	case "Light":
		a.app.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.app.Settings().SetTheme(theme.DarkTheme())
	}

	// show/hide statusbar
	if a.app.Preferences().BoolWithFallback("statusBarVisible", true) == false {
		a.statusBar.Hide()
	}
}

// func (a *App) loadRecent() []fyne.URI {
// 	a.lastOpened = strings.Split(a.app.Preferences().String("lastOpened"), ",")
// 	a.lastOpened = reverseArray(a.lastOpened)
// 	// max. 5 items
// 	if len(a.lastOpened) > 5 {
// 		a.lastOpened = a.lastOpened[:5]
// 	}
// 	// remove dublicates
// 	a.lastOpened = removeDuplicates(a.lastOpened)

// 	recent := []fyne.URI{}
// 	for _, v := range a.lastOpened {
// 		recent = append(recent, storage.NewURI(fyne.CurrentApp().Preferences().String(v)))
// 	}
// 	return recent
// }

func main() {
	a := app.NewWithID("io.github.palexer.image-viewer")
	w := a.NewWindow("Image Viewer")
	a.SetIcon(resourceIconPng)
	w.SetIcon(resourceIconPng)
	ui := &App{app: a, mainWin: w}
	ui.init()
	w.SetContent(ui.loadMainUI())
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("error while opening the file: %v\n", err)
		}
		ui.open(file, true)
	}
	w.Resize(fyne.NewSize(1380, 870))
	w.ShowAndRun()
}
