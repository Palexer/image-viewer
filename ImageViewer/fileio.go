package main

import (
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/disintegration/imageorient"
)

func (a *App) openFileDialog() {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}
		if err == nil && reader == nil {
			return
		}

		file, err := os.Open(reader.URI().String()[7:])
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}

		err = a.open(file, true)
		if err != nil {
			dialog.ShowError(err, a.mainWin)
			return
		}
		defer reader.Close()
	}, a.mainWin)
	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpeg", ".jpg", ".gif"}))
	dialog.Show()
}

func (a *App) open(file *os.File, folder bool) error {
	defer file.Close()

	// decode and update the image + get image path
	var err error
	a.img.OriginalImage, _, err = imageorient.Decode(file)
	if err != nil {
		return fmt.Errorf("Unable to decode image %v", err)
	}
	a.img.Path = file.Name()
	a.image.Image = a.img.OriginalImage
	a.image.Refresh()

	// get and display FileInfo
	a.img.FileData, err = os.Stat(a.img.Path)
	a.imgSize.SetText(fmt.Sprintf("Size: %.2f Mb", float64(a.img.FileData.Size())/1000000))

	a.imgLastMod.SetText(fmt.Sprintf("Last modified: \n%s", a.img.FileData.ModTime().Format("02-01-2006")))

	a.imagePathLabel.SetText("Path: " + a.img.Path)

	// save all images from folder for next/back
	if folder {
		a.img.Directory = filepath.Dir(file.Name())
		openFolder, _ := os.Open(a.img.Directory)
		a.img.ImagesInFolder, _ = openFolder.Readdirnames(0)

		// filter image files
		imgList := []string{}
		for _, v := range a.img.ImagesInFolder {
			if strings.HasSuffix(strings.ToLower(v), ".png") || strings.HasSuffix(strings.ToLower(v), ".jpg") || strings.HasSuffix(strings.ToLower(v), ".jpeg") || strings.HasSuffix(strings.ToLower(v), ".gif") {
				imgList = append(imgList, v)
			}
		}
		a.img.ImagesInFolder = imgList
		sort.Strings(a.img.ImagesInFolder) // sort array alphabetically

		// get first index value
		for i, v := range a.img.ImagesInFolder {
			if filepath.Base(file.Name()) == v {
				a.img.index = i
			}
		}
	}

	a.widthLabel.SetText(fmt.Sprintf("Width:   %dpx", a.img.OriginalImage.Bounds().Max.X))
	a.heightLabel.SetText(fmt.Sprintf("Height: %dpx", a.img.OriginalImage.Bounds().Max.Y))

	a.mainWin.SetTitle(fmt.Sprintf("Image Viewer - %v", (strings.Split(a.img.Path, "/")[len(strings.Split(a.img.Path, "/"))-1])))

	// append to last opened images
	a.lastOpened = append(a.lastOpened, file.Name())
	a.app.Preferences().SetString("lastOpened", strings.Join(a.lastOpened, ","))

	// reset editing history
	a.img.lastFilters = nil
	a.img.lastFiltersUndone = nil

	// activate widgets
	a.reset()
	a.resetBtn.Enable()
	a.leftArrow.Enable()
	a.rightArrow.Enable()
	a.deleteBtn.Enable()
	return nil
}

func (a *App) saveFileDialog() {
	if a.img.OriginalImage == nil {
		dialog.ShowError(errors.New("no image opened"), a.mainWin)
		return
	}
	if a.img.EditedImage == nil {
		a.apply()
	}

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
	case ".jpg":
		jpeg.Encode(writer, a.img.EditedImage, nil)
	case ".png":
		png.Encode(writer, a.img.EditedImage)
	case ".gif":
		gif.Encode(writer, a.img.EditedImage, nil)
	default:
		os.Remove(writer.URI().String()[7:])
		return errors.New("unsupported file extension\n supported extensions: .jpg, .png, .gif")
	}
	return nil
}

func (a *App) deleteFile() {
	if err := os.Remove(a.img.Path); err != nil {
		dialog.NewError(err, a.mainWin)
		return
	}
	if a.img.index == len(a.img.ImagesInFolder)-1 {
		a.nextImage(false, true)
	} else if len(a.img.ImagesInFolder) == 1 {
		a.image.Image = nil
		a.img.EditedImage = nil
		a.img.OriginalImage = nil
		a.rightArrow.Disable()
		a.leftArrow.Disable()
		a.deleteBtn.Disable()
		a.image.Refresh()
	} else {
		a.nextImage(true, true)
	}
}
