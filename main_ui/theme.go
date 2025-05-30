package main_ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/theme"
)

// Background colour for notes based on current theme variane light/dark
func GetThemeColours(themeVarIn theme_variant) AppColours {
	var appColours AppColours
	var err error

	switch themeVarIn {
	case DARK_THEME:
		fmt.Println("Using Dark theme")
		appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourBg)
		if err != nil {
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourNote)
		if err != nil {
			log.Panicln(err)
		}
	case LIGHT_THEME:
		fmt.Println("Using Light theme")
		appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourBg)
		if err != nil {
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourNote)
		if err != nil {
			log.Panicln(err)
		}
	case SYSTEM_THEME:
		themeVariant := mainApp.Settings().ThemeVariant()
		if themeVariant == theme.VariantDark {
			log.Println("using system theme dark")
			appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourBg)
			if err != nil {
				log.Panicln(err)
			}
			appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourNote)
			if err != nil {
				log.Panicln(err)
			}

		} else if themeVariant == theme.VariantLight {
			log.Println("Using system theme light")
			appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourBg)
			if err != nil {
				log.Panicln(err)
			}
			appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourNote)
			if err != nil {
				log.Panicln(err)
			}
		} else {
			log.Println("Unidentified system theme variant!!!")
			appColours.MainBgColour = mainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
			appColours.NoteBgColour = mainApp.Settings().Theme().Color(theme.ColorNameForeground, themeVariant)
		}
	}
	return appColours
}
