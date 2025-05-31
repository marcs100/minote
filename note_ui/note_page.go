package note_ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_ui"
	"github.com/marcs100/minote/note"
	"github.com/marcs100/minote/notes"
)

func (np *NotePage) NewNotePage(retrievedNote *note.NoteData, allowEdit bool, parentWindow fyne.Window) {
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
	np.Entry = NewEntryCustom(func(cs *desktop.CustomShortcut) {
		switch cs.ShortcutName() {
		case main_ui.ScViewMode.ShortcutName():
			SetViewMode(parentWindow)
		case main_ui.ScPinNote.ShortcutName():
			PinNote(noteInfo)
		case main_ui.ScNoteColour.ShortcutName():
			ChangeNoteColour(noteInfo, parentWindow)
		case main_ui.ScShowInfo.ShortcutName():
			ShowProperties(noteInfo)
		}
	}, func() {
		np.NoteInfo.Content = np.Entry.Text
		SaveNote(&np.NoteInfo, retrievedNote, parentWindow)
	},
	)

	np.Entry.Text = np.NoteInfo.Content
	np.Entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(main_ui.AppTheme.NoteBgColour)
	noteColour, _ := main_ui.RGBStringToFyneColor(np.NoteInfo.Colour)

	NoteCanvas.noteBackground = canvas.NewRectangle(noteColour)
	if np.NoteInfo.Colour == "#e7edef" || np.NoteInfo.Colour == "#FFFFFF" || np.NoteInfo.Colour == "#ffffff" || np.NoteInfo.Colour == "#000000" {
		NoteCanvas.noteBackground = canvas.NewRectangle(main_ui.AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(NoteCanvas.noteBackground)

	NoteWidgets.markdownText = widget.NewRichTextFromMarkdown(noteInfo.Content)
	NoteWidgets.markdownText.Wrapping = fyne.TextWrapWord
	NoteWidgets.markdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, NoteWidgets.markdownText)
	NoteContainers.markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(NoteContainers.markdown)
	background := canvas.NewRectangle(AppTheme.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, NoteWidgets.entry)

	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if noteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	NoteWidgets.pinButton = widget.NewButtonWithIcon("", btnIcon, func() {
		PinNote(noteInfo)
	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	changeNotebookBtn := NewChangeNotebookButton(noteInfo, parentWindow)

	colourButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		ChangeNoteColour(noteInfo, parentWindow)
	})

	NoteWidgets.deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteNote(noteInfo, parentWindow)
	})

	tagsBtn := widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func() {
		ToggleTagsNotePanel()
	})

	propertiesButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() { ShowProperties(noteInfo) })

	NoteWidgets.deleteButton.Hide()

	NoteWidgets.modeSelect = widget.NewRadioGroup([]string{EDIT_MODE, VIEW_MODE}, func(value string) {
		switch value {
		case EDIT_MODE:
			if allowEdit {
				SetEditMode(parentWindow)
			}
		case VIEW_MODE:
			SetViewMode(parentWindow)
		}
	})

	if !allowEdit {
		NoteWidgets.modeSelect.Hide()
	}

	NoteContainers.propertiesPanel = NewProperetiesPanel()

	NoteWidgets.modeSelect.SetSelected("View")
	NoteWidgets.modeSelect.Horizontal = true
	toolbar := container.NewHBox(NoteWidgets.modeSelect, spacerLabel, NoteWidgets.pinButton, colourButton, changeNotebookBtn, tagsBtn, propertiesButton, NoteWidgets.deleteButton)

	if err = CreateNotesTagPanel(noteInfo, parentWindow); err != nil {
		dialog.ShowError(err, parentWindow)
		log.Panicln("Error creating tags panel!")
	}
	topVBox := container.NewVBox(toolbar, NoteContainers.tagsPanel)

	NoteContainers.propertiesPanel.Hide()
	NoteContainers.tagsPanel.Hide()

	return container.NewBorder(topVBox, nil, nil, NoteContainers.propertiesPanel, content)
}

func SaveNote(noteInfo *note.NoteInfo, retrievedNote *note.NoteData, parentWindow fyne.Window) {
	var noteChanges note.NoteChanges

	if noteInfo.Deleted {
		go main_ui.UpdateView()
		return
	}

	if noteInfo.NewNote {
		if noteInfo.Content != "" {
			noteChanges.ContentChanged = true
		}
	} else {
		noteChanges = note.CheckChanges(retrievedNote, noteInfo)
	}
	//if contentChanged{
	if noteChanges.ContentChanged || noteChanges.ParamsChanged {
		res, err := note.SaveNote(noteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, parentWindow)
			return
		}

		if res == 0 {
			log.Println("No note was saved (affected rows = 0)")
		} else {
			log.Println("....Note updated successfully....")
			if *retrievedNote, err = notes.GetNote(noteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, parentWindow)
			}
			main_ui.UpdateView()
		}
	} else if noteChanges.PinStatusChanged {
		// we do not want a create or modified time stamp for just pinning/unpinning notes
		res, err := note.SaveNoteNoTimeStamp(noteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, parentWindow)
			return
			//log.Panic()
		}

		if res == 0 {
			log.Println("No note was saved (affected rows = 0)")
		} else {
			log.Println("....Note updated successfully....")
			if *retrievedNote, err = notes.GetNote(noteInfo.Id); err != nil {
				log.Println("Error getting updated note")
				dialog.ShowError(err, parentWindow)
			}
			go main_ui.UpdateView()
		}
	}
}
