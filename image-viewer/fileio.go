package main

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

func (a *App) openFileDialog() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}
		if err == nil && reader == nil {
			return
		}

		err = a.open(reader)
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}
		defer reader.Close()
	}, a.mainWin)
}

func (a *App) open(f fyne.URIReadCloser) error {
	defer f.Close()
	if f == nil {
		return errors.New("cancelled")
	}

	// init Img
	a.img = Img{}
	a.img.init()

	// decode and update the image + get image path
	var err error
	a.img.OriginalImage, _, err = image.Decode(f)
	if err != nil {
		return fmt.Errorf("Unable to decode image %v", err)
	}
	a.img.Path = f.URI().String()[7:]
	a.image.Image = a.img.OriginalImage
	a.image.Refresh()

	// get width and height of the image
	reader, err := os.Open(a.img.Path)
	a.img.OriginalImageData, _, _ = image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf("Unable to get image information %v", err)
	}

	// get and display FileInfo
	a.img.FileData, err = os.Stat(a.img.Path)
	a.imgSize.SetText(fmt.Sprintf("Size: %.2f Mb", float64(a.img.FileData.Size())/1000000))

	modtime := a.img.FileData.ModTime()
	a.imgLastMod.SetText(fmt.Sprintf("Last modified: \n%v-%v-%v", modtime.Year(), int(modtime.Month()), modtime.Day()))

	a.imagePathLabel.SetText("Path: " + a.img.Path)
	a.widthLabel.SetText(fmt.Sprintf("Width:   %dpx", a.img.OriginalImageData.Width))
	a.heightLabel.SetText(fmt.Sprintf("Height: %dpx", a.img.OriginalImageData.Height))

	a.mainWin.SetTitle(fmt.Sprintf("Image Viewer - %v", (strings.Split(a.img.Path, "/")[len(strings.Split(a.img.Path, "/"))-1])))

	// activate widgets
	a.reset()
	a.resetBtn.Enable()
	return nil
}

func (a *App) saveFileDialog() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		err = a.save(writer)
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}
	}, a.mainWin)
}

func (a *App) save(writer fyne.URIWriteCloser) error {
	if writer == nil {
		return nil
	}
	switch writer.URI().Extension() {
	case ".jpeg":
		jpeg.Encode(writer, a.img.EditedImage, nil)
	case ".png":
		png.Encode(writer, a.img.EditedImage)
	case ".gif":
		gif.Encode(writer, a.img.EditedImage, nil)
	default:
		os.Remove(writer.URI().String()[7:])
		return errors.New("unsupported file extension\n supported extensions: .jpeg, .png, .gif")
	}
	return nil
}
