package note_ui

import (
	"fyne.io/fyne/v2"
	"github.com/marcs100/minote/note"
)

func (np *NotePage) NewNotePage(retrievedNote *note.NoteData, allowEdit bool, parentWindow fyne.Window) {
	np.ParentWindow = parentWindow
	np.AllowEdit = allowEdit
	np.NoteInfo.NewNote = false
	if retrievedNote == nil {
		np.NoteInfo.NewNote = true
		np.NoteInfo.Id = 0
	}

}
