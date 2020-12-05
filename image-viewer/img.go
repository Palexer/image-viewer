package main

import (
	"image"

	"github.com/disintegration/gift"
)

// Img is used for the whole image editing process
type Img struct {
	OriginalImage     image.Image
	OriginalImageData image.Config
	EditedImage       *image.RGBA
	gifted            *Gifted
	Path              string

	// saved filters
	brightness gift.Filter
	contrast   gift.Filter
	hue        gift.Filter
	// color balance red, green, blue
	cbRed   gift.Filter
	cbGreen gift.Filter
	cbBlue  gift.Filter
	// transform
	cropWidth  gift.Filter
	cropHeight gift.Filter

	lastFilters []gift.Filter
}

func (i *Img) init() {
	i.gifted = &Gifted{}
	i.gifted.GIFT = gift.New()
}

func (a *App) changeParameter(filterVar *gift.Filter, newFilter gift.Filter) {
	a.img.gifted.Replace(*filterVar, newFilter)
	*filterVar = newFilter
	a.img.lastFilters = append(a.img.lastFilters, newFilter)
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
	a.sliderBrightness.SetValue(0)
	a.sliderContrast.SetValue(0)
	a.sliderHue.SetValue(0)
	a.sliderColorBalanceR.SetValue(0)
	a.sliderColorBalanceG.SetValue(0)
	a.sliderColorBalanceB.SetValue(0)

	a.img.gifted.Empty()
	a.img.EditedImage = nil
	a.image.Image = a.img.OriginalImage
}

func (a *App) undo() {
	if len(a.img.lastFilters) > 0 {
		a.img.gifted.Remove(a.img.lastFilters[len(a.img.lastFilters)-1])
		a.img.lastFilters = a.img.lastFilters[:len(a.img.lastFilters)-1]
		a.apply()
	}
}
