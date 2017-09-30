package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	schema = `
CREATE TABLE tomatoes (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	tag varchar,
	created_at datetime NOT NULL
);
`

	selectQuery = `
SELECT id, tag, created_at FROM tomatoes WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at
`

	tagSummaryQuery = `
SELECT COUNT(tag) as tag_count, tag FROM tomatoes WHERE created_at BETWEEN $1 AND $2
	GROUP BY tag ORDER BY tag_count DESC
`
)

// Tomato is type for `tomatoes` table
type Tomato struct {
	ID        int       `db:"id"`
	Tag       string    `db:"tag"`
	CreatedAt time.Time `db:"created_at"`
}

// TagSummary is type for count per tag
type TagSummary struct {
	Count int    `db:"tag_count"`
	Tag   string `db:"tag"`
}

func initDB() error {
	if isExist(cfg.DataBase) {
		return nil
	}

	db, err := sqlx.Connect("sqlite3", cfg.DataBase)
	if err != nil {
		return err
	}
	defer db.Close()

	db.MustExec(schema)

	return nil
}

func createTomato(tag string) error {
	db, err := sqlx.Connect("sqlite3", cfg.DataBase)
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
	db, err := sqlx.Connect("sqlite3", cfg.DataBase)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tomatoes := []Tomato{}
	err = db.Select(&tomatoes, selectQuery, start, end)
	if err != nil {
		return nil, err
	}

	return tomatoes, nil
}

func selectTagSummary(start time.Time, end time.Time) ([]TagSummary, error) {
	db, err := sqlx.Connect("sqlite3", cfg.DataBase)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tagSummaries := []TagSummary{}
	err = db.Select(&tagSummaries, tagSummaryQuery, start, end)
	if err != nil {
		return nil, err
	}

	return tagSummaries, nil
}
