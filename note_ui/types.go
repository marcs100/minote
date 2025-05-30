package note_ui

import (
	"fyne.io/fyne/v2"
	"github.com/marcs100/minote/note"
)

type NotePage struct {
	Id           uint
	ParentWindow fyne.Window
	NoteInfo     note.NoteInfo
	AllowEdit    bool
}
