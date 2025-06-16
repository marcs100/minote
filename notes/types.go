package notes

type SearchFilter struct {
	Pinned     bool
	WholeWords bool
}

const (
	SORT_NEWEST = iota
	SORT_OLDEST
	SORT_PINNED_FIRST
	SORT_PINNED_LAST
)
