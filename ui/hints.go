package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
)

var custHintsDlg *dialog.CustomDialog

func NewHints(parentWindow fyne.Window) {
	hintsGrid := container.NewGridWithColumns(2,
		widget.NewRichTextFromMarkdown("## When in grid View:"),
		widget.NewLabel(""),
		widget.NewLabel(fmt.Sprintf("%-50s", "mouse left click:")),
		widget.NewLabel("open note in new window"),
		widget.NewLabel(fmt.Sprintf("%-50s", "mouse right click:")),
		widget.NewLabel("open note in single page view"),
		widget.NewRichTextFromMarkdown("## Main window - keyboard shortcuts:"),
		widget.NewLabel(""),
		widget.NewLabel(fmt.Sprintf("%-50s", "Open a new note:")),
		widget.NewLabel(GetScName(main_app.ScOpenNote.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Search:")),
		widget.NewLabel(GetScName(main_app.ScFind.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide pinned notes:")),
		widget.NewLabel(GetScName(main_app.ScViewPinned.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide recent notes:")),
		widget.NewLabel(GetScName(main_app.ScViewRecent.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide tags panel:")),
		widget.NewLabel(GetScName(main_app.ScShowTags.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide notebooks panel:")),
		widget.NewLabel(GetScName(main_app.ScShowNotebooks.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Page forward:")),
		widget.NewLabel(GetScName(main_app.ScPageForward.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Page back:")),
		widget.NewLabel(GetScName(main_app.ScPageBack.ShortcutName())),
		widget.NewRichTextFromMarkdown("## Notes - keyboard shortcuts:"),
		widget.NewLabel(""),
		widget.NewLabel(fmt.Sprintf("%-50s", "Enter edit mode:")),
		widget.NewLabel(GetScName(main_app.ScEditMode.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Exit edit mode:")),
		widget.NewLabel("Escape"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Pin/unpin note:")),
		widget.NewLabel(GetScName(main_app.ScPinNote.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide tags:")),
		widget.NewLabel(GetScName(main_app.ScNoteTags.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Change colour:")),
		widget.NewLabel(GetScName(main_app.ScNoteColour.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Select Notebook:")),
		widget.NewLabel(GetScName(main_app.ScChangeNoteNotebook.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Show/hide proprties panel:")),
		widget.NewLabel(GetScName(main_app.ScShowInfo.ShortcutName())),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert date:")),
		widget.NewLabel("F1"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert time:")),
		widget.NewLabel("F2"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert date and time:")),
		widget.NewLabel("F3"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert custom (F4) snippet:")),
		widget.NewLabel("F4"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert custom (F5) snippet:")),
		widget.NewLabel("F5"),
		widget.NewLabel(fmt.Sprintf("%-50s", "Insert custom (F6) snippet:")),
		widget.NewLabel("F6"))

	scrolledCont := container.NewVScroll(hintsGrid)
	scrolledCont.SetMinSize(fyne.NewSize(400, 900))
	custHintsDlg = dialog.NewCustom("Hints/keyboard shortcuts", "Close", scrolledCont, parentWindow)
}

func GetScName(longName string) string {
	s := strings.Split(longName, ":")
	if len(s) == 2 {
		return s[1]
	}

	return longName
}

func ShowHints() {
	if custHintsDlg != nil {
		custHintsDlg.Show()
	}
}
