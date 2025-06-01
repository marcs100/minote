package note_ui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	//"github.com/marcs100/minote/main_ui"
	"github.com/marcs100/minote/conversions"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
)

func (np *NotePage) NewNotePage(retrievedNote *note.NoteData, allowEdit bool, parentWindow fyne.Window) *fyne.Container {
	np.ParentWindow = parentWindow
	np.AllowEdit = allowEdit
	np.NoteInfo.NewNote = false
	if retrievedNote == nil {
		//New note
		np.NoteInfo = note.NoteInfo{
			Id:           0,
			Notebook:     "General",
			DateCreated:  "",
			DateModified: "",
			Pinned:       false,
			PinnedDate:   "",
			Colour:       "#FFFFFF",
			Content:      "",
			Deleted:      false,
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
	//
	// FOR EACH OD THESE FUNCTIONS BELOW, I WIIL NEED TO PASS an INSTANCE OF THIS STRUCT
	// AS THE STRUCT WILL HOLD THE DATA ELEMENTS THAT WERE ACCESSD BY e.g., NOTEWIDGETS.blah IN SCRIBE-NB
	//ch := make(chan bool)
	np.NotePageWidgets.Entry = NewEntryCustom(func(cs *desktop.CustomShortcut) {
		switch cs.ShortcutName() {
		case main_app.ScViewMode.ShortcutName():
			np.SetViewMode()
		case main_app.ScPinNote.ShortcutName():
			np.PinNote()
		case main_app.ScNoteColour.ShortcutName():
			np.ChangeNoteColour()
		case main_app.ScShowInfo.ShortcutName():
			ShowProperties(np.NoteInfo)
		}
	}, func() {
		np.NoteInfo.Content = np.NotePageWidgets.Entry.Text
		SaveNote(np, retrievedNote)
		//<-ch
	},
	)

	np.NotePageWidgets.Entry.Text = np.NoteInfo.Content
	np.NotePageWidgets.Entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(main_app.AppTheme.NoteBgColour)
	noteColour, _ := conversions.RGBStringToFyneColor(np.NoteInfo.Colour)

	np.NotePageCanvas.NoteBackground = canvas.NewRectangle(noteColour)
	if np.NoteInfo.Colour == "#e7edef" || np.NoteInfo.Colour == "#FFFFFF" || np.NoteInfo.Colour == "#ffffff" || np.NoteInfo.Colour == "#000000" {
		np.NotePageCanvas.NoteBackground = canvas.NewRectangle(main_app.AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(np.NotePageCanvas.NoteBackground)

	np.NotePageWidgets.MarkdownText = widget.NewRichTextFromMarkdown(np.NoteInfo.Content)
	np.NotePageWidgets.MarkdownText.Wrapping = fyne.TextWrapWord
	np.NotePageWidgets.MarkdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, np.NotePageWidgets.MarkdownText)
	np.NotePageContainers.Markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(np.NotePageContainers.Markdown)
	background := canvas.NewRectangle(main_app.AppTheme.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, np.NotePageWidgets.Entry)

	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if np.NoteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	np.NotePageWidgets.PinButton = widget.NewButtonWithIcon("", btnIcon, func() {
		np.PinNote()
	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	//ch = make(chan bool)
	changeNotebookBtn := NewChangeNotebookButton(np)

	colourButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		np.ChangeNoteColour()
	})

	np.NotePageWidgets.DeleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteNote(np)
	})

	tagsBtn := widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func() {
		ToggleTagsNotePanel()
	})

	propertiesButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() { ShowProperties(&np.NoteInfo) })

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

	np.NotePageContainers.PropertiesPanel = NewProperetiesPanel()

	np.NotePageWidgets.ModeSelect.SetSelected("View")
	np.NotePageWidgets.ModeSelect.Horizontal = true
	toolbar := container.NewHBox(np.NotePageWidgets.ModeSelect, spacerLabel, np.NotePageWidgets.PinButton, colourButton, changeNotebookBtn, tagsBtn, propertiesButton, np.NotePageWidgets.DeleteButton)

	if err := CreateNotesTagPanel(&np.NoteInfo, parentWindow); err != nil {
		dialog.ShowError(err, parentWindow)
		log.Panicln("Error creating tags panel!")
	}
	topVBox := container.NewVBox(toolbar, np.NotePageContainers.TagsPanel)

	np.NotePageContainers.PropertiesPanel.Hide()
	np.NotePageContainers.TagsPanel.Hide()

	return container.NewBorder(topVBox, nil, nil, np.NotePageContainers.PropertiesPanel, content)
}

// Thread safe function
//var chn_mut sync.Mutex

