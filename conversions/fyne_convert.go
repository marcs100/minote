package conversions

import (
	"image/color"
	//"fyne.io/fyne/v2"
)

func FyneColourToRGBHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return RGBToHexString(r, g, b)
}

func RGBStringToFyneColor(colourStr string) (color.RGBA, error) {
	var fyneColour color.RGBA
	r, g, b, err := StringToRGBValues(colourStr)

	if err == nil {
		fyneColour = color.RGBA{R: r, G: g, B: b, A: 255}
	}

	return fyneColour, err
}
