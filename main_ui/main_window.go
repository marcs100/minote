package main_ui

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/marcs100/minote/config"
	"github.com/marcs100/minote/conversions"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/note_ui"
	"github.com/marcs100/minote/notes"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartUI(appConfigIn *config.Config, configFile string, version string) {
	main_app.Conf = appConfigIn
	mainApp = app.NewWithID("minote")
	main_app.AppStatus.ConfigFile = configFile
	CreateMainWindow(version)
}

func CreateMainWindow(version string) {

	main_app.AppStatus.NoteSize = fyne.NewSize(main_app.Conf.Settings.NoteWidth, main_app.Conf.Settings.NoteHeight)

	mainWindow = mainApp.NewWindow(fmt.Sprintf("Minote   v%s", version))
	var themeVar main_app.ThemeVariant
	switch main_app.Conf.Settings.ThemeVariant {
	case "light":
		themeVar = main_app.LIGHT_THEME
	case "dark":
		themeVar = main_app.DARK_THEME
	case "system":
		themeVar = main_app.SYSTEM_THEME
	}

	main_app.AppTheme = main_app.GetThemeColours(themeVar)

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(main_app.AppStatus.NoteSize)
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
	main_app.AppStatus.CurrentView = main_app.Conf.Settings.InitialView
	fmt.Println("initial view = " + main_app.Conf.Settings.InitialView)
	main_app.AppStatus.CurrentLayout = main_app.Conf.Settings.InitialLayout

	if err := UpdateView(); err != nil {
		fmt.Println(err)
	}

	//keyboard shortcuts
	AddMainKeyboardShortcuts()

	mainWindow.SetCloseIntercept(func() {
		if len(main_app.AppStatus.OpenNotes) > 0 {
			fmt.Println(fmt.Sprintf("len of open notes array is %d", len(main_app.AppStatus.OpenNotes)))
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
	bgRect := canvas.NewRectangle(main_app.AppTheme.MainBgColour)

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
			if main_app.AppStatus.CurrentLayout != main_app.LAYOUT_GRID {
				main_app.AppStatus.CurrentLayout = main_app.LAYOUT_GRID
				PageView.Reset()
				UpdateView()
			}
		}),
		//show single page view
		widget.NewToolbarAction(theme.FileIcon(), func() {
			if main_app.AppStatus.CurrentLayout != main_app.LAYOUT_PAGE {
				main_app.AppStatus.CurrentLayout = main_app.LAYOUT_PAGE
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
			BackupNotes(main_app.Conf.Settings.Database, mainWindow)
		}),

		//display settings
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			NewSettingsWindow()
		}),
	)

	AppWidgets.Toolbar = toolbar
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
		main_app.AppStatus.CurrentView = main_app.VIEW_PINNED
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
		main_app.AppStatus.CurrentView = main_app.VIEW_RECENT
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
		main_app.AppStatus.CurrentView = main_app.VIEW_TAGS
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
	PageView.Step = main_app.Conf.Settings.GridMaxPages
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
		themeBackground := canvas.NewRectangle(main_app.AppTheme.NoteBgColour)
		noteColour, _ := conversions.RGBStringToFyneColor(notes[i].BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if notes[i].BackgroundColour == "#e7edef" || notes[i].BackgroundColour == "#FFFFFF" || notes[i].BackgroundColour == "#ffffff" || notes[i].BackgroundColour == "#000000" {
			noteBackground = canvas.NewRectangle(main_app.AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
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
	var noteInfo note.NoteInfo
	var retrievedNote note.NoteData
	var err error = nil

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

	var noteId = notes[PageView.CurrentPage-1].Id

	AppWidgets.pageLabel.SetText(PageView.GetLabelText())
	AppWidgets.pageLabel.Show()

	AppContainers.singleNoteStack.RemoveAll()

	if noteId != 0 {
		retrievedNote, err = notes.GetNote(noteId)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	}

	var allowEdit bool = true
	if slices.Contains(main_app.AppStatus.OpenNotes, noteId) {
		dialog.ShowInformation("Warning", "This note is already open in a separate window.\nClose it first if you want to edit it here!", mainWindow)
		allowEdit = false
	}

	var np note_ui.NotePage
	noteContainer := np.NewNotePage(&retrievedNote, allowEdit, mainWindow)
	np.AddNoteKeyboardShortcuts()
	AppContainers.singleNoteStack.Add(noteContainer)
	AppContainers.mainPageContainer.Show()
	AppContainers.mainPageContainer.Refresh()
}

func UpdateView() error {
	//var notes []minotedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch main_app.AppStatus.CurrentView {
	case main_app.VIEW_PINNED:
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
		main_app.AppStatus.Notes, err = notes.GetPinnedNotes()
		main_app.AppStatus.CurrentNotebook = ""
	case main_app.VIEW_RECENT:
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
		main_app.AppStatus.Notes, err = notes.GetRecentNotes(main_app.Conf.Settings.RecentNotesLimit)
		main_app.AppStatus.CurrentNotebook = ""
	case main_app.VIEW_NOTEBOOK:
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}
		AppWidgets.viewLabel.SetText("Notebook - " + main_app.AppStatus.CurrentNotebook)
		main_app.AppStatus.Notes, err = notes.GetNotebook(main_app.AppStatus.CurrentNotebook)
	case main_app.VIEW_TAGS:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		//if len(AppStatus.TagsChecked) > 0 {
		AppWidgets.viewLabel.SetText("Tagged Notes")
		main_app.AppStatus.Notes, err = notes.GetTaggedNotes(main_app.AppStatus.TagsChecked)
		//} else {
		//	AppStatus.Notes = nil
		//}
	case main_app.VIEW_SEARCH:
		if AppContainers.tagsPanel != nil {
			AppContainers.tagsPanel.Hide()
		}
		if len(strings.TrimSpace(AppWidgets.searchEntry.Text)) > 0 {
			main_app.AppStatus.Notes, err = notes.GetSearchResults(AppWidgets.searchEntry.Text, main_app.AppStatus.SearchFilter)
			if err == nil {
				AppWidgets.searchResultsLabel.SetText(fmt.Sprintf("Found (%d) > ", len(main_app.AppStatus.Notes)))
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

	switch main_app.AppStatus.CurrentLayout {
	case main_app.LAYOUT_GRID:
		if len(main_app.AppStatus.Notes) <= main_app.Conf.Settings.GridMaxPages {
			AppWidgets.Toolbar.Items[2].ToolbarObject().Hide()
			AppWidgets.Toolbar.Items[3].ToolbarObject().Hide()
			AppWidgets.pageLabel.Hide()
		} else {
			AppWidgets.Toolbar.Items[2].ToolbarObject().Show()
			AppWidgets.Toolbar.Items[3].ToolbarObject().Show()
			AppWidgets.pageLabel.Show()
		}
		ShowNotesInGrid(main_app.AppStatus.Notes, main_app.AppStatus.NoteSize)
	case main_app.LAYOUT_PAGE:
		AppWidgets.Toolbar.Items[2].ToolbarObject().Show()
		AppWidgets.Toolbar.Items[3].ToolbarObject().Show()
		ShowNotesAsPages(main_app.AppStatus.Notes)
	default:
		err = errors.New("undefined layout")
	}

	return err
}

func CreateNotebooksList() {
	var err error
	main_app.AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}

	AppWidgets.notebooksList = widget.NewList(
		func() int {
			return len(main_app.AppStatus.Notebooks)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("------------Notebooks (xx)------------", func() {})

		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			main_app.AppStatus.Notes, _ = notes.GetNotebook(main_app.AppStatus.Notebooks[id])
			o.(*widget.Button).SetText(fmt.Sprintf("%s (%d)", main_app.AppStatus.Notebooks[id], len(main_app.AppStatus.Notes)))
			o.(*widget.Button).OnTapped = func() {
				//AppStatus.Notes,_ = minotedb.GetNotebook(AppStatus.Notebooks[id])
				main_app.AppStatus.CurrentView = main_app.VIEW_NOTEBOOK
				main_app.AppStatus.CurrentNotebook = main_app.AppStatus.Notebooks[id]
				PageView.Reset()
				UpdateView()
			}
		},
	)

}

func UpdateNotebooksList() {
	var err error
	main_app.AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}
	AppWidgets.notebooksList.Refresh()
}

func ShowNotebooksPanel() {
	UpdateNotebooksList()
	if main_app.AppStatus.CurrentView != main_app.VIEW_NOTEBOOK {
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
	mainWindow.Canvas().AddShortcut(main_app.ScViewPinned, func(shortcut fyne.Shortcut) {
		var err error
		main_app.AppStatus.CurrentView = main_app.VIEW_PINNED
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	//Keyboard shortcut to show Recent notes
	mainWindow.Canvas().AddShortcut(main_app.ScViewRecent, func(shortcut fyne.Shortcut) {
		var err error
		main_app.AppStatus.CurrentView = main_app.VIEW_RECENT
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	mainWindow.Canvas().AddShortcut(main_app.ScPageForward, func(shortcut fyne.Shortcut) {
		if PageView.PageForward() > 0 {
			UpdateView()
		}
	})

	mainWindow.Canvas().AddShortcut(main_app.ScPageBack, func(shortcut fyne.Shortcut) {
		if PageView.PageBack() > 0 {
			UpdateView()
		}
	})

	mainWindow.Canvas().AddShortcut(main_app.ScFind, func(shortcut fyne.Shortcut) {
		ShowSearchPanel()
	})

	//Keyboard shortcut to create a new note
	mainWindow.Canvas().AddShortcut(main_app.ScOpenNote, func(shortcut fyne.Shortcut) {
		//CreateNewNote()
	})

	//Keyboard shortcut to show notebooks list
	mainWindow.Canvas().AddShortcut(main_app.ScShowNotebooks, func(shortcut fyne.Shortcut) {
		ShowNotebooksPanel()
	})
}
