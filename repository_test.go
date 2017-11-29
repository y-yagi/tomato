package tomato

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/y-yagi/goext/osext"
)

func TestInitDB(t *testing.T) {
	pwd, _ := os.Getwd()
	db := filepath.Join(pwd, "testdata", "init.db")

	if osext.IsExist(db) {
		t.Fatalf("Expect database does not exist. %s", db)
	}

	repo := NewRepository(db)

	if err := repo.InitDB(); err != nil {
		t.Fatalf("InitDB failed. %v", err)
	}

	defer os.Remove(db)

	if !osext.IsExist(db) {
		t.Fatalf("Expect database exists. %s", db)
	}
}

func TestCreateAndSelectTomato(t *testing.T) {
	pwd, _ := os.Getwd()
	db := filepath.Join(pwd, "testdata", "tomato.db")
	r := NewRepository(db)
	r.InitDB()
	defer os.Remove(db)

	start := time.Now()
	r.createTomato("test")
	r.createTomato("dummy")
	r.createTomato("test")
	end := time.Now()

	tomatoes, err := r.selectTomatos(start, end)
	if err != nil {
		t.Fatalf("selectTomatos failed. %v", err)
	}

	if len(tomatoes) != 3 {
		t.Fatalf("Expected result count %d got %d", 3, len(tomatoes))
	}

	expectedTags := []string{"test", "dummy", "test"}
	for index, tomato := range tomatoes {
		if tomato.Tag != expectedTags[index] {
			t.Fatalf("Expected %s got %s", expectedTags[index], tomato.Tag)
		}

		id := index + 1
		if tomato.ID != id {
			t.Fatalf("Expected %d got %d", id, tomato.ID)
		}
	}

	summaries, err := r.selectTagSummary(start, end)
	if err != nil {
		t.Fatalf("selectTagSummary failed. %v", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("Expected result count %d got %d", 2, len(summaries))
	}

	cases := []struct {
		wantTag   string
		wantCount int
	}{
		{
			wantTag:   "test",
			wantCount: 2,
		},
		{
			wantTag:   "dummy",
			wantCount: 1,
		},
	}

	for index, summary := range summaries {
		if summary.Tag != cases[index].wantTag {
			t.Fatalf("Expected %s got %s", cases[index].wantTag, summary.Tag)
		}

		if summary.Count != cases[index].wantCount {
			t.Fatalf("Expected %d got %d", cases[index].wantCount, summary.Count)
		}
	}
}
