package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

// NewEditingOption creates a new VBox, that includes an info text and a widget to edit the paramter
func NewEditingOption(infoText string, editingWidget *widget.Slider, onChanged func(float64), defaultValue float64) *widget.Box {
	editingWidget.SetValue(defaultValue)
	editingWidget.OnChanged = func(f float64) { onChanged(f) }
	vbox := widget.NewVBox(
		widget.NewLabel(infoText),
		editingWidget,
	)
	return vbox
}

// App represents the whole application with all its windows, widgets and functions
type App struct {
	app     fyne.App
	mainWin fyne.Window

	autochange bool

	img Img

	image *canvas.Image

	editBrightness    *widget.Box
	editContrast      *widget.Box
	editHue           *widget.Box
	editColorBalanceR *widget.Box
	editColorBalanceG *widget.Box
	editColorBalanceB *widget.Box

	sliderBrightness    *widget.Slider
	sliderContrast      *widget.Slider
	sliderHue           *widget.Slider
	sliderColorBalanceR *widget.Slider
	sliderColorBalanceG *widget.Slider
	sliderColorBalanceB *widget.Slider

	applyBtn             *widget.Button
	resetBtn             *widget.Button
	scrollEditingWidgets *container.Scroll

	informationWidgets *container.Scroll
	widthLabel         *widget.Label
	heightLabel        *widget.Label
	statusBar          *widget.Box
	imagePathLabel     *widget.Label
}

func main() {
	a := app.NewWithID("io.github.palexer")
	w := a.NewWindow("Image Viewer")
	ui := &App{app: a, mainWin: w}
	w.SetContent(ui.loadMainUI())
	w.Resize(fyne.NewSize(1380, 870))
	w.ShowAndRun()
}
