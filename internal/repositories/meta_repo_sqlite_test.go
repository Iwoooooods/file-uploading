package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Iwoooooods/fs-upload-go/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func TestMetaRepositorySQLite_Create(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	//create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS metadata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_id TEXT NOT NULL,
			file_name TEXT NOT NULL,
			md5_hash TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	repo := NewMetaRepositorySQLite(db)
	t.Run("create", func(t *testing.T) {
		err := repo.Create(context.Background(), models.FileMetadata{
			FileId:   "1",
			FileName: "test.txt",
			MD5Hash:  "1234567890",
		})
		if err != nil {
			t.Fatalf("failed to create metadata: %v", err)
		}
		meta, err := repo.Get(context.Background(), "file_id", "1")
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.FileId != "1" {
			t.Fatalf("expected file id to be 1, got %s", meta.FileId)
		}
		t.Log(meta)
	})

	t.Run("update", func(t *testing.T) {
		err := repo.Update(context.Background(), models.FileMetadata{
			FileId:   "1",
			FileName: "test2.txt",
			MD5Hash:  "1234567890",
		})
		if err != nil {
			t.Fatalf("failed to update metadata: %v", err)
		}
		meta, err := repo.Get(context.Background(), "file_id", "1")
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.FileName != "test2.txt" {
			t.Fatalf("expected file name to be test2.txt, got %s", meta.FileName)
		}

	})

	t.Run("delete", func(t *testing.T) {
		err := repo.Delete(context.Background(), "1")
		if err != nil {
			t.Fatalf("failed to delete metadata: %v", err)
		}
		_, err = repo.Get(context.Background(), "file_id", "1")
		if err != sql.ErrNoRows {
			t.Fatalf("expected no rows, got %v", err)
		}
	})
}
