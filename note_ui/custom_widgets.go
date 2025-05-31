package note_ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type scribeNoteText struct {
	widget.RichText
	OnTapped func()
}

// Implement onTapped for this widget
func (sn *scribeNoteText) Tapped(*fyne.PointEvent) {
	if sn.OnTapped != nil {
		sn.OnTapped()
	}
}

func NewScribeNoteText(content string, tapped func()) *scribeNoteText {
	rt := &scribeNoteText{}
	rt.AppendMarkdown(content)
	rt.OnTapped = tapped
	return rt
}

type EntryCustom struct {
	widget.Entry
	onCustomShortCut func(cs *desktop.CustomShortcut)
	onFocusLost      func()
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

func NewEntryCustom(onCustomShortcut func(cs *desktop.CustomShortcut), onFocusLost func()) *EntryCustom {
	e := &EntryCustom{}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.onCustomShortCut = onCustomShortcut
	e.onFocusLost = onFocusLost
	e.ExtendBaseWidget(e)
	return e
}

/*
//The code below is pointless becuase widget.RichText is not focusable.
//Only widgets that are focusable can receive keybard shortcuts
type RichTextFromMarkdownCustom struct {
	widget.RichText
	onCustomShortCut func(cs *desktop.CustomShortcut)
}

func (md *RichTextFromMarkdownCustom) TypedShortcut(s fyne.Shortcut) {
	var ok bool
	var cs *desktop.CustomShortcut
	if cs, ok = s.(*desktop.CustomShortcut); !ok {
		fmt.Println("** Not a custom shortcut!!")
		return
	}
	//var name = cs.ShortcutName()
	md.onCustomShortCut(cs)
}

func NewRichTextFromMarkdownCustom(content string, onCustomShortcut func(cs *desktop.CustomShortcut)) *RichTextFromMarkdownCustom {
	md := &RichTextFromMarkdownCustom{}
	md.onCustomShortCut = onCustomShortcut
	md.AppendMarkdown(content)
	md.ExtendBaseWidget(md)
	return md
	}*/
