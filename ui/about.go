package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
)

var aboutDlg *dialog.FormDialog

func NewAbout(about main_app.About, parentWindow fyne.Window) {

	var UI_Colours = GetAppColours(GetThemeVariant())
	bg := canvas.NewRectangle(UI_Colours.MainBgColour)
	versionLabel := widget.NewLabel("Version")
	versionText := widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", about.Version))
	versionGrid := container.NewGridWithRows(1, versionLabel, versionText)
	aboutStack := container.NewStack(bg, versionGrid)
	formItem := widget.NewFormItem("", aboutStack)
	aboutDlg = dialog.NewForm("About", "Ok", "", []*widget.FormItem{formItem}, func(b bool) {

	}, parentWindow)
}

func ShowAbout() {
	if aboutDlg != nil {
		aboutDlg.Show()
	}
}
