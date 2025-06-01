package main_ui

import (
	"errors"
	"fmt"
	"log"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/marcs100/minote/config"
	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartUI(appConfigIn *config.Config, configFile string, version string) {
	Conf = appConfigIn
	mainApp = app.NewWithID("minote")
	AppStatus.ConfigFile = configFile
	CreateMainWindow(version)
}

func CreateMainWindow(version string) {

	AppStatus.NoteSize = fyne.NewSize(Conf.Settings.NoteWidth, Conf.Settings.NoteHeight)

	mainWindow = mainApp.NewWindow(fmt.Sprintf("Minote   v%s", version))
	var themeVar theme_variant
	switch Conf.Settings.ThemeVariant {
	case "light":
		themeVar = LIGHT_THEME
	case "dark":
		themeVar = DARK_THEME
	case "system":
		themeVar = SYSTEM_THEME
	}

	AppTheme = GetThemeColours(themeVar)

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(AppStatus.NoteSize)
	AppContainers.grid = grid //store to allow interaction in other functions

	singleNotePage := widget.NewRichTextFromMarkdown("")
	AppWidgets.singleNotePage = singleNotePage
	singleNoteStack := container.NewStack()
	AppContainers.singleNoteStack = singleNoteStack

	PageView.CurrentPage = 0
	PageView.NumberOfPages = 0

	//Create The main panel
	main := CreateMainPanel()

	top := CreateTopPanel()

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(top, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000, 1200))

	//set default view and layout`
	AppStatus.CurrentView = Conf.Settings.InitialView
	fmt.Println("initial view = " + Conf.Settings.InitialView)
	AppStatus.CurrentLayout = Conf.Settings.InitialLayout

	if err := UpdateView(); err != nil {
		fmt.Println(err)
	}

	//keyboard shortcuts
	AddMainKeyboardShortcuts()

	mainWindow.SetCloseIntercept(func() {
		if len(AppStatus.OpenNotes) > 0 {
			fmt.Println(fmt.Sprintf("len of open notes array is %d", len(AppStatus.OpenNotes)))
			//do not close if there are notes open
			dlg := dialog.NewInformation("Error", "There are notes open, please close them before closing the application!", mainWindow)
			dlg.Show()
		} else {
			mainWindow.Close()
		}
	})

	mainWindow.ShowAndRun()
}

func CreateMainPanel() *fyne.Container {

	mainGridContainer := container.NewScroll(AppContainers.grid)
	AppContainers.mainGridContainer = mainGridContainer
	mainPageContainer := container.NewScroll(AppContainers.singleNoteStack)
	AppContainers.mainPageContainer = mainPageContainer
	bgRect := canvas.NewRectangle(AppTheme.MainBgColour)

	mainStackedContainer := container.NewStack(bgRect, mainPageContainer, mainGridContainer)

	return mainStackedContainer
}

func CreateTopPanel() *fyne.Container {
	//AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	spacerLabel := widget.NewLabel("                                ")
	AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes      >", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	AppWidgets.pageLabel = widget.NewLabel("Page: ")
	AppWidgets.pageLabel.Hide()

	toolbar := widget.NewToolbar(
		//show grid view
		widget.NewToolbarAction(theme.GridIcon(), func() {
			if AppStatus.CurrentLayout != LAYOUT_GRID {
				AppStatus.CurrentLayout = LAYOUT_GRID
				PageView.Reset()
				UpdateView()
			}
		}),
		//show single page view
		widget.NewToolbarAction(theme.FileIcon(), func() {
			if AppStatus.CurrentLayout != LAYOUT_PAGE {
				AppStatus.CurrentLayout = LAYOUT_PAGE
				PageView.Reset()
				UpdateView()
			}
		}),
		//page forward
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			if PageView.PageBack() > 0 {
				UpdateView()
			}

		}),
		//page back
		widget.NewToolbarAction(theme.NavigateNextIcon(), func() {
			if PageView.PageForward() > 0 {
				UpdateView()
			}

		}),
	)

	settingsBar := widget.NewToolbar(
		//backup database
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			BackupNotes(Conf.Settings.Database, mainWindow)
		}),

		//display settings
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			NewSettingsWindow()
		}),
	)

	AppWidgets.toolbar = toolbar
	topPanel := container.New(layout.NewHBoxLayout(),
		spacerLabel,
		AppWidgets.viewLabel,
		layout.NewSpacer(),
		toolbar,
		AppWidgets.pageLabel,
		layout.NewSpacer(),
		spacerLabel,
		settingsBar,
	)

	return topPanel
}

