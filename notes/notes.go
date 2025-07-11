package notes

import (
	"log"

	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/note"
)

func GetAllTags() []string {
	tags, err := minotedb.GetAllTags()
	if err != nil {
		log.Printf("Error tags - %s\n", err)
	}
	return tags
}

func GetTagsWithSearch(searchQuery string) []string {
	searchQuery = "'%" + searchQuery + "%'"
	tags, err := minotedb.GetTags(searchQuery)
	if err != nil {
		log.Printf("Error tags - %s\n", err)
	}
	return tags
}

func GetPinnedNotes(sortBy int) ([]note.NoteData, error) {
	var pinnedNotes []note.NoteData = nil
	notesDB, err := minotedb.GetPinnedNotes(sortBy)
	if err == nil {
		for _, noteDB := range notesDB {
			pinnedNotes = append(pinnedNotes, note.NoteData(noteDB))
		}
	}
	return pinnedNotes, err
}

func GetRecentNotes(recentNotesLimit int, sortBy int) ([]note.NoteData, error) {
	var recentNotes []note.NoteData = nil
	recentNotesDB, err := minotedb.GetRecentNotes(recentNotesLimit, sortBy)
	if err == nil {
		for _, noteDB := range recentNotesDB {
			recentNotes = append(recentNotes, note.NoteData(noteDB))
		}
	}
	return recentNotes, err
}

func GetNotebook(name string, sortBy int) ([]note.NoteData, error) {
	var notebook_notes []note.NoteData = nil
	notebook_notesDB, err := minotedb.GetNotebook(name, sortBy)
	if err == nil {
		for _, noteDB := range notebook_notesDB {
			notebook_notes = append(notebook_notes, note.NoteData(noteDB))
		}
	}
	return notebook_notes, err
}

func GetTaggedNotes(tags []string, sortBy int) ([]note.NoteData, error) {
	var taggedNotes []note.NoteData = nil
	taggedNotesDB, err := minotedb.GetTaggedNotes(tags, sortBy)
	if err == nil {
		for _, noteDB := range taggedNotesDB {
			taggedNotes = append(taggedNotes, note.NoteData(noteDB))
		}
	}
	return taggedNotes, err
}

func GetSearchResults(searchText string, filter SearchFilter, sortBy int) ([]note.NoteData, error) {
	var searchNotes []note.NoteData = nil
	searchNotesDB, err := minotedb.GetSearchResults(searchText, minotedb.SearchFilter(filter), sortBy)
	if err == nil {
		for _, noteDB := range searchNotesDB {
			searchNotes = append(searchNotes, note.NoteData(noteDB))
		}
	}
	return searchNotes, err
}

func GetNote(id uint) (note.NoteData, error) {
	res, err := minotedb.GetNote(id)
	return note.NoteData(res), err
}

func GetNotebooks() ([]string, error) {
	return minotedb.GetNotebooks()
}

func CheckNotebookExists(notebook string) (bool, error) {
	return minotedb.CheckNotebookExists(notebook)
}

func DeleteNote(id uint) (bool, error) {
	var noteDeleted = false
	res, err := minotedb.DeleteNote(id)
	if res > 0 && err == nil {
		noteDeleted = true
	}

	return noteDeleted, err
}

func UnpinNote(id uint) (bool, error) {
	var unpinned = false
	res, err := minotedb.UnpinNote(id)
	if res > 0 && err == nil {
		unpinned = true
	}

	return unpinned, err
}

func PinNote(id uint) (bool, error) {
	var pinned = false
	res, err := minotedb.PinNote(id)
	if res > 0 && err == nil {
		pinned = true
	}

	return pinned, err
}

func GetPinnedDate(id uint) (string, error) {
	return minotedb.GetPinnedDate(id)
}
