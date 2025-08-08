package ui

//import "fyne.io/fyne/v2"

//var mainWindow fyne.Window

// var AppContainers ApplicationContainers //structure containing pointers to fyne containers for main window
// var AppWidgets ApplicationWidgets       //structure containing pointers to fyne widgets for main window
var PageView PageViewStatus //structure to track page numbers
var SortViews = map[string]int{
	"Modified: new to old":   0,
	"Modified: old to new":   1,
	"Pinned: most recently":  2,
	"Pinned: least recently": 3,
	"Created: new to old":    4,
	"Created: old to new":    5,
}

const (
	LIGHT_THEME ThemeVariant = iota
	DARK_THEME
	AUTO_THEME
)

//var UI_Colours AppColours
