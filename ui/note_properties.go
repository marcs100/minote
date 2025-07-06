package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/note"
)

func CreateProperetiesPanel(np *NotePage) error {
	themeBackground := canvas.NewRectangle(np.UI_Colours.NoteBgColour)
	propertiesTitle := widget.NewRichTextFromMarkdown("**Properties**")
	np.NotePageWidgets.PropertiesText = widget.NewLabel("")
	vbox := container.NewVBox(propertiesTitle, np.NotePageWidgets.PropertiesText)
	propertiesPadded := container.NewPadded(themeBackground, vbox)
	propertiesPanel := container.NewVScroll(propertiesPadded)
	np.NotePageContainers.PropertiesPanel = container.NewStack(propertiesPanel)
	return nil // no error
}

func (np *NotePage) ShowProperties() {
	if np.NotePageContainers.PropertiesPanel.Hidden {
		//fmt.Println("Will show properties panel")
		text := note.GetPropertiesText(&np.NoteInfo)
		np.NotePageWidgets.PropertiesText.SetText(text)
		np.NotePageContainers.PropertiesPanel.Show()
	} else {
		//fmt.Println("Will hide properties panel")
		np.NotePageContainers.PropertiesPanel.Hide()
		np.RefreshWindow()
	}
}

func (np *NotePage) UpdateProperties() {
	if !np.NotePageContainers.PropertiesPanel.Hidden {
		text := note.GetPropertiesText(&np.NoteInfo)
		np.NotePageWidgets.PropertiesText.SetText(text)
		np.NotePageWidgets.PropertiesText.Refresh()
	}
}
