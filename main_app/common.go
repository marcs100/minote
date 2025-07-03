package main_app

import (
	"fyne.io/fyne/v2"
	"github.com/marcs100/minote/config"
)

const (
	VIEW_PINNED   string = "pinned"
	VIEW_RECENT   string = "recent"
	VIEW_NOTEBOOK string = "notebooks"
	VIEW_SEARCH   string = "search"
	VIEW_TAGS     string = "tags"
)

const (
	LAYOUT_GRID string = "grid"
	LAYOUT_PAGE string = "page"
)

const (
	SEARCH_FILT_PINNED     string = "Pinned"
	SEARCH_FILT_WOLE_WORDS string = "Whole words only"
)

var AppStatus ApplicationStatus //structure containing various app status
var Conf *config.Config
var MainApp fyne.App
