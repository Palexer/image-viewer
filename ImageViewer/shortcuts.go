package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func (a *App) loadKeyboardShortcuts() {
	// keyboard shortcuts
	// ctrl+o to open file
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.openFileDialog() })

	// ctrl+s to save file
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.saveFileDialog() })

	// ctrl+z to undo
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyZ,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.undo() })

	// ctrl+y to redo
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyY,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.redo() })

	// ctrl+q to quit application
	a.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: a.mainModKey,
	}, func(shortcut fyne.Shortcut) { a.app.Quit() })

	a.mainWin.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		// move forward/back within the current folder of images
		case fyne.KeyRight:
			a.nextImage(true, false)
		case fyne.KeyLeft:
			a.nextImage(false, false)
		// delete images with delete key
		case fyne.KeyDelete:
			if a.image.Image != nil {
				dialog.ShowConfirm("Delte file?", "Do you really want to delete this image?\n This action can't be undone.", func(b bool) {
					if b {
						a.deleteFile()
					}
				}, a.mainWin)
			}
		// close dialogs with esc key
		case fyne.KeyEscape:
			if len(a.mainWin.Canvas().Overlays().List()) > 0 {
				a.mainWin.Canvas().Overlays().Top().Hide()
			}
		case fyne.KeyF11:
			if a.image.Image == nil {
				return
			}
			a.fullscreenMode()
		case fyne.KeyF2:
			if a.image.Image == nil {
				return
			}
			a.renameDialog()
		}
	})
}

func (a *App) showShortcuts() {
	shortcuts := []string{
		"Ctrl+O", "Ctrl+S", "Ctrl+Z",
		"Ctrl+Y", "Ctrl+Q", "F11",
		"Arrow Right", "Arrow Left", "Delete",
		"F2", "Escape"}
	descriptions := []string{
		"Open File", "Save File", "Undo",
		"Redo", "Quit Application", "Fullscreen View",
		"Next Image", "Last Image", "Delete Image",
		"Rename", "Close dialog"}

	win := a.app.NewWindow("Keyboard Shortcuts")
	table := widget.NewTable(
		func() (int, int) { return len(shortcuts), 2 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row == 0 {
				if id.Col == 0 {
					label.SetText("Description")
					label.TextStyle.Bold = true
				} else {
					label.SetText("Shortcut")
					label.TextStyle.Bold = true
				}
			} else {
				if id.Col == 0 {
					label.SetText(descriptions[id.Row-1])
				} else {
					label.SetText(descriptions[id.Row-1])
				}
			}
		},
	)
	table.SetColumnWidth(0, 250)
	table.SetColumnWidth(1, 250)
	win.SetContent(table)
	win.Resize(fyne.NewSize(500, 500))
	win.Show()
}
