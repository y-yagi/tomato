package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func initDB() (*sql.DB, error) {
	// TODO: file name get from config
	const DB = "goma.db"
	const createTable = `
	CREATE TABLE "tomatoes" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, "tag" varchar, "created_at" datetime NOT NULL);
	`

	dbExist := false

	if isExist(DB) {
		dbExist = true
	}

	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if !dbExist {
		_, err = db.Exec(createTable)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
