package main_ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
)

func CreateSearchPanel() *fyne.Container {

	AppWidgets.searchResultsLabel = widget.NewLabel("")
	filterLabel := widget.NewLabel("Filter: -")
	searchFilter := widget.NewCheckGroup([]string{main_app.SEARCH_FILT_WOLE_WORDS, main_app.SEARCH_FILT_PINNED}, func(selected []string) {
		main_app.AppStatus.SearchFilter.Pinned = false
		main_app.AppStatus.CurrentView = main_app.VIEW_SEARCH
		main_app.AppStatus.SearchFilter.WholeWords = false
		for _, sel := range selected {
			fmt.Println("selected: " + sel)
			if sel == main_app.SEARCH_FILT_PINNED {
				main_app.AppStatus.SearchFilter.Pinned = true
			}

			if sel == main_app.SEARCH_FILT_WOLE_WORDS {
				main_app.AppStatus.SearchFilter.WholeWords = true
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
		main_app.AppStatus.CurrentView = main_app.VIEW_SEARCH
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
