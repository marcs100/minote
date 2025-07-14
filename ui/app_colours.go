package ui

import (
	"fmt"

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
	}

	// fmt.Println(appColours.MainBgColour)
	// fmt.Println(appColours.NoteBgColour)
	// fmt.Println(appColours.MainCtrlsBgColour)

	return appColours
}
