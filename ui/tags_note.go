package ui

import (
	"errors"
	"fmt"
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/note"
)

func CreateNotesTagPanel(np *NotePage) error {
	np.NotePageContainers.TagLabels = container.NewHBox()
	np.NotePageContainers.TagsPanel = container.NewHScroll(np.NotePageContainers.TagLabels)
	np.NotePageWidgets.AddTagButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		tagEntry := widget.NewEntry()
		tagEntryItem := widget.NewFormItem("tag entry", tagEntry)
		tagEntryDialog := dialog.NewForm("      Enter tag name      ", "OK", "Cancel", []*widget.FormItem{tagEntryItem}, func(confirmed bool) {
			if confirmed {
				err := note.WriteTag(tagEntry.Text, np.NoteInfo.Id)
				if err != nil {
					dialog.ShowError(err, np.ParentWindow)
				}
				np.UpdateTags()
				np.UpdateProperties()
				np.RefreshWindow()
			}
		}, np.ParentWindow)
		tagEntryDialog.Show()
	})

	err := np.UpdateTags()
	return err
}

func (np *NotePage) UpdateTags() error {
	var tags []string
	var err error

	if np.NotePageContainers.TagsPanel == nil {
		return errors.New("tags panel is nil")
	}

	np.NotePageContainers.TagLabels.RemoveAll()
	np.NotePageContainers.TagLabels.Add(widget.NewLabel("Tags:  "))

	if np.NoteInfo.NewNote {
		np.NotePageContainers.TagLabels.Add(np.NotePageWidgets.AddTagButton)
		return nil
	}

	if tags, err = note.GetTagsForNote(np.NoteInfo.Id); err == nil {
		//fmt.Println(tags)
		for _, tag := range tags {
			//NoteContainers.tagLabels.Add(container.NewStack(widget.NewLabel(tag)))
			np.NotePageContainers.TagLabels.Add(widget.NewButton(tag, func() {
				tagDialog := dialog.NewConfirm("Delete tag", fmt.Sprint("Do you want to delete tag ", tag, "?"), func(confirmed bool) {
					if confirmed {
						err := note.DeleteTag(tag, np.NoteInfo.Id)
						if err != nil {
							log.Println(err)
						}
						np.UpdateTags()
					}
				}, np.ParentWindow)
				tagDialog.Show()
			}))
		}
		np.NotePageContainers.TagLabels.Refresh()
		np.UpdateProperties()
		np.RefreshWindow()
		if main_app.AppStatus.CurrentView == main_app.VIEW_TAGS {
			np.MainAppWindow.UpdateMainTagsList()
		}
	}

	np.NotePageContainers.TagLabels.Add(np.NotePageWidgets.AddTagButton)

	return err
}

// For notes container
func (np *NotePage) ToggleTagsNotePanel() {
	if np.NotePageContainers.TagsPanel.Visible() {
		np.NotePageContainers.TagsPanel.Hide()
	} else {
		np.NotePageContainers.TagsPanel.Show()
	}
}

func (np *NotePage) TagsButtonDisplay() {
	if np.NoteInfo.NewNote {
		np.NotePageWidgets.TagsButton.Hide()
	} else {
		if np.NotePageWidgets.TagsButton.Hidden {
			np.NotePageWidgets.TagsButton.Show()
		}
	}
}
