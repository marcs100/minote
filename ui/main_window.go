package ui

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
	"github.com/marcs100/minote/conversions"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
	"github.com/marcs100/minote/tracker"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartUI(appConfigIn *config.Config, configFile string, about main_app.About) {
	main_app.Conf = appConfigIn
	main_app.MainApp = app.NewWithID("minote")
	main_app.AppStatus.ConfigFile = configFile
	main_app.AppStatus.CurrentNotebook = "General" // default for new noteooks if note in notrbook view
	createMainWindow(about.Version)
	main_app.MainApp.Run()
}

func createMainWindow(version string) {

	var mw MainWindow

	main_app.AppStatus.NoteSize = fyne.NewSize(main_app.Conf.Settings.NoteWidth, main_app.Conf.Settings.NoteHeight)

	mw.window = main_app.MainApp.NewWindow(fmt.Sprintf("Minote   v%s", version))
	var themeVar ThemeVariant
	switch main_app.Conf.Settings.ThemeVariant {
	case "light":
		themeVar = LIGHT_THEME
	case "dark":
		themeVar = DARK_THEME
	case "auto":
		switch main_app.MainApp.Settings().ThemeVariant() {
		case theme.VariantDark:
			themeVar = DARK_THEME
		case theme.VariantLight:
			themeVar = LIGHT_THEME
		default:
			log.Println("Warning.. Could not auto detect theme variant, will default to dark theme!")
			themeVar = DARK_THEME
		}
	}

	mw.ThemeVariant = themeVar
	mw.UI_Colours = GetAppColours(themeVar)
	// fmt.Println("Will set a custom theme!")
	custTheme := &minoteTheme{
		FontSize:     main_app.Conf.Settings.FontSize,
		BgColour:     mw.UI_Colours.MainBgColour,
		EntryColour:  mw.UI_Colours.NoteBgColour,
		ButtonColour: mw.UI_Colours.ButtonColour,
		AccentColour: mw.UI_Colours.AccentColour,
		FgColour:     mw.UI_Colours.MainFgColour,
	}
	main_app.MainApp.Settings().SetTheme(custTheme)

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(main_app.AppStatus.NoteSize)
	mw.AppContainers.grid = grid //store to allow interaction in other functions

	singleNotePage := widget.NewRichTextFromMarkdown("")
	mw.AppWidgets.singleNotePage = singleNotePage
	singleNoteStack := container.NewStack()
	mw.AppContainers.singleNoteStack = singleNoteStack

	PageView.CurrentPage = 0
	PageView.NumberOfPages = 0

	mw.Tooltip = canvas.NewText(fmt.Sprintf("%-25s", ""), mw.UI_Colours.AccentColour)
	mw.Tooltip.TextStyle = fyne.TextStyle{Monospace: true}

	//Create The main panel
	main := mw.createMainPanel()

	top := mw.createTopPanel()

	side := mw.createSidePanel()

	//layout the main window
	appContainer := container.NewBorder(top, nil, side, nil, main)

	mw.window.SetContent(appContainer)
	mw.window.Resize(fyne.NewSize(1800, 1100))

	//set default view and layout`
	main_app.AppStatus.CurrentView = main_app.Conf.Settings.InitialView
	fmt.Println("initial view = " + main_app.Conf.Settings.InitialView)
	main_app.AppStatus.CurrentLayout = main_app.Conf.Settings.InitialLayout

	mw.setSortOptions(main_app.AppStatus.CurrentView)
	mw.AppWidgets.sortSelect.SetSelectedIndex(0) // this will also initiate UpdateView()

	// if err := UpdateView(); err != nil {
	// 	fmt.Println(err)
	// }

	//keyboard shortcuts
	mw.addMainKeyboardShortcuts()

	mw.window.SetCloseIntercept(func() {
		if tracker.TrackerLen() > 0 {
			fmt.Printf("len of open notes array is %d", tracker.TrackerLen())
			//do not close if there are notes open
			dlg := dialog.NewInformation("Error", "There are notes open, please close them before closing the application!", mw.window)
			dlg.Show()
		} else {
			mw.window.Close()
		}
	})

	mw.window.SetMaster()
	mw.window.Show()
}

