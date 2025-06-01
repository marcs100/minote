package main_app

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/theme"
	"github.com/marcs100/minote/conversions"
)

// Background colour for notes based on current theme variane light/dark
func GetThemeColours(themeVarIn ThemeVariant) AppColours {
	var appColours AppColours
	var err error

	switch themeVarIn {
	case DARK_THEME:
		fmt.Println("Using Dark theme")
		appColours.MainBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourBg)
		if err != nil {
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourNote)
		if err != nil {
			log.Panicln(err)
		}
	case LIGHT_THEME:
		fmt.Println("Using Light theme")
		appColours.MainBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.LightColourBg)
		if err != nil {
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.LightColourNote)
		if err != nil {
			log.Panicln(err)
		}
	case SYSTEM_THEME:
		themeVariant := MainApp.Settings().ThemeVariant()
		if themeVariant == theme.VariantDark {
			log.Println("using system theme dark")
			appColours.MainBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourBg)
			if err != nil {
				log.Panicln(err)
			}
			appColours.NoteBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourNote)
			if err != nil {
				log.Panicln(err)
			}

		} else if themeVariant == theme.VariantLight {
			log.Println("Using system theme light")
			appColours.MainBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.LightColourBg)
			if err != nil {
				log.Panicln(err)
			}
			appColours.NoteBgColour, err = conversions.RGBStringToFyneColor(Conf.Settings.LightColourNote)
			if err != nil {
				log.Panicln(err)
			}
		} else {
			log.Println("Unidentified system theme variant!!!")
			appColours.MainBgColour = MainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
			appColours.NoteBgColour = MainApp.Settings().Theme().Color(theme.ColorNameForeground, themeVariant)
		}
	}
	return appColours
}
