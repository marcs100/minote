package main_app

import (
	"errors"
	"github.com/marcs100/minote/config"
	"os"
)

func ValidateConfig(appConfig *config.Config) error {
	var err error = nil

	if _, err = os.Stat(appConfig.Settings.Database); err != nil {
		return errors.New("config error: Database path error")
	}

	if appConfig.Settings.NoteHeight < 50 || appConfig.Settings.NoteHeight > 2000 {
		return errors.New("config error: NoteHeight range error")
	}

	if appConfig.Settings.NoteWidth < 50 || appConfig.Settings.NoteWidth > 2000 {
		return errors.New("config error: NoteWidth range error")
	}

	if appConfig.Settings.RecentNotesLimit == 0 || appConfig.Settings.RecentNotesLimit > 5000 {
		return errors.New("config error: RecentNotesLimit range error")
	}

	if appConfig.Settings.GridMaxPages == 0 || appConfig.Settings.GridMaxPages > 5000 {
		return errors.New("config error: GridMaxPages range error")
	}

	if appConfig.Settings.FontSize < 5.0 || appConfig.Settings.FontSize > 150.0 {
		return errors.New("config error: FontSize invalid value")
	}

	switch appConfig.Settings.InitialLayout {
	case "grid", "page":
		break
	default:
		return errors.New("config error: unrecognised InitialLayout")
	}

	switch appConfig.Settings.InitialView {
	case "pinned", "recent":
		break
	default:
		return errors.New("config error: unrecognised default view")
	}

	switch appConfig.Settings.ThemeVariant {
	case "auto", "dark", "light":
		break
	default:
		return errors.New("config error: unrecognised theme")
	}

	if err = config.CheckRGBString(appConfig.Settings.DarkColourNote); err != nil {
		return errors.New("config error: DarkColourNote invalid rgb string")
	}

	if err = config.CheckRGBString(appConfig.Settings.LightColourNote); err != nil {
		return errors.New("config error: LightColourNote invalid rgb string")
	}

	if err = config.CheckRGBString(appConfig.Settings.DarkColourBg); err != nil {
		return errors.New("config error: DarkColourBg invalid rgb string")
	}

	if err = config.CheckRGBString(appConfig.Settings.LightColourBg); err != nil {
		return errors.New("config error: LightColourBg invalid rgb string")
	}

	if err = config.CheckRGBString(appConfig.Settings.DarkColourCtBg); err != nil {
		return errors.New("config error: DarkColourCtBg invalid rgb string")
	}

	if err = config.CheckRGBString(appConfig.Settings.LightColourCtBg); err != nil {
		return errors.New("config error: LightColourCtBg invalid rgb string")
	}

	return err
}
