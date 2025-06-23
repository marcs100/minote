package main_app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
)

const (
	EDIT_MODE string = "Edit"
	VIEW_MODE string = "View"
)

type ThemeVariant int

const (
	LIGHT_THEME ThemeVariant = iota
	DARK_THEME
	SYSTEM_THEME
)

type AppColours struct {
	NoteBgColour      color.Color
	MainBgColour      color.Color
	MainCtrlsBgColour color.Color
}

type ApplicationStatus struct {
	ConfigFile          string
	CurrentView         string
	CurrentNotebook     string
	CurrentLayout       string
	CurrentSortSelected string
	Notes               []note.NoteData
	Notebooks           []string
	Tags                []string
	TagsChecked         []string
	NoteSize            fyne.Size
	SearchFilter        notes.SearchFilter
}
