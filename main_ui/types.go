package main_ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
)

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
}

// widgets for main window
type ApplicationWidgets struct {
	toolbar            *widget.Toolbar
	singleNotePage     *widget.RichText
	viewLabel          *widget.Label
	pageLabel          *widget.Label
	notebooksList      *widget.List
	tagsList           *widget.List
	searchEntry        *widget.Entry
	searchResultsLabel *widget.Label
}

type ApplicationStatus struct {
	configFile      string
	currentView     string
	currentNotebook string
	currentLayout   string
	notes           []note.NoteData
	notebooks       []string
	tags            []string
	tagsChecked     []string
	openNotes       []uint //maintain a list of notes that are currently open
	noteSize        fyne.Size
	searchFilter    notes.SearchFilter
}

// widgets for note window
/*type NoteWindowWidgets struct {
pinButton    *widget.Button
markdownText *widget.RichText
//markdownText *RichTextFromMarkdownCustom
deleteButton *widget.Button
entry        *EntryCustom
//entry      *widget.Entry
modeSelect     *widget.RadioGroup
propertiesText *widget.Label
addTagButton   *widget.Button
}*/

/*type NoteWindowContainers struct {
markdown        *fyne.Container
propertiesPanel *fyne.Container
tagsPanel       *container.Scroll
tagLabels       *fyne.Container
}*/

/*type NoteWindowCanvas struct {
noteBackground *canvas.Rectangle
}*/

type AppColours struct {
	NoteBgColour color.Color
	MainBgColour color.Color
}

type theme_variant int
