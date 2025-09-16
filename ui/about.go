package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
)

// var aboutDlg *dialog.FormDialog
var custAboutDlg *dialog.CustomDialog

func NewAbout(about main_app.About, parentWindow fyne.Window) {
	aboutGrid := container.NewGridWithColumns(1,
		container.NewHBox(
			widget.NewLabel("Minote: "),
			widget.NewRichTextFromMarkdown(fmt.Sprintf("**v%s**", about.Version))),
		container.NewHBox(
			widget.NewLabel("Licence: "),
			widget.NewRichTextFromMarkdown(fmt.Sprintf("[%s](%s)", about.Licence, about.LicenceLink))),
		container.NewHBox(
			widget.NewLabel("Website: "),
			widget.NewRichTextFromMarkdown(fmt.Sprintf("[%s](%s)", about.Website, about.Website))),
		container.NewHBox(
			widget.NewLabel("Maintainer: "),
			widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", about.Maintainer))))

	custAboutDlg = dialog.NewCustom("About", "Ok", aboutGrid, parentWindow)
}

func ShowAbout() {
	if custAboutDlg != nil {
		// aboutDlg.Show()
		custAboutDlg.Show()
	}
}
