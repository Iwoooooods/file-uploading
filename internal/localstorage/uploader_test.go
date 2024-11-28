package localstorage

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCheckFileExists(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	// Create metadata table
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

	uploader, err := NewUploader(t.TempDir(), "testuser", db)
	if err != nil {
		t.Fatalf("failed to create uploader: %v", err)
	}

	t.Run("file exists", func(t *testing.T) {
		exists, err := uploader.CheckFileExists(context.Background(), "testmd5")
		if err != nil {
			t.Fatalf("failed to check file exists: %v", err)
		}
		if exists {
			t.Fatal("file should exist")
		}
	})

	t.Run("upload file", func(t *testing.T) {
		err := uploader.UploadFile(context.Background(), bytes.NewReader([]byte("hello world!")), "testfile")
		if err != nil {
			t.Fatalf("failed to upload file: %v", err)
		}
		exists, err := uploader.CheckFileExists(context.Background(), "testmd5")
		if err != nil {
			t.Fatalf("failed to check file exists: %v", err)
		}
		if !exists {
			t.Fatal("file should exist")
		}
	})
}
