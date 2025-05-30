package main_ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// View pinned notes
var scViewPinned = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl,
}

// View recent Notes
var scViewRecent = &desktop.CustomShortcut{
	KeyName:  fyne.KeyR,
	Modifier: fyne.KeyModifierControl,
}

// View Notebooks (toggles notebook list panel)
var scShowNotebooks = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl,
}

// open search panel
var scFind = &desktop.CustomShortcut{
	KeyName:  fyne.KeyF,
	Modifier: fyne.KeyModifierControl,
}

// Open a new note
var scOpenNote = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

var scPageForward = &desktop.CustomShortcut{
	KeyName:  fyne.KeyPeriod,
	Modifier: fyne.KeyModifierControl,
}

var scPageBack = &desktop.CustomShortcut{
	KeyName:  fyne.KeyComma,
	Modifier: fyne.KeyModifierControl,
}

// Set edit mode
var scEditMode = &desktop.CustomShortcut{
	KeyName:  fyne.KeyE,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Set view mode
var scViewMode = &desktop.CustomShortcut{
	KeyName:  fyne.KeyQ,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Pin note
var scPinNote = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Change note colour
var scNoteColour = &desktop.CustomShortcut{
	KeyName:  fyne.KeyC,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Open notebooks menu *** NOT IN USE ******************
var scChangeNoteNotebook = &desktop.CustomShortcut{
	KeyName:  fyne.KeyB,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// show information (properties)
var scShowInfo = &desktop.CustomShortcut{
	KeyName:  fyne.KeyI,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}
