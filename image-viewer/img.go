package main

import (
	"image"
	"os"

	"github.com/disintegration/gift"
)

// Img is used for the whole image editing process
type Img struct {
	OriginalImage     image.Image
	OriginalImageData image.Config
	FileData          os.FileInfo
	EditedImage       *image.RGBA
	gifted            *Gifted
	Path              string

	// saved filters
	// general
	brightness gift.Filter
	contrast   gift.Filter
	hue        gift.Filter
	saturation gift.Filter
	// filters
	grayscale gift.Filter
	sepia     gift.Filter
	pixelate  gift.Filter
	blur      gift.Filter
	// color balance red, green, blue
	cbRed   gift.Filter
	cbGreen gift.Filter
	cbBlue  gift.Filter
	// transform
	cropWidth  gift.Filter
	cropHeight gift.Filter

	lastFilters       []gift.Filter
	lastFiltersUndone []gift.Filter
}

func (i *Img) init() {
	i.gifted = &Gifted{}
	i.gifted.GIFT = gift.New()
}

func (a *App) changeParameter(filterVar *gift.Filter, newFilter gift.Filter) {
	if a.img.OriginalImage == nil {
		return
	}
	a.img.gifted.Replace(*filterVar, newFilter)
	*filterVar = newFilter
	a.img.lastFilters = append(a.img.lastFilters, newFilter)
	go a.apply()
}

func (a *App) addParameter(filter gift.Filter) {
	if a.img.OriginalImage == nil {
		return
	}
	a.img.gifted.Add(filter)
	a.img.lastFilters = append(a.img.lastFilters, filter)
	go a.apply()
}

func (a *App) apply() {
	// apply filters
	a.img.EditedImage = image.NewRGBA(a.img.gifted.Bounds(a.img.OriginalImage.Bounds()))
	a.img.gifted.Draw(a.img.EditedImage, a.img.OriginalImage)

	// show new image
	a.image.Image = a.img.EditedImage
	a.image.Refresh()
}

func (a *App) reset() {
	defer a.image.Refresh()

	// reset values
	a.sliderBrightness.SetValue(0)
	a.sliderContrast.SetValue(0)
	a.sliderHue.SetValue(0)
	a.sliderColorBalanceR.SetValue(0)
	a.sliderColorBalanceG.SetValue(0)
	a.sliderColorBalanceB.SetValue(0)
	a.sliderSepia.SetValue(0)
	a.sliderPixelate.SetValue(0)
	a.sliderSaturation.SetValue(0)
	a.sliderBlur.SetValue(0)

	// clear filters
	a.img.gifted.Empty()
	a.img.lastFilters = nil
	a.img.lastFiltersUndone = nil

	a.img.EditedImage = nil
	a.image.Image = a.img.OriginalImage
}

func (a *App) undo() {
	if len(a.img.lastFilters) > 0 {
		filterToUndo := a.img.lastFilters[len(a.img.lastFilters)-1]
		a.img.gifted.Remove(filterToUndo)
		a.img.lastFiltersUndone = append(a.img.lastFiltersUndone, filterToUndo)
		a.img.lastFilters = a.img.lastFilters[:len(a.img.lastFilters)-1]
		a.apply()
	}
}

func (a *App) redo() {
	if len(a.img.lastFiltersUndone) > 0 {
		filterToRedo := a.img.lastFiltersUndone[len(a.img.lastFiltersUndone)-1]
		a.img.gifted.Add(filterToRedo)
		a.img.lastFilters = append(a.img.lastFilters, filterToRedo)
		a.img.lastFiltersUndone = a.img.lastFiltersUndone[:len(a.img.lastFiltersUndone)-1]
		a.apply()
	}
}
