package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/note"
	"image/color"
)

type ThemeVariant int

type PageViewStatus struct {
	NumberOfPages int
	CurrentPage   int
	Step          int
}

// containers for main window
type ApplicationContainers struct {
	grid              *fyne.Container
	singleNoteStack   *fyne.Container
	mainGridContainer *container.Scroll
	mainPageContainer *container.Scroll
	listPanel         *fyne.Container
	searchPanel       *fyne.Container
	tagsPanel         *fyne.Container
	tagsList          *fyne.Container
}

// widgets for main window
type ApplicationWidgets struct {
	Toolbar        *widget.Toolbar
	singleNotePage *widget.RichText
	viewLabel      *widget.Label
	pageLabel      *widget.Label
	notebooksList  *widget.List
	//tagsList           *widget.List
	searchEntry        *widget.Entry
	searchResultsLabel *widget.Label
	sortSelect         *widget.Select
}

type MainWindow struct {
	window        fyne.Window
	AppWidgets    ApplicationWidgets
	AppContainers ApplicationContainers
	ThemeVariant  ThemeVariant
	UI_Colours    AppColours
}

type NotePage struct {
	Id                 uint
	ParentWindow       fyne.Window
	MainAppWindow      *MainWindow
	NoteInfo           note.NoteInfo
	RetrievedNote      note.NoteData
	AllowEdit          bool
	NewWindowMode      bool
	NotePageWidgets    NotePageWidgets
	NotePageContainers NotePageContainers
	NotePageCanvas     NotePageCanvas
	ThemeVariant       ThemeVariant
	UI_Colours         AppColours
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
	TagsButton     *widget.Button
}

type NotePageContainers struct {
	Markdown        *fyne.Container
	PropertiesPanel *fyne.Container
	TagsPanel       *container.Scroll
	//TagsPanel *fyne.Container
	TagLabels *fyne.Container
}

type NotePageCanvas struct {
	NoteBackground *canvas.Rectangle
}

type AppColours struct {
	NoteBgColour      color.Color
	MainBgColour      color.Color
	MainCtrlsBgColour color.Color
}
