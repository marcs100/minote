package minotedb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//*********** Public functions ************************

func GetNote(id uint) (NoteDataDB, error) {
	var query string = fmt.Sprintf("select * from notes where id = %d", id)
	notes, err := getNotes(query)
	return notes[0], err
}

func GetPinnedNotes() ([]NoteDataDB, error) {
	var query string = "select * from notes where pinned = 1 order by pinnedDate desc"
	return getNotes(query)
}

func GetPinnedDate(noteId uint) (string, error) {
	var query string = fmt.Sprintf("select pinnedDate from notes where id = %d", noteId)
	fields, err := getColumn(query)
	if err == nil && len(fields) == 0 {
		err = errors.New("No value found in db")
		return "", err
	}

	return fields[0], err
}

func GetNotebook(notebookName string) ([]NoteDataDB, error) {
	var query string = fmt.Sprintf("select * from notes where notebook = '%s' order by modified desc", notebookName)
	return getNotes(query)
}

func GetNotebooks() ([]string, error) {
	var query string = fmt.Sprintf("select distinct notebook from notes order by notebook asc")
	return getColumn(query)
}

func CheckNotebookExists(notebook string) (bool, error) {
	var query string = fmt.Sprintf("select notebook from notes where notebook = '%s'", notebook)
	notebooks, err := getColumn(query)
	if err != nil {
		return true, err
	}

	if len(notebooks) > 0 {
		return true, err
	}

	return false, err
}

func GetRecentNotes(noteCount int) ([]NoteDataDB, error) {
	var query string = fmt.Sprintf("select * from notes order by modified desc LIMIT %d", noteCount)
	return getNotes(query)
}

func GetSearchResults(searchQuery string, filter SearchFilter) ([]NoteDataDB, error) {
	return getSearchResults(searchQuery, filter)
}

func GetAllTags() ([]string, error) {
	var query string = "select distinct tag from tags order by tag ASC"
	return getColumn(query)
}

func GetTagsForNote(noteId uint) ([]string, error) {
	query := fmt.Sprintf("select tag from tags where noteId = %d", noteId)
	tags, err := getColumn(query)
	return tags, err
}

func GetTaggedNotes(taggedNotes []string) ([]NoteDataDB, error) {
	var err error = nil
	var notes []NoteDataDB
	var ids []string
	var taggedNotesMod []string

	//Add single quotes to string in each slice
	for _, s := range taggedNotes {
		taggedNotesMod = append(taggedNotesMod, fmt.Sprint("'", s, "'"))
	}

	//create csv list
	taggedCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(taggedNotesMod)), ", "), "[]")

	//Need to get the note id's of all the tagged notes from tags table
	query := fmt.Sprintf("select noteId from tags where tag in (%s)", taggedCSV)
	ids, err = getColumn(query)
	if err != nil {
		return nil, err
	}

	// Now get all the notes that conatain the retreived note id's
	idsCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")
	query = fmt.Sprintf("select * from notes where id in (%s) order by modified desc", idsCSV)
	notes, err = getNotes(query)

	return notes, err
}

func CheckTagExistsForNote(tag string, noteId uint) (bool, error) {
	var result bool = false
	query := fmt.Sprintf("select id from tags where noteId = %d and tag = '%s'", noteId, tag)
	tagsData, err := getColumn(query)
	if err != nil {
		return result, err
	}

	if len(tagsData) > 0 {
		result = true
	}

	return result, err
}

func DoesTagExist(tag string) (bool, error) {
	var result bool = false
	query := fmt.Sprintf("select tag from tags where tag = '%s'", tag)
	tagField, err := getColumn(query)
	if err != nil {
		return result, err
	}

	if len(tagField) > 0 {
		result = true
	}

	return result, err
}

//************ Private functions ************************

// use to get a single column as list of strings
func getColumn(query string) ([]string, error) {
	if !connected {
		return nil, errors.New("GetColumn: database not connected")
	}

	rows, err := db.Query(query)
	var fields []string
	defer rows.Close()
	for rows.Next() {
		var field string
		err := rows.Scan(&field)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, err
}

func getNotes(query string) ([]NoteDataDB, error) {
	if !connected {
		return nil, errors.New("Get Notes: database not connected")
	}

	rows, err := db.Query(query)
	var notes []NoteDataDB
	defer rows.Close()
	for rows.Next() {
		var note NoteDataDB
		err := rows.Scan(&note.Id, &note.Notebook, &note.Content, &note.Created, &note.Modified, &note.Pinned, &note.PinnedDate, &note.BackgroundColour)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	//fmt.Println(notes[1].content)
	//fmt.Println(notes[2].content)
	//fmt.Println(notes[3].content)

	return notes, err
}

func getSearchResults(searchText string, filter SearchFilter) ([]NoteDataDB, error) {
	var st1, st2, st3 string

	if !connected {
		return nil, errors.New("Get Notes: database not connected")
	}

	query := "select * from notes where content like ? "
	searchTerm := "%" + searchText + "%"
	if filter.WholeWords {
		query = "select * from notes where (content like ? or content like ? or content like ?)"
		st1 = fmt.Sprint(searchText, " %")
		st2 = fmt.Sprint("% ", searchText, " %")
		st3 = fmt.Sprint("% ", searchText)
		/*fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(st3)*/
	}

	if filter.Pinned {
		query = fmt.Sprintf("%s and pinned = 1", query)
	}

	query = fmt.Sprintf("%s order by modified desc", query)
	fmt.Println(query)

	var rows *sql.Rows
	var err error
	if filter.WholeWords == false {
		rows, err = db.Query(query, searchTerm)
	} else {
		dbprep, err2 := db.Prepare(query)
		if err2 != nil {
			log.Panicln(err2)
		}
		rows, err = dbprep.Query(st1, st2, st3)
	}

	var notes []NoteDataDB
	defer rows.Close()
	for rows.Next() {
		var note NoteDataDB
		err := rows.Scan(&note.Id, &note.Notebook, &note.Content, &note.Created, &note.Modified, &note.Pinned, &note.PinnedDate, &note.BackgroundColour)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}
	return notes, err
}
