package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

	autochangeSelector := widget.NewCheck("Automatically apply filters to image while editing (slower) ", func(b bool) {
		a.autochange = b
		a.app.Preferences().SetBool("Autochange", b)
	})
	autochangeSelector.SetChecked(a.app.Preferences().BoolWithFallback("Autochange", false))

	winSettings.SetContent(widget.NewVBox(
		widget.NewHBox(
			widget.NewLabel("Theme"),
			themeSelector,
		),
		autochangeSelector,
	))
	winSettings.Resize(fyne.NewSize(320, 240))
	winSettings.Show()
}

func (a *App) loadPreferences() {
	// theme
	switch a.app.Preferences().StringWithFallback("Theme", "Dark") {
	case "Light":
		a.app.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.app.Settings().SetTheme(theme.DarkTheme())
	}

	// show/hide panels
	if a.app.Preferences().BoolWithFallback("informationPanelVisible", true) == false {
		a.informationWidgets.Hide()
	}

	if a.app.Preferences().BoolWithFallback("editorVisible", true) == false {
		a.editControls.Hide()
	}

	if a.app.Preferences().BoolWithFallback("statusBarVisible", true) == false {
		a.statusBar.Hide()
	}

	switch a.app.Preferences().BoolWithFallback("Autochange", false) {
	case true:
		a.autochange = true
	case false:
		a.autochange = false
	}
}
