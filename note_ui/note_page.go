package note_ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
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
