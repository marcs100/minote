package note_ui

import (
	"fmt"
	"log"
	"slices"

	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	//"fyne.io/fyne/v2/data/binding"
	//"github.com/fyne-io/terminal"
)

func NewNoteWindow(noteId uint, parentWindow fyne.Window) {
	var retrievedNote note.NoteData
	var err error

	if noteId != 0 {
		retrievedNote, err = notes.GetNote(noteId)
		if err != nil {
			dialog.ShowError(err, parentWindow)
			log.Panic(err)
		}
	}

	noteWindow := main_app.MainApp.NewWindow("")
	var np NotePage
	noteContainer := np.NewNotePage(&retrievedNote, true, noteWindow)

	//fmt.Println(fmt.Sprintf("************Notebook is %s", "debug"))
	noteWindow.SetTitle(fmt.Sprintf("Notebook: %s", np.NoteInfo.Notebook))
	noteWindow.Resize(fyne.NewSize(900, 750))

	noteWindow.SetContent(noteContainer)
	noteWindow.Canvas().Focus(np.NotePageWidgets.Entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", np.NoteInfo.Id))
		np.SaveNote()

		if index := slices.Index(main_app.AppStatus.OpenNotes, np.NoteInfo.Id); index != -1 {
			main_app.AppStatus.OpenNotes = slices.Delete(main_app.AppStatus.OpenNotes, index, index+1)
		}
	})

	np.AddNoteKeyboardShortcuts()

	if np.NoteInfo.NewNote {
		np.SetEditMode()
	}

	noteWindow.Show()
}
