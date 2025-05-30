package note

type NoteData struct {
	Id               uint
	Notebook         string
	Content          string
	Created          string
	Modified         string
	Pinned           uint
	PinnedDate       string
	BackgroundColour string
}

type NoteInfo struct {
	Id           uint
	NewNote      bool
	Notebook     string
	Pinned       bool
	PinnedDate   string
	Colour       string
	DateCreated  string
	DateModified string
	Content      string
	Hash         string
	Deleted      bool
}

type NoteChanges struct {
	ContentChanged   bool
	ParamsChanged    bool
	PinStatusChanged bool
}
