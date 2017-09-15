package main

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE tomatoes (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	tag varchar,
	created_at datetime NOT NULL
);
`

type Tomato struct {
	Id        int       `db:"id"`
	Tag       string    `db:"tag"`
	CreatedAt time.Time `db:"created_at"`
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// TODO: file name get from config
const DB = "goma.db"

func initDB() error {
	if isExist(DB) {
		return nil
	}

	db, err := sqlx.Connect("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	db.MustExec(schema)

	return nil
}

func createTomato(tag string) error {
	db, err := sqlx.Connect("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tomatoes(tag, created_at) VALUES ($1, $2)", tag, time.Now())
	tx.Commit()

	return nil
}

func selectTomatos(start time.Time, end time.Time) ([]Tomato, error) {
	db, err := sqlx.Connect("sqlite3", DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tomatoes := []Tomato{}
	err = db.Select(&tomatoes, "SELECT id, tag, created_at FROM tomatoes WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at ", start, end)
	if err != nil {
		return nil, err
	}

	return tomatoes, nil
}