func CreateSidePanel() *fyne.Container {
	AppContainers.searchPanel = CreateSearchPanel()
	AppContainers.tagsPanel = CreateTagsPanel()

	newNoteBtn := widget.NewButtonWithIcon("+", theme.DocumentCreateIcon(), func() {
		//CreateNewNote()
	})

	searchBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		//Display the search panel here
		ShowSearchPanel()
	})

	//pinnedBtn := widget.NewButton("P", func(){
	pinnedBtn := widget.NewButtonWithIcon("Pinned", theme.RadioButtonCheckedIcon(), func() {
		AppStatus.CurrentView = VIEW_PINNED
		PageView.Reset()
		err := UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	RecentBtn := widget.NewButtonWithIcon("Recent", theme.HistoryIcon(), func() {
		//AppStatus.Notes,err = minotedb.GetRecentNotes(Conf.Settings.RecentNotesLimit)
		AppStatus.CurrentView = VIEW_RECENT
		PageView.Reset()
		err := UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	tagsBtn := widget.NewButtonWithIcon("Tags", theme.CheckButtonIcon(), func() {
		ToggleMainTagsPanel()
		AppStatus.CurrentView = VIEW_TAGS
		PageView.Reset()
		err := UpdateView()
		if err != nil {
			log.Print("Error getting tagged notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}

	})

	CreateNotebooksList()

	notebooksBtn := widget.NewButtonWithIcon("Notebooks", theme.FolderOpenIcon(), func() {
		ShowNotebooksPanel()
	})

	spacerLabel := widget.NewLabel(" ")

	btnPanel := container.NewVBox(searchBtn, newNoteBtn, spacerLabel, pinnedBtn, RecentBtn, notebooksBtn, tagsBtn)
	AppContainers.listPanel = container.NewStack(AppWidgets.notebooksList)
	AppContainers.listPanel.Hide()
	AppContainers.tagsPanel.Hide()

	sideContainer := container.NewHBox(btnPanel, AppContainers.listPanel, AppContainers.searchPanel, AppContainers.tagsPanel)

	return sideContainer
}

func ShowNotesInGrid(notes []note.NoteData, noteSize fyne.Size) {
	if AppContainers.grid == nil || AppContainers.mainGridContainer == nil {
		return
	}

	if AppContainers.mainPageContainer != nil {
		AppContainers.mainPageContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	PageView.Step = Conf.Settings.GridMaxPages
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	AppContainers.grid.RemoveAll()
	numPages := (PageView.CurrentPage + PageView.Step) - 1
	if numPages > len(notes) {
		numPages = PageView.NumberOfPages
	}

	if AppWidgets.pageLabel.Hidden != true {
		AppWidgets.pageLabel.SetText(PageView.GetGridLabelText())
	}

	for i := PageView.CurrentPage - 1; i < numPages; i++ {
		richText := NewScribeNoteText(notes[i].Content, func() {
			fmt.Println("Note clicked in grid view")
		})
		richText.Wrapping = fyne.TextWrapWord
		themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
		noteColour, _ := RGBStringToFyneColor(notes[i].BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if notes[i].BackgroundColour == "#e7edef" || notes[i].BackgroundColour == "#FFFFFF" || notes[i].BackgroundColour == "#ffffff" || notes[i].BackgroundColour == "#000000" {
			noteBackground = canvas.NewRectangle(AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack := container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		AppContainers.grid.Add(noteStack)
	}
	AppContainers.grid.Refresh()
	AppContainers.mainGridContainer.Show()
}

func ShowNotesAsPages(notes []note.NoteData) {
	//var noteInfo note.NoteInfo
	//var retrievedNote minotedb.NoteData
	//var err error

	if AppContainers.mainGridContainer != nil {
		AppContainers.mainGridContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	PageView.Step = 1
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	if PageView.NumberOfPages == 0 {
		return
	}

	//var noteId = notes[PageView.CurrentPage-1].Id

	AppWidgets.pageLabel.SetText(PageView.GetLabelText())
	AppWidgets.pageLabel.Show()

	AppContainers.singleNoteStack.RemoveAll()

	/*if noteId != 0 {
		retrievedNote, err = minotedb.GetNote(noteId)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	}*/

	/*var allowEdit bool = true
	if slices.Contains(AppStatus.openNotes, noteId) {
		dialog.ShowInformation("Warning", "This note is already open in a separate window.\nClose it first if you want to edit it here!", mainWindow)
		allowEdit = false
		}*/

	//noteContainer = NewNoteContainer(noteId, &noteInfo, &retrievedNote, allowEdit, mainWindow)
	//AddNoteKeyboardShortcuts(&noteInfo, allowEdit, mainWindow)
	//AppContainers.singleNoteStack.Add(noteContainer)
	AppContainers.mainPageContainer.Show()
	AppContainers.mainPageContainer.Refresh()
}

func UpdateView() error {
	//var notes []minotedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch AppStatus.CurrentView {
	case VIEW_PINNED:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}
		AppWidgets.viewLabel.SetText("Pinned Notes")
		AppStatus.Notes, err = notes.GetPinnedNotes()
		AppStatus.CurrentNotebook = ""
	case VIEW_RECENT:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}

		AppWidgets.viewLabel.SetText(("Recent Notes"))
		AppStatus.Notes, err = notes.GetRecentNotes(Conf.Settings.RecentNotesLimit)
		AppStatus.CurrentNotebook = ""
	case VIEW_NOTEBOOK:
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}
		AppWidgets.viewLabel.SetText("Notebook - " + AppStatus.CurrentNotebook)
		AppStatus.Notes, err = notes.GetNotebook(AppStatus.CurrentNotebook)
	case VIEW_TAGS:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		//if len(AppStatus.TagsChecked) > 0 {
		AppWidgets.viewLabel.SetText("Tagged Notes")
		AppStatus.Notes, err = notes.GetTaggedNotes(AppStatus.TagsChecked)
		//} else {
		//	AppStatus.Notes = nil
		//}
	case VIEW_SEARCH:
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}
		if len(strings.TrimSpace(AppWidgets.searchEntry.Text)) > 0 {
			AppStatus.Notes, err = notes.GetSearchResults(AppWidgets.searchEntry.Text, AppStatus.SearchFilter)
			if err == nil {
				AppWidgets.searchResultsLabel.SetText(fmt.Sprintf("Found (%d) > ", len(AppStatus.Notes)))
				AppWidgets.viewLabel.SetText("Search Results")
			}
		}

	default:
		err = errors.New("undefined view")
	}

	if err != nil {
		log.Println("Error in view update!")
		return err
	}

	switch AppStatus.CurrentLayout {
	case LAYOUT_GRID:
		if len(AppStatus.Notes) <= Conf.Settings.GridMaxPages {
			AppWidgets.toolbar.Items[2].ToolbarObject().Hide()
			AppWidgets.toolbar.Items[3].ToolbarObject().Hide()
			AppWidgets.pageLabel.Hide()
		} else {
			AppWidgets.toolbar.Items[2].ToolbarObject().Show()
			AppWidgets.toolbar.Items[3].ToolbarObject().Show()
			AppWidgets.pageLabel.Show()
		}
		ShowNotesInGrid(AppStatus.Notes, AppStatus.NoteSize)
	case LAYOUT_PAGE:
		AppWidgets.toolbar.Items[2].ToolbarObject().Show()
		AppWidgets.toolbar.Items[3].ToolbarObject().Show()
		ShowNotesAsPages(AppStatus.Notes)
	default:
		err = errors.New("undefined layout")
	}

	return err
}

func CreateNotebooksList() {
	var err error
	AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}

	AppWidgets.notebooksList = widget.NewList(
		func() int {
			return len(AppStatus.Notebooks)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("------------Notebooks (xx)------------", func() {})

		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			AppStatus.Notes, _ = notes.GetNotebook(AppStatus.Notebooks[id])
			o.(*widget.Button).SetText(fmt.Sprintf("%s (%d)", AppStatus.Notebooks[id], len(AppStatus.Notes)))
			o.(*widget.Button).OnTapped = func() {
				//AppStatus.Notes,_ = minotedb.GetNotebook(AppStatus.Notebooks[id])
				AppStatus.CurrentView = VIEW_NOTEBOOK
				AppStatus.CurrentNotebook = AppStatus.Notebooks[id]
				PageView.Reset()
				UpdateView()
			}
		},
	)

}

func UpdateNotebooksList() {
	var err error
	AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}
	AppWidgets.notebooksList.Refresh()
}

func ShowNotebooksPanel() {
	UpdateNotebooksList()
	if AppStatus.CurrentView != VIEW_NOTEBOOK {
		AppWidgets.viewLabel.SetText("Notebooks")
	}

	if AppContainers.listPanel != nil {
		if AppContainers.listPanel.Visible() {
			AppContainers.listPanel.Hide()
		} else {
			AppContainers.listPanel.Show()
		}
	}
}

func AddMainKeyboardShortcuts() {
	//Keyboard shortcut to show Pinned Notes
	mainWindow.Canvas().AddShortcut(ScViewPinned, func(shortcut fyne.Shortcut) {
		var err error
		AppStatus.CurrentView = VIEW_PINNED
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	//Keyboard shortcut to show Recent notes
	mainWindow.Canvas().AddShortcut(ScViewRecent, func(shortcut fyne.Shortcut) {
		var err error
		AppStatus.CurrentView = VIEW_RECENT
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	mainWindow.Canvas().AddShortcut(ScPageForward, func(shortcut fyne.Shortcut) {
		if PageView.PageForward() > 0 {
			UpdateView()
		}
	})

	mainWindow.Canvas().AddShortcut(ScPageBack, func(shortcut fyne.Shortcut) {
		if PageView.PageBack() > 0 {
			UpdateView()
		}
	})

	mainWindow.Canvas().AddShortcut(ScFind, func(shortcut fyne.Shortcut) {
		ShowSearchPanel()
	})

	//Keyboard shortcut to create a new note
	mainWindow.Canvas().AddShortcut(ScOpenNote, func(shortcut fyne.Shortcut) {
		//CreateNewNote()
	})

	//Keyboard shortcut to show notebooks list
	mainWindow.Canvas().AddShortcut(ScShowNotebooks, func(shortcut fyne.Shortcut) {
		ShowNotebooksPanel()
	})
}
