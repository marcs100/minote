package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
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
	rt.AppendMarkdown(content)
	rt.OnTapped = tapped
	return rt
}

type EntryCustom struct {
	widget.Entry
	onCustomShortCut func(cs *desktop.CustomShortcut)
	onFocusLost      func()
	OnTypedKey       func(k *fyne.KeyEvent)
}

func (e *EntryCustom) TypedKey(k *fyne.KeyEvent) {
	if e.OnTypedKey != nil {
		if k.Name == fyne.KeyEscape {
			e.OnTypedKey(k)
		} else {
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
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.onCustomShortCut = onCustomShortcut
	e.onFocusLost = onFocusLost
	e.OnTypedKey = onEscapeKey
	e.ExtendBaseWidget(e)
	return e
}

type ButtonCustom struct {
	widget.Button
	onTapped    func()
	icon        *fyne.Resource
	tooltipText string
	// 	tooltip       *widget.PopUp
	mw *MainWindow
}

func (b *ButtonCustom) MouseIn(m *desktop.MouseEvent) {
	fmt.Println("MouseIn event!")
	if b.mw.ToolTip.Hidden {
		fmt.Println("showing tooltip!")
		b.mw.ToolTip.Show()
		b.mw.window.Canvas().Refresh(b.mw.ToolTip)
	}
}

func (b *ButtonCustom) MouseOut() {
	fmt.Println("MouseOut event!")
	if !b.mw.ToolTip.Hidden {
		fmt.Println("hiding tooltip!")
		b.mw.ToolTip.Hide()
		b.mw.window.Canvas().Refresh(b.mw.ToolTip)
	}
}

func (b *ButtonCustom) Tapped(pe *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func NewButtonCustom(label string, icon fyne.Resource, tooltipText string, mw *MainWindow, tapped func()) *ButtonCustom {
	b := &ButtonCustom{}
	b.SetIcon(icon)
	b.SetText(label)
	b.OnTapped = tapped
	b.mw = mw
	b.mw.ToolTip.Text = tooltipText
	b.ExtendBaseWidget(b)
	// text := canvas.NewText(b.tooltipText, color.White)
	// b.tooltip = widget.NewPopUp(text, b.parent) //This popup seems to be the cause of the MouseOut event always firing!
	return b
}
