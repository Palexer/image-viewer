package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

// App represents the whole application with all its windows, widgets and functions
type App struct {
	app     fyne.App
	mainWin fyne.Window

	autochange bool

	img Img

	image *canvas.Image

	sliderBrightness *widget.Slider
	sliderContrast   *widget.Slider
	applyBtn         *widget.Button
	resetBtn         *widget.Button

	editControls       *widget.Box
	informationWidgets *widget.Box
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
