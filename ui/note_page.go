package ui

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/marcs100/minote/conversions"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
	"github.com/marcs100/minote/tracker"
)

func (np *NotePage) NewNotePage(retrievedNote *note.NoteData,
	allowEdit, newWindowMode bool,
	parentWindow fyne.Window,
	mainWindow *MainWindow) *fyne.Container {

	np.ParentWindow = parentWindow
	np.MainAppWindow = mainWindow
	np.AllowEdit = allowEdit
	np.NewWindowMode = newWindowMode
	np.RetrievedNote = *retrievedNote
	np.NoteInfo.NewNote = false
	np.UI_Colours = GetAppColours(mainWindow.ThemeVariant)
	if retrievedNote.Id == 0 {
		//New note
		np.NoteInfo = note.NoteInfo{
			Id:           0,
			Notebook:     retrievedNote.Notebook,
			DateCreated:  "",
			DateModified: "",
			Pinned:       false,
			PinnedDate:   "",
			Colour:       "#FFFFFF",
			Content:      "",
			Deleted:      false,
			NewNote:      true,
		}
	} else {
		np.NoteInfo = note.NoteInfo{
			Id:           retrievedNote.Id,
			Notebook:     retrievedNote.Notebook,
			DateCreated:  retrievedNote.Created,
			DateModified: retrievedNote.Modified,
			PinnedDate:   retrievedNote.PinnedDate,
			Colour:       retrievedNote.BackgroundColour,
			Content:      retrievedNote.Content,
			Deleted:      false,
			NewNote:      false,
		}

		if retrievedNote.Pinned > 0 {
			np.NoteInfo.Pinned = true
		} else {
			np.NoteInfo.Pinned = false
		}
	}

	if retrievedNote.Pinned > 0 {
		np.NoteInfo.Pinned = true
	} else {
		np.NoteInfo.Pinned = false
	}

	//calculate initial note content hash
	note.UpdateHash(&np.NoteInfo)

	//setup keyboard shortcuts
	np.NotePageWidgets.Entry = NewEntryCustom(func(cs *desktop.CustomShortcut) {
		switch cs.ShortcutName() {
		// case main_app.ScViewMode.ShortcutName():
		// 	np.SetViewMode()
		case main_app.ScPinNote.ShortcutName():
			np.PinNote()
		case main_app.ScNoteColour.ShortcutName():
			np.ChangeNoteColour()
		case main_app.ScShowInfo.ShortcutName():
			np.ShowProperties()
		case main_app.ScNoteTags.ShortcutName():
			np.ToggleTagsNotePanel()
		}
	},

		func(k *fyne.KeyEvent) { // Escape key detected
			np.SetViewMode()
		},

		func() { // On focus lost
			if !np.NoteInfo.Deleted {
				np.NoteInfo.Content = np.NotePageWidgets.Entry.Text
				fmt.Println("Focus lost will try and save note")
				np.SaveNote()
			}
		},
	)

	tooltip := canvas.NewText(fmt.Sprintf("%-25s", ""), conversions.RGBStringToFyneColor("#0ed6ea"))
	tooltip.TextStyle = fyne.TextStyle{Monospace: true}
	np.Tooltip = tooltip

	np.NotePageWidgets.Entry.Text = np.NoteInfo.Content
	np.NotePageWidgets.Entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(np.UI_Colours.NoteBgColour)
	topBarBackground := canvas.NewRectangle(np.UI_Colours.MainCtrlsBgColour)
	noteColour := conversions.RGBStringToFyneColor(np.NoteInfo.Colour)

	np.NotePageCanvas.NoteBackground = canvas.NewRectangle(noteColour)
	if np.NoteInfo.Colour == "#e7edef" || np.NoteInfo.Colour == "#FFFFFF" || np.NoteInfo.Colour == "#ffffff" || np.NoteInfo.Colour == "#000000" {
		np.NotePageCanvas.NoteBackground = canvas.NewRectangle(np.UI_Colours.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(np.NotePageCanvas.NoteBackground)

	np.NotePageWidgets.MarkdownText = widget.NewRichTextFromMarkdown(np.NoteInfo.Content)
	np.NotePageWidgets.MarkdownText.Wrapping = fyne.TextWrapWord
	np.NotePageWidgets.MarkdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, np.NotePageWidgets.MarkdownText)
	np.NotePageContainers.Markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(np.NotePageContainers.Markdown)
	background := canvas.NewRectangle(np.UI_Colours.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, np.NotePageWidgets.Entry)

	btnIcon := theme.RadioButtonIcon()
	if np.NoteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	np.NotePageWidgets.PinButton = NewButtonWithTooltip("", btnIcon, "Action: Pin/unpin note", tooltip, np.ParentWindow, func() {
		np.PinNote()
	})

	changeNotebookBtn := np.newChangeNotebookButton()

	colourButton := NewButtonWithTooltip("", theme.ColorPaletteIcon(), "Action: Change colour", tooltip, np.ParentWindow, func() {
		np.ChangeNoteColour()
	})

	np.NotePageWidgets.DeleteButton = NewButtonWithTooltip("", theme.DeleteIcon(), "Action: Delete Note", tooltip, np.ParentWindow, func() {
		np.DeleteNote()
	})

	np.NotePageWidgets.TagsButton = NewButtonWithTooltip("", theme.CheckButtonIcon(), "Action: Tags show/hide", tooltip, np.ParentWindow, func() {
		np.ToggleTagsNotePanel()
	})

	propertiesButton := NewButtonWithTooltip("", theme.InfoIcon(), "Action: Properties show/hide", tooltip, np.ParentWindow, func() { np.ShowProperties() })

	np.NotePageWidgets.DeleteButton.Hide()

	np.NotePageWidgets.ModeSelect = widget.NewRadioGroup([]string{main_app.EDIT_MODE, main_app.VIEW_MODE}, func(value string) {
		switch value {
		case main_app.EDIT_MODE:
			if allowEdit {
				np.SetEditMode()
			}
		case main_app.VIEW_MODE:
			np.SetViewMode()
		}
	})

	if !allowEdit {
		np.NotePageWidgets.ModeSelect.Hide()
	}

	CreateProperetiesPanel(np)

	np.NotePageWidgets.ModeSelect.SetSelected("View")
	np.NotePageWidgets.ModeSelect.Horizontal = true
	toolbar := container.NewHBox(np.NotePageWidgets.ModeSelect,
		spacerLabel, np.NotePageWidgets.PinButton,
		colourButton, changeNotebookBtn,
		np.NotePageWidgets.TagsButton,
		propertiesButton,
		np.NotePageWidgets.DeleteButton,
		tooltip,
	)

	if err := CreateNotesTagPanel(np); err != nil {
		dialog.ShowError(err, np.ParentWindow)
		//log.Panicln("Error creating tags panel")
	}
	np.UpdateTags()

	topVBox := container.NewVBox(toolbar, np.NotePageContainers.TagsPanel)
	topBar := container.NewStack(topBarBackground, topVBox)

	np.NotePageContainers.PropertiesPanel.Hide()
	np.NotePageContainers.TagsPanel.Hide()

	return container.NewBorder(topBar, nil, nil, np.NotePageContainers.PropertiesPanel, content)
}

func (np *NotePage) newChangeNotebookButton() *ButtonWithTooltip {
	changeNotebookBtn := NewButtonWithTooltip("", theme.FolderOpenIcon(), "Action: Select notebook", np.Tooltip, np.ParentWindow, func() {
		var notebooks []string
		var err error
		if notebooks, err = notes.GetNotebooks(); err != nil {
			log.Println("Error getting notebook")
			dialog.ShowError(err, np.ParentWindow)
			log.Panicln(err)
		}
		nbMenu := fyne.NewMenu("Select Notebook")

		//Add new notebook entry to menu
		nbMenuItem := fyne.NewMenuItem("*New*", func() {
			notebookEntry := widget.NewEntry()
			eNotebookEntry := widget.NewFormItem("Name", notebookEntry)
			newNotebookDialog := dialog.NewForm("New Notebook?", "OK", "Cancel", []*widget.FormItem{eNotebookEntry}, func(confirmed bool) {
				if confirmed {
					//check that the notebook does not already exist
					exists, err := notes.CheckNotebookExists(notebookEntry.Text)
					if err == nil {
						if exists == false {
							//chnage notebook to this new notebook
							np.NoteInfo.Notebook = notebookEntry.Text
							np.ParentWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", np.NoteInfo.Notebook, np.NoteInfo.Id))
							_, err = note.SaveNote(&np.NoteInfo)
							if err != nil {
								log.Print("Error saving note: ")
								dialog.ShowError(err, np.ParentWindow)
								//log.Panic(err)
							}
							np.MainAppWindow.UpdateNotebooksList()
							np.UpdateProperties()
						}
					} else {
						dialog.ShowError(err, np.ParentWindow)
						log.Panicf("Error check notebook exists: %s", err)
					}
				}
			}, np.ParentWindow)
			newNotebookDialog.Show()
		})

		nbMenu.Items = append(nbMenu.Items, nbMenuItem)

		//Now add all the existing notebooks to the menu
		for _, notebook := range notebooks {
			menuItem := fyne.NewMenuItem(notebook, func() {
				np.NoteInfo.Notebook = notebook
				//fmt.Println("Change notebook to " + notebook)
				np.ParentWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", np.NoteInfo.Notebook, np.NoteInfo.Id))
				np.UpdateProperties()
			})
			nbMenu.Items = append(nbMenu.Items, menuItem)
		}

		popUpMenu := widget.NewPopUpMenu(nbMenu, np.ParentWindow.Canvas())
		//popUpMenu.Show()
		pos := fyne.NewPos(250, 40)
		popUpMenu.ShowAtPosition(pos)
		//popUpMenu.ShowAtPosition(e.Position.AddXY(150,0))

	})

	//ch <- true
	return changeNotebookBtn
}

func (np *NotePage) DeleteNote() {
	dialog.ShowConfirm("Delete note", "Are you sure?", func(confirm bool) {
		if confirm {
			var res bool
			var err error = nil
			if np.NoteInfo.NewNote {
				res = true
			} else {
				res, err = notes.DeleteNote(np.NoteInfo.Id)
			}

			if res == false || err != nil {
				log.Println("Error deleting note - panic!")
				dialog.ShowError(err, np.ParentWindow)
				//log.Panicln(err)
			} else {
				np.NoteInfo.Deleted = true
				if np.NewWindowMode {
					np.ParentWindow.Close()
				}
				np.MainAppWindow.UpdateView()
			}
		}
	}, np.ParentWindow)
}

func (np *NotePage) PinNote() {
	var res bool
	var err error = nil
	if np.NoteInfo.Pinned {
		if np.NoteInfo.NewNote {
			//new note that hasn't been saved yet'
			res = true
		} else {
			res, err = notes.UnpinNote(np.NoteInfo.Id)
		}

		if err == nil && res == true {
			np.NoteInfo.Pinned = false
			np.NoteInfo.PinnedDate = ""
			if np.NotePageWidgets.PinButton != nil {
				np.NotePageWidgets.PinButton.SetIcon(theme.RadioButtonIcon())
				np.NotePageWidgets.PinButton.Refresh()
			}
		}
	} else {
		if np.NoteInfo.Id == 0 {
			//new note that hasn't been saved yet'
			np.NoteInfo.Pinned = true
			np.NoteInfo.PinnedDate = time.Now().String()[0:19]
			res = true
		} else {
			res, err = notes.PinNote(np.NoteInfo.Id)
			pinnedDate, err := notes.GetPinnedDate(np.NoteInfo.Id)
			if err != nil {
				log.Printf("Error getting pinned date: %s", err)
			}
			np.NoteInfo.PinnedDate = pinnedDate
		}
		if err == nil && res == true {
			np.NoteInfo.Pinned = true
			if np.NotePageWidgets.PinButton != nil {
				np.NotePageWidgets.PinButton.SetIcon(theme.RadioButtonCheckedIcon())
				np.NotePageWidgets.PinButton.Refresh()
			}
		}
	}

	np.UpdateProperties()

	if main_app.AppStatus.CurrentView == main_app.VIEW_PINNED {
		np.MainAppWindow.UpdateView() //updates view on main window c
	}
}

func (np *NotePage) SetEditMode() {
	np.NotePageWidgets.MarkdownText.Hide()
	np.NotePageContainers.Markdown.Hide()
	np.NotePageWidgets.DeleteButton.Show()
	np.TagsButtonDisplay()
	if main_app.AppStatus.CurrentLayout == main_app.LAYOUT_PAGE {
		//Hide page back & forward for edit mode
		//main_ui.AppWidgets.Toolbar.Items[2].ToolbarObject().Hide() // needs review
		//main_ui.AppWidgets.Toolbar.Items[3].ToolbarObject().Hide() // needs review
	}
	np.NotePageWidgets.ModeSelect.SetSelected(main_app.EDIT_MODE)
	np.NotePageWidgets.Entry.Show()
	np.ParentWindow.Canvas().Focus(np.NotePageWidgets.Entry)

}

func (np *NotePage) SetViewMode() {
	np.NotePageWidgets.Entry.Hide()
	np.NotePageWidgets.DeleteButton.Hide()
	np.TagsButtonDisplay()
	if main_app.AppStatus.CurrentLayout == main_app.LAYOUT_PAGE {
		//Show page back & forward for edit mode
		//AppWidgets.Toolbar.Items[2].ToolbarObject().Show() // needs review
		//AppWidgets.Toolbar.Items[3].ToolbarObject().Show() // needs review
	}
	np.NotePageWidgets.MarkdownText.ParseMarkdown(np.NotePageWidgets.Entry.Text)
	np.NotePageWidgets.MarkdownText.Show()
	np.ParentWindow.Canvas().Focus(nil) // this allows the canvas keyboard shortcuts to work rather than the entry widget shortcuts
	np.NotePageContainers.Markdown.Show()
	np.NotePageWidgets.ModeSelect.SetSelected(main_app.VIEW_MODE)
	if np.NewWindowMode {
		np.ParentWindow.Show()
	}
}

func (np *NotePage) ChangeNoteColour() {
	var colourChanged = false
	picker := dialog.NewColorPicker("Note Color", "Pick colour", func(c color.Color) {
		fmt.Println(c)
		hex := conversions.FyneColourToRGBHex(c)
		np.NoteInfo.Colour = fmt.Sprintf("%s%s", "#", hex)
		np.NotePageCanvas.NoteBackground.FillColor = c
		colourChanged = true
	}, np.ParentWindow)
	picker.Advanced = true
	picker.Show()
	if colourChanged {
		np.UpdateProperties()
		np.MainAppWindow.UpdateView()
	}
}

func (np *NotePage) SaveNote() {
	var noteChanges note.NoteChanges
	np.NoteInfo.Content = np.NotePageWidgets.Entry.Text
	if np.NoteInfo.Deleted {
		np.MainAppWindow.UpdateView()
		return
	}

	if np.NoteInfo.NewNote {
		if np.NoteInfo.Content != "" {
			noteChanges.ContentChanged = true
		}
	} else {
		noteChanges = note.CheckChanges(&np.RetrievedNote, &np.NoteInfo)
	}
	//if contentChanged{
	if noteChanges.ContentChanged || noteChanges.ParamsChanged {
		res, err := note.SaveNote(&np.NoteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, np.ParentWindow)
			//ch <- true
			return
		}

		if res == 0 {
			log.Println("No note was saved (affected rows = 0)")
		} else {
			log.Println("....Note updated successfully....")
			if np.RetrievedNote, err = notes.GetNote(np.NoteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, np.ParentWindow)
				return
			}
			//track new note as open.
			// Only works as new notes are always opned in a new window
			if np.NewWindowMode {
				tracker.AddToTracker(np.NoteInfo.Id)
			}

			np.MainAppWindow.UpdateView()
		}
	} else if noteChanges.PinStatusChanged {
		// we do not want a create or modified time stamp for just pinning/unpinning notes
		res, err := note.SaveNoteNoTimeStamp(&np.NoteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, np.ParentWindow)
			return
		}

		if res == 0 {
			log.Println("No note was saved (affected rows = 0)")
		} else {
			log.Println("....Note updated successfully....")
			if np.RetrievedNote, err = notes.GetNote(np.NoteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, np.ParentWindow)
			}
			np.MainAppWindow.UpdateView()
		}
	}
}

