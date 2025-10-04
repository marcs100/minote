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

func ShowSettings(parentWindow fyne.Window, UI_Colours AppColours) {

	newConf := CopySettings()

	bg := canvas.NewRectangle(UI_Colours.MainBgColour)

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
			dialog.ShowInformation("Settings Error", "Recent notes limit must be greater than 0", parentWindow)
			recentNotesLimitEntry.SetText("")
			return
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
			dialog.ShowInformation("Settings Error", "Grid pages limit must be > 1", parentWindow)
			gridLimitEntry.SetText("")
			return
		}
		newConf.Settings.GridMaxPages = i
	}
	gridLimitStack := container.NewStack(gridLimitEntry)
	gridLimitGrid := container.NewGridWithRows(1, gridLimitLabel, gridLimitStack)
	gridLimitEntry.SetText(fmt.Sprintf("%d", main_app.Conf.Settings.GridMaxPages))

	appearanceHeading := widget.NewRichTextFromMarkdown("### Appearance")
	appearanceLabel := widget.NewLabel("  Theme:")
	appearanceSelect := widget.NewSelect([]string{"light", "dark", "auto"}, func(sel string) {
		newConf.Settings.ThemeVariant = sel
	})
	appearanceSelect.Selected = main_app.Conf.Settings.ThemeVariant
	appearanceGrid := container.NewGridWithRows(1, appearanceLabel, appearanceSelect)

	fontSizeLabel := widget.NewLabel("  FontSize:")
	fontSizeEntry := widget.NewEntry()
	fontSizeEntry.SetText(fmt.Sprintf("%.1f", main_app.Conf.Settings.FontSize))
	fontSizeEntry.OnChanged = func(input string) {
		if len(input) > 0 {
			f64, err := strconv.ParseFloat(input, 32)
			if err != nil {
				fontSizeEntry.SetText("")
				return
			}
			i := float32(f64)
			if i < 5 || i > 80 {
				dialog.ShowInformation("Settings Error", "Font size must be between 5 and 80", parentWindow)
				fontSizeEntry.SetText("")
				return
			}
			fontSizeEntry.SetText(fmt.Sprintf("%.1f", i))
			newConf.Settings.FontSize = i
		}
	}
	fontGrid := container.NewGridWithRows(1, fontSizeLabel, fontSizeEntry)

	snippetsHeading := widget.NewRichTextFromMarkdown("### Text Snippets")
	F4TextSnipLabel := widget.NewLabel("  Snippet (F4):")
	f4TextSnipEntry := widget.NewEntry()
	f4TextSnipEntry.SetText(newConf.Settings.F4Snippet)
	f4TextSnipEntry.OnChanged = func(snippet string) {
		newConf.Settings.F4Snippet = snippet
	}
	f4ScrollSnip := container.NewHScroll(f4TextSnipEntry)
	f4SnipGrid := container.NewGridWithRows(1, F4TextSnipLabel, f4ScrollSnip)
	F5TextSnipLabel := widget.NewLabel("  Snippet (F5):")
	f5TextSnipEntry := widget.NewEntry()
	f5TextSnipEntry.SetText(newConf.Settings.F5Snippet)
	f5TextSnipEntry.OnChanged = func(snippet string) {
		newConf.Settings.F5Snippet = snippet
	}
	f5ScrollSnip := container.NewHScroll(f5TextSnipEntry)
	f5SnipGrid := container.NewGridWithRows(1, F5TextSnipLabel, f5ScrollSnip)
	F6TextSnipLabel := widget.NewLabel("  Snippet (F6):")
	f6TextSnipEntry := widget.NewEntry()
	f6TextSnipEntry.SetText(newConf.Settings.F6Snippet)
	f6ScrollSnip := container.NewHScroll(f6TextSnipEntry)
	f6TextSnipEntry.OnChanged = func(snippet string) {
		newConf.Settings.F6Snippet = snippet
	}
	f6SnipGrid := container.NewGridWithRows(1, F6TextSnipLabel, f6ScrollSnip)

	vbox := container.NewVBox(
		viewHeading,
		viewGrid,
		notesLimitGrid,
		layoutHeading,
		layoutGrid,
		gridLimitGrid,
		appearanceHeading,
		appearanceGrid,
		fontGrid,
		snippetsHeading,
		f4SnipGrid,
		f5SnipGrid,
		f6SnipGrid,
		widget.NewLabel(" "),
	)

	stack := container.NewStack(bg, vbox)

	formItem := widget.NewFormItem("", stack)
	d := dialog.NewForm("      Settings      ", "Save", "Cancel", []*widget.FormItem{formItem}, func(confirmed bool) {
		if confirmed {
			if recentNotesLimitEntry.Text == "" || gridLimitEntry.Text == "" || fontSizeEntry.Text == "" {
				dialog.ShowInformation("Settings Error", "blank entries found, settings will NOT be saved!", parentWindow)
				return
			}

			if newConf != *main_app.Conf {
				var err error = nil
				if err = main_app.ValidateConfig(&newConf); err != nil {
					dialog.ShowInformation("Settings error", fmt.Sprint(err, " settings will NOT be saved!"), parentWindow)
					return
				}

				if err = config.WriteConfig(main_app.AppStatus.ConfigFile, newConf); err != nil {
					dialog.ShowError(err, parentWindow)
				}

				main_app.Conf = &newConf
			}
		}
	}, parentWindow)

	d.Show()

}

func CopySettings() config.Config {
	return config.Config{
		Title: main_app.Conf.Title,
		Settings: config.AppSettings{
			Database:          main_app.Conf.Settings.Database,
			RecentNotesLimit:  main_app.Conf.Settings.RecentNotesLimit,
			NoteWidth:         main_app.Conf.Settings.NoteWidth,
			NoteHeight:        main_app.Conf.Settings.NoteHeight,
			InitialView:       main_app.Conf.Settings.InitialView,
			InitialLayout:     main_app.Conf.Settings.InitialLayout,
			GridMaxPages:      main_app.Conf.Settings.GridMaxPages,
			DateFormat:        main_app.Conf.Settings.DateFormat,
			TimeFormat:        main_app.Conf.Settings.TimeFormat,
			DateTimeFormat:    main_app.Conf.Settings.DateTimeFormat,
			F4Snippet:         main_app.Conf.Settings.F4Snippet,
			F5Snippet:         main_app.Conf.Settings.F5Snippet,
			F6Snippet:         main_app.Conf.Settings.F6Snippet,
			ThemeVariant:      main_app.Conf.Settings.ThemeVariant,
			FontSize:          main_app.Conf.Settings.FontSize,
			DarkColourNote:    main_app.Conf.Settings.DarkColourNote,
			LightColourNote:   main_app.Conf.Settings.LightColourNote,
			DarkColourBg:      main_app.Conf.Settings.DarkColourBg,
			LightColourBg:     main_app.Conf.Settings.LightColourBg,
			DarkColourCtBg:    main_app.Conf.Settings.DarkColourCtBg,
			LightColourCtBg:   main_app.Conf.Settings.LightColourCtBg,
			DarkColourButton:  main_app.Conf.Settings.DarkColourButton,
			LightColourButton: main_app.Conf.Settings.LightColourButton,
			DarkColourFg:      main_app.Conf.Settings.DarkColourFg,
			LightColourFg:     main_app.Conf.Settings.LightColourFg,
			DarkColourAccent:  main_app.Conf.Settings.DarkColourAccent,
			LightColourAccent: main_app.Conf.Settings.LightColourAccent,
		},
	}
}
