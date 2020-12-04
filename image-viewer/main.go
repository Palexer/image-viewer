package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
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

// App represents the whole application with all its windows, widgets and functions
type App struct {
	app     fyne.App
	mainWin fyne.Window

	img Img

	image *canvas.Image

	editBrightness    *widget.Box
	editContrast      *widget.Box
	editHue           *widget.Box
	editColorBalanceR *widget.Box
	editColorBalanceG *widget.Box
	editColorBalanceB *widget.Box

	sliderBrightness    *editingSlider
	sliderContrast      *editingSlider
	sliderHue           *editingSlider
	sliderColorBalanceR *editingSlider
	sliderColorBalanceG *editingSlider
	sliderColorBalanceB *editingSlider

	resetBtn             *widget.Button
	scrollEditingWidgets *container.Scroll

	informationWidgets *container.Scroll
	widthLabel         *widget.Label
	heightLabel        *widget.Label
	statusBar          *widget.Box
	imagePathLabel     *widget.Label
}

// newEditingOption creates a new VBox, that includes an info text and a widget to edit the paramter
func newEditingOption(infoText string, slider *editingSlider, defaultValue float64) *widget.Box {
	slider.SetValue(defaultValue)
	vbox := widget.NewVBox(
		widget.NewLabel(infoText),
		slider,
	)
	return vbox
}

func main() {
	a := app.NewWithID("io.github.palexer")
	w := a.NewWindow("Image Viewer")
	ui := &App{app: a, mainWin: w}
	w.SetContent(ui.loadMainUI())
	w.Resize(fyne.NewSize(1380, 870))
	w.ShowAndRun()
}
