package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/marcs100/minote/config"
	"github.com/marcs100/minote/main_app"
)

func NewSettingsWindow() fyne.Window {

	newConf := CopySettings()

	var themeVar main_app.ThemeVariant
	switch main_app.Conf.Settings.ThemeVariant {
	case "light":
		themeVar = main_app.LIGHT_THEME
	case "dark":
		themeVar = main_app.DARK_THEME
	case "system":
		themeVar = main_app.SYSTEM_THEME
	}
	main_app.AppTheme = main_app.GetThemeColours(themeVar)

	settingsWindow := main_app.MainApp.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(500, 450))

	bg := canvas.NewRectangle(main_app.AppTheme.MainBgColour)

	viewHeading := widget.NewRichTextFromMarkdown("### View")
	viewLabel := widget.NewLabel("  Default View:          ")
	viewSelect := widget.NewSelect([]string{"pinned", "recent"}, func(sel string) {
		newConf.Settings.InitialView = sel
	})
	viewSelect.SetSelected(main_app.Conf.Settings.InitialView)
	viewGrid := container.NewGridWithRows(1, viewLabel, viewSelect)

	recentNotesLimitLabel := widget.NewLabel("  Recent Note Limit:")
	recentNotesLimitEntry := widget.NewEntry()
	recentNotesLimitEntry.SetText(fmt.Sprintf("%d", main_app.Conf.Settings.RecentNotesLimit))
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
	layoutSelect.Selected = main_app.Conf.Settings.InitialLayout
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
	gridLimitEntry.SetText(fmt.Sprintf("%d", main_app.Conf.Settings.GridMaxPages))

	appearanceHeading := widget.NewRichTextFromMarkdown("### Appearance")
	appearanceLabel := widget.NewLabel("  Theme:")
	appearanceSelect := widget.NewSelect([]string{"light", "dark", "system"}, func(sel string) {
		newConf.Settings.ThemeVariant = sel
	})
	appearanceSelect.Selected = main_app.Conf.Settings.ThemeVariant
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
		if newConf != *main_app.Conf {
			if err := config.WriteConfig(main_app.AppStatus.ConfigFile, newConf); err != nil {
				dialog.ShowError(err, settingsWindow)
			}
		}
		main_app.AppStatus.SettingsOpen = false
	})

	settingsWindow.SetContent(stack)
	return settingsWindow
}

func CopySettings() config.Config {
	return config.Config{
		Title: main_app.Conf.Title,
		Settings: config.AppSettings{
			Database:         main_app.Conf.Settings.Database,
			RecentNotesLimit: main_app.Conf.Settings.RecentNotesLimit,
			NoteWidth:        main_app.Conf.Settings.NoteWidth,
			NoteHeight:       main_app.Conf.Settings.NoteHeight,
			InitialView:      main_app.Conf.Settings.InitialView,
			InitialLayout:    main_app.Conf.Settings.InitialLayout,
			GridMaxPages:     main_app.Conf.Settings.GridMaxPages,
			ThemeVariant:     main_app.Conf.Settings.ThemeVariant,
			DarkColourNote:   main_app.Conf.Settings.DarkColourNote,
			LightColourNote:  main_app.Conf.Settings.LightColourNote,
			DarkColourBg:     main_app.Conf.Settings.DarkColourBg,
			LightColourBg:    main_app.Conf.Settings.LightColourBg,
			DarkColourCtBg:   main_app.Conf.Settings.DarkColourCtBg,
			LightColourCtBg:  main_app.Conf.Settings.LightColourCtBg,
		},
	}
}
