package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//This is the interface theme must use
/*type Theme interface {
Color(ThemeColorName, ThemeVariant) color.Color
Font(TextStyle) Resource
Icon(ThemeIconName) Resource
Size(ThemeSizeName) float32
}*/

type minoteTheme struct {
	FontSize    float32
	BgColour    color.Color
	EntryColour color.Color
}

//var _ fyne.Theme = (*minoteTheme)(nil)

func (t *minoteTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantDark {
		switch name {
		/*case theme.ColorNameButton:
		col := conversions.RGBStringToFyneColor("#332212")
		return col*/
		case theme.ColorNameBackground:
			return t.BgColour
		case theme.ColorNameInputBackground:
			return t.EntryColour

		default:
			return theme.DefaultTheme().Color(name, variant)
		}

	} else {
		return theme.DefaultTheme().Color(name, variant)

	}
}

func (t *minoteTheme) Icon(name fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(name)
}

func (t *minoteTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *minoteTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return t.FontSize
	}
	return theme.DefaultTheme().Size(name)

}
