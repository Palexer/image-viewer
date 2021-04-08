package main

import (
	"image"
	"os"
	"strconv"

	"github.com/disintegration/gift"
)

// Img is used for the whole image editing process
type Img struct {
	OriginalImage  image.Image
	FileData       os.FileInfo
	EditedImage    *image.RGBA
	gifted         *gifted
	Path           string
	ImagesInFolder []string
	index          int
	Directory      string

	zoom int

	// saved filters
	// general
	brightness gift.Filter
	contrast   gift.Filter
	gamma      gift.Filter
	hue        gift.Filter
	saturation gift.Filter

	// filters
	grayscale gift.Filter
	sepia     gift.Filter
	blur      gift.Filter

	// color balance red, green, blue
	cbRed   gift.Filter
	cbGreen gift.Filter
	cbBlue  gift.Filter

	// transform
	resize gift.Filter

	lastFilters       []gift.Filter
	lastFiltersUndone []gift.Filter
}

func (i *Img) init() {
	i.gifted = &gifted{}
	i.gifted.GIFT = gift.New()
}

func (a *App) changeParameter(filterVar *gift.Filter, newFilter gift.Filter) {
	if a.img.OriginalImage == nil {
		return
	}

	a.img.gifted.replace(*filterVar, newFilter)
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
		a.img.gifted.remove(filterToUndo)
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

func (a *App) zoomImageIn() {
	if a.img.OriginalImage == nil {
		return
	}

	src := a.img.OriginalImage

	my_width := src.Bounds().Dx() - a.img.zoom
	my_height := src.Bounds().Dy() - a.img.zoom

	if a.img.zoom < my_width {
		a.img.zoom += 25
	}
	g := gift.New(
		gift.Crop(image.Rect(a.img.zoom, a.img.zoom, my_width, my_height)),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	// show new image
	a.image.Image = dst
	a.image.Refresh()

	// show zoom level
	a.zoomLabel.SetText(strconv.Itoa(100+a.img.zoom) + "%")
}

func (a *App) zoomImageOut() {
	if a.img.OriginalImage == nil {
		return
	}
	if a.img.zoom > 0 {
		a.img.zoom -= 25

	}
	src := a.img.OriginalImage

	my_width := src.Bounds().Dx() - a.img.zoom
	my_height := src.Bounds().Dy() - a.img.zoom

	g := gift.New(
		gift.Crop(image.Rect(a.img.zoom, a.img.zoom, my_width, my_height)),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	// show new image
	a.image.Image = dst
	a.image.Refresh()

	// show zoom level
	a.zoomLabel.SetText(strconv.Itoa(100+a.img.zoom) + "%")
}

func (a *App) resetZoom() {
	a.img.zoom = 0
	a.zoomLabel.SetText("100%")

	// show new image
	a.image.Image = a.img.OriginalImage
	a.image.Refresh()
}