func NewChangeNotebookButton(np *NotePage) *widget.Button {
	//chn_mut.Lock()
	//defer chn_mut.Unlock()
	changeNotebookBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
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
			//fmt.Println("Need to ask use for new notebook name here!!!!!!!!")
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
							//main_ui.UpdateNotebooksList() // ******* NEED THIS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
							//UpdateProperties(np.NoteInfo) // ******* NEED THIS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
						}
					} else {
						dialog.ShowError(err, np.ParentWindow)
						log.Panicln(fmt.Sprintf("Error check notebook exists: %s", err))
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
				//UpdateProperties(&np.NoteInfo) ********* NEED This!!!!!!!!!!!!!
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

func DeleteNote(np *NotePage) {
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
				np.ParentWindow.Close()
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
			np.NoteInfo.Pinned = false
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
			res = true
		} else {
			res, err = notes.PinNote(np.NoteInfo.Id)
			pinnedDate, err := notes.GetPinnedDate(np.NoteInfo.Id)
			if err != nil {
				log.Println(fmt.Sprintf("Error getting pinned date: %s", err))
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

	//UpdateProperties(np.NoteInfo) // ******* NEED THIS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	if main_app.AppStatus.CurrentView == main_app.VIEW_PINNED {
		//UpdateView() //updates view on main window c
	}
}

func (np *NotePage) SetEditMode() {
	np.NotePageWidgets.MarkdownText.Hide()
	np.NotePageContainers.Markdown.Hide()
	np.NotePageWidgets.DeleteButton.Show()
	if main_app.AppStatus.CurrentLayout == main_app.LAYOUT_PAGE {
		//Hide page back & forward for edit mode
		//main_ui.AppWidgets.Toolbar.Items[2].ToolbarObject().Hide() // needs review
		//main_ui.AppWidgets.Toolbar.Items[3].ToolbarObject().Hide() // needs review
	}
	np.NotePageWidgets.ModeSelect.SetSelected(main_app.EDIT_MODE)
	np.NotePageWidgets.Entry.Show()
	np.ParentWindow.Canvas().Focus(np.NotePageWidgets.Entry)
	np.ParentWindow.Content().Refresh()
}

func (np *NotePage) SetViewMode() {
	np.NotePageWidgets.Entry.Hide()
	np.NotePageWidgets.DeleteButton.Hide()
	if main_app.AppStatus.CurrentLayout == main_app.LAYOUT_PAGE {
		//Show page back & forward for edit mode
		//AppWidgets.Toolbar.Items[2].ToolbarObject().Show() // needs review
		//AppWidgets.Toolbar.Items[3].ToolbarObject().Show() // needs review
	}
	np.NotePageWidgets.MarkdownText.ParseMarkdown(np.NotePageWidgets.Entry.Text)
	np.NotePageWidgets.MarkdownText.Show()
	np.NotePageWidgets.ModeSelect.SetSelected(main_app.VIEW_MODE)
	np.ParentWindow.Canvas().Focus(nil) // this allows the canvas keyboard shortcuts to work rather than the entry widget shortcuts
	np.NotePageContainers.Markdown.Show()
}

func (np *NotePage) ChangeNoteColour() {
	picker := dialog.NewColorPicker("Note Color", "Pick colour", func(c color.Color) {
		fmt.Println(c)
		hex := conversions.FyneColourToRGBHex(c)
		np.NoteInfo.Colour = fmt.Sprintf("%s%s", "#", hex)
		np.NotePageCanvas.NoteBackground.FillColor = c
	}, np.ParentWindow)
	picker.Advanced = true
	picker.Show()
	//UpdateProperties(np) // ************** NEED THIS **********************************
}

//Pass a pointer to the note page - *NotePage
// Make sure this functiion is thread safe, as multiple note instances can be calling this function

//var sn_mut sync.Mutex

func SaveNote(np *NotePage, retrievedNote *note.NoteData) {

	//sn_mut.Lock()
	//defer sn_mut.Unlock()

	var noteChanges note.NoteChanges

	if np.NoteInfo.Deleted {
		//main_ui.UpdateView() // ************** NEED THIS *******************************
		//ch <- true
		return
	}

	if np.NoteInfo.NewNote {
		if np.NoteInfo.Content != "" {
			noteChanges.ContentChanged = true
		}
	} else {
		noteChanges = note.CheckChanges(retrievedNote, &np.NoteInfo)
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
			if *retrievedNote, err = notes.GetNote(np.NoteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, np.ParentWindow)
			}
			//main_ui.UpdateView() // *********** NEED THIS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		}
	} else if noteChanges.PinStatusChanged {
		// we do not want a create or modified time stamp for just pinning/unpinning notes
		res, err := note.SaveNoteNoTimeStamp(&np.NoteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, np.ParentWindow)
			//ch <- true
			return
			//log.Panic()
		}

		if res == 0 {
			log.Println("No note was saved (affected rows = 0)")
		} else {
			log.Println("....Note updated successfully....")
			if *retrievedNote, err = notes.GetNote(np.NoteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, np.ParentWindow)
			}
			//main_ui.UpdateView() // ******************** NEED THIS !!!!!!!!!!!!!!!!!!!!!!!
		}
	}
	//ch <- true
}

func (np *NotePage) AddNoteKeyboardShortcuts() {
	//Keyboard shortcut to set edit mode
	if np.AllowEdit {
		np.ParentWindow.Canvas().AddShortcut(main_app.ScEditMode, func(shortcut fyne.Shortcut) {
			np.SetEditMode()
		})
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
		//ShowProperties(noteInfo) // ***********  NEED THIS ****************************
	})
}

func NoteContainerRefresh() {
	//if noteContainer != nil { // ************ NEED THIS *******************************
	//	noteContainer.Refresh()
	//}
}
