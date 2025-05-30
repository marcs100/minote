package main_ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateSearchPanel() *fyne.Container {

	AppWidgets.searchResultsLabel = widget.NewLabel("")
	filterLabel := widget.NewLabel("Filter: -")
	searchFilter := widget.NewCheckGroup([]string{SEARCH_FILT_WOLE_WORDS, SEARCH_FILT_PINNED}, func(selected []string) {
		AppStatus.searchFilter.Pinned = false
		AppStatus.currentView = VIEW_SEARCH
		AppStatus.searchFilter.WholeWords = false
		for _, sel := range selected {
			fmt.Println("selected: " + sel)
			if sel == SEARCH_FILT_PINNED {
				AppStatus.searchFilter.Pinned = true
			}

			if sel == SEARCH_FILT_WOLE_WORDS {
				AppStatus.searchFilter.WholeWords = true
			}
		}
		var err error = UpdateView()
		if err != nil {
			log.Print("Error getting search results (after setting filter): ")
			dialog.ShowError(err, mainWindow)
			//log.Panic(err)
		}
	})
	searchLabel := widget.NewLabel("               Search:               ")
	AppWidgets.searchEntry = widget.NewEntry()
	AppWidgets.searchEntry.OnSubmitted = func(text string) {
		AppStatus.currentView = VIEW_SEARCH
		var err error = UpdateView()
		mainWindow.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		if err != nil {
			log.Print("Error getting search results: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}

	}

	searchPanel := container.NewVBox(searchLabel, AppWidgets.searchEntry, AppWidgets.searchResultsLabel, filterLabel, searchFilter)
	return searchPanel
}

func ShowSearchPanel() {
	if AppContainers.searchPanel.Hidden {
		AppContainers.searchPanel.Show()
		mainWindow.Canvas().Focus(AppWidgets.searchEntry)
	} else {
		mainWindow.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		AppContainers.searchPanel.Hide()
	}
	mainWindow.Canvas().Refresh(AppContainers.searchPanel)
}
