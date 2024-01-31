package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
	"log"
	"net/url"
)

func MakeMenu(a fyne.App, w fyne.Window, cfg *grid_map.Config) *fyne.MainMenu {
	openItem := fyne.NewMenuItem("Open map", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			cfg.Load(reader)
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})
	saveItem := fyne.NewMenuItem("Save map", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				log.Println("Cancelled")
				return
			}
			cfg.Dump(writer)
		}, w)
	})

	openSettings := func() {
		w := a.NewWindow("Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	label := widget.NewLabel("setting theme")
	label.Alignment = fyne.TextAlignCenter
	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
		}),
	)
	themeItem := fyne.NewMenuItem("Theme", func() {
		w1 := a.NewWindow("Theme Settings")
		w1.SetContent(container.NewVBox(
			label,
			themes,
		))
		w1.Resize(fyne.NewSize(200, 200))
		w1.Show()
	})
	settingMenu := fyne.NewMenu("Settings", settingsItem, themeItem)
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://github.com/gorustyt/go-pathfinding")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
	)
	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", openItem, saveItem)
	main := fyne.NewMainMenu(
		file,
		settingMenu,
		helpMenu,
	)
	return main
}
