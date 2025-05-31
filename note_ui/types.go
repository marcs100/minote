package note_ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/note"
)

type NotePage struct {
	Id                 uint
	ParentWindow       fyne.Window
	NoteInfo           note.NoteInfo
	AllowEdit          bool
	NotePageWidgets    NotePageWidgets
	NotePageContainers NotePageContainers
	NotePageCanvas     NotePageCanvas
}

type NotePageWidgets struct {
	PinButton    *widget.Button
	MarkdownText *widget.RichText
	DeleteButton *widget.Button
	Entry        *EntryCustom
	//entry      *widget.Entry
	ModeSelect     *widget.RadioGroup
	PropertiesText *widget.Label
	AddTagButton   *widget.Button
}

type NotePageContainers struct {
	Markdown        *fyne.Container
	PropertiesPanel *fyne.Container
	TagsPanel       *container.Scroll
	TagLabels       *fyne.Container
}

type NotePageCanvas struct {
	NoteBackground *canvas.Rectangle
}
