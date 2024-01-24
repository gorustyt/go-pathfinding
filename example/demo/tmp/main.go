package main

import (
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	btn        *widget.Button
	on         bool
	maxLayout  *fyne.Container
	closeColor *canvas.Rectangle
	openColor  *canvas.Rectangle
)

func timer() {
	time.Sleep(5 * time.Second)
	btn.SetText("aaaa")
}

func GetCloseColor() *canvas.Rectangle {
	if closeColor == nil {
		closeColor = canvas.NewRectangle(color.NRGBA{R: 242, G: 242, B: 242, A: 255})
	}
	return closeColor
}

func GetOpenColor() *canvas.Rectangle {
	if openColor == nil {
		openColor = canvas.NewRectangle(color.NRGBA{R: 204, G: 255, B: 204, A: 255})
	}
	return openColor
}

func main() {

	os.Setenv("FYNE_FONT", "simhei.ttf")

	a := app.New()
	w := a.NewWindow("模拟设备")
	w.Resize(fyne.NewSize(320, 230))
	setting := a.Settings()
	setting.SetTheme(theme.LightTheme())

	ip := widget.NewLabel("服务器IP:172.16.3.3, 端口:1883")

	btn = widget.NewButton("big button", func() {
		// hello.SetText("Welcome :)")
		if !on {
			on = true
			openColor.Show()
			closeColor.Hide()
			btn.SetText("已打开，点击关闭")
		} else {
			on = false
			openColor.Hide()
			closeColor.Show()
			btn.SetText("已关闭，点击打开")
		}
	})

	btn.Resize(fyne.NewSize(260, 160))
	btn.Move(fyne.NewPos(30-theme.Padding(), 10))

	// closeColor = canvas.NewRectangle(color.NRGBA{R: 242, G: 242, B: 242, A: 255})
	// // closeColor.MinSize()

	// openClolor = canvas.NewRectangle(color.NRGBA{R: 204, G: 255, B: 204, A: 255})
	// // closeColor.MinSize()

	maxLayout = container.New(layout.NewMaxLayout(), GetCloseColor(), GetOpenColor(), btn)
	maxLayout.Resize(fyne.NewSize(260, 160))
	maxLayout.Move(fyne.NewPos((320-260)/2-theme.Padding(), 10))

	openColor.Hide()
	closeColor.Show()

	ly := container.NewWithoutLayout(maxLayout)
	w.SetContent(container.NewVBox(
		ip,
		ly,
	))

	go timer()
	w.ShowAndRun()
}
