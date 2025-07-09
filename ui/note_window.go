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

func NewNoteWindow(noteId uint, parentWindow fyne.Window, mainAppWindow *MainWindow) {
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
		if main_app.AppStatus.CurrentView == main_app.VIEW_NOTEBOOK {
			retrievedNote.Notebook = main_app.AppStatus.CurrentNotebook
		} else {
			retrievedNote.Notebook = "General"
		}
	}

	noteWindow := main_app.MainApp.NewWindow("")
	var np NotePage
	noteContainer := np.NewNotePage(&retrievedNote, true, true, noteWindow, mainAppWindow)

	noteWindow.SetTitle(fmt.Sprintf("Notebook: %s", np.NoteInfo.Notebook))
	noteWindow.Resize(fyne.NewSize(900, 750))

	noteWindow.SetContent(noteContainer)
	noteWindow.Canvas().Focus(np.NotePageWidgets.Entry)
	noteWindow.SetOnClosed(func() {
		fmt.Printf("Closing note %d", np.NoteInfo.Id)
		np.SaveNote()
		tracker.DelFromTracker(np.NoteInfo.Id)
	})

	np.AddNoteKeyboardShortcuts()

	if np.NoteInfo.NewNote {
		np.SetEditMode()
	}

	np.RefreshWindow()

	noteWindow.Show()
}
