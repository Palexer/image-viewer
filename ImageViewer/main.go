package main

import (
	"fmt"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
func newEditingOption(infoText string, slider *editingSlider, defaultValue float64) *fyne.Container {
	data := binding.BindFloat(&defaultValue)
	text := widget.NewLabel(infoText)
	value := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "%.0f"))
	slider.Bind(data)
	slider.Step = 1

	return container.NewVBox(
		container.NewHBox(
			text,
			layout.NewSpacer(),
			value,
		),
		slider,
	)
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
	resetBtn            *widget.Button

	split       *container.Split
	widthLabel  *widget.Label
	heightLabel *widget.Label
	imgSize     *widget.Label
	imgLastMod  *widget.Label

	statusBar      *fyne.Container
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
	w.Resize(fyne.NewSize(1200, 750))
	w.ShowAndRun()
}
