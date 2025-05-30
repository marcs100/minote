package main_ui

import (
	"image/color"

	"github.com/marcs100/minote/conversions"
	//"fyne.io/fyne/v2"
)

func FyneColourToRGBHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return conversions.RGBToHexString(r, g, b)
}

func RGBStringToFyneColor(colourStr string) (color.RGBA, error) {
	var fyneColour color.RGBA
	r, g, b, err := conversions.StringToRGBValues(colourStr)

	if err == nil {
		fyneColour = color.RGBA{R: r, G: g, B: b, A: 255}
	}

	return fyneColour, err
}
