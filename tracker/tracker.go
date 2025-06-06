package tracker

//Conveniant way to manage of list of uints

import (
	"slices"
)

var trackedItems []uint // keep track of items in a list

// returns true if not added to list, Returns false if note id was already iin the list
func AddToTracker(noteId uint) bool {
	var result bool = true
	if slices.Contains(trackedItems, noteId) {
		result = false
	} else {
		trackedItems = append(trackedItems, noteId)
	}
	return result
}

func DelFromTracker(noteId uint) bool {
	var result = false
	if index := slices.Index(trackedItems, noteId); index != -1 {
		trackedItems = slices.Delete(trackedItems, index, index+1)
		result = true
	}

	return result
}

func TrackerCheck(noteId uint) bool {
	return slices.Contains(trackedItems, noteId)
}

func TrackerLen() int {
	return len(trackedItems)
}
