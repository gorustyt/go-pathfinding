package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
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
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
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

		}, w)
	})
	openFolderItem := fyne.NewMenuItem("Open Folder", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}

			children, err := list.List()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			dialog.ShowInformation("Folder Open", out, w)
		}, w)
	})
	openItem.Checked = true

	openSettings := func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut

	themeItem := fyne.NewMenuItem("theme", nil)
	settingMenu := fyne.NewMenu("Settings", settingsItem, themeItem)
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
	)
	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", openFolderItem, openItem, saveItem)
	main := fyne.NewMainMenu(
		file,
		settingMenu,
		helpMenu,
	)
	return main
}
