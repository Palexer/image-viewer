package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (a *App) loadSettingsUI() {
	winSettings := a.app.NewWindow("Settings")

	themeSelector := widget.NewSelect([]string{"Light", "Dark"}, func(selected string) {
		switch selected {
		case "Light":
			a.app.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			a.app.Settings().SetTheme(theme.DarkTheme())
		}
		a.app.Preferences().SetString("Theme", selected)
	})
	themeSelector.SetSelected(a.app.Preferences().StringWithFallback("Theme", "Dark"))

	winSettings.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Theme"),
			themeSelector,
		),
		layout.NewSpacer(),
		widget.NewLabel("v1.1 | License: MIT"),
		widget.NewHyperlink("Github (source code and more information)", parseURL("https://github.com/Palexer/image-viewer")),
	))
	winSettings.Show()
}