func (mw *MainWindow) createMainPanel() *fyne.Container {

	mainGridContainer := container.NewScroll(mw.AppContainers.grid)
	mw.AppContainers.mainGridContainer = mainGridContainer
	mainPageContainer := container.NewScroll(mw.AppContainers.singleNoteStack)
	mw.AppContainers.mainPageContainer = mainPageContainer

	bgRect := canvas.NewRectangle(mw.UI_Colours.MainBgColour)
	mainStackedContainer := container.NewStack(bgRect, mainPageContainer, mainGridContainer)
	return mainStackedContainer
}

func (mw *MainWindow) createTopPanel() *fyne.Container {
	//AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	spacerLabel := widget.NewLabel("  ")
	mw.AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes      >", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mw.AppWidgets.pageLabel = widget.NewLabel("Page: ")
	mw.AppWidgets.pageLabel.Hide()

	gridViewBtn := NewButtonWithTooltip("", theme.GridIcon(), fmt.Sprintf("%-25s", "Grid View"), mw.Tooltip, mw.window, func() {
		if main_app.AppStatus.CurrentLayout != main_app.LAYOUT_GRID {
			main_app.AppStatus.CurrentLayout = main_app.LAYOUT_GRID
			PageView.Reset()
			mw.UpdateView()
		}
	})

	pageViewBtn := NewButtonWithTooltip("", theme.FileAudioIcon(), fmt.Sprintf("%-25s", "Single page View"), mw.Tooltip, mw.window, func() {
		if main_app.AppStatus.CurrentLayout != main_app.LAYOUT_PAGE {
			main_app.AppStatus.CurrentLayout = main_app.LAYOUT_PAGE
			PageView.Reset()
			mw.UpdateView()
		}
	})

	mw.AppWidgets.pageBackBtn = NewButtonWithTooltip("", theme.NavigateBackIcon(), fmt.Sprintf("%-25s", "Page back"), mw.Tooltip, mw.window, func() {
		if PageView.PageBack() > 0 {
			mw.UpdateView()
		}
	})
	mw.AppWidgets.pageBackBtn.Hidden = true

	mw.AppWidgets.pageForwardBtn = NewButtonWithTooltip("", theme.NavigateNextIcon(), fmt.Sprintf("%-25s", "Page forward"), mw.Tooltip, mw.window, func() {
		if PageView.PageForward() > 0 {
			mw.UpdateView()
		}
	})
	mw.AppWidgets.pageForwardBtn.Hidden = true

	viewsHbox := container.NewHBox(gridViewBtn, pageViewBtn, mw.AppWidgets.pageBackBtn, mw.AppWidgets.pageForwardBtn)

	rightToolbar := widget.NewToolbar(
		//backup database
		// widget.NewToolbarAction(theme.DownloadIcon(), func() {
		// 	BackupNotes(main_app.Conf.Settings.Database, mw.window)
		// }),

		//display settings
		widget.NewToolbarAction(theme.MenuIcon(), func() {
			options := fyne.NewMenu("Options")
			settingsMenuItem := fyne.NewMenuItem("Settings", func() {
				ShowSettings(mw.window, mw.UI_Colours)
			})

			backupMenuItem := fyne.NewMenuItem("Backup Notes", func() {
				BackupNotes(main_app.Conf.Settings.Database, mw.window)
			})

			options.Items = append(options.Items, settingsMenuItem, backupMenuItem)

			popUpMenu := widget.NewPopUpMenu(options, mw.window.Canvas())
			pos := fyne.NewPos(225, 40)
			// popUpMenu.ShowAtPosition(pos)
			popUpMenu.ShowAtRelativePosition(pos, mw.AppWidgets.sortSelect)

			//current option - will remove - replace with menu
			//ShowSettings(mw.window, mw.UI_Colours)
		}),
	)

	sortLabel := widget.NewLabel("sort:")
	mw.AppWidgets.sortSelect = widget.NewSelect([]string{""}, func(s string) {
		main_app.AppStatus.CurrentSortSelected = s
		mw.UpdateView()
	})
	mw.AppWidgets.sortSelect.PlaceHolder = "This is how the size "
	mw.AppWidgets.sortSelect.SetSelectedIndex(0)

	//rect := canvas.NewRectangle(UI_Colours.)
	topBar := container.New(layout.NewHBoxLayout(),
		mw.Tooltip,
		spacerLabel,
		mw.AppWidgets.viewLabel,
		layout.NewSpacer(),
		viewsHbox,
		mw.AppWidgets.pageLabel,
		layout.NewSpacer(),
		sortLabel,
		mw.AppWidgets.sortSelect,
		spacerLabel,
		rightToolbar,
	)

	rect := canvas.NewRectangle(mw.UI_Colours.MainCtrlsBgColour)
	topPanel := container.NewStack(rect, topBar)
	return topPanel

}

