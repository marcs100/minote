package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MarkdownCustom struct {
	widget.RichText
	OnTapped func()
}

// Implement onTapped for this widget
func (sn *MarkdownCustom) Tapped(*fyne.PointEvent) {
	if sn.OnTapped != nil {
		sn.OnTapped()
	}
}

func NewMarkdownCustom(content string, tapped func()) *MarkdownCustom {
	rt := &MarkdownCustom{}
	rt.ExtendBaseWidget(rt)
	rt.AppendMarkdown(content)
	rt.OnTapped = tapped
	return rt
}

//##################################################################

type EntryCustom struct {
	widget.Entry
	onCustomShortCut func(cs *desktop.CustomShortcut)
	onFocusLost      func()
	OnTypedKey       func(k *fyne.KeyEvent)
}

func (e *EntryCustom) TypedKey(k *fyne.KeyEvent) {
	if e.OnTypedKey != nil {
		switch k.Name {
		case fyne.KeyEscape, fyne.KeyF1, fyne.KeyF2, fyne.KeyF3, fyne.KeyF4, fyne.KeyF5:
			e.OnTypedKey(k)
		default:
			e.Entry.TypedKey(k)
		}
	}
}

func (e *EntryCustom) TypedShortcut(s fyne.Shortcut) {
	var ok bool
	var cs *desktop.CustomShortcut
	if cs, ok = s.(*desktop.CustomShortcut); !ok {
		//fmt.Printf("shortcut name is %s", cs.ShortcutName())
		e.Entry.TypedShortcut(s) //not a customshort cut - pass through to normal predifined shortcuts
		fmt.Println("** Not a custom shortcut!!")
		return
	}

	e.onCustomShortCut(cs)
}

func (e *EntryCustom) FocusLost() {
	e.Entry.FocusLost()
	e.onFocusLost()
}

func NewEntryCustom(onCustomShortcut func(cs *desktop.CustomShortcut), onEscapeKey func(k *fyne.KeyEvent), onFocusLost func()) *EntryCustom {
	e := &EntryCustom{}
	e.ExtendBaseWidget(e)
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.onCustomShortCut = onCustomShortcut
	e.onFocusLost = onFocusLost
	e.OnTypedKey = onEscapeKey
	return e
}

//################################################################

type ButtonWithTooltip struct {
	widget.Button
	icon        *fyne.Resource
	window      fyne.Window
	tooltipText string
	tooltip     *canvas.Text
}

func (b *ButtonWithTooltip) MouseIn(m *desktop.MouseEvent) {
	if len(strings.TrimSpace(b.tooltip.Text)) == 0 {
		b.tooltip.Text = b.tooltipText
		b.window.Canvas().Refresh(b.tooltip)
	}
	b.Button.MouseIn(m)
}

func (b *ButtonWithTooltip) MouseOut() {
	if len(strings.TrimSpace(b.tooltip.Text)) > 0 {
		b.tooltip.Text = fmt.Sprintf("%-25s", "")
		b.window.Canvas().Refresh(b.tooltip)
	}
	b.Button.MouseOut()
}

func (b *ButtonWithTooltip) Tapped(pe *fyne.PointEvent) {
	if b.OnTapped != nil {
		animateButton(&b.Button)
		b.OnTapped()

	}
}

func NewButtonWithTooltip(label string, icon fyne.Resource, tooltipText string, tooltip *canvas.Text, window fyne.Window, tapped func()) *ButtonWithTooltip {
	b := &ButtonWithTooltip{}
	b.ExtendBaseWidget(b)
	b.SetIcon(icon)
	b.SetText(label)
	b.OnTapped = tapped
	b.tooltip = tooltip
	b.tooltipText = tooltipText
	b.window = window
	// b.tooltip = widget.NewPopUp(text, b.parent) //This popup seems to be the cause of the MouseOut event always firing!
	return b
}

//##################################################################

// type ToolbarActionWithTooltip struct {
// 	widget.ToolbarAction
// 	onTapped    func()
// 	icon        *fyne.Resource
// 	window      fyne.Window
// 	tooltipText string
// 	tooltip     *widget.Label
// }

// func (ta *ToolbarActionWithTooltip) MouseIn(m *desktop.MouseEvent) {
// 	if len(strings.TrimSpace(ta.tooltip.Text)) == 0 {
// 		ta.tooltip.SetText(ta.tooltipText)
// 		ta.window.Canvas().Refresh(ta.tooltip)
// 	}
// 	ta.ToolbarAction.MouseIn(m)
// }

// func (ta *ToolbarActionWithTooltip) MouseOut() {
// 	if len(strings.TrimSpace(ta.tooltip.Text)) > 0 {
// 		ta.tooltip.SetText("                         ")
// 		ta.window.Canvas().Refresh(ta.tooltip)
// 	}
// 	ta.ToolbarAction.MouseOut()
// }

// func (ta *ToolbarActionWithTooltip) Tapped(pe *fyne.PointEvent) {
// 	if ta.OnTapped != nil {
// 		ta.OnTapped()
// 	}
// 	ta.ToolbarAction.Tapped(pe)
// }
// func NewToolbarActionWithTooltip(label string, icon fyne.Resource, tooltipText string, tooltip *widget.Label, window fyne.Window, tapped func()) *ToolbarActionWithTooltip {
// 	ta := &ToolbarActionWithTooltip{}
// 	ta.ExtendBaseWidget(ta)
// 	ta.SetIcon(icon)
// 	ta.SetText(label)
// 	ta.OnTapped = tapped
// 	ta.tooltip = tooltip
// 	ta.tooltipText = tooltipText
// 	ta.window = window
// 	// b.tooltip = widget.NewPopUp(text, b.parent) //This popup seems to be the cause of the MouseOut event always firing!
// 	return ta
// }
