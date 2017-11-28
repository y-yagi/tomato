package tomato

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/y-yagi/goext/osext"
)

func TestInitDB(t *testing.T) {
	pwd, _ := os.Getwd()
	db := filepath.Join(pwd, "testdata", "tomato.db")

	if osext.IsExist(db) {
		t.Errorf("Expect database does not exist. %s", db)
	}

	repo := NewRepository(db)

	if err := repo.InitDB(); err != nil {
		t.Errorf("InitDB failed. %v", err)
	}

	defer os.Remove(db)

	if !osext.IsExist(db) {
		t.Errorf("Expect database exists. %s", db)
	}
}
