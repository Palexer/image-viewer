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
	Filter            *gift.GIFT
	Path              string
}

func (i *Img) init() {
	i.Filter = gift.New()
}

// func (i *Img) changeBrightness(value float64) {
// 	i.Filter.Add(gift.Brightness(float32(value)))
// }

// func (i *Img) changeContrast(percentage float64) {
// 	i.Filter.Add(gift.Contrast(float32(percentage)))
// }

func (a *App) apply() {
	// clear all filters
	a.img.Filter.Empty()

	// add possible filters
	a.img.Filter.Add(gift.Brightness(float32(a.sliderBrightness.Value)))
	a.img.Filter.Add(gift.Contrast(float32(a.sliderContrast.Value)))

	// apply filters
	a.img.EditedImage = image.NewRGBA(a.img.Filter.Bounds(a.img.OriginalImage.Bounds()))
	a.img.Filter.Draw(a.img.EditedImage, a.img.OriginalImage)

	// show new image
	a.image.Image = a.img.EditedImage
	defer a.image.Refresh()
}

func (a *App) reset() {
	a.sliderBrightness.SetValue(0)
	a.sliderContrast.SetValue(0)
	a.img.Filter.Empty()
	a.img.EditedImage = nil
	a.image.Image = a.img.OriginalImage
	defer a.image.Refresh()
}
