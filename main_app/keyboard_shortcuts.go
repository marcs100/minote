package main_app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// View pinned notes
var ScViewPinned = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl,
}

// View recent Notes
var ScViewRecent = &desktop.CustomShortcut{
	KeyName:  fyne.KeyR,
	Modifier: fyne.KeyModifierControl,
}

// View Notebooks (toggles notebook list panel)
var ScShowNotebooks = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl,
}

// View Tags (toggles taggles list panel)
var ScShowTags = &desktop.CustomShortcut{
	KeyName:  fyne.KeyT,
	Modifier: fyne.KeyModifierControl,
}

// open search panel
var ScFind = &desktop.CustomShortcut{
	KeyName:  fyne.KeyF,
	Modifier: fyne.KeyModifierControl,
}

// Open a new note
var ScOpenNote = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

var ScPageForward = &desktop.CustomShortcut{
	KeyName:  fyne.KeyPeriod,
	Modifier: fyne.KeyModifierControl,
}

var ScPageBack = &desktop.CustomShortcut{
	KeyName:  fyne.KeyComma,
	Modifier: fyne.KeyModifierControl,
}

// Set edit mode
// var ScEditMode = &desktop.CustomShortcut{
// 	KeyName:  fyne.KeyE,
// 	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
// }

// Set view mode
// var ScViewMode = &desktop.CustomShortcut{
// 	KeyName:  fyne.KeyQ,
// 	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
// }

// Pin note
var ScPinNote = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Show note tags
var ScNoteTags = &desktop.CustomShortcut{
	KeyName:  fyne.KeyT,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Change note colour
var ScNoteColour = &desktop.CustomShortcut{
	KeyName:  fyne.KeyC,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Open notebooks menu *** NOT IN USE ******************
var ScChangeNoteNotebook = &desktop.CustomShortcut{
	KeyName:  fyne.KeyB,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// show information (properties)
var ScShowInfo = &desktop.CustomShortcut{
	KeyName:  fyne.KeyI,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}
