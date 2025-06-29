package ui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/main_app"
	"github.com/marcs100/minote/notes"
)

// For main window
func (mw *MainWindow) CreateTagsPanel() *fyne.Container {
	tagsList := mw.CreateMainTagsList(notes.GetAllTags())
	tagLabel := widget.NewRichTextFromMarkdown("**Tags:**            ")
	searchEntry := widget.NewEntry()
	searchEntry.OnChanged = func(t string) {
		if len(t) > 1 {
			mw.UpdateTagsFromSearch(t)
		} else {
			mw.UpdateMainTagsList()
		}
	}
	vbox := container.NewVBox(tagLabel, searchEntry)
	mw.AppContainers.tagsList = container.NewStack(tagsList)
	tagsPanel := container.NewBorder(vbox, nil, nil, nil, mw.AppContainers.tagsList)
	return tagsPanel
}

// For main window
func (mw *MainWindow) CreateMainTagsList(tags []string) *widget.List {
	tagsList := widget.NewList(
		func() int {
			return len(tags)
		},

		func() fyne.CanvasObject {
			return widget.NewCheck("__________________________", func(c bool) {})
		},

		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Check).Text = tags[id]
			if slices.Contains(main_app.AppStatus.TagsChecked, tags[id]) {
				o.(*widget.Check).Checked = true
			} else {
				o.(*widget.Check).Checked = false
			}
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
				mw.UpdateView()
			}
			o.Refresh()
		},
	)

	tagsList.Refresh()

	return tagsList
}

func (mw *MainWindow) UpdateMainTagsList() {
	if mw.AppContainers.tagsList != nil && mw.AppContainers.tagsPanel != nil {
		tagsList := mw.CreateMainTagsList(notes.GetAllTags())
		mw.AppContainers.tagsList.RemoveAll()
		mw.AppContainers.tagsList.Add(tagsList)
		mw.AppContainers.tagsPanel.Refresh()
	}
}

func (mw *MainWindow) UpdateTagsFromSearch(search string) {
	tagsList := mw.CreateMainTagsList(notes.GetTagsWithSearch(search))
	mw.AppContainers.tagsList.RemoveAll()
	mw.AppContainers.tagsList.Add(tagsList)
	mw.AppContainers.tagsPanel.Refresh()

}

// For main window
func (mw *MainWindow) ToggleMainTagsPanel() {
	if mw.AppContainers.tagsPanel.Visible() {
		mw.AppContainers.tagsPanel.Hide()
	} else {
		mw.UpdateMainTagsList()
		mw.AppContainers.tagsPanel.Show()
	}
}
