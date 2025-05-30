package main_ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/config"
)

func NewSettingsWindow() {

	newConf := CopySettings()

	var themeVar theme_variant
	switch Conf.Settings.ThemeVariant {
	case "light":
		themeVar = LIGHT_THEME
	case "dark":
		themeVar = DARK_THEME
	case "system":
		themeVar = SYSTEM_THEME
	}
	AppTheme = GetThemeColours(themeVar)

	settingsWindow := mainApp.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(500, 400))

	bg := canvas.NewRectangle(AppTheme.MainBgColour)

	viewHeading := widget.NewRichTextFromMarkdown("### View")
	viewLabel := widget.NewLabel("  Default View:          ")
	viewSelect := widget.NewSelect([]string{"pinned", "recent"}, func(sel string) {
		newConf.Settings.InitialView = sel
	})
	viewSelect.SetSelected(Conf.Settings.InitialView)
	viewGrid := container.NewGridWithRows(1, viewLabel, viewSelect)

	recentNotesLimitLabel := widget.NewLabel("  Recent Note Limit:")
	recentNotesLimitEntry := widget.NewEntry()
	recentNotesLimitEntry.SetText(fmt.Sprintf("%d", Conf.Settings.RecentNotesLimit))
	recentNotesLimitEntry.OnChanged = func(input string) {
		i, err := strconv.Atoi(input)
		if err != nil {
			recentNotesLimitEntry.SetText("")
			return
		}
		if i < 1 {
			dialog.ShowInformation("Setting Error", "Recent notes limit must be > 1", settingsWindow)
		} else {
			newConf.Settings.RecentNotesLimit = i
		}
	}

	notesLimitGrid := container.NewGridWithRows(1, recentNotesLimitLabel, recentNotesLimitEntry)

	layoutHeading := widget.NewRichTextFromMarkdown("### Layout")
	layoutLabel := widget.NewLabel("  Default Layout:")
	layoutSelect := widget.NewSelect([]string{"grid", "page"}, func(sel string) {
		newConf.Settings.InitialLayout = sel
	})
	layoutSelect.Selected = Conf.Settings.InitialLayout
	layoutGrid := container.NewGridWithRows(1, layoutLabel, layoutSelect)

	gridLimitLabel := widget.NewLabel("  Notes per Page Limit:")
	gridLimitEntry := widget.NewEntry()
	gridLimitEntry.OnChanged = func(input string) {
		i, err := strconv.Atoi(input)
		if err != nil {
			gridLimitEntry.SetText("")
			return
		}

		if i < 1 {
			dialog.ShowInformation("Setting Error", "Grid pages limit must be > 1", settingsWindow)
		}
		newConf.Settings.GridMaxPages = i
	}
	gridLimitStack := container.NewStack(gridLimitEntry)
	gridLimitGrid := container.NewGridWithRows(1, gridLimitLabel, gridLimitStack)
	gridLimitEntry.SetText(fmt.Sprintf("%d", Conf.Settings.GridMaxPages))

	appearanceHeading := widget.NewRichTextFromMarkdown("### Appearance")
	appearanceLabel := widget.NewLabel("  Theme:")
	appearanceSelect := widget.NewSelect([]string{"light", "dark", "system"}, func(sel string) {
		newConf.Settings.ThemeVariant = sel
	})
	appearanceSelect.Selected = Conf.Settings.ThemeVariant
	appearanceGrid := container.NewGridWithRows(1, appearanceLabel, appearanceSelect)

	vbox := container.NewVBox(
		viewHeading,
		viewGrid,
		notesLimitGrid,
		layoutHeading,
		layoutGrid,
		gridLimitGrid,
		appearanceHeading,
		appearanceGrid)

	stack := container.NewStack(bg, vbox)

	settingsWindow.SetOnClosed(func() {
		if newConf != *Conf {
			if err := config.WriteConfig(AppStatus.configFile, newConf); err != nil {
				dialog.ShowError(err, settingsWindow)
			}
		}
	})

	settingsWindow.SetContent(stack)
	settingsWindow.Show()
}

func CopySettings() config.Config {
	return config.Config{
		Title: Conf.Title,
		Settings: config.AppSettings{
			Database:         Conf.Settings.Database,
			RecentNotesLimit: Conf.Settings.RecentNotesLimit,
			NoteWidth:        Conf.Settings.NoteWidth,
			NoteHeight:       Conf.Settings.NoteHeight,
			InitialView:      Conf.Settings.InitialView,
			InitialLayout:    Conf.Settings.InitialLayout,
			GridMaxPages:     Conf.Settings.GridMaxPages,
			ThemeVariant:     Conf.Settings.ThemeVariant,
			DarkColourNote:   Conf.Settings.DarkColourNote,
			LightColourNote:  Conf.Settings.LightColourNote,
			DarkColourBg:     Conf.Settings.DarkColourBg,
			LightColourBg:    Conf.Settings.LightColourBg,
		},
	}
}
