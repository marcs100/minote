package minotedb

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func PinNote(id uint) (int64, error) {
	if !connected {
		return 0, errors.New("Pin Note: database not connected")
	}
	var pinDate = time.Now().String()[0:19]
	res, err := db.Exec("update notes set pinned = 1, pinnedDate = ? where id = ?", pinDate, id)
	rows, _ := res.RowsAffected()
	return rows, err
}

func UnpinNote(id uint) (int64, error) {
	if !connected {
		return 0, errors.New("Unpin note: database not connected")
	}
	res, err := db.Exec("update notes set pinned = 0, pinnedDate = '' where id = ?", id)
	rows, _ := res.RowsAffected()
	return rows, err
}

func DeleteNote(id uint) (int64, error) {

	if !connected {
		return 0, errors.New("Delete note: database not connected")
	}
	res, err := db.Exec("delete from notes where id = ?", id)
	rows, _ := res.RowsAffected()

	//also need to remove id from tags table!!!!!!!!!!!!!!
	if err == nil {
		res, err = db.Exec("delete from tags where noteId = ?", id)
	}

	return rows, err
}

func SaveNote(id uint, notebook string, content string, pinned uint, pinnedDate string, colour string) (int64, error) {
	if !connected {
		return 0, errors.New("SaveNote: database not connected")
	}

	var modified = time.Now().String()[0:19]
	//fmt.Printf("fields: %d, %s, %d, %s, %s\n", id, notebook, pinned, modified, colour)

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?, pinnedDate = ?,  modified = ?, BGColour = ? where id = ?", notebook, content, pinned, pinnedDate, modified, colour, id)
	rows, _ := res.RowsAffected()

	return rows, err
}

func SaveNoteNoTimeStamp(id uint, notebook string, content string, pinned uint, pinnedDate, colour string) (int64, error) {
	if !connected {
		return 0, errors.New("SaveNote: database not connected")
	}

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?, pinnedDate = ?, BGColour = ? where id = ?", notebook, content, pinned, pinnedDate, colour, id)
	rows, _ := res.RowsAffected()

	return rows, err
}

func InsertNote(notebook string, content string, pinned uint, pinnedDate string, colour string) (int64, int64, error) {
	var id int64 = -1
	if !connected {
		return 0, 0, errors.New("InsertNote: database not connected")
	}

	//calculate date created and modified
	var created = time.Now().String()[0:19]
	var modified = created
	if pinned > 0 {
		pinnedDate = time.Now().String()[0:19]
	}

	res, err := db.Exec("INSERT INTO notes VALUES(NULL,?,?,?,?,?,?,?)", notebook, content, created, modified, pinned, pinnedDate, colour)
	rows, _ := res.RowsAffected()
	id, _ = res.LastInsertId()

	return rows, id, err
}

func WriteTags(tags []string, noteId uint) (uint, error) {
	var err error = nil
	var tagsWritten uint = 0
	if !connected {
		return 0, errors.New("WriteTags: database not connected")
	}

	for _, tag := range tags {
		//does tag already exist for this noteId
		tagExists, terr := CheckTagExistsForNote(tag, noteId)
		if terr != nil {
			err = terr
			break
		}

		if !tagExists {
			//only write if the tag does not already exist for thsi note
			res, werr := db.Exec("INSERT INTO tags VALUES(NULL,?,?)", tag, noteId)
			if werr != nil {
				err = werr
				break
			}

			if r, ierr := res.RowsAffected(); ierr == nil {
				if r > 0 {
					tagsWritten += 1
				}
			} else {
				err = ierr
			}
		} else {
			err = errors.New(fmt.Sprintf("Tag %s already exists for note %d", tag, noteId))
		}
	}
	return tagsWritten, err
}

func DeleteTags(tags []string, noteId uint) (uint, error) {
	var err error = nil
	var tagsDeleted uint = 0
	if !connected {
		return 0, errors.New("DeleteTags: database not connected")
	}

	for _, tag := range tags {
		//does tag already exist for this noteId
		tagExists, terr := CheckTagExistsForNote(tag, noteId)
		if terr != nil {
			err = terr
			break
		}

		if tagExists {
			//only delete if the tag already exists for this note
			res, werr := db.Exec("DELETE from tags where tag = ? and noteId = ?", tag, noteId)
			if werr != nil {
				err = werr
				break
			}

			if r, ierr := res.RowsAffected(); ierr == nil {
				if r > 0 {
					tagsDeleted += 1
				}
			} else {
				err = ierr
			}
		} else {
			err = errors.New(fmt.Sprintf("Tag %s does not exist, nothing to delete for note %d", tag, noteId))
		}

	}
	return tagsDeleted, err
}
