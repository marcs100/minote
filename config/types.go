package config

type Config struct {
	Title    string
	Settings AppSettings `toml:"settings"`
}

type AppSettings struct {
	Database         string
	RecentNotesLimit int     `toml:"recent_notes_limit"`
	NoteWidth        float32 `toml:"note_width"`
	NoteHeight       float32 `toml:"note_height"`
	InitialView      string  `toml:"initial_view"`
	InitialLayout    string  `toml:"initial_layout"`
	GridMaxPages     int     `toml:"grid_max_pages"`
	ThemeVariant     string  `toml:"theme_variant"`
	DarkColourNote   string  `toml:"dark_colour_note"`
	LightColourNote  string  `toml:"light_colour_note"`
	DarkColourBg     string  `toml:"dark_colour_bg"`
	LightColourBg    string  `toml:"light_colour_bg"`
}
