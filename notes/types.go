package notes

type PageViewStatus struct {
	NumberOfPages int
	CurrentPage   int
	Step          int
}

type SearchFilter struct {
	Pinned     bool
	WholeWords bool
}
