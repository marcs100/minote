package conversions

import (
	"image/color"
	//"fyne.io/fyne/v2"
)

func FyneColourToRGBHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return RGBToHexString(r, g, b)
}

func RGBStringToFyneColor(colourStr string) color.RGBA {
	var fyneColour color.RGBA
	r, g, b, err := StringToRGBValues(colourStr)

	if err != nil {
		fyneColour = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	fyneColour = color.RGBA{R: r, G: g, B: b, A: 255}

	return fyneColour
}
