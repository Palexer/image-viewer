package main

import (
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
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
}

func (a *App) init() {
	a.img = Img{}
	a.img.init()
	a.lastOpened = strings.Split(a.app.Preferences().String("lastOpened"), ",")
}

func main() {
	a := app.NewWithID("io.github.palexer")
	w := a.NewWindow("Image Viewer")
	a.SetIcon(resourceIconPng)
	w.SetIcon(resourceIconPng)
	ui := &App{app: a, mainWin: w}
	ui.init()
	w.SetContent(ui.loadMainUI())
	w.Resize(fyne.NewSize(1380, 870))
	w.ShowAndRun()
}