func (np *NotePage) AddNoteKeyboardShortcuts() {
	//Keyboard shortcut to set edit mode
	if np.AllowEdit {
		np.ParentWindow.Canvas().AddShortcut(main_app.ScEditMode, func(shortcut fyne.Shortcut) {
			np.SetEditMode()
		})
		//
		// np.ParentWindow.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		// 	fmt.Println(ke.Name)
		// 	if ke.Name == "I" {
		// 		ke = &fyne.KeyEvent{}
		// 		np.SetEditMode()
		// 	}
		// })
	}

	//Keyboard shortcut to pin/unpin notes
	np.ParentWindow.Canvas().AddShortcut(main_app.ScPinNote, func(shortcut fyne.Shortcut) {
		np.PinNote()
	})

	//Keyboard shortcut to change note colour
	np.ParentWindow.Canvas().AddShortcut(main_app.ScNoteColour, func(shortcut fyne.Shortcut) {
		np.ChangeNoteColour()
	})

	//Keyboard shortcut to show properties panel
	np.ParentWindow.Canvas().AddShortcut(main_app.ScShowInfo, func(shortcut fyne.Shortcut) {
		np.ShowProperties()
		np.RefreshWindow()
	})

	//Keyboard shortcut to show tags panel
	np.ParentWindow.Canvas().AddShortcut(main_app.ScNoteTags, func(shortcut fyne.Shortcut) {
		np.ToggleTagsNotePanel()
		np.RefreshWindow()
	})

}

func (np *NotePage) RefreshWindow() {
	np.ParentWindow.Content().Refresh()
}
