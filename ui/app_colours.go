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
		appColours.MainFgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourFg)
		appColours.ButtonColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourButton)
		appColours.AccentColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.DarkColourAccent)

	case LIGHT_THEME:
		fmt.Println("Using Light theme")
		appColours.MainBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourBg)
		appColours.NoteBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourNote)
		appColours.MainCtrlsBgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourCtBg)
		appColours.MainFgColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourFg)
		appColours.ButtonColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourButton)
		appColours.AccentColour = conversions.RGBStringToFyneColor(main_app.Conf.Settings.LightColourAccent)
	}

	// fmt.Println(appColours.MainBgColour)
	// fmt.Println(appColours.NoteBgColour)
	// fmt.Println(appColours.MainCtrlsBgColour)

	return appColours
}