func (mw *MainWindow) createSidePanel() *fyne.Container {
	mw.AppContainers.searchPanel = mw.CreateSearchPanel()
	mw.AppContainers.tagsPanel = mw.CreateTagsPanel()

	newNoteBtn := NewButtonWithTooltip("+", theme.DocumentCreateIcon(), fmt.Sprintf("%-25s", "Action: New note"), mw.Tooltip, mw.window, func() {
		NewNoteWindow(0, mw.window, mw)
	})

	searchBtn := NewButtonWithTooltip("", theme.SearchIcon(), fmt.Sprintf("%-25s", "Action: Search"), mw.Tooltip, mw.window, func() {
		//Display the search panel here
		main_app.AppStatus.CurrentView = main_app.VIEW_SEARCH
		mw.setSortOptions(main_app.VIEW_SEARCH)
		mw.AppWidgets.sortSelect.SetSelectedIndex(0)
		mw.ShowSearchPanel()
	})

	//pinnedBtn := widget.NewButton("P", func(){
	pinnedBtn := NewButtonWithTooltip("", theme.RadioButtonCheckedIcon(), fmt.Sprintf("%-25s", "Action: Show pinned notes"), mw.Tooltip, mw.window, func() {
		main_app.AppStatus.CurrentView = main_app.VIEW_PINNED
		mw.setSortOptions(main_app.AppStatus.CurrentView)
		PageView.Reset()
		mw.AppWidgets.sortSelect.SetSelectedIndex(0)
	})

	RecentBtn := NewButtonWithTooltip("", theme.HistoryIcon(), fmt.Sprintf("%-25s", "Action: Show recent notes"), mw.Tooltip, mw.window, func() {
		//AppStatus.Notes,err = minotedb.GetRecentNotes(Conf.Settings.RecentNotesLimit)
		main_app.AppStatus.CurrentView = main_app.VIEW_RECENT
		PageView.Reset()
		mw.setSortOptions(main_app.AppStatus.CurrentView)
		mw.AppWidgets.sortSelect.SetSelectedIndex(0)
	})

	tagsBtn := NewButtonWithTooltip("", theme.CheckButtonIcon(), fmt.Sprintf("%-25s", "Action: Show tags"), mw.Tooltip, mw.window, func() {
		mw.ToggleMainTagsPanel()
		main_app.AppStatus.CurrentView = main_app.VIEW_TAGS
		PageView.Reset()
		mw.setSortOptions(main_app.VIEW_TAGS)
		mw.AppWidgets.sortSelect.SetSelectedIndex(0)
	})

	mw.CreateNotebooksList()

	notebooksBtn := NewButtonWithTooltip("", theme.FolderOpenIcon(), fmt.Sprintf("%-25s", "Action: Show notebooks"), mw.Tooltip, mw.window, func() {
		main_app.AppStatus.CurrentView = main_app.VIEW_NOTEBOOK
		mw.setSortOptions(main_app.AppStatus.CurrentView)
		mw.AppWidgets.sortSelect.SetSelectedIndex(0)
		mw.showNotebooksPanel()
	})

	spacerLabel := widget.NewLabel(" ")

	btnPanel := container.NewVBox(searchBtn, newNoteBtn, spacerLabel, pinnedBtn, RecentBtn, notebooksBtn, tagsBtn)
	mw.AppContainers.listPanel = container.NewStack(mw.AppWidgets.notebooksList)
	mw.AppContainers.listPanel.Hide()
	mw.AppContainers.tagsPanel.Hide()

	sidePanel := container.NewHBox(btnPanel, mw.AppContainers.listPanel, mw.AppContainers.searchPanel, mw.AppContainers.tagsPanel)

	rect := canvas.NewRectangle(mw.UI_Colours.MainCtrlsBgColour)
	sideContainer := container.NewStack(rect, sidePanel)
	return sideContainer
}

