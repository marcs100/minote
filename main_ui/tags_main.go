package main_ui

import (
	"log"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/notes"
)

// For main window
func CreateTagsPanel() *fyne.Container {
	var err error = nil
	AppWidgets.tagsList, err = CreateMainTagsList()
	if err != nil {
		dialog.ShowError(err, mainWindow)
		log.Panicln("Error creating tags list")
	}
	tagLabel := widget.NewRichTextFromMarkdown("**Tags:**            ")
	listCont := container.NewStack(AppWidgets.tagsList)
	tagsPanel := container.NewBorder(tagLabel, nil, nil, nil, listCont)
	return tagsPanel
}

// For main window
func CreateMainTagsList() (*widget.List, error) {
	var err error = nil

	AppStatus.tags = notes.GetAllTags()

	tagsList := widget.NewList(
		func() int {
			return len(AppStatus.tags)
		},

		func() fyne.CanvasObject {
			return widget.NewCheck("__________________________", func(c bool) {})
		},

		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Check).Text = AppStatus.tags[id]
			o.(*widget.Check).OnChanged = func(c bool) {
				AppStatus.currentView = VIEW_TAGS
				t := o.(*widget.Check).Text
				if c {
					//fmt.Printf("tag %s is checked, add to checkedTags\n", t)
					if !slices.Contains(AppStatus.tagsChecked, t) {
						AppStatus.tagsChecked = append(AppStatus.tagsChecked, t)
					}
				} else {
					//fmt.Printf("tag unchecked, remove %s from checkedTags\n", t)
					if i := slices.Index(AppStatus.tagsChecked, t); i > -1 {
						AppStatus.tagsChecked = slices.Delete(AppStatus.tagsChecked, i, i+1)
					}
				}
				//fmt.Println(AppStatus.tagsChecked)
				UpdateView()
			}
			o.Refresh()
		},
	)
	return tagsList, err
}

// For main window
func ToggleMainTagsPanel() {
	if AppContainers.tagsPanel.Visible() {
		AppContainers.tagsPanel.Hide()
	} else {
		CreateMainTagsList()
		AppContainers.tagsPanel.Show()
	}
}
