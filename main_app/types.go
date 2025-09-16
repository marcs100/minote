package main_app

import (
	"fyne.io/fyne/v2"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
)

const (
	EDIT_MODE string = "Edit"
	VIEW_MODE string = "View"
)

type ApplicationStatus struct {
	ConfigFile          string
	CurrentView         string
	CurrentNotebook     string
	CurrentLayout       string
	CurrentSortSelected string
	Notes               []note.NoteData
	Notebooks           []string
	// Tags                []string
	TagsChecked  []string
	NoteSize     fyne.Size
	SearchFilter notes.SearchFilter
}

type About struct {
	Version     string
	Licence     string
	LicenceLink string
	Maintainer  string
	Website     string
}