func (mw *MainWindow) showNotesInGrid(notes []note.NoteData) {
	if mw.AppContainers.grid == nil || mw.AppContainers.mainGridContainer == nil {
		return
	}

	if mw.AppContainers.mainPageContainer != nil {
		mw.AppContainers.mainPageContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	PageView.Step = main_app.Conf.Settings.GridMaxPages
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	mw.AppContainers.grid.RemoveAll()
	numPages := (PageView.CurrentPage + PageView.Step) - 1
	if numPages > len(notes) {
		numPages = PageView.NumberOfPages
	}

	if mw.AppWidgets.pageLabel.Hidden != true {
		mw.AppWidgets.pageLabel.SetText(PageView.GetGridLabelText())
	}

	for i := PageView.CurrentPage - 1; i < numPages; i++ {
		richText := NewMarkdownCustom(notes[i].Content, func() {
			if tracker.AddToTracker(notes[i].Id) {
				NewNoteWindow(notes[i].Id, mw.window, mw)
			} else {
				fmt.Println("note is already open")
			}
		})
		richText.Wrapping = fyne.TextWrapWord
		themeBackground := canvas.NewRectangle(mw.UI_Colours.NoteBgColour) //We need this to overlay the notecolour (border)
		noteColour := conversions.RGBStringToFyneColor(notes[i].BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if notes[i].BackgroundColour == "#e7edef" || notes[i].BackgroundColour == "#FFFFFF" || notes[i].BackgroundColour == "#ffffff" || notes[i].BackgroundColour == "#000000" {
			noteBackground = canvas.NewRectangle(mw.UI_Colours.NoteBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack := container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		mw.AppContainers.grid.Add(noteStack)
	}
	mw.AppContainers.grid.Refresh()
	mw.AppContainers.mainGridContainer.Show()
}

func (mw *MainWindow) showNotesAsPages(notesIn []note.NoteData) {
	//var noteInfo note.NoteInfo
	var retrievedNote note.NoteData
	var err error = nil

	if mw.AppContainers.mainGridContainer != nil {
		mw.AppContainers.mainGridContainer.Hide()
	}

	PageView.NumberOfPages = len(notesIn)
	PageView.Step = 1
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	if PageView.NumberOfPages == 0 {
		return
	}

	var noteId = notesIn[PageView.CurrentPage-1].Id

	mw.AppWidgets.pageLabel.SetText(PageView.GetLabelText())
	mw.AppWidgets.pageLabel.Show()

	mw.AppContainers.singleNoteStack.RemoveAll()

	if noteId != 0 {
		retrievedNote, err = notes.GetNote(noteId)
		if err != nil {
			dialog.ShowError(err, mw.window)
			log.Panic(err)
		}
	}

	var allowEdit bool = true
	if tracker.TrackerCheck(noteId) {
		dialog.ShowInformation("Warning", "This note is already open in a separate window.\nClose it first if you want to edit it here!", mw.window)
		allowEdit = false
	}

	var np NotePage
	noteContainer := np.NewNotePage(&retrievedNote, allowEdit, false, mw.window, mw)
	np.NotePageWidgets.ModeSelect.SetSelected("View")
	np.AddNoteKeyboardShortcuts()
	mw.AppContainers.singleNoteStack.Add(noteContainer)
	mw.AppContainers.mainPageContainer.Show()
	mw.AppContainers.mainPageContainer.Refresh()
}

func (mw *MainWindow) UpdateView() error {
	//var notes []minotedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch main_app.AppStatus.CurrentView {
	case main_app.VIEW_PINNED:
		if mw.AppContainers.listPanel != nil {
			mw.AppContainers.listPanel.Hide()
		}
		if mw.AppContainers.searchPanel != nil {
			mw.AppContainers.searchPanel.Hide()
		}
		if mw.AppContainers.tagsPanel != nil {
			mw.AppContainers.tagsPanel.Hide()
		}
		mw.AppWidgets.viewLabel.SetText("Pinned Notes")
		main_app.AppStatus.Notes, err = notes.GetPinnedNotes(SortViews[main_app.AppStatus.CurrentSortSelected])
		main_app.AppStatus.CurrentNotebook = ""
	case main_app.VIEW_RECENT:
		if mw.AppContainers.listPanel != nil {
			mw.AppContainers.listPanel.Hide()
		}
		if mw.AppContainers.searchPanel != nil {
			mw.AppContainers.searchPanel.Hide()
		}
		if mw.AppContainers.tagsPanel != nil {
			mw.AppContainers.tagsPanel.Hide()
		}

		mw.AppWidgets.viewLabel.SetText(("Recent Notes"))
		main_app.AppStatus.Notes, err = notes.GetRecentNotes(main_app.Conf.Settings.RecentNotesLimit, SortViews[main_app.AppStatus.CurrentSortSelected])
		main_app.AppStatus.CurrentNotebook = ""
	case main_app.VIEW_NOTEBOOK:
		if mw.AppContainers.tagsPanel != nil {
			mw.AppContainers.tagsPanel.Hide()
		}
		mw.AppWidgets.viewLabel.SetText("Notebook - " + main_app.AppStatus.CurrentNotebook)
		main_app.AppStatus.Notes, err = notes.GetNotebook(main_app.AppStatus.CurrentNotebook, SortViews[main_app.AppStatus.CurrentSortSelected])
	case main_app.VIEW_TAGS:
		if mw.AppContainers.listPanel != nil {
			mw.AppContainers.listPanel.Hide()
		}
		if mw.AppContainers.searchPanel != nil {
			mw.AppContainers.searchPanel.Hide()
		}
		//if len(AppStatus.TagsChecked) > 0 {
		mw.AppWidgets.viewLabel.SetText("Tagged Notes")
		main_app.AppStatus.Notes, err = notes.GetTaggedNotes(main_app.AppStatus.TagsChecked, SortViews[main_app.AppStatus.CurrentSortSelected])
		//} else {
		//	AppStatus.Notes = nil
		//}
	case main_app.VIEW_SEARCH:
		if mw.AppContainers.tagsPanel != nil {
			mw.AppContainers.tagsPanel.Hide()
		}
		if len(strings.TrimSpace(mw.AppWidgets.searchEntry.Text)) > 0 {
			main_app.AppStatus.Notes, err = notes.GetSearchResults(mw.AppWidgets.searchEntry.Text, main_app.AppStatus.SearchFilter, SortViews[main_app.AppStatus.CurrentSortSelected])
			if err == nil {
				mw.AppWidgets.searchResultsLabel.SetText(fmt.Sprintf("Found (%d) > ", len(main_app.AppStatus.Notes)))
				mw.AppWidgets.viewLabel.SetText("Search Results")
			}
		}

	default:
		err = errors.New("undefined view")
	}

	if err != nil {
		return err
	}
	err = mw.showCurrentLayout()

	return err
}

func (mw *MainWindow) showCurrentLayout() error {
	var err error = nil
	switch main_app.AppStatus.CurrentLayout {
	case main_app.LAYOUT_GRID:
		if len(main_app.AppStatus.Notes) <= main_app.Conf.Settings.GridMaxPages {
			mw.AppWidgets.pageBackBtn.Hidden = true
			mw.AppWidgets.pageForwardBtn.Hidden = true
			mw.AppWidgets.pageLabel.Hide()
		} else {
			mw.AppWidgets.pageBackBtn.Hidden = false
			mw.AppWidgets.pageForwardBtn.Hidden = false
			mw.AppWidgets.pageLabel.Show()
		}
		mw.showNotesInGrid(main_app.AppStatus.Notes)
	case main_app.LAYOUT_PAGE:
		mw.AppWidgets.pageBackBtn.Hidden = false
		mw.AppWidgets.pageForwardBtn.Hidden = false
		mw.showNotesAsPages(main_app.AppStatus.Notes)
	default:
		err = errors.New("undefined layout")
	}

	return err
}

func (mw *MainWindow) CreateNotebooksList() {
	var err error
	main_app.AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mw.window)
		log.Panic(err)
	}

	mw.AppWidgets.notebooksList = widget.NewList(
		func() int {
			return len(main_app.AppStatus.Notebooks)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("------------Notebooks (xx)------------", func() {})

		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			main_app.AppStatus.Notes, _ = notes.GetNotebook(main_app.AppStatus.Notebooks[id], SortViews[main_app.AppStatus.CurrentSortSelected])
			var name = main_app.AppStatus.Notebooks[id]
			if len(name) > 15 {
				name = fmt.Sprint(name[:len(name)-3], "...")
			}
			o.(*widget.Button).SetText(fmt.Sprintf("%s (%d)", name, len(main_app.AppStatus.Notes)))
			o.(*widget.Button).OnTapped = func() {
				//AppStatus.Notes,_ = minotedb.GetNotebook(AppStatus.Notebooks[id])
				main_app.AppStatus.CurrentView = main_app.VIEW_NOTEBOOK
				main_app.AppStatus.CurrentNotebook = main_app.AppStatus.Notebooks[id]
				PageView.Reset()
				mw.setSortOptions(main_app.VIEW_NOTEBOOK)
				mw.AppWidgets.sortSelect.SetSelectedIndex(0) //this will also trigger UpdateView()
			}
		},
	)

}

func (mw *MainWindow) UpdateNotebooksList() {
	var err error
	main_app.AppStatus.Notebooks, err = minotedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mw.window)
		log.Panic(err)
	}
	mw.AppWidgets.notebooksList.Refresh()
}

