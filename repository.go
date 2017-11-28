package tomato

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/y-yagi/goext/osext"
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

// TagSummary is type for count per tag.
type TagSummary struct {
	Count int    `db:"tag_count"`
	Tag   string `db:"tag"`
}

// Repository is type for database operation.
type Repository struct {
	database string
}

// NewRepository creates a new repository.
func NewRepository(database string) *Repository {
	return &Repository{database: database}
}

// InitDB initialize database.
func (r *Repository) InitDB() error {
	if osext.IsExist(r.database) {
		return nil
	}

	db, err := sqlx.Connect("sqlite3", r.database)
	if err != nil {
		return err
	}
	defer db.Close()

	db.MustExec(schema)

	return nil
}

func (r *Repository) createTomato(tag string) error {
	db, err := sqlx.Connect("sqlite3", r.database)
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tomatoes(tag, created_at) VALUES ($1, $2)", tag, time.Now())
	tx.Commit()

	return nil
}

func (r *Repository) selectTomatos(start time.Time, end time.Time) ([]Tomato, error) {
	db, err := sqlx.Connect("sqlite3", r.database)
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

func (r *Repository) selectTagSummary(start time.Time, end time.Time) ([]TagSummary, error) {
	db, err := sqlx.Connect("sqlite3", r.database)
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
