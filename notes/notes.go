package notes

import (
	"log"

	"github.com/marcs100/minote/minotedb"
	"github.com/marcs100/minote/note"
)

func GetAllTags() []string {
	tags, err := minotedb.GetAllTags()
	if err != nil {
		log.Printf("Error tags - &s\n", err)
	}
	return tags
}

func GetPinnedNotes() ([]note.NoteData, error) {
	var pinnedNotes []note.NoteData = nil
	notesDB, err := minotedb.GetPinnedNotes()
	if err == nil {
		for _, noteDB := range notesDB {
			pinnedNotes = append(pinnedNotes, note.NoteData(noteDB))
		}
	}
	return pinnedNotes, err
}

func GetRecentNotes(recentNotesLimit int) ([]note.NoteData, error) {
	var recentNotes []note.NoteData = nil
	recentNotesDB, err := minotedb.GetRecentNotes(recentNotesLimit)
	if err == nil {
		for _, noteDB := range recentNotesDB {
			recentNotes = append(recentNotes, note.NoteData(noteDB))
		}
	}
	return recentNotes, err
}

func GetNotebook(name string) ([]note.NoteData, error) {
	var notebook_notes []note.NoteData = nil
	notebook_notesDB, err := minotedb.GetNotebook(name)
	if err == nil {
		for _, noteDB := range notebook_notesDB {
			notebook_notes = append(notebook_notes, note.NoteData(noteDB))
		}
	}
	return notebook_notes, err
}

func GetTaggedNotes(tags []string) ([]note.NoteData, error) {
	var taggedNotes []note.NoteData = nil
	taggedNotesDB, err := minotedb.GetTaggedNotes(tags)
	if err == nil {
		for _, noteDB := range taggedNotesDB {
			taggedNotes = append(taggedNotes, note.NoteData(noteDB))
		}
	}
	return taggedNotes, err
}

func GetSearchResults(searchText string, filter SearchFilter) ([]note.NoteData, error) {
	var searchNotes []note.NoteData = nil
	searchNotesDB, err := minotedb.GetSearchResults(searchText, minotedb.SearchFilter(filter))
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
