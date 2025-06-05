package ui

import "fyne.io/fyne/v2"

var mainWindow fyne.Window

var AppContainers ApplicationContainers //structure containing pointers to fyne containers for main window
var AppWidgets ApplicationWidgets       //structure containing pointers to fyne widgets for main window
var PageView PageViewStatus             //structure to track page numbers