func (mw *MainWindow) showNotebooksPanel() {
	mw.UpdateNotebooksList()
	if main_app.AppStatus.CurrentView != main_app.VIEW_NOTEBOOK {
		mw.AppWidgets.viewLabel.SetText("Notebooks")
	}

	if mw.AppContainers.listPanel != nil {
		if mw.AppContainers.listPanel.Visible() {
			mw.AppContainers.listPanel.Hide()
		} else {
			mw.AppContainers.listPanel.Show()
		}
	}
}

func (mw *MainWindow) addMainKeyboardShortcuts() {
	//Keyboard shortcut to show Pinned Notes
	mw.window.Canvas().AddShortcut(main_app.ScViewPinned, func(shortcut fyne.Shortcut) {
		var err error
		main_app.AppStatus.CurrentView = main_app.VIEW_PINNED
		PageView.Reset()
		err = mw.UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mw.window)
			log.Panic(err)
		}
	})

	//Keyboard shortcut to show Recent notes
	mw.window.Canvas().AddShortcut(main_app.ScViewRecent, func(shortcut fyne.Shortcut) {
		var err error
		main_app.AppStatus.CurrentView = main_app.VIEW_RECENT
		PageView.Reset()
		err = mw.UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mw.window)
			log.Panic(err)
		}
	})

	mw.window.Canvas().AddShortcut(main_app.ScPageForward, func(shortcut fyne.Shortcut) {
		if PageView.PageForward() > 0 {
			mw.UpdateView()
		}
	})

	mw.window.Canvas().AddShortcut(main_app.ScPageBack, func(shortcut fyne.Shortcut) {
		if PageView.PageBack() > 0 {
			mw.UpdateView()
		}
	})

	mw.window.Canvas().AddShortcut(main_app.ScFind, func(shortcut fyne.Shortcut) {
		mw.ShowSearchPanel()
	})

	//Keyboard shortcut to create a new note
	mw.window.Canvas().AddShortcut(main_app.ScOpenNote, func(shortcut fyne.Shortcut) {
		//CreateNewNote()
	})

	//Keyboard shortcut to show notebooks list
	mw.window.Canvas().AddShortcut(main_app.ScShowNotebooks, func(shortcut fyne.Shortcut) {
		mw.showNotebooksPanel()
	})

	//Keyboard shortcut to show notebooks list
	mw.window.Canvas().AddShortcut(main_app.ScShowTags, func(shortcut fyne.Shortcut) {
		mw.ToggleMainTagsPanel()
		main_app.AppStatus.CurrentView = main_app.VIEW_TAGS
		PageView.Reset()
		err := mw.UpdateView()
		if err != nil {
			log.Print("Error getting tagged notes: ")
			dialog.ShowError(err, mw.window)
			log.Panic(err)
		}

	})
}

func (mw *MainWindow) setSortOptions(view string) {

	switch view {
	// note - these options must match values in the sortViews map
	case main_app.VIEW_PINNED:
		mw.AppWidgets.sortSelect.Options = []string{"Pinned: most recently", "Pinned: least recently", "Modified: new to old", "Modified: old to new"}
	case main_app.VIEW_RECENT:
		mw.AppWidgets.sortSelect.Options = []string{"Modified: new to old", "Modified: old to new", "Created: new to old", "Created: old to new"}
	case main_app.VIEW_NOTEBOOK:
		mw.AppWidgets.sortSelect.Options = []string{"Modified: new to old", "Modified: old to new", "Created: new to old", "Created: old to new"}
	case main_app.VIEW_TAGS:
		mw.AppWidgets.sortSelect.Options = []string{"Modified: new to old", "Modified: old to new", "Created: new to old", "Created: old to new"}
	case main_app.VIEW_SEARCH:
		mw.AppWidgets.sortSelect.Options = []string{"Modified: new to old", "Modified: old to new", "Created: new to old", "Created: old to new"}

	}

	mw.AppWidgets.sortSelect.ClearSelected()
	mw.AppWidgets.sortSelect.Refresh()
}
