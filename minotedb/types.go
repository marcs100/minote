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

const (
	SORT_NEWEST = iota
	SORT_OLDEST
	SORT_PINNED_NEWER
	SORT_PINNED_OLDER_
	SORT_CREATED_FIRST
	SORT_CREATED_LAST
)
