package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
)

func (mw *MainWindow) CreateSearchPanel() *fyne.Container {

	mw.AppWidgets.searchResultsLabel = widget.NewLabel("")
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
		var err error = mw.UpdateView()
		if err != nil {
			log.Print("Error getting search results (after setting filter): ")
			dialog.ShowError(err, mw.window)
			//log.Panic(err)
		}
	})
	searchLabel := widget.NewLabel("               Search:               ")
	mw.AppWidgets.searchEntry = widget.NewEntry()
	mw.AppWidgets.searchEntry.OnSubmitted = func(text string) {
		main_app.AppStatus.CurrentView = main_app.VIEW_SEARCH
		var err error = mw.UpdateView()
		mw.window.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		if err != nil {
			log.Print("Error getting search results: ")
			dialog.ShowError(err, mw.window)
			log.Panic(err)
		}

	}

	searchPanel := container.NewVBox(searchLabel, mw.AppWidgets.searchEntry, mw.AppWidgets.searchResultsLabel, filterLabel, searchFilter)
	return searchPanel
}

func (mw *MainWindow) ShowSearchPanel() {
	if mw.AppContainers.searchPanel.Hidden {
		mw.AppContainers.searchPanel.Show()
		mw.window.Canvas().Focus(mw.AppWidgets.searchEntry)
	} else {
		mw.window.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		mw.AppContainers.searchPanel.Hide()
	}
	mw.window.Canvas().Refresh(mw.AppContainers.searchPanel)
}
