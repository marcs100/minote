package ui

import (
	"fmt"
	"log"

	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
	"github.com/marcs100/minote/tracker"

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
	} else {
		//new note
		retrievedNote.Id = 0
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
		tracker.DelFromTracker(np.NoteInfo.Id)
	})

	np.AddNoteKeyboardShortcuts()

	if np.NoteInfo.NewNote {
		np.SetEditMode()
	}

	noteWindow.Show()
}
