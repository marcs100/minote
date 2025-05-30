package minotedb

type NoteDataDB struct {
	Id               uint
	Notebook         string
	Content          string
	Created          string
	Modified         string
	Pinned           uint
	PinnedDate       string
	BackgroundColour string
}

type SearchFilter struct {
	Pinned     bool
	WholeWords bool
}
