package ui

import (
	"log"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
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

	main_app.AppStatus.Tags = notes.GetAllTags()

	tagsList := widget.NewList(
		func() int {
			return len(main_app.AppStatus.Tags)
		},

		func() fyne.CanvasObject {
			return widget.NewCheck("__________________________", func(c bool) {})
		},

		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Check).Text = main_app.AppStatus.Tags[id]
			o.(*widget.Check).OnChanged = func(c bool) {
				main_app.AppStatus.CurrentView = main_app.VIEW_TAGS
				t := o.(*widget.Check).Text
				if c {
					//fmt.Printf("tag %s is checked, add to checkedTags\n", t)
					if !slices.Contains(main_app.AppStatus.TagsChecked, t) {
						main_app.AppStatus.TagsChecked = append(main_app.AppStatus.TagsChecked, t)
					}
				} else {
					//fmt.Printf("tag unchecked, remove %s from checkedTags\n", t)
					if i := slices.Index(main_app.AppStatus.TagsChecked, t); i > -1 {
						main_app.AppStatus.TagsChecked = slices.Delete(main_app.AppStatus.TagsChecked, i, i+1)
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
