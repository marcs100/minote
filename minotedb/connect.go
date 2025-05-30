package minotedb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var connected = false
var db *sql.DB = nil

func Open(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Panicln(err)
		connected = false
	} else {
		log.Println("data base open success")
		connected = true
	}

	return err
}

func Close() {
	if connected {
		db.Close()
	}
}
