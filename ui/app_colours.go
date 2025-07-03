package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/theme"
	"github.com/marcs100/minote/conversions"
	"github.com/marcs100/minote/main_app"
)

// Background colour for notes based on current theme variane light/dark
func GetAppColours(themeVarIn ThemeVariant) AppColours {
	var appColours AppColours

	switch themeVarIn {
	case DARK_THEME:
		fmt.Println("Using Dark theme")
		appColours.MainBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourBg)
		appColours.NoteBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourNote)
		appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourCtBg)
	case LIGHT_THEME:
		fmt.Println("Using Light theme")
		appColours.MainBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourBg)
		appColours.NoteBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourNote)
		appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourCtBg)
	case SYSTEM_THEME:
		themeVariant := main_app.MainApp.Settings().ThemeVariant()
		if themeVariant == theme.VariantDark {
			log.Println("using system theme dark")
			appColours.MainBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourBg)
			appColours.NoteBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourNote)
			appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourCtBg)

		} else if themeVariant == theme.VariantLight {
			log.Println("Using system theme light")
			appColours.MainBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourBg)
			appColours.NoteBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourNote)
			appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourCtBg)
		} else {
			log.Panicln("Unidentified system theme variant!!!")
			// appColours.MainBgColour = main_app.MainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
			// appColours.NoteBgColour = main_app.MainApp.Settings().Theme().Color(theme.ColorNameForeground, themeVariant)
		}
	}

	fmt.Println(appColours.MainBgColour)
	fmt.Println(appColours.NoteBgColour)
	fmt.Println(appColours.MainCtrlsBgColour)

	return appColours
}
