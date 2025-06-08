package ui

import "fmt"

func (pv *PageViewStatus) PageForward() int {
	if pv.CurrentPage+pv.Step <= pv.NumberOfPages {
		pv.CurrentPage += pv.Step
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageViewStatus) PageBack() int {
	if pv.CurrentPage > 1 {
		pv.CurrentPage -= pv.Step
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageViewStatus) Reset() {
	pv.CurrentPage = 0
	pv.NumberOfPages = 0
	pv.Step = 1
}

func (pv *PageViewStatus) GetLabelText() string {
	return fmt.Sprintf("Page: %d of %d", pv.CurrentPage, pv.NumberOfPages)
}

func (pv *PageViewStatus) GetGridLabelText() string {
	var pageRange int = 0
	if ((pv.CurrentPage - 1) + pv.Step) > pv.NumberOfPages {
		pageRange = pv.NumberOfPages
	} else {
		pageRange = (pv.CurrentPage - 1) + pv.Step
	}
	return fmt.Sprintf("Showing: %d to %d of %d", pv.CurrentPage, pageRange, pv.NumberOfPages)
}
