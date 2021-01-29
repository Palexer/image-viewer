package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (a *App) loadSettingsUI() {
	winSettings := a.app.NewWindow("Settings")

	themeSelector := widget.NewSelect([]string{"System Default", "Light", "Dark"}, func(selected string) {
		switch selected {
		case "System Default":
			a.app.Settings().SetTheme(theme.DefaultTheme())
		case "Light":
			a.app.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			a.app.Settings().SetTheme(theme.DarkTheme())
		}
		a.app.Preferences().SetString("Theme", selected)
	})
	themeSelector.SetSelected(a.app.Preferences().StringWithFallback("Theme", "System Default"))

	winSettings.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Theme"),
			themeSelector,
		),
	))
	winSettings.Resize(fyne.NewSize(300, 150))
	winSettings.Show()
}
