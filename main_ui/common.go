package main_ui

import "fyne.io/fyne/v2"

var mainWindow fyne.Window

// var noteContainer *fyne.Container = nil
var AppContainers ApplicationContainers //structure containing pointers to fyne containers for main window
// var NoteContainers NoteWindowContainers //structure containing poineters to fyne containers for note window
var AppWidgets ApplicationWidgets //structure containing pointers to fyne widgets for main window
// var NoteWidgets NoteWindowWidgets       //structure containing pointers to fyne widgets for note window
// var NoteCanvas NoteWindowCanvas         //structure containing pinters to canvas object for note window
var PageView PageViewStatus //structure to track page numbers
