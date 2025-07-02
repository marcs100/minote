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

	switch themeVarIn {
	case DARK_THEME:
		fmt.Println("Using Dark theme")
		appColours.MainBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourBg)
		appColours.NoteBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourNote)

		appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourCtBg)
	case LIGHT_THEME:
		fmt.Println("Using Light theme")
		appColours.MainBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourBg)
		appColours.NoteBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourNote)
		appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourCtBg)
	case SYSTEM_THEME:
		themeVariant := MainApp.Settings().ThemeVariant()
		if themeVariant == theme.VariantDark {
			log.Println("using system theme dark")
			appColours.MainBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourBg)
			appColours.NoteBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourNote)

			appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(Conf.Settings.DarkColourCtBg)

		} else if themeVariant == theme.VariantLight {
			log.Println("Using system theme light")
			appColours.MainBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourBg)
			appColours.NoteBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourNote)
			appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(Conf.Settings.LightColourCtBg)
		} else {
			log.Println("Unidentified system theme variant!!!")
			appColours.MainBgColour = MainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
			appColours.NoteBgColour = MainApp.Settings().Theme().Color(theme.ColorNameForeground, themeVariant)
		}
	}
	return appColours
}
