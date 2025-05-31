package note_ui

import (
	"log"
	"sync"

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
	ch := make(chan bool)
	np.NotePageWidgets.Entry = NewEntryCustom(func(cs *desktop.CustomShortcut) {
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
		np.NoteInfo.Content = np.NotePageWidgets.Entry.Text
		SaveNote(np, retrievedNote, ch)
		<-ch
	},
	)

	np.NotePageWidgets.Entry.Text = np.NoteInfo.Content
	np.NotePageWidgets.Entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(main_ui.AppTheme.NoteBgColour)
	noteColour, _ := main_ui.RGBStringToFyneColor(np.NoteInfo.Colour)

	np.NotePageCanvas.NoteBackground = canvas.NewRectangle(noteColour)
	if np.NoteInfo.Colour == "#e7edef" || np.NoteInfo.Colour == "#FFFFFF" || np.NoteInfo.Colour == "#ffffff" || np.NoteInfo.Colour == "#000000" {
		np.NotePageCanvas.NoteBackground = canvas.NewRectangle(main_ui.AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(np.NotePageCanvas.NoteBackground)

	np.NotePageWidgets.MarkdownText = widget.NewRichTextFromMarkdown(np.NoteInfo.Content)
	np.NotePageWidgets.MarkdownText.Wrapping = fyne.TextWrapWord
	np.NotePageWidgets.MarkdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, np.NotePageWidgets.MarkdownText)
	np.NotePageContainers.Markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(np.NotePageContainers.Markdown)
	background := canvas.NewRectangle(main_ui.AppTheme.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, np.NotePageWidgets.Entry)

	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if np.NoteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	np.NotePageWidgets.PinButton = widget.NewButtonWithIcon("", btnIcon, func() {
		PinNote(&np.NoteInfo)
	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	changeNotebookBtn := NewChangeNotebookButton(&np.NoteInfo, parentWindow)

	colourButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		ChangeNoteColour(&np.NoteInfo, parentWindow)
	})

	np.NotePageWidgets.DeleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteNote(&np.NoteInfo, parentWindow)
	})

	tagsBtn := widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func() {
		ToggleTagsNotePanel()
	})

	propertiesButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() { ShowProperties(&np.NoteInfo) })

	np.NotePageWidgets.DeleteButton.Hide()

	np.NotePageWidgets.ModeSelect = widget.NewRadioGroup([]string{main_ui.EDIT_MODE, main_ui.VIEW_MODE}, func(value string) {
		switch value {
		case main_ui.EDIT_MODE:
			if allowEdit {
				SetEditMode(parentWindow)
			}
		case main_ui.VIEW_MODE:
			SetViewMode(parentWindow)
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

//Pass a pointer to the note page - *NotePage
// Make sure this functiion is thread safe, as multiple note instances can be calling this function

func SaveNote(np *NotePage, retrievedNote *note.NoteData, ch chan bool) {
	var mut sync.Mutex
	mut.Lock()
	defer mut.Unlock()

	var noteChanges note.NoteChanges

	if np.NoteInfo.Deleted {
		main_ui.UpdateView()
		ch <- true
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
			ch <- true
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
			main_ui.UpdateView()
		}
	} else if noteChanges.PinStatusChanged {
		// we do not want a create or modified time stamp for just pinning/unpinning notes
		res, err := note.SaveNoteNoTimeStamp(&np.NoteInfo)
		if err != nil {
			log.Println("Error saving note")
			dialog.ShowError(err, np.ParentWindow)
			ch <- true
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
			main_ui.UpdateView()
		}
	}
	ch <- true
}
