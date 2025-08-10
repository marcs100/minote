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
	FontSize     float32
	BgColour     color.Color
	FgColour     color.Color
	EntryColour  color.Color
	ButtonColour color.Color
	AccentColour color.Color
}

//var _ fyne.Theme = (*minoteTheme)(nil)

func (t *minoteTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	/*case theme.ColorNameButton:
	col := conversions.RGBStringToFyneColor("#332212")
	return col*/
	case theme.ColorNameBackground:
		return t.BgColour
	case theme.ColorNameInputBackground:
		return t.EntryColour
	case theme.ColorNameHover:
		return t.AccentColour
	case theme.ColorNameButton:
		return t.ButtonColour
	case theme.ColorNameForeground:
		return t.FgColour
	case theme.ColorNameMenuBackground:
		return t.BgColour
	case theme.ColorNameFocus:
		return t.AccentColour
	case theme.ColorNameOverlayBackground:
		return t.BgColour
	case theme.ColorNamePrimary:
		return t.AccentColour
	default:
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
