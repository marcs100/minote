package ui

import "fyne.io/fyne/v2"

var mainWindow fyne.Window

var AppContainers ApplicationContainers //structure containing pointers to fyne containers for main window
var AppWidgets ApplicationWidgets       //structure containing pointers to fyne widgets for main window
var PageView PageViewStatus             //structure to track page numbers
var SortViews = map[string]int{
	"Newest First":       0,
	"Oldest First":       1,
	"Newly Pinned First": 2,
	"Newly Pinned last":  3,
	"Created First":      4,
	"Created Last":       5,
}
