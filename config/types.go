package config

type Config struct {
	Title    string
	Settings AppSettings `toml:"settings"`
}

type AppSettings struct {
	Database          string
	RecentNotesLimit  int     `toml:"recent_notes_limit"`
	NoteWidth         float32 `toml:"note_width"`
	NoteHeight        float32 `toml:"note_height"`
	InitialView       string  `toml:"initial_view"`
	InitialLayout     string  `toml:"initial_layout"`
	GridMaxPages      int     `toml:"grid_max_pages"`
	FontSize          float32 `toml:"font_size"`
	DateFormat        string  `toml:"date_format"`
	TimeFormat        string  `toml:"time_format"`
	DateTimeFormat    string  `toml:"date_time_format"`
	F4Snippet         string  `toml:"f4_snippet"`
	F5Snippet         string  `toml:"f5_snippet"`
	F6Snippet         string  `toml:"f6_snippet"`
	ThemeVariant      string  `toml:"theme_variant"`
	DarkColourNote    string  `toml:"dark_colour_note"`
	LightColourNote   string  `toml:"light_colour_note"`
	DarkColourBg      string  `toml:"dark_colour_bg"`
	LightColourBg     string  `toml:"light_colour_bg"`
	DarkColourFg      string  `toml:"dark_colour_fg"`
	LightColourFg     string  `toml:"light_colour_fg"`
	DarkColourCtBg    string  `toml:"dark_colour_ct_bg"`
	LightColourCtBg   string  `toml:"light_colour_ct_bg"`
	DarkColourAccent  string  `toml:"dark_colour_accent"`
	LightColourAccent string  `toml:"light_colour_accent"`
	DarkColourButton  string  `toml:"dark_colour_button"`
	LightColourButton string  `toml:"light_colour_button"`
}
